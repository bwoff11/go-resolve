package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/cache"
	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/resolver"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load the cache
	cache, err := cache.New(cfg.DNS.Cache, cfg.DNS.Local)
	if err != nil {
		log.Fatalf("Failed to load cache: %v", err)
	}

	startListeners(cfg, cache)

	select {}
}

func startListeners(config *config.Config, cache *cache.Cache) {

	// Create shared resolver
	resolver := resolver.New(
		config.DNS.Upstream.Servers,
		config.DNS.Upstream.Strategy,
		cache,
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
