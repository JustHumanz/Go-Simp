package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/hako/durafmt"

	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	gRCPconn     pilot.PilotServiceClient
	torTransport = flag.Bool("Tor", false, "Enable multiTor for bot transport")
	agency       = flag.Bool("Agency", false, "Enable scraping for vtuber agency")
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.YoutubeLiveTrackerService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	flag.Parse()
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start main youtube module
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
	Bot = engine.StartBot(*torTransport)

	database.Start(configfile)

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkYtCekJob struct {
	agency  []database.Group
	Reverse bool
	Counter int
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	YoutubeChecker := &checkYtCekJob{
		Counter: 1,
	}

	hostname := engine.GetHostname()

	for {
		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "PastTracker",
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

			YoutubeChecker.agency = engine.UnMarshalPayload(res.VtuberPayload)
			if len(YoutubeChecker.agency) == 0 {
				msg := "vtuber agency was nill,force close the unit"
				pilot.ReportDeadService(msg, ServiceName)
				log.Fatalln(msg)
			}
			YoutubeChecker.Run()

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

		YoutubeChecker.Counter++
		time.Sleep(1 * time.Minute)
	}
}

func (i *checkYtCekJob) Run() {
	if i.Reverse {
		for j := len(i.agency) - 1; j >= 0; j-- {
			Grp := i.agency
			StartCheckYT(Grp[j])
		}

	} else {
		for _, G := range i.agency {
			StartCheckYT(G)
		}
	}
}

