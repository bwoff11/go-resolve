package cache

import (
	"strings"

	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Records   []common.LocalRecord
	Wildcards []Wildcard
}

type Wildcard struct {
	Pattern string
	Type    string
	Target  string
}

func New(cfg config.LocalConfig) *Cache {
	log.Debug().Int("size", len(cfg.Records)).Msg("creating new cache")
	c := &Cache{
		Records:   []common.LocalRecord{},
		Wildcards: []Wildcard{},
	}
	c.addLocalRecords(cfg.Records)
	return c
}

func (c *Cache) addLocalRecords(records []common.LocalRecord) {
	for _, r := range records {
		// Add a "." to the end of the domain if it doesn't already have one
		if !strings.HasSuffix(r.Domain, ".") {
			r.Domain = r.Domain + "."
		}
		if !strings.Contains(r.Domain, "*") {
			c.Records = append(c.Records, r)
		} else {
			c.addWildcard(r)
		}
		log.Debug().Str("domain", r.Domain).Str("type", r.Type).Msg("added local record to cache")
	}
}

func (c *Cache) addWildcard(record common.LocalRecord) {
	c.Wildcards = append(c.Wildcards, Wildcard{
		Pattern: strings.Replace(record.Domain, "*", "%", 1),
		Type:    record.Type,
		Target:  record.Value[0],
	})
	log.Debug().Str("pattern", record.Domain).Str("target", record.Value[0]).Msg("added wildcard to cache")
}

func (c *Cache) Add(rr []dns.RR) {
	for _, r := range rr {
		var record common.LocalRecord
		err := record.FromRR(r)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert RR to cache record")
			return
		}
		c.Records = append(c.Records, record)
		log.Debug().Str("domain", record.Domain).Str("type", record.Type).Msg("added record to cache")
	}
}

func (c *Cache) Query(domain string, recordType uint16) []dns.RR {
	var values []dns.RR
	for _, record := range c.Records {
		if record.Domain == domain && record.Type == dns.TypeToString[recordType] {
			rr, err := record.ToRR()
			if err != nil {
				log.Error().Err(err).Msg("failed to convert cache record to RR")
				continue
			}
			values = append(values, rr)
		}
	}
	return values
}

func (c *Cache) Remove(domain string, recordType uint16) bool {
	for i, record := range c.Records {
		if record.Domain == domain && record.Type == dns.TypeToString[recordType] {
			c.Records = append(c.Records[:i], c.Records[i+1:]...)
			return true
		}
	}
	return false
}
