package main

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	GroupPayload *[]database.Group
	gRCPconn     pilot.PilotServiceClient
)

const (
	ModuleState = "Youtube_Counter"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start main youtube module
func main() {
	var (
		configfile config.ConfigFile
		err        error
	)

	GetPayload := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: ModuleState,
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error(), ModuleState)
			}
			log.Error("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &GroupPayload)
		if err != nil {
			log.Error(err)
		}
	}

	GetPayload()
	configfile.InitConf()

	Bot, err = discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	database.Start(configfile)

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, GetPayload)
	c.AddFunc("0 */13 * * *", func() {
		engine.ExTknList = nil
	})
	log.Info("Enable " + ModuleState)
	go CheckYtByTime()
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}

func CheckYtByTime() {
	var isHaveUpcoming = make(map[string]bool)
	var Counter = 1

	for {
		Counter++
		if Counter == 15 {
			isHaveUpcoming = make(map[string]bool)
			Counter = 1
			log.Info("Reset blacklist")
		}

		for _, Group := range *GroupPayload {
			for _, Member := range Group.Members {
				if Member.YoutubeID == "" || !Member.Active() || !isHaveUpcoming[Member.Name] {
					continue
				}

				log.WithFields(log.Fields{
					"Vtuber": Member.EnName,
					"Group":  Group.GroupName,
				}).Info("Checking Upcoming schedule")
				YoutubeStatus, Key, err := database.YtGetStatus(map[string]interface{}{
					"MemberID":   Member.ID,
					"MemberName": Member.Name,
					"Status":     config.UpcomingStatus,
					"State":      config.Sys,
				})
				if err != nil {
					if err.Error() != "not found any schdule" {
						log.WithFields(log.Fields{
							"Vtuber": Member.EnName,
							"Group":  Group.GroupName,
						}).Error(err)

						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					} else if err.Error() == "not found any schdule" {
						isHaveUpcoming[Member.Name] = true
						log.WithFields(log.Fields{
							"Vtuber": Member.EnName,
							"Group":  Group.GroupName,
						}).Info("Don't have upcoming stream,append to blacklist")
					}
					return
				}

				for _, Youtube := range YoutubeStatus {
					Youtube.AddMember(Member).AddGroup(Group)
					if time.Since(Youtube.Schedul) > time.Until(Youtube.Schedul) {
						log.WithFields(log.Fields{
							"Vtuber":  Member.EnName,
							"Group":   Group.GroupName,
							"VideoID": Youtube.VideoID,
						}).Info("Vtuber upcoming schedule deadline,force change to live")

						Data, err := engine.YtAPI([]string{Youtube.VideoID})
						if err != nil {
							log.Error(err)
							continue
						}
						if len(Data.Items) > 0 {
							if Data.Items[0].Snippet.VideoStatus != "none" {
								err = database.RemoveYtCache(Key, context.Background())
								if err != nil {
									log.Panic(err)
								}

								time.Sleep(10 * time.Second)

								if Data.Items[0].Statistics.ViewCount != "" {
									Youtube.UpdateViewers(Data.Items[0].Statistics.ViewCount)
								} else if Data.Items[0].Statistics.ViewCount == "0" && Youtube.Viewers == "0" || Youtube.Viewers == "" {
									Viewers, err := engine.GetWaiting(Youtube.VideoID)
									if err != nil {
										log.Error(err)
									}
									Youtube.UpdateViewers(Viewers)
								}

								if !Data.Items[0].LiveDetails.ActualStartTime.IsZero() {
									Youtube.UpdateSchdule(Data.Items[0].LiveDetails.ActualStartTime)
								}

								Youtube.UpdateStatus(config.LiveStatus).
									SetState(config.YoutubeLive).
									UpdateYt(config.LiveStatus)

								if Member.BiliRoomID != 0 {
									LiveBili, err := engine.GetRoomStatus(Member.BiliRoomID)
									if err != nil {
										log.Error(err)
									}
									if LiveBili.CheckScheduleLive() {
										Youtube.SetBiliLive(true).UpdateBiliToLive()
									}
								}

								if config.GoSimpConf.Metric {
									bit, err := Youtube.MarshalBinary()
									if err != nil {
										log.Error(err)
									}
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: bit,
										State:      config.LiveStatus,
									})
								}

								isMemberOnly, err := regexp.MatchString("memberonly", strings.ToLower(Youtube.Title))
								if err != nil {
									log.Error(err)
								}

								engine.SendLiveNotif(&Youtube, Bot)
								if isMemberOnly {
									Youtube.UpdateYt(config.PrivateStatus)
								}
							}
						}

						//one vtuber only have one livestream right
						break
					}
					Youtube.
						SetState(config.YoutubeLive).
						UpdateStatus("reminder")
					engine.SendLiveNotif(&Youtube, Bot)
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
