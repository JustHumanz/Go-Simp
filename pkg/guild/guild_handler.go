package main

import (
	"math/rand"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		log.Info("joined unavailable guild", g.Guild.ID)
		return
	}
	New := false
	for _, Guild := range GuildList {
		if Guild != g.ID {
			New = true
		}
	}
	if New {
		log.WithFields(log.Fields{
			"Member": g.Guild.MemberCount,
			"Owner":  g.Guild.OwnerID,
			"Reg":    g.Guild.Region,
		}).Info(g.Guild.Name, " join the battle")

		GuildList = append(GuildList, g.Guild.ID)
		sqlite := OpenLiteDB(PathLiteDB)
		timejoin, err := g.Guild.JoinedAt.Parse()
		if err != nil {
			log.Error(err)
			return
		}
		DataGuild := Guild{
			ID:     g.Guild.ID,
			Name:   g.Guild.Name,
			Join:   timejoin,
			Dbconn: sqlite,
		}
		Info := DataGuild.CheckGuild()
		SendInvite, err := s.UserChannelCreate(config.OwnerDiscordID)
		if err != nil {
			log.Error(err)
		}

		if Info == 0 {
			for _, Channel := range g.Guild.Channels {
				BotPermission, err := s.UserChannelPermissions(BotID.ID, Channel.ID)
				if err != nil {
					log.Error(err)
					return
				}
				if Channel.Type == 0 && BotPermission&2048 != 0 {
					s.ChannelMessageSendEmbed(Channel.ID, engine.NewEmbed().
						SetTitle("Thx for invite me to this server <3 ").
						SetURL("https://go-simp.human-z.tech/Guide/").
						SetThumbnail(config.GoSimpIMG).
						SetImage(H3llcome[rand.Intn(len(H3llcome))]).
						SetColor(14807034).
						AddField("Setup", "You can watch [here](https://go-simp.human-z.tech/Guide/)").
						AddField("Need support?", "Join [dev server](https://discord.com/invite/ydWC5knbJT)").
						InlineAllFields().MessageEmbed)

					//send server name to my discord
					err := DataGuild.InputGuild()
					if err != nil {
						log.Error(err)
						return
					}
					s.ChannelMessageSend(SendInvite.ID, g.Guild.Name+" invited me")
					return
				}
			}
		}
		KillSqlConn(sqlite)
	}
}
