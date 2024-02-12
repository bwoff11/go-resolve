package transport

import (
	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/rs/zerolog/log"
)

const (
	queueBufferSize = 256 // For inbound/outbound queues
	UDPBufferSize   = 512 // For UDP packet size
)

type Transport interface {
	Listen() error
	Close() error
}

// Transports now uses a map for dynamic transport management.
type Transports struct {
	Transports map[common.Protocol]Transport
	Queue      chan QueueItem
}

// Initializes Transports with enabled transports from configuration.
func New(cfg *config.Transport) *Transports {
	ts := &Transports{
		Transports: make(map[common.Protocol]Transport),
		Queue:      make(chan QueueItem, queueBufferSize),
	}

	if cfg.TCP.Enabled {
		tcp, err := NewTCP(cfg.TCP, ts.Queue)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize TCP transport")
		}
		ts.Transports[common.ProtocolTCP] = tcp
		log.Info().Str("common.Protocol", string(common.ProtocolTCP)).Msg("transport enabled")
	}

	if cfg.UDP.Enabled {
		udp, err := NewUDP(cfg.UDP, ts.Queue)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to initialize UDP transport")
		}
		ts.Transports[common.ProtocolUDP] = udp
		log.Info().Str("common.Protocol", string(common.ProtocolUDP)).Msg("transport enabled")
	}

	return ts
}

func (t *Transports) Start() {
	t.Listen()
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
