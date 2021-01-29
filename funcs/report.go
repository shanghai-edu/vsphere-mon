package funcs

import (
	"fmt"
	"math/rand"

	"sort"

	"time"

	"github.com/toolkits/pkg/net/httplib"
	"github.com/toolkits/pkg/str"

	"github.com/shanghai-edu/vsphere-mon/config"
	"github.com/shanghai-edu/vsphere-mon/config/address"

	"github.com/vmware/govmomi/vim25/mo"
)

type hostRegisterForm struct {
	SN      string            `json:"sn"`
	IP      string            `json:"ip"`
	Ident   string            `json:"ident"`
	Name    string            `json:"name"`
	Cate    string            `json:"cate"`
	UniqKey string            `json:"uniqkey"`
	Fields  map[string]string `json:"fields"`
	Digest  string            `json:"digest"`
}

type errRes struct {
	Err string `json:"err"`
}

func getEsxiSn(esxi mo.HostSystem) (sn string) {
	for _, IdentifyingInfo := range esxi.Summary.Hardware.OtherIdentifyingInfo {
		if IdentifyingInfo.IdentifierType.GetElementDescription().Label == "Serial number tag" {
			sn = IdentifyingInfo.IdentifierValue
			return
		}
	}
	if sn == "" {
		sn = esxi.Summary.Hardware.Uuid
	}
	return
}

func report(esxi mo.HostSystem) error {

	fields := map[string]string{
		"cpu":     fmt.Sprintf("%d", esxi.Summary.Hardware.NumCpuCores),
		"mem":     fmt.Sprintf("%.2fG", float64(esxi.Summary.Hardware.MemorySize)/float64(1024*1024*1024)),
		"model":   esxi.Summary.Hardware.Model,
		"version": esxi.Summary.Config.Product.FullName,
		"tenant":  config.Get().Report.Tenant,
	}

	form := hostRegisterForm{
		SN:      getEsxiSn(esxi),
		IP:      esxi.Summary.Config.Name,
		Ident:   esxi.Summary.Config.Name,
		Name:    esxi.Summary.Config.Name,
		Cate:    config.Get().Report.Cate,
		UniqKey: config.Get().Report.UniqKey,
		Fields:  fields,
	}

	content := form.SN + form.IP + form.Ident + form.Name + form.Cate + form.UniqKey
	var keys []string
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		content += fields[key]
	}

	form.Digest = str.MD5(content)

	servers := address.GetHTTPAddresses("ams")
	for _, i := range rand.Perm(len(servers)) {
		url := fmt.Sprintf("http://%s/v1/ams-ce/hosts/register", servers[i])

		var body errRes
		err := httplib.Post(url).JSONBodyQuiet(form).Header("X-Srv-Token", config.Config.Report.Token).SetTimeout(time.Second * 5).ToJSON(&body)
		if err != nil {
			return fmt.Errorf("curl %s fail: %v", url, err)
		}

		if body.Err != "" {
			return fmt.Errorf(body.Err)
		}

		return nil
	}

	return fmt.Errorf("all server instance is dead")
}
