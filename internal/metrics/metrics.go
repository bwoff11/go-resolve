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

	BlocklistDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_blocklist_duration",
			Help:    "Histogram of blocklist query durations.",
			Buckets: []float64{0.0000001, 0.000001, 0.00001, 0.0001, 0.001},
		},
	)

	CacheDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_cache_duration",
			Help:    "Histogram of cache query durations.",
			Buckets: []float64{0.0000001, 0.000001, 0.00001, 0.0001, 0.001},
		},
	)

	UpstreamDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_upstream_duration",
			Help:    "Histogram of upstream query durations.",
			Buckets: []float64{0.01, 0.02, 0.03, 0.04, 0.05},
		},
	)
)

func init() {
	// Register custom metrics with Prometheus
	prometheus.MustRegister(TotalQueries, CacheHits, CacheMisses, BlocklistDuration, UpstreamDuration, CacheDuration)
}
