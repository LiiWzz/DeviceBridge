package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LiiWzz/DeviceBridge/server/control/internal/control"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	c, err := control.LoadConfig(os.LookupEnv)
	if err != nil {
		logger.Error("invalid configuration", "error", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err = control.Run(ctx, c, logger)
	stop()
	if err != nil {
		logger.Error("control API failed", "error", err)
		os.Exit(1)
	}
}
