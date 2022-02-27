package main

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func TwitchMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := configfile.BotPrefix.Twitch
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 1 {
			Payload := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
			if CommandArray[0] == Prefix+Live {
				for _, FindGroupArry := range Payload {
					Agency, err := FindGropName(FindGroupArry)
					if err != nil {
						Member, err := FindVtuber(FindGroupArry)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							LiveTwitch, err := Member.GetTwitchLiveStream(config.LiveStatus)
							if err != nil {
								log.Error(err)
								SendError(map[string]string{
									"ChannelID": m.ChannelID,
									"Username":  m.Author.Username,
									"AvatarURL": m.Author.AvatarURL("128"),
								})
							}
							FixName := engine.FixName(Member.EnName, Member.JpName)
							if LiveTwitch.ID != 0 {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								FanBase := "simps"
								loc := engine.Zawarudo(Member.Region)
								diff := time.Now().In(loc).Sub(LiveTwitch.Schedul.In(loc))
								view, err := strconv.Atoi(LiveTwitch.Viewers)
								if err != nil {
									log.Error(err)
								}

								if Member.Fanbase != "" {
									FanBase = Member.Fanbase
								}

								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(Member.TwitchAvatar).
									SetImage(LiveTwitch.Thumb).
									SetURL("https://twitch.tv/"+Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveTwitch.Game).
									SetColor(Color).
									SetFooter(LiveTwitch.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a livestream right now").
									SetImage(engine.NotFoundIMG()).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						TwitchLive, err := Agency.GetTwitchLiveStream(config.LiveStatus)
						if err != nil {
							log.Error(err)
							SendError(map[string]string{
								"ChannelID": m.ChannelID,
								"Username":  m.Author.Username,
								"AvatarURL": m.Author.AvatarURL("128"),
							})
						}
						if TwitchLive != nil {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
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
								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								view, err := strconv.Atoi(LiveData.Viewers)
								if err != nil {
									log.Error(err)
								}

								if LiveData.Member.Fanbase != "" {
									FanBase = LiveData.Member.Fanbase
								}

								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+Agency.GroupName+"` doesn't have a livestream right now").
								SetImage(engine.NotFoundIMG()).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else if CommandArray[0] == Prefix+Past || CommandArray[0] == Prefix+"last" {
				for _, FindGroupArry := range Payload {
					Agency, err := FindGropName(FindGroupArry)
					if err != nil {
						Member, err := FindVtuber(FindGroupArry)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
							return
						} else {
							TwitchLive, err := Agency.GetTwitchLiveStream(config.PastStatus)
							if err != nil {
								log.Error(err)
								SendError(map[string]string{
									"ChannelID": m.ChannelID,
									"Username":  m.Author.Username,
									"AvatarURL": m.Author.AvatarURL("128"),
								})
							}
							FixName := engine.FixName(Member.JpName, Member.EnName)
							if TwitchLive != nil {
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}

								for _, LiveData := range TwitchLive {
									FanBase := "simps"
									loc := engine.Zawarudo(Member.Region)
									diff := LiveData.Schedul.In(loc).Sub(time.Now().In(loc))
									view, err := strconv.Atoi(LiveData.Viewers)
									if err != nil {
										log.Error(err)
									}

									if Member.Fanbase != "" {
										FanBase = Member.Fanbase
									}

									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(FixName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetThumbnail(Member.TwitchAvatar).
										SetImage(LiveData.Thumb).
										SetURL("https://twitch.tv/"+Member.TwitchName).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
										InlineAllFields().
										AddField("Game", LiveData.Game).
										SetColor(Color).
										SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("It looks like `"+FixName+"` doesn't have a Past livestream right now").
									SetImage(engine.NotFoundIMG()).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						TwitchLive, err := Agency.GetTwitchLiveStream(config.PastStatus)
						if err != nil {
							log.Error(err)
							SendError(map[string]string{
								"ChannelID": m.ChannelID,
								"Username":  m.Author.Username,
								"AvatarURL": m.Author.AvatarURL("128"),
							})
						}
						if TwitchLive != nil {
							Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
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

								diff := time.Now().In(loc).Sub(LiveData.Schedul.In(loc))
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(FixName).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(LiveData.Member.TwitchAvatar).
									SetImage(LiveData.Thumb).
									SetURL("https://twitch.tv/"+LiveData.Member.TwitchName).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Viewers", engine.NearestThousandFormat(float64(view))+" "+FanBase).
									InlineAllFields().
									AddField("Game", LiveData.Game).
									SetColor(Color).
									SetFooter(LiveData.Schedul.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("It looks like `"+Agency.GroupName+"` doesn't have a Past livestream right now").
								SetImage(engine.NotFoundIMG()).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
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
