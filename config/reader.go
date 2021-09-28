package config

import "time"

// ReaderConfig is a config for the core.Reader, holds subscriber informations
type ReaderConfig struct {
	bootstrapServers []string // strings.Split(viper.GetString("kafka.url"), ",")
	topic            string
	groupID          string
	workers          int
	readTimeout      time.Duration
	ConfigInterface
}

// GetGroupID returns a consumer group ID
func (c *ReaderConfig) GetGroupID() string {
	return c.groupID
}

// Workers returns the number of workers for a subscriber
func (c *ReaderConfig) Workers() int {
	return c.workers
}

// GetTopic returns the topic chosed to subscribe to
func (c *ReaderConfig) GetTopic() string {
	return c.topic
}

// GetBootstapServers returns the kafka url the config is pointing to
func (c *ReaderConfig) GetBootstrapServers() []string {
	return c.bootstrapServers
}

// GetReadTimeout returns the time the ReadMessage call will be allowed to wait before timing out (via context)
func (c *ReaderConfig) GetReadTimeout() time.Duration {
	return c.readTimeout
}

// NewReaderConfig creates a new object implementing ReaderConfigInterface
func NewReaderConfig(cfg ConfigInterface, bootstrapServers []string, topic, groupID string, noWorkers int, readTimeout time.Duration) ReaderConfigInterface {
	c := ReaderConfig{
		topic:            topic,
		groupID:          groupID,
		workers:          noWorkers,
		bootstrapServers: bootstrapServers,
		readTimeout:      readTimeout,
		ConfigInterface:  cfg,
	}

	return &c
}
