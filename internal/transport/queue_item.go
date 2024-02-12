package transport

import (
	"github.com/miekg/dns"
)

type QueueItem struct {
	Msg        dns.Msg
	Connection Connection
}

func (qi *QueueItem) Message() *dns.Msg {
	return &qi.Msg
}

func (qi *QueueItem) Question() *dns.Question {
	return &qi.Msg.Question[0]
}

func (qi *QueueItem) Respond(msg *dns.Msg) error {
	return qi.Connection.SendResponse(msg)
}
