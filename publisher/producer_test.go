package publisher

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
)

//go:generate mockgen -package publisher -destination ./mock_monitor.go -source=../monitoring/interfaces.go MonitorInterface

type EnvSpec struct {
	Brokers []string `envconfig:"kafka_cnx_string"`
	Topic   string   `envconfig:"kafka_topic"`
}

type DummyEncoder struct{}

func (e *DummyEncoder) Encode(msg interface{}) ([]byte, error) {
	return []byte("dummy"), nil
}

type Dummy struct {
	Value string `json:"value"`
}

func TestNewProducerWithMultipleConfigs(t *testing.T) {
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writers := make([]config.RefreshableInterface, 0)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

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
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("test-topic", NewMessage("test-1", Dummy{Value: "test"}), NewMessage("test-2", Dummy{Value: "test"}), NewMessage("test-3", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.Nil(err, "No error should be returned")
}

func TestProducerCloseMultipleTimesNotPanics(t *testing.T) {
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

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
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("fake-topic", NewMessage("test", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.NotNil(err, "An error should have been thrown, topic was not supposed to be found")
}

func TestProducerPublishFailsWithKafkaError(t *testing.T) {
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, "random", "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	producer := NewProducer(mockMonitor, writer)

	err := producer.Publish("test-topic", NewMessage("test", Dummy{Value: "test"}))

	assert := assert.New(t)

	assert.NotNil(err, "Topic should have not existed and kafka should have returned an error")
}
