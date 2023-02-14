// following the same pattern in of jwk lib https://github.com/lestrrat-go/jwx/blob/main/jwk/refresh.go
package refresh

//go:generate mockgen -build_flags=--mod=mod -package refresh -destination ./mock_refresh.go -source=interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package refresh -destination ./mock_logging.go -source=../logging/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package refresh -destination ./mock_core.go -source=../core/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package refresh -destination ./mock_config.go -source=../config/interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package refresh -destination ./mock_monitor.go -source=../monitoring/interfaces.go

import (
	"context"
	"sync"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stretchr/testify/assert"

	"testing"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/sasl"
	_tls "github.com/shipperizer/kilo-franz/tls"
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

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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

	cfg := config.NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	defer af.Stop()
	r1, err := af.Object(context.TODO())

	rr0, _ := reader.Get(context.TODO())
	rr1, _ := r1.Get(context.TODO())
	assert := assert.New(t)

	assert.Equal(rr0, rr1, "Refreshable pointers should be different ONLY after Refresh")
	assert.Nil(err, "No error should be thrown")

	r2, err := af.Refresh(context.TODO())
	rr2, _ := r2.Get(context.TODO())
	assert.Nil(err, "No error should be thrown")
	assert.NotEqual(rr2, rr1, "Refreshable pointers should be different after Refresh")
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

	cfg := config.NewConfig(1*time.Hour, &_tls.TLSConfig{UseTLS: false}, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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

	cfg := config.NewConfig(5*time.Millisecond, &_tls.TLSConfig{UseTLS: false}, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	refresh := false
	labels := make(map[string]string)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	mockMetric := NewMockMetricInterface(ctrl)

	mockMonitor.EXPECT().GetService().AnyTimes().Return(groupID)
	mockMonitor.EXPECT().GetMetric("labs_stream_errors").AnyTimes().Return(mockMetric, nil)
	mockMonitor.EXPECT().GetMetric("labs_stream_refresh_subscriber_v1").AnyTimes().Return(mockMetric, nil)

	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes().Do(
		func(tags map[string]string) {
			refresh = true
			labels = tags
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

	cfg := config.NewConfig(5*time.Millisecond, &_tls.TLSConfig{UseTLS: false}, nil, nil)
	writerCfg := config.NewWriterConfig(cfg, specs.Brokers, specs.Topic, "", true, nil)
	writer := core.NewWriter(writerCfg)

	refresh := false
	labels := make(map[string]string)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMetric := NewMockMetricInterface(ctrl)

	mockMonitor.EXPECT().GetService().AnyTimes().Return(groupID)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockMonitor.EXPECT().GetMetric("labs_stream_errors").AnyTimes().Return(mockMetric, nil)
	mockMonitor.EXPECT().GetMetric("labs_stream_refresh_publisher_v1").AnyTimes().Return(mockMetric, nil)

	mockMetric.EXPECT().Inc(gomock.Any()).AnyTimes().Do(
		func(tags map[string]string) {
			refresh = true
			labels = tags
		},
	)

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
	cfg := config.NewConfig(1*time.Hour, &_tls.TLSConfig{UseTLS: false}, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	af := NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&mutex, readerCfg), reader, nil, mockMonitor)
	defer af.Stop()
	assert := assert.New(t)

	assert.NotPanics(func() {
		af.Configure(context.TODO(), &_tls.TLSConfig{UseTLS: true}, sasl.NewSASLConfig("test", "test", sasl.PlainSASL, scram.SHA256, nil, mockLogger))
	})

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
	cfg := config.NewConfig(1*time.Hour, &_tls.TLSConfig{UseTLS: false}, nil, nil)
	readerCfg := config.NewReaderConfig(cfg, specs.BootstrapServers, specs.Topic, groupID, 3, 1*time.Minute)
	reader := core.NewReader(readerCfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

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

func TestNewAutoRefreshXStats(t *testing.T) {
	type EnvSpec struct {
		BootstrapServers []string `envconfig:"kafka_cnx_string"`
		Topic            string   `envconfig:"kafka_topic"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	readerCfg := NewMockReaderConfigInterface(ctrl)
	readerCfg.EXPECT().GetBootstrapServers().Return([]string{"server"}).Times(1)
	readerCfg.EXPECT().GetGroupID().Return("GroupID").Times(1)
	readerCfg.EXPECT().GetTopic().Return("topic").Times(1)
	readerCfg.EXPECT().GetDialer().Return(&kafka.Dialer{}).Times(1)

	reader := core.NewReader(readerCfg)

	mockMonitor := NewMockMonitorInterface(ctrl)
	mockMonitor.EXPECT().AddMetrics(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	autoRefreshXConfig := NewMockAutoRefreshXConfigInterface(ctrl)
	autoRefreshXConfig.EXPECT().GetMutexObj().Return(&sync.RWMutex{}).Times(1)
	autoRefreshXConfig.EXPECT().GetRefreshTimeout().Return(time.Second * 1).Times(1)
	autoRefreshXConfig.EXPECT().GetTLSConfig().Return(&_tls.TLSConfig{}).Times(1)
	autoRefreshXConfig.EXPECT().GetSASLConfig().Return(nil).Times(1)

	af := NewAutoRefreshX(context.TODO(), autoRefreshXConfig, reader, nil, mockMonitor)

	stats := af.Stats()

	assert := assert.New(t)

	assert.NotNil(stats)
}
