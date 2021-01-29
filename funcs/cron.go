package funcs

import (
	"context"
	"net/url"
	"time"

	"github.com/toolkits/pkg/logger"

	"github.com/shanghai-edu/vsphere-mon/config"
	"github.com/shanghai-edu/vsphere-mon/funcs/core"
	"github.com/shanghai-edu/vsphere-mon/funcs/models"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/soap"
)

func Collect() {
	sec := config.Get().Interval
	for _, vcConfig := range config.Get().Vspheres {
		go collect(vcConfig, config.Get().Performance, sec)
	}
}

func collect(vcConfig config.VsphereSection, perfConfig config.PerfSection, sec int64) {
	t := time.NewTicker(time.Second * time.Duration(sec))
	defer t.Stop()
	for {
		collectOnce(vcConfig, perfConfig, sec)
		<-t.C
	}
}

func collectOnce(vcConfig config.VsphereSection, perfConfig config.PerfSection, sec int64) {
	stime := time.Now().Unix()
	ts := time.Now().Unix()

	u, err := soap.ParseURL(vcConfig.Addr)
	if err != nil {
		logger.Errorf("parse vcenter failed, %v", err)
		collectVsphere(sec, ts, vcConfig.Nid, models.VcenterAlive(false, vcConfig.Addr))
		return
	}
	u.User = url.UserPassword(vcConfig.User, vcConfig.Pwd)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		logger.Errorf("create vmomi client failed, %v", err)
		collectVsphere(sec, ts, vcConfig.Nid, models.VcenterAlive(false, vcConfig.Addr))
		return
	}
	defer c.Logout(ctx)
	collectVsphere(sec, ts, vcConfig.Nid, models.VcenterAlive(true, vcConfig.Addr))

	//collect esxi
	esxiList, err := models.EsxiList(ctx, c)
	if err != nil {
		logger.Errorf("get esxi list failed, %v", err)
		return
	}
	for _, esxi := range esxiList {
		if err := report(esxi); err != nil {
			logger.Errorf("report esxi failed, %v", err)
			return
		}
		collectEsxi(sec, ts, models.EsxiPower(esxi))
		collectEsxi(sec, ts, models.EsxiStatus(esxi))
		collectEsxi(sec, ts, models.EsxiUptime(esxi))
		collectEsxi(sec, ts, models.EsxiCPU(esxi))
		collectEsxi(sec, ts, models.EsxiMem(esxi))
		net, err := models.EsxiNet(ctx, c, esxi)
		if err != nil {
			logger.Warningf("get esxi net failed: %v", err)
		} else {
			collectEsxi(sec, ts, net)
		}
		disk, err := models.EsxiDisk(ctx, c, esxi)
		if err != nil {
			logger.Warningf("get esxi disk failed: %v", err)
		} else {
			collectEsxi(sec, ts, disk)
		}
		//collect esxi extend performance
		if vcConfig.EsxiPerf {
			perf, err := models.EsxiExtend(ctx, c, esxi, perfConfig.Esxi)
			if err != nil {
				logger.Warningf("get esxi perfomance failed: %v", err)
			} else {
				collectEsxi(sec, ts, perf)
			}
		}
	}
	//collect vm
	if vcConfig.VM {
		vms, err := models.VsphereVirtualMachines(ctx, c)
		if err != nil {
			logger.Errorf("get vm list failed, %v", err)
			return
		}
		vmList := models.GenVirtualMachinesByName(vms, vcConfig.VmList)
		perfVms := models.GenVirtualMachinesByName(vms, vcConfig.VmPerfList)
		for _, vm := range vmList {
			collectVsphere(sec, ts, vcConfig.Nid, models.VmPower(vm))
			collectVsphere(sec, ts, vcConfig.Nid, models.VmStatus(vm))
			collectVsphere(sec, ts, vcConfig.Nid, models.VmUptime(vm))
			collectVsphere(sec, ts, vcConfig.Nid, models.VmCPU(vm))
			collectVsphere(sec, ts, vcConfig.Nid, models.VmMem(vm))
		}

		//collect vm extend perfomance
		if vcConfig.VmPerf {
			for _, vm := range perfVms {
				perf, err := models.VmExtend(ctx, c, vm, perfConfig.VM)
				if err != nil {
					logger.Warningf("get vm perfomance failed: %v", err)
				} else {
					collectVsphere(sec, ts, vcConfig.Nid, perf)
				}
			}
		}
	}
	etime := time.Now().Unix()
	logger.Infof("vc %s have been collected, time:%d s", vcConfig.Addr, etime-stime)
}

func collectEsxi(sec, ts int64, items []*core.MetricValue) {
	if items == nil || len(items) == 0 {
		return
	}
	metricValues := []*core.MetricValue{}
	for _, item := range items {
		item.Step = sec
		item.Timestamp = ts
		metricValues = append(metricValues, item)
	}
	core.Push(metricValues)
}

func collectVsphere(sec, ts int64, nid string, items []*core.MetricValue) {
	if items == nil || len(items) == 0 {
		return
	}
	metricValues := []*core.MetricValue{}
	for _, item := range items {
		item.Step = sec
		item.Timestamp = ts
		item.Nid = nid
		metricValues = append(metricValues, item)
	}
	core.Push(metricValues)
}
