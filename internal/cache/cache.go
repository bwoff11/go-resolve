package cache

import (
	"sync"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	sync.Mutex
	nextID  uint64
	Records []dns.RR
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) Add(records []dns.RR) {
	c.Lock()
	defer c.Unlock()

	if c.Records == nil {
		c.Records = []dns.RR{}
	}

	for _, record := range records {
		log.Debug().
			Str("msg", "Adding record to cache").
			Str("domain", record.Header().Name).
			Str("type", dns.TypeToString[record.Header().Rrtype]).
			Str("value", getRRValue(record)).
			Int("ttl", int(record.Header().Ttl)).
			Send()
		c.Records = append(c.Records, record)
	}
}

func (c *Cache) Query(question dns.Question) ([]dns.RR, bool) {
	c.Lock()
	defer c.Unlock()

	if c.Records == nil {
		return nil, false
	}

	var records []dns.RR
	for _, record := range c.Records {
		if record.Header().Name == question.Name && record.Header().Rrtype == question.Qtype {
			log.Debug().
				Str("msg", "Found record in cache").
				Str("domain", record.Header().Name).
				Str("type", dns.TypeToString[record.Header().Rrtype]).
				Str("value", getRRValue(record)).
				Int("ttl", int(record.Header().Ttl)).
				Send()
			records = append(records, record)
		}
	}

	return records, len(records) > 0
}
