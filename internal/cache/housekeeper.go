package cache

import (
	"time"

	"github.com/rs/zerolog/log"
)

func (c *Cache) StartHousekeeper(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				c.RemoveExpired()
			}
		}
	}()
}

func (c *Cache) RemoveExpired() {
	now := time.Now()
	rCount := 0
	cCount := 0
	for i, d := range c.DomainRecords {
		if d.ExpiresAt.Before(now) {
			c.DomainRecords = append(c.DomainRecords[:i], c.DomainRecords[i+1:]...)
			rCount += len(d.Records)
			cCount++
		}
	}
	for i, d := range c.CNAMERecords {
		if d.ExpiresAt.Before(now) {
			c.CNAMERecords = append(c.CNAMERecords[:i], c.CNAMERecords[i+1:]...)
			cCount++
		}
	}
	if rCount > 0 || cCount > 0 {
		log.Info().Int("records", rCount).Int("cnames", cCount).Msg("Removed expired records from cache")
	}
}
