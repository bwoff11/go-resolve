package cache

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Records   []Record
	Wildcards []Wildcard
}

type Wildcard struct {
	Pattern string
	Type    uint16
	Target  string
}

func New(cfg config.LocalConfig) *Cache {
	return &Cache{}
}

func (c *Cache) Add(rr []dns.RR) {
	for _, r := range rr {
		record, err := FromRR(r)
		if err != nil {
			log.Error().Err(err).Msg("failed to convert RR to cache record")
			return
		}
		c.Records = append(c.Records, *record)
	}
}

func (c *Cache) Query(domain string, recordType uint16) []dns.RR {
	var values []dns.RR
	for _, record := range c.Records {
		if record.Domain == domain && record.Type == recordType {
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
		if record.Domain == domain && record.Type == recordType {
			c.Records = append(c.Records[:i], c.Records[i+1:]...)
			return true
		}
	}
	return false
}
