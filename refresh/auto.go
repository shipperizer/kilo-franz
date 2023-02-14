// following the same pattern in of jwk lib https://github.com/lestrrat-go/jwx/blob/main/jwk/refresh.go
package refresh

import (
	"context"
	"sync"
	"time"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/monitoring"
)

// channelConfig is used to wrap both tls and sasl configs and send them down the channel
type channelConfig struct {
	config.TLSConfigInterface
	config.SASLConfigInterface
}

// AutoRefreshX is the object taking care of refreshing tls secrets from secretsmanager
// can handle both Writer and Reader structs, due to the core.RefreshableInterface
// it runs a ticker which when triggered will run the machinery to fetch a new set of secrets
// and refresh the refreshable object
// handles the fetching of this object via mutex
type AutoRefreshX struct {
	shutdownCh    chan bool
	refreshTicker *time.Ticker
	// config
	configCh   chan channelConfig
	tlsConfig  config.TLSConfigInterface
	saslConfig config.SASLConfigInterface
	// cached reader/writer object
	obj core.RefreshableInterface

	mutexFetching sync.Mutex
	mutexObj      *sync.RWMutex

	monitor monitoring.MonitorInterface

	logger logging.LoggerInterface
}

// simply a wrapper around fetch
func (af *AutoRefreshX) refresh(ctx context.Context) (time.Time, error) {
	err := af.fetch()

	return time.Now(), err
}

// fetch will ask the object to renew itself
func (af *AutoRefreshX) fetch() error {
	af.mutexFetching.Lock()
	defer af.mutexFetching.Unlock()

	af.mutexObj.Lock()
	defer af.mutexObj.Unlock()

	// close current reader
	af.obj.Close()
	// then refresh it with a new one
	af.obj.Renew(af.tlsConfig, af.saslConfig)

	// TODO @shipperizer evaluate if return is needed
	return nil
}

// autorefresh implements the infinite loop and handles different channels signals
func (af *AutoRefreshX) autorefresh(ctx context.Context) {
	for {
		select {
		case <-af.shutdownCh:
			af.logger.Debug("Terminate refresh loop")
			if af.obj != nil {
				af.obj.Close()
			}
			close(af.shutdownCh)
			return
		case config := <-af.configCh:
			af.logger.Debugf("Reconfigure with %v", config)
			// TODO @shipperizer add mutex Lock
			af.tlsConfig = config.TLSConfigInterface
			af.saslConfig = config.SASLConfigInterface
		case tick := <-af.refreshTicker.C:
			af.logger.Debugf("Tick at %v", tick)

			_, err := af.refresh(ctx)

			if err != nil {
				af.logger.Debugf("refresh has failed: %v", err)

				if m, err := af.monitor.GetMetric("labs_stream_errors"); err != nil {
					af.logger.Debugf("Error fetching metric: %s; keep going....", err)
				} else {
					m.Inc(map[string]string{"task": "tls-refresh", "service": af.monitor.GetService()})
				}
			} else {
				af.incrRefreshMetric()
			}
		}
	}
}

// incrRefreshMetric is a helper for increasing the appropriate monitoring metric
func (af *AutoRefreshX) incrRefreshMetric() {
	c := af.obj.Config()

	if c == nil {
		af.logger.Errorf("object af.obj is of unknown type")

		return
	}

	switch cfg := c.(type) {
	case core.WriterConfigInterface:
		if m, err := af.monitor.GetMetric("labs_stream_refresh_publisher_v1"); err != nil {
			af.logger.Debugf("Error fetching metric: %s; keep going....", err)
		} else {
			m.Inc(map[string]string{"app": af.monitor.GetService(), "topic": cfg.GetTopic()})
		}
	case core.ReaderConfigInterface:
		if m, err := af.monitor.GetMetric("labs_stream_refresh_subscriber_v1"); err != nil {
			af.logger.Debugf("Error fetching metric: %s; keep going....", err)
		} else {
			m.Inc(map[string]string{"app": af.monitor.GetService(), "topic": cfg.GetTopic(), "consumer_group": cfg.GetGroupID()})
		}
	default:
		af.logger.Errorf("object af.obj is of unknown type")
	}
}

// Configure allows to pass a new tls config to the autorefresh
func (af *AutoRefreshX) Configure(ctx context.Context, tlsConfig config.TLSConfigInterface, saslConfig config.SASLConfigInterface) {
	if tlsConfig == nil || saslConfig == nil {
		af.logger.Errorf("empty configs passed - tls: %v - sasl: %v", tlsConfig, saslConfig)

		return
	}

	c := channelConfig{
		TLSConfigInterface:  tlsConfig,
		SASLConfigInterface: saslConfig,
	}

	af.configCh <- c

}

// Object will return the refreshable object (will need to be casted to Writer or Reader)
func (af *AutoRefreshX) Object(ctx context.Context) (core.RefreshableInterface, error) {
	if af.obj == nil {
		af.fetch()
	}

	return af.obj, nil
}

// Refresh force a manual refresh of the secrets
func (af *AutoRefreshX) Refresh(ctx context.Context) (core.RefreshableInterface, error) {
	af.refresh(ctx)

	return af.obj, nil
}

// Stop sends a message on the shutdown channel which will bring the infinite loop to a halt
func (af *AutoRefreshX) Stop() {
	defer func() {
		if recover() == nil {
			return
		}

		af.logger.Warn("AutoRefreshX object has been closed already")
	}()

	af.shutdownCh <- true
}

// Stats returns reader or writer stats (will need to be casted)
func (af *AutoRefreshX) Stats() interface{} {
	// lock reads for af.obj
	af.mutexObj.RLock()
	defer af.mutexObj.RUnlock()
	stats := af.obj.Stats()

	return stats
}

// custom prometheus metrics setup
// ###################################################################################
func (af *AutoRefreshX) registerMetrics() error {
	m := []monitoring.MetricInterface{
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_errors", "task", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_errors_v1", "app", "function", "error"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_refresh_subscriber_v1", "app", "consumer_group", "topic"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_refresh_publisher_v1", "app", "topic"),
	}

	return af.monitor.AddMetrics(m...)
}

// NewAutoRefreshX creates an object implementing AutoRefreshXInterface
func NewAutoRefreshX(ctx context.Context, cfg AutoRefreshXConfigInterface, refreshable core.RefreshableInterface, logger logging.LoggerInterface, monitor monitoring.MonitorInterface) AutoRefreshXInterface {
	af := new(AutoRefreshX)

	af.mutexObj = cfg.GetMutexObj()
	// TODO @shipperizer make refresh time configurable
	af.refreshTicker = time.NewTicker(cfg.GetRefreshTimeout())
	af.configCh = make(chan channelConfig)
	af.shutdownCh = make(chan bool)
	// set refreshable to the one passed in and cast it to core.RefreshableInterface
	// line below will panig
	if refreshable == nil {
		panic("refreshable object is empty, needs to be a core.Reader or core.Writer")
	}

	af.obj = refreshable
	af.tlsConfig = cfg.GetTLSConfig()
	af.saslConfig = cfg.GetSASLConfig()
	af.monitor = monitor
	af.logger = logger

	if logger == nil {
		af.logger = logging.NewLogger()
	}

	_ = af.registerMetrics()

	go af.autorefresh(ctx)

	return af
}
