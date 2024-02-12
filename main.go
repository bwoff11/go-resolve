package main

import (
	"net/http"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/resolver"
	"github.com/bwoff11/go-resolve/internal/transport"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {

	// No error handling on these functions because
	// if they fail, the program should exit.

	cfg, _ := config.Load()

	startMetricsServer(&cfg.Metrics)

	transports := transport.New(&cfg.Transport)
	resolver := resolver.New(cfg, transports.Queue)
	resolver.Start()

	select {} // replace with a signal handler
}

func startMetricsServer(cfg *config.Metrics) {
	http.Handle(cfg.Route, promhttp.Handler())
	go func() {
		log.Info().Str("address", ":"+strconv.Itoa(cfg.Port)).Msg("Starting Prometheus metrics server")
		if err := http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start Prometheus metrics server")
		}
	}()
}
