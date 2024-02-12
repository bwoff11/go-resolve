package transport

import (
	"net"

	"github.com/miekg/dns"
)

type Connection interface {
	SendResponse(msg *dns.Msg) error
}

type UDPConnection struct {
	Addr net.Addr
	Conn net.PacketConn
}

func (uc *UDPConnection) SendResponse(msg *dns.Msg) error {
	data, err := msg.Pack()
	if err != nil {
		return err
	}
	_, err = uc.Conn.WriteTo(data, uc.Addr)
	return err
}

type TCPConnection struct {
	Conn net.Conn
}

func (tc *TCPConnection) SendResponse(msg *dns.Msg) error {
	data, err := msg.Pack()
	if err != nil {
		return err
	}
	_, err = tc.Conn.Write(data)
	return err
}
