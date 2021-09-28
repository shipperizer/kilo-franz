package encoder

import "encoding/json"

// JSONEncoder is the default encoder for publihsers/producers
// it implements the EncoderInterface and simply marshals an interface to json
type JSONEncoder struct{}

// Encode is the only method, needs to return ([]byte, error), in this case simply a wrapper around
// the json.Marshal funcrton
func (e *JSONEncoder) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

// NewJSONEncoder returns an object implementing EncoderInterface
func NewJSONEncoder() EncoderInterface {
	return &JSONEncoder{}
}
