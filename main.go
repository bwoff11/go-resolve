package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/patrickmn/go-cache"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create a shared resolver with a global cache
	var globalCache *cache.Cache
	if cfg.DNS.Cache.Enabled {
		globalCache = cache.New(cfg.DNS.Cache.TTL, cfg.DNS.Cache.PurgeInterval)
	}
	sharedResolver := resolver.New(
		cfg.DNS.Upstream.Servers,
		globalCache,
	)

	// Start listeners with the shared resolver
	startListener(cfg.DNS.Protocols.UDP, "udp", sharedResolver)
	startListener(cfg.DNS.Protocols.TCP, "tcp", sharedResolver)
	startListener(cfg.DNS.Protocols.DOT, "dot", sharedResolver)
	// startWebServer(cfg.Web)
	// startMetricsEndpoint(cfg.Metrics)

	select {} // Keep the main goroutine running
}

func startListener(protocolConfig config.ProtocolConfig, protocolType config.ProtocolType, sharedResolver *resolver.Resolver) {
	if protocolConfig.Enabled {
		l := listener.New(protocolType, protocolConfig.Port, sharedResolver)
		go func() {
			if err := l.Listen(); err != nil {
				log.Fatalf("%s Listener failed: %v", protocolType, err)
			}
		}()
	}
}
