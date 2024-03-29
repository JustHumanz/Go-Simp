package main

import (
	"context"
	"encoding/json"
	"flag"
	"math"
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
	proxy        = flag.Bool("MultiTOR", false, "Enable MultiTOR for scrapping yt rss")
)

const (
	ModuleState = config.YoutubeCheckerModule
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
	Bot = configfile.StartBot()

	database.Start(configfile)

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, GetPayload)
	//c.AddFunc(config.YoutubePrivateSlayer, CheckPrivate) //TODO make new service for private slayer
	c.AddFunc("0 */2 * * *", func() {
		engine.ExTknList = nil
	})
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkYtCekJob struct {
	wg      sync.WaitGroup
	Reverse bool
	Update  bool
	Counter int
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	YoutubeChecker := &checkYtCekJob{
		Counter: 1,
		Update:  true,
	}
	for {

		if YoutubeChecker.Counter == 15 {
			YoutubeChecker.Update = true
			YoutubeChecker.Counter = 1
		}

		res := func() *pilot.RunJob {
			log.WithFields(log.Fields{
				"Service":  ModuleState,
				"Running":  true,
				"YtUpdate": YoutubeChecker.Update,
			}).Info("request for running job")

			if YoutubeChecker.Update {
				tmp, err := client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
					Service: ModuleState,
					Message: "Update",
					Alive:   true,
				})
				if err != nil {
					log.Error(err)
				}
				return tmp
			} else {
				tmp, err := client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
					Service: ModuleState,
					Message: "New",
					Alive:   true,
				})
				if err != nil {
					log.Error(err)
				}
				return tmp
			}
		}()

		if res.Run {
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": true,
			}).Info(res.Message)

			YoutubeChecker.Run()
			_, _ = client.RunModuleJob(context.Background(), &pilot.ServiceMessage{
				Service: ModuleState,
				Message: "Done",
				Alive:   false,
			})
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": false,
			}).Info("reporting job was done")
		} else {
			log.WithFields(log.Fields{
				"Service": ModuleState,
				"Running": false,
			}).Info(res.Message)
		}

		YoutubeChecker.Counter++
		YoutubeChecker.Update = false
		time.Sleep(1 * time.Minute)
	}
}

func (i *checkYtCekJob) Run() {
	if i.Reverse {
		for j := len(*GroupPayload) - 1; j >= 0; j-- {
			i.wg.Add(1)
			Grp := *GroupPayload
			go StartCheckYT(Grp[j], i.Update, &i.wg)
		}
		i.Reverse = false

	} else {
		for _, G := range *GroupPayload {
			i.wg.Add(1)
			StartCheckYT(G, i.Update, &i.wg)
		}
		i.Reverse = true
	}
	i.wg.Wait()
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
				YtData, Key, err := database.YtGetStatus(map[string]interface{}{
					"MemberID":   Member.ID,
					"MemberName": Member.Name,
					"Status":     Status,
					"State":      config.Sys,
				})
				if err != nil {
					log.Error(err)
				}
				for _, Y := range YtData {
					Y.Status = Status
					Check(Y)
					err = Y.RemoveCache(Key)
					if err != nil {
						log.Panic(err)
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	log.Info("Done")
}
