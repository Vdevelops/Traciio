// Package cache provides Redis caching layer with high availability support.
//
// This package implements a production-ready Redis client that supports
// standalone, Sentinel, and Cluster deployment modes. It includes circuit
// breaker pattern for fault tolerance and MessagePack serialization for
// optimal performance.
package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker"
	"github.com/vmihailenco/msgpack/v5"
)

// Cache key prefixes for different data types.
const (
	PrefixUserPermissions = "user:permissions:"
	PrefixRolePermissions = "role:permissions:"
	PrefixSession         = "session:"
	PrefixRateLimitIP     = "ratelimit:ip:"
	PrefixRateLimitUser   = "ratelimit:user:"
)

// Mode represents the Redis deployment topology.
type Mode string

const (
	// ModeStandalone indicates a single Redis instance deployment.
	ModeStandalone Mode = "standalone"
	// ModeSentinel indicates Redis Sentinel deployment for high availability.
	ModeSentinel Mode = "sentinel"
	// ModeCluster indicates Redis Cluster deployment for horizontal scaling.
	ModeCluster Mode = "cluster"
)

// Config holds the configuration for the Redis cache layer.
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
	Enabled  bool

	// High Availability
	Mode             Mode
	SentinelAddrs    []string
	SentinelMaster   string
	ClusterAddrs     []string
	SentinelPassword string

	// Connection Pool
	PoolSize        int
	MinIdleConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolTimeout     time.Duration

	// Circuit Breaker
	CBMaxRequests  uint32
	CBInterval     time.Duration
	CBTimeout      time.Duration
	CBFailureRatio float64
	CBMinRequests  uint32

	// TTL
	PermissionsTTL int
	SessionTTL     int
	GeneralTTL     int

	// Compression
	CompressionThreshold int
}

// Cache provides thread-safe Redis operations with circuit breaker protection.
type Cache struct {
	client  redis.UniversalClient
	cb      *gobreaker.CircuitBreaker
	metrics *Metrics
	cfg     *Config
	ctx     context.Context
	cancel  context.CancelFunc
	enabled bool
	mode    Mode
}

// Metrics tracks cache performance statistics.
type Metrics struct {
	mu sync.RWMutex

	Hits           uint64
	Misses         uint64
	Errors         uint64
	Timeouts       uint64
	CBOpen         uint64
	CBHalfOpen     uint64
	PipelineOps    uint64
	Compressions   uint64
	TotalLatency   uint64
	OperationCount uint64
	MaxLatency     uint64

	PoolHits     uint64
	PoolMisses   uint64
	PoolTimeouts uint64
	TotalConns   uint64
	IdleConns    uint64
	StaleConns   uint64
}

// Client is the global cache instance.
var Client *Cache

