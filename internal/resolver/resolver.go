package resolver

import (
	"github.com/bwoff11/go-resolve/internal/blocklist"
	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Resolver struct {
	Upstreams []upstream.Upstream
	Strategy  config.LoadBalancingStrategy
	Cache     *cache.Cache
	BlockList *blocklist.BlockList
}

func New(cfg *config.Config) *Resolver {

	// Create upstreams from host list
	var upstreams []upstream.Upstream
	for _, host := range cfg.DNS.Upstream.Servers {
		upstreams = append(upstreams, *upstream.New(host))
	}

	return &Resolver{
		Upstreams: upstreams,
		Strategy:  cfg.DNS.Upstream.Strategy,
		Cache:     cache.New(),
		BlockList: blocklist.New(cfg.DNS.BlockList),
	}
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {

	log.Debug().
		Str("msg", "Processing request").
		Str("domain", req.Question[0].Name).
		Str("type", dns.TypeToString[req.Question[0].Qtype]).
		Send()

	// Check blocklist
	if block := r.BlockList.Query(req.Question[0].Name); block != nil {
		log.Debug().
			Str("msg", "Domain is blocked").
			Str("domain", req.Question[0].Name).
			Str("type", dns.TypeToString[req.Question[0].Qtype]).
			Str("category", block.Category).
			Str("reason", block.Reason).
			Send()
		return r.requestToResponse(req, []dns.RR{}, false), nil
	} else {
		log.Debug().
			Str("msg", "Domain is not blocked").
			Str("domain", req.Question[0].Name).
			Str("type", dns.TypeToString[req.Question[0].Qtype]).
			Send()
	}

	// Try cache
	if records, ok := r.Cache.Query(req.Question[0]); ok {
		return r.requestToResponse(req, records, true), nil
	} else {
		log.Debug().
			Str("msg", "No records found in cache").
			Str("domain", req.Question[0].Name).
			Str("type", dns.TypeToString[req.Question[0].Qtype]).
			Send()
	}

	// Try upstream
	upstream := r.selectUpstream()
	if msg, err := upstream.Query(req); err == nil {
		r.Cache.Add(msg.Answer)
		return r.requestToResponse(req, msg.Answer, false), nil
	}

	// Return NXDOMAIN
	log.Debug().
		Str("msg", "No records found").
		Str("domain", req.Question[0].Name).
		Str("type", dns.TypeToString[req.Question[0].Qtype]).
		Send()

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
