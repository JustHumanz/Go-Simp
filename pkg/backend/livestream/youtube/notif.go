package youtube

import (
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (PushData *NotifStruct) SendNude() {
	Status := PushData.YtData.Status
	Avatar := PushData.Member.YoutubeAvatar
	Bot := runner.Bot
	YtChannel := "https://www.youtube.com/channel/" + PushData.Member.YoutubeID + "?sub_confirmation=1"
	YtURL := "https://www.youtube.com/watch?v=" + PushData.YtData.VideoID
	loc := engine.Zawarudo(PushData.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(PushData.Member.EnName, PushData.Member.JpName)
	GroupIcon := PushData.Group.IconURL
	Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
	if err != nil {
		log.Error(err)
	}

	var (
		timestart time.Time
		MsgEmbend *discordgo.MessageEmbed
	)

	if PushData.YtData.Schedul.IsZero() {
		timestart = time.Now().In(loc)
	} else {
		timestart = PushData.YtData.Schedul.In(loc)
	}

	if Status == "upcoming" {
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("New upcoming Livestream").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Start live in", durafmt.Parse(timestart.Sub(expiresAt)).LimitFirstN(2).String()).
			InlineAllFields().
			AddField("Waiting", PushData.YtData.Viewers+" Simps in Room Chat").
			SetFooter(timestart.Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed

	} else if Status == "reminder" {
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Watch "+PushData.Member.EnName+" in "+durafmt.Parse(timestart.Sub(expiresAt)).LimitFirstN(1).String()).
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Start live in", durafmt.Parse(timestart.Sub(expiresAt)).LimitFirstN(2).String()).
			InlineAllFields().
			AddField("Waiting", PushData.YtData.Viewers+" Simps in Room Chat").
			SetFooter(timestart.Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed

	} else if Status == "live" {
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Live right now").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Start live", durafmt.Parse(expiresAt.Sub(timestart)).LimitFirstN(2).String()+" Ago").
			InlineAllFields().
			AddField("Viewers", PushData.YtData.Viewers).
			SetFooter(timestart.Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed

	} else if Status == "past" && PushData.YtData.Type == "Covering" {
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Uploaded a new video").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Upload", durafmt.Parse(expiresAt.Sub(PushData.YtData.Schedul.In(loc))).LimitFirstN(2).String()+" Ago").
			InlineAllFields().
			AddField("Viewers", PushData.YtData.Viewers).
			AddField("Duration", PushData.YtData.Length).
			SetFooter(PushData.YtData.Schedul.In(loc).Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed
	} else if Status == "past" {
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Uploaded a new video").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Upload", durafmt.Parse(expiresAt.Sub(timestart)).LimitFirstN(2).String()+" Ago").
			InlineAllFields().
			AddField("Viewers", PushData.YtData.Viewers).
			SetFooter(timestart.Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed
	}
	id, DiscordChannelID := database.ChannelTag(PushData.Member.ID, 2)
	for i := 0; i < len(DiscordChannelID); i++ {
		UserTagsList := database.GetUserList(id[i], PushData.Member.ID)
		if UserTagsList != nil {
			msg, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], MsgEmbend)
			msg, err = Bot.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
				match, _ := regexp.MatchString("Unknown Channel", err.Error())
				if match {
					log.Info("Delete Discord Channel ", DiscordChannelID[i])
					database.DelChannel(DiscordChannelID[i], PushData.Group.ID)
				}
			}
		}
	}
}
