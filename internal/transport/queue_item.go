package transport

import "github.com/miekg/dns"

type QueueItem interface {
	Message() *dns.Msg
	Question() *dns.Question
	Protocol() string
	Respond(*dns.Msg) error
}
