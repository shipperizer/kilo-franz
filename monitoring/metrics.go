package monitoring

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	GAUGE = iota
	COUNTER
	HISTOGRAM
)

type MetricClass int

type Metric struct {
	counter   *prometheus.CounterVec
	gauge     *prometheus.GaugeVec
	histogram *prometheus.HistogramVec

	name string
}

func (m *Metric) Name() string {
	return m.name
}

func (m *Metric) Collector() prometheus.Collector {
	if m.counter != nil {
		return m.counter
	}

	if m.gauge != nil {
		return m.gauge
	}

	if m.histogram != nil {
		return m.histogram
	}

	return nil
}

func (m *Metric) Inc(tags map[string]string) error {
	if m.counter != nil {
		c, err := m.counter.GetMetricWith(tags)

		if err != nil {
			return err
		}

		c.Inc()
	} else if m.gauge != nil {
		g, err := m.gauge.GetMetricWith(tags)

		if err != nil {
			return err
		}

		g.Inc()
	} else if m.histogram != nil {
		return fmt.Errorf("metric is not a of counter or gauge type")
	}

	return nil
}

func (m *Metric) Dec(tags map[string]string) error {
	if m.gauge == nil {
		return fmt.Errorf("metric is not a of gauge type")
	}

	g, err := m.gauge.GetMetricWith(tags)

	if err != nil {
		return err
	}

	g.Dec()

	return nil
}

func (m *Metric) Add(value float64, tags map[string]string) error {
	if m.counter != nil {
		c, err := m.counter.GetMetricWith(tags)

		if err != nil {
			return err
		}

		c.Add(value)
	} else if m.gauge != nil {
		g, err := m.gauge.GetMetricWith(tags)

		if err != nil {
			return err
		}

		g.Add(value)
	} else if m.histogram != nil {
		return fmt.Errorf("metric is not a of counter or gauge type")
	}

	return nil
}

func (m *Metric) Sub(value float64, tags map[string]string) error {
	if m.gauge == nil {
		return fmt.Errorf("metric is not a of gauge type")
	}

	g, err := m.gauge.GetMetricWith(tags)

	if err != nil {
		return err
	}

	g.Sub(value)

	return nil
}
func (m *Metric) Set(value float64, tags map[string]string) error {
	if m.gauge == nil {
		return fmt.Errorf("metric is not a of gauge type")
	}

	g, err := m.gauge.GetMetricWith(tags)

	if err != nil {
		return err
	}

	g.Set(value)

	return nil
}

func (m *Metric) Observe(value float64, tags map[string]string) error {
	if m.histogram == nil {
		return fmt.Errorf("metric is not an histogram")
	}

	h, err := m.histogram.GetMetricWith(tags)

	if err != nil {
		return err
	}

	h.Observe(value)

	return nil
}

// NewMetric creates a new Metric object, class needs to be one of type MetricClass, like the following
// * GAUGE = 0
// * COUNTER = 1
// * HISTOGRAM = 2
func NewMetric(class MetricClass, name string, labels ...string) *Metric {
	m := new(Metric)

	m.name = name

	if m.name == "" {
		return nil
	}

	switch class {
	case GAUGE:
		m.gauge = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: fmt.Sprintf("LABS metric: %s", name),
			},
			labels,
		)
	case HISTOGRAM:
		m.histogram = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: name,
				Help: fmt.Sprintf("LABS metric: %s", name),
			},
			labels,
		)
	case COUNTER:
		m.counter = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: name,
				Help: fmt.Sprintf("LABS metric: %s", name),
			},
			labels,
		)
	}

	return m

}
