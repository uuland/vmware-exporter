package collector_test

import (
	"regexp"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"vmware-exporter/internal/collector"
)

func TestMetrics(t *testing.T) {
	c := make(chan prometheus.Collector)
	d := make(chan *prometheus.Desc)
	w := make(chan struct{})

	rName, err := regexp.Compile("fqName: \"(.*?)\"")
	if err != nil {
		t.Fatal(err)
	}

	var ks []string
	addKey := func(matches []string) {
		if len(matches) > 1 {
			ks = append(ks, matches[1])
		}
	}

	go func() {
		defer func() {
			close(w)
		}()
		for d2 := range d {
			addKey(rName.FindStringSubmatch(d2.String()))
		}
	}()

	go func() {
		defer func() {
			close(d)
		}()
		for c2 := range c {
			c2.Describe(d)
		}
	}()

	p1 := "summary.quickStats.overallCpuUsage"
	p2 := "summary.hardware.cpuMhz"
	p3 := "summary.hardware.numCpuCores"

	k1 := "vmware_host_cpu_usage"
	k2 := "vmware_host_cpu_max"
	k3 := "vmware_host_cpu_num"

	host := collector.NewMetrics("host")
	host.WithLabels("host_name")
	host.Gauge()

	cpu := host.Subsystem("cpu")
	c <- cpu.Prop("summary.quickStats.overallCpuUsage").Build("usage", "VMWare Host CPU usage in Mhz")
	c <- cpu.Prop("summary.hardware.cpuMhz").Build("max", "VMWare Host CPU max availability in Mhz")
	c <- cpu.Prop("summary.hardware.numCpuCores").Build("num", "VMWare Number of processors in the Host")

	close(c)
	<-w

	assert.Equal(t, []string{p1, p2, p3}, host.Props())
	assert.Equal(t, []string{k1, k2, k3}, ks)
}
