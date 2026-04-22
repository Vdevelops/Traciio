package events

import (
	"log"
	"sync"
)

// NoOpConsumer is a consumer that doesn't actually consume events
// Used when Kafka is disabled or for development/testing
type NoOpConsumer struct {
	logger   *log.Logger
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewNoOpConsumer creates a new NoOpConsumer
func NewNoOpConsumer() *NoOpConsumer {
	return &NoOpConsumer{
		logger:   log.New(log.Writer(), "[EVENT-CONSUMER] ", log.LstdFlags),
		handlers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event type
func (c *NoOpConsumer) Subscribe(eventType string, handler EventHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handlers[eventType] = append(c.handlers[eventType], handler)
	c.logger.Printf("NoOpConsumer - Subscribed to event type: %s", eventType)

	return nil
}

// Start is a no-op for NoOpConsumer
func (c *NoOpConsumer) Start() error {
	c.logger.Println("NoOpConsumer started (no actual consumption)")
	return nil
}

// Stop is a no-op for NoOpConsumer
func (c *NoOpConsumer) Stop() error {
	c.logger.Println("NoOpConsumer stopped")
	return nil
}

// Close is a no-op for NoOpConsumer
func (c *NoOpConsumer) Close() error {
	c.logger.Println("NoOpConsumer closed")
	return nil
}

// KafkaConsumer is a placeholder for future Kafka implementation
// TODO: Implement when Kafka is installed
type KafkaConsumer struct {
	// brokers []string
	// consumerGroup sarama.ConsumerGroup
	// topicPrefix string
	// handlers map[string][]EventHandler
}

// NewKafkaConsumer creates a new KafkaConsumer
// TODO: Implement when Kafka is installed
func NewKafkaConsumer(brokers []string, groupID, topicPrefix string) (*KafkaConsumer, error) {
	// TODO: Initialize Kafka consumer
	// config := sarama.NewConfig()
	// config.Consumer.Return.Errors = true
	// consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	// if err != nil {
	//     return nil, err
	// }
	// return &KafkaConsumer{
	//     brokers: brokers,
	//     consumerGroup: consumerGroup,
	//     topicPrefix: topicPrefix,
	//     handlers: make(map[string][]EventHandler),
	// }, nil
	return nil, nil
}

// Subscribe registers a handler for a specific event type
// TODO: Implement when Kafka is installed
func (c *KafkaConsumer) Subscribe(eventType string, handler EventHandler) error {
	// TODO: Implement Kafka subscription
	// c.handlers[eventType] = append(c.handlers[eventType], handler)
	return nil
}

// Start starts consuming events from Kafka
// TODO: Implement when Kafka is installed
func (c *KafkaConsumer) Start() error {
	// TODO: Start Kafka consumer
	// go func() {
	//     for {
	//         select {
	//         case message := <-c.consumerGroup.Messages():
	//             // Parse event and call handlers
	//         case err := <-c.consumerGroup.Errors():
	//             log.Printf("Kafka consumer error: %v", err)
	//         }
	//     }
	// }()
	return nil
}

// Stop stops consuming events
// TODO: Implement when Kafka is installed
func (c *KafkaConsumer) Stop() error {
	// TODO: Stop Kafka consumer
	return nil
}

// Close closes the Kafka consumer
// TODO: Implement when Kafka is installed
func (c *KafkaConsumer) Close() error {
	// TODO: Close Kafka consumer
	// return c.consumerGroup.Close()
	return nil
}
