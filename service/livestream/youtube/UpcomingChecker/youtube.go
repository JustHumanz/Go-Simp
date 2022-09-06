package main

import (
	"context"
	"encoding/json"
	"flag"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	gRCPconn     pilot.PilotServiceClient
	proxy        = flag.Bool("MultiTOR", false, "Enable MultiTOR for scrapping yt rss")
	torTransport = flag.Bool("Tor", false, "Enable multiTor for bot transport")
	agency       = flag.Bool("Agency", false, "Enable scraping for vtuber agency")
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.YoutubeCheckerService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	flag.Parse()
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

// Start main youtube module
func main() {
	var (
		configfile config.ConfigFile
	)

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
		log.Fatalln(err)
	}

	configfile.InitConf()
	Bot = engine.StartBot(*torTransport)

	database.Start(configfile)

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)
}

type checkYtCekJob struct {
	agency  []database.Group
	Reverse bool
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	YoutubeChecker := &checkYtCekJob{}

	hostname := engine.GetHostname()

	for {
		res, err := client.RequestRunJobsOfService(context.Background(), &pilot.ServiceMessage{
			Service:     ServiceName,
			Message:     "New",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Error(err)
		}

		if res.Run {
			log.WithFields(log.Fields{
				"Running":        res.Run,
				"UUID":           ServiceUUID,
				"Agency Payload": res.VtuberMetadata,
				"IsReverse":      YoutubeChecker.Reverse,
			}).Info(res.Message)

			YoutubeChecker.agency = engine.UnMarshalPayload(res.VtuberPayload)
			if len(YoutubeChecker.agency) == 0 {
				msg := "vtuber agency was nill,force close the unit"
				pilot.ReportDeadService(msg, ServiceName)
				log.Fatalln(msg)
			}
			YoutubeChecker.Run()

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
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info(res.Message)
		}

		time.Sleep(1 * time.Minute)
	}
}

func (i *checkYtCekJob) Run() {
	if i.Reverse {
		for j := len(i.agency) - 1; j >= 0; j-- {
			Grp := i.agency
			StartCheckYT(Grp[j])
		}

	} else {
		for _, G := range i.agency {
			StartCheckYT(G)
		}
	}
}
