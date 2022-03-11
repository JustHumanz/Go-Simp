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
		for _, Member := range Group.Members {
			if Member.TwitterName != "" && Member.Active() {
				Twitter, err := engine.GetTwitterFollow(Member.TwitterName)
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
					continue
				}

				TwFollowDB, err := Member.GetSubsCount()
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
					continue
				}

				SendNotif := func(SubsCount, Tweets, NextCount string, nextdt, score int64) {
					err = Member.RemoveSubsCache()
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					log.WithFields(log.Fields{
						"Vtuber":             Member.Name,
						"Group":              Group.GroupName,
						"TimeSampPrediction": nextdt,
						"Score":              score,
					}).Info("Congratulation for " + SubsCount + " subscriber")

					Avatar := strings.Replace(Twitter.Avatar, "_normal.jpg", ".jpg", -1)
					Color, err := engine.GetColor(config.TmpDir, Avatar)
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22Twitter%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"

					if nextdt != 0 && score != 0 {
						datePredic := time.Unix(nextdt, 0)
						datePredicStr := fmt.Sprintf("%s/%d", datePredic.Month().String(), datePredic.Day())

						SendNude(engine.NewEmbed().
							SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Member.TwitterName).
							SetTitle(engine.FixName(Member.EnName, Member.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+SubsCount+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", Tweets).
							InlineAllFields().
							SetURL("https://twitter.com/"+Member.TwitterName).
							AddField("Graph", Graph).
							AddField("Next Milestone Prediction", datePredicStr+" - "+NextCount).
							SetFooter("Prediction Score "+strconv.Itoa(int(score))+"%").
							SetColor(Color).MessageEmbed, Group, Member)
					} else {
						SendNude(engine.NewEmbed().
							SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Member.TwitterName).
							SetTitle(engine.FixName(Member.EnName, Member.JpName)).
							SetThumbnail(config.TwitterIMG).
							SetDescription("Congratulation for "+SubsCount+" followers").
							SetImage(Avatar).
							AddField("Tweets Count", Tweets).
							InlineAllFields().
							SetURL("https://twitter.com/"+Member.TwitterName).
							AddField("Graph", Graph).
							SetColor(Color).MessageEmbed, Group, Member)
					}

				}

				if Twitter.FollowersCount != TwFollowDB.TwFollow {

					log.WithFields(log.Fields{
						"Past Twitter Follower":    TwFollowDB.TwFollow,
						"Current Twitter Follower": Twitter.FollowersCount,
						"Vtuber":                   Member.Name,
					}).Info("Update Twitter Follower")

					err := TwFollowDB.SetMember(Member).SetGroup(Group).
						UpdateTwitterFollowes(Twitter.FollowersCount).
						UpdateState(config.TwitterArt).
						UpdateSubs()
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					if Twitter.FollowersCount >= 1000000 {
						for i := 0; i < 10000001; i += 1000000 {
							if i == Twitter.FollowersCount {
								NextCount := Twitter.FollowersCount + 1000000
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Member.Name)
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Group.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)

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
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Member.Name)
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Group.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)

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
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Member.Name)
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Group.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)

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
								dt, sc, err := SubsPreDick(NextCount, "Twitter", Member.Name)
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Group.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)

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
				}

				bin, err := TwFollowDB.MarshalBinary()
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
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
