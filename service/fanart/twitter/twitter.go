package twitter

import (
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	configfile  config.ConfigFile
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	b.AddFunc(config.TwitterFanart, CheckNew)
	log.Info("Enable twitter fanart module")
}

//CheckNew Check new fanart
func CheckNew() {
	Scraper := twitterscraper.New()
	Scraper.SetSearchMode(twitterscraper.SearchLatest)
	err := Scraper.SetProxy(config.GoSimpConf.MultiTOR)
	if err != nil {
		log.Error(err)
	}
	wg := new(sync.WaitGroup)
	for _, GroupData := range VtubersData.VtuberData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			Fanarts, err := CreatePayload(Group, Scraper, config.GoSimpConf.LimitConf.TwitterFanart)
			if err != nil {
				log.WithFields(log.Fields{
					"Group": Group.GroupName,
				}).Error(err)
			} else {
				SendFanart(Fanarts, Group)
			}
		}(GroupData, wg)
	}
	wg.Wait()
}
