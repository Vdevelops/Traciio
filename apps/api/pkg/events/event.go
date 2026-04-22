package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event
type Event struct {
	// ID is the unique identifier for this event
	ID string `json:"id"`

	// Type is the event type (e.g., "lead.created", "deal.stage_changed")
	Type string `json:"type"`

	// Timestamp is when the event occurred
	Timestamp time.Time `json:"timestamp"`

	// AggregateID is the ID of the aggregate that generated this event
	AggregateID string `json:"aggregate_id"`

	// AggregateType is the type of aggregate (e.g., "lead", "deal", "activity")
	AggregateType string `json:"aggregate_type"`

	// Payload is the event-specific data (JSON encoded)
	Payload json.RawMessage `json:"payload"`

	// Metadata contains additional context (user_id, trace_id, etc.)
	Metadata map[string]string `json:"metadata,omitempty"`

	// Version is the event schema version
	Version string `json:"version"`
}

// NewEvent creates a new event with generated ID and timestamp
func NewEvent(eventType, aggregateID, aggregateType string, payload interface{}, metadata map[string]string) (*Event, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:            uuid.New().String(),
		Type:          eventType,
		Timestamp:     time.Now(),
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Payload:       payloadBytes,
		Metadata:      metadata,
		Version:       "1.0",
	}, nil
}

// EventProducer defines the interface for publishing events
type EventProducer interface {
	// Publish publishes an event to the event bus
	// Returns error if publishing fails
	Publish(event *Event) error

	// PublishAsync publishes an event asynchronously (fire-and-forget)
	// Errors are logged but not returned
	PublishAsync(event *Event)

	// Close closes the producer and releases resources
	Close() error
}

// EventConsumer defines the interface for consuming events
type EventConsumer interface {
	// Subscribe subscribes to events of a specific type
	// Handler will be called for each event
	Subscribe(eventType string, handler EventHandler) error

	// Start starts consuming events
	Start() error

	// Stop stops consuming events
	Stop() error

	// Close closes the consumer and releases resources
	Close() error
}

// EventHandler is a function that handles an event
type EventHandler func(event *Event) error
