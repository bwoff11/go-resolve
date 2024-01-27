package resolver

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
)

type Resolver struct {
	Upstreams   []string
	LocalCache  *cache.Cache
	RemoteCache *cache.Cache
}

func New(upstreams []string, ttl time.Duration, purgeInterval time.Duration) *Resolver {
	return &Resolver{
		Upstreams:   upstreams,
		LocalCache:  cache.New(cache.NoExpiration, cache.NoExpiration),
		RemoteCache: cache.New(ttl, purgeInterval),
	}
}

func (r *Resolver) Resolve(req dns.Msg) (dns.Msg, error) {
	cacheKey := r.generateCacheKey(req.Question[0])

	if response, found := r.checkCaches(cacheKey, req.Id); found {
		return response, nil
	}

	// Query upstream
	response, err := r.queryUpstream(req)
	if err != nil {
		return dns.Msg{}, err
	}

	// Add to remote cache and return
	r.RemoteCache.Set(cacheKey, response, cache.DefaultExpiration)
	return response, nil
}

func (r *Resolver) checkCaches(key string, reqID uint16) (dns.Msg, bool) {
	// Function to update response ID to match request ID
	updateResponseID := func(resp dns.Msg) dns.Msg {
		resp.Id = reqID
		return resp
	}

	// Check local cache
	if x, found := r.LocalCache.Get(key); found {
		return updateResponseID(x.(dns.Msg)), true
	}

	// Check remote cache
	if x, found := r.RemoteCache.Get(key); found {
		return updateResponseID(x.(dns.Msg)), true
	}

	return dns.Msg{}, false
}

func (r *Resolver) generateCacheKey(q dns.Question) string {
	return fmt.Sprintf("%s_%d", q.Name, q.Qtype)
}

func (r *Resolver) createResponse(req dns.Msg, answers []dns.RR) dns.Msg {
	response := dns.Msg{}
	response.SetReply(&req)
	response.Authoritative = false
	response.Answer = answers
	return response
}
