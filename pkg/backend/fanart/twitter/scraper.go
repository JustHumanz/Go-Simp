package twitter

import (
	"context"

	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

//CurlTwitter start curl twitter api
func (Data *TwitterFanart) CurlTwitter() *TwitterFanart {
	for tweet := range Data.Scraper.SearchTweets(context.Background(), Data.Member.TwitterHashtags+" -filter:retweets -filter:replies", Data.Limit) {
		if tweet.Error != nil {
			log.Error(tweet.Error)
		}

		if Data.Member.CheckMemberFanart(tweet) {
			Data.Fanart = append(Data.Fanart, tweet)
		}
	}
	return Data
}

func (Data *TwitterFanart) RemoveDuplicate() *TwitterFanart {
	keys := make(map[*twitterscraper.Result]bool)
	list := []*twitterscraper.Result{}
	for _, entry := range Data.Fanart {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	Data.Fanart = list
	return Data
}
