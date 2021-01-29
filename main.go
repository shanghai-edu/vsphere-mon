package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shanghai-edu/vsphere-mon/config"
	"github.com/shanghai-edu/vsphere-mon/funcs"
	"github.com/shanghai-edu/vsphere-mon/funcs/core"
	"github.com/shanghai-edu/vsphere-mon/http"

	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/logger"
	"github.com/toolkits/pkg/runner"
)

var (
	vers *bool
	help *bool
	conf *string
)

func init() {
	vers = flag.Bool("v", false, "display the version.")
	help = flag.Bool("h", false, "print this help.")
	conf = flag.String("f", "", "specify configuration file.")
	flag.Parse()

	if *vers {
		fmt.Println("version:", config.Version)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	aconf()
	pconf()
	start()

	cfg := config.Get()
	config.InitLog(cfg.Logger)

	core.InitRpcClients()
	funcs.Collect()

	http.Start()
	ending()
}

// auto detect configuration file
func aconf() {
	if *conf != "" && file.IsExist(*conf) {
		return
	}

	*conf = "etc/vsphere.local.yml"
	if file.IsExist(*conf) {
		return
	}

	*conf = "etc/vsphere.yml"
	if file.IsExist(*conf) {
		return
	}

	fmt.Println("no configuration file for vsphere-mon")
	os.Exit(1)
}

// parse configuration file
func pconf() {
	if err := config.Parse(*conf); err != nil {
		fmt.Println("cannot parse configuration file:", err)
		os.Exit(1)
	}
}

func ending() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-c:
		fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())
	}

	logger.Close()
	http.Shutdown()
	fmt.Println("vsphere-mon stopped successfully")
}

func start() {
	runner.Init()
	fmt.Println("vsphere-mon start, use configuration file:", *conf)
	fmt.Println("runner.cwd:", runner.Cwd)
}
