package main

import (
	"context"
	"encoding/json"
	"flag"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	GroupPayload *[]database.Group
	gRCPconn     pilot.PilotServiceClient
	Proxy        = flag.Bool("MultiTor", false, "Enable multitor for bypass youtube limit TOS")
)

const (
	ModuleState = "Youtube"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	flag.Parse()
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//Start main youtube module
func main() {
	var (
		configfile config.ConfigFile
		err        error
	)

	GetPayload := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: ModuleState,
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error(), ModuleState)
			}
			log.Error("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &GroupPayload)
		if err != nil {
			log.Error(err)
		}
	}

	GetPayload()
	configfile.InitConf()

	Bot, err = discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	database.Start(configfile)

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, GetPayload)
	c.AddFunc(config.YoutubeCheckChannel, CheckYtSchedule)
	c.AddFunc(config.YoutubeCheckUpcomingByTime, CheckYtByTime)
	c.AddFunc(config.YoutubePrivateSlayer, CheckPrivate)
	c.AddFunc("0 */13 * * *", func() {
		exTknList = nil
	})
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}

var Counter = 1
var exTknList []string

func CheckYtSchedule() {
	Update := false
	Counter++
	var (
		wg sync.WaitGroup
	)

	if Counter == 4 {
		Update = true
		Counter = 1
	}

	for _, Group := range *GroupPayload {
		wg.Add(1)
		go StartCheckYT(Group, Update, &wg)
	}
	wg.Wait()
}

func CheckYtByTime() {
	for _, Group := range *GroupPayload {
		for _, Member := range Group.Members {
			if Member.YoutubeID != "" && Member.Active() {
				log.WithFields(log.Fields{
					"Vtuber": Member.EnName,
					"Group":  Group.GroupName,
				}).Info("Checking Upcoming schedule")
				YoutubeStatus, err := database.YtGetStatus(0, Member.ID, config.UpcomingStatus, "", config.Sys)
				if err != nil {
					log.Error(err)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: err.Error(),
						Service: ModuleState,
					})
				}

				for _, Youtube := range YoutubeStatus {
					Youtube.AddMember(Member).AddGroup(Group)
					if time.Since(Youtube.Schedul) > time.Until(Youtube.Schedul) {
						log.WithFields(log.Fields{
							"Vtuber":  Member.EnName,
							"Group":   Group.GroupName,
							"VideoID": Youtube.VideoID,
						}).Info("Vtuber upcoming schedule deadline,force change to live")

						Data, err := YtAPI([]string{Youtube.VideoID})
						if err != nil {
							log.Error(err)
							continue
						}
						if len(Data.Items) > 0 {
							if Data.Items[0].Snippet.VideoStatus != "none" {
								if Data.Items[0].Statistics.ViewCount != "" {
									Youtube.UpdateViewers(Data.Items[0].Statistics.ViewCount)
								} else if Data.Items[0].Statistics.ViewCount == "0" && Youtube.Viewers == "0" || Youtube.Viewers == "" {
									Viewers, err := GetWaiting(Youtube.VideoID)
									if err != nil {
										log.Error(err)
									}
									Youtube.UpdateViewers(Viewers)
								}

								if !Data.Items[0].LiveDetails.ActualStartTime.IsZero() {
									Youtube.UpdateSchdule(Data.Items[0].LiveDetails.ActualStartTime)
								}
								Key := strconv.Itoa(int(Member.ID)) + config.UpcomingStatus + config.Sys
								err = database.RemoveYtCache(Key, context.Background())
								if err != nil {
									log.Panic(err)
								}

								Youtube.UpdateStatus(config.LiveStatus).
									SetState(config.YoutubeLive).
									UpdateYt(config.LiveStatus)

								if Member.BiliRoomID != 0 {
									LiveBili, err := engine.GetRoomStatus(Member.BiliRoomID)
									if err != nil {
										log.Error(err)
									}
									if LiveBili.CheckScheduleLive() {
										Youtube.SetBiliLive(true).UpdateBiliToLive()
									}
								}

								if config.GoSimpConf.Metric {
									bit, err := Youtube.MarshalBinary()
									if err != nil {
										log.Error(err)
									}
									gRCPconn.MetricReport(context.Background(), &pilot.Metric{
										MetricData: bit,
										State:      config.LiveStatus,
									})
								}

								isMemberOnly, err := regexp.MatchString("memberonly", strings.ToLower(Youtube.Title))
								if err != nil {
									log.Error(err)
								}

								engine.SendLiveNotif(&Youtube, Bot)
								if isMemberOnly {
									Youtube.UpdateYt(config.PrivateStatus)
								}
							}
						} else if Data.Items == nil {
							log.WithFields(log.Fields{
								"Vtuber":  Member.EnName,
								"Group":   Group.GroupName,
								"VideoID": Youtube.VideoID,
							}).Warn("Upcoming deleted")
							Key := strconv.Itoa(int(Member.ID)) + config.UpcomingStatus + config.Sys
							err = database.RemoveYtCache(Key, context.Background())
							if err != nil {
								log.Panic(err)
							}

							Youtube.UpdateStatus(config.LiveStatus).
								SetState(config.YoutubeLive).
								UpdateYt(config.PrivateStatus)
						}
						//one vtuber only have one livestream right
						break
					}
					Youtube.
						SetState(config.YoutubeLive).
						UpdateStatus("reminder")
					engine.SendLiveNotif(&Youtube, Bot)
				}
			}
		}
	}
}

func GetWaiting(VideoID string) (string, error) {
	var (
		bit     []byte
		curlerr error
		urls    = "https://www.youtube.com/watch?v=" + VideoID
	)
	bit, curlerr = network.Curl(urls, nil)
	if curlerr != nil || bit == nil {
		bit, curlerr = network.CoolerCurl(urls, nil)
		if curlerr != nil {
			return config.Ytwaiting, curlerr
		}
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return config.Ytwaiting, err
	}
	for _, element := range regexp.MustCompile(`(?m)videoViewCountRenderer.*?text([0-9\s]+).+(isLive\strue)`).FindAllStringSubmatch(reg.ReplaceAllString(string(bit), " "), -1) {
		tmp := strings.Replace(element[1], " ", "", -1)
		if tmp != "" {
			config.Ytwaiting = tmp
		}
	}
	return config.Ytwaiting, nil
}

func CheckPrivate() {
	log.Info("Start Video private slayer")
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
			for _, Member := range Group.Members {
				YtData, err := database.YtGetStatus(0, Member.ID, Status, "", config.Sys)
				if err != nil {
					log.Error(err)
				}
				for _, Y := range YtData {
					Y.Status = Status
					Check(Y)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	log.Info("Done")
}
