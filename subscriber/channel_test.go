package subscriber

import (
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	uuid "github.com/satori/go.uuid"
	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/publisher"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -package subscriber -destination ./mock_monitor.go -source=../monitoring/interfaces.go MonitorInterface
//go:generate mockgen -package subscriber -destination ./mock_service.go . ServiceInterface

func TestNewChannelConsumerReturnsInterfaceImplementation(t *testing.T) {
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

	groupID := fmt.Sprintf("test.%s", uuid.NewV4().String())
	cfg := config.NewConfig(1*time.Hour, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 1, 15*time.Second)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	mockSvc.EXPECT().Flow(gomock.Any(), gomock.Any()).AnyTimes().DoAndReturn(
		func(MessageKey, MessageValue []byte) error {
			fmt.Println("running flow")
			executed = true
			return nil
		},
	)
	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().Incr(gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	assert := assert.New(t)

	c, err := NewChannelConsumer(reader, mockSvc, mockMonitor)
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
	time.Sleep(time.Duration(specs.WaitTime) * time.Millisecond)
	defer c.Stop()

	assert.True(executed, "Service Flow function should have been executed")
}

func TestNewChannelConsumerNotBlockingRefreshIfNoMessages(t *testing.T) {
	type EnvSpec struct {
		Brokers          []string `envconfig:"kafka_cnx_string"`
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
		WaitTime         int      `envconfig:"wait_time_ms" default:"1000"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	refresh := false

	groupID := fmt.Sprintf("test.%s", uuid.NewV4().String())
	cfg := config.NewConfig(5*time.Millisecond, nil, config.NewLogger())
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 10*time.Millisecond)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockSvc := NewMockServiceInterface(ctrl)

	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().GetService().AnyTimes()
	mockMonitor.EXPECT().Gauge(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockSvc.EXPECT().TaskName().AnyTimes().Return("test")
	mockMonitor.EXPECT().Incr(gomock.Any(), gomock.Any()).AnyTimes().MinTimes(2).Do(
		func(metric string, l map[string]string) {
			switch metric {
			case "labs_stream_errors":
				if task, ok := l["task"]; ok {
					if task == "tls-refresh" {
						refresh = true
					}
				}
			case "labs_stream_refresh_subscriber_v1":
				refresh = true
			}
		},
	)

	assert := assert.New(t)

	c, err := NewChannelConsumer(reader, mockSvc, mockMonitor)
	assert.Nil(err, "No error should be thrown")

	c.Start()

	time.Sleep(time.Duration(specs.WaitTime) * time.Millisecond)
	time.Sleep(5000 * time.Millisecond)
	defer c.Stop()

	assert.True(refresh, "autorefresh should have run")
}
