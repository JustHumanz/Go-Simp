package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"./config"
	"./engine"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	engine.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			err = engine.Stop()
			if err != nil {
				log.Print(err)
			}
			log.Printf("captured %v, stopping profiler and exiting..", sig)
			pprof.StopCPUProfile()
			os.Exit(1)
		}
	}()

	<-make(chan struct{})
	return
}
