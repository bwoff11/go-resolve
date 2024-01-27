package listener

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
)

// CreateTCPListener starts a TCP DNS listener on the specified port.
func CreateTCPListener(config *config.Config, resolver *resolver.Resolver) {

	// Create TCP listener
	addr := net.JoinHostPort("", strconv.Itoa(config.DNS.Protocols.TCP.Port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to start TCP listener on %s: %v", addr, err)
	}
	defer listener.Close()

	log.Printf("TCP listener started on %s", addr)

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
		log.Printf("Error reading length: %v", err)
		return
	}
	length := binary.BigEndian.Uint16(lenBuf)

	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		log.Printf("Error reading message: %v", err)
		return
	}

	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		log.Printf("Error unpacking DNS query: %v", err)
		return
	}

	resp, err := res.Resolve(&req)
	if err != nil {
		log.Printf("Error resolving DNS query: %v", err)
		return
	}

	sendTCPResponse(conn, *resp)
}

func sendTCPResponse(conn net.Conn, resp dns.Msg) {
	respBytes, err := resp.Pack()
	if err != nil {
		log.Printf("Error packing DNS response: %v", err)
		return
	}

	lenBuf := []byte{byte(len(respBytes) >> 8), byte(len(respBytes))}
	if _, err := conn.Write(lenBuf); err != nil {
		log.Printf("Error sending length: %v", err)
		return
	}

	if _, err := conn.Write(respBytes); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}
