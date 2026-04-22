// Package circuitbreaker provides circuit breaker pattern for external service calls.
//
// This package implements the circuit breaker pattern to prevent cascading failures
// when external services are down or experiencing high latency. It supports three states:
// - Closed: Normal operation, requests pass through
// - Open: Service is failing, requests fail fast
// - Half-Open: Testing if service recovered
package circuitbreaker

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

var (
	// ErrCircuitOpen is returned when the circuit breaker is open
	ErrCircuitOpen = errors.New("circuit breaker is open")

	// ErrServiceUnavailable is returned when external service is unavailable
	ErrServiceUnavailable = errors.New("external service unavailable")
)

// State represents the state of circuit breaker
type State int

const (
	// StateClosed - normal operation
	StateClosed State = iota
	// StateOpen - failing, reject requests
	StateOpen
	// StateHalfOpen - testing recovery
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Config holds circuit breaker configuration
type Config struct {
	// Name is the identifier for this circuit breaker (used in logs/metrics)
	Name string

	// MaxRequests is the maximum number of requests allowed to pass through
	// when the circuit is half-open
	MaxRequests uint32

	// Interval is the cyclic period of the closed state
	// If 0, it doesn't clear counts during closed state
	Interval time.Duration

	// Timeout is the period of the open state
	// After this period, the circuit moves to half-open
	Timeout time.Duration

	// FailureRatio is the threshold for failure ratio to trip the circuit
	// e.g., 0.6 means if 60% of requests fail, open the circuit
	FailureRatio float64

	// MinRequests is the minimum number of requests required before calculating failure ratio
	MinRequests uint32

	// SuccessOn is the number of consecutive successes needed to close circuit from half-open
	SuccessOn uint32
}

// DefaultConfig returns default circuit breaker configuration
func DefaultConfig(name string) Config {
	return Config{
		Name:         name,
		MaxRequests:  100,
		Interval:     0, // Don't clear counts during closed state
		Timeout:      60 * time.Second,
		FailureRatio: 0.6, // Trip at 60% failure rate
		MinRequests:  10,
		SuccessOn:    5,
	}
}

// CircuitBreaker wraps sony/gobreaker with additional functionality
type CircuitBreaker struct {
	cb     *gobreaker.CircuitBreaker
	config Config
	name   string
	mu     sync.RWMutex
	state  State
}

// Metrics holds circuit breaker metrics
type Metrics struct {
	State                State
	CountRequests        uint32
	CountSuccesses       uint32
	CountFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// New creates a new circuit breaker with the given configuration
func New(cfg Config) *CircuitBreaker {
	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= cfg.MinRequests && failureRatio >= cfg.FailureRatio
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("[CircuitBreaker] %s: state changed from %s to %s", name, from.String(), to.String())
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	return &CircuitBreaker{
		cb:     cb,
		config: cfg,
		name:   cfg.Name,
	}
}

// Execute runs the given function if the circuit allows it
// Returns ErrCircuitOpen if the circuit is open
func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
	result, err := cb.cb.Execute(req)
	if err != nil {
		// Check if it's a circuit breaker error
		if errors.Is(err, gobreaker.ErrOpenState) {
			return nil, fmt.Errorf("%w: %s", ErrCircuitOpen, cb.name)
		}
		return nil, err
	}
	return result, nil
}

// ExecuteWithFallback runs the function if circuit allows, otherwise runs fallback
func (cb *CircuitBreaker) ExecuteWithFallback(
	req func() (interface{}, error),
	fallback func(error) (interface{}, error),
) (interface{}, error) {
	result, err := cb.Execute(req)
	if err != nil {
		if errors.Is(err, ErrCircuitOpen) || errors.Is(err, ErrServiceUnavailable) {
			log.Printf("[CircuitBreaker] %s: executing fallback due to: %v", cb.name, err)
			return fallback(err)
		}
		return nil, err
	}
	return result, nil
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	return State(cb.cb.State())
}

// Name returns the name of this circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// Counts returns the current counts
func (cb *CircuitBreaker) Counts() gobreaker.Counts {
	return cb.cb.Counts()
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.cb.State() == gobreaker.StateOpen
}

// IsClosed returns true if the circuit is closed
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.cb.State() == gobreaker.StateClosed
}

// Registry manages multiple circuit breakers
type Registry struct {
	breakers map[string]*CircuitBreaker
	mu       sync.RWMutex
}

// NewRegistry creates a new circuit breaker registry
func NewRegistry() *Registry {
	return &Registry{
		breakers: make(map[string]*CircuitBreaker),
	}
}

// Register registers a circuit breaker
func (r *Registry) Register(cb *CircuitBreaker) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.breakers[cb.Name()] = cb
}

// Get retrieves a circuit breaker by name
func (r *Registry) Get(name string) (*CircuitBreaker, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cb, ok := r.breakers[name]
	return cb, ok
}

// GetOrCreate gets or creates a circuit breaker with default config
func (r *Registry) GetOrCreate(name string) *CircuitBreaker {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cb, ok := r.breakers[name]; ok {
		return cb
	}

	cb := New(DefaultConfig(name))
	r.breakers[name] = cb
	return cb
}

// GetOrCreateWithConfig gets or creates a circuit breaker with custom config
func (r *Registry) GetOrCreateWithConfig(cfg Config) *CircuitBreaker {
	r.mu.Lock()
	defer r.mu.Unlock()

	if cb, ok := r.breakers[cfg.Name]; ok {
		return cb
	}

	cb := New(cfg)
	r.breakers[cfg.Name] = cb
	return cb
}

// All returns all registered circuit breakers
func (r *Registry) All() map[string]*CircuitBreaker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy
	result := make(map[string]*CircuitBreaker, len(r.breakers))
	for k, v := range r.breakers {
		result[k] = v
	}
	return result
}

// Stats returns statistics for all circuit breakers
func (r *Registry) Stats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	stats := make(map[string]interface{})
	for name, cb := range r.breakers {
		counts := cb.Counts()
		stats[name] = map[string]interface{}{
			"state":                 cb.State().String(),
			"requests":              counts.Requests,
			"total_successes":       counts.TotalSuccesses,
			"total_failures":        counts.TotalFailures,
			"consecutive_successes": counts.ConsecutiveSuccesses,
			"consecutive_failures":  counts.ConsecutiveFailures,
		}
	}
	return stats
}

// global registry instance
var globalRegistry = NewRegistry()

// Register registers a circuit breaker to the global registry
func Register(cb *CircuitBreaker) {
	globalRegistry.Register(cb)
}

// Get retrieves a circuit breaker from the global registry
func Get(name string) (*CircuitBreaker, bool) {
	return globalRegistry.Get(name)
}

// GetOrCreate gets or creates a circuit breaker in the global registry
func GetOrCreate(name string) *CircuitBreaker {
	return globalRegistry.GetOrCreate(name)
}

// Stats returns statistics for all circuit breakers in global registry
func Stats() map[string]interface{} {
	return globalRegistry.Stats()
}
