package resolver

import (
	"github.com/bwoff11/go-resolve/internal/blocklist"
	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/upstream"
)

type Resolver struct {
	Upstreams []upstream.Upstream
	Strategy  config.LoadBalancingStrategy
	Cache     *cache.Cache
	BlockList *blocklist.BlockList
}

// New creates a new Resolver instance.
func New(cfg *config.Config) *Resolver {
	return &Resolver{
		Upstreams: createUpstreams(cfg),
		Strategy:  cfg.DNS.Upstream.Strategy,
		Cache:     cache.New(),
		BlockList: blocklist.New(cfg.DNS.BlockList),
	}
}

// Helper function to create upstreams.
func createUpstreams(cfg *config.Config) []upstream.Upstream {
	var upstreams []upstream.Upstream
	for _, host := range cfg.DNS.Upstream.Servers {
		upstreams = append(upstreams, *upstream.New(host))
	}
	return upstreams
}
