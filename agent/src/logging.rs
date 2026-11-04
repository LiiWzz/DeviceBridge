use crate::config::LogLevel;
use serde::Serialize;
use std::io::{self, Write};
use std::time::{SystemTime, UNIX_EPOCH};

#[derive(Serialize)]
struct Event<'a> {
    timestamp_ms: u128,
    level: LogLevel,
    event: &'a str,
    message: &'a str,
}

pub fn emit(mut output: impl Write, level: LogLevel, event: &str, message: &str) -> io::Result<()> {
    let record = Event {
        timestamp_ms: SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_millis(),
        level,
        event,
        message,
    };
    serde_json::to_writer(&mut output, &record)?;
    output.write_all(b"\n")?;
    output.flush()
}

pub fn log(threshold: LogLevel, level: LogLevel, event: &str, message: &str) -> io::Result<()> {
    if level <= threshold {
        emit(io::stderr().lock(), level, event, message)?;
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn logs_are_one_json_line_with_escaped_input() {
        let mut bytes = Vec::new();
        emit(&mut bytes, LogLevel::Error, "failure", "bad\n\"input").unwrap();
        assert_eq!(bytes.iter().filter(|&&byte| byte == b'\n').count(), 1);
        let record: serde_json::Value = serde_json::from_slice(&bytes).unwrap();
        assert_eq!(record["level"], "error");
        assert_eq!(record["message"], "bad\n\"input");
        assert!(record["timestamp_ms"].is_number());
    }

    #[test]
    fn write_failure_is_returned() {
        struct Broken;
        impl Write for Broken {
            fn write(&mut self, _: &[u8]) -> io::Result<usize> {
                Err(io::ErrorKind::BrokenPipe.into())
            }
            fn flush(&mut self) -> io::Result<()> {
                Ok(())
            }
        }
        assert!(emit(Broken, LogLevel::Info, "start", "hello").is_err());
    }
}
