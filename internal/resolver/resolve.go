package resolver

import (
	"net"

	"github.com/bwoff11/go-resolve/internal/blocklist"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// Resolve processes the DNS query and returns a response.
func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
	logRequest(req)
	if block := r.checkBlockList(req); block != nil {
		return r.blockedResponse(req), nil
	}

	if records, authoritative := r.queryCache(req); authoritative {
		return r.createResponse(req, records, true), nil
	}

	return r.queryUpstream(req)
}

func logRequest(req *dns.Msg) {
	log.Debug().
		Str("msg", "Processing request").
		Str("domain", req.Question[0].Name).
		Str("type", dns.TypeToString[req.Question[0].Qtype]).
		Send()
}

func (r *Resolver) checkBlockList(req *dns.Msg) *blocklist.Block {
	block := r.BlockList.Query(req.Question[0].Name)
	if block != nil {
		log.Debug().
			Str("msg", "Blocked domain requested").
			Str("domain", block.Domain).
			Str("category", block.Category).
			Str("reason", block.Reason).
			Send()
	}
	return block
}

func (r *Resolver) queryCache(req *dns.Msg) ([]dns.RR, bool) {
	records, ok := r.Cache.Query(req.Question[0])
	if ok {
		log.Debug().
			Str("msg", "Found record in cache").
			Str("domain", records[0].Header().Name).
			Str("type", dns.TypeToString[records[0].Header().Rrtype]).
			//Str("value", getRRValue(records[0])).
			Int("ttl", int(records[0].Header().Ttl)).
			Send()
	}
	return records, ok
}

func (r *Resolver) queryUpstream(req *dns.Msg) (*dns.Msg, error) {
	upstream := r.selectUpstream()
	return upstream.Query(req)
}

// selectUpstream selects an upstream server based on the configured strategy.
func (r *Resolver) selectUpstream() *upstream.Upstream {
	// Implement load balancing logic here.
	return &r.Upstreams[0]
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
