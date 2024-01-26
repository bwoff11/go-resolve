package resolver

import (
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstream []string
}

func New() *Resolver {
	return &Resolver{
		Upstream: []string{""},
	}
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

	c := new(dns.Client)

	// Set up a message to query the external DNS server
	m := new(dns.Msg)
	m.SetQuestion(req.Question[0].Name, dns.TypeA)

	// Perform the query
	resp, _, err := c.Exchange(m, "8.8.8.8:53")
	if err != nil {
		// Handle error
		return dns.Msg{}, err
	}

	// Construct the response to the original query
	var response dns.Msg
	response.SetReply(&req)
	response.Authoritative = false // Set false since it's not authoritative
	response.Answer = resp.Answer

	return response, nil
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

func (r *Resolver) createResponse(req dns.Msg, record dns.RR) dns.Msg {
	resp := dns.Msg{}
	resp.SetReply(&req)
	resp.Answer = append(resp.Answer, record)
	return resp
}
