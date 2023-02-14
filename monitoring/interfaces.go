package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shipperizer/kilo-franz/logging"
)

type MonitorInterface interface {
	GetMetric(string) (MetricInterface, error)
	AddMetrics(...MetricInterface) error
	GetService() string
	GetLogger() logging.LoggerInterface
}

type MetricInterface interface {
	Inc(map[string]string) error
	Dec(map[string]string) error
	Add(float64, map[string]string) error
	Sub(float64, map[string]string) error
	Set(float64, map[string]string) error
	Observe(float64, map[string]string) error
	Collector() prometheus.Collector
	Name() string
}
