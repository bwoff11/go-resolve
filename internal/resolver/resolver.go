package resolver

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstreams []upstream.Upstream
	Strategy  config.LoadBalancingStrategy
	Cache     *cache.Cache
}

func New(
	hosts []string,
	strategy config.LoadBalancingStrategy,
	cache *cache.Cache,
) *Resolver {

	log.Printf("Creating resolver with strategy: %s\n", strategy)

	// Create upstreams from host list
	var upstreams []upstream.Upstream
	for _, host := range hosts {
		upstreams = append(upstreams, *upstream.New(host))
	}

	return &Resolver{
		Upstreams: upstreams,
		Strategy:  strategy,
		Cache:     cache,
	}
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {

	question := req.Question[0]
	name := question.Name
	recordType := question.Qtype

	// Check the cache
	if r.Cache != nil {
		record, err := r.Cache.Query(name, recordType)
		if err == nil && record != nil {
			// Cache hit
			log.Printf("Cache hit for %s\n", name)
			return record.ToDNSMsg(int(req.Id)), nil
		} else if err != nil {
			log.Printf("Error querying cache: %v\n", err)
			// Optionally handle cache query error
		}
		// Optionally handle cache miss
	}

	upstream := r.selectUpstream()
	resp, err := upstream.Query(req)
	if err != nil {
		return nil, err
	}

	// Set the response ID to match the request ID
	resp.Id = req.Id

	return resp, nil
}

func (r *Resolver) selectUpstream() *upstream.Upstream {
	return &r.Upstreams[0]
}
