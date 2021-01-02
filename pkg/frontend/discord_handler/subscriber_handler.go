package discordhandler

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
			if CommandArray[0] == Prefix+Subscriber {
				for _, Group := range engine.GroupData {
					for _, Mem := range strings.Split(CommandArray[1], ",") {
						for _, Member := range database.GetMembers(Group.ID) {
							if Mem == strings.ToLower(Member.Name) {
								var (
									embed  *discordgo.MessageEmbed
									Avatar string
									URL    = "https://www.youtube.com/channel/" + Member.YoutubeID + "?sub_confirmation=1"
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
										URL = "https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID)
									} else {
										Avatar = Member.YoutubeAvatar
									}
								}
								Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
								if err != nil {
									log.Error(err)
								}
								if SubsData.BiliFollow != 0 {
									embed = engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetTitle(engine.FixName(Member.EnName, Member.JpName)).
										SetImage(Avatar).
										SetURL(URL).
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
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
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
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Prefix+"` command")
			if err != nil {
				log.Error(err)
			}
		}
	}
}
