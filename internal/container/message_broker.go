// Package container provides message broker-specific dependency injection container using Google Wire.
package container

import (
	"context"
	"errors"
	"fmt"
)

// MessageBrokerConfig represents configuration for message broker container.
type MessageBrokerConfig struct {
	// Future message broker configurations will go here
	// For example:
	// Kafka  KafkaConfig  `yaml:"kafka"`
	// Redis  RedisConfig  `yaml:"redis"`
	// RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

// MessageBrokerContainer holds message broker-specific dependencies.
type MessageBrokerContainer struct {
	*BaseContainer
	// Future message broker services will go here
	// For example:
	// KafkaProducer   kafka.Producer
	// KafkaConsumer   kafka.Consumer
	// RedisPublisher  redis.Publisher
	// RedisSubscriber redis.Subscriber
}

// Close gracefully shuts down the message broker container and all its services.
func (c *MessageBrokerContainer) Close() error {
	var errs []error

	// Future: Shutdown message broker connections
	// if c.KafkaProducer != nil {
	//     if err := c.KafkaProducer.Close(); err != nil {
	//         errs = append(errs, fmt.Errorf("closing kafka producer: %w", err))
	//     }
	// }

	// Close base container (database, etc.)
	if c.BaseContainer != nil {
		if err := c.BaseContainer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("closing base container: %w", err))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

// Start starts the message broker container services.
func (c *MessageBrokerContainer) Start(ctx context.Context) error {
	// Future: Start message broker services
	// if c.KafkaConsumer != nil {
	//     go c.KafkaConsumer.Start(ctx)
	// }

	return nil
}

// Health checks the health of message broker services.
func (c *MessageBrokerContainer) Health(ctx context.Context) error {
	// Future: Health checks for message broker services
	// if c.KafkaProducer != nil {
	//     if err := c.KafkaProducer.Ping(ctx); err != nil {
	//         return fmt.Errorf("kafka producer health check failed: %w", err)
	//     }
	// }

	return nil
}
