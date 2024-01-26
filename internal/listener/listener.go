package listener

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
)

type Listener struct {
	Resolver *resolver.Resolver
	Protocol models.Protocol
	Port     int
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(resolver *resolver.Resolver, protocol models.Protocol, port int) *Listener {
	ctx, cancel := context.WithCancel(context.Background())
	log.Printf("Creating listener on port %d for protocol %s\n", port, protocol)
	return &Listener{
		Resolver: resolver,
		Protocol: protocol,
		Port:     port,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (l *Listener) Listen() error {
	addr := net.JoinHostPort("", strconv.Itoa(l.Port))
	switch l.Protocol {
	case models.TCP:
		return l.listenTCP(addr)
	case models.UDP:
		return l.listenUDP(addr)
	default:
		return fmt.Errorf("unknown protocol: %s", l.Protocol)
	}
}

func (l *Listener) Close() {
	l.cancel()
}

func (l *Listener) handleDNSQuery(req dns.Msg) ([]byte, error) {
	resp, err := l.Resolver.Resolve(req)
	if err != nil {
		return nil, err
	}
	return resp.Pack()
}
