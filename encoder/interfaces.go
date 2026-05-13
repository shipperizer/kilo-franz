package encoder

// EncoderInterface is the interface that all encoders must implement.
// It is used by the publisher/producer to marshal message payloads before
// writing them to Kafka.
//
// The default implementation is [JSONEncoder], which serializes messages
// using encoding/json.
//
// To use Protocol Buffers, implement a custom encoder:
//
//	import (
//		"fmt"
//		"google.golang.org/protobuf/proto"
//	)
//
//	type ProtoEncoder struct{}
//
//	func (e *ProtoEncoder) Encode(msg interface{}) ([]byte, error) {
//		protoMsg, ok := msg.(proto.Message)
//		if !ok {
//			return nil, fmt.Errorf("message does not implement proto.Message: %T", msg)
//		}
//		return proto.Marshal(protoMsg)
//	}
//
//	func NewProtoEncoder() EncoderInterface {
//		return &ProtoEncoder{}
//	}
//
// Then pass the encoder to [config.NewWriterConfig]:
//
//	writerCfg := config.NewWriterConfig(cfg, brokers, topic, nickname, false, NewProtoEncoder())
type EncoderInterface interface {
	Encode(msg interface{}) ([]byte, error)
}
