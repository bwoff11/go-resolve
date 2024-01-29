package cache

import (
	"fmt"
	"time"

	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Records *cache.Cache
}

func New(cfg config.LocalConfig) *Cache {
	c := &Cache{
		Records: cache.New(cache.NoExpiration, 60*time.Second),
	}
	c.AddLocalRecord(cfg.Records)
	return c
}

func (c *Cache) Add(records []dns.RR) {
	for _, record := range records {
		ttl := time.Duration(record.Header().Ttl) * time.Second // Convert TTL to time.Duration
		key := createCacheKey(record.Header().Name, record.Header().Rrtype)
		c.Records.Set(key, record, ttl) // Set the record in the cache with the appropriate TTL
		log.Debug().Str("key", key).Int("ttl", int(ttl.Seconds())).Msg("Added record to cache")
	}
}

func (c *Cache) AddLocalRecord(records []common.LocalRecord) error {
	log.Debug().Msg("Adding local records to cache")
	for _, record := range records {
		// Convert LocalRecord to dns.RR
		rr, err := record.ToRR()
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert LocalRecord to dns.RR")
			return err
		}

		// Create a cache key based on the domain and record type
		key := createCacheKey(rr.Header().Name, rr.Header().Rrtype)

		// Set the record in the cache
		c.Records.Set(key, rr, cache.NoExpiration)
		log.Debug().Str("key", key).Msg("Added local record to cache")
	}
	return nil
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
