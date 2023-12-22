/*
 * Copyright (C) 2023 The Android Open Source Project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

use anyhow::{bail, ensure, Context, Result};
use protobuf::Message;
use std::collections::HashSet;
use std::io::Read;
use std::path::PathBuf;

use crate::codegen::cpp::generate_cpp_code;
use crate::codegen::java::generate_java_code;
use crate::codegen::rust::generate_rust_code;
use crate::codegen::CodegenMode;
use crate::dump::{DumpFormat, DumpPredicate};
use crate::protos::{
    CachedFlagExt, ParsedFlagExt, ProtoCache, ProtoCachedFlag, ProtoFlagMetadata,
    ProtoFlagPermission, ProtoFlagState, ProtoParsedFlag, ProtoParsedFlags, ProtoTracepoint,
};
use crate::storage::generate_storage_files;

pub struct Input {
    pub source: String,
    pub reader: Box<dyn Read>,
}

impl Input {
    fn try_parse_cache(&mut self) -> Result<ProtoCache> {
        let mut buffer = Vec::new();
        self.reader
            .read_to_end(&mut buffer)
            .with_context(|| format!("failed to read {}", self.source))?;
        crate::protos::cache::try_from_binary_proto(&buffer).with_context(|| self.error_context())
    }

    fn error_context(&self) -> String {
        format!("failed to parse {}", self.source)
    }
}

pub struct OutputFile {
    pub path: PathBuf, // relative to some root directory only main knows about
    pub contents: Vec<u8>,
}

pub const DEFAULT_FLAG_STATE: ProtoFlagState = ProtoFlagState::DISABLED;
pub const DEFAULT_FLAG_PERMISSION: ProtoFlagPermission = ProtoFlagPermission::READ_WRITE;

pub fn create_cache(
    package: &str,
    container: Option<&str>,
    declarations: Vec<Input>,
    values: Vec<Input>,
    default_permission: ProtoFlagPermission,
) -> Result<Vec<u8>> {
    let mut cache = ProtoCache {
        package: Some(package.to_string()),
        container: container.map(|s| s.to_string()),
        ..Default::default()
    };

    for mut input in declarations {
        let mut contents = String::new();
        input
            .reader
            .read_to_string(&mut contents)
            .with_context(|| format!("failed to read {}", input.source))?;

        let flag_declarations = crate::protos::flag_declarations::try_from_text_proto(&contents)
            .with_context(|| input.error_context())?;
        ensure!(
            cache.package() == flag_declarations.package(),
            "failed to parse {}: expected package {}, got {}",
            input.source,
            cache.package(),
            flag_declarations.package()
        );
        if let Some(c) = container {
            ensure!(
                c == flag_declarations.container(),
                "failed to parse {}: expected container {}, got {}",
                input.source,
                c,
                flag_declarations.container()
            );
        }
        for flag_declaration in flag_declarations.flag.into_iter() {
            crate::protos::flag_declaration::verify_fields(&flag_declaration)
                .with_context(|| input.error_context())?;

            // create ProtoCachedFlag using ProtoFlagDeclaration and default values
            let permission = if flag_declaration.is_fixed_read_only() {
                Some(ProtoFlagPermission::READ_ONLY.into())
            } else {
                Some(default_permission.into())
            };
            let cached_flag = ProtoCachedFlag {
                name: flag_declaration.name,
                namespace: flag_declaration.namespace,
                description: flag_declaration.description,
                bug: flag_declaration.bug,
                state: Some(DEFAULT_FLAG_STATE.into()),
                permission,
                trace: vec![ProtoTracepoint {
                    source: Some(input.source.clone()),
                    state: Some(DEFAULT_FLAG_STATE.into()),
                    permission,
                    ..Default::default()
                }],
                is_fixed_read_only: flag_declaration.is_fixed_read_only.or(Some(false)),
                is_exported: flag_declaration.is_exported.or(Some(false)),
                metadata: flag_declaration.metadata,
                ..Default::default()
            };

            // verify CachedFlag looks reasonable
            crate::protos::cached_flag::verify_fields(&cached_flag)?;

            // verify ProtoCachedFlag can be added
            ensure!(
                cache.cached_flag.iter().all(|other| other.name() != cached_flag.name()),
                "failed to declare flag {} from {}: flag already declared",
                cached_flag.name(),
                input.source
            );

            // add ProtoCachedFlag to cache
            cache.cached_flag.push(cached_flag);
        }
    }

    for mut input in values {
        let mut contents = String::new();
        input
            .reader
            .read_to_string(&mut contents)
            .with_context(|| format!("failed to read {}", input.source))?;
        let flag_values = crate::protos::flag_values::try_from_text_proto(&contents)
            .with_context(|| input.error_context())?;
        for flag_value in flag_values.flag_value.into_iter() {
            crate::protos::flag_value::verify_fields(&flag_value)
                .with_context(|| input.error_context())?;

            let Some(cached_flag) = cache
                .cached_flag
                .iter_mut()
                .find(|cf| package == flag_value.package() && cf.name() == flag_value.name())
            else {
                // (silently) skip unknown flags
                continue;
            };

            ensure!(
                !cached_flag.is_fixed_read_only()
                    || flag_value.permission() == ProtoFlagPermission::READ_ONLY,
                "failed to set permission of flag {}, since this flag is fixed read only flag",
                flag_value.name()
            );

            cached_flag.set_state(flag_value.state());
            cached_flag.set_permission(flag_value.permission());
            let tracepoint = ProtoTracepoint {
                source: Some(input.source.clone()),
                state: flag_value.state,
                permission: flag_value.permission,
                ..Default::default()
            };
            cached_flag.trace.push(tracepoint);
        }
    }

    // Create a sorted cache
    crate::protos::cache::sort_cached_flags(&mut cache);
    crate::protos::cache::verify_fields(&cache)?;
    let mut output = Vec::new();
    cache.write_to_vec(&mut output)?;
    Ok(output)
}

pub fn create_java_lib(mut input: Input, codegen_mode: CodegenMode) -> Result<Vec<OutputFile>> {
    let cache = input.try_parse_cache()?;
    let package = cache.package().to_string();
    let modified_cached_flags =
        modify_cached_flags_based_on_mode(cache.cached_flag.into_iter(), codegen_mode)?;
    generate_java_code(&package, modified_cached_flags.into_iter(), codegen_mode)
}

pub fn create_cpp_lib(mut input: Input, codegen_mode: CodegenMode) -> Result<Vec<OutputFile>> {
    let cache = input.try_parse_cache()?;
    let package = cache.package().to_string();
    let modified_cached_flags =
        modify_cached_flags_based_on_mode(cache.cached_flag.into_iter(), codegen_mode)?;
    generate_cpp_code(&package, modified_cached_flags.into_iter(), codegen_mode)
}

pub fn create_rust_lib(mut input: Input, codegen_mode: CodegenMode) -> Result<OutputFile> {
    let cache = input.try_parse_cache()?;
    let package = cache.package().to_string();
    let modified_cached_flags =
        modify_cached_flags_based_on_mode(cache.cached_flag.into_iter(), codegen_mode)?;
    generate_rust_code(&package, modified_cached_flags.into_iter(), codegen_mode)
}

pub fn create_storage(inputs: Vec<Input>, container: &str) -> Result<Vec<OutputFile>> {
    let caches: Vec<ProtoCache> = inputs
        .into_iter()
        .map(|mut input| input.try_parse_cache())
        .collect::<Result<Vec<_>>>()?
        .into_iter()
        .filter(|cache| cache.container() == container)
        .collect();
    generate_storage_files(container, caches.iter())
}

pub fn create_device_config_defaults(mut input: Input) -> Result<Vec<u8>> {
    let cache = input.try_parse_cache()?;
    let package = cache.package().to_string();
    let mut output = Vec::new();
    for cached_flag in cache
        .cached_flag
        .into_iter()
        .filter(|pf| pf.permission() == ProtoFlagPermission::READ_WRITE)
    {
        let line = format!(
            "{}:{}={}\n",
            cached_flag.namespace(),
            cached_flag.fully_qualified_name(&package),
            match cached_flag.state() {
                ProtoFlagState::ENABLED => "enabled",
                ProtoFlagState::DISABLED => "disabled",
            }
        );
        output.extend_from_slice(line.as_bytes());
    }
    Ok(output)
}

pub fn create_device_config_sysprops(mut input: Input) -> Result<Vec<u8>> {
    let cache = input.try_parse_cache()?;
    let package = cache.package().to_string();
    let mut output = Vec::new();
    for cached_flag in cache
        .cached_flag
        .into_iter()
        .filter(|pf| pf.permission() == ProtoFlagPermission::READ_WRITE)
    {
        let line = format!(
            "persist.device_config.{}={}\n",
            cached_flag.fully_qualified_name(&package),
            match cached_flag.state() {
                ProtoFlagState::ENABLED => "true",
                ProtoFlagState::DISABLED => "false",
            }
        );
        output.extend_from_slice(line.as_bytes());
    }
    Ok(output)
}

pub fn dump_cache(mut input: Input, format: DumpFormat, filters: &[&str]) -> Result<Vec<u8>> {
    let cache = input.try_parse_cache()?;
    let filters: Vec<Box<DumpPredicate>> = if filters.is_empty() {
        vec![Box::new(|_| true)]
    } else {
        filters
            .iter()
            .map(|f| crate::dump::create_filter_predicate(f))
            .collect::<Result<Vec<_>>>()?
    };
    /*
     * FIXME: add this back
    crate::dump::dump_parsed_flags(
        cache.cached_flag.into_iter().filter(|flag| filters.iter().any(|p| p(flag))),
        format,
    )
    */
    todo!();
}

