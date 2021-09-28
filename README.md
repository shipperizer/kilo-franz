# Labs Stream Library

![test](https://github.com/shipperizer/kilo-franz/workflows/test/badge.svg)
![release](https://github.com/shipperizer/kilo-franz/workflows/release/badge.svg)
[![codecov](https://codecov.io/gh/shipperizer/kilo-franz/branch/master/graph/badge.svg)](https://codecov.io/gh/shipperizer/kilo-franz)

Library used for dealing with Kafka consumers and producers

***Uses segmentio/kafka-go v0.4.17***


## Release process

# DOCS

To have a better look at API reference do `godoc -http=:6060` and then check the browser at `http://localhost:6060/pkg/github.com/shipperizer/kilo-franz/`


# HOWTO

* to create a consumer:

```
// ChannelConsumer is an implementation of the ConsumerInterface
// it will work with 1 goroutine taking care of pulling messages and
// #N workers (defined on constructor)
// Example:

cfg := streamConfig.NewConfig(5*time.Minute, &tlsSetup, nil)
readerCfg := streamConfig.NewReaderConfig(
	cfg,
	strings.Split(viper.GetString("kafka.url"), ","),
	viper.GetString("kafka.consumer.topic"),
	"labs-audit-api.cgroup",
	5,
)
reader := core.NewReader(readerCfg)

consumer, err := subscriber.NewChannelConsumer(
	reader,
	dummy.NewService(
		store.NewStore(
			store.StoreTableConfig{
				Logs: fmt.Sprint(tablePrefix, viper.GetString("dynamodb.tables.audit.logs")),
			},
			dynamoClient,
		),
		monitor,
		readerCfg.GetGroupID(),
	),
	monitor,
)
if err != nil {
	panic(err)
}
consumer.Start()
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
Block until we receive our signal.
<-c
consumer.Stop()
log.Info("Shutting down")
os.Exit(0)
```