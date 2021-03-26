package twitter

import (
	"context"
	"errors"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func CreatePayload(Group database.Group, Scraper *twitterscraper.Scraper, Limit int) ([]database.DataFanart, error) {
	var (
		Hashtags     []string
		LewdHashtags []string
		Fanarts      []database.DataFanart
	)

	CurlTwitter := func(Hashtags []string) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   Group.GroupName,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), "("+strings.ReplaceAll(strings.Join(Hashtags, " OR "), " ", "")+") AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			for _, MemberHashtag := range Group.Members {
				for _, TweetHashtag := range tweet.Hashtags {
					if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterHashtags) && !tweet.IsQuoted && !tweet.IsReply && MemberHashtag.Name != "Kaichou" && len(tweet.Photos) > 0 {
						TweetArt := database.DataFanart{
							PermanentURL: tweet.PermanentURL,
							Author:       "@" + tweet.Username,
							AuthorAvatar: engine.GetAuthorAvatar(tweet.Username),
							TweetID:      tweet.ID,
							Text:         RemoveTwitterShortLink(tweet.Text),
							Photos:       tweet.Photos,
							Likes:        tweet.Likes,
							Member:       MemberHashtag,
							Group:        Group,
							State:        config.TwitterArt,
						}
						if tweet.Videos != nil {
							TweetArt.Videos = tweet.Videos[0].Preview
						}

						New, err := TweetArt.CheckTweetFanArt()
						if err != nil {
							log.Error(err)
						}

						if New {
							Fanarts = append(Fanarts, TweetArt)
						}
					}
				}
			}
		}
	}

	CurlTwitterLewd := func(Hashtags []string) {
		log.WithFields(log.Fields{
			"Hashtag": strings.Join(Hashtags, " OR "),
			"Group":   Group.GroupName,
		}).Info("Start curl twitter")
		for tweet := range Scraper.SearchTweets(context.Background(), "("+strings.ReplaceAll(strings.Join(Hashtags, " OR "), " ", "")+") AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)", Limit) {
			if tweet.Error != nil {
				log.Error(tweet.Error)
			}
			for _, MemberHashtag := range Group.Members {
				for _, TweetHashtag := range tweet.Hashtags {
					if strings.ToLower("#"+TweetHashtag) == strings.ToLower(MemberHashtag.TwitterLewd) && !tweet.IsQuoted && !tweet.IsReply && len(tweet.Photos) > 0 {
						TweetArt := database.DataFanart{
							PermanentURL: tweet.PermanentURL,
							Author:       "@" + tweet.Username,
							AuthorAvatar: engine.GetAuthorAvatar(tweet.Username),
							TweetID:      tweet.ID,
							Text:         RemoveTwitterShortLink(tweet.Text),
							Photos:       tweet.Photos,
							Likes:        tweet.Likes,
							Member:       MemberHashtag,
							Group:        Group,
							State:        config.TwitterArt,
							Lewd:         true,
						}

						if tweet.Videos != nil {
							TweetArt.Videos = tweet.Videos[0].Preview
						}

						New, err := TweetArt.CheckTweetFanArt()
						if err != nil {
							log.Error(err)
						}

						if New {
							Fanarts = append(Fanarts, TweetArt)
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
					CurlTwitter(Hashtags)
					Hashtags = nil
				}
			}

			if lewd && Member.TwitterLewd != "" {
				LewdHashtags = append(LewdHashtags, Member.TwitterLewd)
				if (i%5 == 0) || (i == len(Group.Members)-1) {
					CurlTwitterLewd(LewdHashtags)
					LewdHashtags = nil
				}
			}
		}
	} else {
		for _, Member := range Group.Members {
			if Member.TwitterHashtags != "" {
				Hashtags = append(Hashtags, Member.TwitterHashtags)
			}

			if lewd && Member.TwitterLewd != "" {
				LewdHashtags = append(LewdHashtags, Member.TwitterLewd)
			}
		}

		CurlTwitter(Hashtags)

		if lewd {
			CurlTwitterLewd(LewdHashtags)
		}
	}

	if len(Fanarts) > 0 {
		return Fanarts, nil
	} else {
		return []database.DataFanart{}, errors.New("Still same")
	}
}