pub fn export_flags(mut input: Vec<Input>, dedup: bool) -> Result<Vec<u8>> {
    let mut parsed_flags = ProtoParsedFlags::new();
    let mut seen_flags = HashSet::new();
    let caches =
        input.iter_mut().map(|i: &mut Input| i.try_parse_cache()).collect::<Result<Vec<_>>>()?;
    for cache in caches.into_iter() {
        let package = cache.package().to_string();
        let container = cache.container().to_string();
        for cached_flag in cache.cached_flag.into_iter() {
            let qualified_name = cached_flag.fully_qualified_name(&package);
            if !seen_flags.insert(qualified_name.clone()) {
                if dedup {
                    continue;
                } else {
                    bail!("duplicate flag {}", qualified_name);
                }
            }
            let mut parsed_flag = ProtoParsedFlag {
                package: Some(package.to_string()),
                name: cached_flag.name,
                namespace: cached_flag.namespace,
                description: cached_flag.description,
                bug: cached_flag.bug,
                state: cached_flag.state,
                permission: cached_flag.permission,
                trace: cached_flag.trace,
                is_fixed_read_only: cached_flag.is_fixed_read_only,
                is_exported: cached_flag.is_exported,
                container: Some(container.to_string()),
                ..Default::default()
            };
            let mut metadata = ProtoFlagMetadata::new();
            metadata.set_purpose(cached_flag.metadata.purpose());
            parsed_flag.metadata = Some(metadata).into();
            parsed_flags.parsed_flag.push(parsed_flag);
        }
    }
    parsed_flags.parsed_flag.sort_by_cached_key(|flag| flag.fully_qualified_name());
    crate::protos::parsed_flags::verify_fields(&parsed_flags)?;
    let mut output = Vec::new();
    parsed_flags.write_to_vec(&mut output)?;
    Ok(output)
}

