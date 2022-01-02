package engine

import (
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//Send new fanart to discord channel
func SendFanArtNude(Art database.DataFanart, Bot *discordgo.Session, Color int) {
	Art.Group.RemoveNillIconURL()
	for _, Member := range Art.Group.Members {
		if Art.Member.ID == Member.ID {
			var (
				ChannelData []database.DiscordChannel
				err1        error
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

			for _, Channel := range ChannelData {
				if Art.Group.GroupName == config.Indie && !Channel.IndieNotif {
					//do nothing,like my life
					continue
				} else {
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
			}
		}
	}
}
