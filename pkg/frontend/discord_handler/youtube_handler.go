package discordhandler

import (
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func YoutubeMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	var (
		Region    string
		Prefix    = config.PYoutube
		SendEmbed = func(Data Memberst) {

			Color, err := engine.GetColor("/tmp/yt2.tmp", Data.Thumb)
			if err != nil {
				log.Error(err)
			}

			if Data.VideoID != "" {
				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Data.VTName).
					SetDescription(Data.Title).
					SetImage(Data.Thumb).
					SetThumbnail(Data.YtAvatar).
					SetURL("https://www.youtube.com/watch?v="+Data.VideoID).
					AddField(Data.Msg, Data.Msg1).
					AddField("Viewers", Data.Msg3).
					InlineAllFields().
					AddField("Type", engine.YtFindType(Data.Title)).
					SetFooter(Data.Msg2, config.YoutubeIMG).
					SetColor(Color).MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Data.VTName).
					SetDescription(Data.Msg).
					SetImage(config.WorryIMG).MessageEmbed)
			}
		}
	)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 3 && CommandArray[2] == "-region" {
			Region = strings.ToLower(CommandArray[3])
			if len(Region) > 2 {
				s.ChannelMessageSend(m.ChannelID, "Only support 1 Region,ignoring `"+Region[2:]+"`")
			}
		} else {
			Region = ""
		}
		if strings.ToLower(CommandArray[0]) == Prefix+Upcoming || strings.ToLower(CommandArray[0]) == Prefix+"up" {
			if len(CommandArray) > 1 {
				for _, GroupNameQuery := range strings.Split(strings.TrimSpace(CommandArray[1]), ",") {
					VTuberGroup, err := FindGropName(GroupNameQuery)
					if err != nil {
						VTData := ValidName(GroupNameQuery)
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+GroupNameQuery+"`,Name of Vtuber Group or Vtuber Name was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "upcoming", Region)
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := engine.Zawarudo(Member.Region)
									duration := durafmt.Parse(Member.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
									SendEmbed(Memberst{
										VTName:   VTData.VTName,
										Title:    Member.Title,
										Thumb:    Member.Thumb,
										VideoID:  Member.VideoID,
										YtAvatar: Member.YoutubeAvatar,
										Msg:      "Start live in",
										Msg1:     duration.String(),
										Msg2:     Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:     Member.Viewers + " simps waiting in Room Chat",
									})
								}

							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream schedule for now",
								})
							}
						}
					} else {

						GroupData := database.YtGetStatus(VTuberGroup.ID, 0, "upcoming", Region)
						if GroupData != nil {
							for _, Data := range GroupData {
								loc := engine.Zawarudo(Data.Region)
								duration := durafmt.Parse(Data.Schedul.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2)
								SendEmbed(Memberst{
									VTName:    engine.FixName(Data.NameEN, Data.NameJP),
									Title:     Data.Title,
									Thumb:     Data.Thumb,
									VideoID:   Data.VideoID,
									YtAvatar:  Data.YoutubeAvatar,
									Msg:       "Start in",
									Msg1:      duration.String(),
									Msg2:      Data.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Data.Viewers + " simps waiting in Room Chat",
									YtChannel: Data.ChannelID,
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.NameGroup,
								Msg:    "It seems like `" + strings.Title(GroupNameQuery+" "+Region) + "` doesn't have livestream schedule for now",
							})
						}
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete Upcoming command")
			}
		} else if strings.ToLower(CommandArray[0]) == Prefix+Live {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						VTData := ValidName(FindGroupArry[i])
						if VTData.ID == 0 {
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group or Vtuber Name was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "live", Region)
							if DataMember != nil {
								for _, Member := range DataMember {
									loc := engine.Zawarudo(Member.Region)
									expiresAt := time.Now().In(loc)
									diff := expiresAt.In(loc).Sub(Member.Schedul)
									duration := durafmt.Parse(diff).LimitFirstN(2)

									SendEmbed(Memberst{
										VTName:   VTData.VTName,
										Title:    Member.Title,
										Thumb:    Member.Thumb,
										VideoID:  Member.VideoID,
										YtAvatar: Member.YoutubeAvatar,
										Msg:      "Start live in",
										Msg1:     duration.String() + " Ago",
										Msg2:     Member.Schedul.In(loc).Format(time.RFC822),
										Msg3:     Member.Viewers + " Simps",
									})
								}
							} else {
								SendEmbed(Memberst{
									Msg: "It looks like `" + VTData.VTName + "` doesn't have a livestream right now",
								})
							}
						}
					} else {
						GroupData := database.YtGetStatus(VTuberGroup.ID, 0, "live", Region)
						if GroupData != nil {
							for _, Data := range GroupData {
								loc := engine.Zawarudo(Data.Region)
								diff := time.Now().In(loc).Sub(Data.Schedul.In(loc))
								duration := durafmt.Parse(diff).LimitFirstN(2)
								SendEmbed(Memberst{
									VTName:    engine.FixName(Data.NameEN, Data.NameJP),
									Title:     Data.Title,
									Thumb:     Data.Thumb,
									VideoID:   Data.VideoID,
									YtAvatar:  Data.YoutubeAvatar,
									Msg:       "Start in",
									Msg1:      duration.String() + " Ago",
									Msg2:      Data.Schedul.In(loc).Format(time.RFC822),
									Msg3:      Data.Viewers + " Simps",
									YtChannel: Data.ChannelID,
								})
							}
						} else {
							SendEmbed(Memberst{
								VTName: VTuberGroup.NameGroup,
								Msg:    "It looks like `" + strings.Title(FindGroupArry[i]+" "+Region) + "` doesn't have a livestream right now",
							})
						}
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete Live command")
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
							s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group or Vtuber Name was not found")
						} else {
							DataMember := database.YtGetStatus(0, VTData.ID, "past", Region)
							if DataMember != nil {
								for z := 0; z < len(DataMember); z++ {
									Color, err := engine.GetColor("/tmp/yt3.tmp", DataMember[z].Thumb)
									if err != nil {
										log.Error(err)
									}
									loc := engine.Zawarudo(DataMember[z].Region)
									diff := time.Now().In(loc).Sub(DataMember[z].Schedul)
									duration := durafmt.Parse(diff).LimitFirstN(2)
									diff2 := DataMember[z].End.In(loc).Sub(DataMember[z].Schedul)
									durationlive := durafmt.Parse(diff2).LimitFirstN(2)
									s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
										SetTitle(VTData.VTName).
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
						var (
							Data []database.YtDbData
						)

						if CheckReg(VTuberGroup.NameGroup, Region) {
							Data = database.YtGetStatus(VTuberGroup.ID, 0, "past", Region)
						} else {
							Data = database.YtGetStatus(VTuberGroup.ID, 0, "past", "")
						}

						if Data != nil {
							for ii := 0; ii < len(Data); ii++ {
								Color, err := engine.GetColor("/tmp/yt3.tmp", Data[ii].Thumb)
								if err != nil {
									log.Error(err)
								}

								loc := engine.Zawarudo(Data[ii].Region)
								duration := durafmt.Parse(time.Now().In(loc).Sub(Data[ii].Schedul.In(loc))).LimitFirstN(2)
								durationlive := durafmt.Parse(Data[ii].End.In(loc).Sub(Data[ii].Schedul)).LimitFirstN(2)
								s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetTitle(engine.FixName(Data[ii].NameEN, Data[ii].NameJP)).
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
