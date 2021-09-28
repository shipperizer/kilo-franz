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

func TestNewWriterReturnsKafkaWriter(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w := NewWriter(writerCfg)

	assert.IsType(&Writer{}, w, "A pointer to Writer should be returned")

	ww, err := w.Get(context.TODO())
	assert.Nil(err, "No error should have been thrown")
	assert.IsType(&kafka.Writer{}, ww, "A pointer to kafka.Writer should be returned")
}

func TestWriterRenewsKafkaWriter(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	assert.NotEqual(w1.Renew(nil), w1, "Writer object should have been renewed")

	ww1, _ := w1.Get(context.TODO())
	ww2, _ := w1.Renew(nil).Get(context.TODO())
	assert.NotEqual(ww1, ww2, "kafka.Writer objects should have been renewed")
}
func TestWriterStatsAreWriterStats(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	stats := w1.Stats()
	assert.IsType(kafka.WriterStats{}, stats.(kafka.WriterStats), "Castable WriterStats object should have been returned")

}
func TestWriterConfigIsWriterConfigInterface(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	c1 := w1.Config()
	assert.IsType(&config.WriterConfig{}, c1.(config.WriterConfigInterface), "Castable WriterConfig object should have been returned")
}
