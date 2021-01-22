package main

import (
	"math/rand"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func gacha() bool {
	return rand.Float32() < 0.5
}

//SubsMessage subscriber message handler
func SubsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := config.BotConf.BotPrefix.General
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if len(CommandArray) > 0 {
			if CommandArray[0] == Prefix+Info {
				for _, Group := range engine.GroupData {
					for _, Mem := range strings.Split(CommandArray[1], ",") {
						for _, Member := range database.GetMembers(Group.ID) {
							if Mem == strings.ToLower(Member.Name) {
								var (
									Avatar    string
									GroupIcon string
								)
								SubsData, err := Member.GetSubsCount()
								if err != nil {
									log.Error(err)
								}
								if Group.GroupName != "Independen" {
									GroupIcon = Group.IconURL
								}
								if gacha() {
									Avatar = Member.YoutubeAvatar
								} else {
									if Member.BiliRoomID != 0 {
										Avatar = Member.BiliBiliAvatar
										//URL = "https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID)
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
										SetThumbnail(GroupIcon).
										AddField("Youtube subscriber", engine.NearestThousandFormat(float64(SubsData.YtSubs))).
										AddField("Youtube views", engine.NearestThousandFormat(float64(SubsData.YtViews))).
										AddField("Youtube videos", engine.NearestThousandFormat(float64(SubsData.YtVideos))).
										AddField("BiliBili followers", engine.NearestThousandFormat(float64(SubsData.YtViews))).
										AddField("BiliBili views", engine.NearestThousandFormat(float64(SubsData.BiliViews))).
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
										SetThumbnail(GroupIcon).
										AddField("Youtube subscriber", engine.NearestThousandFormat(float64(SubsData.YtSubs))).
										AddField("Youtube views", engine.NearestThousandFormat(float64(SubsData.YtViews))).
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
							}
						}
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
