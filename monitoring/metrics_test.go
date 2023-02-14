package monitoring

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestNewMetricImplementsMetricInterface(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	assert := assert.New(t)
	assert.Implements((*MetricInterface)(nil), NewMetric(GAUGE, "test"))
}

func TestMetricName(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("test", NewMetric(GAUGE, "test").Name())
}

func TestMetricCollector(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(NewMetric(GAUGE, "test").Collector())
	assert.NotNil(NewMetric(HISTOGRAM, "test").Collector())
	assert.NotNil(NewMetric(COUNTER, "test").Collector())
}

func TestMetricCollectorWhenMetricIsEmpty(t *testing.T) {
	assert := assert.New(t)
	assert.Nil(NewMetric(10, "test").Collector())
}

func TestMetricIncOnInvalidMetricType(t *testing.T) {
	metric := NewMetric(HISTOGRAM, "test")

	assert := assert.New(t)

	assert.NotNil(metric.Inc(nil))
}

func TestMetricIncOnValidMetricType(t *testing.T) {
	g := NewMetric(GAUGE, "test")
	c := NewMetric(COUNTER, "test")

	assert := assert.New(t)

	assert.Nil(g.Inc(nil))
	assert.Nil(c.Inc(nil))

	m := &dto.Metric{}

	gauge, _ := g.Collector().(*prometheus.GaugeVec).GetMetricWith(nil)
	gauge.Write(m)
	assert.Equal(float64(1), m.GetGauge().GetValue())

	counter, _ := c.Collector().(*prometheus.CounterVec).GetMetricWith(nil)
	counter.Write(m)
	assert.Equal(float64(1), m.GetCounter().GetValue())
}

func TestMetricDecOnInvalidMetricType(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewMetric(HISTOGRAM, "test").Dec(nil))
	assert.NotNil(NewMetric(COUNTER, "test").Dec(nil))
}

func TestMetricDecOnValidMetricType(t *testing.T) {
	g := NewMetric(GAUGE, "test")

	assert := assert.New(t)

	assert.Nil(g.Dec(nil))

	m := &dto.Metric{}

	gauge, _ := g.Collector().(*prometheus.GaugeVec).GetMetricWith(nil)
	gauge.Write(m)
	assert.Equal(float64(-1), m.GetGauge().GetValue())
}

func TestMetricAddOnInvalidMetricType(t *testing.T) {
	metric := NewMetric(HISTOGRAM, "test")

	assert := assert.New(t)

	assert.NotNil(metric.Add(float64(5), nil))
}

func TestMetricAddOnValidMetricType(t *testing.T) {
	g := NewMetric(GAUGE, "test")
	c := NewMetric(COUNTER, "test")

	assert := assert.New(t)

	assert.Nil(g.Add(float64(5), nil))
	assert.Nil(c.Add(float64(5), nil))

	m := &dto.Metric{}

	gauge, _ := g.Collector().(*prometheus.GaugeVec).GetMetricWith(nil)
	gauge.Write(m)
	assert.Equal(float64(5), m.GetGauge().GetValue())

	counter, _ := c.Collector().(*prometheus.CounterVec).GetMetricWith(nil)
	counter.Write(m)
	assert.Equal(float64(5), m.GetCounter().GetValue())
}

func TestMetricSubOnInvalidMetricType(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewMetric(HISTOGRAM, "test").Sub(float64(1), nil))
	assert.NotNil(NewMetric(COUNTER, "test").Sub(float64(1), nil))
}

func TestMetricSubOnValidMetricType(t *testing.T) {
	g := NewMetric(GAUGE, "test")

	assert := assert.New(t)

	assert.Nil(g.Sub(float64(5), nil))

	m := &dto.Metric{}

	gauge, _ := g.Collector().(*prometheus.GaugeVec).GetMetricWith(nil)
	gauge.Write(m)
	assert.Equal(float64(-5), m.GetGauge().GetValue())

}

func TestMetricSetOnInvalidMetricType(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewMetric(HISTOGRAM, "test").Set(float64(1), nil))
	assert.NotNil(NewMetric(COUNTER, "test").Set(float64(1), nil))
}

func TestMetricSetOnValidMetricType(t *testing.T) {
	g := NewMetric(GAUGE, "test")

	assert := assert.New(t)

	assert.Nil(g.Set(float64(100), nil))

	m := &dto.Metric{}

	gauge, _ := g.Collector().(*prometheus.GaugeVec).GetMetricWith(nil)
	gauge.Write(m)
	assert.Equal(float64(100), m.GetGauge().GetValue())

}

func TestMetricObserveOnInvalidMetricType(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(NewMetric(GAUGE, "test").Observe(float64(1), nil))
	assert.NotNil(NewMetric(COUNTER, "test").Observe(float64(1), nil))
}

func TestMetricObserveOnValidMetricType(t *testing.T) {
	h := NewMetric(HISTOGRAM, "test")

	assert := assert.New(t)

	assert.Nil(h.Observe(float64(100), nil))

	m := &dto.Metric{}

	histogram, _ := h.Collector().(*prometheus.HistogramVec).GetMetricWith(nil)
	histogram.(prometheus.Metric).Write(m)
	assert.Equal(float64(100), m.GetHistogram().GetSampleSum())

}
