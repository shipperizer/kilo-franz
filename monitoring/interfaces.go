package monitoring

import "github.com/prometheus/client_golang/prometheus"

type MonitorInterface interface {
	GetService() string
	GetMetric(metric string) *prometheus.GaugeVec
	Incr(metric string, opts map[string]string)
	Decr(metric string, opts map[string]string)
	Gauge(metric string, value int64, opts map[string]string)
}
