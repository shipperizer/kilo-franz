package publisher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/sasl"
	"github.com/shipperizer/kilo-franz/vault"
)

//go:generate mockgen -build_flags=--mod=mod -package publisher -destination ./mock_logging.go -source=../logging/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package publisher -destination ./mock_monitor.go -source=../monitoring/interfaces.go

type DummyEncoder struct{}

func (e *DummyEncoder) Encode(msg interface{}) ([]byte, error) {
	return []byte("dummy"), nil
}

type Dummy struct {
	Value string `json:"value"`
}

func TestNewProducerWithMultipleConfigs(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writers := make([]core.RefreshableInterface, 0)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	for i := 0; i < 3; i++ {
		writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, fmt.Sprintf("test-topic-%v", i), true, &DummyEncoder{})
		writers = append(writers, core.NewWriter(writerCfg))

	}

	producer := NewProducer(mockMonitor, writers...)

	assert := assert.New(t)

	assert.IsType(&Producer{}, producer, "A pointer to Producer should be returned")
	assert.Equal(map[string]string{"test-topic-0": specs.Topic, "test-topic-1": specs.Topic, "test-topic-2": specs.Topic}, producer.ListTopics(), "3 writers should have been initialized")
}

func TestProducerPublishSucceeds(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("test-topic", NewMessage("test-1", Dummy{Value: "test"}), NewMessage("test-2", Dummy{Value: "test"}), NewMessage("test-3", Dummy{Value: "test"}), NewMessage("", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.Nil(err, "No error should be returned")
}

func TestProducerSASLPublishSucceeds(t *testing.T) {
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
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	mockLogger := NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any(), gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := sasl.NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := sasl.NewSASLConfig("", "", sasl.PlainSASL, nil, cfgSM, mockLogger)

	writerCfg := config.NewWriterConfig(
		config.NewConfig(1*time.Hour, nil, cfgSASL, mockLogger),
		specs.Brokers,
		specs.Topic,
		"test-topic",
		false,
		nil,
	)

	writer := core.NewWriter(writerCfg)

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("test-topic", NewMessage("test-1", Dummy{Value: "test"}), NewMessage("test-2", Dummy{Value: "test"}), NewMessage("test-3", Dummy{Value: "test"}), NewMessage("", Dummy{Value: "test"}))

	assert := assert.New(t)
	assert.Nil(err, "No error should be returned")
}

func TestProducerCloseMultipleTimesNotPanics(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	producer.Close()

	assert := assert.New(t)
	assert.NotPanics(
		func() {
			for i := 0; i < 10; i++ {
				producer.Close()
			}
		},
	)
}

func TestProducerPublishPanicsIfNoNicknameTopicFound(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("fake-topic", NewMessage("test", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.NotNil(err, "An error should have been thrown, topic was not supposed to be found")
}

func TestProducerPublishFailsWithKafkaError(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, "random", "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)

	err := producer.Publish("test-topic", NewMessage("test", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.NotNil(err, "Topic should have not existed and kafka should have returned an error")
}

func TestProducerPublishFailsWithKafkaStats(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, "random", "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)

	stats := producer.Stats("test-topic")

	assert := assert.New(t)

	assert.NotNil(stats)
}
