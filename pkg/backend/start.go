package main

import (
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/JustHumanz/Go-simp/pkg/backend/cronjob"
	"github.com/JustHumanz/Go-simp/pkg/backend/runner"
	log "github.com/sirupsen/logrus"
)

const configfile = "../../config.toml"

func main() {
	runner.StartInit(configfile)
	cronjob.InitCron()
	log.Info("Backend ready.......")

	chain := make(chan os.Signal, 0)
	signal.Notify(chain, os.Interrupt)
	go func() {
		for sig := range chain {
			log.Warn("captured ", sig, ", stopping profiler and exiting..")
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()

	<-make(chan struct{})
	return
}
