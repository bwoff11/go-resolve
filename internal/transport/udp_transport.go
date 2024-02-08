package transport

import (
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type UDPTransport struct {
	Transports *Transports
	Conn       net.PacketConn
}

func NewUDP(c config.Protocol, ts *Transports) (Transport, error) {
	addr := net.JoinHostPort("", strconv.Itoa(c.Port))
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{
		Transports: ts,
		Conn:       conn,
	}, nil
}

// Bind reads from the connection and enqueues received DNS queries.
func (ut *UDPTransport) Listen() error {

	go func() {
		buf := make([]byte, UDPBufferSize)
		for {
			n, clientAddr, err := ut.Conn.ReadFrom(buf)
			if err != nil {
				log.Error().Err(err).Msg("error reading from udp connection")
				return
			}

			query := make([]byte, n)
			copy(query, buf[:n])

			go ut.handleUDPQuery(query, clientAddr)
		}
	}()

	log.Info().Str("protocol", "udp").Msg("transport listening")
	return nil
}

// handleUDPQuery processes a single UDP packet.
func (ut *UDPTransport) handleUDPQuery(query []byte, clientAddr net.Addr) {
	var req dns.Msg
	if err := req.Unpack(query); err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error unpacking dns query")
		return
	}

	ut.Transports.InboundQueue <- &UDPQueueItem{
		Msg:  req,
		Addr: clientAddr,
		Conn: ut.Conn,
	}
}

func (ut *UDPTransport) Close() error {
	return ut.Conn.Close()
}
