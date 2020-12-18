package twitter

import (
	"encoding/json"
	"net/url"
	"strings"
	"sync"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	network "github.com/JustHumanz/Go-simp/tools/network"
	log "github.com/sirupsen/logrus"
)

//CurlTwitter start curl twitter api
func CurlTwitter(TwitterPayloads []string, Group int64) {
	var (
		TWdata TwitterStruct
	)
	for _, Payload := range TwitterPayloads {
		var (
			body []byte
			err  error
		)
		url := "https://api.twitter.com/1.1/search/tweets.json?" + Payload
		body, err = network.Curl(url, []string{"Authorization", "Bearer " + config.TwitterToken[0]})
		if err != nil {
			if strings.HasPrefix(err.Error(), "401 Unauthorized") {
				body, err = network.CoolerCurl(url, []string{"Authorization", "Bearer " + config.TwitterToken[0]})
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
}

//CreatePayload create payload from twitter by group name
func CreatePayload(Group int64) []string {
	var (
		Hashtags      []string
		FinalTags     []string
		GroupsHashtag = database.GetHashtag(Group)
		params        = url.Values{}
		counter       = 0
	)
	if len(GroupsHashtag) < 7 && GroupsHashtag != nil {
		GroupName := ""
		for _, hashtag := range database.GetHashtag(Group) {
			if hashtag.TwitterHashtags != "" {
				Hashtags = append(Hashtags, hashtag.TwitterHashtags)
				GroupName = hashtag.GroupName
			}
		}
		params.Add("q", strings.Join(Hashtags, " OR ")+" AND "+"-filter:retweets AND -filter:replies")
		log.WithFields(log.Fields{
			"Hashtags":  strings.Join(Hashtags, " OR "),
			"GroupName": GroupName,
		}).Info("Start Curl new art")

		return []string{params.Encode()}
	} else {
		for z, hashtag := range GroupsHashtag {
			Hashtags = append(Hashtags, hashtag.TwitterHashtags)
			if counter == 7 || z == len(GroupsHashtag)-1 {
				log.WithFields(log.Fields{
					"Hashtags":  strings.Join(Hashtags, " OR "),
					"GroupName": hashtag.GroupName,
				}).Info("Start Curl new art")
				params.Add("q", strings.Join(Hashtags, " OR ")+" AND "+"-filter:retweets AND -filter:replies")
				FinalTags = append(FinalTags, params.Encode())
				Hashtags = nil
				counter = 0
			}
			counter++
		}
		return FinalTags
	}
}
