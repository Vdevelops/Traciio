package worker

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
)

// RefreshTokenCleanupWorker handles cleanup of expired refresh tokens
type RefreshTokenCleanupWorker struct {
	refreshTokenRepo interfaces.RefreshTokenRepository
	ticker           *time.Ticker
	stopChan         chan bool
}

// NewRefreshTokenCleanupWorker creates a new refresh token cleanup worker
func NewRefreshTokenCleanupWorker(
	refreshTokenRepo interfaces.RefreshTokenRepository,
	interval time.Duration,
) *RefreshTokenCleanupWorker {
	return &RefreshTokenCleanupWorker{
		refreshTokenRepo: refreshTokenRepo,
		ticker:           time.NewTicker(interval),
		stopChan:         make(chan bool),
	}
}

// Start starts the refresh token cleanup worker
func (w *RefreshTokenCleanupWorker) Start() {
	log.Println("Refresh token cleanup worker started")
	
	go func() {
		for {
			select {
			case <-w.ticker.C:
				w.cleanupExpiredTokens()
			case <-w.stopChan:
				w.ticker.Stop()
				log.Println("Refresh token cleanup worker stopped")
				return
			}
		}
	}()
}

// Stop stops the refresh token cleanup worker
func (w *RefreshTokenCleanupWorker) Stop() {
	w.stopChan <- true
}

// cleanupExpiredTokens deletes expired refresh tokens from the database
func (w *RefreshTokenCleanupWorker) cleanupExpiredTokens() {
	err := w.refreshTokenRepo.DeleteExpired()
	if err != nil {
		log.Printf("Error cleaning up expired refresh tokens: %v", err)
		return
	}
	
	log.Println("Expired refresh tokens cleaned up successfully")
}
