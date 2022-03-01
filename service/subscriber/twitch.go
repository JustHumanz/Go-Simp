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
	for _, Group := range *Payload {
		for _, Name := range Group.Members {
			if Name.TwitchName != "" && Name.Active() {
				res, err := TwitchClient.GetUsers(&helix.UsersParams{Logins: []string{Name.TwitchName}})
				if err != nil {
					log.Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
				}
				TotalFollowers := 0
				TotalViwers := 0
				for _, v := range res.Data.Users {
					tmp, err := TwitchClient.GetUsersFollows(&helix.UsersFollowsParams{ToID: v.ID})
					if err != nil {
						log.Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}
					TotalFollowers = tmp.Data.Total
					TotalViwers = v.ViewCount
					break
				}

				TwitchFollowDB, err := Name.GetSubsCount()
				if err != nil {
					log.Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
					break
				}

				SendNotif := func(SubsCount, Viwers string) {
					err = Name.RemoveSubsCache()
					if err != nil {
						log.Error(err)
					}

					Color, err := engine.GetColor(config.TmpDir, Name.TwitchAvatar)
					if err != nil {
						log.Error(err)
					}

					SendNude(engine.NewEmbed().
						SetAuthor(Group.GroupName, Group.IconURL, "https://twitch.tv/"+Name.TwitchName).
						SetTitle(engine.FixName(Name.EnName, Name.JpName)).
						SetThumbnail(config.TwitchIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Name.TwitchAvatar).
						AddField("Twitch viwers count", Viwers).
						InlineAllFields().
						SetURL("https://twitch.tv/"+Name.TwitchName).
						SetColor(Color).MessageEmbed, Group, Name)

				}
				if TotalFollowers != TwitchFollowDB.TwitchFollow {

					log.WithFields(log.Fields{
						"Past Twitch Follower":    TwitchFollowDB.TwitchFollow,
						"Current Twitch Follower": TotalFollowers,
						"Vtuber":                  Name.Name,
					}).Info("Update Twitch Follower")

					err := TwitchFollowDB.SetMember(Name).SetGroup(Group).
						UpdateTwitchFollowes(TotalFollowers).
						UpdateTwitchViewers(TotalViwers).
						UpdateState(config.TwitchLive).
						UpdateSubs()
					if err != nil {
						log.Error(err)
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
