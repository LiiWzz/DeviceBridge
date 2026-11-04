package control

import (
	"strings"
	"testing"
	"time"
)

func configFrom(values map[string]string) (Config, error) {
	return LoadConfig(func(key string) (string, bool) { v, ok := values[key]; return v, ok })
}

func TestConfigDefaultsAndOverrides(t *testing.T) {
	c, err := configFrom(nil)
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != "127.0.0.1:8080" || c.ReadHeaderTimeout != 5*time.Second || c.ReadTimeout != 10*time.Second || c.WriteTimeout != 10*time.Second || c.IdleTimeout != time.Minute || c.ShutdownTimeout != 10*time.Second {
		t.Fatalf("unexpected defaults: %+v", c)
	}
	values := map[string]string{
		"DEVICEBRIDGE_CONTROL_ADDR":                "[::1]:9090",
		"DEVICEBRIDGE_CONTROL_READ_HEADER_TIMEOUT": "2s",
		"DEVICEBRIDGE_CONTROL_READ_TIMEOUT":        "3s",
		"DEVICEBRIDGE_CONTROL_WRITE_TIMEOUT":       "4s",
		"DEVICEBRIDGE_CONTROL_IDLE_TIMEOUT":        "5s",
		"DEVICEBRIDGE_CONTROL_SHUTDOWN_TIMEOUT":    "6s",
	}
	c, err = configFrom(values)
	if err != nil {
		t.Fatal(err)
	}
	if c.Addr != "[::1]:9090" || c.ReadHeaderTimeout != 2*time.Second || c.ReadTimeout != 3*time.Second || c.WriteTimeout != 4*time.Second || c.IdleTimeout != 5*time.Second || c.ShutdownTimeout != 6*time.Second {
		t.Fatalf("overrides not applied: %+v", c)
	}
}

func TestRejectMalformedConfig(t *testing.T) {
	for _, value := range []string{"", ":8080", "0.0.0.0:8080", "[::]:8080", "192.0.2.1:8080", "localhost:8080", "[::1%lo0]:8080", "127.0.0.1", "127.0.0.1:0", "127.0.0.1:80", "127.0.0.1:65536", "127.0.0.1:http", "127.0.0.1:+8080", "127.0.0.1:08080", "127.0.0.1:8080\n"} {
		t.Run("addr="+value, func(t *testing.T) {
			if _, err := configFrom(map[string]string{"DEVICEBRIDGE_CONTROL_ADDR": value}); err == nil {
				t.Fatalf("accepted invalid address %q", value)
			}
		})
	}
	for _, key := range []string{"READ_HEADER_TIMEOUT", "READ_TIMEOUT", "WRITE_TIMEOUT", "IDLE_TIMEOUT", "SHUTDOWN_TIMEOUT"} {
		for _, value := range []string{"", "secret-do-not-log", "0", "-1s", "1ns", "301s", "999999999999999999h"} {
			t.Run(key+"="+value, func(t *testing.T) {
				_, err := configFrom(map[string]string{"DEVICEBRIDGE_CONTROL_" + key: value})
				if err == nil {
					t.Fatal("accepted invalid duration")
				}
				if strings.Contains(err.Error(), "secret-do-not-log") {
					t.Fatal("configuration value exposed in error")
				}
			})
		}
	}
	if _, err := configFrom(map[string]string{"DEVICEBRIDGE_CONTROL_READ_TIMEOUT": "1s"}); err == nil {
		t.Fatal("accepted read timeout shorter than header timeout")
	}
}
