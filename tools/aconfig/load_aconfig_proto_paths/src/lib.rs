use anyhow::Result;
use std::path::Path;

use std::fs;

pub fn load_paths() -> Result<Vec<String>> {
    let mut result: Vec<String> =
        include!("../partition_aconfig_flags_paths.txt").map(|s| s.to_string()).to_vec();
    for dir in fs::read_dir("/apex")? {
        let dir = dir?;

        // Only scan the currently active version of each mainline module; skip the @version dirs.
        if dir.file_name().as_encoded_bytes().iter().any(|&b| b == b'@') {
            continue;
        }

        let path = format!("/apex/{}/etc/aconfig_flags.pb", dir.file_name().to_string_lossy());
        if Path::new(&path).exists() {
            result.push(path);
        }
    }

    Ok(result)
}
