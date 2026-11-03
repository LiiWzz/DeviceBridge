use crate::config::{Config, LogLevel};
use std::ffi::OsString;
use std::path::PathBuf;

pub const HELP: &str = "DeviceBridge agent — M01 foundation (no transport)\n\
Usage: devicebridge-agent [OPTIONS]\n\
  --config PATH             Optional JSON file, at most 4096 bytes\n\
  --log-level LEVEL         error, info (default), debug\n\
  --shutdown-poll-ms MS      Shutdown polling interval, 10..=1000 (default 100)\n\
  --check-config            Validate configuration and exit\n\
  --help                    Print this help and exit\n\
  --version                 Print version and exit\n\
CLI values override file values; no environment configuration is read.\n";

#[derive(Debug, PartialEq, Eq)]
pub enum Action {
    Run,
    Check,
    Help,
    Version,
}

#[derive(Debug)]
pub struct Cli {
    pub action: Action,
    config_path: Option<PathBuf>,
    log_level: Option<LogLevel>,
    shutdown_poll_ms: Option<u64>,
}

impl Cli {
    pub fn parse(args: impl IntoIterator<Item = OsString>) -> Result<Self, String> {
        let mut cli = Self {
            action: Action::Run,
            config_path: None,
            log_level: None,
            shutdown_poll_ms: None,
        };
        let mut args = args.into_iter();
        let mut count = 0;
        let mut seen = 0_u8;
        while let Some(arg) = args.next() {
            count += 1;
            if count > 10 || arg.len() > 1024 {
                return Err("too many or oversized CLI arguments".into());
            }
            let flag = arg.to_str().ok_or("option must be UTF-8")?;
            let (bit, needs_value) = match flag {
                "--config" => (1, true),
                "--log-level" => (2, true),
                "--shutdown-poll-ms" => (4, true),
                "--check-config" => (8, false),
                "--help" => (16, false),
                "--version" => (32, false),
                _ => return Err("unknown option; use --help".into()),
            };
            if seen & bit != 0 {
                return Err("duplicate option".into());
            }
            seen |= bit;
            if needs_value {
                let value = args.next().ok_or("option requires a value")?;
                count += 1;
                if value.is_empty() || value.len() > 1024 {
                    return Err("empty or oversized option value".into());
                }
                match flag {
                    "--config" => cli.config_path = Some(value.into()),
                    "--log-level" => {
                        cli.log_level =
                            Some(LogLevel::parse(value.to_str().ok_or("invalid level")?)?);
                    }
                    _ => {
                        let text = value.to_str().ok_or("invalid polling interval")?;
                        if !text.bytes().all(|byte| byte.is_ascii_digit()) {
                            return Err("polling interval must be an unsigned integer".into());
                        }
                        cli.shutdown_poll_ms =
                            Some(text.parse().map_err(|_| "invalid polling interval")?);
                    }
                }
            } else {
                if cli.action != Action::Run {
                    return Err("choose only one of --help, --version, --check-config".into());
                }
                cli.action = match flag {
                    "--check-config" => Action::Check,
                    "--help" => Action::Help,
                    _ => Action::Version,
                };
            }
        }
        // Validate CLI values even for informational invocations.
        if let Some(shutdown_poll_ms) = cli.shutdown_poll_ms {
            Config {
                shutdown_poll_ms,
                ..Config::default()
            }
            .validate()?;
        }
        Ok(cli)
    }

    pub fn config(&self) -> Result<Config, String> {
        let mut config = match &self.config_path {
            Some(path) => Config::load(path)?,
            None => Config::default(),
        };
        if let Some(level) = self.log_level {
            config.log_level = level;
        }
        if let Some(interval) = self.shutdown_poll_ms {
            config.shutdown_poll_ms = interval;
        }
        config.validate()?;
        Ok(config)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    fn parse(args: &[&str]) -> Result<Cli, String> {
        Cli::parse(args.iter().map(OsString::from))
    }

    #[test]
    fn defaults_actions_and_overrides() {
        assert_eq!(parse(&[]).unwrap().config().unwrap(), Config::default());
        assert_eq!(parse(&["--help"]).unwrap().action, Action::Help);
        assert_eq!(parse(&["--version"]).unwrap().action, Action::Version);
        let cli = parse(&[
            "--check-config",
            "--log-level",
            "debug",
            "--shutdown-poll-ms",
            "10",
        ])
        .unwrap();
        assert_eq!(cli.action, Action::Check);
        assert_eq!(
            cli.config().unwrap(),
            Config {
                log_level: LogLevel::Debug,
                shutdown_poll_ms: 10
            }
        );
    }

    #[test]
    fn file_values_are_overridden_by_cli() {
        let cli = parse(&[
            "--config",
            concat!(env!("CARGO_MANIFEST_DIR"), "/config.example.json"),
            "--log-level",
            "debug",
            "--shutdown-poll-ms",
            "10",
        ])
        .unwrap();
        assert_eq!(
            cli.config().unwrap(),
            Config {
                log_level: LogLevel::Debug,
                shutdown_poll_ms: 10
            }
        );
    }

    #[cfg(unix)]
    #[test]
    fn invalid_utf8_options_and_values_are_errors() {
        use std::os::unix::ffi::OsStringExt;
        let invalid = OsString::from_vec(vec![0xff]);
        assert!(Cli::parse([invalid.clone()]).is_err());
        assert!(Cli::parse([OsString::from("--log-level"), invalid]).is_err());
    }

    #[test]
    fn rejects_bad_duplicate_or_conflicting_arguments() {
        for args in [
            vec!["--wat"],
            vec!["positional"],
            vec!["--config"],
            vec!["--config", ""],
            vec!["--log-level", "trace"],
            vec!["--log-level", "info", "--log-level", "debug"],
            vec!["--check-config", "--check-config"],
            vec!["--help", "--version"],
            vec!["--shutdown-poll-ms", "-1"],
            vec!["--shutdown-poll-ms", "+10"],
            vec!["--shutdown-poll-ms", "0"],
            vec!["--shutdown-poll-ms", "1001"],
            vec!["--shutdown-poll-ms", "18446744073709551616"],
            vec!["--shutdown-poll-ms", "abc"],
        ] {
            assert!(parse(&args).is_err(), "accepted {args:?}");
        }
        assert!(parse(&["--config", &"x".repeat(1025)]).is_err());
    }
}
