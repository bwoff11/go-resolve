package transport

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/rs/zerolog/log"
)

type Protocol string

const (
	queueBufferSize = 256 // For inbound/outbound queues
	UDPBufferSize   = 512 // For UDP packet size

	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
	ProtocolDOT Protocol = "dot"
	ProtocolDOH Protocol = "doh"
)

type Transport interface {
	Listen() error
	Close() error
}

// Transports now uses a map for dynamic transport management.
type Transports struct {
	Transports map[Protocol]Transport
	Queue      chan QueueItem
}

// Initializes Transports with enabled transports from configuration.
func New(cfg *config.Transport) *Transports {
	ts := &Transports{
		Transports: make(map[Protocol]Transport),
		Queue:      make(chan QueueItem, queueBufferSize),
	}

	if cfg.TCP.Enabled {
		tcp, err := NewTCP(cfg.TCP, ts)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize TCP transport")
		}
		ts.Transports[ProtocolTCP] = tcp
		log.Info().Str("protocol", string(ProtocolTCP)).Msg("transport enabled")
	}

	if cfg.UDP.Enabled {
		udp, err := NewUDP(cfg.UDP, ts)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize UDP transport")
		}
		ts.Transports[ProtocolUDP] = udp
		log.Info().Str("protocol", string(ProtocolUDP)).Msg("transport enabled")
	}

	return ts
}

// Listen on all initialized transports.
func (t *Transports) Listen() error {
	for _, transport := range t.Transports {
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
	for name, transport := range t.Transports {
		if err := transport.Close(); err != nil {
			return err
		}
		log.Info().Str("transport", string(name)).Msg("transport stopped")
	}
	return nil
}
