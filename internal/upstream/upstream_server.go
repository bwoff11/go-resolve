package upstream

import (
	"fmt"
	"net"
	"time"

	"github.com/bwoff11/go-resolve/internal/metrics"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
)

// Upstream represents an upstream DNS server.
type UpstreamServer struct {
	IP      net.IP
	Address string // IP:Port
	Timeout int
	Client  *dns.Client
}

// New creates a new Upstream object with the specified IP address.
func NewUpstreamServer(host string, port int, timeout int) *UpstreamServer {

	// Parse IP address
	ip := net.ParseIP(host)
	if ip == nil {
		log.Fatal().
			Str("msg", "Failed to parse IP address").
			Str("host", host).
			Send()
	}

	// Return new Upstream
	return &UpstreamServer{
		IP:      ip,
		Address: fmt.Sprintf("%s:53", ip.String()),
		Timeout: timeout,
		Client:  new(dns.Client),
	}
}

// Query sends the given DNS query message to the upstream DNS server and returns the response.
func (us *UpstreamServer) Query(msg *dns.Msg) []dns.RR {
	startTime := time.Now()
	resp, _, err := us.Client.Exchange(msg, us.Address)
	if err != nil {
		log.Error().
			Str("msg", "Failed to query upstream DNS server").
			Str("address", us.Address).
			Err(err).
			Send()
		return nil
	}

	metrics.UpstreamDuration.Observe(time.Since(startTime).Seconds())
	return resp.Answer
}
