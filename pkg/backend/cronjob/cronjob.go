package cronjob

import (
	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/bilibili"
	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/twitter"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/live"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/space"
	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	"github.com/JustHumanz/Go-simp/pkg/backend/subscriber"
	"gopkg.in/robfig/cron.v2"
)

//InitCron Start cron job
func InitCron() {

	c := cron.New()
	c.Start()

	c.AddFunc("@every 0h1m30s", twitter.CheckNew)
	c.AddFunc("@every 0h7m0s", bilibili.CheckNew)

	c.AddFunc("@every 0h5m0s", youtube.CheckSchedule)
	c.AddFunc("@every 6h0m0s", youtube.CheckPrivate)
	c.AddFunc("@every 0h4m0s", live.CheckSchedule)
	c.AddFunc("@every 0h6m0s", space.CheckVideo)

	c.AddFunc("@every 1h0m0s", subscriber.CheckYtSubsCount)
	c.AddFunc("@every 0h15m0s", subscriber.CheckTwFollowCount)
	c.AddFunc("@every 0h30m0s", subscriber.CheckBiliFollowCount)
}
