package main

import (
	"net/http"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/listener"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	startMetricsServer(cfg.Metrics)
	startListeners(cfg)

	select {}
}

func startListeners(config *config.Config) {

	// Create shared resolver
	resolver := resolver.New(config)

	if config.Protocols.UDP.Enabled {
		go listener.CreateUDPListener(config, resolver)
	}
	if config.Protocols.TCP.Enabled {
		go listener.CreateTCPListener(config, resolver)
	}
	//if cfg.DNS.Protocols.DOT.Enabled {
	//	go listener.CreateDOTListener(cfg.DNS.Protocols.DOT)
	//}
}

func startMetricsServer(cfg config.Metrics) {
	http.Handle(cfg.Route, promhttp.Handler())
	go func() {
		log.Info().Str("address", ":"+strconv.Itoa(cfg.Port)).Msg("Starting Prometheus metrics server")
		if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start Prometheus metrics server")
		}
	}()
}
