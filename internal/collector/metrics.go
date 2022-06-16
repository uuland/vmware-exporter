package collector

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	cTypeGauge = "gauge"
)

func NewMetrics(subsystem string) *Metrics {
	return &Metrics{
		namespace: "vmware",
		subsystem: subsystem,
		shared: &metricsShared{
			propSet: make(map[string]valueMetric),
		},
		converter: toFloat64,
	}
}

func toFloat64(raw interface{}) float64 {
	switch v := raw.(type) {
	case float64:
		return v
	default:
		fmt.Println("toFloat64: raw type = ", fmt.Sprintf("%T", raw))
		return 0
	}
}

type Metrics struct {
	namespace string
	subsystem string
	labels    []string

	shared *metricsShared // shared data

	vecType   string // gauge
	property  string // vmomi property bind
	converter valueConvert
}

type valueConvert func(raw interface{}) float64

type valueMetric interface {
	WithLabelValues(...string) prometheus.Gauge
}

type valueSet interface {
	Set(float64)
}

type metricsShared struct {
	props   []string               // used property keys
	propSet map[string]valueMetric // prop <=> vec
}

type MetricsOption func(m *Metrics)

func (m *Metrics) WithLabels(keys ...string) *Metrics {
	m.labels = append(m.labels, keys...)
	return m
}

func (m *Metrics) Subsystem(name string) *Metrics {
	m2 := &Metrics{}
	*m2 = *m
	m2.subsystem += "_" + name
	return m2
}

func (m *Metrics) Gauge() *Metrics {
	m.vecType = cTypeGauge
	return m
}

func (m *Metrics) Prop(prop string) *Metrics {
	m.property = prop
	return m
}

func (m *Metrics) Props() []string {
	return m.shared.props
}

func (m *Metrics) Converter(vc valueConvert) *Metrics {
	m.converter = vc
	return m
}

func (m *Metrics) Build(name string, help string, opts ...MetricsOption) prometheus.Collector {
	for _, opt := range opts {
		opt(m)
	}

	pOpts := prometheus.Opts{
		Namespace: m.namespace,
		Subsystem: m.subsystem,
		Name:      name,
		Help:      help,
	}

	var c prometheus.Collector

	switch m.vecType {
	case cTypeGauge:
		c = prometheus.NewGaugeVec(prometheus.GaugeOpts(pOpts), m.labels)
	}

	if c == nil {
		panic("unknown collector type")
	}

	if m.property != "" {
		m.shared.props = append(m.shared.props, m.property)
		m.shared.propSet[m.property] = c.(valueMetric)
		m.property = ""
	}

	return c
}

func (m *Metrics) Process(dps []types.DynamicProperty, labels ...string) {
	for _, p := range dps {
		if vec, has := m.shared.propSet[p.Name]; has {
			vec.WithLabelValues(labels...).Set(m.converter(p.Val))
		}
	}
}
