package publisher

import uuid "github.com/google/uuid"

// Message is an abstraction based around kafka.Message
type Message struct {
	key   []byte
	value interface{}
}

// GetKey returns the key of the message
func (m *Message) GetKey() []byte {
	return m.key
}

// GetValue returns the value of the message
func (m *Message) GetValue() interface{} {
	return m.value
}

// NewMessage creates a new object implementing MessageInterface
func NewMessage(key string, value interface{}) MessageInterface {
	msg := &Message{
		key:   []byte(key),
		value: value,
	}

	if key == "" {
		msg.key = []byte(uuid.New().String())
	}

	return msg
}
