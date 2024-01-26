package listener

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/models"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/miekg/dns"
)

type Listener struct {
	Resolver *resolver.Resolver
	Protocol models.Protocol
	Port     int
}

func New(resolver *resolver.Resolver, protocol models.Protocol, port int) *Listener {
	log.Printf("Creating listener on port %d for protocol %s\n", port, protocol)
	return &Listener{
		Resolver: resolver,
		Protocol: protocol,
		Port:     port,
	}
}

func (l *Listener) Listen() error {
	addr := net.JoinHostPort("", strconv.Itoa(l.Port))
	switch l.Protocol {
	case models.TCP:
		return l.listenTCP(addr)
	case models.UDP:
		return l.listenUDP(addr)
	default:
		return fmt.Errorf("unknown protocol: %s", l.Protocol)
	}
}

func (l *Listener) listenTCP(addr string) error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		_, err := ln.Accept()
		if err != nil {
			// handle error
		}
	}
}

func (l *Listener) listenUDP(addr string) error {
	ln, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	log.Printf("Listening on %s\n", addr)
	for {
		buf := make([]byte, 512)
		n, addr, err := ln.ReadFrom(buf)
		if err != nil {
			return err
		}

		log.Printf("Received %d bytes from %s\n", n, addr)
		log.Printf("Data: %s\n", buf[:n])

		// Packaged the received data into a DNS message...
		var req dns.Msg
		err = req.Unpack(buf[:n])
		if err != nil {
			log.Printf("Error unpacking DNS query: %v\n", err)
			continue
		}

		// Resolve the DNS query...
		msg, err := l.Resolver.Resolve(req)
		if err != nil {
			log.Printf("Error resolving DNS query: %v\n", err)
			continue
		}

		// Convert the response to a byte slice...
		resp, err := msg.Pack()
		if err != nil {
			log.Printf("Error packing DNS response: %v\n", err)
			continue
		}

		// Send the response...
		_, err = ln.WriteTo(resp, addr)
		if err != nil {
			log.Printf("Error sending DNS response: %v\n", err)
			continue
		}

		log.Printf("Sent %d bytes to %s\n", len(resp), addr)
	}
}
