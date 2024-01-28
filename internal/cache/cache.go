package cache

import (
	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// Cache stores DNS records and provides query functionality.
type Cache struct {
	Records []common.LocalRecord // Records stores standard DNS records.
}

// New creates and initializes a new Cache instance.
func New(cfg config.LocalConfig) *Cache {
	log.Debug().Int("size", len(cfg.Records)).Msg("Creating new cache")

	cache := &Cache{
		Records: []common.LocalRecord{},
	}

	return cache
}

// Add inserts new records into the cache.
func (c *Cache) Add(rr []dns.RR) {
	for _, r := range rr {
		record := common.LocalRecord{}
		if err := record.FromRR(r); err != nil {
			log.Error().Err(err).Msg("Failed to convert RR to cache record")
			continue
		}
		c.Records = append(c.Records, record)
		log.Debug().Str("domain", record.Domain).Str("type", record.Type).Msg("Added record to cache")
	}
}

// Query searches for records matching the domain and record type.
func (c *Cache) Query(domain string, recordType uint16) []dns.RR {
	var results []dns.RR
	for _, record := range c.Records {
		if record.Domain == domain && record.Type == dns.TypeToString[recordType] {
			if rr, err := record.ToRR(); err == nil {
				results = append(results, rr)
			}
		}
	}
	return results
}

// Remove deletes a record from the cache.
func (c *Cache) Remove(domain string, recordType uint16) bool {
	for i, record := range c.Records {
		if record.Domain == domain && record.Type == dns.TypeToString[recordType] {
			c.Records = append(c.Records[:i], c.Records[i+1:]...)
			return true
		}
	}
	return false
}
