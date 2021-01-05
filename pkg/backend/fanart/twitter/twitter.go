package twitter

import (
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

//Public variable
var (
	Bot *discordgo.Session
)

//Start start twitter module
func Start(BotInit *discordgo.Session, cronInit *cron.Cron) {
	Bot = BotInit
	cronInit.AddFunc("@every 0h3m0s", CheckNew)
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
	var w sync.WaitGroup
	for _, GroupData := range engine.GroupData {
		go func(wg *sync.WaitGroup, Group database.Group) {
			defer wg.Done()
			Fanarts, err := CreatePayload(database.GetMembers(Group.ID), Group, Scraper, config.BotConf.LimitConf.TwitterFanart)
			if err != nil {
				log.WithFields(log.Fields{
					"Group": Group.GroupName,
				}).Error(err)
			} else {
				SendFanart(Fanarts, Group)
			}
		}(&w, GroupData)
		time.Sleep(30 * time.Second)
	}
	w.Wait()
}
