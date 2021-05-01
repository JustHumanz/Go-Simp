package main

import (
	"context"
	"encoding/json"
	"flag"
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	lewd        bool
	gRCPconn    pilot.PilotServiceClient
)

const (
	ModuleState = "Twitter Fanart"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	Lewd := flag.Bool("LewdFanart", false, "Enable lewd fanart module")
	flag.Parse()
	lewd = *Lewd
	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
}

//main start twitter module
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

		err = json.Unmarshal(res.VtuberPayload, &VtubersData)
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
	c.AddFunc(config.TwitterFanart, CheckNew)
	if lewd {
		log.Info("Enable lewd" + ModuleState)
	} else {
		log.Info("Enable " + ModuleState)
	}

	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}

//CheckNew Check new fanart
func CheckNew() {
	wg := new(sync.WaitGroup)
	for _, GroupData := range VtubersData.VtuberData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			Fanarts, err := engine.CreatePayload(Group, config.Scraper, config.GoSimpConf.LimitConf.TwitterFanart, lewd)
			if err != nil {
				log.WithFields(log.Fields{
					"Group": Group.GroupName,
				}).Error(err)
			} else {
				for _, Art := range Fanarts {
					Color, err := engine.GetColor(config.TmpDir, Art.Photos[0])
					if err != nil {
						log.Error(err)
					}
					if config.GoSimpConf.Metric {
						gRCPconn.MetricReport(context.Background(), &pilot.Metric{
							MetricData: Art.MarshallBin(),
							State:      config.FanartState,
						})
					}

					engine.SendFanArtNude(Art, Bot, Color)
				}
			}
		}(GroupData, wg)
	}
	wg.Wait()
}
