package main

import (
	"context"
	"encoding/json"
	"flag"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot          *discordgo.Session
	lewd         = flag.Bool("LewdFanart", false, "Enable lewd fanart module")
	torTransport = flag.Bool("Tor", false, "Enable multiTor for bot transport")
	gRCPconn     pilot.PilotServiceClient
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.TwitterService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	flag.Parse()
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//main start twitter module
func main() {
	var (
		configfile config.ConfigFile
	)

	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		ServiceUUID: ServiceUUID,
		Service:     ServiceName,
	})
	if err != nil {
		if configfile.Discord != "" {
			pilot.ReportDeadService(err.Error(), ServiceName)
		}
		log.Error("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Fatalln(err)
	}

	configfile.InitConf()

	Bot = engine.StartBot(*torTransport)

	err = Bot.Open()
	if err != nil {
		log.Error(err)
	}

	database.Start(configfile)

	c := cron.New()
	c.Start()

	if *lewd {
		log.Info("Enable lewd" + ServiceName)
	} else {
		log.Info("Enable " + ServiceName)
	}

	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	Twit := &checkTwJob{}

	for {
		log.WithFields(log.Fields{
			"Service": ServiceName,
			"Running": false,
			"UUID":    ServiceUUID,
		}).Info("request for running job")

		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Request",
			ServiceUUID: ServiceUUID,
		})
		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Service": ServiceName,
				"Running": true,
				"UUID":    ServiceUUID,
			}).Info(res.Message)

			Twit.Agency = engine.UnMarshalPayload(res.VtuberPayload)
			Twit.Run()

			_, _ = client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
				Service:     ServiceName,
				Message:     "Done",
				ServiceUUID: ServiceUUID,
			})

			log.WithFields(log.Fields{
				"Service": ServiceName,
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info("reporting job was done")

		} else {
			log.WithFields(log.Fields{
				"Service": ServiceName,
				"Running": res.Run,
				"UUID":    ServiceUUID,
			}).Info(res.Message)
		}

		time.Sleep(1 * time.Minute)
	}
}

type checkTwJob struct {
	wg      sync.WaitGroup
	Agency  []database.Group
	Reverse bool
}

func (i *checkTwJob) Run() {
	Cek := func(Member database.Member, w *sync.WaitGroup) {
		defer w.Done()

		log.WithFields(log.Fields{
			"Hashtag": Member.TwitterHashtag,
			"Vtuber":  Member.Name,
			"Agency":  Member.Group.GroupName,
			"Lewd":    false,
		}).Info("Start curl twitter")

		Fanarts, err := Member.ScrapTwitterFanart(engine.InitTwitterScraper(), false)
		if err != nil {
			log.WithFields(log.Fields{
				"vtuber":  Member.Name,
				"agency":  Member.Group.GroupName,
				"hashtag": Member.TwitterHashtag,
			}).Error(err)
			return
		}
		for _, Art := range Fanarts {
			if config.GoSimpConf.Metric {
				gRCPconn.MetricReport(context.Background(), &pilot.Metric{
					MetricData: Art.MarshallBin(),
					State:      config.FanartState,
				})
			}
			engine.SendFanArtNude(Art, Bot)
		}

		if *lewd && Member.TwitterLewd != "" {

			log.WithFields(log.Fields{
				"Hashtag": Member.TwitterLewd,
				"Vtuber":  Member.Name,
				"Agency":  Member.Group.GroupName,
				"Lewd":    true,
			}).Info("Start curl twitter")

			Fanarts, err := Member.ScrapTwitterFanart(engine.InitTwitterScraper(), true)
			if err != nil {
				log.WithFields(log.Fields{
					"vtuber":  Member.Name,
					"agency":  Member.Group.GroupName,
					"hashtag": Member.TwitterLewd,
				}).Error(err)
				return
			}
			for _, Art := range Fanarts {
				if config.GoSimpConf.Metric {
					gRCPconn.MetricReport(context.Background(), &pilot.Metric{
						MetricData: Art.MarshallBin(),
						State:      config.FanartState,
					})
				}
				engine.SendFanArtNude(Art, Bot)
			}
		}
	}

	if i.Reverse {
		for j := len(i.Agency) - 1; j >= 0; j-- {
			Grp := i.Agency

			for kk, Vtuber := range Grp[j].Members {
				if Vtuber.TwitterHashtag != "" || Vtuber.TwitterLewd != "" {
					Vtuber.Group = Grp[j]
					i.wg.Add(1)
					go Cek(Vtuber, &i.wg)
				}

				if kk%10 == 0 && kk != 0 {
					i.wg.Wait()
				}
			}
		}
		i.Reverse = false
		i.wg.Wait()

	} else {
		for _, G := range i.Agency {
			for kk, Vtuber := range G.Members {
				if Vtuber.TwitterHashtag != "" || Vtuber.TwitterLewd != "" {
					Vtuber.Group = G
					i.wg.Add(1)
					go Cek(Vtuber, &i.wg)
				}

				if kk%10 == 0 && kk != 0 {
					i.wg.Wait()
				}
			}
		}
		i.Reverse = true
		i.wg.Wait()
	}
}
