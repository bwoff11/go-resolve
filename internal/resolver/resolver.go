package resolver

import (
	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstreams     []upstream.Upstream
	Strategy      config.LoadBalancingStrategy
	LocalCache    *cache.LocalCache
	UpstreamCache *cache.UpstreamCache
}

func New(
	hosts []string,
	strategy config.LoadBalancingStrategy,
	localCache *cache.LocalCache,
	upstreamCache *cache.UpstreamCache,
) *Resolver {

	// Create upstreams from host list
	var upstreams []upstream.Upstream
	for _, host := range hosts {
		upstreams = append(upstreams, *upstream.New(host))
	}

	return &Resolver{
		Upstreams: upstreams,
		Strategy:  strategy,
	}
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
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
