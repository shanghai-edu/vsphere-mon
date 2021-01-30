# vsphere-mon

## 功能
适配 [nightingale](https://github.com/didi/nightingale)，采集 vsphere 相关指标
支持 ESXi 和 VM 相关指标监控
支持 ESXi 作为主机资产自动注册的 nightingale 的 ams 中


### 指标
#### ESXi
ESXi 以设备相关的方式上报数据，设备资产数据会自动注册到夜莺中
##### report 信息
|字段|说明|
|--|--|
|sn|硬件序列号，如果取不到会使用 uuid 替代|
|endpoint|esxi 的 name，通常是 ip 地址|
|ip|esxi 的 name，通常是 ip 地址，如果检查不是 ip 则留空|
|name|esxi 的 name，通常是 ip 地址|
|cate|分类，根据配置决定，默认是 physical|
|tenant|租户，根据配置决定，默认是空|
|cpu|物理核心数，不考虑超线程|
|mem|内存，单位是G|
|model|硬件型号，在 ams 中创建扩展字段 model 后可见|
|version|esxi 的 fullname，即类似 VMware ESXi 6.7.0 build-13473784|

##### 基础指标
|metric|说明|
|--|--|
|esxi.power|1:poweredOff，2:poweredOn，3:standBy，4:unknown,可能断开连接或者无响应|
|exsi.status|1:gray,未知状态;2:green，正常;3:red，大毛病;4:yellow，小毛病|
|esxi.uptime|uptime|
|cpu.idle|cpu 空闲率|
|cpu.util|cpu 使用率|
|mem.bytes.total|总内存|
|mem.bytes.used|使用内存|
|mem.bytes.free|空闲内存|
|mem.bytes.used.percent|内存使用率|
|net.in.bits.total|总入流量|
|net.in.bits|单块网卡的入流量，iface=xxx|
|net.out.bits.total|总出流量|
|net.out.bits|单块网卡的出流量，iface=xxx|
|dsik.bytes.free|单块盘（存储）空闲容量，datastore=xxx|
|disk.bytes.total|单块盘（存储）总容量，datastore=xxx|
|disk.bytes.used|单块盘（存储）使用容量，datastore=xxx|
|disk.bytes.used.Percent|单块盘（存储）使用率，datastore=xxx|
|disk.cap.free|存储总空闲量|
|disk.cap.total|存储总量|
|disk.cap.used|存储总使用量|
|disk.cap.used.percent|存储总使用率|

##### 扩展指标
根据 performance 中的配置决定，有啥采啥

#### VM
VM 以设备无关的方式上报数据，虚机的名字以 name=xxx 的方式作为 tag 体现
##### 基础指标
|metric|说明|
|--|--|
|vm.power|1:poweredOff，2:poweredOn，3:standBy，4:unknown,可能断开连接或者无响应|
|vm.status|1:gray,未知状态;2:green，正常;3:red，大毛病;4:yellow，小毛病|
|vm.uptime|uptime|
|cpu.idle|cpu 空闲率|
|cpu.util|cpu 使用率|
|mem.bytes.total|总内存|
|mem.bytes.guest.used|虚机实际使用内存|
|mem.bytes.host.used|分配给虚拟机的内存|
|mem.bytes.guest.used.percent|虚机实际内存使用率|

##### 扩展指标
根据 performance 中的配置决定，有啥采啥

#### 其他
以下指标也已设备无关方式上报

|metric|说明|
|--|--|
|vcetner.alive|vcenter 连接状态，1通0不通|
|datastore.bytes.total|存储容量，ds=xxx,fstype=xxx|
|datastore.bytes.free|存储空闲容量，ds=xxx,fstype=xxx|
|datastore.bytes.used|存储使用容量，ds=xxx,fstype=xxx|
|datastore.used.percent|存储使用率，ds=xxx,fstype=xxx|
### 配置
#### address.yml
```yml
---
transfer:
  http: 0.0.0.0:8008
  rpc: 0.0.0.0:8009
  addresses:
    - 192.168.100.1 # 修改成实际的 n9e 地址

ams:
  http: 0.0.0.0:8002
  addresses:
    - 192.168.100.1
    
vsphere-mon:
  http: 127.0.0.1:2060
```
#### vsphere.yml
```yml
logger:
  dir: logs/
  level: WARNING
  keepHours: 24

# 上报的间隔，注意关注下 info.log 的日志，确保能够在一个周期内完成采集
interval: 300

report:
  # 调用ams的接口上报数据，需要ams的token
  token: ams-builtin-token
  # physical：物理机，virtual：虚拟机，container：容器，switch：交换机
  cate: physical
  # 使用哪个字段作为唯一KEY，即作为where条件更新对应记录，一般使用sn或ip
  uniqkey: ip  
  # 租户，如果配置则直接注册到该租户下
  tenant: 
# 要监控的 vsphere 的配置信息

vspheres:
    # vcenter 的地址
  - addr: https://1.1.1.1/sdk
    # vcenter 的用户名
    user: administrator@vsphere.local
    # vcetner 的密码
    pwd: password
    # 是否开启 esxi 的扩展指标监控，注意这会增加 vcenter 的负担
    esxiperf: true
    # 是否开启虚拟机的监控，注意这会增加 vcenter 的负担
    vm: true
    # 虚拟机监控所在的节点 ID（设备无关）
    nid: 137
    # 采集的虚拟机列表，如果是空数组则采集所有的虚拟机信息。
    # 虚拟机数量的增加不会增加 vcenter 的负担，这里允许控制虚拟机采集数量的目的是可以减少 n9e 的负担，削减指标数量。
    vmlist: ["VC"]
    # 是否开启虚拟机的扩展指标监控，注意这会增加 vcenter 的负担
    vmperf: true
    # 采集虚拟机扩展指标监控的虚机列表，如果是空数组则采集所有虚机的扩展指标
    # 注意这里的虚机数量越多，对 vcenter 的负担越大，建议只对重点关注的虚机开启
    vmperflist: ["VC"] 
  - addr: https://2.2.2.2/sdk
    user: administrator@vsphere.local
    pwd: password
    esxiperf: false
    vm: false
    nid: 138
    vmlist: []
    vmperf: false
    vmperflist: [] 

# 扩展的性能指标，注意采集越多对 vc 的负担就越大    
# 建议根据实际需求配置
# 更多指标和相关含义见 vmware 官网 
# https://vdc-repo.vmware.com/vmwb-repository/dcr-public/790263bc-bd30-48f1-af12-ed36055d718b/e5f17bfc-ecba-40bf-a04f-376bbb11e811/vim.PerformanceManager.html#counterTables
performance:    
  # esxi 宿主机的额外扩展指标
  esxi:
    - cpu.coreUtilization.average
    - cpu.costop.summation
    - cpu.demand.average
    - cpu.idle.summation
    - cpu.latency.average
    - cpu.readiness.average
    - cpu.ready.summation
    - cpu.swapwait.summation
    - cpu.usage.average
    - cpu.usagemhz.average
    - cpu.used.summation
    - cpu.utilization.average
    - cpu.wait.summation
    - disk.deviceReadLatency.average
    - disk.deviceWriteLatency.average
    - disk.kernelReadLatency.average
    - disk.kernelWriteLatency.average
    - disk.numberReadAveraged.average
    - disk.numberWriteAveraged.average
    - disk.read.average
    - disk.totalReadLatency.average
    - disk.totalWriteLatency.average
    - disk.write.average
    - mem.active.average
    - mem.latency.average
    - mem.state.latest
    - mem.swapin.average
    - mem.swapinRate.average
    - mem.swapout.average
    - mem.swapoutRate.average
    - mem.totalCapacity.average
    - mem.usage.average
    - mem.vmmemctl.average
    - net.bytesRx.average
    - net.bytesTx.average
    - net.droppedRx.summation
    - net.droppedTx.summation
    - net.errorsRx.summation
    - net.errorsTx.summation
    - net.usage.average
    - power.power.average
    - storageAdapter.numberReadAveraged.average
    - storageAdapter.numberWriteAveraged.average
    - storageAdapter.read.average
    - storageAdapter.write.average
    - sys.uptime.latest
  # vm 虚拟机的额外扩展指标
  vm:
    - cpu.demand.average
    - cpu.idle.summation
    - cpu.latency.average
    - cpu.readiness.average
    - cpu.ready.summation
    - cpu.run.summation
    - cpu.usagemhz.average
    - cpu.used.summation
    - cpu.wait.summation
    - mem.active.average
    - mem.granted.average
    - mem.latency.average
    - mem.swapin.average
    - mem.swapinRate.average
    - mem.swapout.average
    - mem.swapoutRate.average
    - mem.usage.average
    - mem.vmmemctl.average
    - net.bytesRx.average
    - net.bytesTx.average
    - net.droppedRx.summation
    - net.droppedTx.summation
    - net.usage.average
    - power.power.average
    - virtualDisk.numberReadAveraged.average
    - virtualDisk.numberWriteAveraged.average
    - virtualDisk.read.average
    - virtualDisk.readOIO.latest
    - virtualDisk.throughput.usage.average
    - virtualDisk.totalReadLatency.average
    - virtualDisk.totalWriteLatency.average
    - virtualDisk.write.average
    - virtualDisk.writeOIO.latest
    - sys.uptime.latest
```

## 编译
```
# cd /home
# git clone https://github.com/shanghai-edu/vsphere-mon
# cd vsphere-mon
# ./control build
```
也可以直接在 release 中下载打包好的二进制
## 运行
### 支持 `systemctl` 的操作系统，如 `CentOS7`
执行 `install.sh` 脚本即可，`systemctl` 将托管运行

```
# ./install.sh 
Created symlink from /etc/systemd/system/multi-user.target.wants/vsphere-mon.service to /usr/lib/systemd/system/vsphere-mon.service.
```
后续可通过 `systemctl start/stop/restart vsphere-mon` 来进行服务管理

注意如果没有安装在 `/home` 路径上，则需要修改 `service/vsphere-mon.service` 中的相关路径，否则 `systemctl` 注册时会找不到

### 不支持 systemctl 的操作系统
执行 `./control start` 启动即可
```
# ./control start
vsphere-mon started
```
后续可通过 `./control start/stop/restart` 来进行服务管理