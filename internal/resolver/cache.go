package resolver

import (
	"sync"
	"time"

	"github.com/miekg/dns"
)

type CacheEntry struct {
	Response *dns.Msg
	CachedAt time.Time
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

func (c *Cache) Add(key string, response *dns.Msg) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries[key] = &CacheEntry{
		Response: response,
		CachedAt: time.Now(),
	}
}

func (c *Cache) Get(key string) (*dns.Msg, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, found := c.entries[key]
	if !found {
		return nil, false
	}
	ttl := entry.Response.Answer[0].Header().Ttl
	if time.Since(entry.CachedAt) > time.Duration(ttl)*time.Second {
		delete(c.entries, key)
		return nil, false
	}
	return entry.Response, true
}
