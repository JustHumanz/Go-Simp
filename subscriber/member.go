package subscriber

import (
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

var (
	BiliSession string
	Bot         *discordgo.Session
	Data        []database.GroupName
)

func Start(c *cron.Cron) {
	Bot = engine.BotSession
	BiliSession = config.BiliBiliSes
	Data = database.GetGroup()
	if BiliSession == "" {
		log.Error("BiliBili Session not found")
		os.Exit(1)
	}
	go Bot.AddHandler(SubsMessage)
	c.AddFunc("@every 1h0m0s", CheckYtSubsCount)
	c.AddFunc("@every 0h15m0s", CheckTwFollowCount)
	c.AddFunc("@every 0h20m0s", CheckBiliFollowCount)
	log.Info("Subs&Follow Checker module ready")
}

func SendNude(Embed *discordgo.MessageEmbed, Group database.GroupName) {
	for _, Channel := range Group.GetChannelByGroup() {
		msg, err := Bot.ChannelMessageSendEmbed(Channel, Embed)
		if err != nil {
			log.Error(msg, err)
		}
	}
}
func gacha() bool {
	return rand.Float32() < 0.5
}

func SubsMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := config.PGeneral
	m.Content = strings.ToLower(m.Content)
	CommandArray := strings.Split(m.Content, " ")
	if strings.HasPrefix(m.Content, prefix) {
		if CommandArray[0] == prefix+"subscriber" {
			for _, Group := range Data {
				Members := database.GetName(Group.ID)
				for _, Member := range Members {
					if CommandArray[1] == strings.ToLower(Member.Name) {
						var (
							embed  *discordgo.MessageEmbed
							Avatar string
						)
						SubsData := Member.GetSubsCount()
						if gacha() {
							Avatar = Member.YoutubeAvatar
						} else {
							if Member.BiliRoomID != 0 {
								Avatar = Member.BiliBiliAvatar
							} else {
								Avatar = Member.YoutubeAvatar
							}
						}
						Color, err := engine.GetColor("/tmp/asa3.tmp", m.Author.Avatar)
						if err != nil {
							log.Error(err)
						}
						if SubsData.BiliFollow != 0 {
							embed = engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("80"), "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetImage(Avatar).
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
								SetAuthor(m.Author.Username, m.Author.AvatarURL("80"), "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
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

type Subs struct {
	Success bool   `json:"success"`
	Service string `json:"service"`
	T       int64  `json:"t"`
	Data    struct {
		LvIdentifier string `json:"lv_identifier"`
		Subscribers  int    `json:"subscribers"`
		Videos       int    `json:"videos"`
		Views        int    `json:"views"`
	} `json:"data"`
}

type BiliBiliStat struct {
	LikeView LikeView
	Follow   BiliFollow
	Videos   int
}

type LikeView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Archive struct {
			View int `json:"view"`
		} `json:"archive"`
		Article struct {
			View int `json:"view"`
		} `json:"article"`
		Likes int `json:"likes"`
	} `json:"data"`
}

type BiliFollow struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int `json:"mid"`
		Following int `json:"following"`
		Whisper   int `json:"whisper"`
		Black     int `json:"black"`
		Follower  int `json:"follower"`
	} `json:"data"`
}
