package testkit

import (
	"github.com/prometheus/client_golang/prometheus"
)

// FakeMetric is a simple in-memory metric used in tests.
type FakeMetric struct {
	MetricName string
	NameValue  string
	Value      float64
	Values     []float64
	LastValue  float64
	Tags       map[string]string
	LastTags   map[string]string
	Calls      []string
}

// NewFakeMetric creates a FakeMetric with the given name.
func NewFakeMetric(name string) *FakeMetric {
	return &FakeMetric{MetricName: name, NameValue: name}
}

func (m *FakeMetric) name() string {
	if m == nil {
		return ""
	}
	if m.MetricName != "" {
		return m.MetricName
	}
	return m.NameValue
}

func (m *FakeMetric) Inc(tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.Value++
	m.LastValue = m.Value
	m.Values = append(m.Values, m.Value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Inc")
	return nil
}

func (m *FakeMetric) Dec(tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.Value--
	m.LastValue = m.Value
	m.Values = append(m.Values, m.Value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Dec")
	return nil
}

func (m *FakeMetric) Add(value float64, tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.Value += value
	m.LastValue = value
	m.Values = append(m.Values, value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Add")
	return nil
}

func (m *FakeMetric) Sub(value float64, tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.Value -= value
	m.LastValue = value
	m.Values = append(m.Values, value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Sub")
	return nil
}

func (m *FakeMetric) Set(value float64, tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.Value = value
	m.LastValue = value
	m.Values = append(m.Values, value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Set")
	return nil
}

func (m *FakeMetric) Observe(value float64, tags map[string]string) error {
	if m == nil {
		return nil
	}
	m.LastValue = value
	m.Values = append(m.Values, value)
	m.Tags = tags
	m.LastTags = tags
	m.Calls = append(m.Calls, "Observe")
	return nil
}

func (m *FakeMetric) Collector() prometheus.Collector {
	return nil
}

func (m *FakeMetric) Name() string {
	if m == nil {
		return ""
	}
	return m.name()
}
