package cache

import (
	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// Cache stores DNS records and provides query functionality.
type Cache struct {
	Records   []common.LocalRecord // Records stores standard DNS records.
	Wildcards []Wildcard           // Wildcards stores wildcard DNS records.
}

// Wildcard represents a wildcard DNS record.
type Wildcard struct {
	Pattern string // Pattern is the wildcard pattern (e.g., "*.example.com").
	Type    string // Type is the DNS record type (e.g., "A", "AAAA").
	Target  string // Target is the value associated with the wildcard record.
}

// New creates and initializes a new Cache instance.
func New(cfg config.LocalConfig) *Cache {
	log.Debug().Int("size", len(cfg.Records)).Msg("Creating new cache")

	cache := &Cache{
		Records:   []common.LocalRecord{},
		Wildcards: []Wildcard{},
	}

	cache.addLocalRecords(cfg.Records)
	return cache
}

// addLocalRecords adds local records to the cache.
func (c *Cache) addLocalRecords(records []common.LocalRecord) {
	for _, r := range records {
		r.Domain = ensureTrailingDot(r.Domain)

		if isWildcard(r.Domain) {
			c.addWildcard(r)
		} else {
			c.Records = append(c.Records, r)
			log.Debug().Str("domain", r.Domain).Str("type", r.Type).Msg("Added local record to cache")
		}
	}
}

// addWildcard processes and adds a wildcard record to the cache.
func (c *Cache) addWildcard(record common.LocalRecord) {
	c.Wildcards = append(c.Wildcards, Wildcard{
		Pattern: convertToSQLPattern(record.Domain),
		Type:    record.Type,
		Target:  record.Value[0],
	})
	log.Debug().Str("pattern", record.Domain).Str("target", record.Value[0]).Msg("Added wildcard to cache")
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
	var values []dns.RR
	for _, record := range c.Records {
		if record.Domain == domain && record.Type == dns.TypeToString[recordType] {
			if rr, err := record.ToRR(); err == nil {
				values = append(values, rr)
			} else {
				log.Error().Err(err).Msg("Failed to convert cache record to RR")
			}
		}
	}
	return values
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
