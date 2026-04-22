package google_calendar

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

var (
	ErrInvalidCredentials  = errors.New("invalid Google Calendar credentials")
	ErrEventNotFound        = errors.New("Google Calendar event not found")
	ErrEventCreationFailed  = errors.New("failed to create Google Calendar event")
	ErrEventUpdateFailed    = errors.New("failed to update Google Calendar event")
	ErrEventDeleteFailed    = errors.New("failed to delete Google Calendar event")
	ErrRateLimitExceeded    = errors.New("Google Calendar API rate limit exceeded")
)

// RateLimiter handles rate limiting for Google Calendar API calls
type RateLimiter struct {
	requests     chan struct{}
	lastReset    time.Time
	resetInterval time.Duration
	maxRequests  int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, resetInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests:      make(chan struct{}, maxRequests),
		lastReset:     time.Now(),
		resetInterval: resetInterval,
		maxRequests:   maxRequests,
	}

	for i := 0; i < maxRequests; i++ {
		rl.requests <- struct{}{}
	}

	return rl
}

// Wait waits for an available request slot
func (rl *RateLimiter) Wait(ctx context.Context) error {
	if time.Since(rl.lastReset) >= rl.resetInterval {
		for len(rl.requests) > 0 {
			<-rl.requests
		}
		for i := 0; i < rl.maxRequests; i++ {
			select {
			case rl.requests <- struct{}{}:
			default:
			}
		}
		rl.lastReset = time.Now()
	}

	select {
	case <-rl.requests:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Client represents Google Calendar API client
type Client struct {
	service     *calendar.Service
	rateLimiter *RateLimiter
}

// NewClient creates a new Google Calendar client with OAuth2 token source
// tokenSource handles token refresh automatically
func NewClient(ctx context.Context, tokenSource oauth2.TokenSource) (*Client, error) {
	httpClient := oauth2.NewClient(ctx, tokenSource)

	service, err := calendar.New(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create calendar service: %w", err)
	}

	rateLimiter := NewRateLimiter(1000, 100*time.Second)

	return &Client{
		service:     service,
		rateLimiter: rateLimiter,
	}, nil
}

// CreateEvent creates a new event in Google Calendar
func (c *Client) CreateEvent(ctx context.Context, calendarID string, event *calendar.Event) (*calendar.Event, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	createdEvent, err := c.service.Events.Insert(calendarID, event).Context(ctx).Do()
	if err != nil {
		if err.Error() == "rateLimitExceeded" {
			return nil, ErrRateLimitExceeded
		}
		return nil, fmt.Errorf("%w: %v", ErrEventCreationFailed, err)
	}

	return createdEvent, nil
}

// UpdateEvent updates an existing event in Google Calendar
func (c *Client) UpdateEvent(ctx context.Context, calendarID string, eventID string, event *calendar.Event) (*calendar.Event, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	updatedEvent, err := c.service.Events.Update(calendarID, eventID, event).Context(ctx).Do()
	if err != nil {
		if err.Error() == "rateLimitExceeded" {
			return nil, ErrRateLimitExceeded
		}
		return nil, fmt.Errorf("%w: %v", ErrEventUpdateFailed, err)
	}

	return updatedEvent, nil
}

// DeleteEvent deletes an event from Google Calendar
func (c *Client) DeleteEvent(ctx context.Context, calendarID string, eventID string) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait failed: %w", err)
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	err := c.service.Events.Delete(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		if err.Error() == "rateLimitExceeded" {
			return ErrRateLimitExceeded
		}
		return fmt.Errorf("%w: %v", ErrEventDeleteFailed, err)
	}

	return nil
}

// GetEvent retrieves an event from Google Calendar
func (c *Client) GetEvent(ctx context.Context, calendarID string, eventID string) (*calendar.Event, error) {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}

	if calendarID == "" {
		calendarID = "primary"
	}

	event, err := c.service.Events.Get(calendarID, eventID).Context(ctx).Do()
	if err != nil {
		if err.Error() == "rateLimitExceeded" {
			return nil, ErrRateLimitExceeded
		}
		return nil, fmt.Errorf("%w: %v", ErrEventNotFound, err)
	}

	return event, nil
}

// BuildEventFromSchedule builds a Google Calendar event from schedule data
func BuildEventFromSchedule(title, description string, scheduledAt time.Time, reminderMinutesBefore *int) *calendar.Event {
	event := &calendar.Event{
		Summary:     title,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: scheduledAt.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			DateTime: scheduledAt.Add(1 * time.Hour).Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
	}

	// Set custom reminders if provided, otherwise use calendar defaults
	if reminderMinutesBefore != nil && *reminderMinutesBefore > 0 {
		event.Reminders = &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{
					Method:  "email",
					Minutes: int64(*reminderMinutesBefore),
				},
				{
					Method:  "popup",
					Minutes: int64(*reminderMinutesBefore),
				},
			},
		}
	}

	return event
}

// GetEventURL returns the URL to view the event in Google Calendar.
// Event ID must be base64url encoded.
func GetEventURL(eventID string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(eventID))
	encoded = strings.TrimRight(encoded, "=")
	return fmt.Sprintf("https://calendar.google.com/calendar/event?eid=%s", encoded)
}

