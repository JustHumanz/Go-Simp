package main

import (
	"context"
	"math/rand"
	"strconv"
	"strings"

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
				Member := FindVtuber(CommandArray[1])
				Subs, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}
				RawData, err := PredictionConn.GetSubscriberPrediction(context.Background(), &prediction.Message{
					State: "Twitter",
					Name:  Member.Name,
					Limit: 7,
				})
				if err != nil {
					log.Error(err)
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetDescription("Something error\n"+err.Error()).
						SetImage(engine.NotFoundIMG()).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}
				if RawData.Code == 0 {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetDescription("Current twitter followes "+strconv.Itoa(Subs.TwFollow)+" next 7 days Prediction "+RawData.Prediction+" with Prediction score "+RawData.Score).
						SetImage(Member.TwitchAvatar).MessageEmbed)
					if err != nil {
						log.Error(err)
					}

				} else {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetDescription("Something error\ncan't make prediction,Bot still need more subscriber/followers data").
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
