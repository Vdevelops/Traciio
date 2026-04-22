// Package workerpool provides a worker pool implementation for background task processing.
//
// This package implements a fixed-size worker pool pattern to prevent goroutine explosion
// during high load. It's designed for processing notifications, emails, and other
// background tasks.
package workerpool

import (
	"context"
	"log"
	"sync"
	"time"
)

// Task represents a unit of work to be processed
type Task interface {
	Execute() error
}

// TaskFunc is a function type that implements Task interface
type TaskFunc func() error

// Execute implements Task interface
func (f TaskFunc) Execute() error {
	return f()
}

// Pool represents a worker pool with fixed number of workers
type Pool struct {
	workers     int
	queue       chan Task
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	stopped     bool
	mu          sync.RWMutex
	taskTimeout time.Duration
}

// Config holds configuration for worker pool
type Config struct {
	Workers     int
	QueueSize   int
	TaskTimeout time.Duration
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		Workers:     10,
		QueueSize:   100,
		TaskTimeout: 30 * time.Second,
	}
}

// New creates a new worker pool
func New(cfg Config) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		workers:     cfg.Workers,
		queue:       make(chan Task, cfg.QueueSize),
		ctx:         ctx,
		cancel:      cancel,
		taskTimeout: cfg.TaskTimeout,
	}

	// Start workers
	for i := 0; i < cfg.Workers; i++ {
		pool.wg.Add(1)
		go pool.worker(i)
	}

	log.Printf("[WorkerPool] Started %d workers with queue size %d", cfg.Workers, cfg.QueueSize)

	return pool
}

// worker is the goroutine that processes tasks
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			log.Printf("[WorkerPool] Worker %d stopped", id)
			return

		case task, ok := <-p.queue:
			if !ok {
				log.Printf("[WorkerPool] Worker %d: queue closed", id)
				return
			}

			// Execute task with timeout
			if p.taskTimeout > 0 {
				taskCtx, cancel := context.WithTimeout(p.ctx, p.taskTimeout)
				done := make(chan error, 1)

				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("[WorkerPool] Worker %d: task panicked: %v", id, r)
							done <- nil
						}
					}()
					done <- task.Execute()
				}()

				select {
				case err := <-done:
					if err != nil {
						log.Printf("[WorkerPool] Worker %d: task failed: %v", id, err)
					}
				case <-taskCtx.Done():
					log.Printf("[WorkerPool] Worker %d: task timed out", id)
				}

				cancel()
			} else {
				// Execute without timeout
				err := task.Execute()
				if err != nil {
					log.Printf("[WorkerPool] Worker %d: task failed: %v", id, err)
				}
			}
		}
	}
}

// Submit submits a task to the pool
// Returns error if pool is stopped or queue is full
func (p *Pool) Submit(task Task) error {
	p.mu.RLock()
	if p.stopped {
		p.mu.RUnlock()
		return ErrPoolStopped
	}
	p.mu.RUnlock()

	select {
	case p.queue <- task:
		return nil
	case <-p.ctx.Done():
		return ErrPoolStopped
	default:
		return ErrQueueFull
	}
}

// SubmitAsync submits a task asynchronously (non-blocking)
// Returns immediately with boolean indicating if task was queued
func (p *Pool) SubmitAsync(task Task) bool {
	p.mu.RLock()
	if p.stopped {
		p.mu.RUnlock()
		return false
	}
	p.mu.RUnlock()

	select {
	case p.queue <- task:
		return true
	default:
		return false
	}
}

// Stop gracefully shuts down the worker pool
func (p *Pool) Stop() {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return
	}
	p.stopped = true
	p.mu.Unlock()

	log.Println("[WorkerPool] Stopping...")

	// Cancel context to signal workers to stop
	p.cancel()

	// Wait for all workers to finish
	p.wg.Wait()

	log.Println("[WorkerPool] Stopped")
}

// StopWithTimeout stops the pool with a timeout
func (p *Pool) StopWithTimeout(timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		p.Stop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[WorkerPool] Stopped gracefully")
	case <-time.After(timeout):
		log.Println("[WorkerPool] Stop timeout exceeded, forcing shutdown")
	}
}

// Stats returns current pool statistics
func (p *Pool) Stats() Stats {
	return Stats{
		Workers:   p.workers,
		QueueSize: cap(p.queue),
		QueueLen:  len(p.queue),
		IsStopped: p.stopped,
	}
}

// Stats represents pool statistics
type Stats struct {
	Workers   int  `json:"workers"`
	QueueSize int  `json:"queue_size"`
	QueueLen  int  `json:"queue_len"`
	IsStopped bool `json:"is_stopped"`
}

var (
	// ErrPoolStopped is returned when trying to submit to a stopped pool
	ErrPoolStopped = &PoolError{Message: "worker pool is stopped"}
	// ErrQueueFull is returned when the task queue is full
	ErrQueueFull = &PoolError{Message: "task queue is full"}
)

// PoolError represents pool errors
type PoolError struct {
	Message string
}

func (e *PoolError) Error() string {
	return e.Message
}
