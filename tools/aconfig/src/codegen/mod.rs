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

pub mod cpp;
pub mod java;
pub mod rust;

use crate::commands::CodegenMode;
use crate::protos::{ProtoFlagPermission, ProtoFlagState, ProtoParsedFlag};
use anyhow::{ensure, Result};

pub fn is_valid_name_ident(s: &str) -> bool {
    // Identifiers must match [a-z][a-z0-9_]*, except consecutive underscores are not allowed
    if s.contains("__") {
        return false;
    }
    let mut chars = s.chars();
    let Some(first) = chars.next() else {
        return false;
    };
    if !first.is_ascii_lowercase() {
        return false;
    }
    chars.all(|ch| ch.is_ascii_lowercase() || ch.is_ascii_digit() || ch == '_')
}

pub fn is_valid_package_ident(s: &str) -> bool {
    if !s.contains('.') {
        return false;
    }
    s.split('.').all(is_valid_name_ident)
}

pub fn is_valid_container_ident(s: &str) -> bool {
    is_valid_name_ident(s) || s.split('.').all(is_valid_name_ident)
}

pub fn create_device_config_ident(package: &str, flag_name: &str) -> Result<String> {
    ensure!(is_valid_package_ident(package), "bad package");
    ensure!(is_valid_name_ident(flag_name), "bad flag name");
    Ok(format!("{}.{}", package, flag_name))
}

pub fn process_parsed_flags<I>(
    parsed_flags_iter: I,
    codegen_mode: CodegenMode,
) -> Vec<ProtoParsedFlag>
where
    I: Iterator<Item = ProtoParsedFlag>,
{
    match codegen_mode {
        CodegenMode::Exported => {
            let mut parsed_flags: Vec<_> =
                parsed_flags_iter.filter(|pf| pf.is_exported()).collect();
            for flag in parsed_flags.iter_mut() {
                flag.set_state(ProtoFlagState::DISABLED);
                flag.set_permission(ProtoFlagPermission::READ_WRITE);
                flag.set_is_fixed_read_only(false);
            }
            parsed_flags
        }
        _ => parsed_flags_iter.collect(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_is_valid_name_ident() {
        assert!(is_valid_name_ident("foo"));
        assert!(is_valid_name_ident("foo_bar_123"));
        assert!(is_valid_name_ident("foo_"));

        assert!(!is_valid_name_ident(""));
        assert!(!is_valid_name_ident("123_foo"));
        assert!(!is_valid_name_ident("foo-bar"));
        assert!(!is_valid_name_ident("foo-b\u{00e5}r"));
        assert!(!is_valid_name_ident("foo__bar"));
        assert!(!is_valid_name_ident("_foo"));
    }

    #[test]
    fn test_is_valid_package_ident() {
        assert!(is_valid_package_ident("foo.bar"));
        assert!(is_valid_package_ident("foo.bar_baz"));
        assert!(is_valid_package_ident("foo.bar.a123"));

        assert!(!is_valid_package_ident("foo_bar_123"));
        assert!(!is_valid_package_ident("foo"));
        assert!(!is_valid_package_ident("foo._bar"));
        assert!(!is_valid_package_ident(""));
        assert!(!is_valid_package_ident("123_foo"));
        assert!(!is_valid_package_ident("foo-bar"));
        assert!(!is_valid_package_ident("foo-b\u{00e5}r"));
        assert!(!is_valid_package_ident("foo.bar.123"));
        assert!(!is_valid_package_ident(".foo.bar"));
        assert!(!is_valid_package_ident("foo.bar."));
        assert!(!is_valid_package_ident("."));
        assert!(!is_valid_package_ident(".."));
        assert!(!is_valid_package_ident("foo..bar"));
        assert!(!is_valid_package_ident("foo.__bar"));
    }

    #[test]
    fn test_is_valid_container_ident() {
        assert!(is_valid_container_ident("foo.bar"));
        assert!(is_valid_container_ident("foo.bar_baz"));
        assert!(is_valid_container_ident("foo.bar.a123"));
        assert!(is_valid_container_ident("foo"));
        assert!(is_valid_container_ident("foo_bar_123"));

        assert!(!is_valid_container_ident(""));
        assert!(!is_valid_container_ident("foo._bar"));
        assert!(!is_valid_container_ident("_foo"));
        assert!(!is_valid_container_ident("123_foo"));
        assert!(!is_valid_container_ident("foo-bar"));
        assert!(!is_valid_container_ident("foo-b\u{00e5}r"));
        assert!(!is_valid_container_ident("foo.bar.123"));
        assert!(!is_valid_container_ident(".foo.bar"));
        assert!(!is_valid_container_ident("foo.bar."));
        assert!(!is_valid_container_ident("."));
        assert!(!is_valid_container_ident(".."));
        assert!(!is_valid_container_ident("foo..bar"));
        assert!(!is_valid_container_ident("foo.__bar"));
    }

    #[test]
    fn test_create_device_config_ident() {
        assert_eq!(
            "com.foo.bar.some_flag",
            create_device_config_ident("com.foo.bar", "some_flag").unwrap()
        );
    }

    #[test]
    fn test_process_parsed_flags_prod() {
        let parsed_flags = crate::test::parse_test_flags();
        let p_parsed_flags = process_parsed_flags(
            parsed_flags.parsed_flag.clone().into_iter(),
            CodegenMode::Production,
        );
        assert_eq!(parsed_flags.parsed_flag.len(), p_parsed_flags.len());
        for (i, item) in p_parsed_flags.iter().enumerate() {
            assert!(parsed_flags.parsed_flag[i].eq(item));
        }
    }

    #[test]
    fn test_process_parsed_flags_exported() {
        let parsed_flags = crate::test::parse_test_flags();
        let p_parsed_flags =
            process_parsed_flags(parsed_flags.parsed_flag.into_iter(), CodegenMode::Exported);
<<<<<<< HEAD   (74993c aconfig: add exported mode in c/c++ codegen)
        assert_eq!(2, p_parsed_flags.len());
=======
        assert_eq!(3, p_parsed_flags.len());
>>>>>>> BRANCH (b0192b aconfig: add new testing flag enabled_fixed_ro_exported)
        for flag in p_parsed_flags.iter() {
            assert_eq!(ProtoFlagState::DISABLED, flag.state());
            assert_eq!(ProtoFlagPermission::READ_WRITE, flag.permission());
            assert!(!flag.is_fixed_read_only());
            assert!(flag.is_exported());
        }
    }
}
