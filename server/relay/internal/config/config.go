// Package config validates the relay's foundation-only settings.
package config

import "fmt"

// Config reserves resource ceilings for a future transport. No stream or buffer
// is allocated by this milestone, and these are not a wire protocol contract.
type Config struct {
	Transport      string
	MaxStreams     int
	MaxBufferBytes int
}

func Default() Config {
	return Config{Transport: "disabled", MaxStreams: 16, MaxBufferBytes: 32 * 1024}
}

func (c Config) Validate() error {
	if c.Transport != "disabled" {
		return fmt.Errorf("transport is unimplemented; only disabled is supported")
	}
	if c.MaxStreams < 1 || c.MaxStreams > 1024 {
		return fmt.Errorf("max-streams must be between 1 and 1024")
	}
	if c.MaxBufferBytes < 1 || c.MaxBufferBytes > 1024*1024 {
		return fmt.Errorf("max-buffer-bytes must be between 1 and 1048576")
	}
	return nil
}
