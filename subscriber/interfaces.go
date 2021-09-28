package subscriber

import (
	"github.com/segmentio/kafka-go"
)

type ConsumerInterface interface {
	Stats() kafka.ReaderStats
	Start()
	Stop()
}

type ServiceInterface interface {
	TaskName() string
	Flow(MessageKey, MessageValue []byte) error
}
