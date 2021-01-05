package space

import (
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data CheckSctruct) SendNude() {
	if Data.VideoList != nil {
		Color, err := engine.GetColor(config.TmpDir, Data.MemberFace)
		if err != nil {
			log.Error(err)
		}
		log.WithFields(log.Fields{
			"Vtuber": Data.MemberName,
		}).Info("New video uploaded")

		for _, video := range Data.VideoList {
			ID, DiscordChannelID := database.ChannelTag(Data.MemberID, 2, "LiveOnly")
			for i := 0; i < len(DiscordChannelID); i++ {
				UserTagsList := database.GetUserList(ID[i], Data.MemberID)
				if UserTagsList != nil {
					msg, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], engine.NewEmbed().
						SetAuthor(Data.MemberName, Data.MemberFace, Data.MemberUrl).
						SetTitle("Uploaded new video").
						SetDescription(video.Title).
						SetImage(video.Pic).
						SetThumbnail(Data.GroupIcon).
						SetURL("https://www.bilibili.com/video/"+video.Bvid).
						AddField("Type ", video.VideoType).
						AddField("Duration ", video.Length).
						AddField("Viwers ", engine.NearestThousandFormat(float64(video.Play))).
						SetFooter(durafmt.Parse(time.Now().In(loc).Sub(time.Unix(int64(video.Created), 0).In(loc))).LimitFirstN(2).String()+" Ago", config.BiliBiliIMG).
						InlineAllFields().
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(msg, err)
					} else {
						msg, err = Bot.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							log.Error(msg, err)
						}
					}
				}
			}
		}
	}
}
