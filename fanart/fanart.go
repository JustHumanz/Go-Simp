package fanart

import (
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	bilibili "github.com/JustHumanz/Go-simp/fanart/bilibili"
	twitter "github.com/JustHumanz/Go-simp/fanart/twitter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func Start(c *cron.Cron) {
	bot := engine.BotSession
	db := database.DB
	twitter.Twitterd(bot, db)
	bilibili.BiliBili(bot)

	c.AddFunc("@every 0h1m0s", twitter.CheckNew)
	c.AddFunc("@every 0h2m0s", bilibili.CheckNew)
	log.Info("Fanart module ready")

}
