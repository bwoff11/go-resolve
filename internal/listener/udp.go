package listener

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

const (
	// bufSize defines the size of the buffer used to read incoming UDP packets.
	// 1024 bytes are typically enough for most DNS queries, but this may need
	// to be increased for handling larger queries, like those used with DNSSEC.
	bufSize = 1024
)

// listenUDP starts a UDP listener on the specified address. It handles incoming DNS queries over UDP.
func (l *Listener) listenUDP(addr string) error {
	// Open a UDP network connection
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("UDP listener started on %s", addr)

	// Create a buffer to read incoming packets
	buf := make([]byte, bufSize)

	for {
		select {
		case <-l.ctx.Done():
			// If a shutdown signal is received, stop the listener
			return nil
		default:
			// Read a packet from the connection
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				log.Printf("Error reading UDP data: %v", err)
				continue
			}

			// Unpack the packet into a dns.Msg struct
			var req dns.Msg
			if err := req.Unpack(buf[:n]); err != nil {
				log.Printf("Error unpacking DNS query: %v", err)
				continue
			}

			// Handle the DNS query and get the response
			resp, err := l.handleDNSQuery(req)
			if err != nil {
				log.Printf("Error resolving DNS query: %v", err)
				continue
			}

			// Send the response back to the client
			if _, err := conn.WriteTo(resp, addr); err != nil {
				log.Printf("Error sending DNS response: %v", err)
				continue
			}
		}
	}
}
