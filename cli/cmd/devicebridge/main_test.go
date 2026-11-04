package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name string
		args []string
		code int
		text string
	}{
		{"no args", nil, 0, "Usage:"},
		{"help", []string{"help"}, 0, "transport unimplemented"},
		{"help flag", []string{"--help"}, 0, "Usage:"},
		{"version", []string{"version"}, 0, "devicebridge 0.1.0-dev"},
		{"version flag", []string{"--version"}, 0, "devicebridge 0.1.0-dev"},
		{"check", []string{"check"}, 0, "no connection attempted"},
		{"custom settings", []string{"check", "--relay-url", "https://relay.example:8443", "--timeout", "1s"}, 0, "configuration valid"},
		{"check help", []string{"check", "--help"}, 0, "Usage of check"},
		{"insecure endpoint", []string{"check", "--relay-url", "http://relay.example"}, 2, "HTTPS"},
		{"credential endpoint", []string{"check", "--relay-url", "https://user:secret@relay.example"}, 2, "must not contain credentials"},
		{"bad timeout", []string{"check", "--timeout", "abc"}, 2, "invalid value"},
		{"zero timeout", []string{"check", "--timeout", "0s"}, 2, "timeout"},
		{"unknown flag", []string{"check", "--insecure"}, 2, "flag provided but not defined"},
		{"positional", []string{"check", "device"}, 2, "positional"},
		{"extra version args", []string{"version", "device"}, 2, "does not accept arguments"},
		{"extra help args", []string{"help", "device"}, 2, "does not accept arguments"},
		{"connect unavailable", []string{"connect", "device"}, 2, "unimplemented"},
		{"shell unavailable", []string{"shell", "device", "id"}, 2, "unimplemented"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(tt.args, &stdout, &stderr); code != tt.code {
				t.Fatalf("exit=%d; want %d; stderr=%s", code, tt.code, stderr.String())
			}
			output := stdout.String() + stderr.String()
			if !strings.Contains(output, tt.text) {
				t.Fatalf("missing %q in %q", tt.text, output)
			}
			if strings.Contains(output, "secret") {
				t.Fatal("credential value leaked in output")
			}
		})
	}
}
