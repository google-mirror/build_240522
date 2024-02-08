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
            FlagState::Enabled => write!(f, "T"),
            FlagState::Disabled => write!(f, "F"),
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
            FlagPermission::ReadOnly => write!(f, "R"),
            FlagPermission::ReadWrite => write!(f, "W"),
        }
    }
}

#[derive(Clone)]
struct AconfigDefinition {
    state: FlagState,
    permission: FlagPermission,
}

#[derive(Clone)]
struct Flag {
    namespace: String,
    name: String,
    definitions: Vec<AconfigDefinition>,
    value: String,
    overridden: bool,
}

trait FlagSource {
    fn list_flags() -> Result<Vec<Flag>>;
}

const ABOUT_TEXT: &str = "Tool for reading and writing flags.

Rows in the table from the `list` command follow this format:

  namespace package.flag_name A B CCCC DDDD

  * `A`: flag value.
  * `B`: `O` if there is a manual local override, which prevents server updates, else `-`.
  * `C`: `R` if read only, `W` if read-write, for each of the four partitions:
      1. `system`
      2. `system_ext`
      3. `product`
      4. `vendor`
  * `D`: Again for each of the four partitions, `T` if default-enabled, `F` if
    default-disabled.
";

#[derive(Parser, Debug)]
#[clap(long_about=ABOUT_TEXT)]
struct Cli {
    /// Rows follow this format:
    #[clap(subcommand)]
    command: Command,
}

#[derive(Parser, Debug)]
enum Command {
    /// List all aconfig flags on this device
    List,
}

struct PaddingInfo {
    longest_namespace_col: usize,
    longest_name_col: usize,
    longest_val_col: usize,
}

fn format_flag_row(flag: &Flag, info: &PaddingInfo) -> String {
    let ns = &flag.namespace;
    let p0 = " ".repeat(info.longest_namespace_col - ns.len() + 1);

    let name = &flag.name;
    let p1 = " ".repeat(info.longest_name_col - name.len() + 1);

    let val = &flag.value;
    let p2 = " ".repeat(info.longest_val_col - val.len() + 1);

    let pm_s =
        flag.definitions.first().map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_se =
        flag.definitions.get(1).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_p =
        flag.definitions.get(2).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_v =
        flag.definitions.get(3).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());

    let o = if flag.overridden { "O" } else { "-" };

    let s_s = flag.definitions.first().map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_se = flag.definitions.get(1).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_p = flag.definitions.get(2).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_v = flag.definitions.get(3).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());

    format!("{ns}{p0}{name}{p1}{val}{p2}{o} {pm_s}{pm_se}{pm_p}{pm_v} {s_s}{s_se}{s_p}{s_v}\n")
}

fn list() -> Result<String> {
    let flags = DeviceConfigSource::list_flags()?;
    let padding_info = PaddingInfo {
        longest_namespace_col: flags.iter().map(|f| f.namespace.len()).max().unwrap_or(0),
        longest_name_col: flags.iter().map(|f| f.name.len()).max().unwrap_or(0),
        longest_val_col: flags.iter().map(|f| f.value.len()).max().unwrap_or(0),
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
