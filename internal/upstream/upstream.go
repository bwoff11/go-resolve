package upstream

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// Upstream represents an upstream DNS server.
type Upstream struct {
	IP         net.IP
	totalRTT   time.Duration
	queryCount int
	rttLock    sync.Mutex // Protects access to totalRTT and queryCount
}

// New creates a new Upstream object with the specified IP address.
func New(host string) *Upstream {

	log.Debug().
		Str("msg", "Creating new upstream").
		Str("host", host).
		Send()

	// Parse IP address
	ip := net.ParseIP(host)
	if ip == nil {
		log.Fatal().
			Str("msg", "Failed to parse IP address").
			Str("host", host).
			Send()
	}

	// Return new Upstream
	return &Upstream{
		IP: ip,
	}
}

// Query sends the given DNS query message to the upstream DNS server and returns the response.
func (u *Upstream) Query(msg *dns.Msg) (*dns.Msg, error) {

	log.Debug().
		Str("msg", "Sending request to upstream").
		Str("domain", msg.Question[0].Name).
		Str("type", dns.TypeToString[msg.Question[0].Qtype]).
		Str("upstream", u.IP.String()).
		Send()

	c := new(dns.Client)
	address := fmt.Sprintf("%s:53", u.IP.String()) // Ensure IP is in string format

	resp, rtt, err := c.Exchange(msg, address)
	if err != nil {
		return nil, err
	}

	if len(resp.Answer) > 0 {
		log.Debug().
			Str("msg", "Received response from upstream").
			Str("domain", resp.Question[0].Name).
			Str("type", dns.TypeToString[resp.Question[0].Qtype]).
			Str("value", resp.Answer[0].String()).
			Int("ttl", int(resp.Answer[0].Header().Ttl)).
			Str("upstream", u.IP.String()).
			Str("rtt", rtt.String()).
			Send()
	} else {
		log.Info().
			Str("msg", "Received no response from upstream").
			Str("domain", resp.Question[0].Name).
			Str("type", dns.TypeToString[resp.Question[0].Qtype]).
			Str("upstream", u.IP.String()).
			Str("rtt", rtt.String()).
			Send()
	}

	return resp, nil
}
