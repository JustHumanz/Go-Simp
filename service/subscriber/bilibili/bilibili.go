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

const (
	ServiceName = config.SubscriberService + "BiliBili"
)

var (
	Bot         *discordgo.Session
	configfile  config.ConfigFile
	gRCPconn    pilot.PilotServiceClient
	ServiceUUID = uuid.New().String()
	Agency      []database.Group
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
	c.AddFunc(config.BiliBiliFollowers, CheckBiliBili)
	log.Info("Add bilibili followers to cronjob")

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

func CheckBiliBili() {
	BiliBiliSession := map[string]string{
		"Cookie": "SESSDATA=" + configfile.BiliSess,
	}
	for _, Group := range Agency {
		for _, Member := range Group.Members {
			if Member.BiliBiliID != 0 && Member.Active() {
				var (
					bilistate BiliBiliStat
				)

				body := func() []byte {
					if configfile.MultiTOR != "" {
						body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/relation/stat?vmid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
						if curlerr != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(curlerr)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     curlerr.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body
					} else {
						body, curlerr := network.Curl("https://api.bilibili.com/x/relation/stat?vmid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
						if curlerr != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(curlerr)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     curlerr.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body
					}
				}()

				err := json.Unmarshal(body, &bilistate.Follow)
				if err != nil {
					log.Error(err)
				}

				body2 := func() []byte {
					if configfile.MultiTOR != "" {
						body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
						if curlerr != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(curlerr)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     curlerr.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body
					} else {
						body, curlerr := network.Curl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
						if curlerr != nil {
							log.WithFields(log.Fields{
								"Agency": Group.GroupName,
								"Vtuber": Member.Name,
							}).Error(curlerr)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message:     curlerr.Error(),
								Service:     ServiceName,
								ServiceUUID: ServiceUUID,
							})
						}
						return body
					}
				}()

				err = json.Unmarshal(body2, &bilistate.LikeView)
				if err != nil {
					log.Error(err)
				}

				baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Member.BiliBiliID) + "&ps=100"
				url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
				for f := 0; f < len(url); f++ {
					body3 := func() []byte {
						if configfile.MultiTOR != "" {
							body, curlerr := network.CoolerCurl(url[f], BiliBiliSession)
							if curlerr != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(curlerr)
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message:     curlerr.Error(),
									Service:     ServiceName,
									ServiceUUID: ServiceUUID,
								})
							}
							return body
						} else {
							body, curlerr := network.Curl(url[f], BiliBiliSession)
							if curlerr != nil {
								log.WithFields(log.Fields{
									"Agency": Group.GroupName,
									"Vtuber": Member.Name,
								}).Error(curlerr)
								gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
									Message:     curlerr.Error(),
									Service:     ServiceName,
									ServiceUUID: ServiceUUID,
								})
							}
							return body
						}
					}()
					var video engine.SpaceVideo
					err := json.Unmarshal(body3, &video)
					if err != nil {
						log.Error(err)
					}
					bilistate.Videos += video.Data.Page.Count
				}

				BiliFollowDB, err := Member.GetSubsCount()
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
				}

				BiliFollowDB.SetMember(Member).SetGroup(Group).UpdateState(config.BiliLive)

				if bilistate.Follow.Data.Follower != 0 {
					log.WithFields(log.Fields{
						"Past BiliBili Follower":    BiliFollowDB.BiliFollow,
						"Current BiliBili Follower": bilistate.Follow.Data.Follower,
						"Vtuber":                    Member.Name,
					}).Info("Update BiliBili Follower")

					err := BiliFollowDB.SetMember(Member).SetGroup(Group).
						UpdateBiliBiliFollowers(bilistate.Follow.Data.Follower).
						UpdateBiliBiliVideos(bilistate.Videos).
						UpdateBiliBiliViewers(bilistate.LikeView.Data.Archive.View).
						UpdateState(config.BiliLive).UpdateSubs()
					if err != nil {
						log.Error(err)
					}

					if BiliFollowDB.BiliFollow != bilistate.Follow.Data.Follower {
						if bilistate.Follow.Data.Follower <= 10000 {
							for i := 0; i < 1000001; i += 100000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Member.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.WithFields(log.Fields{
											"Agency": Group.GroupName,
											"Vtuber": Member.Name,
										}).Error(err)
									}

									Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22BiliBili%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"
									engine.SubsSendEmbed(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										SetTitle(engine.FixName(Member.EnName, Member.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Viewers", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
										AddField("Videos", strconv.Itoa(bilistate.Videos)).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										InlineAllFields().
										AddField("Graph", Graph).
										SetColor(Color).MessageEmbed, Group, Member, Bot)
								}
							}
						} else {
							for i := 0; i < 10001; i += 1000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Member.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}

									engine.SubsSendEmbed(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										SetTitle(engine.FixName(Member.EnName, Member.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Views", engine.NearestThousandFormat(float64(bilistate.LikeView.Data.Archive.View))).
										AddField("Videos", engine.NearestThousandFormat(float64(bilistate.Videos))).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										InlineAllFields().
										SetColor(Color).MessageEmbed, Group, Member, Bot)
								}
							}
						}
					}

					bin, err := BiliFollowDB.MarshalBinary()
					if err != nil {
						log.Error(err)
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

type BiliBiliStat struct {
	LikeView LikeView
	Follow   BiliFollow
	Videos   int
}

type LikeView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Archive struct {
			View int `json:"view"`
		} `json:"archive"`
		Article struct {
			View int `json:"view"`
		} `json:"article"`
		Likes int `json:"likes"`
	} `json:"data"`
}

type BiliFollow struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int `json:"mid"`
		Following int `json:"following"`
		Whisper   int `json:"whisper"`
		Black     int `json:"black"`
		Follower  int `json:"follower"`
	} `json:"data"`
}
