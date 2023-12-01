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

use anyhow::Result;
use std::collections::{HashMap, HashSet, hash_map::DefaultHasher};
use std::hash::{Hash, Hasher};
use std::path::PathBuf;

use crate::commands::OutputFile;
use crate::package_table::PackageTable;
use crate::protos::{ProtoParsedFlags, ProtoParsedFlag};

pub const FILE_VERSION: u32 = 1;

pub const HASH_PRIMES: [u32; 29] =
    [7, 13, 29, 53, 97, 193, 389, 769, 1543, 3079, 6151, 12289, 24593,
     49157, 98317, 196613, 393241, 786433, 1572869, 3145739, 6291469,
     12582917, 25165843, 50331653, 100663319, 201326611, 402653189,
     805306457, 1610612741];

pub fn get_table_size(entries : u32) -> Option<u32> {
    for num in HASH_PRIMES {
        if num >= 2*entries {
            return Some(num);
        }
    }
    None
}

pub fn get_bucket_index<T: Hash>(val: &T) -> u64 {
    let mut s = DefaultHasher::new();
    val.hash(&mut s);
    s.finish()
}

pub struct FlagPackage<'a> {
    pub package_name: &'a str,
    pub package_id: u32,
    pub flag_names: HashSet<&'a str>,
    pub boolean_flags: Vec<&'a ProtoParsedFlag>,
    pub boolean_offset: u32,
}

impl<'a> FlagPackage<'a> {
    fn new(package_name: &'a str, package_id : u32) -> Self {
        FlagPackage {
            package_name,
            package_id,
            flag_names: HashSet::new(),
            boolean_flags: vec![],
            boolean_offset: 0,
        }
    }

    fn insert(&mut self, pf: &'a ProtoParsedFlag) {
        if self.flag_names.insert(pf.name()) {
            self.boolean_flags.push(pf);
        }
    }
}

pub fn sort_flags<'a, I>(
    parsed_flags_vec_iter: I
) -> Vec<FlagPackage<'a>>
where
    I: Iterator<Item = &'a ProtoParsedFlags>,
{
    // group flags by package
    let packages: Vec<FlagPackage<'a>> = Vec::new();
    let package_index: HashMap<&'a str, usize> = HashMap::new();
    for parsed_flags in parsed_flags_vec_iter {
        for parsed_flag in parsed_flags.parsed_flag.iter() {
            let index = *(package_index.entry(parsed_flag.package()).or_insert(packages.len()));
            if index == packages.len() {
                packages.push(FlagPackage::new(parsed_flag.package(), index as u32));
            }
            packages[index].insert(parsed_flag);
        }
    }

    // calculate package flag value start offset
    let boolean_offset = 0;
    for p in packages {
        p.boolean_offset = boolean_offset;
        boolean_offset += 2 * p.boolean_flags.len() as u32;
    }

    packages
}


pub fn generate_storage_files<'a, I>(
    container: &str,
    parsed_flags_vec_iter: I,
) -> Result<Vec<OutputFile>>
where
    I: Iterator<Item = &'a ProtoParsedFlags>,
{
    let packages = sort_flags(parsed_flags_vec_iter);

    let package_table = PackageTable::new(container, packages);
    let package_table_file_path = PathBuf::from(r"package.map");
    let package_table_file = OutputFile{
        contents: package_table.as_bytes(),
        path: package_table_file_path,
    };

    Ok(vec![package_table_file])
}
