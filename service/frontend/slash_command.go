package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "basic-command",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Basic command",
		},
		{
			Name:        "options",
			Description: "Command for demonstrating options",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "string-option",
					Description: "String option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "integer-option",
					Description: "Integer option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "bool-option",
					Description: "Boolean option",
					Required:    true,
				},

				// Required options must be listed first since optional parameters
				// always come after when they're used.
				// The same concept applies to Discord's Slash-commands API

				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel-option",
					Description: "Channel option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "User option",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Role option",
					Required:    false,
				},
			},
		},
		{
			Name:        "subcommands",
			Description: "Subcommands and command groups example",
			Options: []*discordgo.ApplicationCommandOption{
				// When a command has subcommands/subcommand groups
				// It must not have top-level options, they aren't accesible in the UI
				// in this case (at least not yet), so if a command has
				// subcommands/subcommand any groups registering top-level options
				// will cause the registration of the command to fail

				{
					Name:        "scmd-grp",
					Description: "Subcommands group",
					Options: []*discordgo.ApplicationCommandOption{
						// Also, subcommand groups aren't capable of
						// containing options, by the name of them, you can see
						// they can only contain subcommands
						{
							Name:        "nst-subcmd",
							Description: "Nested subcommand",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
					},
					Type: discordgo.ApplicationCommandOptionSubCommandGroup,
				},
				// Also, you can create both subcommand groups and subcommands
				// in the command at the same time. But, there's some limits to
				// nesting, count of subcommands (top level and nested) and options.
				// Read the intro of slash-commands docs on Discord dev portal
				// to get more information
				{
					Name:        "subcmd",
					Description: "Top-level subcommand",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
				},
			},
		},
		{
			Name:        "responses",
			Description: "Interaction responses testing initiative",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "resp-type",
					Description: "Response type",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices:     VtuberGroupChoices,
					Required:    true,
				},
			},
		},
		{
			Name:        "followups",
			Description: "Followup messages",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"basic-command": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there! Congratulations, you just executed your first slash command",
				},
			})
		},
		"livestream": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			state := i.ApplicationCommandData().Options[0].IntValue()
			Avatar := i.Member.User.AvatarURL("128")
			Nick := i.Member.User.Username
			var group database.Group
			var member database.Member
			var embed []*discordgo.MessageEmbed
			region := ""
			status := ""

			if i.ApplicationCommandData().Options[1].IntValue() == 1 {
				status = config.LiveStatus
			} else if i.ApplicationCommandData().Options[1].IntValue() == 2 {
				status = config.UpcomingStatus
			} else {
				status = config.PastStatus
			}

			SendMessage := func(e []*discordgo.MessageEmbed) {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: e,
					},
				})
				if err != nil {
					log.Error(err)
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
			for _, v := range *GroupsPayload {
				if v.ID == i.ApplicationCommandData().Options[2].IntValue() {
					group = v
				}
			}
			if len(i.ApplicationCommandData().Options) == 4 {
				for _, v := range *GroupsPayload {
					for _, v2 := range v.Members {
						if v2.ID == i.ApplicationCommandData().Options[3].IntValue() {
							member = v2
						}
					}
				}
			}
			if len(i.ApplicationCommandData().Options) == 5 {
				region = i.ApplicationCommandData().Options[4].StringValue()
			}

			if state == 1 {
				if !member.IsMemberNill() {
					YoutubeData, _, err := database.YtGetStatus(map[string]interface{}{
						"MemberID":   member.ID,
						"MemberName": member.Name,
						"Status":     status,
						"State":      config.Fe,
					})
					if err != nil {
						log.Error(err)
					}
					FixName := engine.FixName(member.EnName, member.JpName)
					if YoutubeData != nil {
						for _, Youtube := range YoutubeData {
							FanBase := "simps"
							loc := engine.Zawarudo(member.Region)
							duration := durafmt.Parse(Youtube.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)

							Color, err := engine.GetColor(config.TmpDir, Youtube.Thumb)
							if err != nil {
								if err.Error() == "Server Error,status get 404 Not Found" {
									Youtube.UpdateYt("private")
								} else {
									log.Error(err)
								}
							}

							if member.Fanbase != "" {
								FanBase = member.Fanbase
							}

							view, err := strconv.Atoi(Youtube.Viewers)
							if err != nil {
								log.Error(err)
							}
							Viewers := engine.NearestThousandFormat(float64(view))

							if status == config.PastStatus {
								durationlive := durafmt.Parse(Youtube.End.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, BotInfo.Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(member.YoutubeAvatar).
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
									SetThumbnail(member.YoutubeAvatar).
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
				} else {
					YoutubeData, _, err := database.YtGetStatus(map[string]interface{}{
						"GroupID":   group.ID,
						"GroupName": group.GroupName,
						"Status":    status,
						"Region":    region,
						"State":     config.Fe,
					})
					if err != nil {
						log.Error(err)
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

							loc := engine.Zawarudo(Youtube.Member.Region)
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
									SetAuthor(Nick, BotInfo.Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Viewers).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							} else {
								expiresAt := time.Now().In(loc)
								duration := durafmt.Parse(expiresAt.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								embed = append(embed, engine.NewEmbed().
									SetAuthor(Nick, Avatar).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(member.YoutubeAvatar).
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
							SetTitle(group.GroupName).
							SetDescription("It looks like `"+group.GroupName+"` doesn't have a `"+status+"` schedule for now").
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
				if !member.IsMemberNill() {
					LiveBili, err := database.BilGet(map[string]interface{}{
						"MemberID": member.ID,
						"Status":   status,
					})
					if err != nil {
						log.Error(err)
					}
					if LiveBili != nil {
						for _, LiveData := range LiveBili {
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
									SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliRoomID)).
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
									SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliRoomID)).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Online", engine.NearestThousandFormat(float64(view))).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
							}
						}
					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetDescription("It looks like `"+group.GroupName+"` doesn't have a `"+status+"` schedule right now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}
					SendMessage(embed)
				} else {
					LiveBili, err := database.BilGet(map[string]interface{}{
						"GroupID": group.ID,
						"Status":  status,
					})
					if err != nil {
						log.Error(err)
					}
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
									SetURL("https://live.bilibili.com/"+strconv.Itoa(Member.BiliRoomID)).
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
									SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveData.Member.BiliRoomID)).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Online", engine.NearestThousandFormat(float64(view))).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
							}
						}
					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetDescription("It looks like `"+group.GroupName+"` doesn't have a `"+status+"` schedule right now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}
					SendMessage(embed)
				}
			} else {
				if status == config.UpcomingStatus {
					NotSupportUp(status, "Twitch")
					return
				}
				if !member.IsMemberNill() {
					LiveTwitch, err := database.TwitchGet(map[string]interface{}{
						"MemberID": member.ID,
						"Status":   status,
					})
					if err != nil {
						log.Error(err)
					}
					FixName := engine.FixName(member.EnName, member.JpName)
					if LiveTwitch != nil {
						Color, err := engine.GetColor(config.TmpDir, Avatar)
						if err != nil {
							log.Error(err)
						}
						for _, LiveData := range LiveTwitch {
							FanBase := "simps"
							loc := engine.Zawarudo(member.Region)
							diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
							view, err := strconv.Atoi(LiveData.Viewers)
							if err != nil {
								log.Error(err)
							}

							if member.Fanbase != "" {
								FanBase = member.Fanbase
							}

							if LiveData.Game == "" {
								LiveData.Game = "???"
							}

							if status == config.PastStatus {
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								embed = append(embed, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(i.Member.Nick, Avatar).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)

							} else {
								embed = append(embed, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(i.Member.Nick, Avatar).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
							}

						}
					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetDescription("It looks like `"+FixName+"` doesn't have a `"+status+"` schedule right now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}
				} else {
					TwitchLive, err := database.TwitchGet(map[string]interface{}{
						"GroupID": group.ID,
						"Status":  status,
					})
					if err != nil {
						log.Error(err)
					}
					if TwitchLive != nil {
						Color, err := engine.GetColor(config.TmpDir, Avatar)
						if err != nil {
							log.Error(err)
						}

						for _, LiveData := range TwitchLive {
							FanBase := "simps"

							Member, err := FindVtuber(LiveData.Member.ID)
							if err != nil {
								log.Error(err)
							}

							LiveData.AddMember(Member)
							loc := engine.Zawarudo(LiveData.Member.Region)
							FixName := engine.FixName(LiveData.Member.EnName, LiveData.Member.JpName)
							view, err := strconv.Atoi(LiveData.Viewers)
							if err != nil {
								log.Error(err)
							}

							if LiveData.Member.Fanbase != "" {
								FanBase = LiveData.Member.Fanbase
							}

							if LiveData.Game == "" {
								LiveData.Game = "???"
							}

							if status == config.PastStatus {
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								embed = append(embed, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(Nick, Avatar).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)

							} else {
								diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
								embed = append(embed, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(Nick, Avatar).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
							}
						}
					} else {
						embed = append(embed, engine.NewEmbed().
							SetAuthor(Nick, Avatar).
							SetDescription("It looks like `"+group.GroupName+"` doesn't have a `"+status+"` schedule right now").
							SetImage(engine.NotFoundIMG()).MessageEmbed)
					}
				}
				SendMessage(embed)
			}
		},
		"art": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var embed *discordgo.MessageEmbed
			var DynamicData DynamicSvr

			Avatar := i.Member.User.AvatarURL("128")
			Nick := i.Member.User.Username

			SendNude := func(Data *database.DataFanart) {
				Color, err := engine.GetColor(config.TmpDir, Avatar)
				if err != nil {
					log.Error(err)
				}

				if Data.State == config.BiliBiliArt {
					body, errcurl := network.CoolerCurl("https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id="+Data.Dynamic_id, nil)
					if errcurl != nil {
						log.Error(errcurl)
					}
					json.Unmarshal(body, &DynamicData)
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
						SetImage(CDN+Data.Photos[0]).
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
			if len(i.ApplicationCommandData().Options) == 2 {
				VtuberName := i.ApplicationCommandData().Options[1].StringValue()
				log.WithFields(log.Fields{
					"Vtuber": VtuberName,
				}).Info("FanArt")
				for _, GroupData := range *GroupsPayload {
					for _, v := range GroupData.Members {
						if strings.EqualFold(v.Name, VtuberName) || strings.EqualFold(v.EnName, VtuberName) || strings.EqualFold(v.JpName, VtuberName) {
							FanArt, err := v.GetRandomFanart()
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
					}
				}
			} else {
				GroupID := i.ApplicationCommandData().Options[0].IntValue()
				log.WithFields(log.Fields{
					"VtuberGroupName": GroupID,
				}).Info("FanArt")
				for _, v := range *GroupsPayload {
					if v.ID == GroupID {
						FanArt, err := v.GetRandomFanart()
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
				}
			}
		},
		"options": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			margs := []interface{}{
				// Here we need to convert raw interface{} value to wanted type.
				// Also, as you can see, here is used utility functions to convert the value
				// to particular type. Yeah, you can use just switch type,
				// but this is much simpler
				i.ApplicationCommandData().Options[0].StringValue(),
				i.ApplicationCommandData().Options[1].IntValue(),
				i.ApplicationCommandData().Options[2].BoolValue(),
			}
			msgformat :=
				` Now you just learned how to use command options. Take a look to the value of which you've just entered:
				> string_option: %s
				> integer_option: %d
				> bool_option: %v
`
			if len(i.ApplicationCommandData().Options) >= 4 {
				margs = append(margs, i.ApplicationCommandData().Options[3].ChannelValue(nil).ID)
				msgformat += "> channel-option: <#%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 5 {
				margs = append(margs, i.ApplicationCommandData().Options[4].UserValue(nil).ID)
				msgformat += "> user-option: <@%s>\n"
			}
			if len(i.ApplicationCommandData().Options) >= 6 {
				margs = append(margs, i.ApplicationCommandData().Options[5].RoleValue(nil, "").ID)
				msgformat += "> role-option: <@&%s>\n"
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				// Ignore type for now, we'll discuss them in "responses" part
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(
						msgformat,
						margs...,
					),
				},
			})
		},
		"subcommands": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			content := ""

			// As you can see, the name of subcommand (nested, top-level) or subcommand group
			// is provided through arguments.
			switch i.ApplicationCommandData().Options[0].Name {
			case "subcmd":
				content =
					"The top-level subcommand is executed. Now try to execute the nested one."
			default:
				if i.ApplicationCommandData().Options[0].Name != "scmd-grp" {
					return
				}
				switch i.ApplicationCommandData().Options[0].Options[0].Name {
				case "nst-subcmd":
					content = "Nice, now you know how to execute nested commands too"
				default:
					// I added this in the case something might go wrong
					content = "Oops, something gone wrong.\n" +
						"Hol' up, you aren't supposed to see this message."
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Responses to a command are very important.
			// First of all, because you need to react to the interaction
			// by sending the response in 3 seconds after receiving, otherwise
			// interaction will be considered invalid and you can no longer
			// use the interaction token and ID for responding to the user's request

			content := ""
			// As you can see, the response type names used here are pretty self-explanatory,
			// but for those who want more information see the official documentation
			switch i.ApplicationCommandData().Options[0].IntValue() {
			case int64(discordgo.InteractionResponseChannelMessageWithSource):
				content =
					"You just responded to an interaction, sent a message and showed the original one. " +
						"Congratulations!"
				content +=
					"\nAlso... you can edit your response, wait 5 seconds and this message will be changed"
			default:
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				})
				if err != nil {
					s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
				}
				return
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseType(i.ApplicationCommandData().Options[0].IntValue()),
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
			if err != nil {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.AfterFunc(time.Second*5, func() {
				err = s.InteractionResponseEdit(s.State.User.ID, i.Interaction, &discordgo.WebhookEdit{
					Content: content + "\n\nWell, now you know how to create and edit responses. " +
						"But you still don't know how to delete them... so... wait 10 seconds and this " +
						"message will be deleted.",
				})
				if err != nil {
					s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
						Content: "Something went wrong",
					})
					return
				}
				time.Sleep(time.Second * 10)
				s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
			})
		},
		"followups": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Followup messages are basically regular messages (you can create as many of them as you wish)
			// but work as they are created by webhooks and their functionality
			// is for handling additional messages after sending a response.

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					// Note: this isn't documented, but you can use that if you want to.
					// This flag just allows you to create messages visible only for the caller of the command
					// (user who triggered the command)
					Flags:   1 << 6,
					Content: "Surprise!",
				},
			})
			msg, err := s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "Followup message has been created, after 5 seconds it will be edited",
			})
			if err != nil {
				s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
					Content: "Something went wrong",
				})
				return
			}
			time.Sleep(time.Second * 5)

			s.FollowupMessageEdit(s.State.User.ID, i.Interaction, msg.ID, &discordgo.WebhookEdit{
				Content: "Now the original message is gone and after 10 seconds this message will ~~self-destruct~~ be deleted.",
			})

			time.Sleep(time.Second * 10)

			s.FollowupMessageDelete(s.State.User.ID, i.Interaction, msg.ID)

			s.FollowupMessageCreate(s.State.User.ID, i.Interaction, true, &discordgo.WebhookParams{
				Content: "For those, who didn't skip anything and followed tutorial along fairly, " +
					"take a unicorn :unicorn: as reward!\n" +
					"Also, as bonus... look at the original interaction response :D",
			})
		},
	}
)
