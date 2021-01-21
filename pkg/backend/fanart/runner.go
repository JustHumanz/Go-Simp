package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/bilibili"
	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/twitter"
	"github.com/JustHumanz/Go-simp/pkg/backend/utility/runfunc"
	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	Twitter := flag.Bool("TwitterFanart", false, "Enable twitter fanart module")
	BiliBili := flag.Bool("BiliBiliFanart", false, "Enable bilibili fanart module")
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

	if *Twitter {
		twitter.Start(Bot, c)
		database.ModuleInfo("TwitterFanart")
	}

	if *BiliBili {
		bilibili.Start(Bot, c)
		database.ModuleInfo("BiliBiliFanart")
	}
	runfunc.Run()
}
