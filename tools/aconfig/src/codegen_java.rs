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
use crate::aconfig::{FlagState, Permission};
use crate::cache::Cache;

pub struct GeneratedFile{
    pub file_content: String,
    pub file_name: String,
}

pub fn generate_java_code(cache: Cache) -> Vec<GeneratedFile> {
    let mut files = vec![];
    for item in cache.iter() {
        let flag_name = item.id.to_ascii_lowercase();
        let class_name = uppercase_first_letter(item.id.to_ascii_lowercase().as_str());
        let method_name = item.id.to_ascii_lowercase();
        let namespace = item.id.to_ascii_lowercase();
        let feature_name = item.id.to_ascii_lowercase();
        let default_value = match item.state {
            FlagState::Enabled => "true",
            FlagState::Disabled => "false",
        };

        let file_content = if item.permission == Permission::ReadOnly {
            java_buildtime_flag_code_gen(&flag_name, &class_name, &method_name, default_value)
        } else {
            java_runtime_flag_code_gen(&flag_name, &class_name, &method_name, &namespace, &feature_name,default_value)
        };

        let file_name =
            format!("{}.java", uppercase_first_letter(item.id.to_ascii_lowercase().as_str()));
        
        files.push(GeneratedFile{file_content, file_name});
    }
    return files;
}

fn java_buildtime_flag_code_gen(
    flag_name: &str,
    class_name: &str,
    method_name: &str,
    default_value: &str) -> String {
    format!(
        include_str!("../templates/java_buildtime.template"),
        flag_name=flag_name,
        class_name=class_name,
        method_name=method_name,
        default_value=default_value
    )
}

fn java_runtime_flag_code_gen(
    flag_name: &str,
    class_name: &str,
    method_name: &str,
    namespace: &str,
    feature_name: &str,
    default_value: &str) -> String {
    format!(
        include_str!("../templates/java_runtime.template"),
        flag_name=flag_name,
        class_name=class_name,
        method_name=method_name,
        namespace=namespace,
        feature_name=feature_name,
        default_value=default_value
    )
}

fn uppercase_first_letter(s: &str) -> String {
    s.chars().enumerate().map(|(index, ch)| {
        if index == 0 {
            ch.to_ascii_uppercase()
        } else {
            ch
        }
    }).collect()
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::aconfig::{Flag, Value};
    use crate::commands::Source;

    #[test]
    fn test_generate_java_code_buildtime() {
        let mut cache = Cache::new(1);
        cache
            .add_flag(
                Source::File("test1.txt".to_string()),
                Flag {
                    id: "test1".to_string(),
                    description: "buildtime enable".to_string(),
                    values: vec![Value::default(FlagState::Enabled, Permission::ReadOnly)],
                },
            )
            .unwrap();
        cache
            .add_flag(
                Source::File("test2.txt".to_string()),
                Flag {
                    id: "test2".to_string(),
                    description: "buildtime disable".to_string(),
                    values: vec![Value::default(FlagState::Disabled, Permission::ReadOnly)],
                },
            )
            .unwrap();
        let test1_expect_content = format!(
            include_str!("../templates/java_buildtime.template"),
            flag_name = "test1",
            class_name = "Test1",
            method_name = "test1",
            default_value = "true"
        );
        let test2_expect_content = format!(
            include_str!("../templates/java_buildtime.template"),
            flag_name = "test2",
            class_name = "Test2",
            method_name = "test2",
            default_value = "false"
        );
        let generated_files = generate_java_code(cache);
        for file in generated_files {
            if file.file_name == "Test1.java" {
                assert_eq!(test1_expect_content, file.file_content);
            } else {
                assert_eq!(test2_expect_content, file.file_content);
            }
        }
    }

    #[test]
    fn test_generate_java_code_runtime() {
        let mut cache = Cache::new(1);
        cache
            .add_flag(
                Source::File("test3.txt".to_string()),
                Flag {
                    id: "test3".to_string(),
                    description: "runtime enable".to_string(),
                    values: vec![Value::default(FlagState::Enabled, Permission::ReadWrite)],
                },
            )
            .unwrap();
        cache
            .add_flag(
                Source::File("test4.txt".to_string()),
                Flag {
                    id: "test4".to_string(),
                    description: "runtime disable".to_string(),
                    values: vec![Value::default(FlagState::Disabled, Permission::ReadWrite)],
                },
            )
            .unwrap();
        let test3_expect_content = format!(
            include_str!("../templates/java_runtime.template"),
            flag_name = "test3",
            class_name = "Test3",
            method_name = "test3",
            namespace = "test3",
            feature_name = "test3",
            default_value = "true"
        );
        let test4_expect_content = format!(
            include_str!("../templates/java_runtime.template"),
            flag_name = "test4",
            class_name = "Test4",
            method_name = "test4",
            namespace = "test4",
            feature_name = "test4",
            default_value = "false"
        );
        let generated_files = generate_java_code(cache);
        for file in generated_files {
            if file.file_name == "Test3.java" {
                assert_eq!(test3_expect_content, file.file_content);
            } else {
                assert_eq!(test4_expect_content, file.file_content);
            }
        }
    }
}
