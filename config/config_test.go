package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigImplementsInterface(t *testing.T) {
	cfg := NewConfig(1*time.Hour, nil, nil)

	assert := assert.New(t)

	assert.Equal(1*time.Hour, cfg.GetRefreshTimeout(), "Refresh timeouts should match")
	assert.Nil(cfg.GetTLSConfig(), "TLS config should be empty")
	assert.NotNil(cfg.GetLogger(), "Default Logger should be *zap.SugaredLogger if nil is passed")
}
