package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LiiWzz/DeviceBridge/cli/internal/config"
)

const version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "Usage: devicebridge <help|version|check>\nM01 developer CLI: transport unimplemented; no remote operations.\n  check [--relay-url https://localhost:443] [--timeout 5s]\n    Validate local settings only; does not contact a server.")
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stdout)
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "help does not accept arguments")
			return 2
		}
		usage(stdout)
		return 0
	case "version", "--version":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "version does not accept arguments")
			return 2
		}
		fmt.Fprintln(stdout, "devicebridge", version)
		return 0
	case "check":
		cfg := config.Default()
		flags := flag.NewFlagSet("check", flag.ContinueOnError)
		// flag's default diagnostics echo untrusted values and option names.
		flags.SetOutput(io.Discard)
		flags.Usage = func() {
			fmt.Fprintln(stderr, "Usage of check:")
			flags.SetOutput(stderr)
			flags.PrintDefaults()
			flags.SetOutput(io.Discard)
		}
		flags.StringVar(&cfg.RelayURL, "relay-url", cfg.RelayURL, "future relay HTTPS endpoint; never contacted")
		flags.DurationVar(&cfg.Timeout, "timeout", cfg.Timeout, "reserved future request timeout (0 < timeout <= 1m)")
		if err := flags.Parse(args[1:]); err != nil {
			if errors.Is(err, flag.ErrHelp) {
				return 0
			}
			fmt.Fprintln(stderr, "invalid value or flag; use --help")
			return 2
		}
		if flags.NArg() != 0 {
			fmt.Fprintln(stderr, "check does not accept positional arguments")
			return 2
		}
		if err := cfg.Validate(); err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		fmt.Fprintln(stdout, "configuration valid; transport unimplemented; no connection attempted")
		return 0
	default:
		fmt.Fprintln(stderr, "command unavailable; transport and remote operations are unimplemented (use help)")
		return 2
	}
}
