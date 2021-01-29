package config

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

type ConfYaml struct {
	Logger      LoggerSection    `yaml:"logger"`
	Interval    int64            `yaml:"interval"`
	Report      ReportSection    `yaml:"report"`
	Vspheres    []VsphereSection `yaml:"vspheres"`
	Performance PerfSection      `yaml:"performance"`
}

type ReportSection struct {
	Token   string `yaml:"token"`
	Cate    string `yaml:"cate"`
	UniqKey string `yaml:"uniqkey"`
	Tenant  string `yaml:"tenant"`
}

type VsphereSection struct {
	Addr       string   `yaml:"addr"`
	User       string   `yaml:"user"`
	Pwd        string   `yaml:"pwd"`
	EsxiPerf   bool     `yaml:"esxiperf"`
	VM         bool     `yaml:"vm"`
	Nid        string   `yaml:"nid"`
	VmList     []string `yaml:"vmlist"`
	VmPerf     bool     `yaml:"vmperf"`
	VmPerfList []string `yaml:"vmperflist"`
}

type PerfSection struct {
	Esxi []string `yaml:"esxi"`
	VM   []string `yaml:"vm"`
}

var (
	Config   *ConfYaml
	lock     = new(sync.RWMutex)
	Endpoint string
	Cwd      string
)

// Get configuration file
func Get() *ConfYaml {
	lock.RLock()
	defer lock.RUnlock()
	return Config
}

func Parse(conf string) error {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	lock.Lock()
	defer lock.Unlock()

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		return fmt.Errorf("cannot read yml[%s]: %v", conf, err)
	}

	err = viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("Unmarshal %v", err)
	}

	return nil
}
