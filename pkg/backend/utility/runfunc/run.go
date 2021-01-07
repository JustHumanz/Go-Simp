package runfunc

import (
	"os"
	"os/signal"
	"runtime/pprof"

	log "github.com/sirupsen/logrus"
)

func Run() {
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
