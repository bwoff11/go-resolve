package listener

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/miekg/dns"
)

// listenTCP starts a TCP listener on the specified address. It handles incoming DNS queries over TCP.
func (l *Listener) listenTCP(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("TCP listener started on %s", addr)

	for {
		select {
		case <-l.ctx.Done():
			// If a shutdown signal is received, stop the listener
			return nil
		default:
			// Accept a new connection
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting TCP connection: %v", err)
				continue
			}
			// Handle the connection in a separate goroutine
			go l.handleTCPConnection(conn)
		}
	}
}

// handleTCPConnection handles an individual TCP connection.
func (l *Listener) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	// Process the connection to read and respond to the DNS query
	if err := l.processTCPConnection(conn); err != nil {
		log.Printf("Error processing TCP connection: %v", err)
	}
}

// processTCPConnection reads a DNS query from the TCP connection, processes it, and sends back the response.
func (l *Listener) processTCPConnection(conn net.Conn) error {
	// Read the first 2 bytes to determine the length of the DNS query
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return err
	}
	length := binary.BigEndian.Uint16(lenBuf)

	// Read the DNS query based on the obtained length
	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		return err
	}

	// Unpack the query into a dns.Msg struct
	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		return err
	}

	// Handle the DNS query and get the response
	resp, err := l.handleDNSQuery(req)
	if err != nil {
		return err
	}

	// Send the response back to the client
	return l.sendTCPResponse(conn, resp)
}

// sendTCPResponse sends a DNS response back to the client over TCP.
func (l *Listener) sendTCPResponse(conn net.Conn, resp []byte) error {
	// The first 2 bytes contain the length of the response
	lenBuf := []byte{byte(len(resp) >> 8), byte(len(resp))}
	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	// Write the actual DNS response
	_, err := conn.Write(resp)
	return err
}
