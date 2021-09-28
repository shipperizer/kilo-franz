// following the same pattern in of jwk lib https://github.com/lestrrat-go/jwx/blob/main/jwk/refresh.go
package refresh

//go:generate mockgen -package refresh -destination ./mock_monitor.go -source=../monitoring/interfaces.go MonitorInterface

import (
	"context"
	"sync"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/tls"
)

func TestNewAutoRefreshXReturnsInterfaceImplementation(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.1"

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	defer af.Stop()
	_, err := af.Object(context.TODO())

	assert := assert.New(t)

	assert.IsType(&AutoRefreshX{}, af, "A pointer to AutoRefreshX should be returned")
	assert.Nil(err, "No error should be thrown")
}

func TestNewAutoRefreshXReturnsNewReaderOnRefresh(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.2"

	cfg := config.NewConfig(1*time.Hour, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	defer af.Stop()
	r1, err := af.Object(context.TODO())
	assert := assert.New(t)
	assert.Equal(reader, r1, "Refreshable pointers should be different after Refresh")

	assert.Nil(err, "No error should be thrown")

	r2, err := af.Refresh(context.TODO())
	assert.Nil(err, "No error should be thrown")
	assert.NotEqual(r2, r1, "Refreshable pointers should be different after Refresh")
}

func TestNewAutoRefreshXReaderWithEmptyReaderPanics(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.2"

	cfg := config.NewConfig(1*time.Hour, &tls.TLSConfig{UseTLS: false}, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	assert := assert.New(t)
	assert.Panics(func() {
		NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), nil, nil, mockMonitor)
	})
}

func TestNewAutoRefreshXRunsInBackgroundWithReader(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.3"

	cfg := config.NewConfig(5*time.Millisecond, &tls.TLSConfig{UseTLS: false}, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	refresh := false
	labels := make(map[string]string)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	mockMonitor.EXPECT().GetService().AnyTimes().Return(groupID)
	mockMonitor.EXPECT().Incr(gomock.Any(), gomock.Any()).AnyTimes().MinTimes(2).Do(
		func(metric string, l map[string]string) {
			switch metric {
			case "errors":
				if task, ok := l["task"]; ok {
					if task == "tls-refresh" {
						refresh = true
						labels = l
					}
				}
			case "refresh_subscriber_v1":
				refresh = true
				labels = l
			}
		},
	)

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	time.Sleep(50 * time.Millisecond)
	af.Stop()

	assert := assert.New(t)

	assert.True(refresh, "autorefresh should have run")
	assert.Equal(map[string]string{"app": groupID, "topic": specs.Topic, "consumer_group": groupID}, labels, "Labels for metrics should match")
}

func TestNewAutoRefreshXRunsInBackgroundWithWriter(t *testing.T) {
	type EnvSpec struct {
		Brokers []string `envconfig:"kafka_cnx_string"`
		Topic   string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.3"

	cfg := config.NewConfig(5*time.Millisecond, &tls.TLSConfig{UseTLS: false}, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "", true, nil)
	writer := core.NewWriter(writerCfg)

	refresh := false
	labels := make(map[string]string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	mockMonitor.EXPECT().GetService().AnyTimes().Return(groupID)
	mockMonitor.EXPECT().Incr(gomock.Any(), gomock.Any()).AnyTimes().MinTimes(2).Do(
		func(metric string, l map[string]string) {
			switch metric {
			case "errors":
				if task, ok := l["task"]; ok {
					if task == "tls-refresh" {
						refresh = true
						labels = l
					}
				}
			case "refresh_publisher_v1":
				refresh = true
				labels = l
			}
		},
	)
	mockMonitor.EXPECT().Incr(gomock.Any(), gomock.Any()).AnyTimes()

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, writerCfg), writer, nil, mockMonitor)
	time.Sleep(50 * time.Millisecond)
	af.Stop()

	assert := assert.New(t)

	assert.True(refresh, "autorefresh should have run")
	assert.Equal(map[string]string{"app": groupID, "topic": specs.Topic}, labels, "Labels for metrics should match")
}

func TestNewAutoRefreshXConfigureNotPanics(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.4"
	cfg := config.NewConfig(1*time.Hour, &tls.TLSConfig{UseTLS: false}, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	defer af.Stop()
	assert := assert.New(t)

	assert.NotPanics(func() { af.Configure(context.TODO(), &tls.TLSConfig{UseTLS: true}) })

}

func TestNewAutoRefreshXStopMultipleTimesNotPanics(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var mutex sync.RWMutex
	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	groupID := "test.4"
	cfg := config.NewConfig(1*time.Hour, &tls.TLSConfig{UseTLS: false}, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	af.Stop()

	assert := assert.New(t)

	assert.NotPanics(
		func() {
			for i := 0; i < 10; i++ {
				af.Stop()
			}
		},
	)
}
