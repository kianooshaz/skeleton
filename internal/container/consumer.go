// Package container provides consumer-specific dependency injection container using Google Wire.
package container

import (
	"context"
	"errors"
	"fmt"
)

// ConsumerConfig represents configuration for consumer container.
type ConsumerConfig struct {
	// Future consumer configurations will go here
	// For example:
	// Kafka    KafkaConfig    `yaml:"kafka"`
	// Redis    RedisConfig    `yaml:"redis"`
	// RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

// ConsumerContainer holds consumer-specific dependencies.
type ConsumerContainer struct {
	*BaseContainer
	// Future consumer services will go here
	// For example:
	// KafkaConsumer   kafka.Consumer
	// RedisSubscriber redis.Subscriber
	// MessageHandler  handler.MessageHandler
}

// Close gracefully shuts down the consumer container and all its services.
func (c *ConsumerContainer) Close() error {
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

// Start starts the consumer container services.
func (c *ConsumerContainer) Start(ctx context.Context) error {
	// Future: Start consumer services
	// if c.KafkaConsumer != nil {
	//     go c.KafkaConsumer.Start(ctx)
	// }

	return nil
}

// Health checks the health of consumer services.
func (c *ConsumerContainer) Health(ctx context.Context) error {
	// Future: Health checks for consumer services
	// if c.KafkaConsumer != nil {
	//     if err := c.KafkaConsumer.Ping(ctx); err != nil {
	//         return fmt.Errorf("kafka consumer health check failed: %w", err)
	//     }
	// }

	return nil
}
