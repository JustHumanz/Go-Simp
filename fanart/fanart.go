package fanart

import (
	bilibili "github.com/JustHumanz/Go-simp/fanart/bilibili"
	twitter "github.com/JustHumanz/Go-simp/fanart/twitter"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func Start(c *cron.Cron) {
	c.AddFunc("@every 0h1m0s", twitter.CheckNew)
	c.AddFunc("@every 0h2m0s", bilibili.CheckNew)
	log.Info("Fanart module ready")

}
