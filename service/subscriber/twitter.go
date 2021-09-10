package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckTwitter() {
	for _, Group := range *Payload {
		for _, Name := range Group.Members {
			if Name.TwitterName != "" && Name.Active() {
				Twitter, err := Name.GetTwitterFollow()
				if err != nil {
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
					log.Error(err)
					break
				}
				TwFollowDB, err := Name.GetSubsCount()
				if err != nil {
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
					log.Error(err)
					break
				}
				SendNotif := func(SubsCount, Tweets, NextCount string, nextdt, score int64) {
					err = Name.RemoveSubsCache()
					if err != nil {
						log.Error(err)
					}

					Avatar := strings.Replace(Twitter.Avatar, "_normal.jpg", ".jpg", -1)
					Color, err := engine.GetColor(config.TmpDir, Avatar)
					if err != nil {
						log.Error(err)
					}

					Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22Twitter%22%2C%20vtuber%3D%22" + Name.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"

					if nextdt != 0 && score != 0 {
						datePredic := time.Unix(nextdt, 0)
						datePredicStr := fmt.Sprintf("%s/%d", datePredic.Month().String(), datePredic.Day())

						SendNude(engine.NewEmbed().
							SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
							SetTitle(engine.FixName(Name.EnName, Name.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+SubsCount+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", Tweets).
							InlineAllFields().
							SetURL("https://twitter.com/"+Name.TwitterName).
							AddField("Graph", Graph).
							AddField("Next Milestone Prediction", datePredicStr+" - "+NextCount).
							SetFooter("Prediction Score "+strconv.Itoa(int(score))+"%").
							SetColor(Color).MessageEmbed, Group, Name)
					} else {
						SendNude(engine.NewEmbed().
							SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Name.TwitterName).
							SetTitle(engine.FixName(Name.EnName, Name.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+SubsCount+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", Tweets).
							InlineAllFields().
							SetURL("https://twitter.com/"+Name.TwitterName).
							AddField("Graph", Graph).
							SetColor(Color).MessageEmbed, Group, Name)
					}

				}
				if Twitter.FollowersCount != TwFollowDB.TwFollow {
					if Twitter.FollowersCount >= 1000000 {
						for i := 0; i < 10000001; i += 1000000 {
							if i == Twitter.FollowersCount {
								NextCount := Twitter.FollowersCount + 1000000
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Name.Name)
								if err != nil {
									log.Error(err)
									gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
										Message: err.Error(),
										Service: ModuleState,
									})
								}
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
									engine.NearestThousandFormat(float64(NextCount)),
									dt, sc)
							}
						}
					} else if Twitter.FollowersCount >= 100000 {
						for i := 0; i < 1000001; i += 100000 {
							if i == Twitter.FollowersCount {
								NextCount := Twitter.FollowersCount + 100000
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Name.Name)
								if err != nil {
									log.Error(err)
									gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
										Message: err.Error(),
										Service: ModuleState,
									})
								}
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
									engine.NearestThousandFormat(float64(NextCount)),
									dt, sc)
							}
						}
					} else if Twitter.FollowersCount >= 10000 {
						for i := 0; i < 100001; i += 10000 {
							if i == Twitter.FollowersCount {
								NextCount := Twitter.FollowersCount + 10000
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Name.Name)
								if err != nil {
									log.Error(err)
									gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
										Message: err.Error(),
										Service: ModuleState,
									})
								}
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
									engine.NearestThousandFormat(float64(NextCount)),
									dt, sc)
							}
						}
					} else if Twitter.FollowersCount >= 1000 {
						for i := 0; i < 10001; i += 1000 {
							if i == Twitter.FollowersCount {
								NextCount := Twitter.FollowersCount + 1000
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Name.Name)
								if err != nil {
									log.Error(err)
									gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
										Message: err.Error(),
										Service: ModuleState,
									})
								}
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
									engine.NearestThousandFormat(float64(NextCount)),
									dt, sc)
							}
						}
					}

					log.WithFields(log.Fields{
						"Past Twitter Follower":    TwFollowDB.TwFollow,
						"Current Twitter Follower": Twitter.FollowersCount,
						"Vtuber":                   Name.EnName,
					}).Info("Update Twitter Follower")

					TwFollowDB.SetMember(Name).SetGroup(Group).
						UptwFollow(Twitter.FollowersCount).
						UpdateState(config.TwitterArt).
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
