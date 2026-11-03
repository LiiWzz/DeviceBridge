package transport

import (
	"context"
	"errors"
	"testing"
)

func TestDisabledAlwaysFailsClosed(t *testing.T) {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	for _, ctx := range []context.Context{context.Background(), cancelled} {
		stream, err := (Disabled{}).Open(ctx)
		if stream != nil || !errors.Is(err, ErrUnimplemented) {
			t.Fatalf("Open() = %v, %v; want nil stream and ErrUnimplemented", stream, err)
		}
	}
}
