package models

import (
	"context"

	"github.com/shanghai-edu/vsphere-mon/funcs/core"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

func VcenterAlive(alive bool, vcAddr string) []*core.MetricValue {
	tags := map[string]string{}
	tags["vcenter"] = vcAddr
	if alive {
		return []*core.MetricValue{core.GaugeValue("", "vcenter.alive", 1, tags)}
	} else {
		return []*core.MetricValue{core.GaugeValue("", "vcenter.alive", 0, tags)}
	}
}

func VsphereDatastores(ctx context.Context, c *govmomi.Client) ([]mo.Datastore, error) {
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return nil, err
	}
	defer v.Destroy(ctx)

	var dss []mo.Datastore
	if err = v.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss); err != nil {
		return nil, err
	}

	if le := len(dss); le != 0 {
		return dss, nil
	}
	return nil, nil
}

//DatastoreMetrics datastore metrics
func DatastoreMetrics(ctx context.Context, c *govmomi.Client) (L []*core.MetricValue, err error) {
	dss, err := VsphereDatastores(ctx, c)
	if err != nil {
		return
	}
	if dss != nil {
		for _, ds := range dss {
			var tags = map[string]string{}
			tags["ds"] = ds.Summary.Name
			tags["fstype"] = ds.Summary.Type
			L = append(L, core.GaugeValue("", "datastore.bytes.total", ds.Summary.Capacity, tags))
			L = append(L, core.GaugeValue("", "datastore.bytes.free", ds.Summary.FreeSpace, tags))
			L = append(L, core.GaugeValue("", "datastore.bytes.used", ds.Summary.Capacity-ds.Summary.FreeSpace, tags))
			if ds.Summary.Capacity > 0 {
				usedPercent := float64(ds.Summary.Capacity-ds.Summary.FreeSpace) / float64(ds.Summary.Capacity) * 100
				L = append(L, core.GaugeValue("", "datastore.used.percent", usedPercent, tags))
			}

		}
	}
	return
}

func VsphereVirtualMachines(ctx context.Context, c *govmomi.Client) ([]mo.VirtualMachine, error) {
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return nil, err
	}
	defer v.Destroy(ctx)

	var vms []mo.VirtualMachine
	if err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms); err != nil {
		return nil, err
	}

	if le := len(vms); le != 0 {
		return vms, nil
	}
	return nil, nil
}

func GenVirtualMachinesByName(vms []mo.VirtualMachine, names []string) (resVms []mo.VirtualMachine) {
	if len(names) == 0 {
		resVms = vms
		return
	}
	for _, vm := range vms {
		if inSliceStr(vm.Summary.Config.Name, names) {
			resVms = append(resVms, vm)
		}
	}
	return
}

func inSliceStr(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}
	return false
}

//VmAlive power status
func VmPower(vm mo.VirtualMachine) []*core.MetricValue {
	/*
		1.0: poweredOff
		2.0: poweredOn
		3.0: standBy
		4.0: unknown
	*/
	var tags = map[string]string{}
	tags["name"] = vm.Summary.Config.Name

	switch vm.Summary.Runtime.PowerState {
	case "poweredOff":
		return []*core.MetricValue{core.GaugeValue("", "vm.power", 1.0, tags)}
	case "poweredOn":
		return []*core.MetricValue{core.GaugeValue("", "vm.power", 2.0, tags)}
	case "standBy":
		return []*core.MetricValue{core.GaugeValue("", "vm.power", 3.0, tags)}
	case "unknown":
		return []*core.MetricValue{core.GaugeValue("", "vm.power", 4.0, tags)}
	default:
		return []*core.MetricValue{core.GaugeValue("", "vm.power", 4.0, tags)}
	}
}

//VmStatus The Status enumeration defines a general "health" value for a managed entity.
func VmStatus(vm mo.VirtualMachine) []*core.MetricValue {
	/*
		1.0: gray,The status is unknown.
		2.0: green,The entity is OK.
		3.0: red,The entity definitely has a problem.
		4.0: yellow,The entity might have a problem.
	*/
	var tags = map[string]string{}
	tags["name"] = vm.Summary.Config.Name
	switch vm.Summary.OverallStatus {
	case "gray":
		return []*core.MetricValue{core.GaugeValue("", "vm.status", 1.0, tags)}
	case "green":
		return []*core.MetricValue{core.GaugeValue("", "vm.status", 2.0, tags)}
	case "red":
		return []*core.MetricValue{core.GaugeValue("", "vm.status", 3.0, tags)}
	case "yellow":
		return []*core.MetricValue{core.GaugeValue("", "vm.status", 4.0, tags)}
	default:
		return []*core.MetricValue{core.GaugeValue("", "vm.status", 1.0, tags)}
	}
}

//VmUptime uptime
func VmUptime(vm mo.VirtualMachine) []*core.MetricValue {
	var tags = map[string]string{}
	tags["name"] = vm.Summary.Config.Name
	return []*core.MetricValue{core.GaugeValue("", "vm.uptime", int64(vm.Summary.QuickStats.UptimeSeconds), tags)}
}

//VmCPU cpu metrics
func VmCPU(vm mo.VirtualMachine) []*core.MetricValue {
	var tags = map[string]string{}
	tags["name"] = vm.Summary.Config.Name
	perf := []*core.MetricValue{}
	if vm.Summary.Runtime.MaxCpuUsage > 0 {
		usePercentCPU := core.GaugeValue("", "cpu.util", float64(vm.Summary.QuickStats.OverallCpuUsage)/float64(vm.Summary.Runtime.MaxCpuUsage)*100, tags)
		freePercentCPU := core.GaugeValue("", "cpu.idle", (float64(vm.Summary.Runtime.MaxCpuUsage)-float64(vm.Summary.QuickStats.OverallCpuUsage))/float64(vm.Summary.Runtime.MaxCpuUsage)*100, tags)
		perf = append(perf, usePercentCPU)
		perf = append(perf, freePercentCPU)
	}
	return perf
}

//VmMem mem metrics
func VmMem(vm mo.VirtualMachine) []*core.MetricValue {
	var tags = map[string]string{}
	tags["name"] = vm.Summary.Config.Name
	totalMem := core.GaugeValue("", "mem.bytes.total", int64(vm.Summary.Runtime.MaxMemoryUsage)*1024*1024, tags)
	guestUseMem := core.GaugeValue("", "mem.bytes.guest.used", int64(vm.Summary.QuickStats.GuestMemoryUsage)*1024*1024, tags)
	hostUsedMem := core.GaugeValue("", "mem.bytes.host.used", int64(vm.Summary.QuickStats.HostMemoryUsage)*1024*1024, tags)
	perf := []*core.MetricValue{totalMem, guestUseMem, hostUsedMem}
	if vm.Summary.Runtime.MaxMemoryUsage > 0 {
		usedMemPer := core.GaugeValue("", "mem.bytes.guest.used.percent", float64(vm.Summary.QuickStats.GuestMemoryUsage)/float64(vm.Summary.Runtime.MaxMemoryUsage)*100, tags)
		perf = append(perf, usedMemPer)
	}
	return perf
}
