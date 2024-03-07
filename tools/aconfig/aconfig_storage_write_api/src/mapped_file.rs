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

use std::collections::HashMap;
use std::fs::{self, File, OpenOptions};
use std::io::{BufReader, Read};
use std::sync::{Arc, Mutex};

use anyhow::anyhow;
use memmap2::MmapMut;
use once_cell::sync::Lazy;

use aconfig_storage_file::protos::{storage_record_pb::try_from_binary_proto, ProtoStorageFiles};
use aconfig_storage_file::AconfigStorageError::{
    self, FileReadFail, MapFileFail, ProtobufParseFail, StorageFileNotFound,
};

/// Cache for already mapped files
static ALL_MAPPED_FILES: Lazy<Mutex<HashMap<String, Arc<MmapMut>>>> = Lazy::new(|| {
    let mapped_files = HashMap::new();
    Mutex::new(mapped_files)
});

/// Find where persistent storage value file is stored for a particular container
fn find_container_persist_flag_value_location(
    location_pb_file: &str,
    container: &str,
) -> Result<String, AconfigStorageError> {
    let file = File::open(location_pb_file).map_err(|errmsg| {
        FileReadFail(anyhow!("Failed to open file {}: {}", location_pb_file, errmsg))
    })?;
    let mut reader = BufReader::new(file);
    let mut bytes = Vec::new();
    reader.read_to_end(&mut bytes).map_err(|errmsg| {
        FileReadFail(anyhow!("Failed to read file {}: {}", location_pb_file, errmsg))
    })?;
    let storage_locations: ProtoStorageFiles = try_from_binary_proto(&bytes).map_err(|errmsg| {
        ProtobufParseFail(anyhow!(
            "Failed to parse storage location pb file {}: {}",
            location_pb_file,
            errmsg
        ))
    })?;
    for location_info in storage_locations.files.iter() {
        if location_info.container() == container {
            return Ok(location_info.flag_val().to_string());
        }
    }
    Err(StorageFileNotFound(anyhow!("Persistent flag value file does not exist for {}", container)))
}

/// Verify the file is read write and then map it
unsafe fn verify_read_write_and_map(file_path: &str) -> Result<MmapMut, AconfigStorageError> {
    // ensure file has read write permission
    let perms = fs::metadata(file_path).unwrap().permissions();
    if perms.readonly() {
        return Err(MapFileFail(anyhow!("fail to map non read write storage file {}", file_path)));
    }

    let file =
        OpenOptions::new().read(true).write(true).open(file_path).map_err(|errmsg| {
            FileReadFail(anyhow!("Failed to open file {}: {}", file_path, errmsg))
        })?;

    unsafe {
        let mapped_file = MmapMut::map_mut(&file).map_err(|errmsg| {
            MapFileFail(anyhow!("fail to map storage file {}: {}", file_path, errmsg))
        })?;
        Ok(mapped_file)
    }
}

/// Get a mapped storage file given the container and file type
pub(crate) unsafe fn get_mapped_file(
    location_pb_file: &str,
    container: &str,
) -> Result<Arc<MmapMut>, AconfigStorageError> {
    let mut all_mapped_files = ALL_MAPPED_FILES.lock().unwrap();
    match all_mapped_files.get(container) {
        Some(mapped_file) => Ok(Arc::clone(mapped_file)),
        None => {
            let value_file =
                find_container_persist_flag_value_location(location_pb_file, container)?;
            unsafe {
                let mapped_file = Arc::new(verify_read_write_and_map(&value_file)?);
                all_mapped_files.insert(container.to_string(), Arc::clone(&mapped_file));
                Ok(mapped_file)
            }
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::test_utils::TestStorageFile;
    use aconfig_storage_file::protos::storage_record_pb::write_proto_to_temp_file;

    #[test]
    fn test_find_persist_flag_value_file_location() {
        let text_proto = r#"
files {
    version: 0
    container: "system"
    package_map: "/system/etc/package.map"
    flag_map: "/system/etc/flag.map"
    flag_val: "/metadata/aconfig/system.val"
    timestamp: 12345
}
files {
    version: 1
    container: "product"
    package_map: "/product/etc/package.map"
    flag_map: "/product/etc/flag.map"
    flag_val: "/metadata/aconfig/product.val"
    timestamp: 54321
}
"#;
        let file = write_proto_to_temp_file(&text_proto).unwrap();
        let file_full_path = file.path().display().to_string();
        let flag_value_file =
            find_container_persist_flag_value_location(&file_full_path, "system").unwrap();
        assert_eq!(flag_value_file, "/metadata/aconfig/system.val");
        let flag_value_file =
            find_container_persist_flag_value_location(&file_full_path, "product").unwrap();
        assert_eq!(flag_value_file, "/metadata/aconfig/product.val");
        let err =
            find_container_persist_flag_value_location(&file_full_path, "vendor").unwrap_err();
        assert_eq!(
            format!("{:?}", err),
            "StorageFileNotFound(Persistent flag value file does not exist for vendor)"
        );
    }

    #[test]
    fn test_mapped_file_contents() {
        let mut rw_file = TestStorageFile::new("./tests/flag.val", false).unwrap();
        let text_proto = format!(
            r#"
files {{
    version: 0
    container: "test_mapped_file_contents"
    package_map: "some_package.map"
    flag_map: "some_flag.map"
    flag_val: "{}"
    timestamp: 12345
}}
"#,
            rw_file.name
        );
        let storage_record_file = write_proto_to_temp_file(&text_proto).unwrap();
        let storage_record_file_path = storage_record_file.path().display().to_string();

        let mut content = Vec::new();
        rw_file.file.read_to_end(&mut content).unwrap();

        unsafe {
            let mmaped_file =
                get_mapped_file(&storage_record_file_path, "test_mapped_file_contents").unwrap();
            assert_eq!(mmaped_file[..], content[..]);
        }
    }

    #[test]
    fn test_mapped_read_only_file() {
        let ro_file = TestStorageFile::new("./tests/flag.val", true).unwrap();
        let text_proto = format!(
            r#"
files {{
    version: 0
    container: "test_mapped_read_only_file"
    package_map: "some_package.map"
    flag_map: "some_flag.map"
    flag_val: "{}"
    timestamp: 12345
}}
"#,
            ro_file.name
        );
        let storage_record_file = write_proto_to_temp_file(&text_proto).unwrap();
        let storage_record_file_path = storage_record_file.path().display().to_string();

        unsafe {
            let error = get_mapped_file(&storage_record_file_path, "test_mapped_read_only_file")
                .unwrap_err();
            assert_eq!(
                format!("{:?}", error),
                format!("MapFileFail(fail to map non read write storage file {})", ro_file.name)
            );
        }
    }
}
