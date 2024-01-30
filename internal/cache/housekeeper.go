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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var count int
	var newRecords []Record
	for _, record := range c.Records {
		if record.Expiry.After(time.Now()) {
			newRecords = append(newRecords, record)
		} else {
			count++
		}
	}
	c.Records = newRecords

	if count > 0 {
		log.Info().Int("count", count).Msg("expired records removed from cache")
	}
}
