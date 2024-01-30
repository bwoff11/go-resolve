package cache

import (
	"time"

	"github.com/miekg/dns"
)

// Query exposes the cache's query methods.
func (c *Cache) Query(q dns.Question) []dns.RR {
	if r := c.queryCNAME(q); r != nil {
		return r
	}
	return c.queryRecords(q)
}

// queryCNAME checks if the domain points to another domain
// and returns both the RR for the CNAME and the RR for the
// domain it points to.
func (c *Cache) queryCNAME(q dns.Question) []dns.RR {
	var results []dns.RR

	for _, cnameRecord := range c.CNAMERecords {
		if cnameRecord.Question.Name == q.Name && cnameRecord.Question.Qtype == dns.TypeCNAME {
			if time.Now().Before(cnameRecord.ExpiresAt) {
				results = append(results, &cnameRecord.Record)
				// Query for the domain the CNAME points to.
				additionalRecords := c.queryRecords(dns.Question{Name: cnameRecord.Record.Target, Qtype: q.Qtype})
				results = append(results, additionalRecords...)
			}
		}
	}

	return results
}

// queryRecords checks if the domain has any records in the cache.
func (c *Cache) queryRecords(q dns.Question) []dns.RR {
	var results []dns.RR

	for _, domainRecord := range c.DomainRecords {
		if domainRecord.Question.Name == q.Name && domainRecord.Question.Qtype == q.Qtype {
			if time.Now().Before(domainRecord.ExpiresAt) {
				results = append(results, domainRecord.Records...)
			}
		}
	}

	return results
}
