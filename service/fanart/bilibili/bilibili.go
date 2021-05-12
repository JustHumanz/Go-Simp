package main

import (
	"context"
	"encoding/json"
	"net/url"

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
	ModuleState = "BiliBili Fanart"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
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
	c.AddFunc(config.BiliBiliFanart, func() {
		for _, Group := range *GroupPayload {
			for _, Member := range Group.Members {
				if Member.BiliBiliHashtags != "" {
					log.WithFields(log.Fields{
						"Group":  Group.GroupName,
						"Vtuber": Member.EnName,
					}).Info("Start crawler bilibili")
					body, errcurl := network.CoolerCurl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(Member.BiliBiliHashtags), nil)
					if errcurl != nil {
						log.Error(errcurl)
					}
					var (
						TB engine.TBiliBili
					)
					_ = json.Unmarshal(body, &TB)
					if len(TB.Data.Cards) > 0 {
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
								}
								if New {
									Color, err := engine.GetColor(config.TmpDir, TBiliData.Photos[0])
									if err != nil {
										log.Error(err)
									}

									if config.GoSimpConf.Metric {
										gRCPconn.MetricReport(context.Background(), &pilot.Metric{
											MetricData: TBiliData.MarshallBin(),
											State:      config.FanartState,
										})
									}

									engine.SendFanArtNude(TBiliData, Bot, Color)
								}
							}
						}
					}
				}
			}
		}
	})
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}
