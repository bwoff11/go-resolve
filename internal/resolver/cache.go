package resolver

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type CacheEntry struct {
	Response  *dns.Msg
	ExpiresAt time.Time
}

type Cache struct {
	mutex   sync.RWMutex
	entries map[string]*CacheEntry
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]*CacheEntry),
	}
}

func (c *Cache) Add(key string, response *dns.Msg, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = &CacheEntry{
		Response:  response,
		ExpiresAt: time.Now().Add(ttl),
	}
}

func (c *Cache) Get(key string) (*dns.Msg, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, found := c.entries[key]
	if !found || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}
	return entry.Response, true
}

// Add a method to periodically clean up expired entries.
