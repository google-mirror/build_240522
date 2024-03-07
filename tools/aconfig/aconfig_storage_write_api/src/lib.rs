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

//! `aconfig_storage_write_api` is a crate that defines write apis to update flag value
//! in storage file. It provides one api to interface with storage files.

pub mod flag_value_update;
pub mod mapped_file;

#[cfg(test)]
mod test_utils;

use crate::flag_value_update::update_boolean_flag_value;
use crate::mapped_file::get_mapped_file;
use aconfig_storage_file::AconfigStorageError;

use anyhow::anyhow;
use std::sync::Arc;

/// Storage file location pb file
pub const STORAGE_LOCATION_FILE: &str = "/metadata/aconfig/persistent_storage_file_records.pb";

/// Set flag value in storage file implementation
pub unsafe fn set_boolean_flag_value_impl(
    pb_file: &str,
    container: &str,
    offset: u32,
    value: bool,
) -> Result<(), AconfigStorageError> {
    unsafe {
        let mut mapped_file_arc = get_mapped_file(pb_file, container)?;
        let mapped_file =
            Arc::get_mut(&mut mapped_file_arc).ok_or(AconfigStorageError::ObtainMappedFileFail(
                anyhow!("fail to get underlying mapped file"),
            ))?;
        update_boolean_flag_value(&mut mapped_file[..], offset, value)?;
        mapped_file.flush().map_err(|errmsg| {
            AconfigStorageError::MapFlushFail(anyhow!("fail to flush storage file: {}", errmsg))
        })?
    }
    Ok(())
}

/// Set the boolean flag value
///
/// This function would map the corresponding flag value file if has not been mapped yet,
/// and then write the target flag value at the specified offset.
///
/// If flag value is successfully set, it returns Ok(()), otherwise it returns the error.
pub unsafe fn set_boolean_flag_value(
    container: &str,
    offset: u32,
    value: bool,
) -> Result<(), AconfigStorageError> {
    set_boolean_flag_value_impl(STORAGE_LOCATION_FILE, container, offset, value)
}
