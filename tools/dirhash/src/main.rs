//! Something


use ring::digest::{Context, Digest, SHA256};
use std::collections::HashSet;
use std::env;
use std::fs;
use std::process;
use walkdir::WalkDir;
use std::io::{BufReader, Error, Read, Write};
use std::path::{Path, PathBuf};
use rayon::prelude::*;

use fs::File;

struct Options {
    help: bool,
    dirs: Vec<String>,
}

fn to_text(buf: &[u8]) -> String {
    let mut result = String::new();
    result.reserve_exact(buf.len());
    for c in buf.iter() {
        let upper = c >> 4;
        result.push(char::from_digit((c >> 4).into(), 16).unwrap());
        result.push(char::from_digit((c & 0x7).into(), 16).unwrap());
    }
    return result;
}

fn parse_options(args: &Vec<String>) -> Result<Options, String> {
    let mut result = Options {
        help: false,
        dirs: Vec::new(),
    };

    let mut arg_index = 1;
    while arg_index < args.len() {
        let arg = &args[arg_index];
        arg_index += 1;
        
        if (arg == "-h" || arg == "--help") {
            result.help = true;
            return Ok(result);
        } else if (arg.as_bytes()[0] == b'-') {
            return Err(format!("unknown argument \"{}\"", arg).to_string())
        } else {
            result.dirs.push(arg.to_string());
        }
    }

    return Ok(result);
}

fn hash_file(filename: &Path) -> Result<String, Error> {
    let input = File::open(filename)?;
    let mut reader = BufReader::new(input);

    let mut context = Context::new(&SHA256);
    let mut buffer = [0; 1024];

    loop {
        let count = reader.read(&mut buffer)?;
        if count == 0 {
            break;
        }
        context.update(&buffer[..count]);
    }

    return Ok(to_text(context.finish().as_ref()))
}

#[derive(Debug)]
struct EntryInfo {
    filename: String,
    symlink: String,
    hash: String,
}


fn process_file(filename: &Path) -> Result<Option<EntryInfo>, Error> {
    let metadata = fs::symlink_metadata(filename)?;
    if metadata.is_dir() {
        return Ok(None)
    } else if metadata.is_symlink() {
        return Ok(Some(EntryInfo {
            filename: String::from(filename.to_str().unwrap()),
            symlink: String::from(fs::read_link(filename).unwrap().to_str().unwrap()),
            hash: String::from(""),
        }));
    } else {
        return Ok(Some(EntryInfo {
            filename: String::from(filename.to_str().unwrap()),
            symlink: String::from(""),
            hash: hash_file(filename).unwrap(),
        }));
    }
}

fn main() {
    env::set_var("RUST_BACKTRACE", "1");

    let options = parse_options(&env::args().collect()).unwrap_or_else(|err|{
        eprintln!("Command line error: {}", err);
        eprintln!("Run 'bit -h' for options.");
        process::exit(1);
    });

    if options.help {
        println!("usage: dirhash DIRS...");
        return;
    }

    let filenames: Vec<PathBuf> = options.dirs.iter()
            .flat_map(|dir| WalkDir::new(dir))
            .map(|result| result.unwrap().path().to_path_buf())
            .collect();

    // TODO: Understand rust traits well enough to know why into_par_iter() can't
    // take the iterator from the map() above.  Anyway, the directory walk isn't
    // that slow, so not parallelizing that isn't such a big deal.

    let infos: Vec<EntryInfo> = filenames
            .into_par_iter()
            .filter_map(|entry| { process_file(&entry).unwrap() })
            .collect();

    for info in infos {
        println!("  {:?}", info);
    }
}

