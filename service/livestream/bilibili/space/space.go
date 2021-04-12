package main

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

var (
	loc         *time.Location
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
)

const (
	ModuleState = "SpaceBiliBili"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
}

//Start start twitter module
func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
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
				pilot.ReportDeadService(err.Error())
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
	c.AddFunc(config.BiliBiliSpace, CheckSpaceVideo)
	log.Info("Enable " + ModuleState)
	go pilot.RunHeartBeat(gRCPconn, ModuleState)
	runfunc.Run(Bot)
}

//CheckSpaceVideo
func CheckSpaceVideo() {
	for _, GroupData := range VtubersData.VtuberData {
		if GroupData.GroupName != "Hololive" {
			wg := new(sync.WaitGroup)
			for i, MemberData := range GroupData.Members {
				wg.Add(1)
				go func(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
					defer wg.Done()
					if Member.BiliBiliID != 0 {
						log.WithFields(log.Fields{
							"Group":      Group.GroupName,
							"Vtuber":     Member.EnName,
							"BiliBiliID": Member.BiliBiliID,
						}).Info("Checking Space BiliBili")

						Group.RemoveNillIconURL()
						SpaceBiliLimit := strconv.Itoa(config.GoSimpConf.LimitConf.SpaceBiliBili)
						Data := &database.LiveStream{
							Member: Member,
							Group:  Group,
						}
						CheckSpace(Data, SpaceBiliLimit)

					}
				}(GroupData, MemberData, wg)
				if i%config.Waiting == 0 {
					wg.Wait()
				}
			}
			wg.Wait()
		}
	}
}
