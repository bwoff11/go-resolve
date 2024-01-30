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

	BlocklistDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_blocklist_duration",
			Help:    "Histogram of blocklist query durations.",
			Buckets: []float64{0.0000001, 0.000001, 0.00001, 0.0001, 0.001},
		},
	)

	BlockedCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dns_blocked_count",
			Help: "Total number of DNS queries blocked by the blocklist.",
		},
	)

	CacheDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_cache_duration",
			Help:    "Histogram of cache query durations.",
			Buckets: []float64{0.0000001, 0.000001, 0.00001, 0.0001, 0.001},
		},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_request_duration",
			Help:    "Duration of request starting from networking layer to response.",
			Buckets: []float64{0.01, 0.1, 0.25},
		},
		[]string{"protocol"},
	)

	ResolutionDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "dns_resolution_duration",
			Help:    "Time taken to process a request after handoff from networking layer.",
			Buckets: []float64{0.01, 0.1, 0.25},
		},
	)

	UpstreamDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dns_upstream_duration",
			Help:    "Time taken to query upstream DNS server.",
			Buckets: []float64{0.01, 0.02, 0.03, 0.04, 0.05},
		},
		[]string{"server"},
	)
)

func init() {
	// Register custom metrics with Prometheus
	prometheus.MustRegister(TotalQueries, CacheHits, BlocklistDuration, UpstreamDuration, CacheDuration, ResolutionDuration, RequestDuration)
}
