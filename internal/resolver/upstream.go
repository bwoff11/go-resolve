package resolver

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func (r *Resolver) queryUpstream(req dns.Msg) (dns.Msg, error) {
	client := new(dns.Client)
	msg := new(dns.Msg)
	msg.SetQuestion(req.Question[0].Name, req.Question[0].Qtype)

	upstream := r.chooseUpstream()
	if upstream == "" {
		return dns.Msg{}, fmt.Errorf("no upstream DNS servers configured")
	}

	// Ensure the upstream address includes a port number
	if !strings.Contains(upstream, ":") {
		upstream = fmt.Sprintf("%s:53", upstream) // Default DNS port is 53
	}

	resp, _, err := client.Exchange(msg, upstream)
	if err != nil {
		return dns.Msg{}, err
	}

	response := r.createResponse(req, resp.Answer)
	return response, nil
}

func (r *Resolver) chooseUpstream() string {
	// Implement load balancing or failover logic if needed
	if len(r.Upstreams) > 0 {
		return r.Upstreams[0] // Simple selection for now
	}
	return ""
}
