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

use anyhow::{anyhow, Error, Result};
use serde::{Deserialize, Serialize};

use crate::protos::{ProtoFlag, ProtoFlags, ProtoOverride, ProtoOverrides};

#[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct Flag {
    pub id: String,
    pub description: String,
    pub value: bool,
}

impl Flag {
    #[allow(dead_code)] // only used in unit tests
    pub fn try_from_text_proto(text_proto: &str) -> Result<Flag> {
        let proto_flag: ProtoFlag = crate::protos::try_from_text_proto(text_proto)?;
        proto_flag.try_into()
    }

    pub fn try_from_text_proto_list(text_proto: &str) -> Result<Vec<Flag>> {
        let proto_flags: ProtoFlags = crate::protos::try_from_text_proto(text_proto)?;
        proto_flags.flag.into_iter().map(|proto_flag| proto_flag.try_into()).collect()
    }
}

impl TryFrom<ProtoFlag> for Flag {
    type Error = Error;

    fn try_from(proto_flag: ProtoFlag) -> Result<Self, Self::Error> {
        if proto_flag.id == String::default() {
            return Err(anyhow!("missing 'id' field"));
        }
        if proto_flag.description == String::default() {
            return Err(anyhow!("missing 'description' field"));
        }
        Ok(Flag { id: proto_flag.id, description: proto_flag.description, value: proto_flag.value })
    }
}

#[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct Override {
    pub id: String,
    pub value: bool,
}

impl Override {
    #[allow(dead_code)] // only used in unit tests
    pub fn try_from_text_proto(text_proto: &str) -> Result<Override> {
        let proto_flag: ProtoOverride = crate::protos::try_from_text_proto(text_proto)?;
        proto_flag.try_into()
    }

    pub fn try_from_text_proto_list(text_proto: &str) -> Result<Vec<Override>> {
        let proto: ProtoOverrides = crate::protos::try_from_text_proto(text_proto)?;
        proto.override_.into_iter().map(|proto_flag| proto_flag.try_into()).collect()
    }
}

impl TryFrom<ProtoOverride> for Override {
    type Error = Error;

    fn try_from(proto: ProtoOverride) -> Result<Self, Self::Error> {
        if proto.id == String::default() {
            return Err(anyhow!("missing 'id' field"));
        }
        Ok(Override { id: proto.id, value: proto.value })
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_flag_try_from_text_proto() {
        let expected = Flag {
            id: "1234".to_owned(),
            description: "Description of the flag".to_owned(),
            value: true,
        };

        let s = r#"
        id: "1234"
        description: "Description of the flag"
        value: true
        "#;
        let actual = Flag::try_from_text_proto(s).unwrap();

        assert_eq!(expected, actual);
    }

    #[test]
    fn test_flag_try_from_text_proto_missing_field() {
        let s = r#"
        description: "Description of the flag"
        value: true
        "#;
        let error = Flag::try_from_text_proto(s).unwrap_err();
        assert_eq!(&format!("{:?}", error), "missing 'id' field");
    }

    #[test]
    fn test_flag_try_from_text_proto_list() {
        let expected = vec![
            Flag { id: "a".to_owned(), description: "A".to_owned(), value: true },
            Flag { id: "b".to_owned(), description: "B".to_owned(), value: false },
        ];

        let s = r#"
        flag {
            id: "a"
            description: "A"
            value: true
        }
        flag {
            id: "b"
            description: "B"
            value: false
        }
        "#;
        let actual = Flag::try_from_text_proto_list(s).unwrap();

        assert_eq!(expected, actual);
    }

    #[test]
    fn test_override_try_from_text_proto_list() {
        let expected = Override { id: "1234".to_owned(), value: true };

        let s = r#"
        id: "1234"
        value: true
        "#;
        let actual = Override::try_from_text_proto(s).unwrap();

        assert_eq!(expected, actual);
    }
}
