package space

import (
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data CheckSctruct) SendNude() {
	var (
		BotSession = engine.BotSession
	)
	if Data.VideoList != nil {
		Color, err := engine.GetColor("/tmp/bilispace.tmps", Data.MemberFace)
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"Vtuber": Data.MemberName,
		}).Info("New video uploaded")

		for _, video := range Data.VideoList {
			Embed := engine.NewEmbed().
				SetAuthor(Data.MemberName, Data.MemberFace, Data.MemberUrl).
				SetTitle("Uploaded new video").
				SetDescription(video.Title).
				SetImage(video.Pic).
				SetThumbnail(Data.GroupIcon).
				SetURL("https://www.bilibili.com/video/"+video.Bvid).
				AddField("Type ", video.VideoType).
				AddField("Duration ", video.Length).
				AddField("Viwers ", strconv.Itoa(video.Play)).
				SetFooter(durafmt.Parse(time.Now().In(loc).Sub(time.Unix(int64(video.Created), 0).In(loc))).LimitFirstN(2).String()+" Ago", config.BiliBiliIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed

			ID, DiscordChannelID := database.ChannelTag(Data.MemberID, 2)
			for i := 0; i < len(DiscordChannelID); i++ {
				UserTagsList := database.GetUserList(ID[i], Data.MemberID)
				if UserTagsList != nil {
					msg, err := BotSession.ChannelMessageSendEmbed(DiscordChannelID[i], Embed)
					if err != nil {
						log.Error(msg, err)
					} else {
						msg, err = BotSession.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							log.Error(msg, err)
						}
					}
				}
			}
		}
	}
}
