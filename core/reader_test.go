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

//go:generate mockgen -build_flags=--mod=mod -package core -destination ./mock_logging.go -source=../logging/interfaces.go

func TestNewReaderReturnsReader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, "test-topic", "my-test", 3, 1*time.Minute)

	r := NewReader(readerCfg)
	defer r.Close()

	assert := assert.New(t)
	assert.IsType(&Reader{}, r, "A pointer to Reader should be returned")

	rr, err := r.Get(context.TODO())
	assert.Nil(err, "No error should have been thrown")
	assert.IsType(&kafka.Reader{}, rr, "A pointer to kafka.Reader should be returned")
}

func TestReaderRenewsKafkaReader(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, "test-topic", "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	rr1, _ := r1.Get(context.TODO())
	r1.Renew(nil, nil)
	rr2, _ := r1.Get(context.TODO())
	assert.NotEqual(t, rr1, rr2, "kafka.Reader objects should have been renewed")
}

func TestReaderStatsAreReaderStats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, "test-topic", "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	stats := r1.Stats()
	assert.IsType(t, kafka.ReaderStats{}, stats.(kafka.ReaderStats), "Castable ReaderStats object should have been returned")
}

func TestReaderConfigIsReaderConfigInterface(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, "test-topic", "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	c1 := r1.Config()
	assert.IsType(t, &config.ReaderConfig{}, c1.(ReaderConfigInterface), "Castable ReaderConfig object should have been returned")
}
