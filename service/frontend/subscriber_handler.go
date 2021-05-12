package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/prediction"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func gacha() bool {
	return rand.Float32() < 0.5
}

//SubsMessage subscriber message handler
func SubsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := configfile.BotPrefix.General
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 0 {
			if CommandArray[0] == Prefix+Info {
				Member := FindVtuber(CommandArray[1])
				var (
					Avatar string
				)
				SubsData, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}

				if gacha() {
					Avatar = Member.YoutubeAvatar
				} else {
					if Member.BiliRoomID != 0 {
						Avatar = Member.BiliBiliAvatar
					} else {
						Avatar = Member.YoutubeAvatar
					}
				}
				Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
				if err != nil {
					log.Error(err)
				}
				YTSubs := "[Youtube](https://www.youtube.com/channel/" + Member.YoutubeID + "?sub_confirmation=1)"
				BiliFollow := "[BiliBili](https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID) + ")"
				TwitterFollow := "[Twitter](https://twitter.com/" + Member.TwitterName + ")"
				if SubsData.BiliFollow != 0 {
					msg, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetTitle(engine.FixName(Member.EnName, Member.JpName)).
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
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err, msg)
					}
				} else {
					msg, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetTitle(engine.FixName(Member.EnName, Member.JpName)).
						SetImage(Avatar).
						AddField("Youtube subscribers", engine.NearestThousandFormat(float64(SubsData.YtSubs))).
						AddField("Youtube viewers", engine.NearestThousandFormat(float64(SubsData.YtViews))).
						AddField("Youtube videos", engine.NearestThousandFormat(float64(SubsData.YtVideos))).
						InlineAllFields().
						AddField("Twitter followers", engine.NearestThousandFormat(float64(SubsData.TwFollow))).
						RemoveInline().
						AddField("<:yt:796023828723269662>", YTSubs).
						AddField("<:bili:796025336542265344>", BiliFollow).
						AddField("<:tw:796025611210588187>", TwitterFollow).
						InlineAllFields().
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err, msg)
					}
				}
			} else if CommandArray[0] == Prefix+Predick {
				if len(CommandArray) > 2 {
					Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
					if err != nil {
						log.Error(err)
					}

					State := CommandArray[1]
					Member := FindVtuber(CommandArray[2])
					Subs, err := Member.GetSubsCount()
					if err != nil {
						log.Error(err)
					}
					var (
						msg    *prediction.Message
						tmp    string
						Avatar string
						Url    string
					)

					if State == "-twitter" || State == "-tw" {
						if Member.IsYtNill() {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetDescription(engine.FixName(Member.EnName, Member.JpName)+" don't have Twitter account").
								SetImage(engine.NotFoundIMG()).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
						msg = &prediction.Message{
							State: "Twitter",
							Name:  Member.Name,
							Limit: 7,
						}
						tmp = strconv.Itoa(Subs.TwFollow)
						Avatar = Member.YoutubeAvatar
						Url = "https://twitter.com/" + Member.TwitterName
					} else if State == "-youtube" || State == "-yt" {
						if Member.IsYtNill() {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetDescription(engine.FixName(Member.EnName, Member.JpName)+" don't have Youtube channel").
								SetImage(engine.NotFoundIMG()).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}

						msg = &prediction.Message{
							State: "Youtube",
							Name:  Member.Name,
							Limit: 7,
						}
						tmp = strconv.Itoa(Subs.YtSubs)
						Avatar = Member.YoutubeAvatar
						Url = "https://www.youtube.com/channel/" + Member.YoutubeID + "?sub_confirmation=1"
					} else if State == "-bilibili" || State == "-bl" {
						if Member.IsBiliNill() {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetDescription(engine.FixName(Member.EnName, Member.JpName)+" don't have bilibili channel").
								SetImage(engine.NotFoundIMG()).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
						msg = &prediction.Message{
							State: "BiliBili",
							Name:  Member.Name,
							Limit: 7,
						}
						tmp = strconv.Itoa(Subs.BiliFollow)
						Avatar = Member.BiliBiliAvatar
						Url = "https://space.bilibili/" + strconv.Itoa(Member.BiliBiliID)
					}

					RawData, err := PredictionConn.GetSubscriberPrediction(context.Background(), msg)
					if err != nil {
						log.Error(err)
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetTitle(engine.FixName(Member.EnName, Member.JpName)).
							SetDescription("Something error\n"+err.Error()).
							SetImage(engine.NotFoundIMG()).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
					if RawData.Code == 0 {
						target := time.Now().AddDate(0, 0, int(msg.Limit))
						dateFormat := fmt.Sprintf("%s/%d", target.Month().String(), target.Day())
						now := time.Now()
						nowFormat := fmt.Sprintf("%s/%d", now.Month().String(), now.Day())
						Graph := "[Views Graph](https://prometheus.humanz.moe/graph?g0.expr=get_subscriber%7Bstate%3D%22" + msg.State + "%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=1w)"
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetTitle(engine.FixName(Member.EnName, Member.JpName)).
							SetURL(Url).
							AddField("Current "+msg.State+" followes/subscriber("+nowFormat+")", tmp).
							AddField("Next 7 days Prediction("+dateFormat+")", RawData.Prediction).
							RemoveInline().
							AddField("Prediction score", RawData.Score).
							AddField("Graph", Graph).
							InlineAllFields().
							SetThumbnail(Avatar).
							SetColor(Color).
							SetFooter("algorithm : Linear regression").MessageEmbed)
						if err != nil {
							log.Error(err)
						}

					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetTitle(engine.FixName(Member.EnName, Member.JpName)).
							SetDescription("Something error\ncan't make prediction,Bot still need more subscriber/followers data").
							SetImage(engine.NotFoundIMG()).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					}
				} else {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetDescription("Incomplete `"+Prefix+Predick+"` command\nUse `"+Prefix+Predick+" [State] [Vtuber nickname]`\nExample `"+Prefix+Predick+" -tw Parerun`").
						SetImage(engine.NotFoundIMG()).MessageEmbed)
					if err != nil {
						log.Error(err)
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
