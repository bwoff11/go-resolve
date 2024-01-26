package resolver

import (
	"net"

	"github.com/miekg/dns"
)

type Resolver struct {
	// a cache or specific configuration
}

func New() *Resolver {
	return &Resolver{}
}

func (r *Resolver) Resolve(req dns.Msg) (dns.Msg, error) {
	switch req.Question[0].Qtype {
	case dns.TypeA:
		return r.resolveA(req)
	case dns.TypeAAAA:
		return r.resolveAAAA(req)
	case dns.TypeCNAME:
		return r.resolveCNAME(req)
	case dns.TypeMX:
		return r.resolveMX(req)
	case dns.TypeNS:
		return r.resolveNS(req)
	case dns.TypeTXT:
		return r.resolveTXT(req)
	default:
		return dns.Msg{}, nil
	}
}

func (r *Resolver) resolveA(req dns.Msg) (dns.Msg, error) {
	var resp dns.Msg
	resp.SetReply(&req)
	resp.Authoritative = true
	resp.Answer = []dns.RR{
		&dns.A{
			Hdr: dns.RR_Header{
				Name:   req.Question[0].Name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    60,
			},
			A: net.IPv4(1, 0, 0, 1),
		},
	}
	return resp, nil
}

func (r *Resolver) resolveAAAA(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveCNAME(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveMX(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveNS(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}

func (r *Resolver) resolveTXT(req dns.Msg) (dns.Msg, error) {
	return dns.Msg{}, nil
}
