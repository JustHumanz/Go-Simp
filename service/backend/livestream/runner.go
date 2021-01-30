package main

import (
	"flag"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/backend/livestream/bilibili/live"
	"github.com/JustHumanz/Go-Simp/service/backend/livestream/bilibili/space"
	"github.com/JustHumanz/Go-Simp/service/backend/livestream/twitch"
	"github.com/JustHumanz/Go-Simp/service/backend/livestream/youtube"
	"github.com/JustHumanz/Go-Simp/service/backend/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	Youtube := flag.Bool("Youtube", false, "Enable youtube module")
	SpaceBiliBili := flag.Bool("SpaceBiliBili", false, "Enable Space.bilibili module")
	LiveBiliBili := flag.Bool("LiveBiliBili", false, "Enable Live.bilibili module")
	Twitch := flag.Bool("Twitch", false, "Enable Twitch module")

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

	if *Twitch {
		twitch.Start(Bot, c)
		database.ModuleInfo("Twitch")
	}

	runfunc.Run(Bot)
}
