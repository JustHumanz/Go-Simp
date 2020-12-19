package twitter

import (
	"sync"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"

	engine "github.com/JustHumanz/Go-simp/tools/engine"
)

//CheckNew Check new fanart
func CheckNew() {
	Scraper := twitterscraper.New()
	err := Scraper.SetProxy(config.MultiTOR)
	if err != nil {
		log.Error(err)
	}
	for _, GroupData := range engine.GroupData {
		wg := new(sync.WaitGroup)
		for _, MemberData := range database.GetMembers(GroupData.ID) {
			wg.Add(1)
			go func(Member database.Member, Group database.Group, wg *sync.WaitGroup) {
				defer wg.Done()
				if Member.TwitterHashtags != "" || Member.Name != "Kaichou" {
					newfanart := TwitterFanart{
						Member:  Member,
						Limit:   5,
						Group:   Group,
						Scraper: Scraper,
					}
					log.WithFields(log.Fields{
						"Name":    Member.EnName,
						"Hashtag": Member.TwitterHashtags,
						"Group":   Group.NameGroup,
					}).Info("Scraping Fanart")

					newfanart.CurlTwitter()
					newfanart.SendNude()
				} else {
					log.Info(Member.EnName + " don't have twitter hashtag")
				}
			}(MemberData, GroupData, wg)
		}
		wg.Wait()
	}
}

type TwitterFanart struct {
	Member  database.Member
	Group   database.Group
	Limit   int
	Fanart  []*twitterscraper.Result
	Scraper *twitterscraper.Scraper
}
