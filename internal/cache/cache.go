package cache

import (
	"fmt"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	Records *cache.Cache
}

func New(cfg config.LocalConfig) *Cache {
	return &Cache{
		Records: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

func (c *Cache) Add(records []dns.RR) {
	for _, record := range records {
		key := createCacheKey(record.Header().Name, record.Header().Rrtype)
		c.Records.Set(key, record, cache.NoExpiration)
	}
}

func (c *Cache) Query(questions []dns.Question) []dns.RR {
	var records []dns.RR
	for _, question := range questions {
		key := createCacheKey(question.Name, question.Qtype)
		if record, found := c.Records.Get(key); found {
			records = append(records, record.(dns.RR))
		}
	}
	return records
}

func createCacheKey(domain string, recordType uint16) string {
	return domain + ":" + dns.TypeToString[recordType]
}

func decodeCacheKey(key string) (string, uint16) {
	var domain string
	var recordType uint16
	_, err := fmt.Sscanf(key, "%s:%s", &domain, &recordType)
	if err != nil {
		return "", 0
	}
	return domain, recordType
}
