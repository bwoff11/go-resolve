package cache

import (
	"fmt"
	"net"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type LocalCache struct {
	cache *cache.Cache
}

// NewLocalCache initializes the local DNS cache with given records.
func NewLocalCache(records config.DNSRecords) *LocalCache {
	c := cache.New(cache.NoExpiration, 10*time.Minute)

	// Populate the cache with A, AAAA, and CNAME records
	for name, ip := range records.A {
		key := createCacheKey(name, dns.TypeA)
		c.Set(key, net.ParseIP(ip), cache.NoExpiration)
	}
	for name, ip := range records.AAAA {
		key := createCacheKey(name, dns.TypeAAAA)
		c.Set(key, net.ParseIP(ip), cache.NoExpiration)
	}
	for alias, target := range records.CNAME {
		key := createCacheKey(alias, dns.TypeCNAME)
		c.Set(key, target, cache.NoExpiration)
	}

	return &LocalCache{cache: c}
}

// createCacheKey generates a cache key based on record name and type.
func createCacheKey(name string, recordType uint16) string {
	return fmt.Sprintf("%s_%d", name, recordType)
}

// QueryA queries for A records in the local cache.
func (lc *LocalCache) QueryA(name string) ([]net.IP, bool) {
	key := createCacheKey(name, dns.TypeA)
	if x, found := lc.cache.Get(key); found {
		return []net.IP{x.(net.IP)}, true
	}
	return nil, false
}

// QueryAAAA queries for AAAA records in the local cache.
func (lc *LocalCache) QueryAAAA(name string) ([]net.IP, bool) {
	key := createCacheKey(name, dns.TypeAAAA)
	if x, found := lc.cache.Get(key); found {
		return []net.IP{x.(net.IP)}, true
	}
	return nil, false
}

// QueryCNAME queries for CNAME records in the local cache.
func (lc *LocalCache) QueryCNAME(name string) (string, bool) {
	key := createCacheKey(name, dns.TypeCNAME)
	if x, found := lc.cache.Get(key); found {
		return x.(string), true
	}
	return "", false
}
