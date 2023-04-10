//! Something

mod build;
mod options;
mod environ;

use std::collections::HashSet;
use std::env;
use std::process;
use build::Build;
use build::BuildInitError;

fn main() {
    env::set_var("RUST_BACKTRACE", "1");

    let options = options::parse(&env::args().collect()).unwrap_or_else(|err|{
        eprintln!("Command line error: {}", err);
        eprintln!("Run 'bit -h' for options.");
        process::exit(1);
    });

    // Get the required environment variables. This will fail and exit if
    // they aren't set.
    let env = environ::default();

    // cd to the root of the source tree, so we never have to worry about
    // current directory again.
    env::set_current_dir(&env.build_top);

    if false {
        println!("-----------------------------------------------------------");
        println!("options: {options:#?}");
        println!("ANDROID_BUILD_TOP={}", env.build_top);
        println!("TARGET_PRODUCT={}", env.build_product);
        println!("TARGET_BUILD_VARIANT={}", env.build_variant);
        println!("TARGET_BUILD_APPS={}", env.build_unbundled);
        println!("ANDROID_PRODUCT_OUT={}", env.product_out);
        println!("-----------------------------------------------------------");
    }

    // Initialize the build system. That includes loading module-info.json, and if
    // that doesn't work, fail with a helpful error message.
    let build = match build::SoongBuild::new(&env) {
        Ok(b) => b,
        Err(err) => match err {
            BuildInitError::Missing => {
                eprintln!("Could not load module-info.json. Did you run a full build first?");
                process::exit(1);
            },
            BuildInitError::Other(msg) => {
                eprintln!("Unexpected error: {}", msg);
                process::exit(1);
            }
        }
    };

    // Validate that the command line provided modules are in module info.
    let mut stop = false;
    for arg in options.targets.iter() {
        if build.get_target(&arg.name).is_none() {
            stop = true;
            eprintln!("Unknown build target: {}", arg.name);
        }
    }
    if stop {
        process::exit(1);
    }

    // Choose which modules to build
    let mut modules = HashSet::<String>::new();
    for arg in options.targets.iter() {
        // Add the requested targets
        modules.insert(arg.name.clone());

        // Add the test dep targets
        if arg.install_test_deps {
            // TODO: Implemented tested_by in soong
            // let build_target = build.get_target(arg.name).unwrap();
            // for test_dep_name in build_target.test_deps {
            //     modules.insert(test_dep_name);
            // }
        }
    }

    let mut target_vec = modules.into_iter().collect::<Vec<_>>();
    target_vec.sort();

    build.run_build(&target_vec);
}

