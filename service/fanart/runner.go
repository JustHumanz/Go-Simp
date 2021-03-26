package main

import (
	"context"
	"encoding/json"
	"flag"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/fanart/bilibili"
	"github.com/JustHumanz/Go-Simp/service/fanart/pixiv"
	"github.com/JustHumanz/Go-Simp/service/fanart/twitter"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/robfig/cron/v3"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

func main() {
	Twitter := flag.Bool("TwitterFanart", false, "Enable twitter fanart module")
	BiliBili := flag.Bool("BiliBiliFanart", false, "Enable bilibili fanart module")
	Lewd := flag.Bool("LewdFanart", false, "Enable lewd fanart module")
	Pixiv := flag.Bool("PixivFanArt", false, "Enable Pixiv fanart module")
	flag.Parse()

	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		configfile config.ConfigFile
		Payload    database.VtubersPayload
	)

	RequestPay := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Fanart",
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

	configfile.InitConf()
	database.Start(configfile)

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, RequestPay)

	if *Twitter {
		twitter.Start(Bot, c, Payload, configfile, *Lewd)
	}

	if *BiliBili {
		bilibili.Start(Bot, c, Payload, configfile)
	}

	if *Pixiv {
		pixiv.Start(Bot, c, Payload, configfile, *Lewd)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "TwitterFanart",
		Enabled: *Twitter,
	})
	if err != nil {
		log.Error(err)
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "BiliBiliFanart",
		Enabled: *BiliBili,
	})
	if err != nil {
		log.Error(err)
	}

	c.AddFunc(config.CheckPayload, RequestPay)

	go pilot.RunHeartBeat(gRCPconn, "Fanart")
	runfunc.Run(Bot)
}
