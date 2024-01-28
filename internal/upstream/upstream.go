package upstream

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
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

	// Parse IP address
	ip := net.ParseIP(host)
	if ip == nil {
		panic(fmt.Sprintf("Invalid IP address: %s", host))
	}

	// Return new Upstream
	return &Upstream{
		IP: ip,
	}
}

// Query sends the given DNS query message to the upstream DNS server and returns the response.
func (u *Upstream) Query(msg *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	address := fmt.Sprintf("%s:53", u.IP.String()) // Ensure IP is in string format

	resp, _, err := c.Exchange(msg, address)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
