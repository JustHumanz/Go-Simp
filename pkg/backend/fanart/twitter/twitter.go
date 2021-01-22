package twitter

import (
	"sync"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot *discordgo.Session
)

//Start start twitter module
func Start(BotInit *discordgo.Session, cronInit *cron.Cron) {
	Bot = BotInit
	cronInit.AddFunc(config.TwitterFanart, CheckNew)
	log.Info("Enable twitter fanart module")
}

//CheckNew Check new fanart
func CheckNew() {
	Scraper := twitterscraper.New()
	Scraper.SetSearchMode(twitterscraper.SearchLatest)
	err := Scraper.SetProxy(config.BotConf.MultiTOR)
	if err != nil {
		log.Error(err)
	}
	wg := new(sync.WaitGroup)
	for _, GroupData := range engine.GroupData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			Fanarts, err := CreatePayload(database.GetMembers(Group.ID), Group, Scraper, config.BotConf.LimitConf.TwitterFanart)
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
