package control

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"
)

type service struct {
	ready atomic.Bool
	http  *http.Server
}

func newService(c Config, logger *slog.Logger) *service {
	s := &service{}
	s.http = &http.Server{
		Addr: c.Addr, Handler: http.HandlerFunc(s.handle),
		ReadHeaderTimeout: c.ReadHeaderTimeout, ReadTimeout: c.ReadTimeout,
		WriteTimeout: c.WriteTimeout, IdleTimeout: c.IdleTimeout,
		MaxHeaderBytes: 8 << 10, ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	return s
}

func (s *service) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	status, body := http.StatusOK, "{\"status\":\"ok\"}\n"
	switch {
	case r.URL.Path != "/healthz" && r.URL.Path != "/readyz":
		status, body = http.StatusNotFound, "{\"error\":\"not_found\"}\n"
	case r.Method != http.MethodGet && r.Method != http.MethodHead:
		w.Header().Set("Allow", "GET, HEAD")
		status, body = http.StatusMethodNotAllowed, "{\"error\":\"method_not_allowed\"}\n"
	case r.URL.Path == "/readyz":
		body = "{\"status\":\"ready\"}\n"
		if !s.ready.Load() {
			status, body = http.StatusServiceUnavailable, "{\"status\":\"shutting_down\"}\n"
		}
	}
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = io.WriteString(w, body)
	}
}

// Run binds only after validation, serves health probes, and drains on cancellation.
func Run(ctx context.Context, c Config, logger *slog.Logger) error {
	if err := c.Validate(); err != nil {
		return err
	}
	listener, err := net.Listen("tcp", c.Addr)
	if err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return newService(c, logger).serve(ctx, listener, c, logger)
}

func (s *service) serve(ctx context.Context, listener net.Listener, c Config, logger *slog.Logger) error {
	defer listener.Close()
	served := make(chan error, 1)
	s.ready.Store(true)
	go func() { served <- s.http.Serve(listener) }()
	logger.Info("control API listening", "addr", listener.Addr().String())
	select {
	case err := <-served:
		s.ready.Store(false)
		// Serve errors do not close any existing connections.
		_ = s.http.Close()
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}
	s.ready.Store(false)
	logger.Info("control API shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), c.ShutdownTimeout)
	defer cancel()
	if err := s.http.Shutdown(shutdownCtx); err != nil {
		_ = s.http.Close()
		<-served
		return fmt.Errorf("shutdown: %w", err)
	}
	if err := <-served; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve: %w", err)
	}
	logger.Info("control API stopped")
	return nil
}
