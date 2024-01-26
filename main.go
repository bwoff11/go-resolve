package main

import (
	"log"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/resolver"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	res := resolver.New()

	startListeners(res, cfg)
	// startDNSOverTLS(cfg)
	// startWebServer(cfg)
	// startMetricsEndpoint(cfg)
	// ...

	select {} // Keep the main goroutine running
}

func startListeners(res *resolver.Resolver, cfg *config.Config) {
	if cfg.DNS.UDP.Enabled {
		startUDPListener(res, cfg.DNS.UDP.Port)
	}

	if cfg.DNS.TCP.Enabled {
		startTCPListener(res, cfg.DNS.TCP.Port)
	}
}

func startUDPListener(res *resolver.Resolver, port int) {
	udpListener := listener.New(res, "udp", port)
	go func() {
		if err := udpListener.Listen(); err != nil {
			log.Fatalf("UDP Listener failed: %v", err)
		}
	}()
}

func startTCPListener(res *resolver.Resolver, port int) {
	tcpListener := listener.New(res, "tcp", port)
	go func() {
		if err := tcpListener.Listen(); err != nil {
			log.Fatalf("TCP Listener failed: %v", err)
		}
	}()
}

// Define other start functions for DNSOverTLS, WebServer, MetricsEndpoint, etc.
