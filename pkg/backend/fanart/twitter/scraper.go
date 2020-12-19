package twitter

import (
	"context"

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
