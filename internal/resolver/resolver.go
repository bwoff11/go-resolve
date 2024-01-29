package resolver

import (
	"net"
	"time"

	"github.com/bwoff11/go-resolve/internal/blocklist"
	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstream  *upstream.Upstream
	Strategy  config.LoadBalancingStrategy
	Cache     *cache.Cache
	BlockList *blocklist.BlockList
}

// New creates a new Resolver instance.
func New(cfg *config.Config) *Resolver {
	return &Resolver{
		Upstream:  upstream.New(cfg.DNS.Upstream),
		Cache:     cache.New(cfg.DNS.Local),
		BlockList: blocklist.New(cfg.DNS.BlockList),
	}
}

// Resolve processes the DNS query and returns a response.
func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
	startTime := time.Now()

	qName := req.Question[0].Name

	// Check block list
	if block := r.BlockList.Query(qName); block != nil {
		return r.blockedResponse(req), nil
	}

	// Check cache
	if records := r.Cache.Query(req.Question); len(records) > 0 {
		return r.createResponse(req, records, true), nil
	}

	// Check upstream
	if records := r.Upstream.Query(req); len(records) > 0 {
		r.Cache.Add(records)
		return r.createResponse(req, records, false), nil
	}

	metrics.ResolutionDuration.Observe(time.Since(startTime).Seconds())
	return r.createResponse(req, []dns.RR{}, false), nil // Need to verify this is correct for NXDOMAIN
}

// createResponse builds a DNS response message.
func (r *Resolver) createResponse(req *dns.Msg, answer []dns.RR, authoritative bool) *dns.Msg {
	return &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Id:                 req.Id,
			Response:           true,
			Opcode:             req.Opcode,
			Authoritative:      authoritative,
			Truncated:          false,
			RecursionDesired:   req.RecursionDesired,
			RecursionAvailable: true,
			Rcode:              dns.RcodeSuccess,
		},
		Compress: false,
		Question: req.Question,
		Answer:   answer,
		Ns:       []dns.RR{}, // Implement if needed
		Extra:    []dns.RR{}, // Implement if needed
	}
}

func (r *Resolver) blockedResponse(req *dns.Msg) *dns.Msg {
	var answer []dns.RR
	switch req.Question[0].Qtype {
	case dns.TypeA:
		answer = []dns.RR{&dns.A{
			Hdr: dns.RR_Header{
				Name:   req.Question[0].Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			A: net.IPv4zero,
		}}
	case dns.TypeAAAA:
		answer = []dns.RR{&dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   req.Question[0].Name,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			AAAA: net.IPv6zero,
		}}
	case dns.TypeCNAME:
		answer = []dns.RR{&dns.CNAME{
			Hdr: dns.RR_Header{
				Name:   req.Question[0].Name,
				Rrtype: dns.TypeCNAME,
				Class:  dns.ClassINET,
				Ttl:    0,
			},
			Target: "blocked.local",
		}}
	default:
		answer = []dns.RR{}
	}

	return r.createResponse(req, answer, true)
}
