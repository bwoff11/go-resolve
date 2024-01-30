package cache

import (
	"time"

	"github.com/miekg/dns"
)

type DomainRecord struct {
	Question  dns.Question
	ExpiresAt time.Time
	Records   []dns.RR
}

func (dr *DomainRecord) IsExpired() bool {
	return time.Now().After(dr.ExpiresAt)
}

func (dr *DomainRecord) Query(q dns.Question) []dns.RR {
	qn := q.Name
	qt := q.Qtype

	rn := dr.Question.Name
	rt := dr.Question.Qtype

	if qn == rn && qt == rt && !dr.IsExpired() {
		return dr.Records
	}

	return nil
}
