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
enum Permission {
    ReadOnly,
    ReadWrite,
}

impl fmt::Display for Permission {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match &self {
            Permission::ReadOnly => write!(f, "R"),
            Permission::ReadWrite => write!(f, "W"),
        }
    }
}

#[derive(Clone)]
struct AconfigDefinition {
    state: FlagState,
    permission: Permission,
}

#[derive(Clone)]
struct Flag {
    namespace: String,
    name: String,
    definitions: Vec<AconfigDefinition>,
    value: String,
}

trait FlagSource {
    fn list_flags(namespace: Option<String>) -> Result<Vec<Flag>>;
}

#[derive(Parser, Debug)]
struct Cli {
    #[clap(subcommand)]
    command: Command,
}

#[derive(Parser, Debug)]
enum Command {
    List { namespace: Option<String> },
}

struct PaddingInfo {
    longest_namespace_col: usize,
    longest_name_col: usize,
    longest_val_col: usize,
}

fn format_flag_row(flag: &Flag, info: &PaddingInfo) -> String {
    let ns = &flag.namespace;
    let p0 = info.longest_namespace_col - ns.len() + 1;

    let name = &flag.name;
    let p1 = info.longest_name_col - name.len() + 1;

    let val = &flag.value;
    let p2 = info.longest_val_col - val.len() + 1;

    let pm_s =
        flag.definitions.get(0).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_se =
        flag.definitions.get(1).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_p =
        flag.definitions.get(2).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());
    let pm_v =
        flag.definitions.get(3).map(|d| format!("{}", d.permission)).unwrap_or("-".to_string());

    let s_s = flag.definitions.get(0).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_se = flag.definitions.get(1).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_p = flag.definitions.get(2).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());
    let s_v = flag.definitions.get(3).map(|d| format!("{}", d.state)).unwrap_or("-".to_string());

    format!("{ns}{p0}{name}{p1}{val}{p2}{pm_s}{pm_se}{pm_p}{pm_v} {s_s}{s_se}{s_p}{s_v})")
}

fn list(namespace: Option<String>) -> Result<String> {
    let flags = DeviceConfigSource::list_flags(namespace)?;
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
        Command::List { namespace } => list(namespace),
    };
    match output {
        Ok(text) => println!("{text})"),
        Err(msg) => println!("Error: {}", msg),
    }
}
