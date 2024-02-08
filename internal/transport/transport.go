package transport

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

const (
	queueBufferSize = 256 // For inbound/outbound queues
	UDPBufferSize   = 512 // For UDP packet size
)

type QueueItem interface {
	Question() dns.Question
}

type Transport interface {
	Listen() error
	Close() error
}

// Transports now uses a map for dynamic transport management.
type Transports struct {
	transports    map[string]Transport
	InboundQueue  chan QueueItem
	OutboundQueue chan QueueItem
}

// Initializes Transports with enabled transports from configuration.
func New(cfg *config.Transport) (*Transports, error) {
	ts := &Transports{
		transports:    make(map[string]Transport),
		InboundQueue:  make(chan QueueItem, queueBufferSize),
		OutboundQueue: make(chan QueueItem, queueBufferSize),
	}

	if cfg.TCP.Enabled {
		tcp, err := NewTCP(cfg.TCP, ts)
		if err != nil {
			return nil, err
		}
		ts.transports["TCP"] = tcp
		log.Info().Str("protocol", "tcp").Msg("transport enabled")
	}

	if cfg.UDP.Enabled {
		udp, err := NewUDP(cfg.UDP, ts)
		if err != nil {
			return nil, err
		}
		ts.transports["UDP"] = udp
		log.Info().Str("protocol", "udp").Msg("transport enabled")
	}

	return ts, nil
}

// Listen on all initialized transports.
func (t *Transports) Listen() error {
	for _, transport := range t.transports {
		go func(transport Transport) {
			if err := transport.Listen(); err != nil {
				log.Error().Err(err).Msg("transport error")
			}
		}(transport)
	}
	return nil
}

// Stop all initialized transports.
func (t *Transports) Stop() error {
	for name, transport := range t.transports {
		if err := transport.Close(); err != nil {
			return err
		}
		log.Info().Str("transport", name).Msg("transport stopped")
	}
	return nil
}