// New creates a new Redis cache instance based on the provided configuration.
// It automatically detects the deployment mode from environment variables
// and configures the client accordingly.
func New(cfg *config.RedisConfig) (*Cache, error) {
	if !cfg.Enabled {
		log.Println("[Cache] Disabled by configuration")
		Client = &Cache{enabled: false}
		return Client, nil
	}

	internalCfg := buildConfig(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	client, err := newClient(internalCfg)
	if err != nil {
		cancel()
		log.Printf("[Cache] Connection failed: %v", err)
		Client = &Cache{enabled: false, ctx: ctx, cancel: cancel}
		return Client, nil
	}

	if err := client.Ping(ctx).Err(); err != nil {
		cancel()
		log.Printf("[Cache] Ping failed: %v", err)
		Client = &Cache{enabled: false, ctx: ctx, cancel: cancel}
		return Client, nil
	}

	Client = &Cache{
		client:  client,
		cb:      newCircuitBreaker(internalCfg),
		metrics: &Metrics{},
		cfg:     internalCfg,
		ctx:     ctx,
		cancel:  cancel,
		enabled: true,
		mode:    internalCfg.Mode,
	}

	go Client.collectPoolStats()

	log.Printf("[Cache] Connected (%s mode)", internalCfg.Mode)
	return Client, nil
}

// buildConfig creates internal configuration from application config.
func buildConfig(cfg *config.RedisConfig) *Config {
	c := &Config{
		Host:           cfg.Host,
		Port:           cfg.Port,
		Password:       cfg.Password,
		DB:             cfg.DB,
		Enabled:        cfg.Enabled,
		PermissionsTTL: cfg.PermissionsTTL,
		SessionTTL:     cfg.SessionTTL,
		GeneralTTL:     cfg.GeneralTTL,
		Mode:           ModeStandalone,

		PoolSize:        200,
		MinIdleConns:    50,
		MaxIdleConns:    100,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		DialTimeout:     5 * time.Second,
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,

		CBMaxRequests:  5,
		CBInterval:     60 * time.Second,
		CBTimeout:      30 * time.Second,
		CBFailureRatio: 0.6,
		CBMinRequests:  10,

		CompressionThreshold: 1024,
	}

	// Sentinel mode configuration
	if addrs := os.Getenv("REDIS_SENTINEL_ADDRS"); addrs != "" {
		c.Mode = ModeSentinel
		c.SentinelAddrs = strings.Split(addrs, ",")
		c.SentinelMaster = getEnvOrDefault("REDIS_SENTINEL_MASTER", "mymaster")
		c.SentinelPassword = getEnvOrDefault("REDIS_SENTINEL_PASSWORD", cfg.Password)
	}

	// Cluster mode configuration
	if addrs := os.Getenv("REDIS_CLUSTER_ADDRS"); addrs != "" {
		c.Mode = ModeCluster
		c.ClusterAddrs = strings.Split(addrs, ",")
	}

	// Override pool settings from environment
	if v := getEnvAsInt("REDIS_POOL_SIZE"); v > 0 {
		c.PoolSize = v
	}
	if v := getEnvAsInt("REDIS_MIN_IDLE_CONNS"); v > 0 {
		c.MinIdleConns = v
	}

	return c
}

// newClient creates the appropriate Redis client based on deployment mode.
func newClient(cfg *Config) (redis.UniversalClient, error) {
	switch cfg.Mode {
	case ModeSentinel:
		return redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:       cfg.SentinelMaster,
			SentinelAddrs:    cfg.SentinelAddrs,
			SentinelPassword: cfg.SentinelPassword,
			Password:         cfg.Password,
			DB:               cfg.DB,
			PoolSize:         cfg.PoolSize,
			MinIdleConns:     cfg.MinIdleConns,
			ConnMaxLifetime:  cfg.ConnMaxLifetime,
			ConnMaxIdleTime:  cfg.ConnMaxIdleTime,
			DialTimeout:      cfg.DialTimeout,
			ReadTimeout:      cfg.ReadTimeout,
			WriteTimeout:     cfg.WriteTimeout,
			PoolTimeout:      cfg.PoolTimeout,
		}), nil

	case ModeCluster:
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           cfg.ClusterAddrs,
			Password:        cfg.Password,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolTimeout:     cfg.PoolTimeout,
		}), nil

	default:
		return redis.NewClient(&redis.Options{
			Addr:            fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password:        cfg.Password,
			DB:              cfg.DB,
			PoolSize:        cfg.PoolSize,
			MinIdleConns:    cfg.MinIdleConns,
			ConnMaxLifetime: cfg.ConnMaxLifetime,
			ConnMaxIdleTime: cfg.ConnMaxIdleTime,
			DialTimeout:     cfg.DialTimeout,
			ReadTimeout:     cfg.ReadTimeout,
			WriteTimeout:    cfg.WriteTimeout,
			PoolTimeout:     cfg.PoolTimeout,
		}), nil
	}
}

// newCircuitBreaker creates a circuit breaker with the specified settings.
func newCircuitBreaker(cfg *Config) *gobreaker.CircuitBreaker {
	return gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "redis",
		MaxRequests: cfg.CBMaxRequests,
		Interval:    cfg.CBInterval,
		Timeout:     cfg.CBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests < cfg.CBMinRequests {
				return false
			}
			return float64(counts.TotalFailures)/float64(counts.Requests) >= cfg.CBFailureRatio
		},
		OnStateChange: func(name string, from, to gobreaker.State) {
			log.Printf("[Cache] Circuit breaker: %s -> %s", from.String(), to.String())
		},
	})
}

