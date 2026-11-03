// Package transport defines a local implementation seam only. It defines no
// framing, authentication, authorization, enrollment, or network protocol.
package transport

import (
	"context"
	"errors"
	"io"
)

var ErrUnimplemented = errors.New("relay transport is unimplemented")

// Connector is deliberately internal and provisional. Production transport
// requires separately reviewed identity, session authorization, and bounds.
type Connector interface {
	Open(context.Context) (io.ReadWriteCloser, error)
}

// Disabled always rejects requests without opening a socket or allocating a
// stream. It must never be replaced with a permissive development transport.
type Disabled struct{}

var _ Connector = Disabled{}

func (Disabled) Open(context.Context) (io.ReadWriteCloser, error) {
	return nil, ErrUnimplemented
}
