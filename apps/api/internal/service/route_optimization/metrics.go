package route_optimization

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	routeOptimizationCacheEventsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "route_optimization_cache_events_total",
		Help: "Route optimization cache events (hit/miss/ineligible).",
	}, []string{"cache", "result"})

	routeOptimizationOSRMSkippedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "route_optimization_osrm_skipped_total",
		Help: "Total number of times OSRM Route call was skipped due to cache hit.",
	})

	routeOptimizationOSRMRouteDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "route_optimization_osrm_route_duration_seconds",
		Help:    "Duration of OSRM Route call in seconds.",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 20, 30},
	})

	routeOptimizationMatrixDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "route_optimization_matrix_duration_seconds",
		Help:    "Duration of distance matrix calculation in seconds.",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 20, 30},
	}, []string{"method"})

	routeOptimizationOptimizeDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "route_optimization_optimize_duration_seconds",
		Help:    "End-to-end duration of Optimize() service call in seconds.",
		Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 20, 30},
	})
)
