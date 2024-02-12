package transport

import (
	"net"

	"github.com/miekg/dns"
)

type TCPQueueItem struct {
	Msg  dns.Msg
	Conn net.Conn
}

func (t *TCPQueueItem) Message() *dns.Msg {
	return &t.Msg
}

func (t *TCPQueueItem) Question() *dns.Question {
	return &t.Msg.Question[0]
}

func (t *TCPQueueItem) Protocol() string {
	return "tcp"
}

func (t *TCPQueueItem) Respond(msg *dns.Msg) error {

	// Convert to wire format
	msgBuf, err := msg.Pack()
	if err != nil {
		return err
	}

	// Write the message to the connection
	_, err = t.Conn.Write(msgBuf)
	return err
}
