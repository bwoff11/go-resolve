package transport

import (
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type UDPTransport struct {
	Conn  net.PacketConn
	Queue chan QueueItem
}

func NewUDP(c config.Protocol, q chan QueueItem) (Transport, error) {
	addr := net.JoinHostPort("", strconv.Itoa(c.Port))
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	return &UDPTransport{
		Conn:  conn,
		Queue: q,
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

func (ut *UDPTransport) handleUDPQuery(query []byte, clientAddr net.Addr) {
	var req dns.Msg
	if err := req.Unpack(query); err != nil {
		log.Error().Err(err).Str("protocol", "udp").Msg("error unpacking dns query")
		return
	}

	// Create a UDPConnection adapter
	udpConn := &UDPConnection{
		Addr: clientAddr,
		Conn: ut.Conn,
	}

	// Enqueue the query with the generic QueueItem structure
	ut.Queue <- QueueItem{
		Msg:        req,
		Connection: udpConn,
	}
}

func (ut *UDPTransport) Close() error {
	return ut.Conn.Close()
}
