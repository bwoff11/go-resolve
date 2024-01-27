package listener

import (
	"context"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
)

type Listener struct {
	Resolver *resolver.Resolver
	Protocol config.ProtocolType
	Port     int
	ctx      context.Context
	cancel   context.CancelFunc
}

// Close stops the listener.
func (l *Listener) Close() {
	l.cancel()
}
