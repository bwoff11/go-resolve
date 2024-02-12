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

	for {
		req, err := tt.readDNSMessage(conn)
		if err != nil {
			log.Error().Err(err).Msg("error handling TCP connection")
			return
		}

		tt.queueDNSRequest(req, conn)
	}
}

func (tt *TCPTransport) readDNSMessage(conn net.Conn) (*dns.Msg, error) {
	lenBuf := make([]byte, 2)
	_, err := io.ReadFull(conn, lenBuf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint16(lenBuf)
	msgBuf := make([]byte, length)
	_, err = io.ReadFull(conn, msgBuf)
	if err != nil {
		return nil, err
	}

	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		return nil, err
	}

	return &req, nil
}

func (tt *TCPTransport) queueDNSRequest(req *dns.Msg, conn net.Conn) {
	tcpConn := &TCPConnection{Conn: conn}
	tt.Transports.Queue <- QueueItem{
		Msg:        *req,
		Connection: tcpConn,
	}
}

func (t *TCPTransport) Close() error {
	return t.Listener.Close()
}
