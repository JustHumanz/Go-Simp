package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	Bot         *discordgo.Session
	gRCPconn    pilot.PilotServiceClient
	ServiceUUID = uuid.New().String()
)

const (
	ServiceName = config.LiveBiliBiliService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start start twitter module
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

	database.Start(configfile)

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkBlLiveeJob struct {
	wg      sync.WaitGroup
	Reverse bool
	Agency  []database.Group
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	Bili := &checkBlLiveeJob{}
	hostname := engine.GetHostname()

	for {

		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Request",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Running":        true,
				"UUID":           ServiceUUID,
				"Agency Payload": res.VtuberMetadata,
			}).Info(res.Message)

			Bili.Agency = engine.UnMarshalPayload(res.VtuberPayload)
			if len(Bili.Agency) == 0 {
				msg := "vtuber agency was nill,force close the unit"
				pilot.ReportDeadService(msg, ServiceName)
				log.Fatalln(msg)
			}
			Bili.Run()

			_, _ = client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
				Service:     ServiceName,
				Message:     "Done",
				ServiceUUID: ServiceUUID,
			})
			log.WithFields(log.Fields{
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info("reporting job was done")

		} else {
			log.WithFields(log.Fields{
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info(res.Message)
		}

		time.Sleep(1 * time.Minute)
	}
}

func (i *checkBlLiveeJob) Run() {
	Cek := func(Agency database.Group, wg *sync.WaitGroup) {
		defer wg.Done()

		for _, v := range []string{config.PastStatus, config.LiveStatus} {
			LiveBili, err := Agency.GetBlLiveStream(v)
			if err != nil {
				log.WithFields(log.Fields{
					"Agency": Agency.GroupName,
				}).Error(err)
			}

			if len(LiveBili) > 0 {
				for _, Bili := range LiveBili {
					for _, Member := range Agency.Members {
						if Bili.Member.ID == Member.ID {
							log.WithFields(log.Fields{
								"Agency": Agency.GroupName,
								"Vtuber": Member.Name,
							}).Info("Checking LiveBiliBili")
							Status, err := engine.GetRoomStatus(Member.BiliBiliRoomID)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Agency.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message: err.Error(),
									Service: ServiceName,
								})
								continue
							}
							var (
								ScheduledStart time.Time
							)
							if Status.CheckScheduleLive() && Bili.Status != config.LiveStatus {
								//Live
								if Status.Data.RoomInfo.LiveStartTime != 0 {
									ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0)
								} else {
									ScheduledStart = time.Now()
								}

								Agency.RemoveNillIconURL()

								log.WithFields(log.Fields{
									"Group":  Agency.GroupName,
									"Vtuber": Member.EnName,
									"Start":  ScheduledStart,
								}).Info("Start live right now")

								Bili.UpdateStatus(config.LiveStatus).
									SetState(config.BiliLive).
									UpdateSchdule(ScheduledStart).
									UpdateViewers(strconv.Itoa(Status.Data.RoomInfo.Online)).
									UpdateThumbnail(Status.Data.RoomInfo.Cover).
									UpdateTitle(Status.Data.RoomInfo.Title).
									UpdateDesc(func() string {
										if Status.Data.RoomInfo.Description == "" {
											return "????"
										} else {
											return Status.Data.RoomInfo.Description
										}
									}())

									//Remove cache before update the status
								err := Bili.RemoveCache(fmt.Sprintf("%d-%s-%s-%d", Member.ID, Member.Name, config.PastStatus, Member.BiliBiliID))
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Agency.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}

								err = Bili.UpdateLiveBili()
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Agency.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}

								if config.GoSimpConf.Metric {
									bit, err := Bili.MarshalBinary()
									if err != nil {
										log.WithFields(log.Fields{
											"Agency": Agency.GroupName,
											"Vtuber": Member.Name,
										}).Error(err)
									}
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: bit,
										State:      config.LiveStatus,
									})
								}

								engine.SendLiveNotif(&Bili, Bot)

							} else if !Status.CheckScheduleLive() && Bili.Status == config.LiveStatus {
								log.WithFields(log.Fields{
									"Group":  Agency.GroupName,
									"Vtuber": Member.EnName,
									"Start":  Bili.Schedul,
								}).Info("Past live stream")
								engine.RemoveEmbed(strconv.Itoa(Bili.Member.BiliBiliRoomID), Bot)
								Bili.UpdateEnd(time.Now()).
									UpdateStatus(config.PastStatus)

									//Remove cache before update the status
								err := Bili.RemoveCache(fmt.Sprintf("%d-%s-%s-%d", Member.ID, Member.Name, config.LiveStatus, Member.BiliBiliID))
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Agency.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}

								err = Bili.UpdateLiveBili()
								if err != nil {
									log.WithFields(log.Fields{
										"Agency": Agency.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}

								if config.GoSimpConf.Metric {
									bit, err := Bili.MarshalBinary()
									if err != nil {
										log.WithFields(log.Fields{
											"Agency": Agency.GroupName,
											"Vtuber": Member.Name,
										}).Error(err)
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
									log.WithFields(log.Fields{
										"Agency": Agency.GroupName,
										"Vtuber": Member.Name,
									}).Error(err)
								}
							}
						}
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
			go Cek(G, &i.wg)
		}
		i.Reverse = true
	}
	i.wg.Wait()
}
