package stream

import (
	engine "github.com/JustHumanz/Go-simp/engine"
	bilibililive "github.com/JustHumanz/Go-simp/livestream/bilibili/live"
	spacebilibili "github.com/JustHumanz/Go-simp/livestream/bilibili/space"
	youtube "github.com/JustHumanz/Go-simp/livestream/youtube"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func Start(c *cron.Cron) {
	Bot := engine.BotSession
	youtube.Start(Bot)
	bilibililive.Start(Bot)
	spacebilibili.Start()
	c.AddFunc("@every 0h5m0s", youtube.CheckSchedule)
	c.AddFunc("@every 4h0m0s", youtube.CheckPrivate)
	c.AddFunc("@every 0h5m0s", bilibililive.CheckSchedule)
	c.AddFunc("@every 0h5m0s", spacebilibili.CheckVideo)
	log.Info("Stream module ready")
}
