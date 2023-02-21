package core

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	gomock "github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/sasl"
	"github.com/shipperizer/kilo-franz/vault"
)

//go:generate mockgen -build_flags=--mod=mod -package core -destination ./mock_logging.go -source=../logging/interfaces.go

func TestNewWriterReturnsKafkaWriter(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	err := envconfig.Process("", &specs)

	assert := assert.New(t)

	assert.Nilf(err, "Env vars are misconfigured: %v", specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
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

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w1 := NewWriter(writerCfg)
	defer w1.Close()

	ww1, _ := w1.Get(context.TODO())
	w1.Renew(nil, nil)
	ww2, _ := w1.Get(context.TODO())
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

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
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

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "my-test", true, nil)

	w1 := NewWriter(writerCfg)

	c1 := w1.Config()
	assert.IsType(&config.WriterConfig{}, c1.(WriterConfigInterface), "Castable WriterConfig object should have been returned")
}

func TestWriterWithSASL(t *testing.T) {
	type EnvSpec struct {
		Brokers  []string `envconfig:"kafka_sasl_cnx_string"`
		Topic    string   `envconfig:"kafka_topic"`
		Region   string   `envconfig:"aws_region"`
		Endpoint string   `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, options ...interface{}) (aws.Endpoint, error) {
		if service == secretsmanager.ServiceID && specs.Endpoint != "" {
			return aws.Endpoint{
				URL:           specs.Endpoint,
				SigningRegion: specs.Region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to it's default resolution.
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, _ := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithRegion(specs.Region),
	)

	client := secretsmanager.NewFromConfig(cfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := sasl.NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := sasl.NewSASLConfig("", "", true, sasl.PlainSASL, nil, cfgSM, mockLogger)

	assert := assert.New(t)

	writerCfg := config.NewWriterConfig(
		config.NewConfig(1*time.Hour, nil, cfgSASL, mockLogger),
		specs.Brokers,
		specs.Topic,
		"my-test",
		true,
		nil,
	)

	w1 := NewWriter(writerCfg)

	stats := w1.Stats()
	assert.IsType(kafka.WriterStats{}, stats.(kafka.WriterStats), "Castable WriterStats object should have been returned")
}
