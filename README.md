# Kilo Franz: opinionated Kafka Consumer and Producer code built on top of segmentio/kafka-go library

![test](https://github.com/shipperizer/kilo-franz/workflows/test/badge.svg)
![release](https://github.com/shipperizer/kilo-franz/workflows/release/badge.svg)
[![codecov](https://codecov.io/gh/shipperizer/kilo-franz/branch/main/graph/badge.svg)](https://codecov.io/gh/shipperizer/kilo-franz)
[![Go Reference](https://pkg.go.dev/badge/github.com/shipperizer/kilo-franz.svg)](https://pkg.go.dev/github.com/shipperizer/kilo-franz)

Kilo Franz is a small Go library that wraps [segmentio/kafka-go](https://github.com/segmentio/kafka-go)
to provide opinionated Kafka consumers and producers. It takes care of TLS/SASL
refresh, reader/writer lifecycle management, pluggable encoders, and a producer
that can target multiple topics through nicknames.

***Uses segmentio/kafka-go v0.4.17***


## Examples

Runnable, minimal examples live under [`examples/`](./examples):

- [`examples/config`](./examples/config) - build a `config.Config` with no external dependencies.
- [`examples/encoder`](./examples/encoder) - encode a struct with `encoder.NewJSONEncoder` and pretty print it.

Run them with:

```sh
go run ./examples/config
go run ./examples/encoder
```


## Release process

# DOCS

To have a better look at API reference do `godoc -http=:6060` and then check the browser at `http://localhost:6060/pkg/github.com/shipperizer/kilo-franz/`


# HOWTO

## Minimal consumer

```go
cfg := config.NewConfig(5*time.Minute, nil, nil, nil)

readerCfg := config.NewReaderConfig(
	cfg,
	[]string{"localhost:9092"},
	"events",
	"my-service.cgroup",
	1,
)

reader := core.NewReader(readerCfg)

consumer, err := subscriber.NewStandardConsumer(reader, myService, monitor)
if err != nil {
	panic(err)
}

consumer.Start()

c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt, syscall.SIGTERM)
<-c

consumer.Stop()
```

## Minimal producer

```go
cfg := config.NewConfig(5*time.Minute, nil, nil, nil)

writerCfg := config.NewWriterConfig(
	cfg,
	[]string{"localhost:9092"},
	"events",
	"events",
	false,
	encoder.NewJSONEncoder(),
)

writer := core.NewWriter(writerCfg)
producer := publisher.NewProducer(monitor, writer)

if err := producer.Publish("events", publisher.NewMessage("", map[string]string{"hello": "world"})); err != nil {
	panic(err)
}

producer.Close()
```

## Consumer (full example)

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
	"test-app.cgroup",
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
