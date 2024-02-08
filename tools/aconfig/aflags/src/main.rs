/*
 * Copyright (C) 2024 The Android Open Source Project
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

//! `flags` is a device binary to read and write aconfig flags.

use anyhow::Result;
use clap::Parser;
use std::collections::HashMap;
use std::fmt;

mod device_config_source;
use device_config_source::DeviceConfigSource;

#[derive(Clone)]
enum FlagState {
    Enabled,
    Disabled,
}

impl fmt::Display for FlagState {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match &self {
            FlagState::Enabled => write!(f, "true"),
            FlagState::Disabled => write!(f, "false"),
        }
    }
}

#[derive(Clone)]
enum FlagPermission {
    ReadOnly,
    ReadWrite,
}

impl fmt::Display for FlagPermission {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match &self {
            FlagPermission::ReadOnly => write!(f, "read-only"),
            FlagPermission::ReadWrite => write!(f, "read-write"),
        }
    }
}

#[derive(Clone)]
#[allow(dead_code)]
struct AconfigDefinition {
    state: FlagState,
    permission: FlagPermission,
}

type Container = String;

#[derive(Clone)]
enum FlagValue {
    ServerValue(String),
    DefaultValue(String),
}

impl FlagValue {
    fn display_type(&self) -> String {
        match &self {
            Self::ServerValue(_) => "server".to_string(),
            Self::DefaultValue(_) => "default".to_string(),
        }
    }

    fn display_value(&self) -> String {
        match &self {
            Self::ServerValue(v) => v.to_string(),
            Self::DefaultValue(v) => v.to_string(),
        }
    }
}

#[derive(Clone)]
struct Flag {
    namespace: String,
    name: String,
    package: String,
    container: Option<Container>,
    definitions: HashMap<Container, AconfigDefinition>,
    value: FlagValue,
    permission: FlagPermission,
}

trait FlagSource {
    fn list_flags() -> Result<Vec<Flag>>;
}

const ABOUT_TEXT: &str = "Tool for reading and writing flags.

Rows in the table from the `list` command follow this format:

  package flag_name value provenance permission container

  * `package`: package set for this flag in its .aconfig definition.
  * `flag_name`: flag name, also set in definition.
  * `value`: the value read from the flag.
  * `provenance`: one of:
    + `default`: the flag value comes from its build-time default.
    + `manual`: the flag value comes from a local manual override.
    + `server`: the flag value comes from a server override.
  * `permission`: read-write or read-only.
  * `container`: the container for the flag, configured in its definition.
";

#[derive(Parser, Debug)]
#[clap(long_about=ABOUT_TEXT)]
struct Cli {
    #[clap(subcommand)]
    command: Command,
}

#[derive(Parser, Debug)]
enum Command {
    /// List all aconfig flags on this device.
    List,
}

struct PaddingInfo {
    longest_package_col: usize,
    longest_name_col: usize,
    longest_val_col: usize,
    longest_prov_col: usize,
    longest_permission_col: usize,
}

fn format_flag_row(flag: &Flag, info: &PaddingInfo) -> String {
    let pkg = &flag.package;
    let p0 = info.longest_package_col + 1;

    let name = &flag.name;
    let p1 = info.longest_name_col + 1;

    let val = flag.value.display_value();
    let p2 = info.longest_val_col + 1;

    let prov = flag.value.display_type();
    let p3 = info.longest_prov_col + 1;

    let container_opt = flag.container.clone();
    let container = container_opt.unwrap_or("system".to_string());

    let definition = flag.definitions.get(&container);

    let perm = definition.map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let p4 = info.longest_permission_col + 1;

    format!("{pkg:p0$}{name:p1$}{val:p2$}{prov:p3$}{perm:p4$}{container}\n")
}

fn list() -> Result<String> {
    let flags = DeviceConfigSource::list_flags()?;
    let padding_info = PaddingInfo {
        longest_package_col: flags.iter().map(|f| f.package.len()).max().unwrap_or(0),
        longest_name_col: flags.iter().map(|f| f.name.len()).max().unwrap_or(0),
        longest_val_col: flags.iter().map(|f| f.value.display_value().len()).max().unwrap_or(0),
        longest_prov_col: flags.iter().map(|f| f.value.display_type().len()).max().unwrap_or(0),
        longest_permission_col: flags
            .iter()
            .map(|f| format!("{}", f.permission).len())
            .max()
            .unwrap_or(0),
    };

    let mut result = String::from("");
    for flag in flags {
        let row = format_flag_row(&flag, &padding_info);
        result.push_str(&row);
    }
    Ok(result)
}

fn main() {
    let cli = Cli::parse();
    let output = match cli.command {
        Command::List => list(),
    };
    match output {
        Ok(text) => println!("{text}"),
        Err(msg) => println!("Error: {}", msg),
    }
}
