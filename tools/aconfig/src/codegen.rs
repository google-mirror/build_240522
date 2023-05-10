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
use std::fmt::Debug;
use std::fs;

use crate::aconfig::{FlagState, Permission};
use crate::cache::Cache;

#[derive(Debug)]
pub struct CodeGenerator {
    cache: Cache,
    writer: Box<dyn CodegenWriter>,
}

impl CodeGenerator {
    pub fn new(cache: Cache, writer: Box<dyn CodegenWriter>) -> CodeGenerator {
        CodeGenerator { cache, writer }
    }

    pub fn generate_java_code(&self, out: &str) -> Result<()> {
        for item in self.cache.iter() {
            let mut default_value = "false".to_string();
            if item.state == FlagState::Enabled {
                default_value = "true".to_string();
            }

            let mut return_value = default_value.clone();
            if item.permission == Permission::ReadWrite {
                return_value = format!(
                    "DeviceConfig.getBoolean(\"{name_space}\", \"{feature_name}__{flag_name}\", {default_value})",
                    name_space=item.id,
                    feature_name=item.id,
                    flag_name=item.id,
                    default_value=default_value
                );
            }
            let file_content = format!(
                include_str!("template/java.tmpl"),
                flag_name = item.id.to_lowercase(),
                class_name = uppercase_first_letter(item.id.to_lowercase().as_str()),
                method_name = item.id.to_lowercase(),
                return_value = return_value
            );
            let file_name =
                format!("{}.java", uppercase_first_letter(item.id.to_lowercase().as_str()));
            self.writer.write_file(&file_content, &file_name, out)?;
        }
        Ok(())
    }

    pub fn generate_c_code(&self) -> Result<()> {
        unimplemented!("Coming soon!!");
    }

    pub fn generate_rust_code(&self) -> Result<()> {
        unimplemented!("Coming soon!!");
    }
}

pub trait CodegenWriter: Debug {
    fn write_file(&self, file_content: &str, file_name: &str, out: &str) -> Result<()>;
}

#[derive(Debug)]
pub struct DevWriter;

impl CodegenWriter for DevWriter {
    fn write_file(&self, file_content: &str, file_name: &str, out: &str) -> Result<()> {
        let dst_path = format!("{}/{}", out, file_name);
        fs::write(dst_path, file_content)?;
        Ok(())
    }
}

fn uppercase_first_letter(s: &str) -> String {
    let mut c_str = s.chars();
    match c_str.next() {
        None => String::new(),
        Some(c) => c.to_uppercase().collect::<String>() + c_str.as_str(),
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;
    use std::default::Default;

    use crate::aconfig::{Flag, Value};
    use crate::commands::Source;

    #[derive(Debug, Default)]
    struct TestWriter {
        pub expect_result: HashMap<String, String>,
    }

    impl CodegenWriter for TestWriter {
        fn write_file(&self, file_content: &str, file_name: &str, out: &str) -> Result<()> {
            let dst_path = format!("{}/{}", out, file_name);
            assert!(self.expect_result.contains_key(&dst_path));
            assert_eq!(self.expect_result.get(&dst_path).unwrap(), &file_content);
            Ok(())
        }
    }

    impl TestWriter {
        fn add_test_expect(&mut self, file_content: &str, file_name: &str, out: &str) {
            let file_name = format!("{}/{}", out, file_name);
            self.expect_result.insert(file_name, file_content.to_string());
        }
    }

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
        let mut writer: Box<TestWriter> = Box::new(Default::default());
        let test1_expect_content = format!(
            include_str!("template/java.tmpl"),
            flag_name = "test1",
            class_name = "Test1",
            method_name = "test1",
            return_value = "true"
        );
        let test2_expect_content = format!(
            include_str!("template/java.tmpl"),
            flag_name = "test2",
            class_name = "Test2",
            method_name = "test2",
            return_value = "false"
        );
        writer.add_test_expect(test1_expect_content.as_str(), "Test1.java", "");
        writer.add_test_expect(test2_expect_content.as_str(), "Test2.java", "");
        let code_generator = CodeGenerator::new(cache, writer);
        code_generator.generate_java_code("").unwrap();
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
        let mut writer: Box<TestWriter> = Box::new(Default::default());
        let test3_expect_content = format!(
            include_str!("template/java.tmpl"),
            flag_name = "test3",
            class_name = "Test3",
            method_name = "test3",
            return_value = "DeviceConfig.getBoolean(\"test3\", \"test3__test3\", true)"
        );
        let test4_expect_content = format!(
            include_str!("template/java.tmpl"),
            flag_name = "test4",
            class_name = "Test4",
            method_name = "test4",
            return_value = "DeviceConfig.getBoolean(\"test4\", \"test4__test4\", false)"
        );
        writer.add_test_expect(test3_expect_content.as_str(), "Test3.java", "");
        writer.add_test_expect(test4_expect_content.as_str(), "Test4.java", "");
        let code_generator = CodeGenerator::new(cache, writer);
        code_generator.generate_java_code("").unwrap();
    }
}
