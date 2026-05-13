package publisher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.uber.org/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/internal/testutil"
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
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := "test-multi-config"

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writers := make([]core.RefreshableInterface, 0)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	for i := 0; i < 3; i++ {
		writerCfg := config.NewWriterConfig(cfg, brokers, topic, fmt.Sprintf("test-topic-%v", i), true, &DummyEncoder{})
		writers = append(writers, core.NewWriter(writerCfg))
	}

	producer := NewProducer(mockMonitor, writers...)

	assert := assert.New(t)
	assert.IsType(&Producer{}, producer, "A pointer to Producer should be returned")
	assert.Equal(map[string]string{"test-topic-0": topic, "test-topic-1": topic, "test-topic-2": topic}, producer.ListTopics(), "3 writers should have been initialized")
}

func TestProducerPublishSucceeds(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := "test-publish"
	testutil.CreateTopic(t, ctx, brokers, topic, 1)

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, topic, "test-topic", false, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("test-topic", NewMessage("test-1", Dummy{Value: "test"}), NewMessage("test-2", Dummy{Value: "test"}), NewMessage("test-3", Dummy{Value: "test"}), NewMessage("", Dummy{Value: "test"}))

	assert := assert.New(t)
	assert.Nil(err, "No error should be returned")
}

func TestProducerCloseMultipleTimesNotPanics(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := "test-close"

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
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

func TestProducerPublishFailsIfNoNicknameTopicFound(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := "test-no-nick"

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, topic, "", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	err := producer.Publish("fake-topic", NewMessage("test", Dummy{Value: "test"}))

	assert := assert.New(t)
	assert.NotNil(err, "An error should have been thrown, topic was not supposed to be found")
}

func TestProducerStats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := "test-stats"

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, brokers, topic, "test-topic", true, nil)
	writer := core.NewWriter(writerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	producer := NewProducer(mockMonitor, writer)
	stats := producer.Stats("test-topic")

	assert := assert.New(t)
	assert.NotNil(stats)
}
