package core

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/shipperizer/kilo-franz/config"
)

// getWriter is a helper function to create a kafka.Writer
func getWriter(cfg WriterConfigInterface) *kafka.Writer {
	return kafka.NewWriter(
		kafka.WriterConfig{
			Brokers:  cfg.GetBrokers(),
			Topic:    cfg.GetTopic(),
			Balancer: &kafka.LeastBytes{},
			Async:    cfg.GetAsync(),
			Dialer:   cfg.GetDialer(),
		},
	)
}

// Writer is an abstraction object on top of kakfa.Writer
// it holds the creation config as attribute and a pointer to the writer itself
// it implements RefreshableInterface so that can be used by AutoRefreshX
type Writer struct {
	writer *kafka.Writer
	cfg    WriterConfigInterface
}

// Renew creates a new kafka.Writer with the new tls config passed in and updates the instance
func (w *Writer) Renew(tlsConfig config.TLSConfigInterface, saslConfig config.SASLConfigInterface, args ...interface{}) {
	cfg := config.NewWriterConfig(
		config.NewConfig(w.cfg.GetRefreshTimeout(), tlsConfig, saslConfig, w.cfg.GetLogger()),
		w.cfg.GetBrokers(),
		w.cfg.GetTopic(),
		w.cfg.GetNickname(),
		w.cfg.GetAsync(),
		w.cfg.GetEncoder(),
	)

	w.writer = getWriter(cfg)
	w.cfg = cfg
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

// NewWriter creates a new object
func NewWriter(cfg WriterConfigInterface) *Writer {
	w := new(Writer)

	w.writer = getWriter(cfg)
	w.cfg = cfg

	return w
}
