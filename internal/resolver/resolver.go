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

	// Try cache
	if r.Cache != nil {
		if answer, auth, err := r.Cache.Query(question); len(answer) > 0 {
			log.Printf("Cache hit for %s", question.Name)
			return r.requestToResponse(req, answer, auth), nil
		} else if err != nil {
			log.Printf("Failed to query cache: %v", err)
			return nil, err
		}
	}

	// Try upstream
	upstream := r.selectUpstream()
	if msg, err := upstream.Query(req); err == nil {
		log.Printf("Upstream hit for %s", question.Name)
		r.Cache.Add(msg.Answer)
		return r.requestToResponse(req, msg.Answer, false), nil
	}

	// Return NXDOMAIN
	return r.requestToResponse(req, []dns.RR{}, false), nil
}

func (r *Resolver) selectUpstream() *upstream.Upstream {
	switch r.Strategy {
	case config.LoadBalancingStrategyLatency:
		return &r.Upstreams[0]
	case config.LoadBalancingStrategyRandom:
		return &r.Upstreams[0]
	case config.LoadBalancingStrategyRoundRobin:
	}
	return &r.Upstreams[0]
}

// Responsible for converting a request to a response by
// modifying all appropriate fields.
func (r *Resolver) requestToResponse(req *dns.Msg, answer []dns.RR, authoritative bool) *dns.Msg {
	return &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:                 req.Id,               // ID copied from request
			Response:           true,                 // This is a response
			Opcode:             req.Opcode,           // Type of query set by client and copied to response
			Authoritative:      authoritative,        // Whether the response is authoritative. Only true for local records
			Truncated:          false,                // Whether the response was truncated. This should never be true.
			RecursionDesired:   req.RecursionDesired, // Whether the client wants recursion
			RecursionAvailable: true,                 // Whether recursion is available. Need to look into this.
			Rcode:              dns.RcodeSuccess,     // Status of the response. 0 = success
		},
		Compress: false, // Whether to compress the response. This should never be true.
		Question: req.Question,
		Answer:   answer,
		Ns:       []dns.RR{}, // todo
		Extra:    []dns.RR{}, // todo
	}
}
