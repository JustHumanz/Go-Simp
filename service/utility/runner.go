package main

import (
	"context"
	"encoding/json"
	"math/rand"
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

var (
	KanoPayload = []string{
		"uEeSOh5jOk0",
		"fokontvrio0",
		"zZ0r3RgZzXA",
		"CD1AvkqS8oE",
		"yXxccEqgAO4",
		"N-ZTFGlD8Rg",
		"Do47UIW_TXw",
		"qxNkMzlV-FU",
		"gkVCUuPuF8I",
		"I8yBbwRurAE",
		"pCa_oSjBU1A",
		"ImSW1g02FUk",
		"GPo-g6tHH_4",
		"paVYNlZ5Xuk",
		"HqIx1CVPBsI",
		"t6o2TpzpPGU",
		"y-bPf-6OHws",
		"_kj5xKz8CDM",
		"CaMKMdkLbck",
		"PIrx5lqQbGU",
		"I7jnsXxHs8k",
		"3Li-FfypZYE",
		"CybFOypDQjY",
		"l3j2Ud8Mo4A",
		"zgqu6_nyRGY",
		"vxZtflYGjA8",
		"RfDN1JMMCM4",
		"eyiYja05RAI",
	}
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	var (
		configfile config.ConfigFile
		GuildList  []string
	)
	res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
		Message: "Send me nude",
		Service: "Utility",
	})
	if err != nil {
		if configfile.Discord != "" {
			pilot.ReportDeadService(err.Error(), "Utility")
		}
		log.Fatalf("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}
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

	GuildCount := func() int {
		for _, GuildID := range Bot.State.Guilds {
			GuildList = append(GuildList, GuildID.ID)
		}
		return len(Bot.State.Guilds)
	}
	GuildCount()
	Donation := configfile.DonationLink
	configfile.InitConf()
	database.Start(configfile)

	Bot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if (m.Emoji.MessageFormat() == configfile.Emoji.Livestream[0] || m.Emoji.MessageFormat() == configfile.Emoji.Livestream[1]) && m.UserID != BotInfo.ID {
			UserState, err := database.GetChannelMessage(m.MessageID)
			if err != nil {
				log.Error(err)
			}
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
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> just added "+UserState.Member.Name+" to their list.")
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
						_, err := s.ChannelMessageSend(m.ChannelID, "<@"+m.UserID+"> just removed "+UserState.Member.Name+" from their list.")
						if err != nil {
							log.Error(err)
						}
					}
				}
			}
		}
	})

	c := cron.New()
	c.Start()
	if configfile.DonationLink != "" {
		c.AddFunc(config.DonationMsg, func() {
			Img := config.GoSimpIMG
			if rand.Float32() < 0.5 {
				if rand.Float32() < 0.5 {
					Img = engine.LewdIMG()
				} else {
					Img = engine.MaintenanceIMG()
				}
				Img = engine.NotFoundIMG()
			} else {
				if rand.Float32() < 0.5 {
					Img = engine.NotFoundIMG()
				} else {
					Img = engine.MaintenanceIMG()
				}
				Img = engine.LewdIMG()
			}

			Music := "https://www.youtube.com/watch?v=" + KanoPayload[engine.RandomNum(0, len(KanoPayload)-1)]
			Bot.ChannelMessageSendEmbed(database.GetRanChannel(), engine.NewEmbed().
				SetTitle("Donate").
				SetURL(Donation).
				SetThumbnail(BotInfo.AvatarURL("128")).
				SetImage(Img).
				SetColor(14807034).
				SetDescription("Enjoy the bot?\nhelp dev to pay server,domain and database for development of "+BotInfo.Username).
				AddField("Ko-Fi", "[Link]("+Donation+")").
				AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[Top.gg]("+configfile.TopGG+")").
				AddField("or help dev simping kano/鹿乃 with listening her music", "[鹿乃チャンネルofficial]("+Music+")\nHope you like her voice ❤️").
				SetFooter("~advertisement").MessageEmbed)
		})
	}
	c.AddFunc(config.CheckServerCount, func() {
		log.Info("POST bot info to top.gg")
		dblClient, err := dbl.NewClient(os.Getenv("TOPGG"))
		if err != nil {
			log.Error(err)
		}

		err = dblClient.PostBotStats(BotInfo.ID, &dbl.BotStatsPayload{
			Shards: []int{GuildCount()},
		})
		if err != nil {
			log.Error(err)
		}
	})

	go pilot.RunHeartBeat(gRCPconn, "Utility")
	runfunc.Run(Bot)
}
