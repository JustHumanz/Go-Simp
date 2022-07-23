package main

import (
	"context"
	"encoding/json"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	configfile   config.ConfigFile
	gRCPconn     pilot.PilotServiceClient
	TwitchClient *helix.Client
	ServiceUUID  = uuid.New().String()
	Agency       []database.Group
)

const (
	ServiceName = config.SubscriberService + "Twitch"
)

//Init service
func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

}

func main() {
	//Get config file from pilot
	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		Service:     ServiceName,
		ServiceUUID: ServiceUUID,
	})
	if err != nil {
		if configfile.Discord != "" {
			pilot.ReportDeadService(err.Error(), ServiceName)
		}
		log.Error("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Error(err)
	}

	configfile.InitConf()
	Bot = engine.StartBot(false)
	TwitchClient = engine.GetTwitchTkn()

	database.Start(configfile)

	resp, err := TwitchClient.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Panic(err)
	}

	TwitchClient.SetAppAccessToken(resp.Data.AccessToken)

	c := cron.New()
	c.Start()
	c.AddFunc(config.TwitchFollowers, CheckTwitch)
	log.Info("Add twitch followers to cronjob")

	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)

	hostname := engine.GetHostname()

	go func() {
		tmp, err := gRCPconn.GetAgencyPayload(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Refresh payload",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(tmp.AgencyVtubers, &Agency)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(1 * time.Hour)
	}()
	runfunc.Run(Bot)
}

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

					Color, err := engine.GetColor(config.TmpDir, Member.TwitchAvatar)
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}

					engine.SubsSendEmbed(engine.NewEmbed().
						SetAuthor(Group.GroupName, Group.IconURL, "https://twitch.tv/"+Member.TwitchName).
						SetTitle(engine.FixName(Member.EnName, Member.JpName)).
						SetThumbnail(config.TwitchIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Member.TwitchAvatar).
						AddField("Twitch viwers count", Viwers).
						InlineAllFields().
						SetURL("https://twitch.tv/"+Member.TwitchName).
						SetColor(Color).MessageEmbed, Group, Member, Bot)

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
