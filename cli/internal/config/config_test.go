package config

import (
	"testing"
	"time"
)

func TestValidateURL(t *testing.T) {
	for _, endpoint := range []string{"https://localhost:443", "https://relay.example", "https://relay.example.:65535", "https://127.0.0.1:8443/", "https://[::1]:443"} {
		t.Run(endpoint, func(t *testing.T) {
			cfg := Default()
			cfg.RelayURL = endpoint
			if err := cfg.Validate(); err != nil {
				t.Fatalf("valid endpoint rejected: %v", err)
			}
		})
	}
	for _, endpoint := range []string{
		"", "relay.example", "http://relay.example", "ssh://relay.example", "https:///",
		"https://user:password@relay.example", "https://relay.example?token=test",
		"https://relay.example?", "https://relay.example/#fragment", "https://relay.example#",
		"https://relay.example/api", "https://relay.example/%2f", "https://relay.example:0",
		"https://relay.example:65536", "https://relay.example:abc", "https://relay.example:",
		"https://relay.example/%zz", "https://relay.example\n", "https:relay.example",
		"https://[::::]:443", "https://[relay.example]:443", "https://[fe80::1%25eth0]",
		"https://::1", "https://-relay.example", "https://relay..example", "https://relay_example",
		"https://.", "https://relay.example../", "https://例.example",
	} {
		t.Run(endpoint, func(t *testing.T) {
			cfg := Default()
			cfg.RelayURL = endpoint
			if err := cfg.Validate(); err == nil {
				t.Fatal("invalid endpoint accepted")
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	for _, timeout := range []time.Duration{-time.Second, 0, time.Minute + 1} {
		cfg := Default()
		cfg.Timeout = timeout
		if err := cfg.Validate(); err == nil {
			t.Errorf("invalid timeout %v accepted", timeout)
		}
	}
	for _, timeout := range []time.Duration{1, time.Minute} {
		cfg := Default()
		cfg.Timeout = timeout
		if err := cfg.Validate(); err != nil {
			t.Errorf("valid timeout %v rejected: %v", timeout, err)
		}
	}
}
