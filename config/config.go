package config

import (
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/shipperizer/kilo-franz/logging"
)

// Config is the core setup of any other config object, holds information about tls, refreshing time of
// TLS secrets and logger used for publisher/subscriber objects (using zap SugaredLogger)
type Config struct {
	tlsConfig      TLSConfigInterface
	refreshTimeout time.Duration
	logger         logging.LoggerInterface
	saslConfig     SASLConfigInterface
}

// GetTLSConfig returns a pointer to a TLSConfig
func (c *Config) GetTLSConfig() TLSConfigInterface {
	return c.tlsConfig
}

// GetSASLConfig returns a pointer to a SASLConfigInterface
func (c *Config) GetSASLConfig() SASLConfigInterface {
	return c.saslConfig
}

// GetRefreshTimeout returns the duration of a refreshing cycle, used by AutoRefreshX
func (c *Config) GetRefreshTimeout() time.Duration {
	return c.refreshTimeout
}

// GetLogger returns internal logger used for operations
func (c *Config) GetLogger() logging.LoggerInterface {
	return c.logger
}

// GetDialer returns internal dialer used for operations
func (c *Config) GetDialer() *kafka.Dialer {
	dialer := new(kafka.Dialer)

	dialer.DualStack = true
	dialer.TLS = c.getTLS()
	dialer.SASLMechanism = c.getSASL()

	return dialer
}

// getTLS returns a crypto/tls Config pointer made from the insatnce attribute c.tlsConfig
func (c *Config) getTLS() *tls.Config {
	if c.tlsConfig == nil {
		return nil
	}

	cfg, err := c.tlsConfig.GetTLS()

	if err != nil {
		c.logger.Errorf("issues fetching TLS config: %v", err)
		return nil
	}

	return cfg
}

// getSASL returns a crypto/tls Config pointer made from the insatnce attribute c.tlsConfig
func (c *Config) getSASL() sasl.Mechanism {
	if c.saslConfig == nil {
		return nil
	}

	return c.saslConfig.GetSASLMechanism()
}

// NewConfig returns an object implementing ConfigInterface
// func NewConfig(refreshTimeout time.Duration, tlsConfig TLSConfigInterface, logger logging.LoggerInterface) ConfigInterface {
func NewConfig(refreshTimeout time.Duration, tlsConfig TLSConfigInterface, saslConfig SASLConfigInterface, logger logging.LoggerInterface) *Config {
	cfg := new(Config)

	cfg.tlsConfig = tlsConfig
	cfg.refreshTimeout = refreshTimeout
	cfg.logger = logger
	cfg.saslConfig = saslConfig

	if logger == nil {
		cfg.logger = logging.NewLogger()
	}

	return cfg
}
