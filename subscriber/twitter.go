package subscriber

import (
	"strconv"
	"time"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func CheckTwFollowCount() {
	for _, Group := range Data {
		for _, Name := range database.GetName(Group.ID) {
			Twitter := Name.GetTwitterFollow()
			if Twitter.FollowersCount <= 100000 {
				for i := 0; i < 1000001; i += 100000 {
					if i == Twitter.FollowersCount && Twitter.FollowersCount != 0 {
						Avatar := engine.GetUserAvatar(Name.TwitterName)
						Color, err := engine.GetColor("/tmp/bili.tmp", Avatar)
						if err != nil {
							log.Error(err)
						}
						SendNude(engine.NewEmbed().
							SetAuthor(Group.NameGroup, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
							SetTitle(engine.FixName(Name.EnName, Name.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", strconv.Itoa(Twitter.FollowersCount)).
							InlineAllFields().
							SetURL("https://twitter.com/"+Name.TwitterName).
							SetColor(Color).MessageEmbed, Group)
					}
				}
			} else {
				for i := 0; i < 10001; i += 1000 {
					if i == Twitter.FollowersCount && Twitter.FollowersCount != 0 {
						Avatar := engine.GetUserAvatar(Name.TwitterName)
						Color, err := engine.GetColor("/tmp/bili.tmp", Avatar)
						if err != nil {
							log.Error(err)
						}
						SendNude(engine.NewEmbed().
							SetAuthor(Group.NameGroup, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
							SetTitle(engine.FixName(Name.EnName, Name.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", strconv.Itoa(Twitter.FollowersCount)).
							InlineAllFields().
							SetURL("https://twitter.com/"+Name.TwitterName).
							SetColor(Color).MessageEmbed, Group)
					}
				}
			}

			TwFollowDB := Name.GetSubsCount()
			log.WithFields(log.Fields{
				"Past Twitter Follower":    TwFollowDB.TwFollow,
				"Current Twitter Follower": Twitter.FollowersCount,
				"Vtuber":                   Name.EnName,
			}).Info("Update Twitter Follower")

			TwFollowDB.TwFollow = Twitter.FollowersCount
			TwFollowDB.UpdateSubs("tw")
			time.Sleep(1 * time.Second)
		}
	}
}
