package upstream

import (
	"math/rand"
	"sync"
	"time"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

type Upstream struct {
	Servers          []*UpstreamServer
	Strategy         config.Strategy
	selectServerFunc func() *UpstreamServer // Function pointer for server selection
	counter          int
	mutex            sync.Mutex
}

// NewUpstream creates a new Upstream instance based on the given config.
func New(cfg config.Upstream) *Upstream {
	servers := make([]*UpstreamServer, 0, len(cfg.Servers))
	for _, server := range cfg.Servers {
		servers = append(servers, NewUpstreamServer(server.IP, server.Port, server.Timeout))
	}

	upstream := &Upstream{
		Servers:  servers,
		Strategy: cfg.Strategy,
	}

	// Assign the server selection function based on strategy
	switch cfg.Strategy {
	case config.StrategyRandom:
		upstream.selectServerFunc = upstream.randomServer
	case config.StrategyRoundRobin:
		upstream.selectServerFunc = upstream.roundRobinServer
	case config.StrategyLatency:
		upstream.selectServerFunc = upstream.latencyServer
	case config.StrategySequential:
		upstream.selectServerFunc = upstream.sequentialServer
	default:
		upstream.selectServerFunc = upstream.randomServer
	}

	return upstream
}

// Query forwards the DNS query to an appropriate upstream server.
func (u *Upstream) Query(msg *dns.Msg) []dns.RR {
	server := u.selectServerFunc() // Use the assigned function
	return server.Query(msg)
}

// randomServer selects a random server from the list.
func (u *Upstream) randomServer() *UpstreamServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(u.Servers))
	return u.Servers[index]
}

// roundRobinServer selects servers in a round-robin fashion.
func (u *Upstream) roundRobinServer() *UpstreamServer {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	server := u.Servers[u.counter%len(u.Servers)]
	u.counter++
	return server
}

// latencyServer selects the server based on latency.
// Placeholder implementation, should be replaced with actual latency measurement logic.
func (u *Upstream) latencyServer() *UpstreamServer {
	var minLatency time.Duration
	var selected *UpstreamServer
	for _, server := range u.Servers {
		if selected == nil || server.Latency < minLatency {
			minLatency = server.Latency
			selected = server
		}
	}
	return selected
}

// sequentialServer selects servers sequentially.
// Placeholder implementation for demonstration.
func (u *Upstream) sequentialServer() *UpstreamServer {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	server := u.Servers[u.counter%len(u.Servers)]
	u.counter++
	return server
}
