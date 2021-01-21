package main

import (
	"flag"

	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/live"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/space"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	"github.com/JustHumanz/Go-simp/pkg/backend/utility/runfunc"
	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	Youtube := flag.Bool("Youtube", false, "Enable youtube module")
	SpaceBiliBili := flag.Bool("SpaceBiliBili", false, "Enable space.bilibili  module")
	LiveBiliBili := flag.Bool("LiveBiliBili", false, "Enable live.bilibili  module")

	flag.Parse()

	conf, err := config.ReadConfig("../../../config.toml")
	if err != nil {
		log.Panic(err)
	}

	Bot, err := discordgo.New("Bot " + config.BotConf.Discord)
	if err != nil {
		log.Error(err)
	}
	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	database.Start(conf.CheckSQL())
	engine.Start()

	c := cron.New()
	c.Start()

	if *Youtube {
		youtube.Start(Bot, c)
		database.ModuleInfo("Youtube")
	}

	if *SpaceBiliBili {
		space.Start(Bot, c)
		database.ModuleInfo("SpaceBiliBili")
	}

	if *LiveBiliBili {
		live.Start(Bot, c)
		database.ModuleInfo("LiveBiliBili")
	}
	runfunc.Run()
}
