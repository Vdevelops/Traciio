package events

import (
	"encoding/json"
	"log"
)

// NoOpProducer is a producer that logs events but doesn't publish them
// Used when Kafka is disabled or for development/testing
type NoOpProducer struct {
	logger *log.Logger
}

// NewNoOpProducer creates a new NoOpProducer
func NewNoOpProducer() *NoOpProducer {
	return &NoOpProducer{
		logger: log.New(log.Writer(), "[EVENT] ", log.LstdFlags),
	}
}

// Publish logs the event but doesn't actually publish it
func (p *NoOpProducer) Publish(event *Event) error {
	eventJSON, err := json.Marshal(event)
	if err != nil {
		p.logger.Printf("ERROR: Failed to marshal event: %v", err)
		return err
	}

	p.logger.Printf("NoOpProducer - Event: %s | Type: %s | AggregateID: %s | Payload: %s",
		event.ID, event.Type, event.AggregateID, string(eventJSON))

	return nil
}

// PublishAsync logs the event asynchronously
func (p *NoOpProducer) PublishAsync(event *Event) {
	go func() {
		if err := p.Publish(event); err != nil {
			p.logger.Printf("ERROR: Failed to publish event asynchronously: %v", err)
		}
	}()
}

// Close is a no-op for NoOpProducer
func (p *NoOpProducer) Close() error {
	p.logger.Println("NoOpProducer closed")
	return nil
}

// KafkaProducer is a placeholder for future Kafka implementation
// TODO: Implement when Kafka is installed
type KafkaProducer struct {
	// brokers []string
	// producer sarama.SyncProducer
	// topicPrefix string
}

// NewKafkaProducer creates a new KafkaProducer
// TODO: Implement when Kafka is installed
func NewKafkaProducer(brokers []string, topicPrefix string) (*KafkaProducer, error) {
	// TODO: Initialize Kafka producer
	// producer, err := sarama.NewSyncProducer(brokers, config)
	// if err != nil {
	//     return nil, err
	// }
	// return &KafkaProducer{
	//     brokers: brokers,
	//     producer: producer,
	//     topicPrefix: topicPrefix,
	// }, nil
	return nil, nil
}

// Publish publishes an event to Kafka
// TODO: Implement when Kafka is installed
func (p *KafkaProducer) Publish(event *Event) error {
	// TODO: Implement Kafka publish
	// topic := p.topicPrefix + "." + event.AggregateType + "." + event.Type
	// message := &sarama.ProducerMessage{
	//     Topic: topic,
	//     Key:   sarama.StringEncoder(event.AggregateID),
	//     Value: sarama.ByteEncoder(eventJSON),
	// }
	// _, _, err := p.producer.SendMessage(message)
	// return err
	return nil
}

// PublishAsync publishes an event to Kafka asynchronously
// TODO: Implement when Kafka is installed
func (p *KafkaProducer) PublishAsync(event *Event) {
	// TODO: Implement async Kafka publish
	go func() {
		_ = p.Publish(event)
	}()
}

// Close closes the Kafka producer
// TODO: Implement when Kafka is installed
func (p *KafkaProducer) Close() error {
	// TODO: Close Kafka producer
	// return p.producer.Close()
	return nil
}
