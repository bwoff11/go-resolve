package transport

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type TCPTransport struct {
	Transports *Transports
	Listener   net.Listener
}

func NewTCP(c config.Protocol, ts *Transports) (*TCPTransport, error) {
	addr := net.JoinHostPort("", strconv.Itoa(c.Port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TCPTransport{
		Transports: ts,
		Listener:   listener,
	}, nil
}

func (t *TCPTransport) Listen() error {

	go func() {
		for {
			conn, err := t.Listener.Accept()
			if err != nil {
				log.Error().Err(err).Msg("error accepting tcp connection")
				return
			}
			go t.handleTCPConnection(conn)
		}
	}()

	log.Info().Str("protocol", "tcp").Msg("transport listening")
	return nil
}

func (tt *TCPTransport) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	lenBuf := make([]byte, 2)
	for {
		// Read the length of the DNS message
		_, err := io.ReadFull(conn, lenBuf)
		if err != nil {
			log.Error().Err(err).Msg("error reading length from TCP connection")
			return
		}
		length := binary.BigEndian.Uint16(lenBuf)

		// Read the DNS message based on the length
		msgBuf := make([]byte, length)
		_, err = io.ReadFull(conn, msgBuf)
		if err != nil {
			log.Error().Err(err).Msg("error reading message from TCP connection")
			return
		}

		// Unpack the DNS message
		var req dns.Msg
		if err := req.Unpack(msgBuf); err != nil {
			log.Error().Err(err).Msg("error unpacking DNS message")
			return
		}

		// Create a TCPConnection adapter
		tcpConn := &TCPConnection{Conn: conn}

		// Queue the item with the generic QueueItem structure
		tt.Transports.Queue <- QueueItem{
			Msg:        req,
			Connection: tcpConn,
		}
	}
}

// Implement processing for OutboundQueue, sending responses back to the client.
// This functionality would require a mapping or association of responses to client connections,
// potentially via a modified QueueItem struct to include a connection reference for TCP responses.

func (t *TCPTransport) Close() error {
	return t.Listener.Close()
}
