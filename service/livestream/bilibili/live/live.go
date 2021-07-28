package main

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	loc          *time.Location
	Bot          *discordgo.Session
	GroupPayload *[]database.Group
	gRCPconn     pilot.PilotServiceClient
)

const (
	ModuleState = "LiveBiliBili"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start start twitter module
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
	c.AddJob("@every 07m", cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(&checkBlLiveeJob{}))

	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}

type checkBlLiveeJob struct {
	wg      sync.WaitGroup
	Reverse bool
}

func (i *checkBlLiveeJob) Run() {
	Cek := func(Group database.Group, wg *sync.WaitGroup) {
		defer wg.Done()

		for _, v := range []string{config.PastStatus, config.LiveStatus} {
			LiveBili, Key, err := database.BilGet(map[string]interface{}{
				"GroupID": Group.ID,
				"Status":  v,
			})
			if err != nil {
				log.Error(err)
			}

			if len(LiveBili) > 0 {
				for _, Bili := range LiveBili {
					for _, Member := range Group.Members {
						if Bili.Member.ID == Member.ID {
							Bili.AddGroup(Group).AddMember(Member)
							log.WithFields(log.Fields{
								"Group":  Group.GroupName,
								"Vtuber": Member.Name,
							}).Info("Checking LiveBiliBili")
							Status, err := engine.GetRoomStatus(Member.BiliRoomID)
							if err != nil {
								log.Error(err)
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message: err.Error(),
									Service: ModuleState,
								})
								continue
							}
							var (
								ScheduledStart time.Time
							)
							if Status.CheckScheduleLive() && Bili.Status != config.LiveStatus {
								//Live
								if Status.Data.RoomInfo.LiveStartTime != 0 {
									ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0).In(loc)
								} else {
									ScheduledStart = time.Now().In(loc)
								}

								Group.RemoveNillIconURL()

								log.WithFields(log.Fields{
									"Group":  Group.GroupName,
									"Vtuber": Member.EnName,
									"Start":  ScheduledStart,
								}).Info("Start live right now")

								Bili.UpdateStatus(config.LiveStatus).
									UpdateSchdule(ScheduledStart).
									UpdateViewers(strconv.Itoa(Status.Data.RoomInfo.Online)).
									UpdateThumbnail(Status.Data.RoomInfo.Cover).
									UpdateTitle(Status.Data.RoomInfo.Title)

								err := Bili.UpdateLiveBili()
								if err != nil {
									log.Error(err)
								}

								err = Bili.RemoveCache(Key, true)
								if err != nil {
									log.Panic(err)
								}

								if config.GoSimpConf.Metric {
									bit, err := Bili.MarshalBinary()
									if err != nil {
										log.Error(err)
									}
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: bit,
										State:      config.LiveStatus,
									})
								}

								engine.SendLiveNotif(&Bili, Bot)

							} else if !Status.CheckScheduleLive() && Bili.Status == config.LiveStatus {
								log.WithFields(log.Fields{
									"Group":  Group.GroupName,
									"Vtuber": Member.EnName,
									"Start":  Bili.Schedul,
								}).Info("Past live stream")
								engine.RemoveEmbed(strconv.Itoa(Bili.Member.BiliRoomID), Bot)
								Bili.UpdateEnd(time.Now().In(loc)).
									UpdateStatus(config.PastStatus)

								err = Bili.UpdateLiveBili()
								if err != nil {
									log.Error(err)
								}

								err = Bili.RemoveCache(Key, true)
								if err != nil {
									log.Panic(err)
								}

								if config.GoSimpConf.Metric {
									bit, err := Bili.MarshalBinary()
									if err != nil {
										log.Error(err)
									}
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: bit,
										State:      config.PastStatus,
									})
								}

							} else {
								Bili.UpdateViewers(strconv.Itoa(Status.Data.RoomInfo.Online))
								err := Bili.UpdateLiveBili()
								if err != nil {
									log.Error(err)
								}
							}
						}
					}
				}
			}
		}
	}

	if i.Reverse {
		for j := len(*GroupPayload) - 1; j >= 0; j-- {
			i.wg.Add(1)
			Grp := *GroupPayload
			go Cek(Grp[j], &i.wg)
		}
		i.Reverse = false

	} else {
		for _, G := range *GroupPayload {
			i.wg.Add(1)
			go Cek(G, &i.wg)
		}
		i.Reverse = true
	}
	i.wg.Wait()
}
