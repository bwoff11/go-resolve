package transport

import (
	"net"

	"github.com/miekg/dns"
)

type UDPQueueItem struct {
	Msg  dns.Msg
	Addr net.Addr
	Conn net.PacketConn
}

func (u *UDPQueueItem) Message() *dns.Msg {
	return &u.Msg
}

func (u *UDPQueueItem) Question() *dns.Question {
	return &u.Msg.Question[0]
}

func (u *UDPQueueItem) Protocol() string {
	return "udp"
}

// Respond serializes the dns.Msg and sends it to the client.
func (u *UDPQueueItem) Respond(msg *dns.Msg) error {

	// Serialize the dns.Msg into wire format.
	data, err := msg.Pack()
	if err != nil {
		return err
	}

	// Send the serialized message to the client.
	_, err = u.Conn.WriteTo(data, u.Addr)
	return err
}
