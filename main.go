package main

import (
	"net/http"
	"strconv"

	"github.com/bwoff11/go-resolve/internal/config"
	"github.com/bwoff11/go-resolve/internal/transport"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	startMetricsServer(cfg.Metrics)
	startTransports(cfg.Transport)

	select {}
}

func startTransports(cfg config.Transport) {
	ts, err := transport.New(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start transports")
	}

	if err := ts.Listen(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start transports")
	}
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
