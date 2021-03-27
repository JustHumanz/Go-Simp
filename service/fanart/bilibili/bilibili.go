package bilibili

import (
	"encoding/json"
	"net/url"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/fanart/notif"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	configfile  config.ConfigFile
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	b.AddFunc(config.BiliBiliFanart, func() {
		for _, Group := range VtubersData.VtuberData {
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
						TB TBiliBili
					)
					_ = json.Unmarshal(body, &TB)
					if len(TB.Data.Cards) > 0 {
						for _, v := range TB.Data.Cards {
							var (
								STB SubTbili
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

									notif.SendNude(TBiliData, Bot, Color)
								}
							}
						}
					}
				}
			}
		}
	})
	log.Info("Enable Bilibili fanart module")
}
