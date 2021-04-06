package notif

import (
	"context"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func SendNude(Art database.DataFanart, Bot *discordgo.Session, Color int) {
	Art.Group.RemoveNillIconURL()
	for _, Member := range Art.Group.Members {
		if Art.Member.ID == Member.ID {
			var (
				ChannelData []database.DiscordChannel
				err1        error
				NewEmbed    *discordgo.MessageEmbed
			)
			if Art.Lewd {
				ChannelData, err1 = database.ChannelTag(Member.ID, 1, config.LewdChannel, Member.Region)
				if err1 != nil {
					log.Error(err1)
				}
			} else {
				ChannelData, err1 = database.ChannelTag(Member.ID, 1, config.Default, Member.Region)
				if err1 != nil {
					log.Error(err1)
				}
			}

			icon := ""
			if Art.State == config.PixivArt {
				icon = config.PixivIMG
			} else if Art.State == config.TwitterArt {
				icon = config.TwitterIMG
			} else {
				icon = config.BiliBiliIMG
			}

			tags := ""
			for i, Channel := range ChannelData {
				Channel.SetMember(Member)
				UserTagsList, err := Channel.GetUserList(context.Background())
				if err != nil {
					log.Error(err)
				} else {
					if UserTagsList != nil {
						tags = strings.Join(UserTagsList, " ")
						NewEmbed = engine.NewEmbed().
							SetAuthor(strings.Title(Art.Group.GroupName), Art.Group.IconURL).
							SetTitle(Art.Author).
							SetURL(Art.PermanentURL).
							SetThumbnail(Art.AuthorAvatar).
							SetDescription(Art.Text).
							SetImage(Art.Photos...).
							AddField("User Tags", tags).
							SetFooter(Art.State, icon).
							InlineAllFields().
							SetColor(Color).MessageEmbed
					} else {
						NewEmbed = engine.NewEmbed().
							SetAuthor(strings.Title(Art.Group.GroupName), Art.Group.IconURL).
							SetTitle(Art.Author).
							SetURL(Art.PermanentURL).
							SetThumbnail(Art.AuthorAvatar).
							SetDescription(Art.Text).
							SetImage(Art.Photos...).
							SetFooter(Art.State, icon).
							InlineAllFields().
							SetColor(Color).MessageEmbed
					}

					if tags == "" && Art.Group.GroupName == config.Indie && !Channel.IndieNotif {
						//do nothing,like my life
					} else {
						tmp, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewEmbed)
						if err != nil {
							log.Error(tmp, err.Error())
							err = Channel.DelChannel(err.Error())
							if err != nil {
								log.Error(err)
							}
						}
						err = engine.Reacting(map[string]string{
							"ChannelID": Channel.ChannelID,
						}, Bot)
						if err != nil {
							log.Error(err)
						}
					}
					if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
						log.WithFields(log.Fields{
							"State": Art.State + " Fanart",
						}).Warn(config.FanartSleep)
						time.Sleep(config.FanartSleep)
					}
				}
			}
		}
	}
}
