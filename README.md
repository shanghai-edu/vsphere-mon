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
    - 192.168.0.100 # 修改成实际的 n9e 服务器地址

probe:
  http: 127.0.0.1:2059
```
#### probe.yml
```yml
logger:
  dir: logs/
  level: INFO
  keepHours: 24

probe:
  # 如果需要区分来自不同区域的探针，可以通过在配置 region 来插入 tag
  #region: default
  timeout: 5 # 探测的超时时间，单位是秒
  limit: 10 # 并发限制
  interval: 30 # 请求的间隔
  headers: # 插入到 http 请求中的 headers，可以多条
    user-agent: Mozilla/5.0 (Linux; Android 6.0.1; Moto G (4)) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36 Edg/87.0.664.66

ping:
  107: # n9e 节点上的 nid 号
    - 114.114.114.114 # 要探测的 ip 地址列表
    - 114.114.115.115

url:
  107: # n9e 节点上的 nid 号
    - https://www.baidu.com # 要探测的 url 地址列表
    - https://www.sjtu.edu.cn/
    - https://bbs.ngacn.cc
    - https://www.163.com
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