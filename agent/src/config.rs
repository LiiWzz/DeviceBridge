use serde::{Deserialize, Serialize};
use std::fs::OpenOptions;
use std::io::Read;
use std::path::Path;

pub const MAX_CONFIG_BYTES: usize = 4096;
pub const MIN_POLL_MS: u64 = 10;
pub const MAX_POLL_MS: u64 = 1000;

#[derive(Clone, Copy, Debug, Deserialize, Serialize, PartialEq, Eq, PartialOrd, Ord)]
#[serde(rename_all = "lowercase")]
pub enum LogLevel {
    Error,
    Info,
    Debug,
}

impl LogLevel {
    pub fn parse(value: &str) -> Result<Self, String> {
        match value {
            "error" => Ok(Self::Error),
            "info" => Ok(Self::Info),
            "debug" => Ok(Self::Debug),
            _ => Err("log level must be error, info or debug".into()),
        }
    }
}

#[derive(Debug, Deserialize, PartialEq, Eq)]
#[serde(default, deny_unknown_fields)]
pub struct Config {
    pub log_level: LogLevel,
    pub shutdown_poll_ms: u64,
}

impl Default for Config {
    fn default() -> Self {
        Self {
            log_level: LogLevel::Info,
            shutdown_poll_ms: 100,
        }
    }
}

impl Config {
    pub fn validate(&self) -> Result<(), String> {
        if !(MIN_POLL_MS..=MAX_POLL_MS).contains(&self.shutdown_poll_ms) {
            return Err(format!(
                "shutdown_poll_ms must be {MIN_POLL_MS}..={MAX_POLL_MS}"
            ));
        }
        Ok(())
    }

    pub fn load(path: &Path) -> Result<Self, String> {
        let mut options = OpenOptions::new();
        options.read(true);
        // Opening a FIFO must not block before the regular-file check.
        #[cfg(unix)]
        {
            use std::os::unix::fs::OpenOptionsExt;
            options.custom_flags(libc::O_NONBLOCK);
        }
        let file = options
            .open(path)
            .map_err(|err| format!("cannot open config: {err}"))?;
        // Reject FIFOs/devices/directories; no potentially unbounded stream reads.
        if !file
            .metadata()
            .map_err(|err| format!("cannot inspect config: {err}"))?
            .is_file()
        {
            return Err("config must be a regular file".into());
        }
        Self::read(file)
    }

    fn read(reader: impl Read) -> Result<Self, String> {
        let mut bytes = Vec::with_capacity(MAX_CONFIG_BYTES + 1);
        reader
            .take((MAX_CONFIG_BYTES + 1) as u64)
            .read_to_end(&mut bytes)
            .map_err(|err| format!("cannot read config: {err}"))?;
        if bytes.len() > MAX_CONFIG_BYTES {
            return Err(format!("config exceeds {MAX_CONFIG_BYTES} bytes"));
        }
        // Serde's struct deserializer also accepts positional arrays. The public
        // configuration format requires a JSON object, including for defaults.
        if bytes.iter().find(|byte| !byte.is_ascii_whitespace()) != Some(&b'{') {
            return Err("config must be a JSON object".into());
        }
        let config: Self = serde_json::from_slice(&bytes).map_err(|err| {
            // Keep attacker-provided values and field names out of diagnostics.
            format!(
                "invalid config JSON at line {}, column {}",
                err.line(),
                err.column()
            )
        })?;
        config.validate()?;
        Ok(config)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use std::io;

    #[test]
    fn defaults_and_boundaries_are_valid() {
        assert_eq!(Config::read(&b"{}"[..]).unwrap(), Config::default());
        for value in [MIN_POLL_MS, MAX_POLL_MS] {
            let input = format!(r#"{{"log_level":"debug","shutdown_poll_ms":{value}}}"#);
            let config = Config::read(input.as_bytes()).unwrap();
            assert_eq!(config.shutdown_poll_ms, value);
            assert_eq!(config.log_level, LogLevel::Debug);
        }
    }

    #[test]
    fn malformed_unknown_duplicate_and_invalid_values_fail_closed() {
        for input in [
            "",
            "[]",
            "null",
            "{} {}",
            "{",
            r#"{"unknown":1}"#,
            r#"{"log_level":"trace"}"#,
            r#"{"log_level":null}"#,
            r#"{"log_level":"info","log_level":"debug"}"#,
            r#"{"shutdown_poll_ms":0}"#,
            r#"{"shutdown_poll_ms":9}"#,
            r#"{"shutdown_poll_ms":1001}"#,
            r#"{"shutdown_poll_ms":-1}"#,
            r#"{"shutdown_poll_ms":1.5}"#,
            r#"{"shutdown_poll_ms":"100"}"#,
            r#"{"shutdown_poll_ms":18446744073709551616}"#,
        ] {
            assert!(Config::read(input.as_bytes()).is_err(), "accepted {input}");
        }
        assert!(Config::read(&[0xff][..]).is_err());
    }

    #[test]
    fn input_bound_is_enforced_and_exact_limit_is_accepted() {
        let mut bytes = vec![b' '; MAX_CONFIG_BYTES];
        bytes[..2].copy_from_slice(b"{}");
        assert!(Config::read(bytes.as_slice()).is_ok());
        bytes.push(b' ');
        assert!(Config::read(bytes.as_slice())
            .unwrap_err()
            .contains("exceeds"));
    }

    #[test]
    fn read_errors_are_reported() {
        struct FailedReader;
        impl Read for FailedReader {
            fn read(&mut self, _: &mut [u8]) -> io::Result<usize> {
                Err(io::Error::other("read failed"))
            }
        }
        assert!(Config::read(FailedReader)
            .unwrap_err()
            .contains("cannot read"));
    }
}