// IsEnabled returns true if the cache is active and connected.
func (c *Cache) IsEnabled() bool {
	return c.enabled && c.client != nil
}

// execute wraps Redis operations with circuit breaker and metrics tracking.
func (c *Cache) execute(op func() (interface{}, error)) (interface{}, error) {
	if !c.IsEnabled() {
		return nil, nil
	}

	start := time.Now()
	result, err := c.cb.Execute(op)
	latency := time.Since(start).Microseconds()

	atomic.AddUint64(&c.metrics.TotalLatency, uint64(latency))
	atomic.AddUint64(&c.metrics.OperationCount, 1)

	// Update max latency using CAS
	for {
		current := atomic.LoadUint64(&c.metrics.MaxLatency)
		if uint64(latency) <= current || atomic.CompareAndSwapUint64(&c.metrics.MaxLatency, current, uint64(latency)) {
			break
		}
	}

	if err != nil {
		switch {
		case errors.Is(err, gobreaker.ErrOpenState):
			atomic.AddUint64(&c.metrics.CBOpen, 1)
		case errors.Is(err, gobreaker.ErrTooManyRequests):
			atomic.AddUint64(&c.metrics.CBHalfOpen, 1)
		default:
			atomic.AddUint64(&c.metrics.Errors, 1)
		}
	}

	return result, err
}

// Set stores a value with the specified TTL.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	_, err := c.execute(func() (interface{}, error) {
		data, err := msgpack.Marshal(value)
		if err != nil {
			return nil, fmt.Errorf("marshal failed: %w", err)
		}

		if len(data) > c.cfg.CompressionThreshold {
			atomic.AddUint64(&c.metrics.Compressions, 1)
		}

		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()

		return nil, c.client.Set(ctx, key, data, ttl).Err()
	})
	return err
}

// Get retrieves a value and unmarshals it into the target.
// Returns true if the key exists, false otherwise.
func (c *Cache) Get(key string, target interface{}) (bool, error) {
	result, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.ReadTimeout)
		defer cancel()

		data, err := c.client.Get(ctx, key).Bytes()
		if err == redis.Nil {
			return false, nil
		}
		if err != nil {
			return nil, fmt.Errorf("get failed: %w", err)
		}

		if err := msgpack.Unmarshal(data, target); err != nil {
			return nil, fmt.Errorf("unmarshal failed: %w", err)
		}

		return true, nil
	})

	if err != nil {
		return false, err
	}

	if result == nil {
		atomic.AddUint64(&c.metrics.Misses, 1)
		return false, nil
	}

	found, _ := result.(bool)
	if found {
		atomic.AddUint64(&c.metrics.Hits, 1)
	} else {
		atomic.AddUint64(&c.metrics.Misses, 1)
	}

	return found, nil
}

// Delete removes a key from the cache.
func (c *Cache) Delete(key string) error {
	_, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()
		return nil, c.client.Del(ctx, key).Err()
	})
	return err
}

// MGet retrieves multiple keys in a single round trip.
func (c *Cache) MGet(keys []string, targets []interface{}) ([]bool, error) {
	if !c.IsEnabled() || len(keys) == 0 {
		return make([]bool, len(keys)), nil
	}

	atomic.AddUint64(&c.metrics.PipelineOps, 1)

	result, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.ReadTimeout)
		defer cancel()

		values, err := c.client.MGet(ctx, keys...).Result()
		if err != nil {
			return nil, err
		}

		found := make([]bool, len(keys))
		for i, val := range values {
			if val == nil {
				found[i] = false
				atomic.AddUint64(&c.metrics.Misses, 1)
				continue
			}

			data, ok := val.(string)
			if !ok {
				found[i] = false
				continue
			}

			if i < len(targets) && targets[i] != nil {
				if err := msgpack.Unmarshal([]byte(data), targets[i]); err != nil {
					found[i] = false
					continue
				}
			}

			found[i] = true
			atomic.AddUint64(&c.metrics.Hits, 1)
		}

		return found, nil
	})

	if err != nil {
		return make([]bool, len(keys)), err
	}

	found, _ := result.([]bool)
	if found == nil {
		return make([]bool, len(keys)), nil
	}

	return found, nil
}

