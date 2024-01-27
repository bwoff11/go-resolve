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
func New(ip net.IP) *Upstream {
	return &Upstream{
		IP: ip,
	}
}

// Query sends the given DNS query message to the upstream DNS server and returns the response.
func (u *Upstream) Query(msg *dns.Msg) (*dns.Msg, error) {
	c := new(dns.Client)
	address := fmt.Sprintf("%s:53", u.IP.String()) // Ensure IP is in string format

	resp, rtt, err := c.Exchange(msg, address)
	if err != nil {
		return nil, err
	}

	// Update RTT statistics
	u.updateRTT(rtt)

	return resp, nil
}

// updateRTT updates the total RTT and query count.
func (u *Upstream) updateRTT(rtt time.Duration) {
	u.rttLock.Lock()
	defer u.rttLock.Unlock()
	u.totalRTT += rtt
	u.queryCount++
}

// GetAverageRTT returns the average RTT for all queries made.
func (u *Upstream) GetAverageRTT() time.Duration {
	u.rttLock.Lock()
	defer u.rttLock.Unlock()
	if u.queryCount == 0 {
		return 0
	}
	return u.totalRTT / time.Duration(u.queryCount)
}
