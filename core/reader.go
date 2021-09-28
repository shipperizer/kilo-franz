package core

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/shipperizer/kilo-franz/config"
	kiloTLS "github.com/shipperizer/kilo-franz/tls"
)

// getReader is a helper function to create a kafka.Reader
func getReader(cfg config.ReaderConfigInterface) *kafka.Reader {
	c := kafka.ReaderConfig{
		Brokers:  cfg.GetBootstrapServers(),
		GroupID:  cfg.GetGroupID(),
		Topic:    cfg.GetTopic(),
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB,
		Dialer: &kafka.Dialer{
			DualStack: true,
			TLS:       cfg.GetTLS(),
		},
	}

	return kafka.NewReader(c)
}

// Reader is an abstraction object on top of kakfa.Reader
// it holds the creation config as attribute and a pointer to the reader itself
// it implements RefreshableInterface so that can be used by AutoRefreshX
type Reader struct {
	reader *kafka.Reader
	cfg    config.ReaderConfigInterface
}

// Renew creates a new Reader (wrapped in a RefreshableInterface) with the new tls config passed in
func (r *Reader) Renew(tlsCfg *kiloTLS.TLSConfig, args ...interface{}) config.RefreshableInterface {
	cfg := config.NewReaderConfig(
		config.NewConfig(r.cfg.GetRefreshTimeout(), tlsCfg, r.cfg.GetLogger()),
		r.cfg.GetBootstrapServers(),
		r.cfg.GetTopic(),
		r.cfg.GetGroupID(),
		r.cfg.Workers(),
		r.cfg.GetRefreshTimeout(),
	)

	return NewReader(cfg)
}

// Stats returns a copy of kafka.ReaderStats (will need to be casted)
func (r *Reader) Stats() interface{} {
	return r.reader.Stats()
}

// Config returns the internal ReaderConfigInterface (will need to be casted)
func (r *Reader) Config() interface{} {
	return r.cfg
}

// Close makes sure the kafka.Reader.Close function is called
func (r *Reader) Close() {
	r.reader.Close()
}

// Get returns the internal reader object (will need to be casted) if present
func (r *Reader) Get(ctx context.Context) (interface{}, error) {
	if r.reader == nil {
		r.reader = getReader(r.cfg)
	}

	return r.reader, nil
}

// NewReader creates a new object implementing config.RefreshableInterface
func NewReader(cfg config.ReaderConfigInterface) config.RefreshableInterface {
	return &Reader{
		reader: getReader(cfg),
		cfg:    cfg,
	}
}
