package config

import (
	"sync"
)

// AutoRefreshXConfig is the config needed for the autorefresh object
// it's composed by a ConfigInterface object and therefore inherits all the methods of ConfigInterface
type AutoRefreshXConfig struct {
	ConfigInterface
	mutexObj *sync.RWMutex
}

// GetMutexObj returns a pointer to a sync.RWMutex object
func (c *AutoRefreshXConfig) GetMutexObj() *sync.RWMutex {
	return c.mutexObj
}

// NewAutoRefreshXConfig returns an object implementing AutoRefreshXConfigInterface
func NewAutoRefreshXConfig(mutexObj *sync.RWMutex, config ConfigInterface) AutoRefreshXConfigInterface {
	return &AutoRefreshXConfig{
		ConfigInterface: config,
		mutexObj:        mutexObj,
	}
}
