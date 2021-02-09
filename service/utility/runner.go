package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/top-gg/go-dbl"
)

func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	var (
		configfile config.ConfigFile
	)
	RequestPay := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Livestream",
		})
		if err != nil {
			log.Fatalf("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Panic(err)
		}
	}
	RequestPay()

	Bot, err := discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}
	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	BotInfo, err := Bot.User("@me")
	if err != nil {
		log.Panic(err)
	}
	Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if (m.Emoji.MessageFormat() == configfile.Emoji.Livestream[0] || m.Emoji.MessageFormat() == configfile.Emoji.Livestream[1]) && m.UserID != BotInfo.ID {
			UserState := database.GetChannelMessage(m.MessageID)
			if UserState != nil {
				if m.Emoji.MessageFormat() == configfile.Emoji.Livestream[0] {
					UserInfo, err := s.User(m.MessageReaction.UserID)
					if err != nil {
						log.Error(err)
					}
					log.WithFields(log.Fields{
						"UserID":    UserInfo.ID,
						"UserName":  UserInfo.Username,
						"ChannelID": m.ChannelID,
						"Group":     UserState.Group.GroupName,
						"Vtuber":    UserState.Member.Name,
					}).Info("New user add from reac")
					UserState.SetDiscordID(UserInfo.ID).
						SetDiscordUserName(UserInfo.Username)
					err = UserState.Adduser()
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> "+err.Error())
						if err != nil {
							log.Error(err)
						}
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> Just add "+UserState.Member.Name)
						if err != nil {
							log.Error(err)
						}
					}
				} else if m.Emoji.MessageFormat() == configfile.Emoji.Livestream[1] {
					UserInfo, err := s.User(m.MessageReaction.UserID)
					if err != nil {
						log.Error(err)
					}
					log.WithFields(log.Fields{
						"UserID":    UserInfo.ID,
						"UserName":  UserInfo.Username,
						"ChannelID": m.ChannelID,
						"Group":     UserState.Group.GroupName,
						"Vtuber":    UserState.Member.Name,
					}).Info("New user del from reac")
					UserState.SetDiscordID(UserInfo.ID).
						SetDiscordUserName(UserInfo.Username)
					err = UserState.Deluser()
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> "+err.Error())
						if err != nil {
							log.Error(err)
						}
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> Just remove "+UserState.Member.Name)
						if err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	})

	Donation := configfile.DonationLink
	configfile.InitConf()
	database.Start(configfile)

	c := cron.New()
	c.Start()
	if configfile.DonationLink != "" {
		c.AddFunc(config.DonationMsg, func() {
			Bot.ChannelMessageSendEmbed(database.GetRanChannel(), engine.NewEmbed().
				SetTitle("Donate").
				SetURL(Donation).
				SetThumbnail(BotInfo.AvatarURL("128")).
				SetImage(config.GoSimpIMG).
				SetColor(14807034).
				SetDescription("Enjoy the bot?\ndon't forget to support this bot and dev").
				AddField("Ko-Fi", "[Link]("+Donation+")").
				AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[Top.gg]("+configfile.TopGG+")").
				AddField("or give some star on github", "[Github](https://github.com/JustHumanz/Go-Simp)").MessageEmbed)
		})
	}
	c.AddFunc(config.CheckServerCount, func() {
		log.Info("POST bot info to top.gg")
		dblClient, err := dbl.NewClient(os.Getenv("TOPGG"))
		if err != nil {
			log.Error(err)
		}

		err = dblClient.PostBotStats(BotInfo.ID, &dbl.BotStatsPayload{
			Shards: []int{database.GetGuildsCount()},
		})
		if err != nil {
			log.Error(err)
		}
	})

	go pilot.RunHeartBeat(gRCPconn, "Utility")
	runfunc.Run(Bot)
}
