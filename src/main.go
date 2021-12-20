package main

import (
	"flag"

	"autoSignIn/src/config"
	"autoSignIn/src/crawler"
)

var configPath = flag.String("config", "", "Path to the configuration file.")

func main() {
	flag.Parse()
	config.InitConfig(*configPath)
	crawler.NotifyChanify(crawler.RunTask())
}
