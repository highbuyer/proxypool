package main

import (
	"flag"
	_ "net/http/pprof"
	"os"

	"github.com/highbuyer/proxypool/api"
	"github.com/highbuyer/proxypool/config"
	"github.com/highbuyer/proxypool/internal/app"
	"github.com/highbuyer/proxypool/internal/cron"
	"github.com/highbuyer/proxypool/internal/database"
	"github.com/highbuyer/proxypool/log"
	"github.com/highbuyer/proxypool/pkg/geoIp"
)

var debugMode = false

func main() {
	var configFilePath = ""

	//go func() {
	//	http.ListenAndServe("0.0.0.0:6060", nil)
	//}()

	flag.StringVar(&configFilePath, "c", "", "path to config file: config.yaml")
	flag.BoolVar(&debugMode, "d", false, "debug output")
	flag.Parse()

	log.SetLevel(log.INFO)
	if debugMode {
		log.SetLevel(log.DEBUG)
		log.Debugln("=======Debug Mode=======")
	}
	if configFilePath == "" {
		configFilePath = os.Getenv("CONFIG_FILE")
	}
	if configFilePath == "" {
		configFilePath = "config\\config.yaml"
	}

	config.SetFilePath(configFilePath)

	err := app.InitConfigAndGetters()
	if err != nil {
		log.Errorln("Configuration init error: %s", err.Error())
		panic(err)
	}

	exe, _ := os.Executable()
	log.Infoln("Running image path: %s", exe)

	database.InitTables()
	// init GeoIp db reader and map between emoji's and countries
	// return: struct geoIp (dbreader, emojimap)
	err = geoIp.InitGeoIpDB()
	if err != nil {
		os.Exit(1)
	}
	log.Infoln("Do the first crawl...")
	go app.CrawlGo() // 抓取主程序
	go cron.Cron()   // 定时运行
	api.Run()        // Web Serve
}
