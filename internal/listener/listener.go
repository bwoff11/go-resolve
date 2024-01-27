package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

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

// New creates a new DNS listener with a resolver configured with upstreams and cache settings.
func New(protocol config.ProtocolType, port int, resolver *resolver.Resolver) *Listener {
	ctx, cancel := context.WithCancel(context.Background())

	log.Printf("Creating %s listener on port %d\n", protocol, port)
	return &Listener{
		Resolver: resolver,
		Protocol: protocol,
		Port:     port,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Listen starts the DNS listener based on the configured protocol.
func (l *Listener) Listen() error {
	addr := net.JoinHostPort("", strconv.Itoa(l.Port))
	switch l.Protocol {
	case "tcp":
		return l.listenTCP(addr)
	case "udp":
		return l.listenUDP(addr)
	default:
		return fmt.Errorf("unknown protocol: %s", l.Protocol)
	}
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
