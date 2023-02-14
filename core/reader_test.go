package core

import (
	"context"
	"testing"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
)

func TestNewReaderReturnsReader(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, "my-test", 3, 1*time.Minute)

	r := NewReader(readerCfg)
	defer r.Close()

	assert.IsType(&Reader{}, r, "A pointer to Reader should be returned")

	rr, err := r.Get(context.TODO())
	assert.Nil(err, "No error should have been thrown")
	assert.IsType(&kafka.Reader{}, rr, "A pointer to kafka.Reader should be returned")
}

func TestReaderRenewsKafkaReader(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	rr1, _ := r1.Get(context.TODO())
	r1.Renew(nil, nil)
	rr2, _ := r1.Get(context.TODO())
	assert.NotEqual(rr1, rr2, "kafka.Reader objects should have been renewed")
}
func TestReaderStatsAreReaderStats(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	stats := r1.Stats()
	assert.IsType(kafka.ReaderStats{}, stats.(kafka.ReaderStats), "Castable ReaderStats object should have been returned")

}

func TestReaderConfigIsWriterConfigInterface(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, "my-test", 3, 1*time.Minute)

	r1 := NewReader(readerCfg)

	c1 := r1.Config()
	assert.IsType(&config.ReaderConfig{}, c1.(ReaderConfigInterface), "Castable ReaderConfig object should have been returned")
}
