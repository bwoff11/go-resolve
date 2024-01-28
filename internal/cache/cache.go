package cache

import (
	"sync"

	"github.com/miekg/dns"
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
			records = append(records, record)
		}
	}

	return records, len(records) > 0
}