pub fn dump_flags(mut input: Input) -> Result<Vec<u8>> {
    let parsed_flags: ProtoParsedFlags = protobuf::Message::parse_from_reader(&mut input.reader)
        .with_context(|| format!("failed to parse {}", input.source))?;
    let mut output = vec![];
    for flag in parsed_flags.parsed_flag.into_iter() {
        let str = format!("{}\n", flag.fully_qualified_name());
        output.extend_from_slice(str.as_bytes());
    }
    Ok(output)
}

pub fn modify_cached_flags_based_on_mode<I>(
    iter: I,
    codegen_mode: CodegenMode,
) -> Result<Vec<ProtoCachedFlag>>
where
    I: Iterator<Item = ProtoCachedFlag>,
{
    fn exported_mode_flag_modifier(mut parsed_flag: ProtoCachedFlag) -> ProtoCachedFlag {
        parsed_flag.set_state(ProtoFlagState::DISABLED);
        parsed_flag.set_permission(ProtoFlagPermission::READ_WRITE);
        parsed_flag.set_is_fixed_read_only(false);
        parsed_flag
    }

    let modified_parsed_flags: Vec<_> = match codegen_mode {
        CodegenMode::Exported => {
            iter.filter(|pf| pf.is_exported()).map(exported_mode_flag_modifier).collect()
        }
        CodegenMode::Production | CodegenMode::Test => iter.collect(),
    };
    if modified_parsed_flags.is_empty() {
        bail!("{codegen_mode} library contains no {codegen_mode} flags");
    }

    Ok(modified_parsed_flags)
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::protos::ProtoFlagPurpose;

    #[test]
    fn test_create_cache() {
        let cache = crate::test::create_test_cache(); // calls create_cache
        crate::protos::cache::verify_fields(&cache).unwrap();

        let enabled_ro = cache.cached_flag.iter().find(|flag| flag.name() == "enabled_ro").unwrap();
        assert!(crate::protos::cached_flag::verify_fields(enabled_ro).is_ok());
        assert_eq!("enabled_ro", enabled_ro.name());
        assert_eq!("This flag is ENABLED + READ_ONLY", enabled_ro.description());
        assert_eq!(ProtoFlagState::ENABLED, enabled_ro.state());
        assert_eq!(ProtoFlagPermission::READ_ONLY, enabled_ro.permission());
        assert_eq!(ProtoFlagPurpose::PURPOSE_BUGFIX, enabled_ro.metadata.purpose());
        assert_eq!(3, enabled_ro.trace.len());
        assert!(!enabled_ro.is_fixed_read_only());
        assert_eq!("tests/test.aconfig", enabled_ro.trace[0].source());
        assert_eq!(ProtoFlagState::DISABLED, enabled_ro.trace[0].state());
        assert_eq!(ProtoFlagPermission::READ_WRITE, enabled_ro.trace[0].permission());
        assert_eq!("tests/first.values", enabled_ro.trace[1].source());
        assert_eq!(ProtoFlagState::DISABLED, enabled_ro.trace[1].state());
        assert_eq!(ProtoFlagPermission::READ_WRITE, enabled_ro.trace[1].permission());
        assert_eq!("tests/second.values", enabled_ro.trace[2].source());
        assert_eq!(ProtoFlagState::ENABLED, enabled_ro.trace[2].state());
        assert_eq!(ProtoFlagPermission::READ_ONLY, enabled_ro.trace[2].permission());

        assert_eq!(9, cache.cached_flag.len());
        for flag in cache.cached_flag.iter() {
            if flag.name().starts_with("enabled_fixed_ro") {
                continue;
            }
            let first = flag.trace.first().unwrap();
            assert_eq!(DEFAULT_FLAG_STATE, first.state());
            assert_eq!(DEFAULT_FLAG_PERMISSION, first.permission());

            let last = flag.trace.last().unwrap();
            assert_eq!(flag.state(), last.state());
            assert_eq!(flag.permission(), last.permission());
        }

        let enabled_fixed_ro =
            cache.cached_flag.iter().find(|pf| pf.name() == "enabled_fixed_ro").unwrap();
        assert!(enabled_fixed_ro.is_fixed_read_only());
        assert_eq!(ProtoFlagState::ENABLED, enabled_fixed_ro.state());
        assert_eq!(ProtoFlagPermission::READ_ONLY, enabled_fixed_ro.permission());
        assert_eq!(2, enabled_fixed_ro.trace.len());
        assert_eq!(ProtoFlagPermission::READ_ONLY, enabled_fixed_ro.trace[0].permission());
        assert_eq!(ProtoFlagPermission::READ_ONLY, enabled_fixed_ro.trace[1].permission());
    }

    #[test]
    fn test_create_cache_setting_default() {
        let first_flag = r#"
        package: "com.first"
        flag {
            name: "first"
            namespace: "first_ns"
            description: "This is the description of the first flag."
            bug: "123"
        }
        "#;
        let declaration =
            vec![Input { source: "memory".to_string(), reader: Box::new(first_flag.as_bytes()) }];
        let value: Vec<Input> = vec![];

        let flags_bytes = crate::commands::create_cache(
            "com.first",
            None,
            declaration,
            value,
            ProtoFlagPermission::READ_ONLY,
        )
        .unwrap();
        let cache = crate::protos::cache::try_from_binary_proto(&flags_bytes).unwrap();
        assert_eq!(1, cache.cached_flag.len());
        let cached_flag = cache.cached_flag.first().unwrap();
        assert_eq!(ProtoFlagState::DISABLED, cached_flag.state());
        assert_eq!(ProtoFlagPermission::READ_ONLY, cached_flag.permission());
    }

    #[test]
    fn test_parse_flags_package_mismatch_between_declaration_and_command_line() {
        let first_flag = r#"
        package: "com.declaration.package"
        container: "first.container"
        flag {
            name: "first"
            namespace: "first_ns"
            description: "This is the description of the first flag."
            bug: "123"
        }
        "#;
        let declaration =
            vec![Input { source: "memory".to_string(), reader: Box::new(first_flag.as_bytes()) }];

        let value: Vec<Input> = vec![];

        let error = crate::commands::create_cache(
            "com.argument.package",
            Some("first.container"),
            declaration,
            value,
            ProtoFlagPermission::READ_WRITE,
        )
        .unwrap_err();
        assert_eq!(
            format!("{:?}", error),
            "failed to parse memory: expected package com.argument.package, got com.declaration.package"
        );
    }

    #[test]
    fn test_parse_flags_container_mismatch_between_declaration_and_command_line() {
        let first_flag = r#"
        package: "com.first"
        container: "declaration.container"
        flag {
            name: "first"
            namespace: "first_ns"
            description: "This is the description of the first flag."
            bug: "123"
        }
        "#;
        let declaration =
            vec![Input { source: "memory".to_string(), reader: Box::new(first_flag.as_bytes()) }];

        let value: Vec<Input> = vec![];

        let error = crate::commands::create_cache(
            "com.first",
            Some("argument.container"),
            declaration,
            value,
            ProtoFlagPermission::READ_WRITE,
        )
        .unwrap_err();
        assert_eq!(
            format!("{:?}", error),
            "failed to parse memory: expected container argument.container, got declaration.container"
        );
    }

    #[test]
    fn test_parse_flags_override_fixed_read_only() {
        let first_flag = r#"
        package: "com.first"
        container: "com.first.container"
        flag {
            name: "first"
            namespace: "first_ns"
            description: "This is the description of the first flag."
            bug: "123"
            is_fixed_read_only: true
        }
        "#;
        let declaration =
            vec![Input { source: "memory".to_string(), reader: Box::new(first_flag.as_bytes()) }];

        let first_flag_value = r#"
        flag_value {
            package: "com.first"
            name: "first"
            state: DISABLED
            permission: READ_WRITE
        }
        "#;
        let value = vec![Input {
            source: "memory".to_string(),
            reader: Box::new(first_flag_value.as_bytes()),
        }];
        let error = crate::commands::create_cache(
            "com.first",
            Some("com.first.container"),
            declaration,
            value,
            ProtoFlagPermission::READ_WRITE,
        )
        .unwrap_err();
        assert_eq!(
            format!("{:?}", error),
            "failed to set permission of flag first, since this flag is fixed read only flag"
        );
    }

    #[test]
    fn test_parse_flags_metadata() {
        let metadata_flag = r#"
        package: "com.first"
        flag {
            name: "first"
            namespace: "first_ns"
            description: "This is the description of this feature flag."
            bug: "123"
            metadata {
                purpose: PURPOSE_FEATURE
            }
        }
        "#;
        let declaration = vec![Input {
            source: "memory".to_string(),
            reader: Box::new(metadata_flag.as_bytes()),
        }];
        let value: Vec<Input> = vec![];

        let flags_bytes = crate::commands::create_cache(
            "com.first",
            None,
            declaration,
            value,
            ProtoFlagPermission::READ_ONLY,
        )
        .unwrap();
        let cache = crate::protos::cache::try_from_binary_proto(&flags_bytes).unwrap();
        assert_eq!(1, cache.cached_flag.len());
        let cached_flag = cache.cached_flag.first().unwrap();
        assert_eq!(ProtoFlagPurpose::PURPOSE_FEATURE, cached_flag.metadata.purpose());
    }

    #[test]
    fn test_create_device_config_defaults() {
        let input = create_test_cache_as_input();
        let bytes = create_device_config_defaults(input).unwrap();
        let text = std::str::from_utf8(&bytes).unwrap();
        assert_eq!("aconfig_test:com.android.aconfig.test.disabled_rw=disabled\naconfig_test:com.android.aconfig.test.disabled_rw_exported=disabled\nother_namespace:com.android.aconfig.test.disabled_rw_in_other_namespace=disabled\naconfig_test:com.android.aconfig.test.enabled_rw=enabled\n", text);
    }

    #[test]
    fn test_create_device_config_sysprops() {
        let input = create_test_cache_as_input();
        let bytes = create_device_config_sysprops(input).unwrap();
        let text = std::str::from_utf8(&bytes).unwrap();
        assert_eq!("persist.device_config.com.android.aconfig.test.disabled_rw=false\npersist.device_config.com.android.aconfig.test.disabled_rw_exported=false\npersist.device_config.com.android.aconfig.test.disabled_rw_in_other_namespace=false\npersist.device_config.com.android.aconfig.test.enabled_rw=true\n", text);
    }

    #[test]
    fn test_dump() {
        let input = create_test_cache_as_input();
        let bytes =
            dump_cache(input, DumpFormat::Custom("{fully_qualified_name}".to_string()), &[])
                .unwrap();
        let text = std::str::from_utf8(&bytes).unwrap();
        assert!(text.contains("com.android.aconfig.test.disabled_ro"));
    }

    #[test]
    fn test_export_flags() {
        let input = create_test_cache_as_input();
        let bytes = export_flags(vec![input], false).unwrap();
        let expected = protobuf::text_format::parse_from_str::<ProtoParsedFlags>(
            crate::test::TEST_FLAGS_TEXTPROTO,
        )
        .unwrap()
        .write_to_bytes()
        .unwrap();
        assert_eq!(bytes, expected);
    }

    #[test]
    fn test_export_flags_dedup() {
        let input = create_test_cache_as_input();
        let input_copy = create_test_cache_as_input();
        let bytes = export_flags(vec![input, input_copy], true).unwrap();
        let expected = protobuf::text_format::parse_from_str::<ProtoParsedFlags>(
            crate::test::TEST_FLAGS_TEXTPROTO,
        )
        .unwrap()
        .write_to_bytes()
        .unwrap();
        assert_eq!(bytes, expected);
    }

    #[test]
    fn test_export_flags_multiple_packages() {
        fn create_test_cache_as_input_2() -> Input {
            let textproto = r#"
package: "zzz.zzz"
container: "zzz"
cached_flag {
    name: "zzz"
    namespace: "zzz"
    bug: "zzz"
    description: "zzz"
    state: ENABLED
    permission: READ_WRITE
    trace {
        source: "memory"
        state: DISABLED
        permission: READ_ONLY
    }
    is_fixed_read_only: false
    is_exported: false
}
"#;
            let cache = crate::protos::cache::try_from_text_proto(textproto).unwrap();
            let binary_proto = cache.write_to_bytes().unwrap();
            let cursor = std::io::Cursor::new(binary_proto);
            let reader = Box::new(cursor);
            Input { source: "test.data".to_string(), reader }
        }

        // input order A, B
        let bytes1 =
            export_flags(vec![create_test_cache_as_input(), create_test_cache_as_input_2()], false)
                .unwrap();

        // input order B, A
        let bytes2 =
            export_flags(vec![create_test_cache_as_input_2(), create_test_cache_as_input()], false)
                .unwrap();
        assert_eq!(bytes1, bytes2);

        let parsed_flags: ProtoParsedFlags = protobuf::Message::parse_from_bytes(&bytes1).unwrap();
        assert_eq!(parsed_flags.parsed_flag[0].package(), "com.android.aconfig.test");
        assert_eq!(
            parsed_flags.parsed_flag[parsed_flags.parsed_flag.len() - 1].package(),
            "zzz.zzz"
        );
    }

    fn create_test_cache_as_input() -> Input {
        let cache = crate::test::create_test_cache();
        let binary_proto = cache.write_to_bytes().unwrap();
        let cursor = std::io::Cursor::new(binary_proto);
        let reader = Box::new(cursor);
        Input { source: "test.data".to_string(), reader }
    }

    #[test]
    fn test_modify_cached_flags_based_on_mode_prod() {
        let cache = crate::test::create_test_cache();
        let p_parsed_flags = modify_cached_flags_based_on_mode(
            cache.cached_flag.iter().cloned(),
            CodegenMode::Production,
        )
        .unwrap();
        assert_eq!(cache.cached_flag.len(), p_parsed_flags.len());
        for (i, item) in p_parsed_flags.iter().enumerate() {
            assert!(cache.cached_flag[i].eq(item));
        }
    }

    #[test]
    fn test_modify_cached_flags_based_on_mode_exported() {
        let cache = crate::test::create_test_cache();
        let p_parsed_flags = modify_cached_flags_based_on_mode(
            cache.cached_flag.iter().cloned(),
            CodegenMode::Exported,
        )
        .unwrap();
        assert_eq!(3, p_parsed_flags.len());
        for flag in p_parsed_flags.iter() {
            assert_eq!(ProtoFlagState::DISABLED, flag.state());
            assert_eq!(ProtoFlagPermission::READ_WRITE, flag.permission());
            assert!(!flag.is_fixed_read_only());
            assert!(flag.is_exported());
        }

        let mut cache = crate::test::create_test_cache();
        cache.cached_flag.retain(|pf| !pf.is_exported());
        let error = modify_cached_flags_based_on_mode(
            cache.cached_flag.iter().cloned(),
            CodegenMode::Exported,
        )
        .unwrap_err();
        assert_eq!("exported library contains no exported flags", format!("{:?}", error));
    }
}
