package space

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data CheckSctruct) SendNude() {
	if Data.VideoList != nil {
		Color, err := engine.GetColor(config.TmpDir, Data.Member.BiliBiliAvatar)
		if err != nil {
			log.Error(err)
		}
		expiresAt := time.Now().In(loc)
		log.WithFields(log.Fields{
			"Vtuber": Data.Member.Name,
		}).Info("New video uploaded")

		for _, video := range Data.VideoList {
			//ID, DiscordChannelID
			var (
				wg          sync.WaitGroup
				ChannelData = database.ChannelTag(Data.Member.ID, 2, "NotLiveOnly")
			)
			for i, v := range ChannelData {
				wg.Add(1)
				go func(Channel database.DiscordChannel, wg *sync.WaitGroup) {
					defer wg.Done()
					UserTagsList := Channel.GetUserList(context.Background()) //database.GetUserList(Channel.ID, Data.MemberID)
					if UserTagsList != nil {
						msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
							SetAuthor(Data.Member.Name, Data.Member.BiliBiliAvatar, "https://space.bilibili/"+strconv.Itoa(Data.Member.BiliBiliID)).
							SetTitle("Uploaded new video").
							SetDescription(video.Title).
							SetImage(video.Pic).
							SetThumbnail(Data.Group.IconURL).
							SetURL("https://www.bilibili.com/video/"+video.Bvid).
							AddField("Type ", video.VideoType).
							AddField("Duration ", video.Length).
							InlineAllFields().
							AddField("Viwers ", engine.NearestThousandFormat(float64(video.Play))).
							SetFooter(durafmt.Parse(expiresAt.Sub(time.Unix(int64(video.Created), 0).In(loc))).LimitFirstN(2).String()+" Ago", config.BiliBiliIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(msg, err)
						} else {
							msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "UserTags: "+strings.Join(UserTagsList, " "))
							if err != nil {
								log.Error(msg, err)
							}
						}
					}
				}(v, &wg)
				//Wait every ge 5 discord channel
				if i%config.Waiting == 0 && configfile.LowResources {
					log.WithFields(log.Fields{
						"Func":  "BiliBili space",
						"Value": config.Waiting,
					}).Warn("Waiting send message")
					wg.Wait()
					expiresAt = time.Now().In(loc)
				}
			}
			wg.Wait()
		}
	}
}
