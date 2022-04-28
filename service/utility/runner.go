package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/top-gg/go-dbl"
)

const ServiceName = "Utility"

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
		"_PVZr0iJiug",
		"cLJL6uRezSI",
		"6XZek8E_SiE",
	}

	ServiceUUID = uuid.New().String()
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
	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		Service:     ServiceName,
		ServiceUUID: ServiceUUID,
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
	configfile.InitConf()
	database.Start(configfile)

	Bot := engine.StartBot(false)

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
	c.AddFunc(config.DonationMsg, func() {
		Img := config.GoSimpIMG
		Num := engine.RandomNum(1, 4)
		if Num == 1 {
			Img = engine.LewdIMG()
		} else if Num == 2 {
			Img = engine.MaintenanceIMG()
		} else if Num == 3 {
			Img = engine.NotFoundIMG()
		} else {
			Img = engine.Gif()
		}

		Music := "https://www.youtube.com/watch?v=" + KanoPayload[engine.RandomNum(0, len(KanoPayload)-1)]

		if configfile.DonationLink != "" {
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

		} else {
			Bot.ChannelMessageSendEmbed(database.GetRanChannel(), engine.NewEmbed().
				SetTitle("Donate").
				SetURL(Donation).
				SetThumbnail(BotInfo.AvatarURL("128")).
				SetImage(Img).
				SetColor(14807034).
				SetDescription("Enjoy the bot?\nhelp dev to pay server,domain and database for development of "+BotInfo.Username).
				AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[Top.gg]("+configfile.TopGG+")").
				AddField("or help dev simping kano/鹿乃 with listening her music", "[鹿乃チャンネルofficial]("+Music+")\nHope you like her voice ❤️").
				SetFooter("~advertisement").MessageEmbed)

		}
	})

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

	c.AddFunc(config.YoutubePrivateSlayer, func() {
		log.Info("Start Video private slayer")
		var GroupPayload *[]database.Group
		req, err := gRCPconn.GetAgencyPayload(context.Background(), &pilot.ServiceMessage{
			Message:     "Request agency data",
			Service:     ServiceName,
			ServiceUUID: ServiceUUID,
		})
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(req.AgencyVtubers, &GroupPayload)
		if err != nil {
			log.Error(err)
		}

		Check := func(Youtube database.LiveStream) {
			if Youtube.Status == "upcoming" && time.Since(Youtube.Schedul) > time.Until(Youtube.Schedul) {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Member only video")
				Youtube.UpdateYt("past")
				engine.RemoveEmbed(Youtube.VideoID, Bot)
			} else if Youtube.Status == "live" && Youtube.Viewers == "0" || Youtube.Status == "live" && int(math.Round(time.Since(Youtube.Schedul).Hours())) > 30 {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Member only video")
				Youtube.UpdateYt("past")
				engine.RemoveEmbed(Youtube.VideoID, Bot)
			}

			_, err := network.Curl("https://i3.ytimg.com/vi/"+Youtube.VideoID+"/hqdefault.jpg", nil)
			if err != nil && strings.HasPrefix(err.Error(), "404") && Youtube.Status != "private" {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Private Video")
				Youtube.UpdateYt("private")
			} else if err == nil && Youtube.Status == "private" {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("From Private Video to past")
				Youtube.UpdateYt("past")
			} else {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Video was daijobu")
			}
		}

		log.Info("Start Check Private video")
		for _, Status := range []string{config.UpcomingStatus, config.PastStatus, config.LiveStatus, config.PrivateStatus} {
			for _, Group := range *GroupPayload {
				YtData, err := Group.GetYtLiveStream(Status, "")
				if err != nil {
					log.Error(err)
				}
				for _, Y := range YtData {
					Y.Status = Status
					Check(Y)
					err = Y.RemoveCache(fmt.Sprintf("%d-%s-%s", Group.ID, Group.GroupName, Status))
					if err != nil {
						log.Panic(err)
					}
				}
			}
			time.Sleep(10 * time.Second)
		}
		log.Info("Done")
	})

	c.AddFunc(config.CheckUser, func() {
		log.Info("Start check deleted user")
		UsersDeleted := []string{}
		Users := database.GetAllUser()
		for _, v := range Users {
			User, err := Bot.User(v)
			if err != nil {
				log.Error(err, v)
			}
			if strings.HasPrefix(User.Username, "Deleted") {
				log.WithFields(log.Fields{
					"userId": v,
				}).Info("User not found or user deleted ther account")
				UsersDeleted = append(UsersDeleted, v)
			}
		}

		log.Info("Start deleting user ", UsersDeleted)
		database.DeleteDeletedUser(UsersDeleted)
	})

	log.Info("Start check deleted user")
	UsersDeleted := []string{}
	Users := database.GetAllUser()
	for _, v := range Users {
		User, err := Bot.User(v)
		if err != nil {
			log.Error(err, v)
		}
		if strings.HasPrefix(User.Username, "Deleted") {
			log.WithFields(log.Fields{
				"userId": v,
			}).Info("User not found or user deleted ther account")
			UsersDeleted = append(UsersDeleted, v)
		}
	}

	log.Info("Start deleting user ", UsersDeleted)
	database.DeleteDeletedUser(UsersDeleted)

	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	runfunc.Run(Bot)
}