// MSet stores multiple key-value pairs in a single round trip.
func (c *Cache) MSet(items map[string]interface{}, ttl time.Duration) error {
	if !c.IsEnabled() || len(items) == 0 {
		return nil
	}

	atomic.AddUint64(&c.metrics.PipelineOps, 1)

	_, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()

		pipe := c.client.Pipeline()
		for key, value := range items {
			data, err := msgpack.Marshal(value)
			if err != nil {
				continue
			}
			pipe.Set(ctx, key, data, ttl)
		}

		_, err := pipe.Exec(ctx)
		return nil, err
	})

	return err
}

// MDelete removes multiple keys in a single round trip.
func (c *Cache) MDelete(keys []string) error {
	if !c.IsEnabled() || len(keys) == 0 {
		return nil
	}

	atomic.AddUint64(&c.metrics.PipelineOps, 1)

	_, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()
		return nil, c.client.Del(ctx, keys...).Err()
	})

	return err
}

// DeletePattern removes all keys matching the given pattern.
// Uses SCAN to avoid blocking operations on large datasets.
func (c *Cache) DeletePattern(pattern string) error {
	if !c.IsEnabled() {
		return nil
	}

	_, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
		defer cancel()

		var cursor uint64
		var allKeys []string

		for {
			keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
			if err != nil {
				return nil, err
			}

			allKeys = append(allKeys, keys...)
			cursor = nextCursor

			if cursor == 0 || len(allKeys) > 100000 {
				break
			}
		}

		if len(allKeys) == 0 {
			return nil, nil
		}

		// Delete in batches to avoid memory issues
		const batchSize = 1000
		for i := 0; i < len(allKeys); i += batchSize {
			end := i + batchSize
			if end > len(allKeys) {
				end = len(allKeys)
			}

			pipe := c.client.Pipeline()
			for _, key := range allKeys[i:end] {
				pipe.Del(ctx, key)
			}
			if _, err := pipe.Exec(ctx); err != nil {
				return nil, err
			}
		}

		return nil, nil
	})

	return err
}

// Exists checks if a key exists in the cache.
func (c *Cache) Exists(key string) (bool, error) {
	result, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.ReadTimeout)
		defer cancel()

		count, err := c.client.Exists(ctx, key).Result()
		return count > 0, err
	})

	if err != nil {
		return false, err
	}

	exists, _ := result.(bool)
	return exists, nil
}

// Increment atomically increments a key and returns the new value.
func (c *Cache) Increment(key string) (int64, error) {
	result, err := c.execute(func() (interface{}, error) {
		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()
		return c.client.Incr(ctx, key).Result()
	})

	if err != nil {
		return 0, err
	}

	val, _ := result.(int64)
	return val, nil
}

// SetNX sets a value only if the key does not exist.
// Returns true if the key was set, false if it already existed.
func (c *Cache) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	result, err := c.execute(func() (interface{}, error) {
		data, err := msgpack.Marshal(value)
		if err != nil {
			return false, fmt.Errorf("marshal failed: %w", err)
		}

		ctx, cancel := context.WithTimeout(c.ctx, c.cfg.WriteTimeout)
		defer cancel()

		return c.client.SetNX(ctx, key, data, ttl).Result()
	})

	if err != nil {
		return false, err
	}

	set, _ := result.(bool)
	return set, nil
}

