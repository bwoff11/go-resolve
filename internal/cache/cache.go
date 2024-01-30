package cache

import (
	"sync"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	mutex   sync.RWMutex
	Records []Record
}

type Record struct {
	Question dns.Question
	Answer   []dns.RR
	Expiry   time.Time
}

func New(cfg config.Cache) *Cache {
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
	log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("Added record to cache")
	metrics.CacheSize.Set(float64(len(c.Records)))
}

func (c *Cache) Query(q dns.Question) []dns.RR {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, record := range c.Records {
		if record.Question.Name == q.Name && record.Question.Qtype == q.Qtype {
			log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("Found record in cache")
			metrics.CacheHits.Inc()
			return record.Answer
		}
	}

	log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[q.Qtype]).Msg("Record not found in cache")
	metrics.CacheMisses.Inc()

	return []dns.RR{}
}
