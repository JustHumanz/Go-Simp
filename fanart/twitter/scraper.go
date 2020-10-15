package twitter

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func Changebearer() {
	twbearerold := twbearer
	for _, Token := range config.TwitterToken {
		if Token != twbearer {
			twbearer = Token
		}
	}
	log.WithFields(log.Fields{
		"Old": twbearerold,
		"New": twbearer,
	}).Info("Change Twitter bearer")
}

//Start get data from twitter by group name
func CurlTwitter(Group int64) {
	var (
		Hashtags  []string
		TWdata    TwitterStruct
		GroupName string
		body      []byte
		err       error
	)
	for _, hashtag := range database.GetHashtag(Group) {
		if hashtag.TwitterHashtags != "" {
			Hashtags = append(Hashtags, hashtag.TwitterHashtags)
			GroupName = hashtag.GroupName
		}
	}

	log.WithFields(log.Fields{
		"GroupName": GroupName,
	}).Info(strings.Join(Hashtags, " OR "))

	url := "https://api.twitter.com/1.1/search/tweets.json?q=" + url.QueryEscape(strings.Join(Hashtags, " OR ")+" "+"filter:links -filter:replies filter:media -filter:retweets")
	body, err = engine.Curl(url, []string{"Authorization", "Bearer " + twbearer})
	if err != nil {
		if strings.HasPrefix(err.Error(), "401 Unauthorized") {
			Changebearer()
			body, err = engine.CoolerCurl(url, []string{"Authorization", "Bearer " + twbearer})
			if err != nil {
				log.Error(err, string(body))
				return
			}
		} else {
			log.Error(err, string(body))
			return
		}
	}
	err = json.Unmarshal(body, &TWdata)
	if err != nil {
		log.Error(err)
	}
	wg := new(sync.WaitGroup)
	for _, Statuses := range TWdata.CheckNew() {
		wg.Add(1)
		go Statuses.CheckHashTag(database.GetHashtag(Group), wg)
	}
	wg.Wait()
}
