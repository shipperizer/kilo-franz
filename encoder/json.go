package encoder

import "encoding/json"

// JSONEncoder is the default encoder for publishers/producers.
// It implements [EncoderInterface] by marshaling messages to JSON.
type JSONEncoder struct{}

// Encode serializes the given message to JSON bytes.
func (e *JSONEncoder) Encode(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

// NewJSONEncoder returns a new [JSONEncoder] implementing [EncoderInterface].
func NewJSONEncoder() EncoderInterface {
	return &JSONEncoder{}
}
