package subscriber

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/shipperizer/kilo-franz/monitoring"
	"go.uber.org/zap"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/refresh"
)

// ChannelConsumer is an implementation of the ConsumerInterface
// it will work with 1 goroutine taking care of pulling messages and
// #N workers (defined on constructor)
// Example:
//
// cfg := streamConfig.NewConfig(5*time.Minute, &tlsSetup, nil)
//
// readerCfg := streamConfig.NewReaderConfig(
// 	cfg,
// 	strings.Split(viper.GetString("kafka.url"), ","),
// 	viper.GetString("kafka.consumer.topic"),
// 	"test-app.cgroup",
// 	5,
// )
//
// reader := core.NewReader(readerCfg)

// consumer, err := subscriber.NewChannelConsumer(
// 	reader,
// 	dummy.NewService(
// 		store.NewStore(
// 			store.StoreTableConfig{
// 				Logs: fmt.Sprint(tablePrefix, viper.GetString("dynamodb.tables.audit.logs")),
// 			},
// 			dynamoClient,
// 		),
// 		monitor,
// 		readerCfg.GetGroupID(),
// 	),
// 	monitor,
// )
//
// if err != nil {
// 	panic(err)
// }
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
type ChannelConsumer struct {
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

	messageCh  chan kafka.Message
	shutdownCh chan bool

	// external dependencies
	svc     ServiceInterface
	monitor monitoring.MonitorInterface

	logger *zap.SugaredLogger
}

// unwrapReader is a helper methods to remove the different interfaces and reach the final kafka.Reader
func (c *ChannelConsumer) unwrapReader() (*kafka.Reader, error) {
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

// consume takes care of running an infinite loop and run the service.Flow method
func (c *ChannelConsumer) consume(workerID int) {
	c.logger.Infof("starting worker %v", workerID)
	for {
		select {
		case <-c.shutdownCh:
			c.logger.Infof("shutting down worker %v", workerID)
			c.wg.Done()
			return
		case msg := <-c.messageCh:
			err := c.svc.Flow(msg.Key, msg.Value)

			if err != nil {
				c.monitor.Incr("errors", map[string]string{"task": c.svc.TaskName(), "service": c.groupID})
				c.logger.Warn("Moving on...")
				continue
			}

			c.monitor.Incr("count", map[string]string{"task": c.svc.TaskName(), "service": c.groupID})
		}
	}
}

// read is the infinite loop for the kafka reading part of the consumer
func (c *ChannelConsumer) read() {
	c.logger.Infof("Listening to kafka on topic %s", c.Stats().Topic)

	metricLabels := map[string]string{
		"topic":          strings.ReplaceAll(c.Stats().Topic, ".", "_"),
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

			c.monitor.Gauge("lag", c.Stats().Lag, metricLabels)

			if err == context.DeadlineExceeded {
				c.logger.Debug("context timed out reading from kafka, releasing lock")
				continue
			}

			if err != nil {
				c.logger.Errorf("failed fetching kafka message: %s", err)
				c.logger.Debug("demand a refresh of the kafka reader")
				c.af.Refresh(context.TODO())
				continue
			}

			c.monitor.Incr("events", metricLabels)

			c.messageCh <- msg

		}
	}
}

// Stop makes sure goroutines for read and consume are being gracefully stopped
func (c *ChannelConsumer) Stop() {
	// send reader shutdown
	c.shutdownCh <- true

	// send #workers shutdowns
	for i := 0; i < c.workers; i++ {
		c.shutdownCh <- true
	}

	c.wg.Wait()
	defer c.logger.Desugar().Sync()
}

// Start creates all the goroutines for read and consume
func (c *ChannelConsumer) Start() {

	for i := 0; i < c.workers; i++ {
		c.wg.Add(1)
		go c.consume(i + 1)
	}
	c.wg.Add(1)
	go c.read()
}

// Stats returns kafka.ReaderStats
func (c *ChannelConsumer) Stats() kafka.ReaderStats {
	if stats, ok := c.af.Stats().(kafka.ReaderStats); ok {
		return stats
	}

	panic("consumer reader is setup incorrectly, not of type *kafka.Reader")
}

// NewChannelConsumer creates an object implementing ConsumerInterface
func NewChannelConsumer(refreshable config.RefreshableInterface, svc ServiceInterface, monitor monitoring.MonitorInterface) (ConsumerInterface, error) {
	consumer := new(ChannelConsumer)

	c := refreshable.Config()

	if c == nil {
		return nil, fmt.Errorf("config is empty: %v", c)
	}

	cfg, ok := c.(config.ReaderConfigInterface)

	if !ok {
		return nil, fmt.Errorf("config is non parsable, wrong type, should be config.ReaderConfigInterface: %v", c)
	}

	consumer.svc = svc
	consumer.monitor = monitor
	consumer.groupID = cfg.GetGroupID()
	consumer.workers = cfg.Workers()
	consumer.logger = cfg.GetLogger()
	consumer.readTimeout = cfg.GetRefreshTimeout()

	consumer.messageCh = make(chan kafka.Message, consumer.workers)
	consumer.shutdownCh = make(chan bool)

	consumer.af = refresh.NewAutoRefreshX(context.TODO(), config.NewAutoRefreshXConfig(&consumer.mutexReader, cfg), refreshable, consumer.logger, monitor)

	consumer.logger.Debug("config:", cfg)

	return consumer, nil
}
