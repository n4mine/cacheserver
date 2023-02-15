package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/n4mine/cacheserver/cache"
	"github.com/n4mine/cacheserver/config"
	"github.com/n4mine/cacheserver/models"
	"github.com/n4mine/cacheserver/rpc"
	"github.com/n4mine/cacheserver/web"
)

func main() {
	fmt.Printf("build from git commit %s, golang version %s\n", models.GitVer, runtime.Version())

	configFile := flag.String("c", "etc/config.cfg", "config file path")
	version := flag.Bool("v", false, "show version")
	flag.Parse()

	if *version {
		fmt.Println("version:", models.GitVer)
		fmt.Println("build time:", models.BuildTime)
		os.Exit(0)
	}

	config := config.LoadConfig(*configFile)

	// logger
	initLog()

	// init
	cache.InitCaches()

	// business logic
	go web.Start(config.Web)
	go rpc.Start(config.Rpc)
	go cache.GC(config.GC)

	fmt.Println("start ok")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM)
	s := <-sc
	fmt.Println("catch signal: ", s)
	fmt.Println("stop ok")
}

func initLog() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
