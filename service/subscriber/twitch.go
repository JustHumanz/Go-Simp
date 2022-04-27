package main

import (
	"context"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

func CheckTwitch() {
	for _, Group := range Agency {
		for _, Member := range Group.Members {
			if Member.TwitchName != "" && Member.Active() {
				res, err := TwitchClient.GetUsers(&helix.UsersParams{Logins: []string{Member.TwitchName}})
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message:     err.Error(),
						Service:     ServiceName,
						ServiceUUID: ServiceUUID,
					})
				}
				TotalFollowers := 0
				TotalViwers := 0
				for _, v := range res.Data.Users {
					tmp, err := TwitchClient.GetUsersFollows(&helix.UsersFollowsParams{ToID: v.ID})
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message:     err.Error(),
							Service:     ServiceName,
							ServiceUUID: ServiceUUID,
						})
					}
					TotalFollowers = tmp.Data.Total
					TotalViwers = v.ViewCount
					break
				}

				TwitchFollowDB, err := Member.GetSubsCount()
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message:     err.Error(),
						Service:     ServiceName,
						ServiceUUID: ServiceUUID,
					})
					break
				}

				SendNotif := func(SubsCount, Viwers string) {
					err = Member.RemoveSubsCache()
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					Color, err := engine.GetColor(config.TmpDir, Member.TwitchAvatar)
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					SendNude(engine.NewEmbed().
						SetAuthor(Group.GroupName, Group.IconURL, "https://twitch.tv/"+Member.TwitchName).
						SetTitle(engine.FixName(Member.EnName, Member.JpName)).
						SetThumbnail(config.TwitchIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Member.TwitchAvatar).
						AddField("Twitch viwers count", Viwers).
						InlineAllFields().
						SetURL("https://twitch.tv/"+Member.TwitchName).
						SetColor(Color).MessageEmbed, Group, Member)

				}
				if TotalFollowers != TwitchFollowDB.TwitchFollow {

					log.WithFields(log.Fields{
						"Past Twitch Follower":    TwitchFollowDB.TwitchFollow,
						"Current Twitch Follower": TotalFollowers,
						"Vtuber":                  Member.Name,
					}).Info("Update Twitch Follower")

					err := TwitchFollowDB.SetMember(Member).SetGroup(Group).
						UpdateTwitchFollowes(TotalFollowers).
						UpdateTwitchViewers(TotalViwers).
						UpdateState(config.TwitchLive).
						UpdateSubs()
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					if TotalFollowers >= 1000000 {
						for i := 0; i < 10000001; i += 1000000 {
							if i == TotalFollowers {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(TotalViwers)))
							}
						}
					} else if TotalFollowers >= 100000 {
						for i := 0; i < 1000001; i += 100000 {
							if i == TotalFollowers {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(TotalViwers)))
							}
						}
					} else if TotalFollowers >= 10000 {
						for i := 0; i < 100001; i += 10000 {
							if i == TotalFollowers {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(TotalViwers)))
							}
						}
					} else if TotalFollowers >= 1000 {
						for i := 0; i < 10001; i += 1000 {
							if i == TotalFollowers {
								SendNotif(engine.NearestThousandFormat(float64(i)), engine.NearestThousandFormat(float64(TotalViwers)))
							}
						}
					}
				}

				bin, err := TwitchFollowDB.MarshalBinary()
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
