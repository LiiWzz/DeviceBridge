package config

import "testing"

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want bool
	}{
		{"defaults", Default(), true},
		{"minimum bounds", Config{"disabled", 1, 1}, true},
		{"maximum bounds", Config{"disabled", 1024, 1048576}, true},
		{"zero value", Config{}, false},
		{"tls unavailable", Config{"tls", 16, 32768}, false},
		{"empty transport", Config{"", 16, 32768}, false},
		{"zero streams", Config{"disabled", 0, 32768}, false},
		{"negative streams", Config{"disabled", -1, 32768}, false},
		{"too many streams", Config{"disabled", 1025, 32768}, false},
		{"zero buffer", Config{"disabled", 16, 0}, false},
		{"negative buffer", Config{"disabled", 16, -1}, false},
		{"oversized buffer", Config{"disabled", 16, 1048577}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.cfg.Validate(); (err == nil) != tt.want {
				t.Fatalf("Validate() = %v; want valid=%v", err, tt.want)
			}
		})
	}
}
