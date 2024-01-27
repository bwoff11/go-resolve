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
	// bufSize defines the size of the buffer used to read incoming UDP packets.
	bufSize = 1024 // Adjust this size based on expected query sizes
)

// CreateUDPListener starts a UDP DNS listener on the specified port.
func CreateUDPListener(
	protoConfig config.ProtocolConfig,
	cacheConfig config.CacheConfig,
	upstreamConfig config.UpstreamConfig,
) {

	// Create resolver
	res := resolver.New(
		upstreamConfig.Servers,
		cacheConfig.TTL,
		cacheConfig.PurgeInterval,
	)

	// Create UDP listener
	addr := net.JoinHostPort("", strconv.Itoa(protoConfig.Port))
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to start UDP listener on %s: %v", addr, err)
	}
	defer conn.Close()

	log.Printf("UDP listener started on %s", addr)

	buf := make([]byte, bufSize)

	// Accept connections
	for {
		n, clientAddr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Printf("Error reading UDP data: %v", err)
			continue
		}

		var req dns.Msg
		if err := req.Unpack(buf[:n]); err != nil {
			log.Printf("Error unpacking DNS query: %v", err)
			continue
		}

		resp, err := res.Resolve(req)
		if err != nil {
			log.Printf("Error resolving DNS query: %v", err)
			continue
		}

		respBytes, err := resp.Pack()
		if err != nil {
			log.Printf("Error packing DNS response: %v", err)
			continue
		}

		if _, err := conn.WriteTo(respBytes, clientAddr); err != nil {
			log.Printf("Error sending DNS response: %v", err)
		}
	}
}
