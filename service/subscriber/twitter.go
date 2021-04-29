package main

import (
	"context"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckTwitter() {
	for _, Group := range Payload.VtuberData {
		for _, Name := range Group.Members {
			if Name.TwitterName != "" {
				Twitter, err := Name.GetTwitterFollow()
				if err != nil {
					log.Error(err)
					break
				}
				TwFollowDB, err := Name.GetSubsCount()
				if err != nil {
					log.Error(err)
					break
				}
				SendNotif := func(SubsCount, Tweets string) {
					Avatar := strings.Replace(Twitter.Avatar, "_normal.jpg", ".jpg", -1)
					Color, err := engine.GetColor(config.TmpDir, Avatar)
					if err != nil {
						log.Error(err)
					}
					SendNude(engine.NewEmbed().
						SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
						SetTitle(engine.FixName(Name.EnName, Name.JpName)).
						SetThumbnail(config.TwitterIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Avatar).
						AddField("Tweets Count", Tweets).
						InlineAllFields().
						SetURL("https://twitter.com/"+Name.TwitterName).
						SetColor(Color).MessageEmbed, Group, Name)
				}
				if Twitter.FollowersCount != TwFollowDB.TwFollow {
					if Twitter.FollowersCount >= 1000000 {
						for i := 0; i < 10000001; i += 1000000 {
							if i == Twitter.FollowersCount {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(Twitter.TweetsCount)))
							}
						}
					} else if Twitter.FollowersCount >= 100000 {
						for i := 0; i < 1000001; i += 100000 {
							if i == Twitter.FollowersCount {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(Twitter.TweetsCount)))
							}
						}
					} else if Twitter.FollowersCount >= 10000 {
						for i := 0; i < 100001; i += 10000 {
							if i == Twitter.FollowersCount {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(Twitter.TweetsCount)))
							}
						}
					} else if Twitter.FollowersCount >= 1000 {
						for i := 0; i < 10001; i += 1000 {
							if i == Twitter.FollowersCount {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(Twitter.TweetsCount)))
							}
						}
					}

					log.WithFields(log.Fields{
						"Past Twitter Follower":    TwFollowDB.TwFollow,
						"Current Twitter Follower": Twitter.FollowersCount,
						"Vtuber":                   Name.EnName,
					}).Info("Update Twitter Follower")

					NewSubs := Twitter.FollowersCount - TwFollowDB.TwFollow
					TwFollowDB.SetMember(Name).SetGroup(Group).
						UptwFollow(Twitter.FollowersCount).
						UpdateState(config.TwitterArt).
						AddNewSubs(NewSubs).
						UpdateSubs()

					bin, err := TwFollowDB.MarshalBinary()
					if err != nil {
						log.Error(err)
					}
					if config.GoSimpConf.Metric {
						gRCPconn.MetricReport(context.Background(), &pilot.Metric{
							MetricData: bin,
							State:      config.SubsState,
						})
					}
				}
			}
		}
	}
}
