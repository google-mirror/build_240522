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

//! `printflags` is a device binary to print feature flags.

use aconfig_protos::aconfig::Flag_permission;
use aconfig_protos::aconfig::Flag_state;
use aconfig_protos::aconfig::Parsed_flags as ProtoParsedFlags;
use anyhow::{anyhow, Context, Result};
use itertools::Itertools;
use once_cell::sync::Lazy;
use regex::Regex;
use std::fmt::Write;
use std::process::Command;
use std::{fmt, fs, str};

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum Permission {
    ReadWrite,
    ReadOnly,
}

impl From<Flag_permission> for Permission {
    fn from(perm: Flag_permission) -> Permission {
        match perm {
            Flag_permission::READ_WRITE => Permission::ReadWrite,
            Flag_permission::READ_ONLY => Permission::ReadOnly,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum State {
    Enabled,
    Disabled,
}

impl From<Flag_state> for State {
    fn from(state: Flag_state) -> State {
        match state {
            Flag_state::ENABLED => State::Enabled,
            Flag_state::DISABLED => State::Disabled,
        }
    }
}

#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
enum Source {
    System,
    SystemExt,
    Product,
    Vendor,
    DeviceConfig,
    SysProp,
}

impl Source {
    /// Sources whose states can be modified at runtime
    fn is_runtime(&self) -> bool {
        matches!(self, Source::DeviceConfig | Source::SysProp)
    }
}

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord)]
struct FqName {
    /// The gantry namespace
    namespace: String,
    /// The gantry flagname
    name: String,
}

impl FqName {
    fn new(namespace: &str, name: &str) -> FqName {
        FqName { namespace: namespace.to_string(), name: name.to_string() }
    }
}

/// Represents the state of an aconfig flag in a particular source
/// Permission is always ReadWrite for runtime flags
#[derive(Debug, Clone, Copy, PartialEq, Eq, PartialOrd, Ord)]
struct SourceState(Source, State, Permission);

impl fmt::Display for SourceState {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(
            f,
            "{:?}: {}{}",
            self.0,
            match self.1 {
                State::Enabled => "Enabled",
                State::Disabled => "Disabled",
            },
            match self.2 {
                Permission::ReadOnly => " (RO)",
                _ => "",
            }
        )
    }
}

static CONFIG_RE: Lazy<Regex> =
    Lazy::new(|| Regex::new(r"(?m)^([[[:alnum:]]_]+)/([[[:alnum:]]_\.]+)=(true|false)$").unwrap());
static SYSPROP_RE: Lazy<Regex> = Lazy::new(|| {
    Regex::new(concat!(r"(?m)^\[persist\.device_config\.aconfig_flags\.([[[:alnum:]]_]+)\.",
        r"([[[:alnum:]]_\.]+)\]: \[(true|false)\]$")).unwrap()
});

fn parse_with_regex(raw: &str, source: Source) -> Vec<(FqName, SourceState)> {
    match source {
        Source::DeviceConfig => &CONFIG_RE,
        Source::SysProp => &SYSPROP_RE,
        _ => panic!("Invalid source for regex parse"),
    }
    .captures_iter(raw)
    .map(move |capture| {
        let unwrap = |x| capture.get(x).expect("Invalid regex").as_str();
        let (namespace, name, state) = (unwrap(1), unwrap(2), unwrap(3));
        let value = match state {
            "true" => State::Enabled,
            "false" => State::Disabled,
            _ => panic!("Invalid regex"),
        };

        (FqName::new(namespace, name), SourceState(source, value, Permission::ReadWrite))
    })
    .collect_vec()
}

fn get_proto_flags(partition: Source) -> Result<Vec<(FqName, SourceState)>> {
    let path = format!(
        "/{}/etc/aconfig_flags.pb",
        match partition {
            Source::System => "system",
            Source::SystemExt => "system_ext",
            Source::Product => "product",
            Source::Vendor => "vendor",
            _ => panic!("Invalid source"),
        }
    );

    let data = fs::read(&path).with_context(|| format!("Failed to read flags from {}", path))?;

    let proto_flags: ProtoParsedFlags =
        protobuf::Message::parse_from_bytes(&data).with_context(|| {
            format!("failed to parse {} ({}, {} byte(s))", path, xxd(&data), data.len())
        })?;

    Ok(proto_flags
        .parsed_flag
        .into_iter()
        .map(move |mut flag| {
            (
                FqName {
                    namespace: flag.take_namespace(),
                    name: flag.take_package() + "." + flag.name(),
                },
                SourceState(partition, flag.state().into(), flag.permission().into()),
            )
        })
        .sorted()
        .collect_vec())
}

fn run_command(cmd: &mut Command) -> Result<String> {
    let output = cmd.output()?;
    if output.status.success() {
        str::from_utf8(&output.stdout).map_err(anyhow::Error::new)
    } else {
        match output.status.code() {
            Some(code) => Err(anyhow!("exit code {}", code)),
            None => Err(anyhow!("terminated by signal".to_string())),
        }
    }
    .with_context(|| format!("cmd {:?} failed", cmd))
    .map(|s| s.to_string())
}

fn xxd(bytes: &[u8]) -> String {
    let n = 8.min(bytes.len());
    let mut v = Vec::with_capacity(n);
    for byte in bytes.iter().take(n) {
        v.push(format!("{:02x}", byte));
    }
    let trailer = match bytes.len() {
        0..=8 => "",
        _ => " ..",
    };
    format!("[{}{}]", v.join(" "), trailer)
}

