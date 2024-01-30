package upstream

import (
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/miekg/dns"
)

type Upstream struct {
	Servers  []*UpstreamServer `yaml:"upstreams"`
	Strategy config.Strategy   `yaml:"strategy"`
}

func New(cfg config.Upstream) *Upstream {

	// Create upstream servers
	var servers []*UpstreamServer
	for _, server := range cfg.Servers {
		servers = append(servers, NewUpstreamServer(server.IP, server.Port, server.Timeout))
	}

	// Return new upstream
	return &Upstream{
		Servers:  servers,
		Strategy: cfg.Strategy,
	}
}

func (u *Upstream) Query(msg *dns.Msg) []dns.RR {
	server := u.selectServer()
	return server.Query(msg)
}

func (u *Upstream) selectServer() *UpstreamServer {
	switch u.Strategy {
	case config.StrategyRandom:
		return u.randomServer()
	case config.StrategyRoundRobin:
		return u.roundRobinServer()
	case config.StrategyLatency:
		return u.latencyServer()
	case config.StrategySequential:
		return u.sequentialServer()
	default:
		return u.randomServer()
	}
}

func (u *Upstream) randomServer() *UpstreamServer {
	// Unimplemented
	return u.Servers[0]
}

func (u *Upstream) roundRobinServer() *UpstreamServer {
	// Unimplemented
	return u.Servers[0]
}

func (u *Upstream) latencyServer() *UpstreamServer {
	// Unimplemented
	return u.Servers[0]
}

func (u *Upstream) sequentialServer() *UpstreamServer {
	// Unimplemented
	return u.Servers[0]
}
