package publisher

import (
	"github.com/segmentio/kafka-go"
)

// MessageInterface is an abstraction on top of kafka.Message
type MessageInterface interface {
	GetKey() []byte
	GetValue() interface{}
}

// ProducerInterface defines the methods exposed by publishers/producers
type ProducerInterface interface {
	ListTopics() map[string]string
	Publish(topicNickname string, messages ...MessageInterface) error
	Stats(topicNickname string) kafka.WriterStats
	Close()
}
