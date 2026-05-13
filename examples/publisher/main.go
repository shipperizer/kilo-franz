package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/encoder"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/publisher"
)

const (
	defaultKafkaURL   = "kafka:9092"
	defaultTopic      = "test"
	defaultNickname   = "test"
	defaultInterval   = 5
	defaultRefreshMin = 5
)

func main() {
	logger := logging.NewLogger()
	defer logger.Desugar().Sync()

	kafkaURL := envOrDefault("KAFKA_URL", defaultKafkaURL)
	topic := envOrDefault("KAFKA_TOPIC", defaultTopic)
	nickname := envOrDefault("KAFKA_TOPIC_NICKNAME", defaultNickname)

	cfg := config.NewConfig(
		time.Duration(defaultRefreshMin)*time.Minute,
		nil,
		nil,
		logger,
	)

	writerCfg := config.NewWriterConfig(
		cfg,
		strings.Split(kafkaURL, ","),
		topic,
		nickname,
		false,
		encoder.NewJSONEncoder(),
	)

	writer := core.NewWriter(writerCfg)

	monitor := NewMonitor("example-publisher", logger)

	producer := publisher.NewProducer(monitor, writer)
	defer producer.Close()

	ticker := time.NewTicker(time.Duration(defaultInterval) * time.Second)
	defer ticker.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	logger.Infof("starting publisher on topic=%s nickname=%s interval=%ds", topic, nickname, defaultInterval)

	counter := 0
	for {
		select {
		case <-ticker.C:
			counter++
			msg := publisher.NewMessage("", map[string]interface{}{
				"counter":   counter,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"source":    "example-publisher",
			})

			if err := producer.Publish(nickname, msg); err != nil {
				logger.Errorf("failed to publish message: %v", err)
				continue
			}

			logger.Infof("published message #%d to topic=%s", counter, topic)

		case <-c:
			logger.Info("shutting down publisher")
			return
		}
	}
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
