package config

import (
	"testing"
	"time"

	enc "github.com/shipperizer/kilo-franz/encoder"
	"github.com/stretchr/testify/assert"
)

func TestNewWriterConfigImplementsInterface(t *testing.T) {
	brokers := []string{"test.1.svc.local:9094", "test.2.svc.local:9094", "test.3.svc.local:9094"}
	topic := "test.v1"
	nickname := "test"
	async := true

	cfg := NewConfig(1*time.Hour, nil, nil)
	writerCfg := NewWriterConfig(cfg, brokers, topic, nickname, async, nil)

	assert := assert.New(t)

	assert.Equal(brokers, writerCfg.GetBrokers(), "Brokers should match")
	assert.Equal(topic, writerCfg.GetTopic(), "Topics should match")
	assert.Equal(nickname, writerCfg.GetNickname(), "Nicknames should match")
	assert.True(writerCfg.GetAsync(), "Async should be set to true")
	assert.Nil(writerCfg.GetTLS(), "TLS config should be empty")
	assert.NotNil(writerCfg.GetEncoder(), "Default Encoder should be JSONEncoder if nil is passed")
}

func TestNewWriterConfigDefaultEncoderJSON(t *testing.T) {
	brokers := []string{"test.1.svc.local:9094", "test.2.svc.local:9094", "test.3.svc.local:9094"}
	topic := "test.v1"
	nickname := "test"
	async := true

	cfg := NewConfig(1*time.Hour, nil, nil)
	writerCfg := NewWriterConfig(cfg, brokers, topic, nickname, async, nil)

	assert := assert.New(t)

	assert.IsType(&enc.JSONEncoder{}, writerCfg.GetEncoder(), "Default Encoder should be of type *JSONEncoder if nil is passed")
}

func TestNewWriterConfigDefaultNicknameIsTopic(t *testing.T) {
	brokers := []string{"test.1.svc.local:9094", "test.2.svc.local:9094", "test.3.svc.local:9094"}
	topic := "test.v1"
	async := true

	cfg := NewConfig(1*time.Hour, nil, nil)
	writerCfg := NewWriterConfig(cfg, brokers, topic, "", async, nil)

	assert := assert.New(t)

	assert.Equal(writerCfg.GetTopic(), writerCfg.GetNickname(), "Nickname should default to topic if left empty")
}
