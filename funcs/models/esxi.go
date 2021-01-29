package models

import (
	"context"

	"github.com/shanghai-edu/vsphere-mon/funcs/core"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

//EsxiList get esxi list
func EsxiList(ctx context.Context, c *govmomi.Client) (esxiList []mo.HostSystem, err error) {
	m := view.NewManager(c.Client)
	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
	if err != nil {
		return
	}
	defer v.Destroy(ctx)

	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary", "datastore"}, &esxiList)
	return
}

//EsxiAlive power status
func EsxiPower(esxi mo.HostSystem) []*core.MetricValue {
	/*
		1.0: poweredOff
		2.0: poweredOn
		3.0: standBy
		4.0: unknown
	*/
	switch esxi.Summary.Runtime.PowerState {
	case "poweredOff":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.power", 1.0, nil)}
	case "poweredOn":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.power", 2.0, nil)}
	case "standBy":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.power", 3.0, nil)}
	case "unknown":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.power", 4.0, nil)}
	default:
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.power", 4.0, nil)}
	}
}

//EsxiStatus The Status enumeration defines a general "health" value for a managed entity.
func EsxiStatus(esxi mo.HostSystem) []*core.MetricValue {
	/*
		1.0: gray,The status is unknown.
		2.0: green,The entity is OK.
		3.0: red,The entity definitely has a problem.
		4.0: yellow,The entity might have a problem.
	*/
	switch esxi.Summary.OverallStatus {
	case "gray":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.status", 1.0, nil)}
	case "green":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.status", 2.0, nil)}
	case "red":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.status", 3.0, nil)}
	case "yellow":
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.status", 4.0, nil)}
	default:
		return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.status", 1.0, nil)}
	}
}

//EsxiUptime uptime
func EsxiUptime(esxi mo.HostSystem) []*core.MetricValue {
	return []*core.MetricValue{core.GaugeValue(esxi.Summary.Config.Name, "esxi.uptime", int64(esxi.Summary.QuickStats.Uptime), nil)}
}

//EsxiCPU cpu metrics
func EsxiCPU(esxi mo.HostSystem) []*core.MetricValue {
	var total = int64(esxi.Summary.Hardware.CpuMhz) * int64(esxi.Summary.Hardware.NumCpuCores)
	perf := []*core.MetricValue{}
	if total > 0 {
		usePercentCPU := core.GaugeValue(esxi.Summary.Config.Name, "cpu.util", float64(esxi.Summary.QuickStats.OverallCpuUsage)/float64(total)*100, nil)
		perf = append(perf, usePercentCPU)
		freePercentCPU := core.GaugeValue(esxi.Summary.Config.Name, "cpu.idle", (float64(total)-float64(esxi.Summary.QuickStats.OverallCpuUsage))/float64(total)*100, nil)
		perf = append(perf, freePercentCPU)
	}

	return perf
}

//EsxiMem mem metrics
func EsxiMem(esxi mo.HostSystem) []*core.MetricValue {
	var total = esxi.Summary.Hardware.MemorySize
	var free = int64(esxi.Summary.Hardware.MemorySize) - (int64(esxi.Summary.QuickStats.OverallMemoryUsage) * 1024 * 1024)

	totalMem := core.GaugeValue(esxi.Summary.Config.Name, "mem.bytes.total", total, nil)
	useMem := core.GaugeValue(esxi.Summary.Config.Name, "mem.bytes.used", int64(esxi.Summary.QuickStats.OverallMemoryUsage)*1024*1024, nil)
	freeMem := core.GaugeValue(esxi.Summary.Config.Name, "mem.bytes.free", free, nil)
	perf := []*core.MetricValue{totalMem, useMem, freeMem}
	if total > 0 {
		usedMemPer := core.GaugeValue(esxi.Summary.Config.Name, "mem.bytes.used.percent", float64(total-free)/float64(total)*100, nil)
		perf = append(perf, usedMemPer)
	}

	return perf
}

//EsxiNet net metrics
func EsxiNet(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) ([]*core.MetricValue, error) {
	var netPerf []*core.MetricValue
	counterNameID, err := CounterWithID(ctx, c)
	if err != nil {
		return nil, err
	}
	var EsxiNetExtend = []string{"net.transmitted.average", "net.received.average"}
	extendID := CounterIDByName(counterNameID, EsxiNetExtend)
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err != nil {
			return nil, err
		}
		for _, each := range metricPerf {
			var tags = map[string]string{}
			if each.Metric == "net.received.average" {
				if each.Instance == "" {
					netPerf = append(netPerf, core.GaugeValue(esxi.Summary.Config.Name, "net.in.bits.total", each.Value*1024*8, tags))
				} else {
					tags["iface"] = each.Instance
					netPerf = append(netPerf, core.GaugeValue(esxi.Summary.Config.Name, "net.in.bits", each.Value*1024*8, tags))
				}
			}
			if each.Metric == "net.transmitted.average" {
				if each.Instance == "" {
					netPerf = append(netPerf, core.GaugeValue(esxi.Summary.Config.Name, "net.out.bits.total", each.Value*1024*8, tags))
				} else {
					tags["iface"] = each.Instance
					netPerf = append(netPerf, core.GaugeValue(esxi.Summary.Config.Name, "net.out.bits", each.Value*1024*8, tags))
				}
			}
		}

	}
	return netPerf, nil
}

//EsxiDisk disk metrics
func EsxiDisk(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem) ([]*core.MetricValue, error) {
	var diskPerf []*core.MetricValue
	pc := property.DefaultCollector(c.Client)
	dss := []mo.Datastore{}
	err := pc.Retrieve(ctx, esxi.Datastore, []string{"summary"}, &dss)
	if err != nil {
		return nil, err
	}
	var (
		freeAll  int64
		totalAll int64
		usedAll  int64
	)
	for _, ds := range dss {
		var tags = map[string]string{}
		tags["datastore="] = ds.Summary.Name
		var free = ds.Summary.FreeSpace
		var total = ds.Summary.Capacity
		var used = total - free
		freeAll += free
		totalAll += total
		usedAll += used

		diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "dsik.bytes.free", free, tags))
		diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.bytes.total", total, tags))
		diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.bytes.used", used, tags))
		if total > 0 {
			usedPercent := float64(used) / float64(total) * 100
			diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.bytes.used.Percent", usedPercent, tags))
		}

	}

	diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.cap.free", freeAll, nil))
	diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.cap.total", totalAll, nil))
	diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.cap.used", usedAll, nil))
	if totalAll > 0 {
		usedAllPercent := float64(usedAll) / float64(totalAll) * 100
		diskPerf = append(diskPerf, core.GaugeValue(esxi.Summary.Config.Name, "disk.cap.used.percent", usedAllPercent, nil))
	}

	return diskPerf, nil
}
