package subscriber

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/monitoring"
	"github.com/shipperizer/kilo-franz/refresh"
)

// StandardConsumer is an implementation of the ConsumerInterface
// it will work with single function taking care of pulling messages and
// processing it inline
// Example:
//
// cfg := streamConfig.NewConfig(5*time.Minute, &tlsSetup, nil)
//
// readerCfg := streamConfig.NewReaderConfig(
// 	cfg,
// 	strings.Split(viper.GetString("kafka.url"), ","),
// 	viper.GetString("kafka.consumer.topic"),
// 	"labs-audit-api.cgroup",
// 	0,
// )
//
// reader := core.NewReader(readerCfg)

// consumer, err := subscriber.NewStandardConsumer(
//
//	reader,
//	dummy.NewService(
//		store.NewStore(
//			store.StoreTableConfig{
//				Logs: fmt.Sprint(tablePrefix, viper.GetString("dynamodb.tables.audit.logs")),
//			},
//			dynamoClient,
//		),
//		monitor,
//		readerCfg.GetGroupID(),
//	),
//	monitor,
//
// )
//
//	if err != nil {
//		panic(err)
//	}
//
// consumer.Start()
//
// c := make(chan os.Signal, 1)
// signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
//
// // Block until we receive our signal.
// <-c
//
// consumer.Stop()
// log.Info("Shutting down")
// os.Exit(0)
type StandardConsumer struct {
	// embed AutoRefresh object to simplify mutex setup
	af refresh.AutoRefreshXInterface

	// mutex
	mutexReader sync.RWMutex

	// goroutine control
	wg sync.WaitGroup

	// generic config
	workers     int
	groupID     string
	readTimeout time.Duration

	shutdownCh chan bool

	// external dependencies
	svc     ServiceInterface
	monitor monitoring.MonitorInterface

	logger logging.LoggerInterface
}

// unwrapReader is a helper methods to remove the different interfaces and reach the final kafka.Reader
func (c *StandardConsumer) unwrapReader() (*kafka.Reader, error) {
	r, err := c.af.Object(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed fetching kafka reader: %s", err)

	}

	rr, ok := r.(*core.Reader)

	if !ok {
		return nil, fmt.Errorf("kafka reader is not of type core.Reader %v", rr)
	}

	rrI, err := r.Get(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("failed fetching kafka reader interface: %s", err)

	}

	reader, ok := rrI.(*kafka.Reader)

	if !ok {
		return nil, fmt.Errorf("reader is not of type *kafka.Reader, casting failed")
	}

	return reader, nil
}

