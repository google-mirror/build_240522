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

use crate::{AconfigDefinition, Flag, FlagSource, FlagState, Permission};
use aconfig_protos::ProtoFlagPermission as ProtoPermission;
use aconfig_protos::ProtoFlagState as ProtoState;
use aconfig_protos::ProtoParsedFlag;
use aconfig_protos::ProtoParsedFlags;
use anyhow::{bail, Result};
use regex::Regex;
use std::collections::BTreeMap;
use std::collections::HashMap;
use std::process::Command;
use std::{fs, str};

pub struct DeviceConfigSource {}

fn convert_definition(flag: &ProtoParsedFlag) -> AconfigDefinition {
    let state = match flag.state() {
        ProtoState::ENABLED => FlagState::Enabled,
        ProtoState::DISABLED => FlagState::Disabled,
    };
    let permission = match flag.permission() {
        ProtoPermission::READ_ONLY => Permission::ReadOnly,
        ProtoPermission::READ_WRITE => Permission::ReadWrite,
    };
    AconfigDefinition { state, permission }
}

fn convert_parsed_flag(flag: &ProtoParsedFlag) -> Flag {
    let namespace = flag.namespace().to_string();
    let name = format!("{}.{}", flag.package(), flag.name());
    Flag { namespace, name, value: "".to_string(), definitions: vec![], overridden: false }
}

fn read_pb_files() -> Result<Vec<Flag>> {
    let mut flags: BTreeMap<String, Flag> = BTreeMap::new();
    for partition in ["system", "system_ext", "product", "vendor"] {
        let path = format!("/{}/etc/aconfig_flags.pb", partition);
        let Ok(bytes) = fs::read(&path) else {
            eprintln!("warning: failed to read {}", path);
            continue;
        };
        let parsed_flags: ProtoParsedFlags = protobuf::Message::parse_from_bytes(&bytes)?;
        for flag in parsed_flags.parsed_flag {
            let key = format!("{}/{}.{}", flag.namespace(), flag.package(), flag.name());
            flags
                .entry(key)
                .or_insert(convert_parsed_flag(&flag))
                .definitions
                .push(convert_definition(&flag));
        }
    }
    Ok(flags.values().cloned().collect())
}

fn parse_device_config(raw: &str) -> HashMap<String, String> {
    let mut flags = HashMap::new();
    let regex = Regex::new(r"(?m)^([[[:alnum:]]_]+/[[[:alnum:]]_\.]+)=(true|false)$").unwrap();
    for capture in regex.captures_iter(raw) {
        let key = capture.get(1).unwrap().as_str().to_string();
        let value = match capture.get(2).unwrap().as_str() {
            "true" => FlagState::Enabled.to_string(),
            "false" => FlagState::Disabled.to_string(),
            _ => panic!(),
        };
        flags.insert(key, value);
    }
    flags
}

fn read_device_config_output(command: &str) -> Result<String> {
    let output = Command::new("/system/bin/device_config").arg(command).output()?;
    if !output.status.success() {
        let reason = match output.status.code() {
            Some(code) => format!("exit code {}", code),
            None => "terminated by signal".to_string(),
        };
        bail!("failed to execute device_config: {}", reason);
    }
    Ok(str::from_utf8(&output.stdout)?.to_string())
}

fn read_device_config_flags(
    overrides: &HashMap<String, String>,
) -> Result<HashMap<String, String>> {
    let list_output = read_device_config_output("list")?;
    let mut flags = parse_device_config(&list_output);
    for (key, value) in overrides.iter() {
        if flags.contains_key(key) {
            flags.insert(key.to_string(), value.to_string());
        }
    }
    Ok(flags)
}

fn reconcile(pb_flags: &[Flag], dc_flags: HashMap<String, String>) -> Vec<Flag> {
    pb_flags
        .iter()
        .map(|f| {
            dc_flags
                .get(&format!("{}/{}", f.namespace, f.name))
                .map(|value| Flag { value: value.to_string(), ..f.clone() })
                .unwrap_or(f.clone())
        })
        .collect()
}

impl FlagSource for DeviceConfigSource {
    fn list_flags(namespace: Option<String>) -> Result<Vec<Flag>> {
        let pb_flags = read_pb_files()?;

        let dc_overrides_output = read_device_config_output("list_local_overrides")?;
        let dc_overrides = parse_device_config(&dc_overrides_output);
        let dc_flags = read_device_config_flags(&dc_overrides)?;

        let mut flags = reconcile(&pb_flags, dc_flags);
        for flag in &mut flags {
            if dc_overrides.contains_key(&(flag.namespace.clone() + "/" + &flag.name)) {
                flag.overridden = true;
            }
        }

        let result = match namespace {
            Some(n) => flags.iter().filter(|f| f.namespace == n).cloned().collect(),
            None => flags,
        };
        Ok(result)
    }
}
