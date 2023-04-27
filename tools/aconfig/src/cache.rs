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

use anyhow::{anyhow, Result};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::io::{Read, Write};

use crate::aconfig::{Flag, Override};
use crate::commands::Source;

#[derive(Serialize, Deserialize)]
enum FlagOrOverride {
    Flag(Flag),
    Override(Override),
}

#[derive(Serialize, Deserialize)]
struct Value {
    source: Source,
    value: FlagOrOverride,
}

#[derive(Serialize, Deserialize)]
pub struct Cache {
    values: HashMap<String, Vec<Value>>,
}

impl Cache {
    pub fn new() -> Cache {
        Cache { values: HashMap::new() }
    }

    pub fn read_from_reader(reader: impl Read) -> Result<Cache> {
        serde_json::from_reader(reader).map_err(|e| e.into())
    }

    pub fn write_to_writer(&self, writer: impl Write) -> Result<()> {
        serde_json::to_writer(writer, self).map_err(|e| e.into())
    }

    pub fn add_flag(&mut self, source: Source, flag: Flag) -> Result<()> {
        if let Some(values) = self.values.get(&flag.id) {
            let first_source = &values[0].source;
            return Err(anyhow!(
                "failed to add flag {} from {}: already added from {}",
                flag.id,
                source,
                first_source
            ));
        }
        self.values
            .insert(flag.id.clone(), vec![Value { source, value: FlagOrOverride::Flag(flag) }]);
        Ok(())
    }

    pub fn add_override(&mut self, source: Source, override_: Override) -> Result<()> {
        if let Some(values) = self.values.get_mut(&override_.id) {
            values.push(Value { source, value: FlagOrOverride::Override(override_) });
            Ok(())
        } else {
            Err(anyhow!("failed to override flag {}: unknown flag", override_.id))
        }
    }

    pub fn keys(&self) -> Vec<String> {
        self.values.keys().cloned().collect()
    }

    pub fn value(&self, id: &str) -> Option<bool> {
        self.values.get(id).map(|values| match &values.last().unwrap().value {
            FlagOrOverride::Flag(flag) => flag.value,
            FlagOrOverride::Override(override_) => override_.value,
        })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_add_flag() {
        let mut cache = Cache::new();
        cache
            .add_flag(
                Source::File("first.txt".to_string()),
                Flag { id: "foo".to_string(), description: "desc".to_string(), value: true },
            )
            .unwrap();
        let error = cache
            .add_flag(
                Source::File("second.txt".to_string()),
                Flag { id: "foo".to_string(), description: "desc".to_string(), value: false },
            )
            .unwrap_err();
        assert_eq!(
            &format!("{:?}", error),
            "failed to add flag foo from file second.txt: already added from file first.txt"
        );
    }

    #[test]
    fn test_add_override() {
        let mut cache = Cache::new();
        let error = cache
            .add_override(Source::Memory, Override { id: "foo".to_string(), value: false })
            .unwrap_err();
        assert_eq!(&format!("{:?}", error), "failed to override flag foo: unknown flag");

        cache
            .add_flag(
                Source::File("first.txt".to_string()),
                Flag { id: "foo".to_string(), description: "desc".to_string(), value: true },
            )
            .unwrap();
        assert_eq!(Some(true), cache.value("foo"));

        cache
            .add_override(Source::Memory, Override { id: "foo".to_string(), value: false })
            .unwrap();
        assert_eq!(Some(false), cache.value("foo"));

        cache
            .add_override(Source::Memory, Override { id: "foo".to_string(), value: true })
            .unwrap();
        assert_eq!(Some(true), cache.value("foo"));
    }
}
