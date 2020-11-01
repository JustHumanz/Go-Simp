package main

import (
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/JustHumanz/Go-simp/pkg/frontend/discord_handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	discord_handler.StartInit("../../config.toml")
	log.Info("Frontend ready.......")

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
