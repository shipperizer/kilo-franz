// Package testutil provides shared test utilities for integration tests.
package testutil

import (
	"context"
	"fmt"
	"net"
	"testing"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
)

// KafkaContainer starts a Kafka container and returns the broker address.
// It automatically terminates the container when the test completes.
func KafkaContainer(t *testing.T, ctx context.Context) []string {
	t.Helper()

	kafkaContainer, err := kafka.Run(ctx, "confluentinc/confluent-local:7.6.0")
	if err != nil {
		t.Fatalf("failed to start kafka container: %v", err)
	}

	t.Cleanup(func() {
		if err := kafkaContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate kafka container: %v", err)
		}
	})

	brokers, err := kafkaContainer.Brokers(ctx)
	if err != nil {
		t.Fatalf("failed to get kafka brokers: %v", err)
	}

	return brokers
}

// CreateTopic creates a Kafka topic with the given number of partitions.
func CreateTopic(t *testing.T, ctx context.Context, brokers []string, topic string, partitions int) {
	t.Helper()

	conn, err := kafkago.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		t.Fatalf("failed to dial kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		t.Fatalf("failed to get controller: %v", err)
	}

	controllerConn, err := kafkago.Dial("tcp", net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port)))
	if err != nil {
		// Fall back to using the broker address directly
		controllerConn, err = kafkago.Dial("tcp", brokers[0])
		if err != nil {
			t.Fatalf("failed to dial controller: %v", err)
		}
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
	if err != nil {
		t.Fatalf("failed to create topic %s: %v", topic, err)
	}
}
