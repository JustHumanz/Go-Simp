package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		H3llcome   = []string{config.Bonjour, config.Howdy, config.Guten, config.Koni, config.Selamat, config.Assalamu, config.Approaching}
		configfile config.ConfigFile
		GuildList  []string
	)
	res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
		Message: "Send me nude",
		Service: "Guild",
	})
	if err != nil {
		log.Fatalf("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}

	Bot := configfile.StartBot()

	err = Bot.Open()
	if err != nil {
		log.Error(err)
	}

	BotInfo, err := Bot.User("@me")
	if err != nil {
		log.Error(err)
	}

	configfile.InitConf()
	database.Start(configfile)

	for _, GuildID := range Bot.State.Guilds {
		GuildList = append(GuildList, GuildID.ID)
	}

	Bot.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		if g.Unavailable {
			log.Error("joined unavailable guild", g.Guild.ID)
			return
		}

		for _, v := range GuildList {
			if v == g.ID {
				return
			}
		}

		GuildList = append(GuildList, g.ID)

		Join, err := g.JoinedAt.Parse()
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"GuildName": g.Name,
			"OwnerID":   g.OwnerID,
			"JoinDate":  Join.Format(time.RFC822),
		}).Info("New invite")

		_, err = s.ChannelMessageSendEmbed(configfile.InviteLog, engine.NewEmbed().
			SetTitle("A Guild Invited "+BotInfo.Username).
			SetColor(14807034).
			AddField("GuildName", g.Name).
			AddField("OwnerID", g.OwnerID).
			AddField("JoinDate", Join.Format(time.RFC822)).
			AddField("Member Count", strconv.Itoa(len(g.Members))).
			InlineAllFields().MessageEmbed)
		if err != nil {
			log.Error(err)
		}

		for _, Channel := range g.Guild.Channels {
			if Channel.ID == g.Guild.ID && Channel.Type == 0 {
				Donation := configfile.DonationLink
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
						AddField("Support "+BotInfo.Username, "[Ko-Fi]("+Donation+")").
						AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[top.gg]("+config.GoSimpConf.TopGG+")").
						AddField("give some star on github", "[Github]("+config.GuildSupport+")").MessageEmbed)
				}
				return
			}
		}

	})
	log.Info("Guild handler ready.......")
	go pilot.RunHeartBeat(gRCPconn, "Guild")
	runfunc.Run(Bot)
}
