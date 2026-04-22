package handlers

import (
	"net/http"
	"net/http/pprof"
	"runtime"
	"time"

	"github.com/gilabs/crm-healthcare/api/pkg/cache"
	"github.com/gilabs/crm-healthcare/api/pkg/circuitbreaker"
	"github.com/gilabs/crm-healthcare/api/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler handles system and cache health endpoints.
type HealthHandler struct{}

// NewHealthHandler creates a new health handler.
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CacheHealthResp represents the cache health response payload.
type CacheHealthResp struct {
	Status         string           `json:"status"`
	Mode           string           `json:"mode"`
	CircuitBreaker string           `json:"circuit_breaker"`
	Latency        string           `json:"latency"`
	HitRate        float64          `json:"hit_rate_percent"`
	AvgLatency     float64          `json:"avg_latency_us"`
	Uptime         string           `json:"uptime"`
	Metrics        *CacheMetricsDTO `json:"metrics"`
	Pool           *PoolStatsDTO    `json:"pool,omitempty"`
}

// CacheMetricsDTO represents cache metrics data.
type CacheMetricsDTO struct {
	Hits              uint64              `json:"hits"`
	Misses            uint64              `json:"misses"`
	Errors            uint64              `json:"errors"`
	Timeouts          uint64              `json:"timeouts"`
	CBOpen            uint64              `json:"circuit_breaker_open_count"`
	CBHalfOpen        uint64              `json:"circuit_breaker_half_open_count"`
	PipelineOps       uint64              `json:"pipeline_operations"`
	Compressions      uint64              `json:"compressions"`
	TotalOps          uint64              `json:"total_operations"`
	MaxLatency        uint64              `json:"max_latency_us"`
	EnterpriseMetrics *AdvancedMetricsDTO `json:"advanced_metrics,omitempty"`
}

// PoolStatsDTO represents connection pool statistics.
type PoolStatsDTO struct {
	Hits       uint64 `json:"hits"`
	Misses     uint64 `json:"misses"`
	Timeouts   uint64 `json:"timeouts"`
	TotalConns uint64 `json:"total_connections"`
	IdleConns  uint64 `json:"idle_connections"`
	StaleConns uint64 `json:"stale_connections"`
}

// SystemHealthResp represents overall system health response.
type SystemHealthResp struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]interface{} `json:"services"`
	System    *SystemInfoDTO         `json:"system"`
}

// SystemInfoDTO represents system runtime information.
type SystemInfoDTO struct {
	GoVersion    string `json:"go_version"`
	NumGoroutine int    `json:"goroutines"`
	NumCPU       int    `json:"cpus"`
	MemAlloc     uint64 `json:"memory_alloc_mb"`
	MemSys       uint64 `json:"memory_sys_mb"`
	GCCycles     uint32 `json:"gc_cycles"`
}

var startTime = time.Now()

// GetCacheHealth returns the health status of the Redis cache.
func (h *HealthHandler) GetCacheHealth(c *gin.Context) {
	if cache.Client == nil || !cache.Client.IsEnabled() {
		response.SuccessResponse(c, CacheHealthResp{
			Status:  "disabled",
			Mode:    "none",
			Uptime:  time.Since(startTime).String(),
			Metrics: &CacheMetricsDTO{},
		}, nil)
		return
	}

	healthStatus, err := cache.Client.HealthCheck()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, CacheHealthResp{
			Status:         "unhealthy",
			Mode:           healthStatus.Mode,
			CircuitBreaker: healthStatus.CBState,
			Uptime:         time.Since(startTime).String(),
			Metrics:        toMetricsDTO(&healthStatus.Metrics),
		})
		return
	}

	resp := CacheHealthResp{
		Status:         "healthy",
		Mode:           healthStatus.Mode,
		CircuitBreaker: healthStatus.CBState,
		Latency:        healthStatus.Latency.String(),
		HitRate:        healthStatus.HitRate,
		AvgLatency:     healthStatus.AvgLat,
		Uptime:         time.Since(startTime).String(),
		Metrics:        toMetricsDTO(&healthStatus.Metrics),
		Pool: &PoolStatsDTO{
			Hits:       healthStatus.Metrics.PoolHits,
			Misses:     healthStatus.Metrics.PoolMisses,
			Timeouts:   healthStatus.Metrics.PoolTimeouts,
			TotalConns: healthStatus.Metrics.TotalConns,
			IdleConns:  healthStatus.Metrics.IdleConns,
			StaleConns: healthStatus.Metrics.StaleConns,
		},
	}

	response.SuccessResponse(c, resp, nil)
}

