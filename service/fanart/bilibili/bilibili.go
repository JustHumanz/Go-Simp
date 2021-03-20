package bilibili

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
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
								TBiliData := database.DataFanart{
									PermanentURL: "https://t.bilibili.com/" + v.Desc.DynamicIDStr + "?tab=2",
									Author:       v.Desc.UserProfile.Info.Uname,
									Likes:        v.Desc.Like,
									Photos:       img,
									Dynamic_id:   v.Desc.DynamicIDStr,
									Text:         STB.Item.Description,
									Member:       Member,
								}

								New, err := TBiliData.CheckTBiliBiliFanArt()
								if err != nil {
									log.Error(err)
								}
								if New {
									GroupData := Group.RemoveNillIconURL()
									ChannelData, err := database.ChannelTag(TBiliData.Member.ID, 1, config.Default, TBiliData.Member.Region)
									if err != nil {
										log.Error(err)
									}

									Color, err := engine.GetColor(config.TmpDir, TBiliData.Photos[0])
									if err != nil {
										log.Error(err)
									}
									tags := ""
									for i, Channel := range ChannelData {
										Channel.SetMember(TBiliData.Member)
										ctx := context.Background()
										UserTagsList, err := Channel.GetUserList(ctx)
										if err != nil {
											log.Error(err)
											break
										}
										if UserTagsList != nil {
											tags = strings.Join(UserTagsList, " ")
										} else {
											tags = "_"
										}

										if tags == "_" && Group.GroupName == config.Indie && !Channel.IndieNotif {
											//do nothing,like my life
										} else {
											tmp, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
												SetAuthor(strings.Title(GroupData.GroupName), GroupData.IconURL).
												SetTitle(TBiliData.Author).
												SetURL(TBiliData.PermanentURL).
												SetThumbnail(v.Desc.UserProfile.Info.Face).
												SetDescription(TBiliData.Text).
												SetImage(TBiliData.Photos[0]).
												AddField("User Tags", tags).
												//AddField("Similar art", msg).
												SetFooter("1/"+strconv.Itoa(len(TBiliData.Photos))+" photos", config.BiliBiliIMG).
												InlineAllFields().
												SetColor(Color).MessageEmbed)
											if err != nil {
												log.Error(tmp, err.Error())
												err = Channel.DelChannel(err.Error())
												if err != nil {
													log.Error(err)
												}
											}
											err = engine.Reacting(map[string]string{
												"ChannelID": Channel.ChannelID,
											}, Bot)
											if err != nil {
												log.Error(err)
											}
										}
										if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
											log.WithFields(log.Fields{
												"Func": "BiliBili Fanart",
											}).Warn(config.FanartSleep)
											time.Sleep(config.FanartSleep)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	})
	log.Info("Enable bilibili fanart module")
}
