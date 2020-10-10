package youtube

import (
	"regexp"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (PushData NotifStruct) SendNuke(Status string) {
	YtChannlID := strings.Split(PushData.Member.YoutubeID, "\n")[0]
	Avatar := PushData.Member.YoutubeAvatar
	YtChannel := "https://www.youtube.com/channel/" + YtChannlID + "?sub_confirmation=1"
	YtURL := "https://www.youtube.com/watch?v=" + PushData.Data.VideoID
	loc := Zawarudo(PushData.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(PushData.Member.EnName, PushData.Member.JpName)
	GroupIcon := PushData.Group.IconURL
	Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
	if err != nil {
		log.Error(err)
	}

	var (
		msg, msg1, msg2, msg3, msg4, msg5 string
		wg                                sync.WaitGroup
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

	Embed := engine.NewEmbed().
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

	id, DiscordChannelID := database.ChannelTag(PushData.Member.ID, 2)
	for i := 0; i < len(DiscordChannelID); i++ {
		UserTagsList := database.GetUserList(id[i], PushData.Member.ID)

		go func(DiscordChannel string, wg *sync.WaitGroup) {
			wg.Add(1)
			if UserTagsList != nil {
				msg, err := BotSession.ChannelMessageSendEmbed(DiscordChannel, Embed)
				if err != nil {
					log.Error(msg, err)
				}

				msg, err = BotSession.ChannelMessageSend(DiscordChannel, "UserTags: "+strings.Join(UserTagsList, " "))
				if err != nil {
					log.Error(msg, err)
					match, _ := regexp.MatchString("Unknown Channel", err.Error())
					if match {
						log.Info("Delete Discord Channel ", DiscordChannel)
						database.DelChannel(DiscordChannel, PushData.Member.ID)
					}
				}
			}
		}(DiscordChannelID[i], &wg)
	}
	wg.Wait()
}

func Zawarudo(Region string) *time.Location {
	if Region == "ID" {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		return loc
	} else if Region == "JP" {
		loc, _ := time.LoadLocation("Asia/Tokyo")
		return loc
	} else if Region == "CN" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		return loc
	} else {
		loc, _ := time.LoadLocation("UTC")
		return loc
	}
}
