package subscriber

import (
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func CheckTwFollowCount() {
	for _, Group := range engine.GroupData {
		for _, Name := range database.GetName(Group.ID) {
			if Name.TwitterName != "" {
				Twitter := Name.GetTwitterFollow()
				TwFollowDB := Name.GetSubsCount()
				SendNotif := func(SubsCount, Tweets string) {
					Avatar := strings.Replace(Twitter.ProfileImageURLHTTPS, "_normal.jpg", ".jpg", -1)
					Color, err := engine.GetColor("/tmp/bili.tmp", Avatar)
					if err != nil {
						log.Error(err)
					}
					SendNude(engine.NewEmbed().
						SetAuthor(Group.NameGroup, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
						SetTitle(engine.FixName(Name.EnName, Name.JpName)).
						SetThumbnail(config.TwitterIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Avatar).
						AddField("Tweets Count", Tweets).
						InlineAllFields().
						SetURL("https://twitter.com/"+Name.TwitterName).
						SetColor(Color).MessageEmbed, Group)
				}
				if Twitter.FollowersCount != TwFollowDB.TwFollow {
					if Twitter.FollowersCount >= 1000000 {
						for i := 0; i < 10000001; i += 1000000 {
							if i == Twitter.FollowersCount {
								SendNotif(strconv.Itoa(i), strconv.Itoa(Twitter.StatusesCount))
							}
						}
					} else if Twitter.FollowersCount >= 100000 {
						for i := 0; i < 1000001; i += 100000 {
							if i == Twitter.FollowersCount {
								SendNotif(strconv.Itoa(i), strconv.Itoa(Twitter.StatusesCount))
							}
						}
					} else if Twitter.FollowersCount >= 10000 {
						for i := 0; i < 100001; i += 10000 {
							if i == Twitter.FollowersCount {
								SendNotif(strconv.Itoa(i), strconv.Itoa(Twitter.StatusesCount))
							}
						}
					} else if Twitter.FollowersCount >= 1000 {
						for i := 0; i < 10001; i += 1000 {
							if i == Twitter.FollowersCount {
								SendNotif(strconv.Itoa(i), strconv.Itoa(Twitter.StatusesCount))
							}
						}
					}
				}
				log.WithFields(log.Fields{
					"Past Twitter Follower":    TwFollowDB.TwFollow,
					"Current Twitter Follower": Twitter.FollowersCount,
					"Vtuber":                   Name.EnName,
				}).Info("Update Twitter Follower")

				TwFollowDB.UptwFollow(Twitter.FollowersCount).
					UpdateSubs("tw")
			}
		}
	}
}
