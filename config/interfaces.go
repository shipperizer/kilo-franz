package config

import (
	"crypto/tls"
	"time"

	"github.com/segmentio/kafka-go/sasl"

	"github.com/segmentio/kafka-go"
	"github.com/shipperizer/kilo-franz/logging"
)

// ConfigInterface is the interface for a core config
type ConfigInterface interface {
	GetLogger() logging.LoggerInterface
	GetTLSConfig() TLSConfigInterface
	GetRefreshTimeout() time.Duration
	GetSASLConfig() SASLConfigInterface
	GetDialer() *kafka.Dialer
}

type TLSConfigInterface interface {
	GetTLS() (*tls.Config, error)
}

// SASLConfigInterface is the interface for a core config
type SASLConfigInterface interface {
	GetSASLMechanism() sasl.Mechanism
}
