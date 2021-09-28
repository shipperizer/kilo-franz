package core

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/shipperizer/kilo-franz/config"
	kiloTLS "github.com/shipperizer/kilo-franz/tls"
)

// getWriter is a helper function to create a kafka.Writer
func getWriter(cfg config.WriterConfigInterface) *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  cfg.GetBrokers(),
			Topic:    cfg.GetTopic(),
			Balancer: &kafka.LeastBytes{},
			Async:    cfg.GetAsync(),
			Dialer: &kafka.Dialer{
				DualStack: true,
				TLS:       cfg.GetTLS(),
			},
		},
	)
}

// Writer is an abstraction object on top of kakfa.Writer
// it holds the creation config as attribute and a pointer to the writer itself
// it implements RefreshableInterface so that can be used by AutoRefreshX
type Writer struct {
	writer *kafka.Writer
	cfg    config.WriterConfigInterface
}

// Renew creates a new Writer (wrapped in a RefreshableInterface) with the new tls config passed in
func (w *Writer) Renew(tlsCfg *kiloTLS.TLSConfig, args ...interface{}) config.RefreshableInterface {
	cfg := config.NewWriterConfig(
		config.NewConfig(w.cfg.GetRefreshTimeout(), tlsCfg, w.cfg.GetLogger()),
		w.cfg.GetBrokers(),
		w.cfg.GetTopic(),
		w.cfg.GetNickname(),
		w.cfg.GetAsync(),
		w.cfg.GetEncoder(),
	)

	return NewWriter(cfg)
}

// Stats returns a copy of kafka.ReaderStats (will need to be casted)
func (w *Writer) Stats() interface{} {
	return w.writer.Stats()
}

// Config returns the internal ReaderConfigInterface (will need to be casted)
func (w *Writer) Config() interface{} {
	return w.cfg
}

// Close makes sure the kafka.Reader.Close function is called
func (w *Writer) Close() {
	w.writer.Close() // flush all messages
}

// Get returns the internal reader object (will need to be casted) if present
func (w *Writer) Get(ctx context.Context) (interface{}, error) {
	if w.writer == nil {
		w.writer = getWriter(w.cfg)
	}

	return w.writer, nil
}

// NewWriter creates a new object implementing config.RefreshableInterface
func NewWriter(cfg config.WriterConfigInterface) config.RefreshableInterface {
	return &Writer{
		writer: getWriter(cfg),
		cfg:    cfg,
	}
}
