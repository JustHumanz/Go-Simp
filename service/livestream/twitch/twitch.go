package main

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	TwitchClient *helix.Client
	gRCPconn     pilot.PilotServiceClient
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.TwitchService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//main start twitter module
func main() {
	var (
		configfile config.ConfigFile
	)

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
		log.Fatalln(err)
	}

	configfile.InitConf()
	Bot = engine.StartBot(false)
	TwitchClient = engine.GetTwitchTkn()

	database.Start(configfile)

	resp, err := TwitchClient.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
			Message: err.Error(),
			Service: ServiceName,
		})
		log.Panic(err)
	}

	TwitchClient.SetAppAccessToken(resp.Data.AccessToken)

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkTwcJob struct {
	wg      sync.WaitGroup
	Agency  []database.Group
	Reverse bool
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	Twitch := &checkTwcJob{}

	for {
		log.WithFields(log.Fields{
			"Service": ServiceName,
			"Running": false,
			"UUID":    ServiceUUID,
		}).Info("request for running job")

		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Request",
			ServiceUUID: ServiceUUID,
		})
		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Service": ServiceName,
				"Running": true,
				"UUID":    ServiceUUID,
			}).Info(res.Message)

			Twitch.Agency = engine.UnMarshalPayload(res.VtuberPayload)
			Twitch.Run()

			_, _ = client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
				Service:     ServiceName,
				Message:     "Done",
				ServiceUUID: ServiceUUID,
			})

			log.WithFields(log.Fields{
				"Service": ServiceName,
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info("reporting job was done")
		}
		time.Sleep(1 * time.Minute)
	}
}

func (i *checkTwcJob) Run() {

	Cek := func(Group database.Group, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, Member := range Group.Members {
			if Member.TwitchName != "" && Member.Active() {
				log.WithFields(log.Fields{
					"Agency": Group.GroupName,
					"Vtuber": Member.Name,
				}).Info("Checking Twitch")

				result, err := TwitchClient.GetStreams(&helix.StreamsParams{
					UserLogins: []string{Member.TwitchName},
				})

				if err != nil || result.ErrorMessage != "" {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err, result.ErrorMessage)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message:     err.Error() + " " + result.ErrorMessage,
						Service:     ServiceName,
						ServiceUUID: ServiceUUID,
					})
					return
				}

				ResultDB, err := database.GetTwitch(Member.ID)
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
					return
				}

				ResultDB.AddMember(Member).AddGroup(Group).SetState(config.TwitchLive)

				if len(result.Data.Streams) > 0 {
					for _, Stream := range result.Data.Streams {
						if ResultDB.Status == config.PastStatus && Stream.Type == config.LiveStatus {
							GameResult, err := TwitchClient.GetGames(&helix.GamesParams{
								IDs: []string{Stream.GameID},
							})
							if err != nil || GameResult.ErrorMessage != "" {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err, GameResult.ErrorMessage)
							}

							Stream.ThumbnailURL = strings.Replace(Stream.ThumbnailURL, "{width}", "1280", -1)
							Stream.ThumbnailURL = strings.Replace(Stream.ThumbnailURL, "{height}", "720", -1)

							ResultDB.UpdateStatus(config.LiveStatus).
								UpdateViewers(strconv.Itoa(Stream.ViewerCount)).
								UpdateThumbnail(Stream.ThumbnailURL).
								SetState(config.TwitchLive).
								UpdateSchdule(Stream.StartedAt)

							if len(GameResult.Data.Games) > 0 {
								ResultDB.UpdateGame(GameResult.Data.Games[0].Name)
							} else {
								ResultDB.UpdateGame("-")
							}

							err = ResultDB.UpdateTwitch()
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							if config.GoSimpConf.Metric {
								bit, err := ResultDB.MarshalBinary()
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Group.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}
								gRCPconn.MetricReport(context.Background(), &pilot.Metric{
									MetricData: bit,
									State:      config.LiveStatus,
								})
							}

							engine.SendLiveNotif(ResultDB, Bot)

							log.WithFields(log.Fields{
								"Group":      Group.GroupName,
								"VtuberName": Member.Name,
							}).Info("Change Twitch status to Live")
						} else if Stream.Type == config.LiveStatus && ResultDB.Status == config.LiveStatus {
							log.WithFields(log.Fields{
								"Group":      Group.GroupName,
								"VtuberName": Member.Name,
								"Viewers":    Stream.ViewerCount,
							}).Info("Update Viewers")

							ResultDB.UpdateViewers(strconv.Itoa(Stream.ViewerCount)).UpdateTwitch()
						}
					}
				} else if ResultDB.Status == config.LiveStatus && len(result.Data.Streams) == 0 {
					ResultDB.UpdateEnd(time.Now()).UpdateStatus(config.PastStatus)
					err = ResultDB.UpdateTwitch()
					if err != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)
					}
					log.WithFields(log.Fields{
						"Group":      Group.GroupName,
						"VtuberName": Member.Name,
					}).Info("Change Twitch status to Past")

					engine.RemoveEmbed("Twitch"+Member.TwitchName, Bot)

					if config.GoSimpConf.Metric {
						bit, err := ResultDB.MarshalBinary()
						if err != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(err)
						}
						gRCPconn.MetricReport(context.Background(), &pilot.Metric{
							MetricData: bit,
							State:      config.PastStatus,
						})
					}
				}
			}

		}
	}
	if i.Reverse {
		for j := len(i.Agency) - 1; j >= 0; j-- {
			i.wg.Add(1)

			Grp := i.Agency
			go Cek(Grp[j], &i.wg)
		}
		i.Reverse = false

	} else {
		for _, G := range i.Agency {
			i.wg.Add(1)
			Cek(G, &i.wg)
		}
		i.Reverse = true
	}
	i.wg.Wait()
}
