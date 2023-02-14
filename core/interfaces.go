package core

import (
	"context"
	"time"

	"github.com/shipperizer/kilo-franz/config"
	"github.com/shipperizer/kilo-franz/encoder"
)

// ReaderConfigInterface is the core.Reader config interface, embeds config.ConfigInterface
type ReaderConfigInterface interface {
	config.ConfigInterface
	GetBootstrapServers() []string
	GetTopic() string
	GetGroupID() string
	GetReadTimeout() time.Duration
	Workers() int
}

// WriterConfigInterface is the core.Writer config interface, embeds config.ConfigInterface
type WriterConfigInterface interface {
	config.ConfigInterface
	GetBrokers() []string
	GetTopic() string
	GetNickname() string
	GetAsync() bool
	GetEncoder() encoder.EncoderInterface
}

// RefreshableInterface is implemented by core.Writer and core.Reader so that they can be refreshed by AutoRefreshX
type RefreshableInterface interface {
	Get(context.Context) (interface{}, error)
	Config() interface{}
	Close()
	Stats() interface{}
	Renew(config.TLSConfigInterface, config.SASLConfigInterface, ...interface{})
}
