package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

func TestRunCommands(t *testing.T) {
	tests := []struct {
		name string
		args []string
		code int
		text string
	}{
		{"help", []string{"--help"}, 0, "no listeners"},
		{"version", []string{"--version"}, 0, "devicebridge-relay 0.1.0-dev"},
		{"check", []string{"--check"}, 0, "transport unimplemented and disabled"},
		{"transport", []string{"--transport", "tls", "--check"}, 2, "transport is unimplemented"},
		{"zero streams", []string{"--max-streams", "0", "--check"}, 2, "max-streams"},
		{"oversized buffer", []string{"--max-buffer-bytes", "1048577", "--check"}, 2, "max-buffer-bytes"},
		{"malformed integer", []string{"--max-streams", "abc"}, 2, "invalid value or flag"},
		{"unknown flag", []string{"--listen", "0.0.0.0:443"}, 2, "invalid value or flag"},
		{"unexpected arguments", []string{"shell"}, 2, "unexpected positional"},
		{"invalid config even with version", []string{"--version", "--transport", "tls"}, 2, "transport is unimplemented"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(context.Background(), tt.args, &stdout, &stderr); code != tt.code {
				t.Fatalf("exit=%d; want %d; stderr=%s", code, tt.code, stderr.String())
			}
			if !strings.Contains(stdout.String()+stderr.String(), tt.text) {
				t.Fatalf("missing %q in stdout=%q stderr=%q", tt.text, stdout.String(), stderr.String())
			}
		})
	}
}

func TestIdleStopsOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var stdout, stderr bytes.Buffer
	done := make(chan int, 1)
	// Wait until the run loop reports it has started, without racing on buffers.
	started := make(chan struct{})
	go func() {
		done <- run(ctx, nil, startWriter{&stdout, started}, &stderr)
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("idle relay did not start")
	}
	cancel()
	select {
	case code := <-done:
		if code != 0 || !strings.Contains(stdout.String(), "relay stopped") || stderr.Len() != 0 {
			t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
		}
	case <-time.After(2 * time.Second):
		t.Fatal("relay did not shut down after cancellation")
	}
}

type startWriter struct {
	buf     *bytes.Buffer
	started chan struct{}
}

func (w startWriter) Write(p []byte) (int, error) {
	n, err := w.buf.Write(p)
	if strings.Contains(string(p), "relay idle") {
		close(w.started)
	}
	return n, err
}

func TestAlreadyCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	var stdout, stderr bytes.Buffer
	if code := run(ctx, nil, &stdout, &stderr); code != 0 || stdout.Len() != 0 || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestMalformedArgumentsDoNotLeakValues(t *testing.T) {
	const marker = "db007-sensitive-input"
	for _, args := range [][]string{
		{"--max-streams", marker},
		{"--max-buffer-bytes", marker},
		{"--check=" + marker},
		{"--" + marker},
		{"--transport", marker, "--check"},
	} {
		var stdout, stderr bytes.Buffer
		if code := run(context.Background(), args, &stdout, &stderr); code != 2 {
			t.Fatalf("exit=%d; want 2", code)
		}
		if strings.Contains(stdout.String()+stderr.String(), marker) {
			t.Fatal("sensitive input leaked into diagnostics")
		}
	}
}
