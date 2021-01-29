package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"

	"github.com/vmware/govmomi"
)

const (
	USER   = "administrator@vsphere.local"
	PASS   = "password"
	VcAddr = "https://1.1.1.1/sdk"
	VcIP   = "1.1.1.1"
)

var ctx = context.Background()
var c *govmomi.Client
var esxiList []mo.HostSystem
var vmList []mo.VirtualMachine
var vmnames = []string{"n9e-v3-centos7-92.18", "VSAN-VC", "splunk-demo1-10.26"}

func init() {
	u, err := soap.ParseURL(VcAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	u.User = url.UserPassword(USER, PASS)
	c, err = govmomi.NewClient(ctx, u, true)

	if err != nil {
		fmt.Println(err)
		return
	}
}
func Test_Esxi(t *testing.T) {
	if c.IsVC() {
		t.Log("connection successful!")
	}
	var err error
	esxiList, err = EsxiList(ctx, c)

	if err != nil {
		t.Error(err)
		return
	}
	bs, _ := json.Marshal(esxiList[0])
	t.Log(string(bs))

}

func Test_EsxiPower(t *testing.T) {
	for _, esxi := range esxiList {
		res := EsxiPower(esxi)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}

	}
}

func Test_EsxiStatus(t *testing.T) {
	for _, esxi := range esxiList {
		res := EsxiStatus(esxi)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_EsxiUptime(t *testing.T) {
	for _, esxi := range esxiList {
		res := EsxiUptime(esxi)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_EsxiMem(t *testing.T) {
	for _, esxi := range esxiList {
		res := EsxiMem(esxi)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_EsxiNet(t *testing.T) {
	for _, esxi := range esxiList {
		res, err := EsxiNet(ctx, c, esxi)
		if err != nil {
			t.Error(err)
			return
		}
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_EsxiDisk(t *testing.T) {
	for _, esxi := range esxiList {
		res, err := EsxiDisk(ctx, c, esxi)
		if err != nil {
			t.Error(err)
			return
		}
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_EsxiExtend(t *testing.T) {
	extend := []string{
		"cpu.coreUtilization.average",
		"cpu.costop.summation",
		"cpu.demand.average",
		"cpu.idle.summation",
		"cpu.latency.average",
		"cpu.readiness.average",
		"cpu.ready.summation",
		"cpu.swapwait.summation",
		"cpu.usage.average",
		"cpu.usagemhz.average",
		"cpu.used.summation",
		"cpu.utilization.average",
		"cpu.wait.summation",
		"disk.deviceReadLatency.average",
		"disk.deviceWriteLatency.average",
		"disk.kernelReadLatency.average",
		"disk.kernelWriteLatency.average",
		"disk.numberReadAveraged.average",
		"disk.numberWriteAveraged.average",
		"disk.read.average",
		"disk.totalReadLatency.average",
		"disk.totalWriteLatency.average",
		"disk.write.average",
		"mem.active.average",
		"mem.latency.average",
		"mem.state.latest",
		"mem.swapin.average",
		"mem.swapinRate.average",
		"mem.swapout.average",
		"mem.swapoutRate.average",
		"mem.totalCapacity.average",
		"mem.usage.average",
		"mem.vmmemctl.average",
		"net.bytesRx.average",
		"net.bytesTx.average",
		"net.droppedRx.summation",
		"net.droppedTx.summation",
		"net.errorsRx.summation",
		"net.errorsTx.summation",
		"net.usage.average",
		"power.power.average",
		"storageAdapter.numberReadAveraged.average",
		"storageAdapter.numberWriteAveraged.average",
		"storageAdapter.read.average",
		"storageAdapter.write.average",
		"sys.uptime.latest",
	}
	for _, esxi := range esxiList {
		res, err := EsxiExtend(ctx, c, esxi, extend)
		if err != nil {
			t.Error(err)
			return
		}
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_VsphereDatastores(t *testing.T) {
	res, err := VsphereDatastores(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range res {
		bs, _ := json.Marshal(r)
		t.Log(string(bs))
	}
}

func Test_DatastoreMetrics(t *testing.T) {
	res, err := DatastoreMetrics(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}
	for _, r := range res {
		bs, _ := json.Marshal(r)
		t.Log(string(bs))
	}
}

func Test_VsphereVirtualMachines(t *testing.T) {
	var err error
	vmList, err = VsphereVirtualMachines(ctx, c)
	if err != nil {
		t.Error(err)
		return
	}
}

func Test_VmPower(t *testing.T) {
	for _, vm := range vmList {
		res := VmPower(vm)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_VmStatus(t *testing.T) {
	for _, vm := range vmList {
		res := VmStatus(vm)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}
func Test_VmUptime(t *testing.T) {
	for _, vm := range vmList {
		res := VmUptime(vm)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}
func Test_VmCpu(t *testing.T) {
	for _, vm := range vmList {
		res := VmCPU(vm)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}
func Test_VmMem(t *testing.T) {
	for _, vm := range vmList {
		res := VmMem(vm)
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_VmExtend(t *testing.T) {

	extend := []string{
		"cpu.demand.average",
		"cpu.idle.summation",
		"cpu.latency.average",
		"cpu.readiness.average",
		"cpu.ready.summation",
		"cpu.run.summation",
		"cpu.usagemhz.average",
		"cpu.used.summation",
		"cpu.wait.summation",
		"mem.active.average",
		"mem.granted.average",
		"mem.latency.average",
		"mem.swapin.average",
		"mem.swapinRate.average",
		"mem.swapout.average",
		"mem.swapoutRate.average",
		"mem.usage.average",
		"mem.vmmemctl.average",
		"net.bytesRx.average",
		"net.bytesTx.average",
		"net.droppedRx.summation",
		"net.droppedTx.summation",
		"net.usage.average",
		"power.power.average",
		"virtualDisk.numberReadAveraged.average",
		"virtualDisk.numberWriteAveraged.average",
		"virtualDisk.read.average",
		"virtualDisk.readOIO.latest",
		"virtualDisk.throughput.usage.average",
		"virtualDisk.totalReadLatency.average",
		"virtualDisk.totalWriteLatency.average",
		"virtualDisk.write.average",
		"virtualDisk.writeOIO.latest",
		"sys.uptime.latest",
	}
	vms := GenVirtualMachinesByName(vmList, vmnames)
	for _, vm := range vms {
		bs, _ := json.Marshal(vm)
		t.Log(string(bs))
		res, err := VmExtend(ctx, c, vm, extend)
		if err != nil {
			t.Error(err)
			return
		}
		for _, r := range res {
			bs, _ := json.Marshal(r)
			t.Log(string(bs))
		}
	}
}

func Test_Logout(t *testing.T) {
	err := c.Logout(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("logout success!")
}
