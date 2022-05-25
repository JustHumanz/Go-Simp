package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	Bot         *discordgo.Session
	configfile  config.ConfigFile
	gRCPconn    pilot.PilotServiceClient
	ServiceUUID = uuid.New().String()
	Agency      []database.Group
)

const (
	ServiceName = config.SubscriberService + "Youtube"
)

//Init service
func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

}

func main() {
	//Get config file from pilot
	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		Service:     ServiceName,
		ServiceUUID: ServiceUUID,
	})
	if err != nil {
		if configfile.Discord != "" {
			pilot.ReportDeadService(err.Error(), ServiceName)
		}
		log.Error("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Error(err)
	}

	configfile.InitConf()
	Bot = engine.StartBot(false)

	database.Start(configfile)

	c := cron.New()
	c.Start()

	c.AddFunc(config.YoutubeSubscriber, CheckYoutube)
	log.Info("Add youtube subscriber to cronjob")

	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)

	hostname := engine.GetHostname()

	go func() {
		tmp, err := gRCPconn.GetAgencyPayload(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Refresh payload",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(tmp.AgencyVtubers, &Agency)
		if err != nil {
			log.Error(err)
		}

		time.Sleep(1 * time.Hour)
	}()
	runfunc.Run(Bot)
}

func CheckYoutube() {
	var YTstate Subs
	Token := engine.GetYtToken()
	for _, Group := range Agency {
		for _, Member := range Group.Members {
			if !Member.IsYtNill() && Member.Active() {
				body, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+Member.YoutubeID+"&key="+*Token, nil)
				if err != nil {
					if err.Error() == "403 Forbidden" {
						log.Error(err, string(body))
					} else {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						}).Error(err)

						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message:     err.Error(),
							Service:     ServiceName,
							ServiceUUID: ServiceUUID,
						})
					}
				}

				err = json.Unmarshal(body, &YTstate)
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
				}

				for _, Item := range YTstate.Items {
					if Member.YoutubeID == Item.ID && !Item.Statistics.HiddenSubscriberCount {
						YtSubsDB, err := Member.GetSubsCount()
						if err != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(err)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     err.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}

						YtSubsDB.SetMember(Member).SetGroup(Group).UpdateState(config.YoutubeLive)

						YTSubscriberCount, err := strconv.Atoi(Item.Statistics.SubscriberCount)
						if err != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(err)
						}

						SendNotif := func(SubsCount string) {
							err = Member.RemoveSubsCache()
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							log.WithFields(log.Fields{
								"Vtuber": Member.Name,
								"Group":  Group.GroupName,
							}).Info("Congratulation for " + SubsCount + " subscriber")

							Color, err := engine.GetColor(config.TmpDir, Member.YoutubeAvatar)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							VideoCount, err := strconv.Atoi(Item.Statistics.VideoCount)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							ViewCount, err := strconv.Atoi(Item.Statistics.ViewCount)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22Youtube%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"

							engine.SubsSendEmbed(engine.NewEmbed().
								SetAuthor(Group.GroupName, Group.IconURL, "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetThumbnail(config.YoutubeIMG).
								SetDescription("Congratulation for "+SubsCount+" subscriber").
								SetImage(Member.YoutubeAvatar).
								AddField("Viewers", engine.NearestThousandFormat(float64(ViewCount))).
								AddField("Videos", engine.NearestThousandFormat(float64(VideoCount))).
								InlineAllFields().
								AddField("Graph", Graph).
								SetURL("https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
								SetColor(Color).MessageEmbed, Group, Member, Bot)

						}
						if YtSubsDB.YtSubs != YTSubscriberCount {

							log.WithFields(log.Fields{
								"Past Youtube subscriber":    YtSubsDB.YtSubs,
								"Current Youtube subscriber": YTSubscriberCount,
								"Vtuber":                     Member.Name,
							}).Info("Update Youtube subscriber")

							VideoCount, err := strconv.Atoi(Item.Statistics.VideoCount)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							ViewCount, err := strconv.Atoi(Item.Statistics.ViewCount)
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							err = YtSubsDB.UpdateYoutubeSubs(YTSubscriberCount).
								UpdateYoutubeVideos(VideoCount).
								UpdateYoutubeViewers(ViewCount).
								UpdateSubs()
							if err != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(err)
							}

							if YTSubscriberCount >= 1000000 {
								for i := 0; i < 10000001; i += 1000000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
										)
									}
								}
							} else if YTSubscriberCount >= 100000 {
								for i := 0; i < 1000001; i += 100000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
										)
									}
								}
							} else if YTSubscriberCount >= 10000 {
								for i := 0; i < 100001; i += 10000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
										)
									}
								}
							} else if YTSubscriberCount >= 1000 {
								for i := 0; i < 10001; i += 1000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
										)
									}
								}
							}
						}

						bin, err := YtSubsDB.MarshalBinary()
						if err != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(err)
						}

						if config.GoSimpConf.Metric {
							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bin,
								State:      config.SubsState,
							})
						}
					}
				}
			}
		}
	}
}

type Subs struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind       string `json:"kind"`
		Etag       string `json:"etag"`
		ID         string `json:"id"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			CommentCount          string `json:"commentCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}
