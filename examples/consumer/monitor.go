package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/monitoring"
)

// Monitor implements monitoring.MonitorInterface for the examples
type Monitor struct {
	logger  logging.LoggerInterface
	metrics map[string]monitoring.MetricInterface
	service string
}

// NewMonitor creates a new Monitor instance
func NewMonitor(service string, logger logging.LoggerInterface) *Monitor {
	return &Monitor{
		logger:  logger,
		metrics: make(map[string]monitoring.MetricInterface),
		service: service,
	}
}

// GetMetric returns a metric by name
func (m *Monitor) GetMetric(name string) (monitoring.MetricInterface, error) {
	if metric, ok := m.metrics[name]; ok {
		return metric, nil
	}
	return nil, fmt.Errorf("metric %s not found", name)
}

// AddMetrics registers metrics
func (m *Monitor) AddMetrics(metrics ...monitoring.MetricInterface) error {
	for _, metric := range metrics {
		m.metrics[metric.Name()] = metric
		if err := prometheus.Register(metric.Collector()); err != nil {
			m.logger.Debugf("metric %s already registered: %v", metric.Name(), err)
		}
	}
	return nil
}

// GetService returns the service name
func (m *Monitor) GetService() string {
	return m.service
}

// GetLogger returns the logger
func (m *Monitor) GetLogger() logging.LoggerInterface {
	return m.logger
}
