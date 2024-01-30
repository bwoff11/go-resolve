package cache

import (
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

/*
	The cache is responsible for storing local and remote DNS records.

	Questions (keys) consist of a domain name and record type.
	Answers (values) consist of a slice of DNS resource records.
*/

type Cache struct {
	DomainRecords []DomainRecord
	CNAMERecords  []CNAMERecord
}

func New(cfg config.LocalConfig) *Cache {
	c := &Cache{
		DomainRecords: []DomainRecord{},
		CNAMERecords:  []CNAMERecord{},
	}

	purgeInterval := 1 * time.Second
	c.StartHousekeeper(purgeInterval)
	return c
}

// AddRecords accepts one question, and a slice of resource records.
// The RRs are divided into CNAME and domain records, then added to
// their respective caches.
func (c *Cache) AddRecords(q dns.Question, rs []dns.RR) {
	for _, r := range rs {
		switch r.Header().Rrtype {
		case dns.TypeCNAME:
			if cr, ok := r.(*dns.CNAME); ok {
				c.addCNAME(q, cr)
			} else {
				log.Error().Msg("Failed to cast RR to CNAME")
			}
		default:
			c.addDomain(q, r)
		}
	}
}

// addCNAME adds a CNAME record to the cache.
func (c *Cache) addCNAME(q dns.Question, r dns.RR) {
	c.CNAMERecords = append(c.CNAMERecords, CNAMERecord{
		Question:  q,
		ExpiresAt: time.Now().Add(time.Duration(r.Header().Ttl) * time.Second),
		Record:    *r.(*dns.CNAME),
	})
	log.Debug().Str("domain", q.Name).Str("target", r.(*dns.CNAME).Target).Msg("Added CNAME record to cache")
}

// addDomain adds a domain record to the cache.
func (c *Cache) addDomain(q dns.Question, r dns.RR) {
	c.DomainRecords = append(c.DomainRecords, DomainRecord{
		Question:  q,
		ExpiresAt: time.Now().Add(time.Duration(r.Header().Ttl) * time.Second),
		Records:   []dns.RR{r},
	})
	log.Debug().Str("domain", q.Name).Str("type", dns.TypeToString[r.Header().Rrtype]).Msg("Added domain record to cache")
}
