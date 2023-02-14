package config

import (
	enc "github.com/shipperizer/kilo-franz/encoder"
)

// WriterConfig is a config for the core.Writer, holds subscriber informations
type WriterConfig struct {
	brokers  []string // strings.Split(viper.GetString("kafka.url"), ",")
	topic    string
	nickname string
	async    bool
	encoder  enc.EncoderInterface
	ConfigInterface
}

// GetAsync returns if the publisher is supposed to be setup as async or sync
func (c *WriterConfig) GetAsync() bool {
	return c.async
}

// GetNickname returns the nickname for the associated topic
func (c *WriterConfig) GetNickname() string {
	return c.nickname
}

// GetTopic returns the configured topic
func (c *WriterConfig) GetTopic() string {
	return c.topic
}

// GetBrokers returns the configured list of brokers
func (c *WriterConfig) GetBrokers() []string {
	return c.brokers
}

// GetEncoder returns the configured encoder
func (c *WriterConfig) GetEncoder() enc.EncoderInterface {
	return c.encoder
}

// NewWriterConfig creates a new object implementing WriterConfigInterface
func NewWriterConfig(cfg ConfigInterface, brokers []string, topic, nickname string, async bool, encoder enc.EncoderInterface) *WriterConfig {
	c := WriterConfig{
		topic:           topic,
		nickname:        nickname,
		async:           async,
		brokers:         brokers,
		ConfigInterface: cfg,
		encoder:         encoder,
	}

	if nickname == "" {
		c.nickname = topic
	}

	if encoder == nil {
		c.encoder = enc.NewJSONEncoder()
	}

	return &c
}
