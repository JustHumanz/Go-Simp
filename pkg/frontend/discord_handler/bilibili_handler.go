package discordhandler

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

//BiliBiliMessage message handler
func BiliBiliMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := config.PBilibili
	loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 1 {
			Payload := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
			if CommandArray[0] == Prefix+Live {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						VTData := ValidName(FindGroupArry)
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
						} else {
							DataMembers := database.BilGet(0, VTData.ID, "Live")
							if len(DataMembers) > 0 {
								for _, DataMember := range DataMembers {
									diff := time.Now().In(loc).Sub(DataMember.ScheduledStart.In(loc))
									Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(VTData.VTName).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription(DataMember.Description).
										SetThumbnail(DataMember.Avatar).
										SetImage(DataMember.Thumbnail).
										SetURL("https://live.bilibili.com/"+strconv.Itoa(DataMember.LiveRoomID)).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Online", strconv.Itoa(DataMember.Online)).
										SetColor(Color).
										SetFooter(DataMember.ScheduledStart.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetDescription("It looks like `"+VTData.VTName+"` doesn't have a livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						DataGroups := database.BilGet(VTuberGroup.ID, 0, "Live")
						if len(DataGroups) > 0 {
							for _, DataGroup := range DataGroups {
								diff := time.Now().In(loc).Sub(DataGroup.ScheduledStart.In(loc))
								Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle(engine.FixName(DataGroup.EnName, DataGroup.JpName)).
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription(DataGroup.Description).
									SetThumbnail(DataGroup.Avatar).
									SetImage(DataGroup.Thumbnail).
									SetURL("https://live.bilibili.com/"+strconv.Itoa(DataGroup.LiveRoomID)).
									AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Online", strconv.Itoa(DataGroup.Online)).
									SetColor(Color).
									SetFooter(DataGroup.ScheduledStart.In(loc).Format(time.RFC822)).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
								SetDescription("It looks like `"+VTuberGroup.NameGroup+"` doesn't have a livestream right now").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}

			} else if CommandArray[0] == Prefix+Past || CommandArray[0] == Prefix+"last" {
				for _, FindGroupArry := range Payload {
					VTuberGroup, err := FindGropName(FindGroupArry)
					if err != nil {
						VTData := ValidName(FindGroupArry)
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
						} else {
							DataMembers := database.BilGet(0, VTData.ID, "Past")
							if len(DataMembers) > 0 {
								for _, DataMember := range DataMembers {
									diff := DataMember.ScheduledStart.In(loc).Sub(time.Now().In(loc))
									Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
									if err != nil {
										log.Error(err)
									}
									_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetTitle(engine.FixName(DataMember.EnName, DataMember.JpName)).
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription(DataMember.Description).
										SetThumbnail(DataMember.Avatar).
										SetImage(DataMember.Thumbnail).
										SetURL("https://live.bilibili.com/"+strconv.Itoa(DataMember.LiveRoomID)).
										AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
										AddField("Online", strconv.Itoa(DataMember.Online)).
										SetColor(Color).
										SetFooter(DataMember.ScheduledStart.In(loc).Format(time.RFC822)).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetDescription("It looks like `"+VTData.VTName+"` doesn't have a Past livestream right now").
									SetImage(config.WorryIMG).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					} else {
						for _, LiveBili := range database.BilGet(VTuberGroup.ID, 0, "Past") {
							diff := time.Now().In(loc).Sub(LiveBili.ScheduledStart.In(loc))
							Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetTitle(engine.FixName(LiveBili.EnName, LiveBili.JpName)).
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription(LiveBili.Description).
								SetThumbnail(LiveBili.Avatar).
								SetImage(LiveBili.Thumbnail).
								SetURL("https://live.bilibili.com/"+strconv.Itoa(LiveBili.LiveRoomID)).
								AddField("Start live", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Online", strconv.Itoa(LiveBili.Online)).
								SetColor(Color).
								SetFooter(LiveBili.ScheduledStart.In(loc).Format(time.RFC822)).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete command")
			if err != nil {
				log.Error(err)
			}
		}
	}
}

//BiliBiliSpace message hadler
func BiliBiliSpace(s *discordgo.Session, m *discordgo.MessageCreate) {
	loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
	m.Content = strings.ToLower(m.Content)
	Prefix := "sp_" + config.PBilibili
	if strings.HasPrefix(m.Content, Prefix) {
		Payload := m.Content[len(Prefix):]
		if Payload != "" {
			for _, FindGroupArry := range strings.Split(strings.TrimSpace(Payload), ",") {
				VTuberGroup, err := FindGropName(FindGroupArry)
				if err != nil {
					VTData := ValidName(FindGroupArry)
					if VTData.ID == 0 {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry+"`,Name of Vtuber Group or Vtuber Name was not found")
					} else {
						DataMembers := database.SpaceGet(0, VTData.ID)
						if len(DataMembers) > 0 {
							for _, DataMember := range DataMembers {
								diff := time.Now().In(loc).Sub(DataMember.UploadDate.In(loc))
								Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetTitle(engine.FixName(DataMember.EnName, DataMember.JpName)).
									SetDescription(DataMember.Title).
									SetImage(DataMember.Thumbnail).
									SetThumbnail(DataMember.Avatar).
									SetURL("https://www.bilibili.com/video/"+DataMember.VideoID).
									AddField("Type", DataMember.Type).
									AddField("Video uploaded", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
									AddField("Duration", DataMember.Length).
									AddField("Viewers now", strconv.Itoa(DataMember.Viewers)).
									SetFooter(DataMember.UploadDate.In(loc).Format(time.RFC822), config.BiliBiliIMG).
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
								SetDescription("It looks like `"+VTData.VTName+"` doesn't have a video in space.bilibili").
								SetImage(config.WorryIMG).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
					}
					break
				} else {
					DataMembers := database.SpaceGet(VTuberGroup.ID, 0)
					if len(DataMembers) > 0 {
						for _, DataMember := range DataMembers {
							diff := time.Now().In(loc).Sub(DataMember.UploadDate.In(loc))
							Color, err := engine.GetColor("/tmp/bil1.tmp", m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}

							_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
								SetTitle(engine.FixName(DataMember.EnName, DataMember.JpName)).
								SetDescription(DataMember.Title).
								SetImage(DataMember.Thumbnail).
								SetThumbnail(DataMember.Avatar).
								SetURL("https://www.bilibili.com/video/"+DataMember.VideoID).
								AddField("Type", DataMember.Type).
								AddField("Video uploaded", durafmt.Parse(diff).LimitFirstN(2).String()+" Ago").
								AddField("Duration", DataMember.Length).
								AddField("Viewers now", strconv.Itoa(DataMember.Viewers)).
								SetFooter(DataMember.UploadDate.In(loc).Format(time.RFC822), config.BiliBiliIMG).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}

						}
					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
							SetDescription("It looks like `"+VTuberGroup.NameGroup+"` doesn't have a video in space.bilibili").
							SetImage(config.WorryIMG).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
			}
		}
	}
}
