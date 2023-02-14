package subscriber

type ServiceInterface interface {
	TaskName() string
	Flow(MessageKey, MessageValue []byte) error
}
