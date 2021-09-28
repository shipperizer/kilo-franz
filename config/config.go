package config

import (
	"crypto/tls"
	"time"

	"go.uber.org/zap"

	kiloTLS "github.com/shipperizer/kilo-franz/tls"
)

// Config is the core setup of any other config object, holds information about tls, refreshing time of
// TLS secrets and logger used for publisher/subscriber objects (using zap SugaredLogger)
type Config struct {
	tlsConfig      *kiloTLS.TLSConfig
	refreshTimeout time.Duration
	logger         *zap.SugaredLogger
}

// GetTLSConfig returns a pointer to a TLSConfig
func (c *Config) GetTLSConfig() *kiloTLS.TLSConfig {
	return c.tlsConfig
}

// GetTLS returns a crypto/tls Config pointer made from the insatnce attribute c.tlsConfig
func (c *Config) GetTLS() *tls.Config {
	if c.tlsConfig != nil {

		cfg, err := kiloTLS.GetTLS(*c.tlsConfig)

		if err != nil {
			c.logger.Errorf("issues fetching TLS config: %v", err)
			return nil
		}

		return cfg
	}

	return nil
}

// GetRefreshTimeout returns the duration of a refreshing cycle, used by AutoRefreshX
func (c *Config) GetRefreshTimeout() time.Duration {
	return c.refreshTimeout
}

// GetLogger returns internal logger used for operations
func (c *Config) GetLogger() *zap.SugaredLogger {
	return c.logger
}

// NewConfig returns an object implementing ConfigInterface
func NewConfig(refreshTimeout time.Duration, tlsConfig *kiloTLS.TLSConfig, logger *zap.SugaredLogger) ConfigInterface {
	cfg := &Config{
		tlsConfig:      tlsConfig,
		refreshTimeout: refreshTimeout,
		logger:         logger,
	}

	if logger == nil {
		cfg.logger = NewLogger()
	}

	return cfg
}
