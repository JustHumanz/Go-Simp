package main

import (
	"encoding/json"
	"math/rand"

	"github.com/JustHumanz/Go-simp/tools/network"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		log.Error("joined unavailable guild", g.Guild.ID)
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
		timejoin, err := g.Guild.JoinedAt.Parse()
		if err != nil {
			log.Error(err)
			return
		}
		DataGuild := database.Guild{
			ID:   g.Guild.ID,
			Name: g.Guild.Name,
			Join: timejoin,
		}
		Info := DataGuild.CheckGuild()
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

					//Save discord name to database
					err := DataGuild.InputGuild()
					if err != nil {
						log.Error(err)
						return
					}

					PayloadBytes, err := json.Marshal(map[string]interface{}{
						"embeds": []interface{}{
							map[string]interface{}{
								"description": "A Guild Invited me",
								"fields": []interface{}{
									map[string]interface{}{
										"name":   "GuildName",
										"value":  g.Guild.Name,
										"inline": true,
									},
									map[string]interface{}{
										"name":   "OwnerID",
										"value":  g.Guild.OwnerID,
										"inline": true,
									},
									map[string]interface{}{
										"name":   "Member Count",
										"value":  g.Guild.MemberCount,
										"inline": true,
									},
									map[string]interface{}{
										"name":   "Join Date",
										"value":  timejoin.String(),
										"inline": true,
									},
									map[string]interface{}{
										"name":   "Region",
										"value":  g.Guild.Region,
										"inline": true,
									},
								},
							},
						},
					})
					if err != nil {
						log.Error(err)
					}
					err = network.CurlPost(config.DiscordWebHook, PayloadBytes)
					if err != nil {
						log.Error(err)
					}
					return
				}
			}
		}
	}
}
