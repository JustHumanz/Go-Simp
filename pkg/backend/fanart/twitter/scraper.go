package twitter

import (
	"context"
	"errors"
	"strings"

	"github.com/JustHumanz/Go-simp/tools/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func CreatePayload(Data []database.Member, Group database.Group, Scraper *twitterscraper.Scraper, Limit int) ([]Fanart, error) {
	var (
		Hashtags []string
		Fanarts  []Fanart
	)
	CurlTwitter := func(Hashtags []string) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   Group.NameGroup,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), strings.Join(Hashtags, " OR ")+" -filter:replies -filter:retweets", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			//Find hashtags members
			for _, MemberHashtag := range Data {
				for _, TweetHashtag := range tweet.Hashtags {
					if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterHashtags) {
						if MemberHashtag.CheckMemberFanart(tweet) {
							Fanarts = append(Fanarts, Fanart{
								Member: MemberHashtag,
								Tweet:  tweet,
							})
						}
					}
				}
			}
		}
	}
	if len(Data) > 7 {
		for i, Member := range Data {
			Hashtags = append(Hashtags, Member.TwitterHashtags)
			if i == len(Data)/2 || i == len(Data)-1 {
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
		return []Fanart{}, errors.New("No new fanart")
	}
}

type Fanart struct {
	Member database.Member
	Tweet  *twitterscraper.Result
}
