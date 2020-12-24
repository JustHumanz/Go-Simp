package twitter

import (
	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"

	engine "github.com/JustHumanz/Go-simp/tools/engine"
)

//CheckNew Check new fanart
func CheckNew() {
	Scraper := twitterscraper.New()
	Scraper.SetSearchMode(twitterscraper.SearchLatest)
	err := Scraper.SetProxy(config.MultiTOR)
	if err != nil {
		log.Error(err)
	}
	for _, GroupData := range engine.GroupData {
		Fanarts, err := CreatePayload(database.GetHashtag(GroupData.ID), GroupData, Scraper, config.FanartLimit)
		if err != nil {
			log.WithFields(log.Fields{
				"Group": GroupData.NameGroup,
			}).Error(err)
		} else {
			SendFanart(Fanarts, GroupData)
		}
	}
}
