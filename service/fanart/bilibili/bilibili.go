package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/url"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
)

//Public variable
var (
	Bot          *discordgo.Session
	gRCPconn     pilot.PilotServiceClient
	torTransport = flag.Bool("Tor", false, "Enable multiTor for bot transport")
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.TBiliBiliService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	flag.Parse()
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start start tbilibili module
func main() {
	var (
		configfile config.ConfigFile
	)

	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		ServiceUUID: ServiceUUID,
		Service:     ServiceName,
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

	err = Bot.Open()
	if err != nil {
		log.Error(err)
	}

	database.Start(configfile)

	c := cron.New()
	c.Start()

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkBlJob struct {
	wg          sync.WaitGroup
	mutex       sync.Mutex
	Agency      []database.Group
	Reverse     bool
	FanArtIDTMP map[string]string
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	Bili := &checkBlJob{}

	for {
		log.WithFields(log.Fields{
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
				"Running":        res.Run,
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
		}
		time.Sleep(1 * time.Minute)
	}
}

func (i *checkBlJob) AddPostID(Member, VideoID string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.FanArtIDTMP[Member] = VideoID
}

func (i *checkBlJob) CekFirstPostID(Member, PostID string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.FanArtIDTMP[Member] == PostID {
		log.WithFields(log.Fields{
			"Vtuber": Member,
		}).Warn("Post Still same")
		return true
	}
	return false
}

func (k *checkBlJob) Run() {
	Cek := func(wg *sync.WaitGroup, Group database.Group) {
		defer wg.Done()

		for _, Member := range Group.Members {
			if Member.BiliBiliHashtag != "" {
				log.WithFields(log.Fields{
					"Agency": Group.GroupName,
					"Vtuber": Member.Name,
				}).Info("Start crawler bilibili")

				body := func() []byte {
					if config.GoSimpConf.MultiTOR != "" {
						body, errcurl := network.CoolerCurl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(Member.BiliBiliHashtag), nil)
						if errcurl != nil {
							log.Error(errcurl)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     errcurl.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body

					} else {
						body, errcurl := network.Curl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(Member.BiliBiliHashtag), nil)
						if errcurl != nil {
							log.Error(errcurl)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     errcurl.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body
					}

				}()
				var (
					TB engine.TBiliBili
				)
				_ = json.Unmarshal(body, &TB)
				if len(TB.Data.Cards) > 0 {
					if k.CekFirstPostID(Member.Name, TB.Data.Cards[0].Desc.DynamicIDStr) {
						continue
					}

					for _, v := range TB.Data.Cards {
						var (
							STB engine.SubTbili
							img []string
						)
						err := json.Unmarshal([]byte(v.Card), &STB)
						if err != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(err)
						}
						if STB.Item.Pictures != nil && v.Desc.Type == 2 { //type 2 is picture post (prob,heheheh)
							for _, pic := range STB.Item.Pictures {
								img = append(img, pic.ImgSrc)
							}

							TBiliData := database.DataFanart{
								PermanentURL: "https://t.bilibili.com/" + v.Desc.DynamicIDStr + "?tab=2",
								Author:       v.Desc.UserProfile.Info.Uname,
								AuthorAvatar: v.Desc.UserProfile.Info.Face,
								Likes:        v.Desc.Like,
								Photos:       img,
								Dynamic_id:   v.Desc.DynamicIDStr,
								Text:         STB.Item.Description,
								Member:       Member,
								Group:        Group,
								State:        config.BiliBiliArt,
							}

							New, err := TBiliData.CheckTBiliBiliFanArt()
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
							if New {

								if config.GoSimpConf.Metric {
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: TBiliData.MarshallBin(),
										State:      config.FanartState,
									})
								}

								engine.SendFanArtNude(TBiliData, Bot)
							}
						}
					}
					k.CekFirstPostID(Member.Name, TB.Data.Cards[0].Desc.DynamicIDStr)
				}
			}
		}
	}

	if k.Reverse {
		for j := len(k.Agency) - 1; j >= 0; j-- {
			Group := k.Agency
			k.wg.Add(1)
			go Cek(&k.wg, Group[j])
		}
		k.Reverse = false

	} else {
		for _, G := range k.Agency {
			k.wg.Add(1)
			go Cek(&k.wg, G)
		}
		k.Reverse = true
	}
	k.wg.Wait()

}
