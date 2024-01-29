package cache

import (
	"fmt"
	"time"

	"github.com/bwoff11/go-resolve/internal/common"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	Records *cache.Cache
	CNAMEs  *cache.Cache
}

func New(cfg config.LocalConfig) *Cache {
	c := &Cache{
		Records: cache.New(10*time.Minute, 60*time.Second),
	}
	c.AddLocalRecord(cfg.Records)
	return c
}

func (c *Cache) Add(records []dns.RR) {
	for _, record := range records {
		ttl := time.Duration(record.Header().Ttl) * time.Second // Convert TTL to time.Duration
		key := createCacheKey(record.Header().Name, record.Header().Rrtype)
		c.Records.Set(key, record, ttl) // Set the record in the cache with the appropriate TTL
		log.Debug().Str("key", key).Str("record", record.String()).Int("ttl", int(ttl.Seconds())).Msg("Added record to cache")
	}
}

func (c *Cache) AddLocalRecord(lr []common.LocalRecord) error {
	log.Debug().Msg("Adding local records to cache")
	var records []dns.RR
	for _, r := range lr {
		// Convert LocalRecord to dns.RR
		rr, err := r.ToRR()
		if err != nil {
			log.Error().Err(err).Msg("Failed to convert LocalRecord to dns.RR")
			return err
		}
		records = append(records, rr)
	}
	c.Add(records)
	return nil
}

func (c *Cache) Query(questions []dns.Question) []dns.RR {
	startTime := time.Now()
	log.Debug().Str("questions", fmt.Sprintf("%v", questions)).Msg("Querying cache")
	var records []dns.RR

	for _, question := range questions {
		// First, check for a CNAME record
		cnameKey := createCacheKey(question.Name, dns.TypeCNAME)
		log.Debug().Str("cnameKey", cnameKey).Msg("Checking cache for CNAME record")

		if cnameRecord, found := c.Records.Get(cnameKey); found {
			log.Debug().Str("cnameKey", cnameKey).Msg("Found CNAME record in cache")
			cnameRR := cnameRecord.(dns.RR)

			// Append the CNAME record and then resolve its target
			records = append(records, cnameRR)
			targetRecords := c.resolveCNAME(cnameRR.(*dns.CNAME).Target, question.Qtype)
			records = append(records, targetRecords...)

		} else {
			// If no CNAME record, check for the requested type
			key := createCacheKey(question.Name, question.Qtype)
			log.Debug().Str("key", key).Msg("Checking cache for record")

			if record, found := c.Records.Get(key); found {
				log.Debug().Str("key", key).Msg("Found record in cache")
				rr := record.(dns.RR)
				records = append(records, rr)
			}
		}
	}

	if len(records) == 0 {
		log.Debug().Msg("No records found in cache")
	}

	metrics.CacheDuration.Observe(time.Since(startTime).Seconds())
	return records
}

// resolveCNAME resolves the target of a CNAME record.
func (c *Cache) resolveCNAME(target string, qtype uint16) []dns.RR {
	var result []dns.RR
	key := createCacheKey(target, qtype)
	if record, found := c.Records.Get(key); found {
		result = append(result, record.(dns.RR))
	}
	return result
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
