package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/performance"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/types"

	"vmware-exporter/internal/collector"
)

func init() {
	collector.Registry.Add("host_perf", newHostPerf)
}

func newHostPerf() collector.Collector {
	m := collector.NewMetrics("host_perf")
	m.WithLabels("host_name")

	//cpu := m.Subsystem("cpu").Gauge()
	//mem := m.Subsystem("mem").Gauge()
	//disk := m.Subsystem("disk").Gauge()
	//net := m.Subsystem("net").Gauge()

	cs := []prometheus.Collector{
		//cpu,
		//mem,
		//disk,
		//net,
	}

	return &hostPerf{metrics: m, collectors: cs}
}

type hostPerf struct {
	metrics    *collector.Metrics
	collectors []prometheus.Collector
}

func (c *hostPerf) Startup() error {
	return nil
}

func (c *hostPerf) Shutdown() error {
	return nil
}

func (c *hostPerf) Collectors() []prometheus.Collector {
	return c.collectors
}

func (c *hostPerf) Scrape(cc *collector.Context) error {
	ctx := cc.Context
	cli := cc.Client.Client

	// Get virtual machines references
	m := view.NewManager(cli)

	v, err := m.CreateContainerView(ctx, cli.ServiceContent.RootFolder, nil, true)
	if err != nil {
		return err
	}

	defer v.Destroy(ctx)

	vmsRefs, err := v.Find(ctx, []string{"VirtualMachine"}, nil)
	if err != nil {
		return err
	}

	// Create a PerfManager
	perfManager := performance.NewManager(cli)

	// Retrieve counters name list
	counters, err := perfManager.CounterInfoByName(ctx)
	if err != nil {
		return err
	}

	var names []string
	for name := range counters {
		names = append(names, name)
	}

	// Create PerfQuerySpec
	spec := types.PerfQuerySpec{
		MaxSample:  1,
		MetricId:   []types.PerfMetricId{{Instance: "*"}},
		IntervalId: 20,
	}

	// Query metrics
	sample, err := perfManager.SampleByName(ctx, spec, names, vmsRefs)
	if err != nil {
		return err
	}

	result, err := perfManager.ToMetricSeries(ctx, sample)
	if err != nil {
		return err
	}

	// Read result
	for _, metric := range result {
		name := metric.Entity

		for _, v := range metric.Value {
			counter := counters[v.Name]
			units := counter.UnitInfo.GetElementDescription().Label

			instance := v.Instance
			if instance == "" {
				instance = "-"
			}

			if len(v.Value) != 0 {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n",
					name, instance, v.Name, v.ValueCSV(), units)
			}
		}
	}
	return nil
}
