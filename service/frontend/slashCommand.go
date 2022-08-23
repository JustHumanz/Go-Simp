package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

var SendSimpleMsg = func(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		log.Error(err)
	}
}

var (
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mytags": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			list, err := database.UserStatus(i.Member.User.ID, i.ChannelID)
			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("mytags")

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.AvatarURL("128"))
			if err != nil {
				log.Error(err)
			}

			if list != nil {
				tableString := &strings.Builder{}
				table := tablewriter.NewWriter(tableString)
				table.SetAutoWrapText(false)
				table.SetAutoFormatHeaders(true)
				table.SetCenterSeparator("")
				table.SetColumnSeparator("")
				table.SetRowSeparator("")
				table.SetHeaderLine(true)
				table.SetBorder(false)
				table.SetTablePadding("\t")
				table.SetNoWhiteSpace(true)
				table.SetHeader([]string{"Vtuber Group", "Vtuber Name", "Reminder"})
				table.AppendBulk(list)
				table.Render()

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
							SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
							SetThumbnail(i.Member.User.AvatarURL("128")).
							SetDescription("```" + tableString.String() + "```").
							SetColor(Color).MessageEmbed},
					},
				})
				if err != nil {
					log.Error(err)
				}
			} else {

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
							SetDescription("Your tag list is empty.").
							SetTitle("404 Not found").
							SetImage(engine.NotFoundIMG()).
							SetColor(Color).MessageEmbed},
					},
				})
				if err != nil {
					log.Error(err)
				}
			}
		},
		"livestream": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var Agency database.Group
			var Member database.Member
			var embed []*discordgo.MessageEmbed

			state := i.ApplicationCommandData().Options[0].IntValue()
			Avatar := i.Member.User.AvatarURL("128")
			Nick := i.Member.User.Username
			SelectedAgency := ""
			SelectedVtuber := ""
			region := ""
			status := ""

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}

				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}

				if v.Name == "status" {
					status = func() string {
						if v.IntValue() == 1 {
							return config.LiveStatus
						} else if v.IntValue() == 2 {
							return config.UpcomingStatus
						} else {
							return config.PastStatus
						}
					}()
				}

				if v.Name == "region" {
					region = v.StringValue()
				}
			}

			SendMessage := func(e []*discordgo.MessageEmbed) {
				log.WithFields(log.Fields{
					"User":        i.Member.User.Username,
					"Channel":     i.ChannelID,
					"Status":      status,
					"Payload len": len(e),
					"Agency":      Agency.GroupName,
					"Vtuber":      Member.Name,
				}).Info("livestream")

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "On progress~",
					},
				})
				if err != nil {
					log.Error(err)
				}

				if len(e) == 0 {
					_, err := s.ChannelMessageSendEmbed(i.ChannelID, engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetDescription("It looks like doesn't have a livestream schedule for now").
						SetImage(engine.NotFoundIMG()).MessageEmbed)
					if err != nil {
						log.Error(err)
					}

				} else {
					for _, v := range e {
						_, err := s.ChannelMessageSendEmbed(i.ChannelID, v)
						if err != nil {
							log.Error(err)
						}
						time.Sleep(1 * time.Second)
					}
				}
			}

			NotSupportUp := func(aa, b string) {
				a := engine.NewEmbed().
					SetAuthor(Nick, Avatar).
					SetTitle("Still development").
					SetDescription(aa + " command in " + b + " is still development").
					SetImage(engine.NotFoundIMG()).SetFooter("honestly,developer doesn't know how to get upcoming schedule from " + b).MessageEmbed

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{a},
					},
				})
				if err != nil {
					log.Error(err)
				}
			}

			for _, v := range GroupsPayload {
				if strings.EqualFold(v.GroupName, SelectedAgency) {
					Agency = v
					if SelectedVtuber != "" {
						for _, v2 := range v.Members {
							if engine.CheckVtuberName(v2, SelectedVtuber) {
								Member = v2
							}
						}

						//Retrun invalid vtuber name
						if Member.IsMemberNill() {
							SendSimpleMsg(s, i, fmt.Sprintf("Invalid vtuber name %s", SelectedVtuber))
							log.WithFields(log.Fields{
								"ChannelID":     i.ChannelID,
								"SelectedGroup": SelectedVtuber,
								"State":         "Livestream",
							}).Warn("Invalid vtuber name")
							return
						}
					}
				}
			}

			if Agency.IsNull() {
				SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", SelectedAgency))
				log.WithFields(log.Fields{
					"ChannelID":     i.ChannelID,
					"SelectedGroup": SelectedAgency,
					"State":         "Livestream",
				}).Warn("Invalid agency name")
				return
			}

			if state == 1 {
				if !Member.IsMemberNill() {
					YoutubeData, err := Member.GetYtLiveStream(status)
					if err != nil {
						log.Error(err)
					}

					FixName := engine.FixName(Member.EnName, Member.JpName)
					if YoutubeData != nil {
						for _, Youtube := range YoutubeData {
							FanBase := "simps"

							loc, err := engine.Zawarudo(Member.Region)
							if err != nil {
								log.Error(err)
							}

							duration := durafmt.Parse(Youtube.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)

							Color, err := engine.GetColor(config.TmpDir, Youtube.Thumb)
							if err != nil {
								if err.Error() == "Server Error,status get 404 Not Found" {
									Youtube.UpdateYt("private")
								} else {
									log.Error(err)
								}
							}

							if Member.Fanbase != "" {
								FanBase = Member.Fanbase
							}

							view, err := strconv.Atoi(Youtube.Viewers)
							if err != nil {
								log.Error(err)
							}
							Viewers := engine.NearestThousandFormat(float64(view))

							if status == config.PastStatus {
								durationlive := durafmt.Parse(Youtube.End.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Viewers).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							} else {
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Start live in", duration.String()).
									AddField("Viewers", Viewers+" "+FanBase).
									InlineAllFields().
									AddField("Type", engine.YtFindType(Youtube.Title)).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							}
						}

					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetTitle(FixName).
							SetDescription("It looks like `"+FixName+"` doesn't have a `"+status+"` schedule for now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}

					SendMessage(embed)

				} else if !Agency.IsNull() {
					YoutubeData, err := Agency.GetYtLiveStream(status, region)
					if err != nil {
						log.Error(err)
					}

					if len(YoutubeData) > 10 {
						SendSimpleMsg(s, i, "On progress~")
					}

					if YoutubeData != nil {
						for _, Youtube := range YoutubeData {
							Member, err := FindVtuber(Youtube.Member.ID)
							if err != nil {
								log.Error(err)
							}
							Youtube.AddMember(Member)
							FixName := engine.FixName(Youtube.Member.EnName, Youtube.Member.JpName)
							FanBase := "simps"

							loc, err := engine.Zawarudo(Youtube.Member.Region)
							if err != nil {
								log.Error(err)
							}

							Color, err := engine.GetColor(config.TmpDir, Youtube.Thumb)
							if err != nil {
								if err.Error() == "Server Error,status get 404 Not Found" {
									Youtube.UpdateYt("private")
								} else {
									log.Error(err)
								}
							}

							if Youtube.Member.Fanbase != "" {
								FanBase = Youtube.Member.Fanbase
							}

							view, err := strconv.Atoi(Youtube.Viewers)
							if err != nil {
								log.Error(err)
							}
							Viewers := engine.NearestThousandFormat(float64(view))
							if status == config.PastStatus {
								duration := durafmt.Parse(Youtube.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
								durationlive := durafmt.Parse(Youtube.End.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Viewers).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							} else if status == config.LiveStatus {
								expiresAt := time.Now().In(loc)
								duration := durafmt.Parse(expiresAt.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Start live in", duration.String()).
									AddField("Viewers", Viewers+" "+FanBase).
									InlineAllFields().
									AddField("Type", engine.YtFindType(Youtube.Title)).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							} else {
								duration := durafmt.Parse(Youtube.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Start live in", duration.String()).
									AddField("Viewers", Viewers+" "+FanBase).
									InlineAllFields().
									AddField("Type", engine.YtFindType(Youtube.Title)).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							}
						}
					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetTitle(Agency.GroupName).
							SetDescription("It looks like `"+Agency.GroupName+"` doesn't have a `"+status+"` schedule for now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}

					SendMessage(embed)
				}

			} else if state == 2 {
				if status == config.UpcomingStatus {
					NotSupportUp(status, "BiliBili")
					return
				}

				loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
				if !Member.IsMemberNill() {
					LiveData, err := Member.GetBlLiveStream(status)
					if err != nil {
						log.Error(err)
					}

					if LiveData.ID != 0 {
						Color, err := engine.GetColor(config.TmpDir, LiveData.Thumb)
						if err != nil {
							log.Error(err)
						}

						Member, err := FindVtuber(LiveData.Member.ID)
						if err != nil {
							log.Error(err)
						}

						LiveData.AddMember(Member)
						FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
						view, err := strconv.Atoi(LiveData.Viewers)
						if err != nil {
							log.Error(err)
						}
						if status == config.PastStatus {
							diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
							embed = append(embed, engine.NewEmbed().
								SetTitle(FixName).
								SetAuthor(Nick, Avatar).
								SetDescription(LiveData.Desc).
								SetThumbnail(Member.BiliBiliAvatar).
								SetImage(LiveData.Thumb).
								SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliBiliRoomID)).
								AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Online", engine.NearestThousandFormat(float64(view))).
								SetColor(Color).
								SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
						} else {
							diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
							embed = append(embed, engine.NewEmbed().
								SetTitle(FixName).
								SetAuthor(Nick, Avatar).
								SetDescription(LiveData.Desc).
								SetThumbnail(LiveData.Member.BiliBiliAvatar).
								SetImage(LiveData.Thumb).
								SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliBiliRoomID)).
								AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Online", engine.NearestThousandFormat(float64(view))).
								SetColor(Color).
								SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
						}

					}
				} else {
					embed = append(embed, engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetDescription("It looks like `"+Agency.GroupName+"` doesn't have a `"+status+"` schedule right now").
						SetImage(engine.NotFoundIMG()).MessageEmbed)
				}

				SendMessage(embed)

			} else {
				LiveBili, err := Agency.GetBlLiveStream(status)
				if err != nil {
					log.Error(err)
				}

				if len(LiveBili) > 10 {
					SendSimpleMsg(s, i, "On progress~")
				}

				loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/

				if LiveBili != nil {
					for _, LiveData := range LiveBili {
						Member, err := FindVtuber(LiveData.Member.ID)
						if err != nil {
							log.Error(err)
						}

						Color, err := engine.GetColor(config.TmpDir, LiveData.Thumb)
						if err != nil {
							log.Error(err)
						}

						LiveData.AddMember(Member)
						FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
						view, err := strconv.Atoi(LiveData.Viewers)
						if err != nil {
							log.Error(err)
						}

						if status == config.PastStatus {
							diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
							embed = append(embed, engine.NewEmbed().
								SetTitle(FixName).
								SetAuthor(Nick, Avatar).
								SetDescription(LiveData.Desc).
								SetThumbnail(Member.BiliBiliAvatar).
								SetImage(LiveData.Thumb).
								SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliBiliRoomID)).
								AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Online", engine.NearestThousandFormat(float64(view))).
								SetColor(Color).
								SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
						} else {
							diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
							embed = append(embed, engine.NewEmbed().
								SetTitle(FixName).
								SetAuthor(Nick, Avatar).
								SetDescription(LiveData.Desc).
								SetThumbnail(LiveData.Member.BiliBiliAvatar).
								SetImage(LiveData.Thumb).
								SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliBiliRoomID)).
								AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Online", engine.NearestThousandFormat(float64(view))).
								SetColor(Color).
								SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
						}
					}
				} else {
					embed = append(embed, engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetDescription("It looks like `"+Agency.GroupName+"` doesn't have a `"+status+"` schedule right now").
						SetImage(engine.NotFoundIMG()).MessageEmbed)
				}
				SendMessage(embed)
			}

		},
		"art": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var embed *discordgo.MessageEmbed
			var DynamicData DynamicSvr
			var SelectedAgency string
			var SelectedVtuber string

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}

				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}
			}

			Avatar := i.Member.User.AvatarURL("128")
			Nick := i.Member.User.Username

			SendNude := func(Data *database.DataFanart) {

				log.WithFields(log.Fields{
					"User":    i.Member.User.Username,
					"Channel": i.ChannelID,
					"Vtuber":  SelectedVtuber,
					"Agency":  SelectedAgency,
				}).Info("art")

				Color, err := engine.GetColor(config.TmpDir, Avatar)
				if err != nil {
					log.Error(err)
				}

				if Data.State == config.BiliBiliArt {
					if config.GoSimpConf.MultiTOR != "" {
						body, errcurl := network.CoolerCurl("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id="+Data.Dynamic_id, nil)
						if errcurl != nil {
							log.Error(errcurl)
						}

						err := json.Unmarshal(body, &DynamicData)
						if err != nil {
							log.Error(err)
						}
					} else {
						body, errcurl := network.Curl("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id="+Data.Dynamic_id, nil)
						if errcurl != nil {
							log.Error(errcurl)
						}

						err := json.Unmarshal(body, &DynamicData)
						if err != nil {
							log.Error(err)
						}
					}

					embed = engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetTitle(Data.Author).
						SetThumbnail(DynamicData.GetUserAvatar()).
						SetDescription(RemovePic(Data.Text)).
						SetURL(Data.PermanentURL).
						SetImage(Data.Photos...).
						SetColor(Color).
						InlineAllFields().
						SetFooter(Data.State, config.BiliBiliIMG).MessageEmbed
				} else if Data.State == config.TwitterArt {
					embed = engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetTitle(Data.Author).
						SetThumbnail(engine.GetAuthorAvatar(Data.Author)).
						SetDescription(RemovePic(Data.Text)).
						SetURL(Data.PermanentURL).
						SetImage(Data.Photos...).
						SetColor(Color).
						InlineAllFields().
						SetFooter(Data.State, config.TwitterIMG).MessageEmbed
				} else {
					embed = engine.NewEmbed().
						SetAuthor(Nick, Avatar).
						SetTitle(Data.Author).
						SetDescription(Data.Text).
						SetURL(Data.PermanentURL).
						SetImage(config.PixivProxy+Data.Photos[0]).
						SetColor(Color).
						InlineAllFields().
						SetFooter(Data.State, config.PixivIMG).MessageEmbed

				}
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							embed,
						},
					},
				})
				if err != nil {
					log.Error(err)
				}

			}

			var Agency database.Group
			var VtuberName database.Member
			for _, v := range GroupsPayload {
				if strings.EqualFold(SelectedAgency, v.GroupName) {
					Agency = v
					if SelectedVtuber != "" {
						for _, k := range v.Members {
							if engine.CheckVtuberName(k, SelectedVtuber) {
								VtuberName = k
							}
						}

						if VtuberName.IsMemberNill() {
							SendSimpleMsg(s, i, fmt.Sprintf("Invalid vtuber name %s ", SelectedVtuber))
							log.WithFields(log.Fields{
								"ChannelID":     i.ChannelID,
								"SelectedGroup": SelectedAgency,
								"State":         "Art",
							}).Warn("Invalid vtuber name")
							return
						}
					}
				}
			}

			if Agency.IsNull() {
				SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", SelectedAgency))
				log.WithFields(log.Fields{
					"ChannelID":     i.ChannelID,
					"SelectedGroup": SelectedAgency,
					"State":         "Art",
				}).Warn("Invalid agency name")
				return
			}

			if !VtuberName.IsMemberNill() {
				FanArt, err := VtuberName.GetRandomFanart()

				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Oops,something error\n" + err.Error(),
						},
					})
				} else {
					SendNude(FanArt)
				}

			} else if !Agency.IsNull() {
				FanArt, err := Agency.GetRandomFanart()
				if err != nil {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Oops,something error\n" + err.Error(),
						},
					})
				} else {
					SendNude(FanArt)
				}
			}
		},
		"lewd": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			ChannelRaw, err := s.Channel(i.ChannelID)
			if err != nil {
				log.Error(err)
			}
			if ChannelRaw.NSFW {
				var embed *discordgo.MessageEmbed
				Avatar := i.Member.User.AvatarURL("128")
				Nick := i.Member.User.Username
				var SelectedAgency string
				var SelectedVtuber string

				for _, v := range i.ApplicationCommandData().Options {
					if v.Name == engine.AgencyOption {
						SelectedAgency = v.StringValue()
					}

					if v.Name == engine.VtuberOption {
						SelectedVtuber = v.StringValue()
					}
				}

				SendNude := func(Data *database.DataFanart) {
					log.WithFields(log.Fields{
						"User":    i.Member.User.Username,
						"Channel": i.ChannelID,
						"Vtuber":  SelectedVtuber,
						"Agency":  SelectedAgency,
					}).Info("lewd")

					Color, err := engine.GetColor(config.TmpDir, Avatar)
					if err != nil {
						log.Error(err)
					}

					if Data.State == config.TwitterArt {
						embed = engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetTitle(Data.Author).
							SetThumbnail(engine.GetAuthorAvatar(Data.Author)).
							SetDescription(RemovePic(Data.Text)).
							SetURL(Data.PermanentURL).
							SetImage(Data.Photos...).
							SetColor(Color).
							InlineAllFields().
							SetFooter(Data.State, config.TwitterIMG).MessageEmbed
					} else {
						embed = engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetTitle(Data.Author).
							SetDescription(Data.Text).
							SetURL(Data.PermanentURL).
							SetImage(config.PixivProxy+Data.Photos[0]).
							SetColor(Color).
							InlineAllFields().
							SetFooter(Data.State, config.PixivIMG).MessageEmbed

					}
					err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Embeds: []*discordgo.MessageEmbed{
								embed,
							},
						},
					})
					if err != nil {
						log.Error(err)
					}

				}

				var Agency database.Group
				var VtuberName database.Member
				for _, v := range GroupsPayload {
					if strings.EqualFold(SelectedAgency, v.GroupName) {
						Agency = v
						if SelectedVtuber != "" {
							for _, k := range v.Members {
								if engine.CheckVtuberName(k, SelectedVtuber) {
									VtuberName = k
								}
							}

							if VtuberName.IsMemberNill() {
								SendSimpleMsg(s, i, fmt.Sprintf("Invalid vtuber name %s", SelectedVtuber))
								log.WithFields(log.Fields{
									"ChannelID":     i.ChannelID,
									"SelectedGroup": SelectedVtuber,
									"State":         "Lewd",
								}).Warn("Invalid vtuber name")
								return
							}
						}
					}
				}

				if Agency.IsNull() {
					SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", Agency.GroupName))
					log.WithFields(log.Fields{
						"ChannelID":     i.ChannelID,
						"SelectedGroup": SelectedAgency,
						"State":         "Lewd",
					}).Warn("Invalid agency name")
					return
				}

				if !VtuberName.IsMemberNill() {
					FanArt, err := VtuberName.GetRandomLewd()
					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Oops,something error\n" + err.Error(),
							},
						})
					} else {
						SendNude(FanArt)
					}

				} else {
					FanArt, err := Agency.GetRandomFanart()
					if err != nil {
						s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "Oops,something error\n" + err.Error(),
							},
						})
					} else {
						SendNude(FanArt)
					}

				}
			} else {
				err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							engine.NewEmbed().
								SetDescription("i know you horny,but this channel was not a NSFW channel").
								SetImage(engine.LewdIMG()).MessageEmbed,
						},
					},
				})
				if err != nil {
					log.Error(err)
				}
			}
		},
		"info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			VtuberInput := i.ApplicationCommandData().Options[0].StringValue()
			for _, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if engine.CheckVtuberName(v2, VtuberInput) {
						var (
							Avatar string
						)

						log.WithFields(log.Fields{
							"User":    i.Member.User.Username,
							"Channel": i.ChannelID,
							"Vtuber":  v2.Name,
						}).Info("info")

						SubsData, err := v2.GetSubsCount()
						if err != nil {
							log.Error(err)
						}

						if gacha() {
							Avatar = v2.YoutubeAvatar
						} else {
							if v2.BiliBiliRoomID != 0 {
								Avatar = v2.BiliBiliAvatar
							} else {
								Avatar = v2.YoutubeAvatar
							}
						}

						Color, err := engine.GetColor(config.TmpDir, i.Member.User.Avatar)
						if err != nil {
							log.Error(err)
						}
						YTSubs := "[Youtube](https://www.youtube.com/channel/" + v2.YoutubeID + "?sub_confirmation=1)"
						BiliFollow := "[BiliBili](https://space.bilibili.com/" + strconv.Itoa(v2.BiliBiliID) + ")"
						TwitterFollow := "[Twitter](https://twitter.com/" + v2.TwitterName + ")"
						if SubsData.BiliFollow != 0 && SubsData.BiliViews != 0 && SubsData.BiliVideos != 0 {
							err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Embeds: []*discordgo.MessageEmbed{
										engine.NewEmbed().
											SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
											SetTitle(engine.FixName(v2.EnName, v2.JpName)).
											SetImage(Avatar).
											AddField("Youtube subscribers", engine.NearestThousandFormat(float64(SubsData.YtSubs))).
											AddField("Youtube viewers", engine.NearestThousandFormat(float64(SubsData.YtViews))).
											AddField("Youtube videos", engine.NearestThousandFormat(float64(SubsData.YtVideos))).
											AddField("BiliBili followers", engine.NearestThousandFormat(float64(SubsData.YtViews))).
											AddField("BiliBili viewers", engine.NearestThousandFormat(float64(SubsData.BiliViews))).
											AddField("BiliBili videos", engine.NearestThousandFormat(float64(SubsData.BiliVideos))).
											InlineAllFields().
											AddField("Twitter followers", engine.NearestThousandFormat(float64(SubsData.TwFollow))).
											RemoveInline().
											AddField("<:yt:796023828723269662>", YTSubs).
											AddField("<:bili:796025336542265344>", BiliFollow).
											AddField("<:tw:796025611210588187>", TwitterFollow).
											InlineAllFields().
											SetColor(Color).MessageEmbed,
									},
								},
							})
							if err != nil {
								log.Error(err)
							}

						} else {
							err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Embeds: []*discordgo.MessageEmbed{
										engine.NewEmbed().
											SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
											SetTitle(engine.FixName(v2.EnName, v2.JpName)).
											SetImage(Avatar).
											AddField("Youtube subscribers", engine.NearestThousandFormat(float64(SubsData.YtSubs))).
											AddField("Youtube viewers", engine.NearestThousandFormat(float64(SubsData.YtViews))).
											AddField("Youtube videos", engine.NearestThousandFormat(float64(SubsData.YtVideos))).
											InlineAllFields().
											AddField("Twitter followers", engine.NearestThousandFormat(float64(SubsData.TwFollow))).
											RemoveInline().
											AddField("<:yt:796023828723269662>", YTSubs).
											AddField("<:tw:796025611210588187>", TwitterFollow).
											InlineAllFields().
											SetColor(Color).MessageEmbed,
									},
								},
							})
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}

			SendSimpleMsg(s, i, fmt.Sprintf("Invalid vtuber name %s", VtuberInput))
		},
		"tag-me": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if len(i.ApplicationCommandData().Options) < 1 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Null payload,make sure your command",
					},
				})
				return
			}

			var (
				Already        []string
				Done           []string
				SelectedAgency string
				Reminder       int
				SelectedVtuber string
			)

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("tag-me")

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}
				if v.Name == "reminder" {
					Reminder = int(v.IntValue())
				}
				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}
			}

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.Avatar)
			if err != nil {
				log.Error(err)
			}

			One := true
			for _, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if (SelectedAgency == "" && engine.CheckVtuberName(v2, SelectedVtuber)) || (strings.EqualFold(SelectedAgency, v.GroupName) && SelectedAgency != "") {
						if database.CheckChannelEnable(i.ChannelID, v2.Name, v2.Group.ID) {
							if (Reminder > 60 && Reminder < 10) && Reminder != 0 {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Cannot set remnder over 60 minutes",
									},
								})
								return
							}

							User := &database.UserStruct{
								DiscordID:       i.Member.User.ID,
								DiscordUserName: i.Member.User.Username,
								Channel_ID:      i.ChannelID,
								Human:           true,
								Reminder:        Reminder,
								Group:           v,
							}
							err := User.SetMember(v2).Adduser()
							if err != nil {
								log.Error(err)
								Already = append(Already, v2.Name)
							} else {
								Done = append(Done, v2.Name)
							}
							if One {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Processing",
									},
								})
								One = false
							}
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "look like this channel not enable " + v.GroupName,
								},
							})
							return
						}
					}
				}
			}
			if Already != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(i.Member.User.Mention() + " Already Added\n" + strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

			if Done != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(i.Member.User.Mention() + " notifications have been added to these members\n" + strings.Join(Done, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}
		},
		"del-tag": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var (
				Already        []string
				Done           []string
				SelectedAgency string
				SelectedVtuber string
			)
			One := true

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("add-tag")

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.Avatar)
			if err != nil {
				log.Error(err)
			}

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}
				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}
			}

			for _, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if (SelectedAgency == "" && engine.CheckVtuberName(v2, SelectedVtuber)) || (strings.EqualFold(SelectedAgency, v.GroupName) && SelectedAgency != "") {
						if database.CheckChannelEnable(i.ChannelID, v2.Name, v2.Group.ID) {

							User := &database.UserStruct{
								DiscordID:       i.Member.User.ID,
								DiscordUserName: i.Member.User.Username,
								Channel_ID:      i.ChannelID,
								Human:           true,
								Group:           v,
							}
							err := User.SetMember(v2).Deluser()
							if err != nil {
								Already = append(Already, v2.Name)
							} else {
								Done = append(Done, v2.Name)
							}

							if One {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Processing",
									},
								})
								One = false
							}
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "look like this channel not enable " + v.GroupName,
								},
							})
							return
						}
					}
				}
			}
			if Already != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(i.Member.User.Mention() + " You already removed this Group/Member from your list, or you never added them.\n" + strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

			if Done != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(i.Member.User.Mention() + " You removed these Members from your list.\n" + strings.Join(Done, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}
		},
		"tag-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Admin, err := MemberHasPermission(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Error(err)
			}
			if !Admin {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have permission to enable/disable/update,make sure you have `Manage Channels` Role permission",
					},
				})
				return
			}
			if len(i.ApplicationCommandData().Options) < 1 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Null payload,make sure your command",
					},
				})
				return
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("tag-role")

			var (
				Already        []string
				Done           []string
				SelectedAgency string
				Reminder       int
				SelectedVtuber string
				RoleState      *discordgo.Role
			)

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}
				if v.Name == "reminder" {
					Reminder = int(v.IntValue())
				}
				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}
				if v.Name == "role-name" {
					RoleState = v.RoleValue(nil, "")
				}
			}

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.Avatar)
			if err != nil {
				log.Error(err)
			}

			One := true
			for _, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if (SelectedAgency == "" && engine.CheckVtuberName(v2, SelectedVtuber)) || (strings.EqualFold(SelectedAgency, v.GroupName) && SelectedAgency != "") {
						if database.CheckChannelEnable(i.ChannelID, v2.Name, v2.Group.ID) {
							if (Reminder > 60 && Reminder < 10) && Reminder != 0 {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Cannot set remnder over 60 minutes",
									},
								})
								return
							}

							User := &database.UserStruct{
								DiscordID:       RoleState.ID,
								DiscordUserName: RoleState.Name,
								Channel_ID:      i.ChannelID,
								Human:           false,
								Reminder:        Reminder,
								Group:           v,
							}
							err := User.SetMember(v2).Adduser()
							if err != nil {
								log.Error(err)
								Already = append(Already, v2.Name)
							} else {
								Done = append(Done, v2.Name)
							}
							if One {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Processing",
									},
								})
								One = false
							}
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "look like this channel not enable " + v.GroupName,
								},
							})
							return
						}
					}
				}
			}
			if Already != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(RoleState.Mention() + " Already Added\n" + strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

			if Done != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(RoleState.Mention() + " notifications have been added to these members\n" + strings.Join(Done, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}
		},
		"del-role": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Admin, err := MemberHasPermission(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Error(err)
			}
			if !Admin {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have permission to enable/disable/update,make sure you have `Manage Channels` Role permission",
					},
				})
				return
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("del-role")

			One := true

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.Avatar)
			if err != nil {
				log.Error(err)
			}

			var (
				Already        []string
				Done           []string
				SelectedAgency string
				SelectedVtuber string
				RoleState      *discordgo.Role
			)

			for _, v := range i.ApplicationCommandData().Options {
				if v.Name == engine.AgencyOption {
					SelectedAgency = v.StringValue()
				}
				if v.Name == engine.VtuberOption {
					SelectedVtuber = v.StringValue()
				}
				if v.Name == "role-name" {
					RoleState = v.RoleValue(nil, "")
				}
			}

			for _, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if (SelectedAgency == "" && engine.CheckVtuberName(v2, SelectedVtuber)) || (strings.EqualFold(v.GroupName, SelectedAgency) && SelectedAgency != "") {
						if database.CheckChannelEnable(i.ChannelID, v2.Name, v2.Group.ID) {

							User := &database.UserStruct{
								DiscordID:       RoleState.ID,
								DiscordUserName: RoleState.Name,
								Channel_ID:      i.ChannelID,
								Human:           false,
								Group:           v,
							}
							err := User.SetMember(v2).Deluser()
							if err != nil {
								Already = append(Already, v2.Name)
							} else {
								Done = append(Done, v2.Name)
							}

							if One {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Processing",
									},
								})
								One = false
							}
						} else {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: "look like this channel not enable " + v.GroupName,
								},
							})
							return
						}
					}
				}
			}
			if Already != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(RoleState.Mention() + " You already removed this Group/Member from your list, or you never added them.\n" + strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}

			if Done != nil {
				_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
					Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
						SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
						SetDescription(RoleState.Mention() + " You removed these Members from your list.\n" + strings.Join(Done, " ")).
						SetThumbnail(config.GoSimpIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed},
				})
				if err != nil {
					s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
			}
		},
		"role-info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			role := i.ApplicationCommandData().Options[0].RoleValue(nil, "")

			list, err := database.UserStatus(role.ID, i.ChannelID)
			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("role-info")

			Color, err := engine.GetColor(config.TmpDir, i.Member.User.AvatarURL("128"))
			if err != nil {
				log.Error(err)
			}

			if list != nil {
				tableString := &strings.Builder{}
				table := tablewriter.NewWriter(tableString)
				table.SetAutoWrapText(false)
				table.SetAutoFormatHeaders(true)
				table.SetCenterSeparator("")
				table.SetColumnSeparator("")
				table.SetRowSeparator("")
				table.SetHeaderLine(true)
				table.SetBorder(false)
				table.SetTablePadding("\t")
				table.SetNoWhiteSpace(true)
				table.SetHeader([]string{"Vtuber Group", "Vtuber Name", "Reminder"})
				table.AppendBulk(list)
				table.Render()

				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
							SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
							SetThumbnail(i.Member.User.AvatarURL("128")).
							SetDescription("```" + tableString.String() + "```").
							SetColor(Color).MessageEmbed},
					},
				})
				if err != nil {
					log.Error(err)
				}
			} else {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
							SetDescription("Your tag list is empty.").
							SetTitle("404 Not found").
							SetImage(engine.NotFoundIMG()).
							SetColor(Color).MessageEmbed},
					},
				})
				if err != nil {
					log.Error(err)
				}
			}
		},
		"setup": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Admin, err := MemberHasPermission(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Error(err)
			}
			if !Admin {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have permission to enable/disable/update,make sure you have `Manage Channels` Role permission",
					},
				})
				return
			}
			log.WithFields(log.Fields{
				"Admin":     Admin,
				"ChannelID": i.ChannelID,
			}).Info("request for setup channel")

			Add := func(ChannelData *database.DiscordChannel, Group database.Group) {
				if ChannelData.ChannelCheck() {
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Already setup `" + Group.GroupName + "`,for add/del region use `Update` command",
						},
					})
					return
				}

				err := ChannelData.AddChannel()
				if err != nil {
					log.Error(err)
					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: err.Error(),
						},
					})
					return
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Done",
					},
				})

			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": i.ChannelID,
			}).Info("setup")

			for _, v := range i.ApplicationCommandData().Options[0].Options {
				var (
					Channel                                                            *discordgo.Channel
					Group                                                              database.Group
					liveonly, newupcoming, dynamic, liteMode, indieNotif, fanart, lewd bool
					region                                                             string
				)
				if v.Name == "livestream" {
					for _, v2 := range v.Options {
						if v2.Name == "channel-name" {
							Channel = v2.ChannelValue(nil)
						}

						if v2.Name == engine.AgencyOption {
							for _, v3 := range GroupsPayload {
								if strings.EqualFold(v3.GroupName, v2.StringValue()) {
									Group = v3
									break
								}
							}
						}

						if Group.IsNull() {
							SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", Group.GroupName))
							log.WithFields(log.Fields{
								"Admin":         Admin,
								"ChannelID":     i.ChannelID,
								"SelectedGroup": Group.GroupName,
								"State":         "Livestream",
							}).Warn("Invalid agency name")
							return
						}

						log.WithFields(log.Fields{
							"Admin":     Admin,
							"ChannelID": i.ChannelID,
							"Group":     Group.GroupName,
							"State":     "Livestream",
						}).Info("setup livestream")

						if v2.Name == "liveonly" {
							liveonly = v2.BoolValue()
						}

						if v2.Name == "newupcoming" {
							newupcoming = v2.BoolValue()
						}

						if v2.Name == "dynamic" {
							dynamic = v2.BoolValue()
						}

						if v2.Name == "lite-mode" {
							liteMode = v2.BoolValue()
						}

						if v2.Name == "indie-notif" {
							indieNotif = v2.BoolValue()
						}

						if v2.Name == "fanart" {
							fanart = v2.BoolValue()
						}
					}

					for Key, Val := range RegList {
						if Key == Group.GroupName {
							region = Val
						}
					}

					ChannelData := &database.DiscordChannel{
						ChannelID: Channel.ID,
						TypeTag: func() int {
							if fanart {
								return 3
							} else {
								return 2
							}
						}(),
						LiveOnly:    liveonly,
						NewUpcoming: newupcoming,
						Dynamic:     dynamic,
						LiteMode:    liteMode,
						IndieNotif: func() bool {
							if Group.GroupName != config.Indie {
								return false
							}
							return indieNotif
						}(),
						Group:  Group,
						Region: region,
					}

					Add(ChannelData, Group)

				} else if v.Name == "fanart" {
					for _, v2 := range v.Options {
						if v2.Name == "channel-name" {
							Channel = v2.ChannelValue(nil)
						}

						if v2.Name == engine.AgencyOption {
							for _, v3 := range GroupsPayload {
								if strings.EqualFold(v3.GroupName, v2.StringValue()) {
									Group = v3
									break
								}
							}
						}

						if Group.IsNull() {
							SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", Group.GroupName))
							log.WithFields(log.Fields{
								"Admin":         Admin,
								"ChannelID":     i.ChannelID,
								"SelectedGroup": Group.GroupName,
								"State":         "Setup",
							}).Warn("Invalid agency name")
							return
						}

						if v2.Name == "lewd" {
							lewd = v2.BoolValue()
						}

						log.WithFields(log.Fields{
							"Admin":     Admin,
							"ChannelID": i.ChannelID,
							"Group":     Group.GroupName,
							"State":     "Fanart",
							"Lewd":      lewd,
						}).Info("setup fanart")

						if lewd {
							ChannelRaw, err := s.Channel(Channel.ID)
							if err != nil {
								log.Error(err)
							}
							if !ChannelRaw.NSFW {
								err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Embeds: []*discordgo.MessageEmbed{
											engine.NewEmbed().
												SetDescription("i know you horny,but this channel was not a NSFW channel").
												SetImage(engine.LewdIMG()).MessageEmbed,
										},
									},
								})
								if err != nil {
									log.Error(err)
								}
							}
						}

						ChannelData := &database.DiscordChannel{
							ChannelID: Channel.ID,
							TypeTag: func() int {
								if lewd {
									return 70
								} else {
									return 1
								}
							}(),
							LiveOnly:    liveonly,
							NewUpcoming: newupcoming,
							Dynamic:     dynamic,
							LiteMode:    liteMode,
							IndieNotif: func() bool {
								if Group.GroupName != config.Indie {
									return false
								}
								return indieNotif
							}(),
							Group:  Group,
							Region: region,
						}
						Add(ChannelData, Group)
					}
				} else {
					for _, v2 := range v.Options {
						if v2.Name == "channel-name" {
							Channel = v2.ChannelValue(nil)
						}

						if v2.Name == engine.AgencyOption {
							for _, v3 := range GroupsPayload {
								if strings.EqualFold(v3.GroupName, v2.StringValue()) {
									Group = v3
									break
								}
							}
						}

						if Group.IsNull() {
							SendSimpleMsg(s, i, fmt.Sprintf("Invalid agency name %s", Group.GroupName))
							log.WithFields(log.Fields{
								"Admin":         Admin,
								"ChannelID":     i.ChannelID,
								"SelectedGroup": Group.GroupName,
								"State":         "Setup",
							}).Warn("Invalid agency name")
							return
						}

						log.WithFields(log.Fields{
							"Admin":     Admin,
							"ChannelID": i.ChannelID,
							"Group":     Group.GroupName,
							"State":     "Fanart",
							"Lewd":      lewd,
						}).Info("setup lewd")

						ChannelRaw, err := s.Channel(Channel.ID)
						if err != nil {
							log.Error(err)
						}
						if !ChannelRaw.NSFW {
							err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Embeds: []*discordgo.MessageEmbed{
										engine.NewEmbed().
											SetDescription("i know you horny,but this channel was not a NSFW channel").
											SetImage(engine.LewdIMG()).MessageEmbed,
									},
								},
							})
							if err != nil {
								log.Error(err)
							}
						}
						ChannelData := &database.DiscordChannel{
							ChannelID:   Channel.ID,
							TypeTag:     69,
							LiveOnly:    liveonly,
							NewUpcoming: newupcoming,
							Dynamic:     dynamic,
							LiteMode:    liteMode,
							IndieNotif: func() bool {
								if Group.GroupName != config.Indie {
									return false
								}
								return indieNotif
							}(),
							Group:  Group,
							Region: region,
						}
						Add(ChannelData, Group)
					}
				}
			}
		},
		"channel-state": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Color, err := engine.GetColor(config.TmpDir, i.Member.User.AvatarURL("128"))
			if err != nil {
				log.Error(err)
			}

			var (
				Typestr     string
				LiveOnly    = config.No
				NewUpcoming = config.No
				Dynamic     = config.No
				LiteMode    = config.No
				Indie       = ""
				Region      = "All"
				embed       []*discordgo.MessageEmbed
				Channel     = i.ApplicationCommandData().Options[0].ChannelValue(nil)
			)
			ChannelData, err := database.ChannelStatus(Channel.ID)
			if err != nil {
				log.Error(err)
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Processing",
				},
			})

			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": Channel.ID,
			}).Info("channel-state")

			if len(ChannelData) > 0 {
				for _, Channel := range ChannelData {
					ChannelRaw, err := s.Channel(Channel.ChannelID)
					if err != nil {
						log.Error(err)
					}

					if Channel.Region != "" {
						Region = Channel.Region
					}
					if Channel.IndieNotif && Channel.Group.GroupName == config.Indie {
						Indie = config.Ok
					} else if Channel.Group.GroupName != config.Indie {
						Indie = "-"
					} else {
						Indie = config.No
					}

					if Channel.IsFanart() && !Channel.IsLewd() && !Channel.IsLive() {
						Typestr = "Fanart"
					} else if !Channel.IsFanart() && !Channel.IsLewd() && Channel.IsLive() {
						Typestr = "Live"
					} else if Channel.IsFanart() && !Channel.IsLewd() && Channel.IsLive() {
						Typestr = "Fanart & Livestream"
					} else if Channel.IsLewd() && !Channel.IsFanart() {
						Typestr = "Lewd"
					} else if Channel.IsLewd() && Channel.IsFanart() {
						Typestr = "Fanart & Lewd"
					}

					if Channel.LiveOnly {
						LiveOnly = config.Ok
					}

					if Channel.NewUpcoming {
						NewUpcoming = config.Ok
					}

					if Channel.Dynamic {
						Dynamic = config.Ok
					}

					if Channel.LiteMode {
						LiteMode = config.Ok
					}

					if Channel.IsFanart() && !Channel.IsLewd() && !Channel.IsLive() {
						if Channel.Group.GroupName == config.Indie {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("Regions", Region).
								AddField("Independent notif", Indie).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						} else {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("Regions", Region).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						}

					} else if Channel.IsLewd() && !Channel.IsFanart() {
						if Channel.Group.GroupName == config.Indie {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("Regions", Region).
								AddField("Independent notif", Indie).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						} else {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("Regions", Region).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						}

					} else {
						if Channel.Group.GroupName == config.Indie {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("LiveOnly", LiveOnly).
								AddField("Dynamic", Dynamic).
								AddField("Upcoming", NewUpcoming).
								AddField("Lite", LiteMode).
								AddField("Regions", Region).
								AddField("Independent notif", Indie).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						} else {
							embed = append(embed, engine.NewEmbed().
								SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
								SetThumbnail(config.GoSimpIMG).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle(ChannelRaw.Name).
								AddField("Type", Typestr).
								AddField("LiveOnly", LiveOnly).
								AddField("Dynamic", Dynamic).
								AddField("Upcoming", NewUpcoming).
								AddField("Lite", LiteMode).
								AddField("Regions", Region).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
						}
					}
				}
				if len(ChannelData) >= 10 {
					for _, v := range embed {
						_, err := s.ChannelMessageSendEmbed(i.ChannelID, v)
						if err != nil {
							log.Error(err)
						}
					}
				} else {
					_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Embeds: embed,
					})
					if err != nil {
						s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
							Content: "Something went wrong",
						})
						return
					}
				}
			} else {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
							SetTitle("404 Not found").
							SetThumbnail(config.GoSimpIMG).
							SetImage(engine.NotFoundIMG()).
							SetColor(Color).MessageEmbed},
					},
				})
			}
		},
		"channel-update": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Admin, err := MemberHasPermission(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Error(err)
			}
			if !Admin {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have permission to enable/disable/update,make sure you have `Manage Channels` Role permission",
					},
				})
				return
			}

			Channel := i.ApplicationCommandData().Options[0].ChannelValue(nil)

			var (
				Typestr     string
				LiveOnly    = config.No
				NewUpcoming = config.No
				Dynamic     = config.No
				LiteMode    = config.No
				Indie       = ""
				Region      = "All"
			)
			ChannelData, err := database.ChannelStatus(Channel.ID)
			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"User":    i.Member.User.Username,
				"Channel": Channel.ID,
			}).Info("channel-Update")

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Processing",
				},
			})
			if err != nil {
				log.Error(err)
			}

			if len(ChannelData) > 0 {
				for _, Channel := range ChannelData {
					if Channel.Region != "" {
						Region = Channel.Region
					}

					if Channel.IsFanart() {
						Typestr = "FanArt"
					}

					if Channel.IsLive() {
						Typestr = "Live"
					}

					if Channel.IsFanart() && Channel.IsLive() {
						Typestr = "FanArt & Livestream"
					}

					if Channel.IsLewd() {
						Typestr = "Lewd"
					}

					if Channel.IsLewd() && Channel.IsFanart() {
						Typestr = "FanArt & Lewd"
					}

					if Channel.LiveOnly {
						LiveOnly = config.Ok
					}

					if Channel.NewUpcoming {
						NewUpcoming = config.Ok
					}

					if Channel.Dynamic {
						Dynamic = config.Ok
					}

					if Channel.LiteMode {
						LiteMode = config.Ok
					}

					if Channel.Group.GroupName == config.Indie {
						if Channel.IndieNotif {
							Indie = config.Ok
						} else if Channel.Group.GroupName != config.Indie {
							Indie = "-"
						} else {
							Indie = config.No
						}
						Channel.Group.RemoveNillIconURL()

						_, err = s.ChannelMessageSendEmbed(i.ChannelID, engine.NewEmbed().
							SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
							SetThumbnail(Channel.Group.IconURL).
							SetDescription("Channel States of "+Channel.Group.GroupName).
							SetTitle("ID "+strconv.Itoa(int(Channel.ID))).
							AddField("Type", Typestr).
							AddField("LiveOnly", LiveOnly).
							AddField("Dynamic", Dynamic).
							AddField("Upcoming", NewUpcoming).
							AddField("Lite", LiteMode).
							AddField("Regions", Region).
							AddField("Independent notif", Indie).
							InlineAllFields().MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					} else {
						_, err = s.ChannelMessageSendEmbed(i.ChannelID, engine.NewEmbed().
							SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
							SetThumbnail(Channel.Group.IconURL).
							SetDescription("Channel States of "+Channel.Group.GroupName).
							SetTitle("ID "+strconv.Itoa(int(Channel.ID))).
							AddField("Type", Typestr).
							AddField("LiveOnly", LiveOnly).
							AddField("Dynamic", Dynamic).
							AddField("Upcoming", NewUpcoming).
							AddField("Lite", LiteMode).
							AddField("Regions", Region).
							InlineAllFields().MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSendEmbed(i.ChannelID, engine.NewEmbed().
					SetTitle("404 Not found,use `/setup` first").
					SetThumbnail(config.GoSimpIMG).
					SetImage(engine.NotFoundIMG()).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
				return
			}

			AdminID := i.Member.User.ID
			Register := &ChannelRegister{
				AdminID:       AdminID,
				State:         UpdateState,
				ChannelStates: ChannelData,
			}
			_, err = s.ChannelMessageSend(i.ChannelID, "Select ID : ")
			if err != nil {
				log.Error(err)
			}
			Counter := 0
			Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
				if m.Author.ID == Register.AdminID && Register.State == UpdateState {
					Counter++
					if strings.ToLower(m.Content) == "exit" {
						Clear(Register)
						return
					}

					tableString := &strings.Builder{}
					table := tablewriter.NewWriter(tableString)
					table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
					table.SetCenterSeparator("|")

					tmp, err := strconv.Atoi(m.Content)
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "Worng input ID")
						if err != nil {
							log.Error(err)
						}
						return
					} else {
						for _, ChannelState := range Register.ChannelStates {
							if int(ChannelState.ID) == tmp {
								Register.SetChannel(ChannelState)
							}
						}
						if Register.ChannelState.ID != 0 {
							Register.SetChannelID(m.ChannelID)
							_, err := s.ChannelMessageSend(m.ChannelID, "You selectd `"+Register.ChannelState.Group.GroupName+"` with ID `"+strconv.Itoa(int(Register.ChannelState.ID))+"`")
							if err != nil {
								log.Error(err)
							}
							table.SetHeader([]string{"Menu"})
							table.Append([]string{"Update Channel state"})
							table.Append([]string{"Add region in this channel"})
							table.Append([]string{"Delete region in this channel"})

							if Register.ChannelState.TypeTag == 2 || Register.ChannelState.TypeTag == 3 {
								table.Append([]string{"Change Livestream state"})
							}

							table.Render()
							MsgText, err := s.ChannelMessageSend(m.ChannelID, "```"+tableString.String()+"```")
							if err != nil {
								log.Error(err)
							}

							if Register.ChannelState.TypeTag == 2 || Register.ChannelState.TypeTag == 3 {
								err = engine.Reacting(map[string]string{
									"ChannelID": m.ChannelID,
									"State":     "Menu2",
									"MessageID": MsgText.ID,
								}, s)
								if err != nil {
									log.Error(err)
								}
							} else {
								err = engine.Reacting(map[string]string{
									"ChannelID": m.ChannelID,
									"State":     "Menu",
									"MessageID": MsgText.ID,
								}, s)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSend(m.ChannelID, "Channel ID not found")
							if err != nil {
								log.Error(err)
							}
							if Counter == 5 {
								Clear(Register)
							}
							return
						}
					}
				}
			})

			Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
				if Register != nil && Register.AdminID == m.UserID {
					EmojiUpdate(Register, s, m)
					EmojiHandler(Register, s, m)
				}
			})
		},
		"channel-delete": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			Admin, err := MemberHasPermission(i.GuildID, i.Member.User.ID)
			if err != nil {
				log.Error(err)
			}
			if !Admin {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You don't have permission to enable/disable/update,make sure you have `Manage Channels` Role permission",
					},
				})
				return
			}

			Channel := i.ApplicationCommandData().Options[0].ChannelValue(nil)
			ChannelData, err := database.ChannelStatus(Channel.ID)
			if err != nil {
				log.Error(err)
			}
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   1 << 6,
					Content: "Processing",
				},
			})
			if err != nil {
				log.Error(err)
			}

			for _, v := range ChannelData {
				if v.ChannelCheck() {
					log.WithFields(log.Fields{
						"User":    i.Member.User.Username,
						"Channel": v.ChannelID,
					}).Info("channel-delete")

					err := v.DelChannel()
					if err != nil {
						log.Error(err)
						_, err := s.ChannelMessageSend(i.ChannelID, "Something error XD")
						if err != nil {
							log.Error(err)
						}
						return
					}
					_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
						Content: "<@" + i.Member.User.ID + "> is disabled " + v.Group.GroupName + " from this channel",
					})
					if err != nil {
						s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
							Content: "Something went wrong",
						})
						return
					}
				}
			}
		},
		"prediction": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var tw, bl, yt bool

			if i.ApplicationCommandData().Options[0].StringValue() == "tw" {
				tw = true
			}

			if i.ApplicationCommandData().Options[0].StringValue() == "bl" {
				bl = true
			}

			if i.ApplicationCommandData().Options[0].StringValue() == "yt" {
				yt = true
			}

			VtuberInput := i.ApplicationCommandData().Options[1].StringValue()

			PredictionCount := func() int {
				if len(i.ApplicationCommandData().Options) > 2 {
					return int(i.ApplicationCommandData().Options[2].IntValue())
				} else {
					return 7
				}
			}()

			if PredictionCount < 0 || PredictionCount > 100 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Invalid number,must lower than 100",
					},
				})
				return
			}

			for k, v := range GroupsPayload {
				for _, v2 := range v.Members {
					if strings.EqualFold(VtuberInput, v2.EnName) || strings.EqualFold(VtuberInput, v2.JpName) || strings.EqualFold(VtuberInput, v2.Name) {
						var (
							tmp            string
							Avatar         string
							Url            string
							PredictionData int
							State          string
						)

						Color, err := engine.GetColor(config.TmpDir, i.Member.User.AvatarURL("128"))
						if err != nil {
							log.Error(err)
						}

						Data, err := v2.GetSubsCount()
						if err != nil {
							log.Error(err)
						}

						log.WithFields(log.Fields{
							"VtuberName": v2.Name,
							"User":       i.Member.User.Username,
							"State":      i.ApplicationCommandData().Options[0].StringValue(),
							"Limit":      PredictionCount,
						}).Info("Prediction")

						if tw && v2.TwitterName == "" {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: VtuberInput + " don't have twitter account",
								},
							})
							return
						} else if tw && v2.TwitterName != "" {
							PredictionData, err = engine.Prediction(v2, "Twitter", PredictionCount)
							if err != nil {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Something error\n" + err.Error(),
									},
								})
								return
							}
							tmp = strconv.Itoa(Data.TwFollow)
							Avatar = v2.YoutubeAvatar
							Url = "https://twitter.com/" + v2.TwitterName
							State = "Twitter"
						}

						if bl && v2.BiliBiliID == 0 {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: VtuberInput + " don't have bilibili account",
								},
							})
							return
						} else if bl && v2.BiliBiliID != 0 {
							PredictionData, err = engine.Prediction(v2, "BiliBili", PredictionCount)
							if err != nil {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Something error\n" + err.Error(),
									},
								})
								return
							}
							tmp = strconv.Itoa(Data.BiliFollow)
							Avatar = v2.BiliBiliAvatar
							Url = "https://space.bilibili/" + strconv.Itoa(v2.BiliBiliID)
							State = "BiliBili"
						}

						if yt && v2.YoutubeID == "" {
							s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
								Type: discordgo.InteractionResponseChannelMessageWithSource,
								Data: &discordgo.InteractionResponseData{
									Content: VtuberInput + " don't have youtube account",
								},
							})
							return
						} else if yt && v2.YoutubeID != "" {
							PredictionData, err = engine.Prediction(v2, "Youtube", PredictionCount)
							if err != nil {
								s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
									Type: discordgo.InteractionResponseChannelMessageWithSource,
									Data: &discordgo.InteractionResponseData{
										Content: "Something error\n" + err.Error(),
									},
								})
								return
							}

							tmp = strconv.Itoa(Data.YtSubs)
							Avatar = v2.YoutubeAvatar
							Url = "https://www.youtube.com/channel/" + v2.YoutubeID + "?sub_confirmation=1"
							State = "Youtube"

						}

						target := time.Now().AddDate(0, 0, PredictionCount)
						dateFormat := fmt.Sprintf("%s/%d", target.Month().String(), target.Day())
						now := time.Now()
						nowFormat := fmt.Sprintf("%s/%d", now.Month().String(), now.Day())
						Graph := "[Views Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22" + State + "%22%2C%20vtuber%3D%22" + v2.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=1w)"
						err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds: []*discordgo.MessageEmbed{engine.NewEmbed().
									SetAuthor(i.Member.User.Username, i.Member.User.AvatarURL("128")).
									SetTitle(engine.FixName(v2.EnName, v2.JpName)).
									SetURL(Url).
									AddField("Current "+State+" followes/subscriber("+nowFormat+")", tmp).
									AddField("Next "+strconv.Itoa(PredictionCount)+" days Prediction("+dateFormat+")", strconv.Itoa(PredictionData)).
									RemoveInline().
									AddField("Graph", Graph).
									InlineAllFields().
									SetThumbnail(Avatar).
									SetColor(Color).
									SetFooter("algorithm : Linear regression").MessageEmbed},
							},
						})
						if err != nil {
							log.Error(err)
						}
					}
				}
				if k == len(GroupsPayload)-1 {
					err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "Invalid vtuber name",
						},
					})
					if err != nil {
						log.Error(err)
					}
					return
				}
			}
		},
	}
)
