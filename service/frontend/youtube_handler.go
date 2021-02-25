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
		Region    string
		Prefix    = configfile.BotPrefix.Youtube
		SendEmbed = func(Data Memberst) {
			if !Data.YtData.IsEmpty() {
				YTData := Data.YtData
				Color, err := engine.GetColor(config.TmpDir, YTData.Thumb)
				if err != nil {
					if err.Error() == "Server Error,status get 404 Not Found" {
						YTData.UpdateYt("private")
					} else {
						log.Error(err)
					}
				}

				_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Data.VTName).
					SetDescription(YTData.Title).
					SetImage(YTData.Thumb).
					SetThumbnail(YTData.YoutubeAvatar).
					SetURL("https://www.youtube.com/watch?v="+YTData.VideoID).
					AddField(Data.Msg, Data.Msg1).
					AddField("Viewers", Data.Msg3).
					InlineAllFields().
					AddField("Type", engine.YtFindType(YTData.Title)).
					SetFooter(Data.Msg2, config.YoutubeIMG).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.Error(err)
				}

			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Data.VTName).
					SetDescription(Data.Msg).
					SetImage(config.WorryIMG).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			}
		}
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
						VTData := ValidName(GroupNameQuery)
						if VTData.ID == 0 {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+GroupNameQuery+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
						} else {
							DataMember, err := database.YtGetStatus(0, VTData.ID, "upcoming", Region)
							if err != nil {
								log.Error(err)
							}
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := engine.Zawarudo(Member.Region)
									duration := durafmt.Parse(Member.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
									Views, err := strconv.Atoi(Member.Viewers)
									if err != nil {
										log.Error(err)
									}
									Member.Viewers = engine.NearestThousandFormat(float64(Views))
									SendEmbed(Memberst{
										VTName: VTData.VTName,
										YtData: Member,
										Msg:    "Start live in",
										Msg1:   duration.String(),
										Msg2:   Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:   Member.Viewers + " simps waiting in Room Chat",
									})
								}

							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream schedule for now",
								})
							}
						}
					} else {
						GroupData, err := database.YtGetStatus(VTuberGroup.ID, 0, "upcoming", Region)
						if err != nil {
							log.Error(err)
						}
						if GroupData != nil {
							for _, Member := range GroupData {
								loc := engine.Zawarudo(Member.Region)
								duration := durafmt.Parse(Member.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
								Views, err := strconv.Atoi(Member.Viewers)
								if err != nil {
									log.Error(err)
								}
								Member.Viewers = engine.NearestThousandFormat(float64(Views))
								SendEmbed(Memberst{
									VTName:    engine.FixName(Member.NameEN, Member.NameJP),
									YtData:    Member,
									Msg:       "Start in",
									Msg1:      duration.String(),
									Msg2:      Member.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Member.Viewers + " simps waiting in Room Chat",
									YtChannel: Member.ChannelID,
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.GroupName,
								Msg:    "It seems like `" + strings.Title(GroupNameQuery+" "+Region) + "` doesn't have livestream schedule for now",
							})
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Upcoming+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+Live {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
						} else {
							DataMember, err := database.YtGetStatus(0, VTData.ID, "live", Region)
							if err != nil {
								log.Error(err)
							}
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := engine.Zawarudo(Member.Region)
									expiresAt := time.Now().In(loc)
									duration := durafmt.Parse(expiresAt.In(loc).Sub(Member.Schedul)).LimitFirstN(2)
									Views, err := strconv.Atoi(Member.Viewers)
									if err != nil {
										log.Error(err)
									}
									Member.Viewers = engine.NearestThousandFormat(float64(Views))
									SendEmbed(Memberst{
										VTName: VTData.VTName,
										YtData: Member,
										Msg:    "Start live in",
										Msg1:   duration.String() + " Ago",
										Msg2:   Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:   Member.Viewers + " Simps",
									})
								}
							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream right now",
								})
							}
						}
					} else {
						GroupData, err := database.YtGetStatus(VTuberGroup.ID, 0, "live", Region)
						if err != nil {
							log.Error(err)
						}
						if GroupData != nil {
							for _, Member := range GroupData {
								loc := engine.Zawarudo(Member.Region)
								duration := durafmt.Parse(time.Now().In(loc).Sub(Member.Schedul.In(loc))).LimitFirstN(2)
								Views, err := strconv.Atoi(Member.Viewers)
								if err != nil {
									log.Error(err)
								}
								Member.Viewers = engine.NearestThousandFormat(float64(Views))
								SendEmbed(Memberst{
									VTName:    engine.FixName(Member.NameEN, Member.NameJP),
									YtData:    Member,
									Msg:       "Start in",
									Msg1:      duration.String() + " Ago",
									Msg2:      Member.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Member.Viewers + " Simps",
									YtChannel: Member.ChannelID,
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.GroupName,
								Msg:    "It looks like `" + strings.Title(FindGroupArry[i]+" "+Region) + "` doesn't have a livestream right now",
							})
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Live+"` command")
				if err != nil {
					log.Error(err)
				}
				return
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+"last" || strings.ToLower(CommandArray[0]) == Prefix+Past {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group or Vtuber Name was not found")
							if err != nil {
								log.Error(err)
							}
						} else {
							DataMember, err := database.YtGetStatus(0, VTData.ID, "past", Region)
							if err != nil {
								log.Error(err)
							}
							if DataMember != nil {
								for _, Member := range DataMember {
									Color, err := engine.GetColor(config.TmpDir, Member.Thumb)
									if err != nil {
										log.Error(err)
									}
									Views, err := strconv.Atoi(Member.Viewers)
									if err != nil {
										log.Error(err)
									}
									Member.Viewers = engine.NearestThousandFormat(float64(Views))

									loc := engine.Zawarudo(Member.Region)
									duration := durafmt.Parse(time.Now().In(loc).Sub(Member.Schedul)).LimitFirstN(2)
									durationlive := durafmt.Parse(Member.End.In(loc).Sub(Member.Schedul)).LimitFirstN(2)
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetTitle(VTData.VTName).
										SetDescription(Member.Title).
										SetImage(Member.Thumb).
										SetThumbnail(VTData.YtData.YoutubeAvatar).
										SetURL("https://www.youtube.com/watch?v="+Member.VideoID).
										AddField("Live duration", durationlive.String()).
										AddField("Live ended", duration.String()+" Ago").
										InlineAllFields().
										AddField("Viewers", Member.Viewers).
										SetFooter(Member.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
										SetColor(Color).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								SendEmbed(Memberst{
									VTName: VTuberGroup.GroupName,
									Msg:    "Internal error XD",
								})
							}
						}
					} else {
						var (
							Data []database.YtDbData
						)

						if CheckReg(VTuberGroup.GroupName, Region) {
							Data, err = database.YtGetStatus(VTuberGroup.ID, 0, "past", Region)
							if err != nil {
								log.Error(err)
							}
						} else {
							Data, err = database.YtGetStatus(VTuberGroup.ID, 0, "past", "")
							if err != nil {
								log.Error(err)
							}
						}

						if Data != nil {
							for _, Member := range Data {
								Color, err := engine.GetColor(config.TmpDir, Member.Thumb)
								if err != nil {
									log.Error(err)
								}
								Views, err := strconv.Atoi(Member.Viewers)
								if err != nil {
									log.Error(err)
								}
								Member.Viewers = engine.NearestThousandFormat(float64(Views))

								loc := engine.Zawarudo(Member.Region)
								duration := durafmt.Parse(time.Now().In(loc).Sub(Member.Schedul.In(loc))).LimitFirstN(2)
								durationlive := durafmt.Parse(Member.End.In(loc).Sub(Member.Schedul)).LimitFirstN(2)
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetTitle(engine.FixName(Member.NameEN, Member.NameJP)).
									SetThumbnail(Member.YoutubeAvatar).
									SetDescription(Member.Title).
									SetImage(Member.Thumb).
									SetURL("https://www.youtube.com/watch?v="+Member.VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Member.Viewers).
									SetFooter(Member.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.GroupName,
								Msg:    "Internal error XD",
							})
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Prefix+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