// GetMetrics returns a snapshot of current cache metrics.
func (c *Cache) GetMetrics() Metrics {
	c.metrics.mu.RLock()
	defer c.metrics.mu.RUnlock()

	return Metrics{
		Hits:           atomic.LoadUint64(&c.metrics.Hits),
		Misses:         atomic.LoadUint64(&c.metrics.Misses),
		Errors:         atomic.LoadUint64(&c.metrics.Errors),
		Timeouts:       atomic.LoadUint64(&c.metrics.Timeouts),
		CBOpen:         atomic.LoadUint64(&c.metrics.CBOpen),
		CBHalfOpen:     atomic.LoadUint64(&c.metrics.CBHalfOpen),
		PipelineOps:    atomic.LoadUint64(&c.metrics.PipelineOps),
		Compressions:   atomic.LoadUint64(&c.metrics.Compressions),
		TotalLatency:   atomic.LoadUint64(&c.metrics.TotalLatency),
		OperationCount: atomic.LoadUint64(&c.metrics.OperationCount),
		MaxLatency:     atomic.LoadUint64(&c.metrics.MaxLatency),
		PoolHits:       c.metrics.PoolHits,
		PoolMisses:     c.metrics.PoolMisses,
		PoolTimeouts:   c.metrics.PoolTimeouts,
		TotalConns:     c.metrics.TotalConns,
		IdleConns:      c.metrics.IdleConns,
		StaleConns:     c.metrics.StaleConns,
	}
}

// GetHitRate returns the cache hit rate as a percentage.
func (c *Cache) GetHitRate() float64 {
	hits := atomic.LoadUint64(&c.metrics.Hits)
	misses := atomic.LoadUint64(&c.metrics.Misses)
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total) * 100
}

// GetAvgLatency returns the average operation latency in microseconds.
func (c *Cache) GetAvgLatency() float64 {
	total := atomic.LoadUint64(&c.metrics.TotalLatency)
	count := atomic.LoadUint64(&c.metrics.OperationCount)
	if count == 0 {
		return 0
	}
	return float64(total) / float64(count)
}

// collectPoolStats periodically collects connection pool statistics.
func (c *Cache) collectPoolStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if !c.IsEnabled() {
				continue
			}

			var stats *redis.PoolStats
			switch client := c.client.(type) {
			case *redis.Client:
				stats = client.PoolStats()
			default:
				continue
			}

			if stats != nil {
				c.metrics.mu.Lock()
				c.metrics.PoolHits = uint64(stats.Hits)
				c.metrics.PoolMisses = uint64(stats.Misses)
				c.metrics.PoolTimeouts = uint64(stats.Timeouts)
				c.metrics.TotalConns = uint64(stats.TotalConns)
				c.metrics.IdleConns = uint64(stats.IdleConns)
				c.metrics.StaleConns = uint64(stats.StaleConns)
				c.metrics.mu.Unlock()
			}
		}
	}
}

// Health represents the health status of the cache.
type Health struct {
	Healthy bool          `json:"healthy"`
	Message string        `json:"message"`
	Mode    string        `json:"mode"`
	CBState string        `json:"circuit_breaker_state"`
	Latency time.Duration `json:"latency"`
	HitRate float64       `json:"hit_rate_percent"`
	AvgLat  float64       `json:"avg_latency_us"`
	Metrics Metrics       `json:"metrics"`
	Time    time.Time     `json:"timestamp"`
}

// HealthCheck performs a comprehensive health check on the cache.
func (c *Cache) HealthCheck() (*Health, error) {
	h := &Health{
		Mode: string(c.mode),
		Time: time.Now(),
	}

	if !c.IsEnabled() {
		h.Message = "cache disabled"
		return h, nil
	}

	cbState := c.cb.State()
	h.CBState = cbState.String()

	if cbState == gobreaker.StateOpen {
		h.Message = "circuit breaker open"
		return h, fmt.Errorf("circuit breaker open")
	}

	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	start := time.Now()
	if err := c.client.Ping(ctx).Err(); err != nil {
		h.Message = fmt.Sprintf("ping failed: %v", err)
		return h, err
	}
	h.Latency = time.Since(start)

	h.Metrics = c.GetMetrics()
	h.HitRate = c.GetHitRate()
	h.AvgLat = c.GetAvgLatency()
	h.Healthy = true
	h.Message = "ok"

	return h, nil
}

// Close gracefully shuts down the cache connection.
func (c *Cache) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// GetClient returns the underlying Redis client for advanced operations.
func (c *Cache) GetClient() redis.UniversalClient {
	return c.client
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return defaultValue
}

func getEnvAsInt(key string) int {
	v := os.Getenv(key)
	if v == "" {
		return 0
	}
	var result int
	fmt.Sscanf(v, "%d", &result)
	return result
}
