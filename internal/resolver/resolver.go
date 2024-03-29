package resolver

import (
	"net"
	"time"

	"github.com/bwoff11/go-resolve/internal/blocklist"
	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/local"
	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/bwoff11/go-resolve/internal/transport"
	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

type Resolver struct {
	BlockList *blocklist.BlockList
	Cache     *cache.Cache
	Local     *local.LocalRecords
	Upstream  *upstream.Upstream
	Queue     chan transport.QueueItem
}

// New creates a new Resolver instance.
func New(cfg *config.Config, q chan transport.QueueItem) *Resolver {
	return &Resolver{
		Upstream:  upstream.New(cfg.Upstream),
		Local:     local.New(&cfg.Local),
		Cache:     cache.New(),
		BlockList: blocklist.New(cfg.BlockLists),
		Queue:     q,
	}
}

func (r *Resolver) Start() {
	go func() {
		for item := range r.Queue {
			req := item.Message()
			resp, err := r.Resolve(req)
			if err != nil {
				log.Error().Err(err).Msg("Failed to resolve query")
				continue
			}
			item.Respond(resp)
		}
	}()
	log.Info().Msg("Resolver started and listening on the inbound queue")
}

// Resolve processes the DNS query and returns a response.
func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
	log.Debug().Str("domain", req.Question[0].Name).Msg("resolving domain")
	startTime := time.Now()

	q := &req.Question[0] // Only support one question
	qName := req.Question[0].Name

	// Check block list
	if block := r.BlockList.Query(qName); block != nil {
		return r.blockedResponse(req, startTime), nil
	}

	// Check local records
	if records := r.Local.Query(q); len(records) > 0 {
		return r.createResponse(req, records, true, startTime), nil
	}

	// Check cache
	if records := r.Cache.Query(q); len(records) > 0 {
		return r.createResponse(req, records, true, startTime), nil
	}

	// Check upstream
	if records := r.Upstream.Query(req); len(records) > 0 {
		r.Cache.Add(q, records)
		return r.createResponse(req, records, false, startTime), nil
	}

	log.Info().Str("domain", qName).Msg("domain not found in local, cache, or upstream")
	return r.createResponse(req, []dns.RR{}, false, startTime), nil // Need to verify this is correct for NXDOMAIN
}

// createResponse builds a DNS response message.
func (r *Resolver) createResponse(req *dns.Msg, answer []dns.RR, authoritative bool, startTime time.Time) *dns.Msg {
	msg := &dns.Msg{
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
	log.Debug().
		Str("domain", req.Question[0].Name).
		Str("type", dns.TypeToString[req.Question[0].Qtype]).
		//Str("answer", answer[0].String()). //possibly nil
		Msg("created response")
	metrics.ResolutionDuration.Observe(time.Since(startTime).Seconds())
	return msg
}

func (r *Resolver) blockedResponse(req *dns.Msg, startTime time.Time) *dns.Msg {
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

	return r.createResponse(req, answer, true, startTime)
}
