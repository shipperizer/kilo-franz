package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/subscriber"
)

const (
	defaultKafkaURL   = "kafka:9092"
	defaultTopic      = "test"
	defaultGroupID    = "example-consumer.cgroup"
	defaultWorkers    = 3
	defaultRefreshMin = 5
)

// DummyService implements subscriber.ServiceInterface
type DummyService struct {
	logger logging.LoggerInterface
}

// TaskName returns the name of the task
func (s *DummyService) TaskName() string {
	return "dummy-consumer"
}

// Flow processes a message
func (s *DummyService) Flow(key, value []byte) error {
	s.logger.Infof("received message: key=%s value=%s", string(key), string(value))
	return nil
}

func main() {
	logger := logging.NewLogger()
	defer logger.Desugar().Sync()

	kafkaURL := envOrDefault("KAFKA_URL", defaultKafkaURL)
	topic := envOrDefault("KAFKA_TOPIC", defaultTopic)
	groupID := envOrDefault("KAFKA_GROUP_ID", defaultGroupID)

	cfg := config.NewConfig(
		time.Duration(defaultRefreshMin)*time.Minute,
		nil,
		nil,
		logger,
	)

	readerCfg := config.NewReaderConfig(
		cfg,
		strings.Split(kafkaURL, ","),
		topic,
		groupID,
		defaultWorkers,
		30*time.Second,
	)

	reader := core.NewReader(readerCfg)

	monitor := NewMonitor("example-consumer", logger)

	svc := &DummyService{logger: logger}

	consumer, err := subscriber.NewChannelConsumer(reader, svc, monitor)
	if err != nil {
		logger.Fatalf("failed to create consumer: %v", err)
	}

	logger.Infof("starting consumer on topic=%s group=%s", topic, groupID)
	consumer.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	logger.Info("shutting down consumer")
	consumer.Stop()
	logger.Info("consumer stopped")
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

