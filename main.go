package main

import (
	"flag"
	"fmt"
	"github.com/ZeaLoVe/falcon-winAgent/cron"
	"github.com/ZeaLoVe/falcon-winAgent/funcs"
	"github.com/ZeaLoVe/falcon-winAgent/g"
	"github.com/ZeaLoVe/falcon-winAgent/http"
	"os"
)

func main() {

	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	g.ParseConfig(*cfg)

	g.InitRootDir()
	g.InitLocalIps()
	g.InitRpcClients()

	funcs.BuildMappers()

	go cron.InitDataHistory()

	cron.ReportAgentStatus()
	cron.SyncMinePlugins()
	cron.SyncBuiltinMetrics()
	cron.SyncTrustableIps()
	cron.Collect()

	go http.Start()
	select {}

}
