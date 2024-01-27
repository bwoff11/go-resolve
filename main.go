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

	if config.DNS.Protocols.UDP.Enabled {
		go listener.CreateUDPListener(config)
	}
	if config.DNS.Protocols.TCP.Enabled {
		go listener.CreateTCPListener(config)
	}
	//if cfg.DNS.Protocols.DOT.Enabled {
	//	go listener.CreateDOTListener(cfg.DNS.Protocols.DOT)
	//}
}
