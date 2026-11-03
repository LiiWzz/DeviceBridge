mod cli;
mod config;
mod logging;

use cli::{Action, Cli};
use config::{Config, LogLevel};
use std::io::{self, Write};
use std::process::ExitCode;

fn main() -> ExitCode {
    match execute() {
        Ok(()) => ExitCode::SUCCESS,
        Err((code, message)) => {
            // A closed log sink must not panic, including while reporting failure.
            let _ = logging::emit(
                io::stderr().lock(),
                LogLevel::Error,
                "agent_error",
                &message,
            );
            ExitCode::from(code)
        }
    }
}

fn execute() -> Result<(), (u8, String)> {
    let cli = Cli::parse(std::env::args_os().skip(1)).map_err(|err| (2, err))?;
    match cli.action {
        Action::Help => return print(cli::HELP),
        Action::Version => {
            return print(concat!(
                "devicebridge-agent ",
                env!("CARGO_PKG_VERSION"),
                "\n"
            ))
        }
        _ => {}
    }
    let config = cli.config().map_err(|err| (2, err))?;
    if cli.action == Action::Check {
        return print("configuration valid\n");
    }
    run(&config).map_err(|err| (1, err.to_string()))
}

fn print(text: &str) -> Result<(), (u8, String)> {
    io::stdout()
        .lock()
        .write_all(text.as_bytes())
        .map_err(|err| (1, format!("cannot write output: {err}")))
}

#[cfg(unix)]
fn run(config: &Config) -> io::Result<()> {
    use signal_hook::consts::{SIGINT, SIGTERM};
    use std::sync::atomic::{AtomicBool, Ordering};
    use std::sync::Arc;
    use std::time::Duration;

    let stopping = Arc::new(AtomicBool::new(false));
    let interrupt = signal_hook::flag::register(SIGINT, Arc::clone(&stopping))?;
    let terminate = match signal_hook::flag::register(SIGTERM, Arc::clone(&stopping)) {
        Ok(id) => id,
        Err(err) => {
            signal_hook::low_level::unregister(interrupt);
            return Err(err);
        }
    };
    let result = (|| {
        logging::log(
            config.log_level,
            LogLevel::Info,
            "agent_started",
            "M01 foundation; no transport",
        )?;
        while !stopping.load(Ordering::Relaxed) {
            // No allocations, logs, network or disk writes in the idle loop.
            std::thread::sleep(Duration::from_millis(config.shutdown_poll_ms));
        }
        logging::log(
            config.log_level,
            LogLevel::Info,
            "agent_stopped",
            "shutdown requested",
        )
    })();
    signal_hook::low_level::unregister(interrupt);
    signal_hook::low_level::unregister(terminate);
    result
}

#[cfg(not(unix))]
fn run(_: &Config) -> io::Result<()> {
    Err(io::Error::new(
        io::ErrorKind::Unsupported,
        "runtime requires Unix SIGINT/SIGTERM support",
    ))
}
