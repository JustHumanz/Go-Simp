package cronjob

import (
	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/bilibili"
	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/twitter"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/live"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/space"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	"github.com/JustHumanz/Go-simp/pkg/backend/runner"
	"github.com/JustHumanz/Go-simp/pkg/backend/subscriber"
	"github.com/JustHumanz/Go-simp/tools/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

//InitCron Start cron job
func InitCron() {
	c := cron.New()
	c.Start()

	Flags := runner.Flags
	if *Flags["twitterfanart"].(*bool) {
		c.AddFunc("@every 0h3m0s", twitter.CheckNew)
		log.Info("Add twitter fanart to cronjob")
	}

	if *Flags["bilibilifanart"].(*bool) {
		c.AddFunc("@every 0h7m0s", bilibili.CheckNew)
		log.Info("Add bilibili fanart to cronjob")
	}

	if *Flags["youtube"].(*bool) {
		c.AddFunc("@every 0h5m0s", youtube.CheckSchedule)
		c.AddFunc("@every 2h0m0s", youtube.CheckPrivate)
		log.Info("Add youtube to cronjob")
	}

	if *Flags["bilibili"].(*bool) {
		c.AddFunc("@every 0h4m0s", live.CheckSchedule)
		log.Info("Add bilibili to cronjob")
	}

	if *Flags["space.bilibili"].(*bool) {
		c.AddFunc("@every 0h6m0s", space.CheckVideo)
		log.Info("Add space.bilibili to cronjob")
	}

	if *Flags["subscriber"].(*bool) {
		c.AddFunc("@every 1h0m0s", subscriber.CheckYtSubsCount)
		c.AddFunc("@every 0h15m0s", subscriber.CheckTwFollowCount)
		c.AddFunc("@every 0h30m0s", subscriber.CheckBiliFollowCount)
		log.Info("Add subscriber to cronjob")
	}
	if config.BotConf.DonationLink != "" {
		c.AddFunc("@every 0h30m0s", runner.Donate)
	}
}
