package transport

import (
	"net"

	"github.com/miekg/dns"
)

type TCPQueueItem struct {
	Msg  dns.Msg
	Conn net.Conn
}

func (t *TCPQueueItem) Question() dns.Question {
	return t.Msg.Question[0]
}

func (t *TCPQueueItem) Respond(msg []byte) error {
	_, err := t.Conn.Write(msg)
	return err
}
