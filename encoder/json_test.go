package encoder

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Event struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func TestJSONEncoderMarshals(t *testing.T) {
	encoder := NewJSONEncoder()

	event := &Event{Key: "test", Value: 1}
	expected, _ := json.Marshal(event)

	assert := assert.New(t)

	res, err := encoder.Encode(event)

	assert.Nil(err, "No error should be thrown")
	assert.Equal(expected, res, "Mashalled results should match")

}
