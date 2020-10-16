package space

import (
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data CheckSctruct) SendNude(Color int) {
	var (
		BotSession = engine.BotSession
		wg         = new(sync.WaitGroup)
	)
	if Data.VideoList != nil {
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
				wg.Add(1)
				go func(DiscordChannel string, wg *sync.WaitGroup) {
					defer wg.Done()
					if UserTagsList != nil {
						msg, err := BotSession.ChannelMessageSendEmbed(DiscordChannel, Embed)
						msg, err = BotSession.ChannelMessageSend(DiscordChannel, "UserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							log.Error(msg, err)
							match, _ := regexp.MatchString("Unknown Channel", err.Error())
							if match {
								log.Info("Delete Discord Channel ", DiscordChannel)
								database.DelChannel(DiscordChannel, Data.MemberID)
							}
						}
					}
				}(DiscordChannelID[i], wg)
			}
			wg.Wait()
		}
	}
}
