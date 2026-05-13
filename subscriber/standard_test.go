package subscriber

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/internal/testutil"
	"github.com/shipperizer/kilo-franz/publisher"
)

//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_service.go . ServiceInterface
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_core.go -source=../core/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_monitor.go -source=../monitoring/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_config.go -source=../config/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_refresh.go -source=../refresh/interfaces.go

type DummyStandard struct {
	Value string `json:"value"`
}

func TestNewStandardConsumerReturnsInterfaceImplementation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := fmt.Sprintf("test-std-%s", uuid.New().String()[:8])
	testutil.CreateTopic(t, ctx, brokers, topic, 1)

	var executed atomic.Bool

	groupID := fmt.Sprintf("test.%s", uuid.New().String())
	cfg := config.NewConfig(5*time.Minute, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, topic, groupID, 3, 5*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetMetric(gomock.Any()).AnyTimes().Return(mockMetric, nil)
	mockMonitor.EXPECT().GetService().AnyTimes()

	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Observe(gomock.Any(), gomock.Any()).AnyTimes()

	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(MessageKey, MessageValue []byte) error {
			executed.Store(true)
			return nil
		},
	)

	c, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	require.NoError(t, err)

	c.Start()

	// Wait for consumer group to stabilize
	time.Sleep(10 * time.Second)

	writerCfg := config.NewWriterConfig(cfg, brokers, topic, "test", false, nil)
	writer := core.NewWriter(writerCfg)

	pub := publisher.NewProducer(mockMonitor, writer)
	for g := 0; g < 5; g++ {
		msgs := make([]publisher.MessageInterface, 0)
		for i := 0; i < 5; i++ {
			msgs = append(msgs, publisher.NewMessage("test", DummyStandard{Value: "test"}))
		}
		err = pub.Publish("test", msgs...)
		require.NoError(t, err)
	}

	pub.Close()

	assert.Eventually(t, func() bool {
		return executed.Load()
	}, 60*time.Second, 500*time.Millisecond, "Service Flow function should have been executed")
}

func TestStandardConsumerRefreshableConfigFail(t *testing.T) {
	ctrl := gomock.NewController(t)

	reader := NewMockRefreshableInterface(ctrl)
	reader.EXPECT().Config().Return(nil).Times(1)

	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	_, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	assert.NotNil(t, err)
}

func TestStandardConsumerInterfaceTypeFails(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRefresher := NewMockRefreshableInterface(ctrl)
	mockRefresher.EXPECT().Config().Return(struct{}{}).Times(1)

	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	_, err := NewStandardConsumer(mockRefresher, mockSvc, mockMonitor)
	assert.NotNil(t, err)
}

func TestStandardConsumerStatsPanics(t *testing.T) {
	ctrl := gomock.NewController(t)

	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Stats().Return(struct{}{}).Times(1)

	c := StandardConsumer{af: autoRefresh}

	assert.Panics(t, func() { c.Stats() })
}

func TestStandardConsumerUnwrapReaderFail1(t *testing.T) {
	ctrl := gomock.NewController(t)

	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Object(gomock.Any()).Return(&core.Reader{}, fmt.Errorf("boom")).Times(1)

	c := StandardConsumer{af: autoRefresh}

	_, err := c.unwrapReader()
	assert.NotNil(t, err)
}

func TestStandardConsumerUnwrapReaderFail2(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockRefresher := NewMockRefreshableInterface(ctrl)
	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Object(gomock.Any()).Return(mockRefresher, nil).Times(1)

	c := StandardConsumer{af: autoRefresh}

	_, err := c.unwrapReader()
	assert.NotNil(t, err)
}

func TestNewStandardConsumerNotBlockingRefreshIfNoMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := fmt.Sprintf("test-std-refresh-%s", uuid.New().String()[:8])
	testutil.CreateTopic(t, ctx, brokers, topic, 1)

	var refreshed atomic.Bool

	groupID := fmt.Sprintf("test.%s", uuid.New().String())
	cfg := config.NewConfig(5*time.Millisecond, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, topic, groupID, 3, 10*time.Millisecond)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)
	mockMetricRefresh := NewMockMetricInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetService().AnyTimes()
	mockMonitor.EXPECT().GetMetric("labs_stream_refresh_subscriber_v1").AnyTimes().Return(mockMetricRefresh, nil)
	mockMonitor.EXPECT().GetMetric(gomock.Any()).AnyTimes().Return(mockMetric, nil)
	mockMonitor.EXPECT().GetService().AnyTimes()
	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Observe(gomock.Any(), gomock.Any()).AnyTimes()

	mockMetricRefresh.EXPECT().Inc(gomock.Any()).AnyTimes().Do(
		func(tags map[string]string) {
			refreshed.Store(true)
		},
	)

	c, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	require.NoError(t, err)

	c.Start()

	assert.Eventually(t, func() bool {
		return refreshed.Load()
	}, 30*time.Second, 500*time.Millisecond, "autorefresh should have run")

	c.Stop()
}

func TestNewStandardConsumerWithKafkaMessages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	brokers := testutil.KafkaContainer(t, ctx)
	topic := fmt.Sprintf("test-std-msgs-%s", uuid.New().String()[:8])
	testutil.CreateTopic(t, ctx, brokers, topic, 1)

	var messageCount atomic.Int32

	groupID := fmt.Sprintf("test.%s", uuid.New().String())
	// Use 5-minute refresh to prevent AutoRefreshX from disrupting the consumer group
	cfg := config.NewConfig(5*time.Minute, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, brokers, topic, groupID, 1, 5*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetMetric(gomock.Any()).AnyTimes().Return(mockMetric, nil)
	mockMonitor.EXPECT().GetService().AnyTimes()
	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Observe(gomock.Any(), gomock.Any()).AnyTimes()

	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(key, value []byte) error {
			messageCount.Add(1)
			return nil
		},
	)

	c, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	require.NoError(t, err)

	c.Start()

	// Wait for consumer group to join and get partition assignment
	time.Sleep(10 * time.Second)

	// Produce messages directly via kafka-go
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
	}

	msgs := make([]kafka.Message, 10)
	for i := range msgs {
		msgs[i] = kafka.Message{
			Key:   []byte(fmt.Sprintf("key-%d", i)),
			Value: []byte(fmt.Sprintf("value-%d", i)),
		}
	}

	err = writer.WriteMessages(ctx, msgs...)
	require.NoError(t, err)
	writer.Close()

	assert.Eventually(t, func() bool {
		return messageCount.Load() >= 10
	}, 60*time.Second, 500*time.Millisecond, "all messages should be consumed")

	// Note: Stop() will block until the current ReadMessage timeout expires.
	// Since we use a 5-minute timeout, we don't call Stop() here and let
	// the container cleanup handle teardown.
}
