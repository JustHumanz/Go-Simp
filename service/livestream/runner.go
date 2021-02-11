package main

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/livestream/bilibili/live"
	"github.com/JustHumanz/Go-Simp/service/livestream/bilibili/space"
	"github.com/JustHumanz/Go-Simp/service/livestream/twitch"
	"github.com/JustHumanz/Go-Simp/service/livestream/youtube"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

func main() {
	Youtube := flag.Bool("Youtube", false, "Enable youtube module")
	SpaceBiliBili := flag.Bool("SpaceBiliBili", false, "Enable Space.bilibili module")
	LiveBiliBili := flag.Bool("LiveBiliBili", false, "Enable Live.bilibili module")
	Twitch := flag.Bool("Twitch", false, "Enable Twitch module")

	flag.Parse()

	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		configfile config.ConfigFile
		Payload    database.VtubersPayload
	)

	RequestPay := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Livestream",
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error())
			}
			log.Fatalf("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}
	}
	RequestPay()

	Bot, err := discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	configfile.InitConf()
	database.Start(configfile)

	c := cron.New()
	c.Start()

	c.AddFunc(config.CheckPayload, RequestPay)

	if *Youtube {
		youtube.Start(Bot, c, Payload, configfile)
	}

	if *SpaceBiliBili {
		space.Start(Bot, c, Payload, configfile)
	}

	if *LiveBiliBili {
		live.Start(Bot, c, Payload, configfile)
	}

	if *Twitch {
		twitch.Start(Bot, c, Payload, configfile)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "Youtube",
		Enabled: *Youtube,
	})
	if err != nil {
		log.Error(err)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "SpaceBiliBili",
		Enabled: *SpaceBiliBili,
	})
	if err != nil {
		log.Error(err)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "LiveBiliBili",
		Enabled: *LiveBiliBili,
	})
	if err != nil {
		log.Error(err)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "Twitch",
		Enabled: *Twitch,
	})
	if err != nil {
		log.Error(err)
	}

	go pilot.RunHeartBeat(gRCPconn, "Livestream")
	runfunc.Run(Bot)
}
