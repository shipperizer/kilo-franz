// Package encoder provides the EncoderInterface for serializing messages
// before they are published to Kafka.
//
// The default implementation is [JSONEncoder], which marshals messages using
// encoding/json. Custom encoders (e.g., for Protocol Buffers) can be created
// by implementing the [EncoderInterface].
//
// # Protobuf Encoder Example
//
// To use Protocol Buffers for message serialization:
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
//	func NewProtoEncoder() encoder.EncoderInterface {
//		return &ProtoEncoder{}
//	}
//
// Pass the encoder when creating a writer config:
//
//	writerCfg := config.NewWriterConfig(cfg, brokers, topic, nickname, false, NewProtoEncoder())
package encoder
