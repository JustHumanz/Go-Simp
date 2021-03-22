package twitch

import (
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	TwitchClient *helix.Client
	VtubersData  database.VtubersPayload
	configfile   config.ConfigFile
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	b.AddFunc(config.Twitch, CheckTwitch)
	var err error
	TwitchClient, err = helix.NewClient(&helix.Options{
		ClientID:     config.GoSimpConf.Twitch.ClientID,
		ClientSecret: config.GoSimpConf.Twitch.ClientSecret,
	})
	if err != nil {
		log.Error(err)
	}
	TwitchClient.SetUserAccessToken(config.GoSimpConf.GetTwitchAccessToken())
	log.Info("Enable Twitch module")
}

func CheckTwitch() {
	for _, Group := range VtubersData.VtuberData {
		for _, Member := range Group.Members {
			if Member.TwitchName != "" {
				log.WithFields(log.Fields{
					"Group":      Group.GroupName,
					"VtuberName": Member.Name,
				}).Info("Checking Twitch")

				result, err := TwitchClient.GetStreams(&helix.StreamsParams{
					UserLogins: []string{Member.TwitchName},
				})
				if err != nil {
					log.Error(err)
				}

				ResultDB, err := database.GetTwitch(Member.ID)
				if err != nil {
					log.Error(err)
				}

				if len(result.Data.Streams) > 0 {
					for _, Stream := range result.Data.Streams {
						if ResultDB.Status == config.PastStatus && Stream.Type == config.LiveStatus {
							if strings.ToLower(Stream.UserName) == strings.ToLower(Member.TwitchName) {
								GameResult, err := TwitchClient.GetGames(&helix.GamesParams{
									IDs: []string{Stream.GameID},
								})
								if err != nil {
									log.Error(err)
								}

								Stream.ThumbnailURL = strings.Replace(Stream.ThumbnailURL, "{width}", "1280", -1)
								Stream.ThumbnailURL = strings.Replace(Stream.ThumbnailURL, "{height}", "720", -1)

								ResultDB.UpdateStatus(config.LiveStatus).
									UpdateViewers(strconv.Itoa(Stream.ViewerCount)).
									UpdateThumbnail(Stream.ThumbnailURL).
									UpdateSchdule(Stream.StartedAt)

								if len(GameResult.Data.Games) > 0 {
									ResultDB.UpdateGame(GameResult.Data.Games[0].Name)
								} else {
									ResultDB.UpdateGame("-")
								}

								err = ResultDB.UpdateTwitch()
								if err != nil {
									log.Error(err)
								}

								err = SendNotif(*ResultDB)
								if err != nil {
									log.Error(err)
								}

								log.WithFields(log.Fields{
									"Group":      Group.GroupName,
									"VtuberName": Member.Name,
								}).Info("Change Twitch status to Live")

							}
						} else if Stream.Type == config.LiveStatus && ResultDB.Status == config.LiveStatus {
							log.WithFields(log.Fields{
								"Group":      Group.GroupName,
								"VtuberName": Member.Name,
								"Viewers":    Stream.ViewerCount,
							}).Info("Update Viewers")

							ResultDB.UpdateViewers(strconv.Itoa(Stream.ViewerCount)).UpdateTwitch()
						}
					}
				} else if ResultDB.Status == config.LiveStatus && len(result.Data.Streams) == 0 {
					ResultDB.UpdateStatus(config.PastStatus)
					err = ResultDB.UpdateTwitch()
					if err != nil {
						log.Error(err)
					}
					log.WithFields(log.Fields{
						"Group":      Group.GroupName,
						"VtuberName": Member.Name,
					}).Info("Change Twitch status to Past")
					engine.RemoveEmbed("Twitch"+Member.TwitchName, Bot)
				}
			}
		}
	}
}
