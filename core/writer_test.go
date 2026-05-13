package core

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/internal/testutil"
)

func TestNewWriterReturnsKafkaWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, "test-topic", "my-test", true, nil)

	w := NewWriter(writerCfg)

	assert := assert.New(t)
	assert.IsType(&Writer{}, w, "A pointer to Writer should be returned")

	ww, err := w.Get(context.TODO())
	assert.Nil(err, "No error should have been thrown")
	assert.IsType(&kafka.Writer{}, ww, "A pointer to kafka.Writer should be returned")
}

func TestWriterRenewsKafkaWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, "test-topic", "my-test", true, nil)

	w1 := NewWriter(writerCfg)
	defer w1.Close()

	ww1, _ := w1.Get(context.TODO())
	w1.Renew(nil, nil)
	ww2, _ := w1.Get(context.TODO())
	assert.NotEqual(t, ww1, ww2, "kafka.Writer objects should have been renewed")
}

func TestWriterStatsAreWriterStats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, "test-topic", "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	stats := w1.Stats()
	assert.IsType(t, kafka.WriterStats{}, stats.(kafka.WriterStats), "Castable WriterStats object should have been returned")
}

func TestWriterConfigIsWriterConfigInterface(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, "test-topic", "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	c1 := w1.Config()
	assert.IsType(t, &config.WriterConfig{}, c1.(WriterConfigInterface), "Castable WriterConfig object should have been returned")
}
