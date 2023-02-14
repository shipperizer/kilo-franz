package config

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewAutoRefreshXConfigImplementsInterface(t *testing.T) {
	var mutexR sync.RWMutex

	cfg := NewAutoRefreshXConfig(&mutexR, NewConfig(1*time.Hour, nil, nil, nil))

	assert := assert.New(t)

	assert.Equal(&mutexR, cfg.GetMutexObj(), "Reader mutexes should match")
	assert.Equal(1*time.Hour, cfg.GetRefreshTimeout(), "Refresh timeouts should match")
	assert.Nil(cfg.GetTLSConfig(), "TLS config should be empty")
	assert.NotNil(cfg.GetLogger(), "Default Logger should be *zap.SugaredLogger if nil is passed")
}
