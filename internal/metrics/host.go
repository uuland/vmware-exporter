package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/types"

	"vmware-exporter/internal/collector"
)

func init() {
	collector.Registry.Add("host", newHost)
}

func newHost() collector.Collector {
	m := collector.NewMetrics("host")
	m.WithLabels("host_name")

	boot := m.Subsystem("boot").Gauge()
	cpu := m.Subsystem("cpu").Gauge()
	mem := m.Subsystem("mem").Gauge()

	cs := []prometheus.Collector{
		// boot
		boot.Build("seconds", "VMWare Host boot time in seconds", prop("runtime.bootTime"), bootTime),
		// cpu
		cpu.Build("usage", "VMWare Host CPU usage in Mhz", prop("summary.quickStats.overallCpuUsage")),
		cpu.Build("max", "VMWare Host CPU max availability in Mhz", prop("summary.hardware.cpuMhz"), cpuTotalMhz),
		cpu.Build("num", "VMWare Number of processors in the Host", prop("summary.hardware.numCpuCores")),
		// memory
		mem.Build("usage", "VMWare Host Memory usage in Mbytes", prop("summary.quickStats.overallMemoryUsage")),
		mem.Build("max", "VMWare Host Memory Max availability in Mbytes", prop("summary.hardware.memorySize"), toMBytes),
	}

	return &host{metrics: m, collectors: cs}
}

type host struct {
	metrics    *collector.Metrics
	collectors []prometheus.Collector
}

func (c *host) Startup() error {
	return nil
}

func (c *host) Shutdown() error {
	return nil
}

func (c *host) Collectors() []prometheus.Collector {
	return c.collectors
}

func (c *host) Scrape(ctx *collector.Context) error {
	v, err := view.NewManager(ctx.Client.Client).CreateContainerView(ctx.Context, ctx.Client.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return err
	}
	defer v.Destroy(ctx.Context)

	var ret []types.ObjectContent
	if err := v.Retrieve(ctx.Context, []string{"HostSystem"}, c.metrics.Props(), &ret); err != nil {
		return err
	}

	for _, r := range ret {
		c.metrics.Process(r.PropSet, "dc1")
	}

	return nil
}
