package youtube

import (
	"regexp"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (PushData NotifStruct) SendNude() {
	id, DiscordChannelID := database.ChannelTag(PushData.Member.ID, 2)
	for i := 0; i < len(DiscordChannelID); i++ {
		UserTagsList := database.GetUserList(id[i], PushData.Member.ID)
		if UserTagsList != nil {
			msg, err := BotSession.ChannelMessageSendEmbed(DiscordChannelID[i], PushData.Embed)
			if err != nil {
				log.Error(msg, err)
				match, _ := regexp.MatchString("Unknown Channel", err.Error())
				if match {
					log.Info("Delete Discord Channel ", DiscordChannelID[i])
					database.DelChannel(DiscordChannelID[i], PushData.Group.ID)
				}
			} else {
				msg, err = BotSession.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
			}
		}
	}
}

func (PushData NotifStruct) GetEmbed(Status string) NotifStruct {
	Avatar := PushData.Member.YoutubeAvatar
	YtChannel := "https://www.youtube.com/channel/" + PushData.Member.YoutubeID + "?sub_confirmation=1"
	YtURL := "https://www.youtube.com/watch?v=" + PushData.Data.VideoID
	loc := engine.Zawarudo(PushData.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(PushData.Member.EnName, PushData.Member.JpName)
	GroupIcon := PushData.Group.IconURL
	Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
	if err != nil {
		log.Error(err)
	}

	var (
		msg, msg1, msg2, msg3, msg4, msg5 string
	)
	if Status == "upcoming" {
		log.Info("New upcoming live stream")
		msg = "Start live in"
		msg1 = durafmt.Parse(PushData.Data.Schedul.In(loc).Sub(expiresAt)).LimitFirstN(2).String()
		msg2 = "New upcoming live stream"
		msg3 = PushData.Data.Schedul.In(loc).Format(time.RFC822)
		msg4 = "Waiting"
		msg5 = PushData.Data.Viewers + " Simps in Room Chat"

	} else if Status == "reminder" {
		msg = "Start live in"
		msg1 = durafmt.Parse(PushData.Data.Schedul.In(loc).Sub(expiresAt)).LimitFirstN(2).String()
		msg2 = "Reminder"
		msg3 = PushData.Data.Schedul.In(loc).Format(time.RFC822)
		msg4 = "Waiting"
		msg5 = PushData.Data.Viewers + " Simps in Room Chat"

	} else if Status == "live" {
		log.Info("New live stream")
		msg = "Start live"
		msg1 = durafmt.Parse(expiresAt.Sub(PushData.Data.Schedul.In(loc))).LimitFirstN(2).String() + " Ago"
		msg2 = "Live right now"
		msg3 = PushData.Data.Schedul.In(loc).Format(time.RFC822)
		msg4 = "Viewers"
		msg5 = PushData.Data.Viewers

	} else if Status == "past" && PushData.Data.Type == "Covering" {
		log.Info("New cover has uploaded")
		msg = "Upload"
		msg1 = durafmt.Parse(expiresAt.Sub(PushData.Data.Schedul.In(loc))).LimitFirstN(2).String() + " Ago"
		msg2 = "Uploaded new video"
		msg3 = PushData.Data.Schedul.In(loc).Format(time.RFC822)
		msg4 = "Viewers"
		msg5 = PushData.Data.Viewers
	} else if Status == "past" {
		log.Info("Suddenly upload new video")
		msg = "Upload"
		msg1 = durafmt.Parse(expiresAt.Sub(PushData.Data.Schedul.In(loc))).LimitFirstN(2).String() + " Ago"
		msg2 = "Uploaded new video"
		msg3 = PushData.Data.Schedul.In(loc).Format(time.RFC822)
		msg4 = "Viewers"
		msg5 = PushData.Data.Viewers
	}
	PushData.Embed = engine.NewEmbed().
		SetAuthor(VtuberName, Avatar, YtChannel).
		SetTitle(msg2).
		SetDescription(PushData.Data.Title).
		SetImage(PushData.Data.Thumb).
		SetThumbnail(GroupIcon).
		SetURL(YtURL).
		AddField("Type ", PushData.Data.Type).
		AddField(msg, msg1).
		InlineAllFields().
		AddField(msg4, msg5).
		SetFooter(msg3, config.YoutubeIMG).
		SetColor(Color).MessageEmbed
	return PushData
}
