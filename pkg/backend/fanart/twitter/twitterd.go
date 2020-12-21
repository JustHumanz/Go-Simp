package twitter

import (
	"math/rand"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"

	engine "github.com/JustHumanz/Go-simp/tools/engine"
)

//CheckNew Check new fanart
func CheckNew() {
	Scraper := twitterscraper.New()
	Scraper.SearchLive(true)
	err := Scraper.SetProxy(config.MultiTOR)
	if err != nil {
		log.Error(err)
	}
	for _, GroupData := range engine.GroupData {
		wg := new(sync.WaitGroup)
		MembersData := database.GetMembers(GroupData.ID)
		for _, MemberData := range MembersData {
			wg.Add(1)
			go func(Member database.Member, Group database.Group, wg *sync.WaitGroup) {
				defer wg.Done()
				if Member.TwitterHashtags != "" {
					if Member.Name != "Kaichou" {
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

						newfanart.CurlTwitter().SendNude()
					}
				} else {
					log.Info(Member.EnName + " don't have twitter hashtag")
				}
			}(MemberData, GroupData, wg)
			if len(MembersData) > 7 {
				time.Sleep(time.Duration(rand.Intn(config.RandomSleep-400)+400) * time.Millisecond)
			}
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
