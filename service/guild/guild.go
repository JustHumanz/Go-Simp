package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/backend/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	H3llcome  = []string{config.Bonjour, config.Howdy, config.Guten, config.Koni, config.Selamat, config.Assalamu, config.Approaching}
	GuildList []string
	BotID     *discordgo.User
)

func main() {
	conf, err := config.ReadConfig("../../config.toml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Bot, err := discordgo.New("Bot " + config.BotConf.Discord)
	if err != nil {
		log.Error(err)
	}
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	BotID, err = Bot.User("@me")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	database.Start(conf.CheckSQL())

	for _, GuildID := range Bot.State.Guilds {
		GuildList = append(GuildList, GuildID.ID)
	}

	Bot.AddHandler(Guild)
	log.Info("Guild handler ready.......")

	runfunc.Run(Bot)
}

func Guild(s *discordgo.Session, g *discordgo.GuildCreate) {
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
				}
				if Channel.Type == 0 && BotPermission&2048 != 0 {
					Donation := config.BotConf.DonationLink
					if Donation == "" {
						s.ChannelMessageSendEmbed(Channel.ID, engine.NewEmbed().
							SetTitle("Thx for invite me to this server <3 ").
							SetURL(config.GuideURL).
							SetThumbnail(config.GoSimpIMG).
							SetImage(H3llcome[rand.Intn(len(H3llcome))]).
							SetColor(14807034).
							AddField("Setup", "You can see [here]("+config.GuideURL+")").
							AddField("Need support?", "Join [dev server]("+config.GuildSupport+")").
							InlineAllFields().MessageEmbed)
					} else {
						s.ChannelMessageSendEmbed(Channel.ID, engine.NewEmbed().
							SetTitle("Thx for invite me to this server <3 ").
							SetURL(config.GuideURL).
							SetThumbnail(config.GoSimpIMG).
							SetImage(H3llcome[rand.Intn(len(H3llcome))]).
							SetColor(14807034).
							AddField("Setup", "You can see [here]("+config.GuideURL+")").
							AddField("Need support?", "Join [dev server]("+config.GuildSupport+")").
							InlineAllFields().
							AddField("Support "+BotID.Username, "[Ko-Fi]("+Donation+")").
							AddField("if you a broke gang,you can upvote "+BotID.Username, "[top.gg]("+config.BotConf.TopGG+")").
							AddField("give some star on github", "[Github]("+config.GuildSupport+")").MessageEmbed)
					}

					//Save discord name to database
					err := DataGuild.InputGuild()
					if err != nil {
						log.Error(err)
					}

					PayloadBytes, err := json.Marshal(map[string]interface{}{
						"embeds": []interface{}{
							map[string]interface{}{
								"description": "A Guild Invited " + BotID.Username,
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
					err = network.CurlPost(config.BotConf.DiscordWebHook, PayloadBytes)
					if err != nil {
						log.Error(err)
					}
				}
				break
			}
		}
	}
}
