package youtube

import (
	"math"
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

	var (
		MsgEmbend *discordgo.MessageEmbed
		Timestart time.Time
	)

	if !PushData.YtData.Schedul.IsZero() {
		Timestart = PushData.YtData.Schedul
	} else if PushData.YtData.Schedul.IsZero() && !PushData.YtData.Published.IsZero() {
		Timestart = PushData.YtData.Published
	} else if PushData.YtData.Schedul.IsZero() && PushData.YtData.Published.IsZero() {
		Timestart = time.Now()
	}

	if Status == "upcoming" {
		Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
		if err != nil {
			log.Error(err)
		}
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("New upcoming Livestream").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(2).String()).
			InlineAllFields().
			AddField("Waiting", PushData.YtData.Viewers+" Simps in Room Chat").
			SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed

	} else if Status == "live" {
		Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
		if err != nil {
			log.Error(err)
		}
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Live right now").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Start live", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(2).String()+" Ago").
			InlineAllFields().
			AddField("Viewers", PushData.YtData.Viewers).
			SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed

	} else if Status == "past" && PushData.YtData.Type == "Covering" {
		Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
		if err != nil {
			log.Error(err)
		}
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Uploaded a new video").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Upload", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(2).String()+" Ago").
			AddField("Viewers", PushData.YtData.Viewers).
			AddField("Duration", PushData.YtData.Length).
			InlineAllFields().
			SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed
	} else if Status == "past" {
		Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
		if err != nil {
			log.Error(err)
		}
		MsgEmbend = engine.NewEmbed().
			SetAuthor(VtuberName, Avatar, YtChannel).
			SetTitle("Uploaded a new video").
			SetDescription(PushData.YtData.Title).
			SetImage(PushData.YtData.Thumb).
			SetThumbnail(GroupIcon).
			SetURL(YtURL).
			AddField("Type ", PushData.YtData.Type).
			AddField("Upload", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(2).String()+" Ago").
			AddField("Viewers", PushData.YtData.Viewers).
			AddField("Duration", PushData.YtData.Length).
			InlineAllFields().
			SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
			SetColor(Color).MessageEmbed
	}

	if Status == "reminder" {
		Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
		if err != nil {
			log.Error(err)
		}
		id, DiscordChannelID := database.ChannelTag(PushData.Member.ID, 2)
		UpcominginMinutes := int(math.Round(Timestart.In(loc).Sub(time.Now().In(loc)).Minutes()))
		for i := 0; i < len(DiscordChannelID); i++ {
			k := 0
			for ii := 0; ii < 70; ii += 10 {
				k = ii + 6
				if UpcominginMinutes > ii && UpcominginMinutes < k {
					UserTagsList := database.GetUserReminderList(id[i], PushData.Member.ID, k)
					if UserTagsList != nil {
						msg, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], engine.NewEmbed().
							SetAuthor(VtuberName, Avatar, YtChannel).
							SetTitle(PushData.Member.EnName+" Live in "+durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()).
							SetDescription(PushData.YtData.Title).
							SetImage(PushData.YtData.Thumb).
							SetThumbnail(GroupIcon).
							SetURL(YtURL).
							AddField("Type ", PushData.YtData.Type).
							AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(2).String()).
							InlineAllFields().
							AddField("Waiting", PushData.YtData.Viewers+" Simps in Room Chat").
							SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(msg, err)
							match, _ := regexp.MatchString("Unknown Channel", err.Error())
							if match {
								log.Info("Delete Discord Channel ", DiscordChannelID[i])
								database.DelChannel(DiscordChannelID[i], PushData.Group.ID)
							}
						}
						msg, err = Bot.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
					}
				}
			}
		}
		return
	}

	id, DiscordChannelID := database.ChannelTag(PushData.Member.ID, 2)
	for i := 0; i < len(DiscordChannelID); i++ {
		UserTagsList := database.GetUserList(id[i], PushData.Member.ID)
		if UserTagsList != nil {
			msg, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], MsgEmbend)
			if err != nil {
				log.Error(msg, err)
				match, _ := regexp.MatchString("Unknown Channel", err.Error())
				if match {
					log.Info("Delete Discord Channel ", DiscordChannelID[i])
					database.DelChannel(DiscordChannelID[i], PushData.Group.ID)
				}
			}
			msg, err = Bot.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
		}
	}
}
