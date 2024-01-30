package listener

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// CreateTCPListener starts a TCP DNS listener on the specified port.
func CreateTCPListener(config *config.Config, resolver *resolver.Resolver) {

	// Create TCP listener
	addr := net.JoinHostPort("", strconv.Itoa(config.Protocols.TCP.Port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating TCP listener")
	}
	defer listener.Close()

	log.Debug().Str("protocol", "tcp").Str("address", addr).Msg("Listening")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %v", err)
			continue
		}
		go handleTCPConnection(conn, resolver)
	}
}

func handleTCPConnection(conn net.Conn, res *resolver.Resolver) {
	defer conn.Close()
	processTCPConnection(conn, res)
}

func processTCPConnection(conn net.Conn, res *resolver.Resolver) {
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error reading length")
		return
	}
	length := binary.BigEndian.Uint16(lenBuf)

	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error reading message")
		return
	}

	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error unpacking DNS message")
		return
	}

	resp, err := res.Resolve(&req)
	if err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error resolving DNS message")
		return
	}

	sendTCPResponse(conn, *resp)
}

func sendTCPResponse(conn net.Conn, resp dns.Msg) {
	respBytes, err := resp.Pack()
	if err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error packing DNS message")
		return
	}

	lenBuf := []byte{byte(len(respBytes) >> 8), byte(len(respBytes))}
	if _, err := conn.Write(lenBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error sending length")
		return
	}

	if _, err := conn.Write(respBytes); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("Error sending message")
		return
	}
}
