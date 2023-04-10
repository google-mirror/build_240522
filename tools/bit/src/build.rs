//! Encapsulates all access to the build system

use crate::environ;
use serde::{Deserialize, Serialize};
use serde_json;
use std::collections::HashMap;
use std::io::BufReader;
use std::io::ErrorKind;
use std::fs::File;
use std::process::Command;

/// Something that can be built
#[derive(Debug)]
pub struct BuildTarget {
    name: String,
    class: Vec<String>,
    installed: Vec<String>,
}

/// The build system
pub trait Build {
    fn get_target(&self, name: &String) -> Option<BuildTarget>;
    fn run_build(&self, targets: &Vec<String>); // TODO: Return errors or something
}

/// A soong based build
pub struct SoongBuild<'a> {
    env: &'a environ::Env,
    module_info: ModuleInfo,
}

impl<'a> SoongBuild<'a> {
    /// Get a new soong build system
    pub fn new(env: &'a environ::Env) -> Result<Self, BuildInitError> {
        let module_info = match load_module_info(&format!("{}/module-info.json", env.product_out)) {
            Ok(val) => val,
            Err(err) => return Err(err),
        };

        Ok(SoongBuild {
            env: env,
            module_info: module_info,
        })
    }
}

impl<'a> Build for SoongBuild<'a> {
    /// foo
    fn get_target(&self, name: &String) -> Option<BuildTarget> {
        match self.module_info.modules.get(name) {
            Some(module) => Some(BuildTarget {
                name: name.clone(),
                class: module.class.to_owned(),
                installed: module.installed.to_owned(),
            }),
            None => None,
        }
    }

    /// foo
    fn run_build(&self, targets: &Vec<String>) {
        // soong_ui.bash
        let mut cmd = Command::new("/usr/bin/echo");
        cmd.arg("/build/soong/soong_ui.bash");
        cmd.arg("--build-mode");
        cmd.arg("all-modules");
        cmd.arg(format!("--dir=\"{}\"", self.env.build_top));

        // Build targets
        for target in targets {
            cmd.arg(target);
        }

        // Run it
        let status = cmd.status();
    }
}

fn load_module_info(filename: &String) -> Result<ModuleInfo, BuildInitError> {
    let file = match File::open(filename) {
        Ok(val) => val,
        Err(err) => match (err.kind()) {
            ErrorKind::NotFound => return Err(BuildInitError::Missing),
            _ => return Err(BuildInitError::Other(format!("Error opening {}: {}",
                                                          filename,
                                                          err.to_string()))),
        }
    };
    let reader = BufReader::new(file);
    
    let module_info: ModuleInfo = match serde_json::from_reader(reader) {
        Ok(val) => val,
        Err(err) => return Err(BuildInitError::Other(format!("Error parsing {}: {}",
                                                          filename,
                                                          err.to_string()))),
    };

    Ok(module_info)
}

#[derive(Debug, Deserialize, Serialize)]
struct ModuleInfo {
    #[serde(flatten)]
    modules: HashMap<String, BuildModule>,
}

#[derive(Debug, Deserialize, Serialize)]
struct BuildModule {
    class: Vec<String>,
    installed: Vec<String>,
}

pub enum BuildInitError {
    Missing,
    Other(String),
}

/*
  "AVFHostTestCases": {
      "class": ["JAVA_LIBRARIES"],
      "path": ["packages/modules/Virtualization/tests/benchmark_hostside"],
      "tags": ["tests"],
      "installed": ["out/host/linux-x86/framework/AVFHostTestCases.jar"],
      "compatibility_suites": ["general-tests"],
      "auto_test_config": [],
      "module_name": "AVFHostTestCases",
      "test_config": ["packages/modules/Virtualization/tests/benchmark_hostside/AndroidTest.xml"],
      "dependencies": [],
      "shared_libs": [],
      "system_shared_libs": ["none"],
      "srcs": [],
      "srcjars": [],
      "classes_jar": ["out/host/common/obj/JAVA_LIBRARIES/AVFHostTestCases_intermediates/classes.jar"],
      "test_mainline_modules": [],
      "is_unit_test": "",
      "test_options_tags": [],
      "data": [],
      "runtime_dependencies": [],
      "static_dependencies": [],
      "data_dependencies": [],
      "supported_variants": ["HOST"],
      "host_dependencies": [],
      "target_dependencies": [] },
*/
