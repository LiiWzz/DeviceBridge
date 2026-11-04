// Package config validates CLI settings without contacting the relay.
package config

import (
	"fmt"
	"net/netip"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	RelayURL string
	Timeout  time.Duration
}

func Default() Config {
	return Config{RelayURL: "https://localhost:443", Timeout: 5 * time.Second}
}

func (c Config) Validate() error {
	u, err := url.Parse(c.RelayURL)
	if err != nil || u.Scheme != "https" || u.Hostname() == "" || u.Opaque != "" {
		return fmt.Errorf("relay-url must be an absolute HTTPS URL")
	}
	if u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(c.RelayURL, "#") || (u.Path != "" && u.Path != "/") {
		return fmt.Errorf("relay-url must not contain credentials, query, fragment, or a non-root path")
	}
	if strings.HasSuffix(u.Host, ":") {
		return fmt.Errorf("relay-url port must not be empty")
	}
	if !validHost(u.Hostname(), strings.HasPrefix(u.Host, "[")) {
		return fmt.Errorf("relay-url must contain a valid ASCII DNS hostname or IP address")
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return fmt.Errorf("relay-url port must be between 1 and 65535")
		}
	}
	if c.Timeout <= 0 || c.Timeout > time.Minute {
		return fmt.Errorf("timeout must be greater than zero and at most 1m")
	}
	return nil
}

func validHost(host string, bracketed bool) bool {
	if addr, err := netip.ParseAddr(host); err == nil {
		return addr.Zone() == "" && bracketed == addr.Is6()
	}
	if bracketed {
		return false
	}
	host = strings.TrimSuffix(host, ".")
	if len(host) == 0 || len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, c := range label {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
				return false
			}
		}
	}
	return true
}
