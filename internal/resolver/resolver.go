package resolver

import (
	"fmt"
	"net"

	"github.com/bwoff11/go-resolve/internal/upstream"
	"github.com/miekg/dns"
)

type Resolver struct {
	Upstreams []upstream.Upstream
}

func New(hosts []string) *Resolver {

	// Parse upstreams as IP addresses
	var ips []net.IP
	for _, host := range hosts {
		ip := net.ParseIP(host)
		if ip == nil {
			panic(fmt.Sprintf("Invalid IP address: %s", host))
		}
		ips = append(ips, ip)
	}

	// Convert ip addresses to upstreams
	var upstreams []upstream.Upstream
	for _, ip := range ips {
		upstreams = append(upstreams, *upstream.New(ip))
	}

	return &Resolver{
		Upstreams: upstreams,
	}
}

func (r *Resolver) Resolve(req *dns.Msg) (*dns.Msg, error) {
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
