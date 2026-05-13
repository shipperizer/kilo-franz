package encoder

// EncoderInterface is an interface each Encoder will have to abide to
// it will be used by the publisher/producer object to marshal payloads
// check default implementation of JSONEncoder for more info
// Example for protobuf:
//
// import "google.golang.org/protobuf/proto"
//
// type ProtoEncoder struct{}
//
//	func (e *ProtoEncoder) Encode(msg any) ([]byte, error) {
//		return proto.Marshal(msg)
//	}
//
//	func NewProtoEncoder() EncoderInterface {
//		return &ProtoEncoder{}
//	}
type EncoderInterface interface {
	Encode(msg any) ([]byte, error)
}
