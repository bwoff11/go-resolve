package cache

import (
	"log"

	"github.com/miekg/dns"
)

// Query exposes the cache's query methods.
func (c *Cache) Query(q dns.Question) []dns.RR {
	if r := c.queryCNAME(q); len(r) > 0 {
		return r
	}
	if r := c.queryDomain(q); len(r) > 0 {
		return r
	}
	return nil
}

func (c *Cache) queryCNAME(q dns.Question) []dns.RR {
	var results []dns.RR
	for _, cr := range c.CNAMERecords {
		if r := cr.Query(q); len(r) > 0 {
			results = append(results, r...)
			newQuestion := dns.Question{
				Name:   cr.Record.Target,
				Qtype:  q.Qtype,
				Qclass: q.Qclass,
			}
			log.Println("newQuestion", newQuestion)
			if r := c.queryDomain(newQuestion); len(r) > 0 {
				results = append(results, r...)
			}
		}
	}
	return results
}

func (c *Cache) queryDomain(q dns.Question) []dns.RR {
	for _, dr := range c.DomainRecords {
		if r := dr.Query(q); len(r) > 0 {
			return r
		}
	}
	return nil
}
