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

use crate::create_storage::{self, FlagPackage};

pub struct PackageTableHeader<'a> {
    pub version: u32,
    pub container: &'a str,
    pub file_size: u32,
    pub num_packages: u32,
    pub bucket_offset: u32,
    pub node_offset: u32,
}

impl<'a> PackageTableHeader<'a> {
    fn new(container: &'a str, num_packages: u32) -> Self {
        PackageTableHeader {
            version: create_storage::FILE_VERSION,
            container,
            file_size: 0,
            num_packages,
            bucket_offset: 0,
            node_offset: 0,
        }
    }

    fn as_bytes(&self) -> Vec<u8> {
        let result = Vec::new();
        result.extend_from_slice(&self.version.to_be_bytes());
        let container_bytes = self.container.as_bytes();
        result.extend_from_slice(&(container_bytes.len() as u32).to_be_bytes());
        result.extend_from_slice(container_bytes);
        result.extend_from_slice(&self.file_size.to_be_bytes());
        result.extend_from_slice(&self.num_packages.to_be_bytes());
        result.extend_from_slice(&self.bucket_offset.to_be_bytes());
        result.extend_from_slice(&self.node_offset.to_be_bytes());
        result
    }
}

pub struct PackageTableNode<'a> {
    pub package_name: &'a str,
    pub package_id: u32,
    pub boolean_offset: u32,
    pub next_offset: u32,
    pub start_offset: u32,
    pub bucket_index: u32,
}

impl<'a> PackageTableNode<'a> {
    fn new(package: &'a FlagPackage, num_buckets: u32) -> Self {
        let bucket_index =
            create_storage::get_bucket_index(&package.package_name.to_string(), num_buckets) as u32;
        PackageTableNode {
            package_name: package.package_name,
            package_id: package.package_id,
            boolean_offset: package.boolean_offset,
            next_offset: 0,
            start_offset: 0,
            bucket_index,
        }
    }

    fn as_bytes(&self) -> Vec<u8> {
        let result = Vec::new();
        let name_bytes = self.package_name.as_bytes();
        result.extend_from_slice(&(name_bytes.len() as u32).to_be_bytes());
        result.extend_from_slice(name_bytes);
        result.extend_from_slice(&self.package_id.to_be_bytes());
        result.extend_from_slice(&self.boolean_offset.to_be_bytes());
        result.extend_from_slice(&self.next_offset.to_be_bytes());
        result
    }
}

pub struct PackageTable<'a> {
    pub header: PackageTableHeader<'a>,
    pub buckets: Vec<Option<u32>>,
    pub nodes: Vec<PackageTableNode<'a>>,
}

impl<'a> PackageTable<'a> {
    pub fn new(container: &'a str, packages: Vec<FlagPackage>) -> Self {
        // create table
        let num_packages = packages.len() as u32;
        let num_buckets =
            create_storage::get_table_size(num_packages).expect("Number of packages is too large");
        let table = PackageTable {
            header: PackageTableHeader::new(container, num_packages),
            buckets: vec![None; num_buckets as usize],
            nodes: packages.iter().map(|pkg| PackageTableNode::new(pkg, num_buckets)).collect(),
        };

        // sort nodes by bucket index for efficiency
        table.nodes.sort_by(|a, b| a.bucket_index.cmp(&b.bucket_index));

        // fill all node offset
        let offset = 0;
        for (i, node) in table.nodes.iter().enumerate() {
            node.start_offset = offset;
            offset += node.as_bytes().len() as u32;
            if i + 1 < table.nodes.len() && node.bucket_index == table.nodes[i + 1].bucket_index {
                node.next_offset = offset;
            }
            if table.buckets[node.bucket_index as usize].is_none() {
                table.buckets[node.bucket_index as usize] = Some(node.start_offset);
            }
        }

        // fill table region offset
        table.header.bucket_offset = table.header.as_bytes().len() as u32;
        table.header.node_offset = table.header.bucket_offset + num_buckets * 4;
        table.header.file_size = table.header.node_offset
            + table.nodes.iter().map(|x| x.as_bytes().len()).sum::<usize>() as u32;

        table
    }

    pub fn as_bytes(&self) -> Vec<u8> {
        vec![
            self.header.as_bytes(),
            self.buckets.iter().map(|v| v.unwrap_or(0).to_be_bytes()).collect::<Vec<_>>().concat(),
            self.nodes.iter().map(|v| v.as_bytes()).collect::<Vec<_>>().concat(),
        ]
        .concat()
    }
}
