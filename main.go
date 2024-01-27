package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
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

	cacheConfig := config.DNS.Cache
	upstreamConfig := config.DNS.Upstream

	if config.DNS.Protocols.UDP.Enabled {
		go listener.CreateUDPListener(config.DNS.Protocols.UDP, cacheConfig, upstreamConfig)
	}
	if config.DNS.Protocols.TCP.Enabled {
		go listener.CreateTCPListener(config.DNS.Protocols.TCP, cacheConfig, upstreamConfig)
	}
	//if cfg.DNS.Protocols.DOT.Enabled {
	//	go listener.CreateDOTListener(cfg.DNS.Protocols.DOT)
	//}
}
