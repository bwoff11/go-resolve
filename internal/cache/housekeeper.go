package cache

import (
	"time"
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
}
