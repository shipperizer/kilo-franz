package testkit

import (
	"testing"

	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/monitoring"
)

func TestFakesSatisfyInterfaces(t *testing.T) {
	var _ logging.LoggerInterface = (*FakeLogger)(nil)
	var _ monitoring.MetricInterface = (*FakeMetric)(nil)
	var _ monitoring.MonitorInterface = (*FakeMonitor)(nil)
}
