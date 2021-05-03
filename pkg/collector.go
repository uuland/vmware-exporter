package pkg

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

const (
	defaultTimeout = 5 * time.Second
)

func NewCollector(client *govmomi.Client, logger *log.Logger) *Collector {
	return &Collector{
		host:   client.URL().Hostname(),
		vmomi:  client,
		logger: logger,
	}
}

type Collector struct {
	host   string
	vmomi  *govmomi.Client
	logger *log.Logger

	hostView  *view.ContainerView
	storeView *view.ContainerView
	vmsView   *view.ContainerView
}

func (c *Collector) Start() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()

	var err error

	if c.hostView, err = view.NewManager(c.vmomi.Client).CreateContainerView(ctx,
		c.vmomi.ServiceContent.RootFolder, []string{"HostSystem"}, true); err != nil {
		return err
	}

	if c.storeView, err = view.NewManager(c.vmomi.Client).CreateContainerView(ctx,
		c.vmomi.ServiceContent.RootFolder, []string{"Datastore"}, true); err != nil {
		return err
	}

	if c.vmsView, err = view.NewManager(c.vmomi.Client).CreateContainerView(ctx,
		c.vmomi.ServiceContent.RootFolder, []string{"VirtualMachine"}, true); err != nil {
		return err
	}

	return nil
}

func (c *Collector) Stop() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()

	_ = c.hostView.Destroy(ctx)
	_ = c.storeView.Destroy(ctx)
	_ = c.vmsView.Destroy(ctx)

	return nil
}

func (c *Collector) Scrape() error {
	if err := c.hostMetrics(); err != nil {
		return err
	}

	if err := c.storeMetrics(); err != nil {
		return err
	}

	if err := c.vmsMetrics(); err != nil {
		return err
	}

	return nil
}

func (c *Collector) hostMetrics() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()

	var hss []mo.HostSystem
	if err := c.hostView.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss); err != nil {
		return err
	}

	for _, hs := range hss {
		hostPower.WithLabelValues(c.host).Set(powerState(hs.Summary.Runtime.PowerState))
		hostBootTime.WithLabelValues(c.host).Set(float64(hs.Summary.Runtime.BootTime.Unix()))
		hostCpuTotal.WithLabelValues(c.host).Set(totalCpu(hs))
		hostCpuUsage.WithLabelValues(c.host).Set(float64(hs.Summary.QuickStats.OverallCpuUsage))
		hostMemTotal.WithLabelValues(c.host).Set(float64(hs.Summary.Hardware.MemorySize))
		hostMemUsage.WithLabelValues(c.host).Set(float64(hs.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)
	}

	hs, err := find.NewFinder(c.vmomi.Client).DefaultHostSystem(ctx)
	if err != nil {
		return err
	}

	ss, err := hs.ConfigManager().StorageSystem(ctx)
	if err != nil {
		return err
	}

	var hostss mo.HostStorageSystem
	err = ss.Properties(ctx, ss.Reference(), nil, &hostss)
	if err != nil {
		return err
	}

	for _, e := range hostss.StorageDeviceInfo.ScsiLun {
		lun := e.GetScsiLun()
		ok := 1.0
		for _, s := range lun.OperationalState {
			if s != "ok" {
				ok = 0
				break
			}
		}
		hostDiskState.WithLabelValues(c.host, lun.DeviceName).Set(ok)
	}

	return nil
}

func (c *Collector) storeMetrics() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()

	var dss []mo.Datastore
	if err := c.storeView.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss); err != nil {
		return err
	}

	for _, ds := range dss {
		name := ds.Summary.Name
		storeCapacity.WithLabelValues(name, c.host).Set(float64(ds.Summary.Capacity))
		storeFreespace.WithLabelValues(name, c.host).Set(float64(ds.Summary.FreeSpace))
	}

	return nil
}

func (c *Collector) vmsMetrics() error {
	ctx, cancel := context.WithTimeout(context.TODO(), defaultTimeout)
	defer cancel()

	var vms []mo.VirtualMachine
	if err := c.vmsView.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms); err != nil {
		return err
	}

	for _, vm := range vms {
		name := vm.Summary.Config.Name
		vmBootTime.WithLabelValues(name, c.host).Set(convertTime(vm))
		vmCpuNum.WithLabelValues(name, c.host).Set(float64(vm.Summary.Config.NumCpu))
		vmCpuTotal.WithLabelValues(name, c.host).Set(float64(vm.Summary.Runtime.MaxCpuUsage) * 1000 * 1000)
		vmCpuUsage.WithLabelValues(name, c.host).Set(float64(vm.Summary.QuickStats.OverallCpuUsage) * 1000 * 1000)
		vmMemTotal.WithLabelValues(name, c.host).Set(float64(vm.Summary.Config.MemorySizeMB))
		vmMemUsage.WithLabelValues(name, c.host).Set(float64(vm.Summary.QuickStats.GuestMemoryUsage) * 1024 * 1024)
	}

	return nil
}
