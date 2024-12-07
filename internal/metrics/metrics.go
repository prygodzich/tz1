package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	FileRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_requests_total",
			Help: "The total number of file requests",
		},
		[]string{"status"},
	)

	FileRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "file_request_duration_seconds",
			Help:    "The duration of file requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"status"},
	)

	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "The total number of cache hits",
		},
		[]string{"cache_type"},
	)

	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "The total number of cache misses",
		},
		[]string{"cache_type"},
	)
)
