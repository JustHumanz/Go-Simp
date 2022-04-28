package main

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

var (
	Bot          *discordgo.Session
	GroupPayload *[]database.Group
	gRCPconn     = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	ServiceUUID  = uuid.New().String()
)

const (
	ServiceName = config.SpaceBiliBiliService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

//Start start twitter module
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
	Bot = engine.StartBot(false)

	database.Start(configfile)

	log.Info("Enable " + ServiceName)
	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go ReqRunningJob(gRCPconn)
	runfunc.Run(Bot)

}

type checkBlSpaceJob struct {
	wg         sync.WaitGroup
	mutex      sync.Mutex
	Reverse    bool
	VideoIDTMP map[string]string
	agency     []database.Group
}

func (i *checkBlSpaceJob) AddVideoID(Member, VideoID string) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.VideoIDTMP[Member] = VideoID
}

func (i *checkBlSpaceJob) CekFirstVideoID(Member, VideoID string) bool {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if i.VideoIDTMP[Member] == VideoID {
		log.WithFields(log.Fields{
			"Vtuber": Member,
		}).Warn("Video Still same")
		return true
	}
	return false
}

func ReqRunningJob(client pilot.PilotServiceClient) {
	SpaceBiliBili := &checkBlSpaceJob{
		VideoIDTMP: make(map[string]string),
	}

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
			}).Info(res.Message)

			SpaceBiliBili.agency = engine.UnMarshalPayload(res.VtuberPayload)
			SpaceBiliBili.Run()

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

		}

		time.Sleep(1 * time.Minute)
	}
}

func (i *checkBlSpaceJob) Run() {
	Cek := func(Group database.Group, wg *sync.WaitGroup) {
		defer wg.Done()
		for _, Member := range Group.Members {
			if Member.BiliBiliID != 0 && Member.Active() {
				log.WithFields(log.Fields{
					"Agency": Group.GroupName,
					"Vtuber": Member.Name,
				}).Info("Checking Space BiliBili")

				Group.RemoveNillIconURL()
				Data := &database.LiveStream{
					Member: Member,
					Group:  Group,
				}
				var (
					Videotype string
					PushVideo engine.SpaceVideo
				)

				body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/arc/search?mid="+strconv.Itoa(Data.Member.BiliBiliID)+"&ps="+strconv.Itoa(config.GoSimpConf.LimitConf.SpaceBiliBili), nil)
				if curlerr != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(curlerr)
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message: curlerr.Error(),
						Service: ServiceName,
					})
					return
				}

				err := json.Unmarshal(body, &PushVideo)
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					}).Error(err)
				}

				if len(PushVideo.Data.List.Vlist) > 0 {
					FirstVideoID := PushVideo.Data.List.Vlist[0].Bvid

					if i.CekFirstVideoID(Member.Name, FirstVideoID) {
						continue
					}

					for _, video := range PushVideo.Data.List.Vlist {
						if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|翻唱|mv|歌曲)", strings.ToLower(video.Title)); Cover || video.Typeid == 31 {
							Videotype = "Covering"
						} else {
							Videotype = "Streaming"
						}

						Data.AddVideoID(video.Bvid).SetType(Videotype).
							UpdateTitle(video.Title).
							UpdateThumbnail(video.Pic).UpdateSchdule(time.Unix(int64(video.Created), 0)).
							UpdateViewers(strconv.Itoa(video.Play)).UpdateLength(video.Length).SetState(config.SpaceBili)
						new, id := Data.CheckVideo()
						if new {
							log.WithFields(log.Fields{
								"Vtuber": Data.Member.Name,
							}).Info("New video uploaded")

							Data.InputSpaceVideo()
							video.VideoType = Videotype
							engine.SendLiveNotif(Data, Bot)
						} else {
							Data.UpdateSpaceViews(id)
						}
					}
					i.AddVideoID(Member.Name, FirstVideoID)
				}
			}
		}
	}

	if i.Reverse {
		for j := len(*GroupPayload) - 1; j >= 0; j-- {
			i.wg.Add(1)
			Grp := *GroupPayload
			go Cek(Grp[j], &i.wg)
		}
		i.Reverse = false

	} else {
		for _, G := range *GroupPayload {
			i.wg.Add(1)
			go Cek(G, &i.wg)
		}
		i.Reverse = true
	}
	i.wg.Wait()
}
