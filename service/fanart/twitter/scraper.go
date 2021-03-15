package twitter

import (
	"context"
	"errors"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func CreatePayload(Group database.Group, Scraper *twitterscraper.Scraper, Limit int) ([]Fanart, error) {
	var (
		Hashtags []string
		Fanarts  []Fanart
	)

	CurlTwitter := func(Hashtags []string, GroupData database.Group) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   GroupData.GroupName,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), "("+strings.ReplaceAll(strings.Join(Hashtags, " OR "), " ", "")+") AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			for _, MemberHashtag := range Group.Members {
				for _, TweetHashtag := range tweet.Hashtags {
					if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterHashtags) && !tweet.IsQuoted && !tweet.IsReply && MemberHashtag.Name != "Kaichou" && len(tweet.Photos) > 0 {
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
					}
				}
			}
		}
	}

	if len(Group.Members) > 7 {
		for i, Member := range Group.Members {
			if Member.TwitterHashtags != "" {
				Hashtags = append(Hashtags, Member.TwitterHashtags)
				if (i%5 == 0) || (i == len(Group.Members)-1) {
					CurlTwitter(Hashtags, Group)
					Hashtags = nil
				}
			}
		}
	} else {
		for _, Member := range Group.Members {
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
