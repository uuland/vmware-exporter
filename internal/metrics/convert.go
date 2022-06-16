package metrics

import "vmware-exporter/internal/collector"

func bootTime(m *collector.Metrics) {
}

func cpuTotalMhz(m *collector.Metrics) {

}

func toMBytes(m *collector.Metrics) {

}

func prop(name string) collector.MetricsOption {
	return func(m *collector.Metrics) {
		m.Prop(name)
	}
}
