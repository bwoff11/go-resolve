package cache

import (
	"sync"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

/*
	The cache is responsible for storing local and remote DNS records.

	Questions (keys) consist of a domain name and record type.
	Answers (values) consist of a slice of DNS resource records.
*/

type Cache struct {
	mutex   sync.RWMutex
	Records []Record
}

type Record struct {
	Question dns.Question
	Answer   []dns.RR
	Expiry   time.Time
}

func New(cfg config.LocalConfig) *Cache {
	c := &Cache{}

	purgeInterval := 1 * time.Second
	c.StartHousekeeper(purgeInterval)
	return c
}

func (c *Cache) Add(q dns.Question, records []dns.RR) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	ttl := time.Duration(records[0].Header().Ttl) * time.Second

	c.Records = append(c.Records, Record{
		Question: q,
		Answer:   records,
		Expiry:   time.Now().Add(ttl),
	})
}

func (c *Cache) Query(q dns.Question) []dns.RR {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, record := range c.Records {
		if record.Question.Name == q.Name && record.Question.Qtype == q.Qtype {
			return record.Answer
		}
	}
	return []dns.RR{}
}