// GetSystemHealth returns overall system health including all services.
func (h *HealthHandler) GetSystemHealth(c *gin.Context) {
	services := make(map[string]interface{})
	overallStatus := "healthy"

	if cache.Client != nil && cache.Client.IsEnabled() {
		healthStatus, err := cache.Client.HealthCheck()
		if err != nil {
			services["redis"] = map[string]interface{}{
				"status":          "unhealthy",
				"mode":            healthStatus.Mode,
				"circuit_breaker": healthStatus.CBState,
				"error":           err.Error(),
			}
			overallStatus = "degraded"
		} else {
			services["redis"] = map[string]interface{}{
				"status":          "healthy",
				"mode":            healthStatus.Mode,
				"circuit_breaker": healthStatus.CBState,
				"latency":         healthStatus.Latency.String(),
				"hit_rate":        healthStatus.HitRate,
			}
		}
	} else {
		services["redis"] = map[string]interface{}{
			"status": "disabled",
		}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	resp := SystemHealthResp{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Services:  services,
		System: &SystemInfoDTO{
			GoVersion:    runtime.Version(),
			NumGoroutine: runtime.NumGoroutine(),
			NumCPU:       runtime.NumCPU(),
			MemAlloc:     memStats.Alloc / 1024 / 1024,
			MemSys:       memStats.Sys / 1024 / 1024,
			GCCycles:     memStats.NumGC,
		},
	}

	response.SuccessResponse(c, resp, nil)
}

// GetCacheMetrics returns detailed cache metrics.
func (h *HealthHandler) GetCacheMetrics(c *gin.Context) {
	if cache.Client == nil || !cache.Client.IsEnabled() {
		response.SuccessResponse(c, CacheMetricsDTO{}, nil)
		return
	}

	metrics := cache.Client.GetMetrics()
	advancedMetrics := cache.Advanced().GetAdvancedMetrics()

	dto := toMetricsDTO(&metrics)
	dto.EnterpriseMetrics = &AdvancedMetricsDTO{
		WarmupTasks:    advancedMetrics.WarmupTasks,
		WarmupSuccess:  advancedMetrics.WarmupSuccess,
		WarmupFailures: advancedMetrics.WarmupFailures,
		BatchOps:       advancedMetrics.BatchOps,
		CacheAsideHit:  advancedMetrics.CacheAsideHit,
		CacheAsideMiss: advancedMetrics.CacheAsideMiss,
		Invalidations:  advancedMetrics.Invalidations,
		PatternDeletes: advancedMetrics.PatternDeletes,
	}

	response.SuccessResponse(c, dto, nil)
}

// ResetCacheMetrics resets cache metrics counters.
func (h *HealthHandler) ResetCacheMetrics(c *gin.Context) {
	response.SuccessResponse(c, map[string]string{
		"message": "Metrics are cumulative and cannot be reset",
	}, nil)
}

// AdvancedMetricsDTO represents advanced cache metrics
type AdvancedMetricsDTO struct {
	WarmupTasks    uint64 `json:"warmup_tasks"`
	WarmupSuccess  uint64 `json:"warmup_success"`
	WarmupFailures uint64 `json:"warmup_failures"`
	BatchOps       uint64 `json:"batch_ops"`
	CacheAsideHit  uint64 `json:"cache_aside_hit"`
	CacheAsideMiss uint64 `json:"cache_aside_miss"`
	Invalidations  uint64 `json:"invalidations"`
	PatternDeletes uint64 `json:"pattern_deletes"`
}

func toMetricsDTO(m *cache.Metrics) *CacheMetricsDTO {
	return &CacheMetricsDTO{
		Hits:         m.Hits,
		Misses:       m.Misses,
		Errors:       m.Errors,
		Timeouts:     m.Timeouts,
		CBOpen:       m.CBOpen,
		CBHalfOpen:   m.CBHalfOpen,
		PipelineOps:  m.PipelineOps,
		Compressions: m.Compressions,
		TotalOps:     m.OperationCount,
		MaxLatency:   m.MaxLatency,
	}
}

// GetCircuitBreakerStats returns circuit breaker statistics for all external services
func (h *HealthHandler) GetCircuitBreakerStats(c *gin.Context) {
	stats := circuitbreaker.Stats()

	response.SuccessResponse(c, map[string]interface{}{
		"timestamp":        time.Now(),
		"circuit_breakers": stats,
	}, nil)
}

// GetRuntimeStats returns Go runtime statistics for profiling
func (h *HealthHandler) GetRuntimeStats(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats := map[string]interface{}{
		"timestamp":  time.Now(),
		"go_version": runtime.Version(),
		"goroutines": runtime.NumGoroutine(),
		"cpus":       runtime.NumCPU(),
		"memory": map[string]interface{}{
			"alloc_mb":         memStats.Alloc / 1024 / 1024,
			"sys_mb":           memStats.Sys / 1024 / 1024,
			"heap_alloc_mb":    memStats.HeapAlloc / 1024 / 1024,
			"heap_sys_mb":      memStats.HeapSys / 1024 / 1024,
			"heap_inuse_mb":    memStats.HeapInuse / 1024 / 1024,
			"heap_idle_mb":     memStats.HeapIdle / 1024 / 1024,
			"heap_released_mb": memStats.HeapReleased / 1024 / 1024,
			"total_alloc_mb":   memStats.TotalAlloc / 1024 / 1024,
		},
		"gc": map[string]interface{}{
			"cycles":          memStats.NumGC,
			"last_gc_ns":      memStats.LastGC,
			"pause_total_ms":  memStats.PauseTotalNs / 1e6,
			"gc_cpu_fraction": memStats.GCCPUFraction,
		},
	}

	response.SuccessResponse(c, stats, nil)
}

// PprofGroup returns pprof handlers for profiling
func (h *HealthHandler) PprofGroup() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/debug/pprof/":
			pprof.Index(w, r)
		case "/debug/pprof/cmdline":
			pprof.Cmdline(w, r)
		case "/debug/pprof/profile":
			pprof.Profile(w, r)
		case "/debug/pprof/symbol":
			pprof.Symbol(w, r)
		case "/debug/pprof/trace":
			pprof.Trace(w, r)
		case "/debug/pprof/goroutine":
			pprof.Handler("goroutine").ServeHTTP(w, r)
		case "/debug/pprof/heap":
			pprof.Handler("heap").ServeHTTP(w, r)
		case "/debug/pprof/threadcreate":
			pprof.Handler("threadcreate").ServeHTTP(w, r)
		case "/debug/pprof/block":
			pprof.Handler("block").ServeHTTP(w, r)
		case "/debug/pprof/mutex":
			pprof.Handler("mutex").ServeHTTP(w, r)
		case "/debug/pprof/allocs":
			pprof.Handler("allocs").ServeHTTP(w, r)
		default:
			pprof.Index(w, r)
		}
	})
}
