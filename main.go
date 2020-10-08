package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"
	fanart "github.com/JustHumanz/Go-simp/fanart"
	stream "github.com/JustHumanz/Go-simp/livestream"
	subscriber "github.com/JustHumanz/Go-simp/subscriber"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	db, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Bot, _ := discordgo.New("Bot " + config.Token)

	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	u, err := Bot.User("@me")
	c := cron.New()
	c.Start()

	database.Start(db)
	engine.Start(Bot, u.ID)
	fanart.Start(c)
	stream.Start(c)
	subscriber.Start(c)

	log.Info("Bombing bays ready.......")

	chain := make(chan os.Signal, 0)
	signal.Notify(chain, os.Interrupt)
	go func() {
		for sig := range chain {
			log.Warn("Stopping Cron")
			c.Stop()
			log.Warn("Shutdown DB")
			db.Close()
			log.Warn("captured ", sig, ", stopping profiler and exiting..")
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()

	<-make(chan struct{})
	return
}
