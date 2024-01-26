package listener

import (
	"log"
	"net"

	"github.com/miekg/dns"
)

const (
	bufSize = 1024
)

func (l *Listener) listenUDP(addr string) error {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Printf("UDP listener started on %s", addr)
	buf := make([]byte, bufSize)
	for {
		select {
		case <-l.ctx.Done():
			return nil
		default:
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				log.Printf("Error reading UDP data: %v", err)
				continue
			}

			var req dns.Msg
			if err := req.Unpack(buf[:n]); err != nil {
				log.Printf("Error unpacking DNS query: %v", err)
				continue
			}

			resp, err := l.handleDNSQuery(req)
			if err != nil {
				log.Printf("Error resolving DNS query: %v", err)
				continue
			}

			if _, err := conn.WriteTo(resp, addr); err != nil {
				log.Printf("Error sending DNS response: %v", err)
				continue
			}
		}
	}
}