// start is the infinite loop for the kafka reading and the execution of the service.Flow method
func (c *StandardConsumer) start() {
	c.logger.Infof("Listening to kafka on topic %s", c.Stats().Topic)

	tags := map[string]string{
		"topic":          strings.Join(strings.Split(c.Stats().Topic, "."), "_"),
		"consumer_group": c.groupID,
		"service":        c.groupID,
	}

	for {
		select {
		case <-c.shutdownCh:
			c.logger.Info("shutting down reader")
			c.af.Stop()
			c.wg.Done()
			return
		default:
			c.mutexReader.RLock()

			reader, err := c.unwrapReader()

			if err != nil {
				c.logger.Errorf("failed unwrapping kafka reader: %s", err)
				c.mutexReader.RUnlock()
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), c.readTimeout)

			msg, err := reader.ReadMessage(ctx)
			// RUnlock as first command, as af.Stats blocks as well with an RLock
			c.mutexReader.RUnlock()
			cancel()

			if m, err := c.monitor.GetMetric("labs_stream_lag"); err != nil {
				c.logger.Debugf("Error fetching metric: %s; keep going....", err)
			} else {
				m.Set(float64(c.Stats().Lag), tags)
			}

			if err == context.DeadlineExceeded {
				c.logger.Debug("context timed out reading from kafka, releasing lock")
				continue
			}

			if m, err := c.monitor.GetMetric("labs_stream_events"); err != nil {
				c.logger.Debugf("Error fetching metric: %s; keep going....", err)
			} else {
				m.Inc(tags)
			}

			if err != nil {
				c.logger.Errorf("failed fetching kafka message: %s", err)
				c.logger.Debug("demand a refresh of the kafka reader")
				c.af.Refresh(context.TODO())
				// set message to an empty one (might already be that)
				msg = kafka.Message{}
			}

			startTime := time.Now()

			err = c.svc.Flow(msg.Key, msg.Value)

			if m, e := c.monitor.GetMetric("labs_stream_handling_seconds_v1"); e == nil {
				status := "ok"

				if err != nil {
					status = "error"
				}

				m.Observe(time.Since(startTime).Seconds(), map[string]string{"task": c.svc.TaskName(), "service": c.groupID, "status": status})
			}

			if err != nil {
				c.logger.Error(err)

				if m, err := c.monitor.GetMetric("labs_stream_errors"); err != nil {
					c.logger.Debugf("Error fetching metric: %s; keep going....", err)
				} else {
					m.Inc(tags)
				}

				c.logger.Warn("Moving on...")
				continue
			}

			if m, err := c.monitor.GetMetric("labs_stream_count"); err != nil {
				c.logger.Debugf("Error fetching metric: %s; keep going....", err)
			} else {
				m.Inc(tags)
			}
		}
	}
}

// Start creates all the goroutine for consuming and reading
func (c *StandardConsumer) Start() {
	// add the Start function as something to wait for by the waitgroup
	// will be marked Done when Stop() is called and autorefresh has been called off
	c.wg.Add(1)
	go c.start()
}

// Stop makes sure goroutines for read and consume are being gracefully stopped
func (c *StandardConsumer) Stop() {
	// send reader shutdown
	c.logger.Error("shutting down")
	c.shutdownCh <- true
	c.wg.Wait()
}

// Stats returns kafka.ReaderStats
func (c *StandardConsumer) Stats() kafka.ReaderStats {
	if stats, ok := c.af.Stats().(kafka.ReaderStats); ok {
		return stats
	}

	panic("consumer reader is setup incorrectly, not of type *kafka.Reader")
}

// custom prometheus metrics setup
// ###################################################################################
func (c *StandardConsumer) registerMetrics() error {
	m := []monitoring.MetricInterface{
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_lag", "topic", "consumer_group", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_events", "topic", "consumer_group", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_count", "task", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_time", "task", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_errors", "task", "service"),
		monitoring.NewMetric(monitoring.GAUGE, "labs_stream_batch_timeout", "task", "service"),
		monitoring.NewMetric(monitoring.HISTOGRAM, "labs_stream_handling_seconds_v1", "task", "service", "status"),
	}

	return c.monitor.AddMetrics(m...)
}

// NewStandardConsumer creates an object implementing ConsumerInterface
func NewStandardConsumer(refreshable core.RefreshableInterface, svc ServiceInterface, monitor monitoring.MonitorInterface) (*StandardConsumer, error) {
	consumer := new(StandardConsumer)

	c := refreshable.Config()

	if c == nil {
		return nil, fmt.Errorf("config is empty: %v", c)
	}

	cfg, ok := c.(core.ReaderConfigInterface)

	if !ok {
		return nil, fmt.Errorf("config is non parsable, wrong type, should be config.ReaderConfigInterface: %v", c)
	}

	consumer.svc = svc
	consumer.monitor = monitor
	consumer.groupID = cfg.GetGroupID()
	consumer.workers = 1
	consumer.logger = cfg.GetLogger()
	consumer.shutdownCh = make(chan bool)
	consumer.readTimeout = cfg.GetRefreshTimeout()

	consumer.af = refresh.NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&consumer.mutexReader, cfg), refreshable, consumer.logger, monitor)

	consumer.logger.Debug("config:", cfg)

	_ = consumer.registerMetrics()

	return consumer, nil
}
