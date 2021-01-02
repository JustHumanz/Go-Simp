package twitter

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/go-redis/redis/v8"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func CreatePayload(Data []database.Member, Group database.Group, Scraper *twitterscraper.Scraper, Limit int) ([]Fanart, error) {
	var (
		Hashtags []string
		Fanarts  []Fanart
		ctx      = context.Background()
		rdb      = database.FanartCache
	)
	CurlTwitter := func(Hashtags []string) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   Group.GroupName,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), "("+strings.Join(Hashtags, " OR ")+") AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			//Find hashtags members
			for _, MemberHashtag := range Data {
				for _, TweetHashtag := range tweet.Hashtags {
					if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterHashtags) && MemberHashtag.EnName != "Kaichou" {
						_, err := rdb.Get(ctx, tweet.ID).Result()
						if err == redis.Nil {
							if MemberHashtag.CheckMemberFanart(tweet) {
								Fanarts = append(Fanarts, Fanart{
									Member: MemberHashtag,
									Tweet:  tweet,
								})
							}
							err := rdb.Set(ctx, tweet.ID, MemberHashtag.Name, 30*time.Minute).Err()
							if err != nil {
								log.Error(err)
							}
						} else if err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	}

	if len(Data) > 7 {
		for i, Member := range Data {
			Hashtags = append(Hashtags, Member.TwitterHashtags)
			if (i%5 == 0) || (i == len(Data)-1) {
				CurlTwitter(Hashtags)
				Hashtags = nil
			}
		}
	} else {
		for _, Member := range Data {
			Hashtags = append(Hashtags, Member.TwitterHashtags)
		}
		CurlTwitter(Hashtags)
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
