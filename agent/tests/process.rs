use std::fs::{self, OpenOptions};
use std::io::Write;
use std::path::PathBuf;
use std::process::{Command, Output};
use std::sync::atomic::{AtomicUsize, Ordering};

static NEXT: AtomicUsize = AtomicUsize::new(0);

struct Fixture(PathBuf);

impl Fixture {
    fn new(bytes: &[u8]) -> Self {
        let path = std::env::temp_dir().join(format!(
            "devicebridge-agent-test-{}-{}",
            std::process::id(),
            NEXT.fetch_add(1, Ordering::Relaxed)
        ));
        let mut file = OpenOptions::new()
            .write(true)
            .create_new(true)
            .open(&path)
            .unwrap();
        file.write_all(bytes).unwrap();
        Self(path)
    }
}

impl Drop for Fixture {
    fn drop(&mut self) {
        let _ = fs::remove_file(&self.0);
    }
}

fn agent() -> Command {
    Command::new(env!("CARGO_BIN_EXE_devicebridge-agent"))
}

fn assert_config_error(output: Output) {
    assert_eq!(output.status.code(), Some(2));
    assert!(output.stdout.is_empty());
    let record: serde_json::Value = serde_json::from_slice(&output.stderr).unwrap();
    assert_eq!(record["event"], "agent_error");
    assert_eq!(record["level"], "error");
}

#[test]
fn informational_commands_and_invalid_options() {
    for option in ["--help", "--version", "--check-config"] {
        let output = agent().arg(option).output().unwrap();
        assert!(output.status.success());
        assert!(!output.stdout.is_empty());
        assert!(output.stderr.is_empty());
    }
    assert_config_error(agent().arg("--remote-command").output().unwrap());
}

#[test]
fn file_cli_precedence_and_invalid_file_cannot_be_masked() {
    let valid = Fixture::new(br#"{"log_level":"error","shutdown_poll_ms":1000}"#);
    let output = agent()
        .arg("--config")
        .arg(&valid.0)
        .args([
            "--shutdown-poll-ms",
            "10",
            "--log-level",
            "debug",
            "--check-config",
        ])
        .output()
        .unwrap();
    assert!(output.status.success());
    let invalid = Fixture::new(br#"{"shutdown_poll_ms":0}"#);
    assert_config_error(
        agent()
            .arg("--config")
            .arg(&invalid.0)
            .args(["--shutdown-poll-ms", "10", "--check-config"])
            .output()
            .unwrap(),
    );
}

#[test]
fn missing_directory_oversized_and_malformed_config_are_rejected() {
    let missing = Fixture::new(b"{}");
    fs::remove_file(&missing.0).unwrap();
    assert_config_error(agent().arg("--config").arg(&missing.0).output().unwrap());
    assert_config_error(
        agent()
            .arg("--config")
            .arg(std::env::temp_dir())
            .output()
            .unwrap(),
    );
    for input in [vec![b' '; 4097], b"{bad}".to_vec(), vec![0xff]] {
        let fixture = Fixture::new(&input);
        assert_config_error(agent().arg("--config").arg(&fixture.0).output().unwrap());
    }
}

#[cfg(unix)]
mod unix {
    use super::*;
    use std::io::{BufRead, BufReader, Read};
    use std::process::{Child, Stdio};
    use std::time::{Duration, Instant};

    struct Running(Child);
    impl Drop for Running {
        fn drop(&mut self) {
            let _ = self.0.kill();
            let _ = self.0.wait();
        }
    }

    fn wait(child: &mut Child) -> std::process::ExitStatus {
        let deadline = Instant::now() + Duration::from_secs(3);
        loop {
            if let Some(status) = child.try_wait().unwrap() {
                return status;
            }
            assert!(
                Instant::now() < deadline,
                "child did not exit within three seconds"
            );
            std::thread::sleep(Duration::from_millis(10));
        }
    }

    #[test]
    fn sigint_and_sigterm_exit_successfully_with_structured_shutdown() {
        for signal in ["INT", "TERM"] {
            let mut child = Running(
                agent()
                    .args(["--shutdown-poll-ms", "10"])
                    .stdout(Stdio::null())
                    .stderr(Stdio::piped())
                    .spawn()
                    .unwrap(),
            );
            let stderr = child.0.stderr.take().unwrap();
            let (tx, rx) = std::sync::mpsc::channel();
            let reader = std::thread::spawn(move || {
                let mut reader = BufReader::new(stderr);
                let mut line = String::new();
                reader.read_line(&mut line).unwrap();
                let _ = tx.send((line, reader));
            });
            let (line, mut stderr) = rx
                .recv_timeout(Duration::from_secs(3))
                .expect("startup event missing");
            let record: serde_json::Value = serde_json::from_str(&line).unwrap();
            assert_eq!(record["event"], "agent_started");
            assert!(Command::new("kill")
                .args(["-s", signal, &child.0.id().to_string()])
                .status()
                .unwrap()
                .success());
            assert!(wait(&mut child.0).success());
            let mut rest = String::new();
            stderr.read_to_string(&mut rest).unwrap();
            let record: serde_json::Value = serde_json::from_str(&rest).unwrap();
            assert_eq!(record["event"], "agent_stopped");
            reader.join().unwrap();
        }
    }

    #[test]
    fn fifo_config_is_rejected_without_waiting_for_a_writer() {
        let fifo = Fixture::new(b"");
        fs::remove_file(&fifo.0).unwrap();
        assert!(Command::new("mkfifo")
            .arg(&fifo.0)
            .status()
            .unwrap()
            .success());
        let mut child = Running(
            agent()
                .arg("--config")
                .arg(&fifo.0)
                .stdout(Stdio::null())
                .stderr(Stdio::null())
                .spawn()
                .unwrap(),
        );
        assert_eq!(wait(&mut child.0).code(), Some(2));
    }
}
