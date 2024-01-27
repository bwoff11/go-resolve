package listener

import (
	"context"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
)

type Listener struct {
	Resolver *resolver.Resolver
	Protocol config.ProtocolType
	Port     int
	ctx      context.Context
	cancel   context.CancelFunc
}

// handleDNSQuery processes a DNS query using the resolver.
func (l *Listener) handleDNSQuery(req dns.Msg) ([]byte, error) {
	resp, err := l.Resolver.Resolve(req)
	if err != nil {
		return nil, err
	}
	return resp.Pack()
}

// Close stops the listener.
func (l *Listener) Close() {
	l.cancel()
}
