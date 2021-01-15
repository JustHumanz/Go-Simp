package twitter

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/JustHumanz/Go-simp/tools/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func CreatePayload(Data []database.Member, Group database.Group, Scraper *twitterscraper.Scraper, Limit int) ([]Fanart, error) {
	var (
		Hashtags []string
		Fanarts  []Fanart
		rdb      = database.FanartCache
	)

	CurlTwitter := func(Hashtags []string, GroupData database.Group) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   GroupData.GroupName,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), "("+strings.Join(Hashtags, " OR ")+") AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			//Find hashtags members
			val, err := rdb.Get(context.Background(), tweet.ID).Result()
			if val == "" || err != nil {
				for _, MemberHashtag := range Data {
					for _, TweetHashtag := range tweet.Hashtags {
						if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterHashtags) && !tweet.IsQuoted && !tweet.IsReply && MemberHashtag.Name != "Kaichou" {
							New, err := MemberHashtag.CheckMemberFanart(tweet)
							if err != nil {
								log.Error(err)
							}
							if New {
								Fanarts = append(Fanarts, Fanart{
									Member: MemberHashtag,
									Tweet:  tweet,
								})
							}
							bit, err := json.Marshal(tweet)
							if err != nil {
								log.Error(err)
							}
							err = rdb.Set(context.Background(), tweet.ID, bit, 30*time.Minute).Err()
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}
		}
	}

	if len(Data) > 7 {
		for i, Member := range Data {
			if Member.TwitterHashtags != "" {
				Hashtags = append(Hashtags, Member.TwitterHashtags)
				if (i%5 == 0) || (i == len(Data)-1) {
					CurlTwitter(Hashtags, Group)
					Hashtags = nil
				}
			}
		}
	} else {
		for _, Member := range Data {
			if Member.TwitterHashtags != "" {
				Hashtags = append(Hashtags, Member.TwitterHashtags)
			}
		}
		CurlTwitter(Hashtags, Group)
	}

	if len(Fanarts) > 0 {
		return Fanarts, nil
	} else {
		return []Fanart{}, errors.New("Still same")
	}
}

type Fanart struct {
	Member database.Member
	Tweet  *twitterscraper.Result
}
