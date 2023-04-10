//! Something

use std::env;
use std::process;

pub struct Env {
    pub build_top: String,
    pub build_product: String,
    pub build_variant: String,
    pub build_unbundled: String,
    pub product_out: String,
}

pub fn get_required_env(name: &str) -> String {
    match env::var(name) {
        Ok(val) => val,
        Err(e) => {
            eprintln!("{name} not set. Did you source build/envsetup.sh, run lunch and \
                do a build?");
            process::exit(1);
        }
    }
}

pub fn get_optional_env(name: &str) -> String {
    match env::var(name) {
        Ok(val) => val,
        Err(e) => String::from(""),
    }
}

pub fn default() -> Env {
    Env {
        build_top: get_required_env("ANDROID_BUILD_TOP"),
        build_product: get_required_env("TARGET_PRODUCT"),
        build_variant: get_required_env("TARGET_BUILD_VARIANT"),
        build_unbundled: get_optional_env("TARGET_BUILD_APPS"),
        product_out: get_required_env("ANDROID_PRODUCT_OUT"),
    }
}
