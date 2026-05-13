# Kilo Franz

![test](https://github.com/shipperizer/kilo-franz/workflows/test/badge.svg)
![release](https://github.com/shipperizer/kilo-franz/workflows/release/badge.svg)
[![codecov](https://codecov.io/gh/shipperizer/kilo-franz/branch/main/graph/badge.svg)](https://codecov.io/gh/shipperizer/kilo-franz)
[![Go Reference](https://pkg.go.dev/badge/github.com/shipperizer/kilo-franz.svg)](https://pkg.go.dev/github.com/shipperizer/kilo-franz)

Opinionated Kafka consumer and producer library built on top of [segmentio/kafka-go](https://github.com/segmentio/kafka-go).

## Features

- Channel-based consumer with configurable worker pools
- Standard (single-goroutine) consumer for simpler use cases
- Producer with support for multiple topics via nicknames
- Pluggable encoder interface (JSON by default, protobuf-ready)
- Automatic TLS certificate refresh
- SASL authentication support
- Prometheus metrics integration

## Installation

```bash
go get github.com/shipperizer/kilo-franz
```

## Documentation

For the full API reference, run:

```bash
godoc -http=:6060
```

Then visit `http://localhost:6060/pkg/github.com/shipperizer/kilo-franz/` in your browser.

## Usage

### Creating a Consumer

The `ChannelConsumer` uses a dedicated goroutine for reading messages from Kafka
and distributes work across N worker goroutines via a channel:

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/subscriber"
)

func main() {
	// Create base config with TLS and a 5-minute refresh interval
	cfg := config.NewConfig(5*time.Minute, &tlsSetup, nil, nil)

	// Create reader config specifying brokers, topic, consumer group, workers, and read timeout
	readerCfg := config.NewReaderConfig(
		cfg,
		strings.Split("kafka-broker-1:9092,kafka-broker-2:9092", ","),
		"my-topic",
		"my-app.consumer-group",
		5,               // number of worker goroutines
		30*time.Second,  // read timeout
	)

	// Create the reader (implements RefreshableInterface)
	reader := core.NewReader(readerCfg)

	// Create the channel consumer
	consumer, err := subscriber.NewChannelConsumer(reader, myService, monitor)
	if err != nil {
		panic(err)
	}

	// Start consuming
	consumer.Start()

	// Wait for termination signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// Graceful shutdown
	consumer.Stop()
	fmt.Println("Shutting down")
}
```

The `StandardConsumer` reads and processes messages sequentially in a single goroutine:

```go
consumer, err := subscriber.NewStandardConsumer(reader, myService, monitor)
if err != nil {
    panic(err)
}

consumer.Start()
```

### Creating a Producer

```go
package main

import (
	"fmt"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/core"
	"github.com/shipperizer/kilo-franz/encoder"
	"github.com/shipperizer/kilo-franz/publisher"
)

func main() {
	cfg := config.NewConfig(5*time.Minute, &tlsSetup, nil, nil)

	// Create a writer config with the default JSON encoder
	writerCfg := config.NewWriterConfig(
		cfg,
		[]string{"kafka-broker-1:9092", "kafka-broker-2:9092"},
		"my-topic",
		"my-topic-nickname", // used to identify this writer when publishing
		false,               // async mode
		nil,                 // nil defaults to JSONEncoder
	)

	// Create the writer (implements RefreshableInterface)
	writer := core.NewWriter(writerCfg)

	// Create the producer with one or more writers
	producer := publisher.NewProducer(monitor, writer)
	defer producer.Close()

	// Publish messages using the topic nickname
	msg := publisher.NewMessage("message-key", map[string]string{"hello": "world"})
	if err := producer.Publish("my-topic-nickname", msg); err != nil {
		fmt.Printf("failed to publish: %v\n", err)
	}
}
```

### Custom Encoders

The `EncoderInterface` allows you to define custom serialization for messages.
By default, a JSON encoder is used. To use Protocol Buffers or any other format,
implement the `EncoderInterface`:

```go
package main

import (
	"google.golang.org/protobuf/proto"

	"github.com/shipperizer/kilo-franz/encoder"
)

// ProtoEncoder encodes messages using Protocol Buffers.
// The message passed to Encode must implement proto.Message.
type ProtoEncoder struct{}

// Encode marshals the given message using protobuf serialization.
// Returns the binary representation of the message or an error if
// the message does not implement proto.Message or marshaling fails.
func (e *ProtoEncoder) Encode(msg interface{}) ([]byte, error) {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("message does not implement proto.Message: %T", msg)
	}
	return proto.Marshal(protoMsg)
}

// NewProtoEncoder returns an encoder implementing EncoderInterface
// that serializes messages using Protocol Buffers.
func NewProtoEncoder() encoder.EncoderInterface {
	return &ProtoEncoder{}
}
```

Then pass it to `NewWriterConfig`:

```go
writerCfg := config.NewWriterConfig(
    cfg,
    brokers,
    "my-topic",
    "my-topic-nickname",
    false,
    NewProtoEncoder(), // use protobuf encoding
)
```

### Implementing the ServiceInterface

Consumers require a `ServiceInterface` implementation that processes each message:

```go
// ServiceInterface is what the consumer calls for each message.
type ServiceInterface interface {
    TaskName() string
    Flow(MessageKey, MessageValue []byte) error
}
```

Example implementation:

```go
type MyService struct{}

func (s *MyService) TaskName() string {
    return "my-processing-task"
}

func (s *MyService) Flow(key, value []byte) error {
    // Process the message
    fmt.Printf("Received message: key=%s value=%s\n", key, value)
    return nil
}
```

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                    Publisher                         │
│  ┌───────────┐    ┌───────────┐    ┌───────────┐   │
│  │  Writer 1 │    │  Writer 2 │    │  Writer N │   │
│  │ (topic A) │    │ (topic B) │    │ (topic N) │   │
│  └───────────┘    └───────────┘    └───────────┘   │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│               ChannelConsumer                        │
│  ┌────────┐     ┌──────────────────────────────┐   │
│  │ Reader │────▶│  Channel  ──▶ Worker 1       │   │
│  │        │     │            ──▶ Worker 2       │   │
│  │        │     │            ──▶ Worker N       │   │
│  └────────┘     └──────────────────────────────┘   │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│              StandardConsumer                        │
│  ┌────────┐     ┌──────────────────────────────┐   │
│  │ Reader │────▶│  Process inline (1 routine)  │   │
│  └────────┘     └──────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

## Package Structure

| Package | Description |
|---------|-------------|
| `config` | Configuration objects for readers, writers, TLS, and SASL |
| `core` | Reader and Writer wrappers around kafka-go with auto-refresh support |
| `encoder` | Encoder interface and default JSON implementation |
| `logging` | Logger interface abstraction |
| `monitoring` | Prometheus metrics interface and helpers |
| `publisher` | Kafka producer with multi-topic support |
| `refresh` | Auto-refresh mechanism for TLS credentials |
| `sasl` | SASL authentication configuration |
| `subscriber` | Channel-based and standard Kafka consumers |
| `tls` | TLS configuration helpers |
| `vault` | Secret retrieval from AWS Secrets Manager |

## License

See [LICENSE](LICENSE) for details.
