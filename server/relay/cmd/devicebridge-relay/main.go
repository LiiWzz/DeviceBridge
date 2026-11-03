package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/LiiWzz/DeviceBridge/server/relay/internal/config"
)

const version = "0.1.0-dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	cfg := config.Default()
	flags := flag.NewFlagSet("devicebridge-relay", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showVersion := flags.Bool("version", false, "print version and exit")
	check := flags.Bool("check", false, "validate settings and exit without starting")
	flags.StringVar(&cfg.Transport, "transport", cfg.Transport, "transport (disabled only; unimplemented)")
	flags.IntVar(&cfg.MaxStreams, "max-streams", cfg.MaxStreams, "reserved future stream ceiling (1..1024)")
	flags.IntVar(&cfg.MaxBufferBytes, "max-buffer-bytes", cfg.MaxBufferBytes, "reserved future buffer ceiling (1..1048576 bytes)")
	flags.Usage = func() {
		fmt.Fprintln(stderr, "Usage: devicebridge-relay [flags]\nM01 idle skeleton: no listeners; transport unimplemented.")
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected positional arguments")
		return 2
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if *showVersion {
		fmt.Fprintln(stdout, "devicebridge-relay", version)
		return 0
	}
	if *check {
		fmt.Fprintln(stdout, "configuration valid; transport unimplemented and disabled")
		return 0
	}
	if ctx.Err() != nil {
		return 0
	}
	fmt.Fprintln(stdout, "relay idle; no listeners; transport unimplemented and disabled")
	<-ctx.Done()
	fmt.Fprintln(stdout, "relay stopped")
	return 0
}
