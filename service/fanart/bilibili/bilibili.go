package main

import (
	"context"
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
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
	GroupPayload *[]database.Group
	gRCPconn     pilot.PilotServiceClient
)

const (
	ModuleState = config.TBiliBiliModule
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start start tbilibili module
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
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkBlJob struct {
	wg          sync.WaitGroup
	mutex       sync.Mutex
	Reverse     bool
	FanArtIDTMP map[string]string
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	Bili := &checkBlJob{}

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

			Bili.Run()
			_, _ = client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
				Service: ModuleState,
				Message: "Done",
				Alive:   false,
			})
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": false,
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
			if Member.BiliBiliHashtags != "" {
				log.WithFields(log.Fields{
					"Group":  Group.GroupName,
					"Vtuber": Member.EnName,
				}).Info("Start crawler bilibili")

				body, errcurl := network.CoolerCurl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(Member.BiliBiliHashtags), nil)
				if errcurl != nil {
					log.Error(errcurl)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: errcurl.Error(),
						Service: ModuleState,
					})
				}
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
							log.Error(err)
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
								log.Error(err)
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message: err.Error(),
									Service: ModuleState,
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
		for j := len(*GroupPayload) - 1; j >= 0; j-- {
			Group := *GroupPayload
			k.wg.Add(1)
			go Cek(&k.wg, Group[j])
		}
		k.Reverse = false

	} else {
		for _, G := range *GroupPayload {
			k.wg.Add(1)
			go Cek(&k.wg, G)
		}
		k.Reverse = true
	}
	k.wg.Wait()

}
