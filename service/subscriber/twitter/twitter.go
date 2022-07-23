package main

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	Bot         *discordgo.Session
	configfile  config.ConfigFile
	gRCPconn    pilot.PilotServiceClient
	ServiceUUID = uuid.New().String()
	Agency      []database.Group
)

const (
	ServiceName = config.SubscriberService + "Twitter"
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
	database.Start(configfile)

	c := cron.New()
	c.Start()

	c.AddFunc(config.TwitterFollowers, CheckTwitter)
	log.Info("Add twitter followers to cronjob")

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

func CheckTwitter() {
	for _, Group := range Agency {
		for _, Member := range Group.Members {
			if Member.TwitterName != "" && Member.Active() {
				Twitter, err := engine.GetTwitterFollow(Member.TwitterName)
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
					continue
				}

				TwFollowDB, err := Member.GetSubsCount()
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
					continue
				}

				TwFollowDB.SetMember(Member).SetGroup(Group).UpdateState(config.TwitterArt)

				SendNotif := func(SubsCount, Tweets string) {
					log.WithFields(log.Fields{
						"Vtuber": Member.Name,
						"Group":  Group.GroupName,
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
					engine.SubsSendEmbed(engine.NewEmbed().
						SetAuthor(Group.GroupName, Group.IconURL, "https://twitter.com/"+Member.TwitterName).
						SetTitle(engine.FixName(Member.EnName, Member.JpName)).
						SetThumbnail(config.TwitterIMG).
						SetDescription("Congratulation for "+SubsCount+" followers").
						SetImage(Avatar).
						AddField("Tweets Count", Tweets).
						InlineAllFields().
						SetURL("https://twitter.com/"+Member.TwitterName).
						AddField("Graph", Graph).
						SetColor(Color).MessageEmbed, Group, Member, Bot)

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
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
								)
							}
						}
					} else if Twitter.FollowersCount >= 100000 {
						for i := 0; i < 1000001; i += 100000 {
							if i == Twitter.FollowersCount {
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
								)
							}
						}
					} else if Twitter.FollowersCount >= 10000 {
						for i := 0; i < 100001; i += 10000 {
							if i == Twitter.FollowersCount {
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
								)
							}
						}
					} else if Twitter.FollowersCount >= 1000 {
						for i := 0; i < 10001; i += 1000 {
							if i == Twitter.FollowersCount {
								SendNotif(
									engine.NearestThousandFormat(float64(i)),
									engine.NearestThousandFormat(float64(Twitter.TweetsCount)),
								)
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
