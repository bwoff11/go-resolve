package cache

import (
	"time"

	"github.com/miekg/dns"
)

type CNAMERecord struct {
	Question  dns.Question
	ExpiresAt time.Time
	Record    dns.CNAME
}

func (cr *CNAMERecord) IsExpired() bool {
	return time.Now().After(cr.ExpiresAt)
}

func (cr *CNAMERecord) Query(q dns.Question) []dns.RR {
	qn := q.Name
	qt := q.Qtype

	rn := cr.Question.Name
	rt := cr.Question.Qtype

	if qn == rn && qt == rt && !cr.IsExpired() {
		return []dns.RR{&cr.Record}
	}

	return nil
}
