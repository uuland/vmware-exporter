package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "vmware"
)

var (
	// host metrics
	hostPower = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "power_state",
		Help:      "poweredOn 1, poweredOff 2, standBy 3, other 0",
	}, []string{"host_name"})
	hostBootTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "boot_timestamp",
		Help:      "Uptime host",
	}, []string{"host_name"})
	hostCpuTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "cpu_max",
		Help:      "CPU total",
	}, []string{"host_name"})
	hostCpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "cpu_usage",
		Help:      "CPU Usage",
	}, []string{"host_name"})
	hostMemTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "memory_max",
		Help:      "Memory max",
	}, []string{"host_name"})
	hostMemUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "memory_usage",
		Help:      "Memory Usage",
	}, []string{"host_name"})
	hostDiskState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "host",
		Name:      "disk_ok",
		Help:      "Disk is working normally",
	}, []string{"host_name", "device"})

	// datastore metrics
	storeCapacity = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "datastore",
		Name:      "capacity_size",
		Help:      "Datastore total",
	}, []string{"ds_name", "host_name"})
	storeFreespace = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "datastore",
		Name:      "freespace_size",
		Help:      "Datastore free",
	}, []string{"ds_name", "host_name"})

	// virtual machine metrics
	vmBootTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "boot_timestamp",
		Help:      "VMWare VM boot time in seconds",
	}, []string{"vm_name", "host_name"})
	vmCpuTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "cpu_available_mhz",
		Help:      "VMWare VM total CPU",
	}, []string{"vm_name", "host_name"})
	vmCpuUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "cpu_usage_mhz",
		Help:      "VMWare VM usage CPU",
	}, []string{"vm_name", "host_name"})
	vmCpuNum = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "num_cpu",
		Help:      "Available number of cores",
	}, []string{"vm_name", "host_name"})
	vmMemTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "mem_available",
		Help:      "Available memory",
	}, []string{"vm_name", "host_name"})
	vmMemUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "mem_usage",
		Help:      "Usage memory",
	}, []string{"vm_name", "host_name"})
	vmNetRX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "net_rx",
		Help:      "Network RX bytes",
	}, []string{"vm_name", "host_name", "interface"})
	vmNetTX = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "vm",
		Name:      "net_tx",
		Help:      "Network TX bytes",
	}, []string{"vm_name", "host_name", "interface"})
)

func init() {
	prometheus.MustRegister(
		hostPower,
		hostBootTime,
		hostCpuTotal,
		hostCpuUsage,
		hostMemTotal,
		hostMemUsage,
		hostDiskState,
		storeCapacity,
		storeFreespace,
		vmBootTime,
		vmCpuTotal,
		vmCpuNum,
		vmMemTotal,
		vmMemUsage,
		vmCpuUsage,
		vmNetRX,
		vmNetTX,
	)
}
