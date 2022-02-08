package engine

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//Send new fanart to discord channel
func SendFanArtNude(Art database.DataFanart, Bot *discordgo.Session) {
	Art.Group.RemoveNillIconURL()
	for _, Member := range Art.Group.Members {
		if Art.Member.ID == Member.ID {

			ChannelData := func() []database.DiscordChannel {
				if Art.Lewd {
					Data, err1 := database.ChannelTag(Member.ID, 1, config.LewdChannel, Member.Region)
					if err1 != nil {
						log.Error(err1)
					}
					return Data
				} else {
					Data, err1 := database.ChannelTag(Member.ID, 1, config.Default, Member.Region)
					if err1 != nil {
						log.Error(err1)
					}
					return Data
				}
			}()

			Color := func() int {
				if Art.FilePath != "" {
					Color, err := GetColor("", Art.FilePath)
					if err != nil {
						log.Error(err)
					}
					return Color
				} else {
					Color, err := GetColor(config.TmpDir, Art.Photos[0])
					if err != nil {
						log.Error(err)
					}
					return Color
				}
			}()

			icon := ""
			if Art.State == config.PixivArt {
				icon = config.PixivIMG
			} else if Art.State == config.TwitterArt {
				icon = config.TwitterIMG
			} else {
				icon = config.BiliBiliIMG
			}

			var (
				wgg sync.WaitGroup
			)

			for i, C := range ChannelData {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
				defer cancel()
				wgg.Add(1)

				go func(ctx context.Context, Channel database.DiscordChannel, wg *sync.WaitGroup) {
					defer wg.Done()

					done := make(chan struct{})

					go func() {
						if Art.Group.GroupName == config.Indie && !Channel.IndieNotif {
							//do nothing,like my life
							return
						} else {
							log.WithFields(log.Fields{
								"channelId": Channel.ChannelID,
								"vtuber":    Member.Name,
							}).Info("Send pic")
							tmp, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewEmbed().
								SetAuthor(strings.Title(Art.Group.GroupName), Art.Group.IconURL).
								SetTitle(Art.Author).
								SetURL(Art.PermanentURL).
								SetThumbnail(Art.AuthorAvatar).
								SetDescription(Art.Text).
								SetImage(Art.Photos...).
								SetFooter(Art.State, icon).
								InlineAllFields().
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(tmp, err.Error())
								err = Channel.DelChannel(err.Error())
								if err != nil {
									log.Error(err)
								}
							}
							err = Reacting(map[string]string{
								"ChannelID": Channel.ChannelID,
							}, Bot)
							if err != nil {
								log.Error(err)
							}
						}
						done <- struct{}{}
					}()

					select {
					case <-done:
						{
						}
					case <-ctx.Done():
						{
							log.WithFields(log.Fields{
								"channelID":      Channel.ID,
								"discordChannel": Channel.ChannelID,
								"vtuber":         Art.Member.Name,
							}).Error(ctx.Err())
						}
					}

				}(ctx, C, &wgg)

				if i%10 == 0 && config.GoSimpConf.LowResources && i != 0 {
					tidur := 10 * time.Second
					log.WithFields(log.Fields{
						"Type":  "Sleep",
						"Value": tidur,
					}).Info("Waiting send message")
					time.Sleep(tidur)
				}

			}

			log.WithFields(log.Fields{
				"Type": "Wait",
			}).Info("Waiting send message")
			wgg.Wait()
		}
	}
}
