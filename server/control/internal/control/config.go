package control

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
	"time"
)

// Config contains the foundation service's bounded listener and HTTP settings.
type Config struct {
	Addr              string
	ReadHeaderTimeout time.Duration
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

// LoadConfig treats an explicitly empty value as an error, not as a default.
// Values are never included in errors so configuration cannot leak into logs.
func LoadConfig(lookup func(string) (string, bool)) (Config, error) {
	c := Config{
		Addr: "127.0.0.1:8080", ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second, ShutdownTimeout: 10 * time.Second,
	}
	if value, ok := lookup("DEVICEBRIDGE_CONTROL_ADDR"); ok {
		c.Addr = value
	}
	for _, field := range []struct {
		name string
		dest *time.Duration
	}{
		{"READ_HEADER_TIMEOUT", &c.ReadHeaderTimeout},
		{"READ_TIMEOUT", &c.ReadTimeout},
		{"WRITE_TIMEOUT", &c.WriteTimeout},
		{"IDLE_TIMEOUT", &c.IdleTimeout},
		{"SHUTDOWN_TIMEOUT", &c.ShutdownTimeout},
	} {
		name := "DEVICEBRIDGE_CONTROL_" + field.name
		if value, ok := lookup(name); ok {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return Config{}, fmt.Errorf("%s must be a duration", name)
			}
			*field.dest = duration
		}
	}
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c Config) Validate() error {
	host, port, err := net.SplitHostPort(c.Addr)
	if err != nil {
		return fmt.Errorf("DEVICEBRIDGE_CONTROL_ADDR must be a loopback IP and port")
	}
	ip, err := netip.ParseAddr(host)
	if err != nil || ip.Zone() != "" || !ip.Unmap().IsLoopback() {
		return fmt.Errorf("DEVICEBRIDGE_CONTROL_ADDR must use a loopback IP literal")
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1024 || n > 65535 || strconv.Itoa(n) != port {
		return fmt.Errorf("DEVICEBRIDGE_CONTROL_ADDR port must be an integer from 1024 to 65535")
	}
	for _, field := range []struct {
		name  string
		value time.Duration
	}{
		{"READ_HEADER_TIMEOUT", c.ReadHeaderTimeout},
		{"READ_TIMEOUT", c.ReadTimeout},
		{"WRITE_TIMEOUT", c.WriteTimeout},
		{"IDLE_TIMEOUT", c.IdleTimeout},
		{"SHUTDOWN_TIMEOUT", c.ShutdownTimeout},
	} {
		if field.value < time.Millisecond || field.value > 5*time.Minute {
			return fmt.Errorf("DEVICEBRIDGE_CONTROL_%s must be between 1ms and 5m", field.name)
		}
	}
	if c.ReadHeaderTimeout > c.ReadTimeout {
		return fmt.Errorf("DEVICEBRIDGE_CONTROL_READ_HEADER_TIMEOUT must not exceed READ_TIMEOUT")
	}
	return nil
}
