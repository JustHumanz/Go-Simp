package main

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

//YoutubeMessage Youtube message handler
func YoutubeMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	var (
		Region string
		Prefix = configfile.BotPrefix.Youtube
	)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 3 && CommandArray[2] == "-region" {
			Region = strings.ToLower(CommandArray[3])
			if len(Region) > 2 {
				_, err := s.ChannelMessageSend(m.ChannelID, "Only support 1 Region,ignoring `"+Region[2:]+"`")
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			Region = ""
		}
		if (strings.ToLower(CommandArray[0]) == Prefix+Upcoming) || (strings.ToLower(CommandArray[0]) == Prefix+"up") {
			if len(CommandArray) > 1 {
				for _, GroupNameQuery := range strings.Split(strings.TrimSpace(CommandArray[1]), ",") {
					VTuberGroup, err := FindGropName(GroupNameQuery)
					if err != nil {
						Member := FindVtuber(GroupNameQuery)
						if Member.IsMemberNill() {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+GroupNameQuery+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
							return
						} else {
							YoutubeData, err := database.YtGetStatus(0, Member.ID, config.UpcomingStatus, Region, config.Fe)
							if err != nil {
								log.Error(err)
							}
							FixName := engine.FixName(Member.EnName, Member.JpName)
							if YoutubeData != nil {
								for _, Youtube := range YoutubeData {
									FanBase := "simps"
									loc := engine.Zawarudo(Member.Region)
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
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
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
									if err != nil {
										log.Error(err)
									}
								}

							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription("It looks like `"+FixName+"` doesn't have a livestream schedule for now").
									SetImage(engine.NotFoundIMG()).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						GroupData, err := database.YtGetStatus(VTuberGroup.ID, 0, config.UpcomingStatus, Region, config.Fe)
						if err != nil {
							log.Error(err)
						}
						if GroupData != nil {
							for _, Youtube := range GroupData {
								Youtube.AddMember(FindVtuber(Youtube.Member.ID))
								FixName := engine.FixName(Youtube.Member.EnName, Youtube.Member.JpName)
								FanBase := "simps"

								loc := engine.Zawarudo(Youtube.Member.Region)
								duration := durafmt.Parse(Youtube.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)

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
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Youtube.Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Start live in", duration.String()).
									AddField("Viewers", Viewers+" "+FanBase).
									InlineAllFields().
									AddField("Type", engine.YtFindType(Youtube.Title)).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(VTuberGroup.GroupName).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a livestream schedule for now").
								SetImage(engine.NotFoundIMG()).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetDescription("Incomplete `"+Upcoming+"` command").
					SetImage(engine.NotFoundIMG()).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+Live {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for _, GroupName := range FindGroupArry {
					VTuberGroup, err := FindGropName(GroupName)
					if err != nil {
						Member := FindVtuber(GroupName)
						if Member.IsMemberNill() {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+GroupName+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
						} else {
							YoutubeData, err := database.YtGetStatus(0, Member.ID, config.LiveStatus, Region, config.Fe)
							if err != nil {
								log.Error(err)
							}
							FixName := engine.FixName(Member.EnName, Member.JpName)

							if YoutubeData != nil {
								for _, Youtube := range YoutubeData {
									FanBase := "simps"
									loc := engine.Zawarudo(Member.Region)
									expiresAt := time.Now().In(loc)
									duration := durafmt.Parse(expiresAt.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)

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

									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetTitle(FixName).
										SetDescription(Youtube.Title).
										SetImage(Youtube.Thumb).
										SetThumbnail(Youtube.Member.YoutubeAvatar).
										SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
										AddField("Start live in", duration.String()+" Ago").
										AddField("Viewers", Viewers+" "+FanBase).
										InlineAllFields().
										AddField("Type", engine.YtFindType(Youtube.Title)).
										SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
										SetColor(Color).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription("It looks like `"+FixName+"` doesn't have a livestream schedule for now").
									SetImage(engine.NotFoundIMG()).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						YoutubeData, err := database.YtGetStatus(VTuberGroup.ID, 0, config.LiveStatus, Region, config.Fe)
						if err != nil {
							log.Error(err)
						}
						if YoutubeData != nil {
							for _, Youtube := range YoutubeData {
								Youtube.AddMember(FindVtuber(Youtube.Member.ID))

								FanBase := "simps"
								loc := engine.Zawarudo(Youtube.Member.Region)
								duration := durafmt.Parse(time.Now().In(loc).Sub(Youtube.Schedul.In(loc))).LimitFirstN(2)

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
								FixName := engine.FixName(Youtube.Member.EnName, Youtube.Member.JpName)

								view, err := strconv.Atoi(Youtube.Viewers)
								if err != nil {
									log.Error(err)
								}
								Viewers := engine.NearestThousandFormat(float64(view))
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetThumbnail(Youtube.Member.YoutubeAvatar).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Start in", duration.String()+" Ago").
									AddField("Viewers", Viewers+" "+FanBase).
									InlineAllFields().
									AddField("Type", engine.YtFindType(Youtube.Title)).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(VTuberGroup.GroupName).
								SetDescription("It looks like `"+VTuberGroup.GroupName+"` doesn't have a livestream schedule for now").
								SetImage(engine.NotFoundIMG()).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetDescription("Incomplete `"+Live+"` command").
					SetImage(engine.NotFoundIMG()).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
				return
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+"last" || strings.ToLower(CommandArray[0]) == Prefix+Past {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for _, GroupName := range FindGroupArry {
					VTuberGroup, err := FindGropName(GroupName)
					if err != nil {
						Member := FindVtuber(GroupName)
						if Member.IsMemberNill() {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+GroupName+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
							return
						} else {
							YoutubeData, err := database.YtGetStatus(0, Member.ID, config.PastStatus, Region, config.Fe)
							if err != nil {
								log.Error(err)
							}
							FixName := engine.FixName(Member.EnName, Member.JpName)
							if YoutubeData != nil {
								for _, Youtube := range YoutubeData {
									Color, err := engine.GetColor(config.TmpDir, Youtube.Thumb)
									if err != nil {
										log.Error(err)
									}

									view, err := strconv.Atoi(Youtube.Viewers)
									if err != nil {
										log.Error(err)
									}
									Viewers := engine.NearestThousandFormat(float64(view))

									loc := engine.Zawarudo(Member.Region)
									duration := durafmt.Parse(time.Now().In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
									durationlive := durafmt.Parse(Youtube.End.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)

									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
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
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetDescription("Internal error XD").
									SetImage(engine.NotFoundIMG()).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						var (
							YoutubeData []database.LiveStream
						)

						if CheckReg(VTuberGroup.GroupName, Region) {
							YoutubeData, err = database.YtGetStatus(VTuberGroup.ID, 0, config.PastStatus, Region, config.Fe)
							if err != nil {
								log.Error(err)
							}
						} else {
							YoutubeData, err = database.YtGetStatus(VTuberGroup.ID, 0, config.PastStatus, "", config.Fe)
							if err != nil {
								log.Error(err)
							}
						}

						if YoutubeData != nil {
							for _, Youtube := range YoutubeData {
								Youtube.AddMember(FindVtuber(Youtube.Member.ID))

								Color, err := engine.GetColor(config.TmpDir, Youtube.Thumb)
								if err != nil {
									log.Error(err)
								}
								view, err := strconv.Atoi(Youtube.Viewers)
								if err != nil {
									log.Error(err)
								}
								Viewers := engine.NearestThousandFormat(float64(view))
								FixName := engine.FixName(Youtube.Member.EnName, Youtube.Member.JpName)
								loc := engine.Zawarudo(Youtube.Member.Region)

								duration := durafmt.Parse(time.Now().In(loc).Sub(Youtube.Schedul.In(loc))).LimitFirstN(2)
								durationlive := durafmt.Parse(Youtube.End.In(loc).Sub(Youtube.Schedul)).LimitFirstN(2)
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(FixName).
									SetThumbnail(Youtube.Member.YoutubeAvatar).
									SetDescription(Youtube.Title).
									SetImage(Youtube.Thumb).
									SetURL("https://www.youtube.com/watch?v="+Youtube.VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Viewers).
									SetFooter(Youtube.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(VTuberGroup.GroupName).
								SetDescription("Internal error XD").
								SetImage(engine.NotFoundIMG()).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetDescription("Incomplete `"+Prefix+"` command").
					SetImage(engine.NotFoundIMG()).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
