package testkit

import (
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/monitoring"
)

// FakeMonitor is a simple in-memory monitor used in tests.
type FakeMonitor struct {
	Service     string
	ServiceName string
	Logger      *FakeLogger
	Metrics     map[string]monitoring.MetricInterface
	MetricList  []monitoring.MetricInterface
}

// NewFakeMonitor creates a FakeMonitor with the given service name and logger.
func NewFakeMonitor(service string, logger *FakeLogger) *FakeMonitor {
	if logger == nil {
		logger = &FakeLogger{}
	}
	return &FakeMonitor{
		Service:     service,
		ServiceName: service,
		Logger:      logger,
		Metrics:     map[string]monitoring.MetricInterface{},
	}
}

func (m *FakeMonitor) ensureMetrics() {
	if m == nil {
		return
	}
	if m.Metrics == nil {
		m.Metrics = map[string]monitoring.MetricInterface{}
	}
}

func (m *FakeMonitor) GetMetric(name string) (monitoring.MetricInterface, error) {
	if m == nil {
		return NewFakeMetric(name), nil
	}
	m.ensureMetrics()
	if metric, ok := m.Metrics[name]; ok {
		return metric, nil
	}
	metric := NewFakeMetric(name)
	m.Metrics[name] = metric
	m.MetricList = append(m.MetricList, metric)
	return metric, nil
}

func (m *FakeMonitor) AddMetrics(metrics ...monitoring.MetricInterface) error {
	if m == nil {
		return nil
	}
	m.ensureMetrics()
	for _, metric := range metrics {
		if metric == nil {
			continue
		}
		m.Metrics[metric.Name()] = metric
		m.MetricList = append(m.MetricList, metric)
	}
	return nil
}

func (m *FakeMonitor) GetService() string {
	if m == nil {
		return ""
	}
	if m.Service != "" {
		return m.Service
	}
	return m.ServiceName
}

func (m *FakeMonitor) GetLogger() logging.LoggerInterface {
	if m == nil {
		return &FakeLogger{}
	}
	if m.Logger == nil {
		m.Logger = &FakeLogger{}
	}
	return m.Logger
}
