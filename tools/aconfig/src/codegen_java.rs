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
use tinytemplate::TinyTemplate;

use crate::cache::Cache;
use crate::codegen_context::{generate_codegen_context, GeneratedFile};

pub fn generate_java_code(cache: &Cache) -> Result<GeneratedFile> {
    let mut context = generate_codegen_context(cache);
    context.namespace = uppercase_first_letter(&context.namespace);
    let mut template = TinyTemplate::new();
    template.add_template("java_code_gen", include_str!("../templates/java.template"))?;
    let file_content = template.render("java_code_gen", &context)?;
    Ok(GeneratedFile { file_content, file_name: format!("{}.java", context.namespace) })
}

fn uppercase_first_letter(s: &str) -> String {
    s[0..1].to_uppercase() + &s[1..]
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::aconfig::{Flag, FlagState, Permission, Value};
    use crate::commands::Source;

    #[test]
    fn test_generate_java_code() {
        let namespace = "TeSTFlaG";
        let mut cache = Cache::new(1, namespace.to_string());
        cache
            .add_flag(
                Source::File("test.txt".to_string()),
                Flag {
                    name: "test".to_string(),
                    description: "buildtime enable".to_string(),
                    values: vec![Value::default(FlagState::Enabled, Permission::ReadOnly)],
                },
            )
            .unwrap();
        cache
            .add_flag(
                Source::File("test2.txt".to_string()),
                Flag {
                    name: "test2".to_string(),
                    description: "runtime disable".to_string(),
                    values: vec![Value::default(FlagState::Disabled, Permission::ReadWrite)],
                },
            )
            .unwrap();
        let expect_content = "package com.android.aconfig;

        import android.provider.DeviceConfig;

        public final class Testflag {

            public static boolean test() {
                return true;
            }

            public static boolean test2() {
                return DeviceConfig.getBoolean(
                    \"Testflag\",
                    \"test2__test2\",
                    false
                );
            }

        }
        ";
        let expected_file_name = String::from("Testflag.java");
        let generated_file = generate_java_code(&cache).unwrap();
        assert_eq!(expected_file_name, generated_file.file_name);
        assert_eq!(expect_content.replace(' ', ""), generated_file.file_content.replace(' ', ""));
    }
}
