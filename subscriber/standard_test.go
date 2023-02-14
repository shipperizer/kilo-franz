package subscriber

import (
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/publisher"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_service.go . ServiceInterface
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_core.go -source=../core/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_monitor.go -source=../monitoring/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_config.go -source=../config/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package subscriber -destination ./mock_refresh.go -source=../refresh/interfaces.go

func TestNewStandardConsumerReturnsInterfaceImplementation(t *testing.T) {
	type EnvSpec struct {
		Brokers          []string `envconfig:"kafka_cnx_string"`
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
		WaitTime         int      `envconfig:"wait_time_ms" default:"1000"`
	}

	type Dummy struct {
		Value string `json:"value"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	executed := false

	groupID := fmt.Sprintf("test.%s", uuid.New().String())
	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 15*time.Second)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetMetric(gomock.Any()).AnyTimes().Return(mockMetric, nil)

	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Observe(gomock.Any(), gomock.Any()).AnyTimes()

	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(MessageKey, MessageValue []byte) error {
			executed = true
			return nil
		},
	)

	assert := assert.New(t)

	c, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	assert.Nil(err, "No error should be thrown")

	c.Start()

	time.Sleep(time.Duration(specs.WaitTime) * time.Millisecond)

	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "test", false, nil)
	writer := core.NewWriter(writerCfg)

	pub := publisher.NewProducer(mockMonitor, writer)
	for g := 0; g < 5; g++ {
		msgs := make([]publisher.MessageInterface, 0)
		for i := 0; i < 5; i++ {
			msgs = append(msgs, publisher.NewMessage("test", Dummy{Value: "test"}))

		}
		err = pub.Publish("test", msgs...)
		assert.Nil(err, "No error should be thrown")
	}

	pub.Close()
	c.Stop()

	time.Sleep(time.Duration(specs.WaitTime) * time.Millisecond)

	assert.True(executed, "Service Flow function should have been executed")
}

func TestNewStandardConsumerNotBlockingRefreshIfNoMessages(t *testing.T) {
	type EnvSpec struct {
		Brokers          []string `envconfig:"kafka_cnx_string"`
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
		WaitTime         int      `envconfig:"wait_time_ms" default:"1000"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	refresh := false

	groupID := fmt.Sprintf("test.%s", uuid.New().String())
	cfg := config.NewConfig(5*time.Millisecond, nil, nil, logging.NewLogger())
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 10*time.Millisecond)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)
	mockMetricRefresh := NewMockMetricInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	// autorefresher
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetService().AnyTimes()
	mockMonitor.EXPECT().GetMetric("labs_stream_refresh_subscriber_v1").AnyTimes().Return(mockMetricRefresh, nil)
	mockMonitor.EXPECT().GetMetric(gomock.Any()).AnyTimes().Return(mockMetric, nil)
	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Set(gomock.Any(), gomock.Any()).AnyTimes()
	mockMetric.EXPECT().Observe(gomock.Any(), gomock.Any()).AnyTimes()

	mockMetricRefresh.EXPECT().Inc(gomock.Any()).AnyTimes().Do(
		func(tags map[string]string) {
			refresh = true
		},
	)

	assert := assert.New(t)

	c, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	assert.Nil(err, "No error should be thrown")

	c.Start()

	time.Sleep(time.Duration(specs.WaitTime) * time.Millisecond)

	c.Stop()
	time.Sleep(5000 * time.Millisecond)

	assert.True(refresh, "autorefresh should have run")
}

func TestStandardConsumerRefreshableConfigFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockRefreshableInterface(ctrl)
	reader.EXPECT().Config().Return(nil).Times(1)

	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	assert := assert.New(t)

	_, err := NewStandardConsumer(reader, mockSvc, mockMonitor)
	assert.NotNil(err)
}

func TestStandardConsumerInterfaceTypeFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRefresher := NewMockRefreshableInterface(ctrl)
	mockRefresher.EXPECT().Config().Return(struct{}{}).Times(1)

	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	assert := assert.New(t)

	_, err := NewStandardConsumer(mockRefresher, mockSvc, mockMonitor)
	assert.NotNil(err)
}

func TestStandardConsumerStatsPanics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assert := assert.New(t)

	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Stats().Return(struct{}{}).Times(1)

	c := StandardConsumer{af: autoRefresh}

	assert.Panics(func() { c.Stats() })
}

func TestStandardConsumerUnwrapReaderFail1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assert := assert.New(t)

	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Object(gomock.Any()).Return(&core.Reader{}, fmt.Errorf("boom")).Times(1)

	c := StandardConsumer{af: autoRefresh}

	_, err := c.unwrapReader()

	assert.NotNil(err)
}

func TestStandardConsumerUnwrapReaderFail2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	assert := assert.New(t)

	mockRefresher := NewMockRefreshableInterface(ctrl)
	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Object(gomock.Any()).Return(mockRefresher, nil).Times(1)

	c := StandardConsumer{af: autoRefresh}

	_, err := c.unwrapReader()

	assert.NotNil(err)
}

func TestStandardConsumerUnwrapReaderFail3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	assert := assert.New(t)

	mockRefresher := NewMockRefreshableInterface(ctrl)
	autoRefresh := NewMockAutoRefreshXInterface(ctrl)
	autoRefresh.EXPECT().Object(gomock.Any()).Return(mockRefresher, nil).Times(1)

	c := StandardConsumer{af: autoRefresh}

	_, err := c.unwrapReader()

	assert.NotNil(err)
}
