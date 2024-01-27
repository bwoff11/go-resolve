package listener

import (
	"log"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
)

const (
	bufSize = 1024 // Size of the buffer to read UDP packets
)

// CreateUDPListener starts a UDP DNS listener on the specified port.
func CreateUDPListener(config *config.Config) {

	// Create resolver
	resolver := resolver.New(config.DNS.Upstream.Servers)

	addr := net.JoinHostPort("", strconv.Itoa(config.DNS.Protocols.UDP.Port))
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to start UDP listener on %s: %v", addr, err)
	}
	defer conn.Close()

	log.Printf("UDP listener started on %s", addr)
	handleUDPConnections(conn, resolver)
}

// handleUDPConnections listens for incoming UDP packets and processes them.
func handleUDPConnections(conn net.PacketConn, res *resolver.Resolver) {
	buf := make([]byte, bufSize)

	for {
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading UDP data: %v", err)
			continue
		}

		go processUDPQuery(buf[:n], conn, clientAddr, res)
	}
}

// processUDPQuery processes a single UDP query and sends back a response.
func processUDPQuery(query []byte, conn net.PacketConn, addr net.Addr, res *resolver.Resolver) {
	var req dns.Msg
	if err := req.Unpack(query); err != nil {
		log.Printf("Error unpacking DNS query: %v", err)
		return
	}

	resp, err := res.Resolve(&req)
	if err != nil {
		log.Printf("Error resolving DNS query: %v", err)
		return
	}

	respBytes, err := resp.Pack()
	if err != nil {
		log.Printf("Error packing DNS response: %v", err)
		return
	}

	if _, err := conn.WriteTo(respBytes, addr); err != nil {
		log.Printf("Error sending DNS response: %v", err)
	}
}
