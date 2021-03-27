package twitter

import (
	"regexp"
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/fanart/notif"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	configfile  config.ConfigFile
	lewd        bool
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile, e bool) {
	Bot = a
	VtubersData = c
	configfile = d
	lewd = e
	b.AddFunc(config.TwitterFanart, CheckNew)
	if lewd {
		log.Info("Enable Twitter lewd fanart module")
	} else {
		log.Info("Enable Twitter fanart module")
	}
}

//CheckNew Check new fanart
func CheckNew() {
	wg := new(sync.WaitGroup)
	for _, GroupData := range VtubersData.VtuberData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			Fanarts, err := CreatePayload(Group, config.Scraper, config.GoSimpConf.LimitConf.TwitterFanart)
			if err != nil {
				log.WithFields(log.Fields{
					"Group": Group.GroupName,
				}).Error(err)
			} else {
				for _, Art := range Fanarts {
					Color, err := engine.GetColor(config.TmpDir, Art.Photos[0])
					if err != nil {
						log.Error(err)
					}
					notif.SendNude(Art, Bot, Color)
				}
			}
		}(GroupData, wg)
	}
	wg.Wait()
}

//RemoveTwitterShortLink remove twitter shotlink
func RemoveTwitterShortLink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}
