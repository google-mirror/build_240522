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

use crate::protos::ProtoStorageFiles;
use anyhow::Result;
use protobuf::Message;
use std::fs;
use std::io::Write;
use tempfile::NamedTempFile;

pub(crate) fn get_binary_storage_proto_bytes(text_proto: &str) -> Result<Vec<u8>> {
    let storage_files: ProtoStorageFiles = protobuf::text_format::parse_from_str(text_proto)?;
    let mut binary_proto = Vec::new();
    storage_files.write_to_vec(&mut binary_proto)?;
    Ok(binary_proto)
}

pub(crate) fn write_storage_text_to_temp_file(text_proto: &str) -> Result<NamedTempFile> {
    let bytes = get_binary_storage_proto_bytes(text_proto).unwrap();
    let mut file = NamedTempFile::new()?;
    let _ = file.write_all(&bytes);
    Ok(file)
}

pub(crate) fn copy_to_temp_read_only_file(source_file: &str) -> Result<NamedTempFile> {
    let file = NamedTempFile::new()?;
    fs::copy(source_file, file.path()).unwrap();
    let mut perms = fs::metadata(file.path()).unwrap().permissions();
    if !perms.readonly() {
        perms.set_readonly(true);
        fs::set_permissions(file.path(), perms).unwrap();
    }
    Ok(file)
}

pub(crate) fn copy_to_temp_read_write_file(source_file: &str) -> Result<NamedTempFile> {
    let file = NamedTempFile::new()?;
    fs::copy(source_file, file.path()).unwrap();
    let mut perms = fs::metadata(file.path()).unwrap().permissions();
    if perms.readonly() {
        perms.set_readonly(false);
        fs::set_permissions(file.path(), perms).unwrap();
    }
    Ok(file)
}
