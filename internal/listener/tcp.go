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
		log.Fatal().Err(err).Msg("error creating tcp listener")
	}
	defer listener.Close()

	log.Debug().Str("protocol", "tcp").Str("address", addr).Msg("listening")

	// Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting TCP connection: %v", err)
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
		log.Error().Err(err).Str("protocol", "tcp").Msg("error reading length")
		return
	}
	length := binary.BigEndian.Uint16(lenBuf)

	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error reading message")
		return
	}

	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error unpacking dns message")
		return
	}

	resp, err := res.Resolve(&req)
	if err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error resolving dns message")
		return
	}

	sendTCPResponse(conn, *resp)
}

func sendTCPResponse(conn net.Conn, resp dns.Msg) {
	respBytes, err := resp.Pack()
	if err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error packing dns message")
		return
	}

	lenBuf := []byte{byte(len(respBytes) >> 8), byte(len(respBytes))}
	if _, err := conn.Write(lenBuf); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error sending length")
		return
	}

	if _, err := conn.Write(respBytes); err != nil {
		log.Error().Err(err).Str("protocol", "tcp").Msg("error sending message")
		return
	}
}
