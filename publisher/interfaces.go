package publisher

// MessageInterface is an abstraction on top of kafka.Message
type MessageInterface interface {
	GetKey() []byte
	GetValue() interface{}
}
