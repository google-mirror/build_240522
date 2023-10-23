//! `metadata` is a build time tool to generate metadata for the test_run rules.

use anyhow::Result;
use clap::{Arg, ArgAction, Command};

fn cli() -> Command {
    Command::new("metadata")
        .subcommand_required(true)
        .subcommand(
            Command::new("generate-metadata")
                .arg(Arg::new("inputFilenames").long("inputFilenames").action(ArgAction::Append))
                .arg(Arg::new("outputFilenames").long("outputFilenames").action(ArgAction::Append)))
}

fn main() -> Result<()> {
    let matches = cli().get_matches();
    match matches.subcommand() {
        Some(("metadata", sub_matches)) => {
            for path in sub_matches.get_many::<String>("inputFilenames").unwrap_or_default() {
                println!("Generating metadata for input files: {:?}", path);
            }
        }
        _ => unreachable!(),
    }
    Ok(())
}