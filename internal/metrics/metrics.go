package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalQueries = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_total_queries",
			Help: "Total number of DNS queries received.",
		},
	)

	CacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_hits",
			Help: "Total number of DNS queries answered from cache.",
		},
	)

	CacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_cache_misses",
			Help: "Total number of DNS queries that resulted in a cache miss.",
		},
	)
)

func init() {
	// Register custom metrics with Prometheus
	prometheus.MustRegister(TotalQueries, CacheHits, CacheMisses)
}