//StartCheckYT Youtube rss and API
func StartCheckYT(Group database.Group) {

	//check vtuber agency youtube channel
	if Group.YoutubeChannels != nil && *agency {
		for _, YtChan := range Group.YoutubeChannels {
			log.WithFields(log.Fields{
				"Agency":    Group.GroupName,
				"ChannelID": YtChan.YtChannel,
				"Region":    YtChan.Region,
			}).Info("Checking agency channel")

			Yt, err := Group.GetYtLiveStream(config.LiveStatus, "")
			if err != nil {
				log.Error(err)
			}

			for _, YoutubeCache := range Yt {
				if YoutubeCache.ID != 0 {
					log.WithFields(log.Fields{
						"Group":   Group.GroupName,
						"VideoID": YoutubeCache.VideoID,
					}).Info("Update VideoID")

					YoutubeData := &YoutubeCache

					Data, err := engine.YtAPI([]string{YoutubeCache.VideoID})
					if err != nil {
						log.WithFields(log.Fields{
							"Agency":    Group.GroupName,
							"ChannelID": YtChan.YtChannel,
							"Region":    YtChan.Region,
						}).Error(err)

						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message:     err.Error(),
							Service:     ServiceName,
							ServiceUUID: ServiceUUID,
						})
					}

					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
						continue
					}

					Items := Data.Items[0]

					YoutubeData.SetState(config.YoutubeLive).AddGroup(Group)

					if Items.Snippet.VideoStatus == "none" && YoutubeData.Status == config.LiveStatus {
						log.WithFields(log.Fields{
							"VideoData ID": YoutubeCache.VideoID,
							"Status":       config.PastStatus,
						}).Info("Update video status from " + YoutubeData.Status + " to past")
						YoutubeData.UpdateEnd(Items.LiveDetails.EndTime).
							UpdateViewers(Items.Statistics.ViewCount).
							UpdateLength(durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String()).
							UpdateGroupYt(config.PastStatus)

						engine.RemoveEmbed(YoutubeCache.VideoID, Bot)

						if config.GoSimpConf.Metric {
							bit, err := YoutubeData.MarshalBinary()
							if err != nil {
								log.WithFields(log.Fields{
									"Agency":    Group.GroupName,
									"ChannelID": YtChan.YtChannel,
									"Region":    YtChan.Region,
								}).Error(err)
							}

							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bit,
								State:      config.PastStatus,
							})
						}

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus && YoutubeData.Status == config.PastStatus {
						log.WithFields(log.Fields{
							"VideoData ID": YoutubeCache.VideoID,
							"Status":       Items.Snippet.VideoStatus,
						}).Info("maybe yt error or human error")

						YoutubeData.UpdateStatus(config.UpcomingStatus)
					} else if Items.Snippet.VideoStatus == "none" && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": YoutubeCache.VideoID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.PastStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateViewers(Items.Statistics.ViewCount).UpdateGroupYt(config.PastStatus)

					} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": YoutubeCache.VideoID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.LiveStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateViewers(Items.Statistics.ViewCount).UpdateGroupYt(config.LiveStatus)

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus {
						if Items.LiveDetails.StartTime != YoutubeData.Schedul {
							log.WithFields(log.Fields{
								"VideoData ID": YoutubeCache.VideoID,
								"Old schdule":  YoutubeData.Schedul,
								"New schdule":  Items.LiveDetails.StartTime,
								"Status":       config.UpcomingStatus,
							}).Info("Livestream schdule changed")

							YoutubeData.UpdateSchdule(Items.LiveDetails.StartTime)
							YoutubeData.UpdateGroupYt(config.UpcomingStatus)
						}
					}
				}
			}
		}
	}

	//check vtuber agency members youtube channel
	var (
		wg      sync.WaitGroup
		counter = 0
	)
	for _, v := range Group.Members {
		if !v.IsYtNill() && v.Active() {
			wg.Add(1)
			go func(Member database.Member, w *sync.WaitGroup) {
				defer w.Done()
				log.WithFields(log.Fields{
					"Vtuber": Member.Name,
					"Agency": Group.GroupName,
				}).Info("Checking Vtuber channel")

				Yt, err := Member.GetYtLiveStream(config.LiveStatus)
				if err != nil {
					log.Error(err)
				}
				for _, YoutubeCache := range Yt {
					if YoutubeCache.ID != 0 {
						log.WithFields(log.Fields{
							"Group":   Group.GroupName,
							"Member":  Member.Name,
							"VideoID": YoutubeCache.VideoID,
						}).Info("Update VideoID")

						Data, err := engine.YtAPI([]string{YoutubeCache.VideoID})
						if err != nil {
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     err.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
							log.WithFields(log.Fields{
								"Vtuber": Member.Name,
								"Agency": Group.GroupName,
							}).Error(err)
						}

						if len(Data.Items) == 0 {
							fmt.Println("Opps something error\n", Data)
							continue
						}

						YoutubeData := &YoutubeCache
						YoutubeData.GetYtVideoDetail()

						Items := Data.Items[0]

						YoutubeData.AddMember(Member).AddGroup(Group).SetState(config.YoutubeLive)

						if Items.Snippet.VideoStatus == "none" && YoutubeData.Status == config.LiveStatus {
							log.WithFields(log.Fields{
								"VideoData ID": YoutubeCache.VideoID,
								"Status":       config.PastStatus,
							}).Info("Update video status from " + YoutubeData.Status + " to past")
							YoutubeData.UpdateEnd(Items.LiveDetails.EndTime).
								UpdateViewers(Items.Statistics.ViewCount).
								UpdateLength(durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String()).UpdateYt(config.PastStatus)

							engine.RemoveEmbed(YoutubeCache.VideoID, Bot)

							if config.GoSimpConf.Metric {
								bit, err := YoutubeData.MarshalBinary()
								if err != nil {
									log.WithFields(log.Fields{
										"Vtuber": Member.Name,
										"Agency": Group.GroupName,
									}).Error(err)
								}
								gRCPconn.MetricReport(context.Background(), &pilot.Metric{
									MetricData: bit,
									State:      config.PastStatus,
								})
							}

						} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Status == config.UpcomingStatus {
							log.WithFields(log.Fields{
								"VideoData ID": YoutubeCache.VideoID,
								"Status":       config.LiveStatus,
							}).Info("Update video status from " + YoutubeData.Status + " to live")
							YoutubeData.UpdateStatus(config.LiveStatus)

							log.Info("Update database")
							if !Items.LiveDetails.ActualStartTime.IsZero() {
								YoutubeData.UpdateSchdule(Items.LiveDetails.ActualStartTime)
							}

							YoutubeData.UpdateViewers(Items.Statistics.ViewCount).UpdateYt(YoutubeData.Status)
							engine.SendLiveNotif(YoutubeData, Bot)

						} else if Items.Snippet.VideoStatus == "none" && YoutubeData.Viewers != Items.Statistics.ViewCount {
							log.WithFields(log.Fields{
								"VideoData ID": YoutubeCache.VideoID,
								"Viwers past":  YoutubeData.Viewers,
								"Viwers now":   Items.Statistics.ViewCount,
								"Status":       config.PastStatus,
							}).Info("Update Viwers")
							YoutubeData.UpdateViewers(Items.Statistics.ViewCount).UpdateYt(config.PastStatus)

						} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Viewers != Items.Statistics.ViewCount {
							log.WithFields(log.Fields{
								"VideoData ID": YoutubeCache.VideoID,
								"Viwers past":  YoutubeData.Viewers,
								"Viwers now":   Items.Statistics.ViewCount,
								"Status":       config.LiveStatus,
							}).Info("Update Viwers")
							YoutubeData.UpdateViewers(Items.Statistics.ViewCount).UpdateYt(config.LiveStatus)

						} else if Items.Snippet.VideoStatus == config.UpcomingStatus {
							if Items.LiveDetails.StartTime != YoutubeData.Schedul {
								log.WithFields(log.Fields{
									"VideoData ID": YoutubeCache.VideoID,
									"Old schdule":  YoutubeData.Schedul,
									"New schdule":  Items.LiveDetails.StartTime,
									"Status":       config.UpcomingStatus,
								}).Info("Livestream schdule changed")

								YoutubeData.UpdateSchdule(Items.LiveDetails.StartTime)
								YoutubeData.UpdateYt(config.UpcomingStatus)
							}
						}
					}
				}
			}(v, &wg)
			counter++
		}
		if counter%10 == 0 && counter != 0 {
			log.WithFields(log.Fields{
				"Wait wg": 10,
				"Counter": counter,
			}).Info("Waiting waitgroup")
			wg.Wait()
		}
	}
	wg.Wait()
}
