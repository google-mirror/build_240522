//! Command line options

/// Individual target (aka soong module) specified on the command line.
#[derive(Clone, Debug, PartialEq, PartialOrd)]
pub struct ArgTarget {
    pub name: String,
    pub install: bool,
    pub install_build_deps: bool,
    pub install_test_deps: bool,
    pub reboot: bool,
}

// When adding more optins, add them here, in the beginning of parse(), and
// in test::new_options. The defaults in both places must match for the tests
// to pass.
/// Arguments parsed from command line.
#[derive(Debug, PartialEq, PartialOrd)]
pub struct Options {
    pub disable_verity: bool,
    pub dry_run: bool,
    pub help: bool,
    pub ignore_installed_check: bool,
    pub targets: Vec<ArgTarget>,
}

/// Parse arguments and return a new Options, or None.
pub fn parse(args: &Vec<String>) -> Result<Options, String> {
    let mut result = Options {
        disable_verity: true,
        dry_run: false,
        help: false,
        ignore_installed_check: false,
        targets: Vec::new(),
    };

    // Skip first argument
    let mut arg_index = 1;
    let mut install = true;
    let mut install_build_deps = true;
    let mut install_test_deps = true;
    let mut reboot = true;

    while arg_index < args.len() {
        let arg = &args[arg_index];
        arg_index += 1;

        if arg.len() == 0 {
            return Err("unknown empty argument: \"\"".to_string());
        }

        let bytes = arg.as_bytes();
        if bytes[0] == b'-' || bytes[0] == b'+' {
            // '+' --> adding == true
            // '-' --> adding == false
            let mut adding = false;

            // Arg is just '-' or '+'
            if bytes.len() < 2 {
                return Err(format!("unknown argument (+/- missing flag): {}", arg).to_string());
            }

            // Long args
            if bytes[1] == b'-' {
                if arg == "--dry-run" {
                    result.dry_run = true;
                } else if arg == "--help" {
                    result.help = true;
                } else if arg == "--ignore-installed-check" {
                    result.ignore_installed_check = true;
                } else if arg == "--install" {
                    install = true;
                } else if arg == "--install-deps" {
                    install_build_deps = true;
                } else if arg == "--install-test-deps" {
                    install_test_deps = true;
                } else if arg == "--no-disable-verity" {
                    result.disable_verity = false;
                } else if arg == "--no-install" {
                    install = false;
                } else if arg == "--no-install-deps" {
                    install_build_deps = false;
                } else if arg == "--no-install-test-deps" {
                    install_test_deps = false;
                } else {
                    return Err(format!("unknown argument: {}", arg).to_string())
                }
                continue;
            }

            // Short args
            for b in bytes {
                if *b == b'+' {
                    adding = true;
                } else if *b == b'-' {
                    adding = false;
                } else if *b == b'd' {
                    install_build_deps = adding;
                } else if *b == b'i' {
                    install = adding;
                } else if *b == b'r' {
                    reboot = adding;
                } else if *b == b't' {
                    install_test_deps = adding;
                } else {
                    return Err(format!("unknown flags '{}' in argument: {}", b, arg).to_string())
                }
            }
        } else {
            // Arg is a target
            result.targets.push(ArgTarget{
                name: arg.clone(),
                install: install,
                install_build_deps: install_build_deps && install,
                install_test_deps: install_test_deps && install,
                reboot: reboot,
            });
        }

    }

    if result.targets.len() == 0 {
        return Err(String::from("No targets specified."));
    }

    Ok(result)
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    // Note: defaults here match default behavior of parse
    fn new_arg_target(name: &str) -> ArgTarget {
        return ArgTarget {
            name: String::from(name),
            install: true,
            install_build_deps: true,
            install_test_deps: true,
            reboot: true,
        }
    }

    fn new_options(arg_targets: &[&str]) -> Options {
        Options {
            disable_verity: true,
            dry_run: false,
            help: false,
            ignore_installed_check: false,
            targets: arg_targets.iter().map(|&name| new_arg_target(name)).collect(),
        }
    }

    fn wrap_parse_ok(args: &[&str]) -> Options {
        let mut v = args.iter().map(|&a| String::from(a)).collect::<Vec<_>>();
        v.insert(0, String::from("_unused_"));
        return parse(&v).expect("expected parse to return Ok");
    }

    #[test]
    fn test_parse_no_zero_arg() {
        let actual = parse(&Vec::new()).expect("parse failed");

        let expected = new_options(&[]);
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_parse_no_args() {
        let actual = wrap_parse_ok(&[]);

        let expected = new_options(&[]);
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_parse_single_module() {
        let actual = wrap_parse_ok(&["Module"]);

        let expected = new_options(&["Module"]);
        assert_eq!(actual, expected);
    }

    #[test]
    fn test_flags() {
        let actual = wrap_parse_ok(&["Zero", "-i", "One", "+id", "Two", "-r", "Three",
                "+r-t", "Four", "Five", "-d", "Six"]);

        let expected = Options {
            disable_verity: true,
            dry_run: false,
            help: false,
            ignore_installed_check: false,
            targets: vec![
                ArgTarget {
                    name: String::from("Zero"),
                    install: true,
                    install_build_deps: true,
                    install_test_deps: true,
                    reboot: true,
                },
                ArgTarget {
                    name: String::from("One"),
                    install: false,
                    install_build_deps: false,
                    install_test_deps: false,
                    reboot: false,
                },
                ArgTarget {
                    name: String::from("Two"),
                    install: true,
                    install_build_deps: true,
                    install_test_deps: true,
                    reboot: true,
                },
                ArgTarget {
                    name: String::from("Three"),
                    install: true,
                    install_build_deps: true,
                    install_test_deps: true,
                    reboot: false,
                },
                ArgTarget {
                    name: String::from("Four"),
                    install: true,
                    install_build_deps: true,
                    install_test_deps: false,
                    reboot: true,
                },
                ArgTarget {
                    name: String::from("Five"),
                    install: true,
                    install_build_deps: true,
                    install_test_deps: false,
                    reboot: true,
                },
                ArgTarget {
                    name: String::from("Six"),
                    install: true,
                    install_build_deps: false,
                    install_test_deps: false,
                    reboot: true,
                },

            ],
        };
        assert_eq!(actual, expected);
    }
}
