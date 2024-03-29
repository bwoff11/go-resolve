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

	Latency time.Duration
}

// New creates a new Upstream object with the specified IP address.
func NewUpstreamServer(host string, port int, timeout int) *UpstreamServer {

	// Parse IP address
	ip := net.ParseIP(host)
	if ip == nil {
		log.Fatal().
			Str("msg", "failed to parse IP address").
			Str("host", host).
			Send()
	}

	// Return new Upstream
	return &UpstreamServer{
		IP:      ip,
		Address: fmt.Sprintf("%s:53", ip.String()),
		Timeout: timeout,
	}
}

// Query sends the given DNS query message to the upstream DNS server and returns the response.
func (us *UpstreamServer) Query(msg *dns.Msg) (response []dns.RR) {
	startTime := time.Now()
	defer func() { metrics.UpstreamDuration.Observe(time.Since(startTime).Seconds()) }()

	client := &dns.Client{
		Net:     "udp",
		Timeout: time.Duration(us.Timeout) * time.Second,
	}

	resp, rtt, err := client.Exchange(msg, us.Address)
	metrics.UpstreamRTT.WithLabelValues(us.Address).Observe(rtt.Seconds())
	us.Latency = rtt

	if err != nil {
		log.Error().Str("msg", "Failed to query upstream DNS server").Str("address", us.Address).Err(err).Send()
		return nil
	}

	if resp == nil {
		return nil
	}

	if resp.Rcode != dns.RcodeSuccess {
		log.Debug().Str("msg", "Upstream DNS server responded with error").Str("address", us.Address).Int("rcode", resp.Rcode).Send()
		return nil
	}

	log.Debug().Str("msg", "Upstream DNS server responded").Str("address", us.Address).Send()
	return resp.Answer
}
