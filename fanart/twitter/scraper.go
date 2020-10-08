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

//Start get data from twitter by group name
func CurlTwitter(Group int64) {
	var (
		Hashtags  []string
		TWdata    TwitterStruct
		GroupName string
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

	query := url.QueryEscape(strings.Join(Hashtags, " OR ") + " " + "filter:links -filter:replies filter:media -filter:retweets")
	body, err := engine.Curl("https://api.twitter.com/1.1/search/tweets.json?q="+query, []string{"Authorization", "Bearer " + config.TwitterToken})
	if err != nil {
		log.Error(err, string(body))
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
