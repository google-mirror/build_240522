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

use serde::Serialize;

use crate::aconfig::{FlagState, Permission};
use crate::cache::{Cache, Item};

pub struct GeneratedFile {
    pub file_content: String,
    pub file_name: String,
}

pub fn generate_codegen_context(cache: &Cache) -> Context {
    let class_elements: Vec<ClassElement> = cache.iter().map(create_class_element).collect();
    let readwrite = class_elements.iter().any(|item| item.readwrite);
    let namespace = cache.iter()
                         .find(|item| !item.namespace.is_empty())
                         .unwrap().namespace.to_lowercase();
    Context { namespace: namespace.clone(), readwrite, class_elements }
}

#[derive(Serialize)]
pub struct Context {
    pub namespace: String,
    pub readwrite: bool,
    pub class_elements: Vec<ClassElement>,
}

#[derive(Serialize)]
pub struct ClassElement {
    pub method_name: String,
    pub readwrite: bool,
    pub default_value: String,
    pub feature_name: String,
    pub flag_name: String,
}

fn create_class_element(item: &Item) -> ClassElement {
    ClassElement {
        method_name: item.name.clone(),
        readwrite: item.permission == Permission::ReadWrite,
        default_value: if item.state == FlagState::Enabled {
            "true".to_string()
        } else {
            "false".to_string()
        },
        feature_name: item.name.clone(),
        flag_name: item.name.clone(),
    }
}

