package config

import (
	"context"
	"crypto/tls"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/shipperizer/kilo-franz/encoder"
	kiloTLS "github.com/shipperizer/kilo-franz/tls"
)

// ConfigInterface is the interface for a core config
type ConfigInterface interface {
	GetLogger() *zap.SugaredLogger
	GetTLSConfig() *kiloTLS.TLSConfig
	GetRefreshTimeout() time.Duration
	GetTLS() *tls.Config
}

// ReaderConfigInterface is the core.Reader config interface, embeds ConfigInterface
type ReaderConfigInterface interface {
	ConfigInterface
	GetBootstrapServers() []string
	GetTopic() string
	GetGroupID() string
	Workers() int
}

// WriterConfigInterface is the core.Writer config interface, embeds ConfigInterface
type WriterConfigInterface interface {
	ConfigInterface
	GetBrokers() []string
	GetTopic() string
	GetNickname() string
	GetAsync() bool
	GetEncoder() encoder.EncoderInterface
}

// AutoRefreshXConfigInterface is the refresh.AutoRefreshX config interface, embeds ConfigInterface
type AutoRefreshXConfigInterface interface {
	ConfigInterface
	GetMutexObj() *sync.RWMutex
}

// RefreshableInterface is implemented by core.Writer and core.Reader so that they can be refreshed by AutoRefreshX
type RefreshableInterface interface {
	Get(ctx context.Context) (interface{}, error)
	Config() interface{}
	Close()
	Stats() interface{}
	Renew(tlsCfg *kiloTLS.TLSConfig, args ...interface{}) RefreshableInterface
}
