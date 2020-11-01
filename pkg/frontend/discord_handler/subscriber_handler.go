package discord_handler

import (
	"math/rand"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func gacha() bool {
	return rand.Float32() < 0.5
}

func SubsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := config.PGeneral
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, prefix) {
		CommandArray := strings.Split(m.Content, " ")
		if CommandArray[0] == prefix+Subscriber {
			for _, Group := range engine.GroupData {
				for _, Mem := range strings.Split(CommandArray[1], ",") {
					for _, Member := range database.GetName(Group.ID) {
						if Mem == strings.ToLower(Member.Name) {
							var (
								embed  *discordgo.MessageEmbed
								Avatar string
								Url    = "https://www.youtube.com/channel/" + Member.YoutubeID + "?sub_confirmation=1"
							)
							SubsData := Member.GetSubsCount()
							if gacha() {
								Avatar = Member.YoutubeAvatar
							} else {
								if Member.BiliRoomID != 0 {
									Avatar = Member.BiliBiliAvatar
									Url = "https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID)
								} else {
									Avatar = Member.YoutubeAvatar
								}
							}
							Color, err := engine.GetColor("/tmp/bilia.tmp", m.Author.AvatarURL("128"))
							if err != nil {
								log.Error(err)
							}
							if SubsData.BiliFollow != 0 {
								embed = engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetTitle(engine.FixName(Member.EnName, Member.JpName)).
									SetImage(Avatar).
									SetURL(Url).
									AddField("Youtube subscriber", strconv.Itoa(SubsData.YtSubs)).
									AddField("Youtube views", strconv.Itoa(SubsData.YtViews)).
									AddField("Youtube videos", strconv.Itoa(SubsData.YtVideos)).
									AddField("BiliBili followers", strconv.Itoa(SubsData.BiliFollow)).
									AddField("BiliBili views", strconv.Itoa(SubsData.BiliViews)).
									AddField("BiliBili videos", strconv.Itoa(SubsData.BiliVideos)).
									AddField("Twitter followers", strconv.Itoa(SubsData.TwFollow)).
									InlineAllFields().
									SetColor(Color).MessageEmbed
							} else {
								embed = engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
									SetTitle(engine.FixName(Member.EnName, Member.JpName)).
									SetImage(Avatar).
									AddField("Youtube subscriber", strconv.Itoa(SubsData.YtSubs)).
									AddField("Youtube views", strconv.Itoa(SubsData.YtViews)).
									AddField("Youtube videos", strconv.Itoa(SubsData.YtVideos)).
									AddField("Twitter followers", strconv.Itoa(SubsData.TwFollow)).
									InlineAllFields().
									SetColor(Color).MessageEmbed
							}
							msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
							if err != nil {
								log.Error(err, msg)
							}
						}
					}
				}
			}
		}
	}
}
