package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewReaderConfigImplementsInterface(t *testing.T) {
	bootstrapServers := []string{"test.1.svc.local:9094", "test.2.svc.local:9094", "test.3.svc.local:9094"}
	topic := "test.v1"
	groupID := "test"

	timeout := 1 * time.Minute

	cfg := NewConfig(1*time.Hour, nil, nil, nil)
	readerCfg := NewReaderConfig(cfg, bootstrapServers, topic, groupID, 3, timeout)

	assert := assert.New(t)

	assert.Equal(bootstrapServers, readerCfg.GetBootstrapServers(), "Bootstrap servers should match")
	assert.Equal(3, readerCfg.Workers(), "Workers should match")
	assert.Equal(topic, readerCfg.GetTopic(), "Topics should match")
	assert.Equal(groupID, readerCfg.GetGroupID(), "GroupIDs should match")
	assert.Equal(timeout, readerCfg.GetReadTimeout(), "timeout should match")
	assert.Nil(readerCfg.GetTLSConfig(), "TLS config should be empty")
}
