package main

import (
	"flag"
	"tinderMatchingSystem/config"
	app "tinderMatchingSystem/internal/application"

	"os"
	"os/signal"
	"syscall"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf-folder", "../../", "config path, eg: -conf-folder ./")

}

func handleSignals(server *app.Application) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	server.Logger.Infof("signal %s received", <-sigs)
	server.Shutdown()
}

// @title         Certificate Service
// @version       0.0.1
// @description   Swagger API.
// @host          localhost:9030
// @contact.name  Skyler

func main() {

	flag.Parse()
	config.LoadConf([]string{flagconf}, config.GetConfig())

	server := app.Default()
	go handleSignals(server)
	server.Run()
}