fn sorted_group_fold<T, U, K, V>(iters: T) -> impl Iterator<Item = (K, Vec<V>)>
where
    T: Iterator<Item = U> + Sized,
    U: IntoIterator<Item = (K, V)>,
    K: std::cmp::Ord,
    V: std::cmp::Ord,
{
    iters.kmerge_by(|x, y| x.0 < y.0).peekable().batching(|it| {
        let (key, val) = it.next()?;
        let mut vals = vec![val];
        vals.extend(it.peeking_take_while(|(nkey, _)| nkey == &key).map(|(_, y)| y));
        vals.sort();
        Some((key, vals))
    })
}

fn main() -> Result<()> {

    let runtime_iter = [
        (run_command(Command::new("/system/bin/device_config").arg("list")), Source::DeviceConfig),
        (run_command(&mut Command::new("/system/bin/getprop")), Source::SysProp),
    ]
    .into_iter()
    .filter_map(|(raw, source)| raw.map(|it| parse_with_regex(it.as_ref(), source)).ok());

    let proto_iter = [Source::System, Source::SystemExt, Source::Product, Source::Vendor]
        .into_iter()
        .filter_map(|x| get_proto_flags(x).ok());

    sorted_group_fold(proto_iter.chain(runtime_iter))
        .filter(|(_, state_list)| {
            state_list.iter().any(|SourceState(source, _, _)| !source.is_runtime())
        })
        .for_each(|(flagname, value)| {
            let mut acc = value.get(0).map(|x| x.to_string()).unwrap_or_default();
            for elem in value.iter().skip(1) {
                write!(&mut acc, ", {}", elem).expect("Writing to string");
            }
            println!("{}/{}: [{}]", flagname.namespace, flagname.name, acc)
        });
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_parse_device_config() {
        let input = r#"
namespace_one/com.foo.bar.flag_one=true
namespace_one/com.foo.bar.flag_two=false
random_noise;
namespace_two/android.flag_one=true
namespace_two/android.flag_two=nonsense
"#;
        let expected = vec![
            (
                FqName::new("namespace_one", "com.foo.bar.flag_one"),
                SourceState(Source::DeviceConfig, State::Enabled, Permission::ReadWrite),
            ),
            (
                FqName::new("namespace_one", "com.foo.bar.flag_two"),
                SourceState(Source::DeviceConfig, State::Disabled, Permission::ReadWrite),
            ),
            (
                FqName::new("namespace_two", "android.flag_one"),
                SourceState(Source::DeviceConfig, State::Enabled, Permission::ReadWrite),
            ),
        ];
        let actual = parse_with_regex(input, Source::DeviceConfig);
        assert_eq!(expected, actual);
    }

    #[test]
    fn test_parse_sysprop() {
        let input = r#"
[persist.device_config.aconfig_flags.namespace_one.com.package1.flag_1]: [false]
[persist.device_config.aconfig_flags.namespace_one.com.package1.flag_2]: [true]
[persist.device_config.aconfig_flags.namespace_two.com.package2.flag_1]: [true]
[persist.device_config.aconfig_flags.namespace_two.com.package2.flag_2]: [0]
[persist.device_config.aconfig_flags.namespace_two]: [true]
[persist.device_config.something_else.namespace.flag]: [false]
"#;
        let expected = vec![
            (
                FqName::new("namespace_one", "com.package1.flag_1"),
                SourceState(Source::SysProp, State::Disabled, Permission::ReadWrite),
            ),
            (
                FqName::new("namespace_one", "com.package1.flag_2"),
                SourceState(Source::SysProp, State::Enabled, Permission::ReadWrite),
            ),
            (
                FqName::new("namespace_two", "com.package2.flag_1"),
                SourceState(Source::SysProp, State::Enabled, Permission::ReadWrite),
            ),
        ];
        let actual = parse_with_regex(input, Source::SysProp);
        assert_eq!(expected, actual);
    }

    #[test]
    fn test_xxd() {
        let input = [0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9];
        assert_eq!("[]", &xxd(&input[0..0]));
        assert_eq!("[00]", &xxd(&input[0..1]));
        assert_eq!("[00 01]", &xxd(&input[0..2]));
        assert_eq!("[00 01 02 03 04 05 06]", &xxd(&input[0..7]));
        assert_eq!("[00 01 02 03 04 05 06 07]", &xxd(&input[0..8]));
        assert_eq!("[00 01 02 03 04 05 06 07 ..]", &xxd(&input[0..9]));
        assert_eq!("[00 01 02 03 04 05 06 07 ..]", &xxd(&input));
    }

    #[test]
    fn test_sorted_group_fold() {
        let input = [
            vec![("abc", 1), ("cde", 2)],
            vec![("efg", 5)],
            vec![("abc", 3), ("def", 4), ("efg", 7)],
            vec![("cde", 3), ("def", 7)],
            vec![("abc", 0), ("efg", 6)],
        ];

        let mut expected = vec![
            ("abc", vec![0, 1, 3]),
            ("cde", vec![2, 3]),
            ("def", vec![4, 7]),
            ("efg", vec![5, 6, 7]),
        ];
        expected.sort();
        let actual: Vec<_> = sorted_group_fold(input.into_iter()).collect();
        assert_eq!(expected, actual);
    }
}
