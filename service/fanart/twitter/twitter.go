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
	like         = flag.Bool("Like", false, "Update like fanart")
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
	hostname := engine.GetHostname()

	for {

		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "Request",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Running":        true,
				"UUID":           ServiceUUID,
				"Agency Payload": res.VtuberMetadata,
			}).Info(res.Message)

			Twit.Agency = engine.UnMarshalPayload(res.VtuberPayload)
			if len(Twit.Agency) == 0 {
				msg := "vtuber agency was nill,force close the unit"
				pilot.ReportDeadService(msg, ServiceName)
				log.Fatalln(msg)
			}
			Twit.Run()

			_, _ = client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
				Service:     ServiceName,
				Message:     "Done",
				ServiceUUID: ServiceUUID,
			})

			log.WithFields(log.Fields{
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info("reporting job was done")

		} else {
			log.WithFields(log.Fields{
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

		Fanarts, err := Member.ScrapTwitterFanart(engine.InitTwitterScraper(), false, *like)
		if err != nil {
			log.WithFields(log.Fields{
				"Vtuber":  Member.Name,
				"Agency":  Member.Group.GroupName,
				"Hashtag": Member.TwitterHashtag,
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

			Fanarts, err := Member.ScrapTwitterFanart(engine.InitTwitterScraper(), true, *like)
			if err != nil {
				log.WithFields(log.Fields{
					"Vtuber":  Member.Name,
					"Agency":  Member.Group.GroupName,
					"Hashtag": Member.TwitterLewd,
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
