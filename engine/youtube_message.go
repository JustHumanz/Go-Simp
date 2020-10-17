package engine

import (
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func YoutubeMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PYoutube
	SendEmbed := func(Data Memberst) {
		Color, err := GetColor("/tmp/yt2.tmp", m.Author.AvatarURL("80"))
		if err != nil {
			log.Error(err)
		}
		if Data.Thumb != "" {
			Color, err = GetColor("/tmp/yt2.tmp", Data.Thumb)
			if err != nil {
				log.Error(err)
			}
		}
		if Data.VideoID != "" {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle(Data.VTName).
				SetDescription(Data.Desc).
				SetImage(Data.Thumb).
				SetThumbnail(Data.YtAvatar).
				SetURL("https://www.youtube.com/watch?v="+Data.VideoID).
				AddField(Data.Msg, Data.Msg1).
				AddField("Viewers", Data.Msg3).
				SetFooter(Data.Msg2, config.YoutubeIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed)
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle(Data.VTName).
				SetDescription(Data.Msg).
				SetImage(config.WorryIMG).MessageEmbed)
		}
	}

	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if strings.ToLower(CommandArray[0]) == Prefix+"upcoming" || strings.ToLower(CommandArray[0]) == Prefix+"up" {
			if len(CommandArray) > 1 {
				GroupName := strings.TrimSpace(CommandArray[1])
				FindGroupArry := strings.Split(GroupName, ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "upcoming")
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := Zawarudo(Member.Region)
									duration := durafmt.Parse(Member.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
									SendEmbed(Memberst{
										VTName:  "【" + VTData.VTName + "】",
										Desc:    Member.Title,
										Thumb:   Member.Thumb,
										VideoID: Member.VideoID,
										Msg:     "Start live in",
										Msg1:    duration.String(),
										Msg2:    Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:    Member.Viewers + " simps waiting in Room Chat",
									})
								}

							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream schedule for now",
								})
							}
						}
					} else {
						GroupData := database.YtGetStatus(VTuberGroup.ID, 0, "upcoming")
						if GroupData != nil {
							for _, Data := range GroupData {
								loc := Zawarudo(Data.Region)
								duration := durafmt.Parse(Data.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
								SendEmbed(Memberst{
									VTName:    "【" + FixName(Data.NameEN, Data.NameJP) + "】",
									Desc:      Data.Title,
									Thumb:     Data.Thumb,
									VideoID:   Data.VideoID,
									Msg:       "Start in",
									Msg1:      duration.String(),
									Msg2:      Data.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Data.Viewers + " simps waiting in Room Chat",
									YtChannel: strings.Split(Data.ChannelID, "\n"),
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.NameGroup,
								Msg:    "It seems like `" + strings.Title(FindGroupArry[i]) + "` doesn't have livestream schedule for now",
							})
						}
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete Upcoming command")
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+"live" {
			if len(CommandArray) > 1 {
				GroupName := strings.TrimSpace(CommandArray[1])
				FindGroupArry := strings.Split(GroupName, ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "live")
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := Zawarudo(Member.Region)
									expiresAt := time.Now().In(loc)
									diff := expiresAt.In(loc).Sub(Member.Schedul)
									duration := durafmt.Parse(diff).LimitFirstN(2)

									SendEmbed(Memberst{
										VTName:  "【" + VTData.VTName + "】",
										Desc:    Member.Title,
										Thumb:   Member.Thumb,
										VideoID: Member.VideoID,
										Msg:     "Start live in",
										Msg1:    duration.String() + " Ago",
										Msg2:    Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:    Member.Viewers,
									})
								}
							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream right now",
								})
							}
						}
					} else {
						DataGroup := database.YtGetStatus(VTuberGroup.ID, 0, "live")
						if DataGroup != nil {
							for _, Data := range DataGroup {
								loc := Zawarudo(Data.Region)
								diff := time.Now().In(loc).Sub(Data.Schedul.In(loc))
								duration := durafmt.Parse(diff).LimitFirstN(2)
								SendEmbed(Memberst{
									VTName:    "【" + FixName(Data.NameEN, Data.NameJP) + "】",
									Desc:      Data.Title,
									Thumb:     Data.Thumb,
									VideoID:   Data.VideoID,
									Msg:       "Start in",
									Msg1:      duration.String() + " Ago",
									Msg2:      Data.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Data.Viewers,
									YtChannel: strings.Split(Data.ChannelID, "\n"),
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.NameGroup,
								Msg:    "It looks like " + strings.Title(FindGroupArry[i]) + " doesn't have a livestream right now",
							})
						}
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete Live command")
				return
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+"last" || strings.ToLower(CommandArray[0]) == Prefix+"past" {
			if len(CommandArray) > 1 {
				GroupName := strings.TrimSpace(CommandArray[1])
				FindGroupArry := strings.Split(GroupName, ",")

				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "past")
							if DataMember != nil {
								for z := 0; z < len(DataMember); z++ {
									Color, err := GetColor("/tmp/yt3.tmp", DataMember[z].Thumb)
									if err != nil {
										log.Error(err)
									}
									loc := Zawarudo(DataMember[z].Region)
									diff := time.Now().In(loc).Sub(DataMember[z].Schedul)
									duration := durafmt.Parse(diff).LimitFirstN(2)
									diff2 := DataMember[z].End.In(loc).Sub(DataMember[z].Schedul)
									durationlive := durafmt.Parse(diff2).LimitFirstN(2)
									s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
										SetTitle("【"+VTData.VTName+"】").
										SetDescription(DataMember[z].Title).
										SetImage(DataMember[z].Thumb).
										SetThumbnail(VTData.YtAvatar).
										SetURL("https://www.youtube.com/watch?v="+DataMember[z].VideoID).
										AddField("Live duration", durationlive.String()).
										AddField("Live ended", duration.String()+" Ago").
										InlineAllFields().
										AddField("Viewers", DataMember[z].Viewers).
										SetFooter(DataMember[z].Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
										SetColor(Color).MessageEmbed)
								}
							} else {
								SendEmbed(Memberst{
									VTName: VTuberGroup.NameGroup,
									Msg:    "Internal error XD",
								})
							}
						}
					} else {
						Data := database.YtGetStatus(VTuberGroup.ID, 0, "last")
						if Data != nil {
							for ii := 0; ii < len(Data); ii++ {
								Color, err := GetColor("/tmp/yt3.tmp", Data[ii].Thumb)
								if err != nil {
									log.Error(err)
								}

								loc := Zawarudo(Data[ii].Region)
								diff := time.Now().In(loc).Sub(Data[ii].Schedul.In(loc))
								duration := durafmt.Parse(diff).LimitFirstN(2)

								diff2 := Data[ii].End.In(loc).Sub(Data[ii].Schedul)
								durationlive := durafmt.Parse(diff2).LimitFirstN(2)
								s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetTitle(FixName(Data[ii].NameEN, Data[ii].NameJP)).
									SetThumbnail(Data[ii].YoutubeAvatar).
									SetDescription(Data[ii].Title).
									SetImage(Data[ii].Thumb).
									SetURL("https://www.youtube.com/watch?v="+Data[ii].VideoID).
									AddField("Live duration", durationlive.String()).
									AddField("Live ended", duration.String()+" Ago").
									InlineAllFields().
									AddField("Viewers", Data[ii].Viewers).
									SetFooter(Data[ii].Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
									SetColor(Color).MessageEmbed)
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.NameGroup,
								Msg:    "Internal error XD",
							})
						}
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete Last command")
			}
		}
	}
}

func Zawarudo(Region string) *time.Location {
	if Region == "ID" {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		return loc
	} else if Region == "JP" {
		loc, _ := time.LoadLocation("Asia/Tokyo")
		return loc
	} else if Region == "CN" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		return loc
	} else if Region == "KR" {
		loc, _ := time.LoadLocation("Asia/Seoul")
		return loc
	} else {
		loc, _ := time.LoadLocation("UTC")
		return loc
	}
}
