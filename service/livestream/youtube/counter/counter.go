package main

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	"github.com/robfig/cron/v3"

	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	GroupPayload *[]database.Group
	gRCPconn     pilot.PilotServiceClient
)

const (
	ModuleState = config.YoutubeCounterModule
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start main youtube module
func main() {
	var (
		configfile config.ConfigFile
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
	Bot = engine.StartBot()

	database.Start(configfile)

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, GetPayload)
	c.AddFunc("0 */2 * * *", func() {
		engine.ExTknList = nil
	})
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkYtJob struct {
	wg         sync.WaitGroup
	mutex      sync.Mutex
	CekCounter map[string]bool
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	YoutubeCounter := &checkYtJob{
		CekCounter: make(map[string]bool),
	}
	for {
		log.WithFields(log.Fields{
			"Service": ModuleState,
			"Running": true,
		}).Info("request for running job")

		res, err := client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
			Service: ModuleState,
			Message: "Request",
			Alive:   true,
		})

		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": false,
			}).Info(res.Message)

			YoutubeCounter.Run()
			_, _ = client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
				Service: ModuleState,
				Message: "Done",
				Alive:   false,
			})
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": false,
			}).Info("reporting job was done")
		} else {
			log.WithFields(log.Fields{
				"Service": ModuleState,
			}).Info(res.Message)
		}

		time.Sleep(1 * time.Minute)
	}
}

func (i *checkYtJob) Run() {
	upcomingLiveStream, err := database.GetUpcomingFromCache()
	if err != nil {
		log.Error(err)
	}

	for key, val := range upcomingLiveStream {
		log.WithFields(log.Fields{
			"VideoID": key,
		}).Info("Count video")

		i.wg.Add(1)
		go func(Youtube database.LiveStream, Key string, w *sync.WaitGroup) {
			defer w.Done()
			for _, Group := range *GroupPayload {
				for _, Member := range Group.Members {
					if Member.ID == Youtube.Member.ID {
						Youtube.AddMember(Member).AddGroup(Group)
						if time.Since(Youtube.Schedul) > time.Until(Youtube.Schedul) {
							if i.CekCounterCount(Youtube.VideoID) {
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message: ModuleState + " error,loop notif detect,force remove cache & update livestream status to live,VideoID=" + Youtube.VideoID,
									Alive:   true,
									Service: ModuleState,
								})

								err = Youtube.RemoveUpcomingCache(Key)
								if err != nil {
									log.Panic(err)
								}

								Youtube.UpdateYt(config.LiveStatus)
								return
							}
							log.WithFields(log.Fields{
								"Vtuber":  Member.EnName,
								"Group":   Group.GroupName,
								"VideoID": Youtube.VideoID,
							}).Info("Vtuber upcoming schedule deadline,force change to live")

							Data, err := engine.YtAPI([]string{Youtube.VideoID})
							if err != nil {
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message: err.Error(),
									Service: ModuleState,
								})
								log.Error(err)
								return
							}
							if len(Data.Items) > 0 {
								if Data.Items[0].Snippet.VideoStatus != "none" {
									err = Youtube.RemoveUpcomingCache(Key)
									if err != nil {
										log.Panic(err)
									}

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

									isMemberOnly, err := regexp.MatchString("(memberonly|member|メン限)", strings.ToLower(Youtube.Title))
									if err != nil {
										log.Error(err)
									}

									engine.SendLiveNotif(&Youtube, Bot)
									if isMemberOnly {
										Youtube.UpdateYt(config.PrivateStatus)
									}
								} else if Data.Items[0].Snippet.VideoStatus == "none" {
									err = Youtube.RemoveUpcomingCache(Key)
									if err != nil {
										log.Panic(err)
									}

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

									log.WithFields(log.Fields{
										"VideoID": Youtube.VideoID,
										"Status":  config.PastStatus,
									}).Info("Update video status from " + Data.Items[0].Snippet.VideoStatus + " to past")

									LiveDur := durafmt.Parse(engine.ParseDuration(Data.Items[0].ContentDetails.Duration)).String()

									Youtube.UpdateLength(LiveDur).UpdateYt(config.PastStatus)
									engine.RemoveEmbed(Youtube.VideoID, Bot)

									if config.GoSimpConf.Metric {
										bit, err := Youtube.MarshalBinary()
										if err != nil {
											log.Error(err)
										}
										gRCPconn.MetricReport(context.Background(), &pilot.Metric{
											MetricData: bit,
											State:      config.PastStatus,
										})
									}
								}

							} else if Data.Items == nil || len(Data.Items) == 0 {
								log.WithFields(log.Fields{
									"VideoID": Youtube.VideoID,
									"Status":  config.PastStatus,
								}).Info("Update video status from " + config.UpcomingStatus + " to private")
								Youtube.UpdateYt(config.PrivateStatus)
								err = Youtube.RemoveUpcomingCache(Key)
								if err != nil {
									log.Panic(err)
								}
							}

							i.UpCek(Youtube.VideoID)
						}
						Youtube.
							SetState(config.YoutubeLive).
							UpdateStatus("reminder")
						engine.SendLiveNotif(&Youtube, Bot)
					}
				}
			}
		}(val.(database.LiveStream), key, &i.wg)
	}

	i.wg.Wait()

}

func (i *checkYtJob) UpCek(VideoID string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.CekCounter[VideoID] = true
}

func (i *checkYtJob) CekCounterCount(VideoID string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.CekCounter[VideoID] {
		return true
	}
	return false
}
