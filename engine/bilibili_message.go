package engine

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func BiliBiliSpace(s *discordgo.Session, m *discordgo.MessageCreate) {
	loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
	Prefix := config.PBilibili
	m.Content = strings.ToLower(m.Content)
	SendEmbed := func(Data Memberst) {
		if Data.VideoID != "" {
			Color, err := GetColor("/tmp/bil1.tmp", Data.Thumb)
			if err != nil {
				log.Error(err)
			}
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
				SetTitle(Data.VTName).
				SetDescription(Data.Desc).
				SetImage(Data.Thumb).
				SetThumbnail(Data.BiliAvatar).
				SetURL(Data.VideoID).
				AddField(Data.Msg, Data.Msg1).
				AddField("Type", Data.Msg3).
				AddField("Duration", Data.Length).
				SetFooter(Data.Msg2, config.BiliBiliIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed)
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
				SetTitle(Data.VTName).
				SetDescription(Data.Msg).
				SetImage(config.WorryIMG).MessageEmbed)
		}
	}

	Prefix2 := "sp_" + Prefix
	if strings.HasPrefix(m.Content, Prefix2) {
		Payload := m.Content[len(Prefix2):]
		if Payload != "" {
			GroupName := strings.TrimSpace(Payload)
			FindGroupArry := strings.Split(GroupName, ",")
			for i := 0; i < len(FindGroupArry); i++ {
				VTuberGroup, err := FindGropName(FindGroupArry[i])
				if err != nil {
					VTData := ValidName(FindGroupArry[i])
					if VTData.ID == 0 {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
					} else {
						DataMember := database.SpaceGet(0, VTData.ID) //database.BilGet(0, VTData.ID, "Upcoming")
						if DataMember != nil {
							for z := 0; z < len(DataMember); z++ {
								expiresAt := time.Now().In(loc)
								diff := expiresAt.In(loc).Sub(DataMember[z].UploadDate)
								duration := durafmt.Parse(diff).LimitFirstN(2)

								VTData.Desc = DataMember[z].Title
								VTData.VideoID = "https://www.bilibili.com/video/" + DataMember[z].VideoID
								VTData.Msg = "Video upload "
								VTData.Msg1 = duration.String() + " Ago"
								VTData.Msg2 = "Viewers now " + strconv.Itoa(DataMember[z].Viewers)
								VTData.Msg3 = DataMember[z].Type
								VTData.Length = DataMember[z].Length
								SendEmbed(VTData)
							}
						} else {
							SendEmbed(Memberst{
								Msg: "It looks like `" + VTData.VTName + "` doesn't have a schedule livestream for now",
							})
						}
					}
					break
				} else {
					Data := database.SpaceGet(VTuberGroup.ID, 0) //database.BilGet(VTuberGroup.ID, 0, "Upcoming")
					if Data != nil {
						for ii := 0; ii < len(Data); ii++ {
							diff := time.Now().In(loc).Sub(Data[ii].UploadDate)
							SendEmbed(Memberst{
								VTName:     FixName(Data[ii].EnName, Data[ii].JpName),
								BiliAvatar: Data[ii].Avatar,
								Desc:       Data[ii].Title,
								Thumb:      Data[ii].Thumbnail,
								VideoID:    "https://www.bilibili.com/video/" + Data[ii].VideoID,
								Msg:        "Video upload ",
								Msg1:       durafmt.Parse(diff).LimitFirstN(2).String() + " Ago",
								Msg2:       "Viewers now " + strconv.Itoa(Data[ii].Viewers),
								Msg3:       Data[ii].Type,
								Length:     Data[ii].Length,
							})
						}
					} else {
						SendEmbed(Memberst{
							Msg: "It looks like `" + VTuberGroup.NameGroup + "` doesn't have a schedule livestream for now",
						})
					}
				}
			}
		}
	}

}

func BiliBiliMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := config.PBilibili
	loc, _ := time.LoadLocation("Asia/Shanghai") /*Use CST*/
	m.Content = strings.ToLower(m.Content)
	SendEmbed := func(Data Memberst) {
		if Data.VideoID != "" {
			Color, err := GetColor("/tmp/bil1.tmp", Data.Thumb)
			if err != nil {
				log.Error(err)
			}
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
				SetTitle(Data.VTName).
				SetDescription(Data.Desc).
				SetImage(Data.Thumb).
				SetThumbnail(Data.BiliAvatar).
				SetURL(Data.VideoID).
				AddField(Data.Msg, Data.Msg1).
				SetFooter(Data.Msg2).
				SetColor(Color).MessageEmbed)
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
				SetTitle(Data.VTName).
				SetDescription(Data.Msg).
				SetImage(config.WorryIMG).MessageEmbed)
		}
	}

	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 1 {
			CommandArray[0] = strings.ToLower(CommandArray[0])
			GroupName := strings.TrimSpace(CommandArray[1])
			FindGroupArry := strings.Split(GroupName, ",")
			if CommandArray[0] == Prefix+"live" {
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
						} else {
							DataMember := database.BilGet(0, VTData.ID, "Live")
							if DataMember != nil {
								for z := 0; z < len(DataMember); z++ {
									diff := time.Now().In(loc).Sub(DataMember[i].ScheduledStart.In(loc))
									VTData.Desc = DataMember[z].Description
									VTData.VideoID = "https://live.bilibili.com/" + strconv.Itoa(DataMember[z].LiveRoomID)
									VTData.Msg = "Start live"
									VTData.Msg1 = durafmt.Parse(diff).LimitFirstN(2).String() + " Ago"
									VTData.Msg2 = "Online : " + strconv.Itoa(DataMember[z].Online)
									SendEmbed(VTData)
								}
							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream right now",
								})
							}
						}
					} else {
						Data := database.BilGet(VTuberGroup.ID, 0, "Live")
						if Data != nil {
							for ii := 0; ii < len(Data); ii++ {
								diff := time.Now().In(loc).Sub(Data[ii].ScheduledStart.In(loc))
								SendEmbed(Memberst{
									VTName:     FixName(Data[ii].EnName, Data[ii].JpName),
									BiliAvatar: Data[ii].Avatar,
									Desc:       Data[ii].Description,
									Thumb:      Data[ii].Thumbnail,
									VideoID:    "https://live.bilibili.com/" + strconv.Itoa(Data[ii].LiveRoomID),
									Msg:        "Start live",
									Msg1:       durafmt.Parse(diff).LimitFirstN(2).String() + " Ago",
									Msg2:       "Online : " + strconv.Itoa(Data[ii].Online),
								})
							}
						} else {
							SendEmbed(Memberst{
								Msg: "It looks like `" + VTuberGroup.NameGroup + "` doesn't have a livestream right now",
							})
						}
					}
				}
			} else if CommandArray[0] == Prefix+"past" || CommandArray[0] == Prefix+"last" {
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not found")
						} else {
							DataMember := database.BilGet(0, VTData.ID, "Past")
							if DataMember != nil {
								for z := 0; z < len(DataMember); z++ {
									diff := DataMember[i].ScheduledStart.In(loc).Sub(time.Now().In(loc))
									VTData.Desc = DataMember[z].Description
									VTData.VideoID = "https://live.bilibili.com/" + strconv.Itoa(DataMember[z].LiveRoomID)
									VTData.Msg = "Start live"
									VTData.Msg1 = durafmt.Parse(diff).LimitFirstN(2).String() + " Ago"
									VTData.Msg2 = "Online : " + strconv.Itoa(DataMember[z].Online)
									SendEmbed(VTData)
								}
							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a Past livestream right now",
								})
							}
						}
					} else {
						Data := database.BilGet(VTuberGroup.ID, 0, "Past")
						if Data != nil {
							for ii := 0; ii < len(Data); ii++ {
								diff := time.Now().In(loc).Sub(Data[ii].ScheduledStart.In(loc))
								SendEmbed(Memberst{
									VTName:     FixName(Data[ii].EnName, Data[ii].JpName),
									BiliAvatar: Data[ii].Avatar,
									Desc:       Data[ii].Description,
									Thumb:      Data[ii].Thumbnail,
									VideoID:    "https://live.bilibili.com/" + strconv.Itoa(Data[ii].LiveRoomID),
									Msg:        "Start live ",
									Msg1:       durafmt.Parse(diff).LimitFirstN(2).String() + " Ago",
									Msg2:       "Online : " + strconv.Itoa(Data[ii].Online),
								})
							}
						} else {
							SendEmbed(Memberst{
								Msg: "It looks like `" + VTuberGroup.NameGroup + "` doesn't have a Past livestream right now",
							})
						}
					}
				}
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Incomplete command")
		}
	}
}
