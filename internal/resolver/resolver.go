package resolver

import (
	"database/sql"
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstreams []upstream.Upstream
	Strategy  config.LoadBalancingStrategy
	Cache     *sql.DB
}

func New(
	hosts []string,
	strategy config.LoadBalancingStrategy,
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
	}
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {

	// check cache here

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
