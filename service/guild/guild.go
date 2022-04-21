package main

import (
	"context"
	"encoding/json"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		configfile config.ConfigFile
	)
	res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
		Message: "Send me nude",
		Service: "Guild",
	})
	if err != nil {
		log.Fatalf("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}

	var GroupsPayload *[]database.Group
	GetPayload := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Guild",
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error(), "Guild")
			}
			log.Error("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &GroupsPayload)
		if err != nil {
			log.Error(err)
		}
	}
	GetPayload()
	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, GetPayload)

	Bot, err := discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	err = Bot.Open()
	if err != nil {
		log.Error(err)
	}

	configfile.InitConf()

	Bot.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		if g.Unavailable {
			log.Error("joined unavailable guild", g.Guild.ID)
			return
		}

		Join, err := g.JoinedAt.Parse()
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"GuildName": g.Name,
			"OwnerID":   g.OwnerID,
			"JoinDate":  Join.Format(time.RFC822),
		}).Info("New invite")
		engine.InitSlash(Bot, *GroupsPayload, g.Guild)

	})
	log.Info("Guild handler ready.......")
	go pilot.RunHeartBeat(gRCPconn, "Guild")
	runfunc.Run(Bot)
}
