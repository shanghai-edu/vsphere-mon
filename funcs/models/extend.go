package models

import (
	"context"
	"strings"

	"github.com/shanghai-edu/vsphere-mon/funcs/core"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
)

func EsxiExtend(ctx context.Context, c *govmomi.Client, esxi mo.HostSystem, extend []string) ([]*core.MetricValue, error) {
	counterNameID, err := CounterWithID(ctx, c)
	if err != nil {
		return nil, err
	}
	extendID := CounterIDByName(counterNameID, extend)
	perf := []*core.MetricValue{}
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, esxi.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags = map[string]string{}
				if each.Instance != "" {
					sp := strings.Split(each.Metric, ".")
					tags[sp[0]] = each.Instance
				}
				perf = append(perf, core.GaugeValue(esxi.Summary.Config.Name, each.Metric, each.Value, tags))
			}
		}
	}
	return perf, nil
}

func VmExtend(ctx context.Context, c *govmomi.Client, vm mo.VirtualMachine, extend []string) ([]*core.MetricValue, error) {
	counterNameID, err := CounterWithID(ctx, c)
	if err != nil {
		return nil, err
	}
	extendID := CounterIDByName(counterNameID, extend)
	perf := []*core.MetricValue{}
	for _, k := range extendID {
		metricPerf, err := Performance(ctx, c, vm.Self, k)
		if err == nil {
			for _, each := range metricPerf {
				var tags = map[string]string{}
				if each.Instance != "" {
					sp := strings.Split(each.Metric, ".")
					tags[sp[0]] = each.Instance
				}
				tags["name"] = vm.Summary.Config.Name
				perf = append(perf, core.GaugeValue("", each.Metric, each.Value, tags))
			}
		}
	}
	return perf, nil
}
