package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/resolver"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	startListeners(config)

	select {}
}

func startListeners(config *config.Config) {

	// Create shared local cache
	var localCache *cache.LocalCache
	if config.DNS.LocalDNSConfig.Enabled {
		localCache = cache.NewLocalCache(config.DNS.LocalDNSConfig.Records)
	}

	// Create shared upstream cache
	var upstreamCache *cache.UpstreamCache
	if config.DNS.Upstream.Enabled {
		upstreamCache = cache.NewUpstreamCache()
	}

	// Create shared resolver
	resolver := resolver.New(
		config.DNS.Upstream.Servers,
		config.DNS.Upstream.Strategy,
		localCache,
		upstreamCache,
	)

	if config.DNS.Protocols.UDP.Enabled {
		go listener.CreateUDPListener(config, resolver)
	}
	if config.DNS.Protocols.TCP.Enabled {
		go listener.CreateTCPListener(config, resolver)
	}
	//if cfg.DNS.Protocols.DOT.Enabled {
	//	go listener.CreateDOTListener(cfg.DNS.Protocols.DOT)
	//}
}
