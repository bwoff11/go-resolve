package listener

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"github.com/miekg/dns"
)

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
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("Error accepting TCP connection: %v", err)
				continue
			}
			go l.handleTCPConnection(conn)
		}
	}
}

func (l *Listener) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	if err := l.processTCPConnection(conn); err != nil {
		log.Printf("Error processing TCP connection: %v", err)
	}
}

func (l *Listener) processTCPConnection(conn net.Conn) error {
	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, lenBuf); err != nil {
		return err
	}
	length := binary.BigEndian.Uint16(lenBuf)

	msgBuf := make([]byte, length)
	if _, err := io.ReadFull(conn, msgBuf); err != nil {
		return err
	}

	var req dns.Msg
	if err := req.Unpack(msgBuf); err != nil {
		return err
	}

	resp, err := l.handleDNSQuery(req)
	if err != nil {
		return err
	}

	return l.sendTCPResponse(conn, resp)
}

func (l *Listener) sendTCPResponse(conn net.Conn, resp []byte) error {
	lenBuf := []byte{byte(len(resp) >> 8), byte(len(resp))}
	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	_, err := conn.Write(resp)
	return err
}
