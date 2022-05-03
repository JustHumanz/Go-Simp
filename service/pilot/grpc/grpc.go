package pilot

import (
	context "context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/metric"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/network"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	confByte      []byte
	WeebHookURL   string
	VtubersAgency []database.Group
	//Payload       = make(map[string]interface{})
)

func Start() {
	configfile, err := config.ReadConfig("./config.toml")
	if err != nil {
		log.Panic(err)
	}

	c := cron.New()
	c.Start()
	configfile.InitConf()
	database.Start(configfile)

	GetGroups := func() {
		log.Info("Start get groups from database")
		VtubersAgency, err = database.GetGroups()
		if err != nil {
			log.Error(err)
		}

		/*
			for _, v := range GroupData {
				v.Members, err = database.GetMembers(v.ID)
				if err != nil {
					log.Error(err)
				}
				VtubersAgency = append(VtubersAgency, v)
			}


				tmp, err := json.Marshal(Grp)
				VtubersByte = &tmp
				if err != nil {
					log.Error(err)
				}
		*/
	}
	GetGroups()

	c.AddFunc(config.PilotGetGroups, GetGroups)

	WeebHookURL = configfile.PilotReporting
	confByte, err = json.Marshal(configfile)
	if err != nil {
		log.Error(err)
	}
}

type Service struct {
	Name    string
	Unit    []*UnitService
	CronJob int
	Note    string
	Run     bool
}

type UnitService struct {
	UUID    string
	Counter int
	Payload []database.Group
}

type Server struct {
	ServiceBootstrap []ServiceMessage
	Service          []*Service
	UnimplementedPilotServiceServer
}

func (s *Server) GetAgencyPayload(ctx context.Context, in *ServiceMessage) (*AgencyPayload, error) {
	log.WithFields(log.Fields{
		"Service": in.Service,
		"Message": in.Message,
		"UUID":    in.ServiceUUID,
	}).Info("Request agency payload")

	BytePayload, err := json.Marshal(VtubersAgency)
	if err != nil {
		return nil, err
	}
	return &AgencyPayload{
		AgencyVtubers: BytePayload,
	}, nil
}

func (s *Server) GetBotPayload(ctx context.Context, in *ServiceMessage) (*ServiceInit, error) {
	log.WithFields(log.Fields{
		"Service": in.Service,
		"Message": in.Message,
		"UUID":    in.ServiceUUID,
	}).Info("Request bot payload")

	return &ServiceInit{
		ConfigFile: confByte,
	}, nil
}

func (i *Service) SetRun(a bool) *Service {
	i.Run = a
	return i
}

func (i *Service) SetNote(n string) *Service {
	i.Note = n
	return i
}

func (s *Server) IsYtCheckerRunning() bool {
	for _, v := range s.Service {
		if v.Name == config.YoutubeCheckerService {
			if v.Run && v.Note == "Update" {
				return true
			}
		}
	}
	return false
}

//var ModuleWatcher = make(map[string]int)

func (s *Server) CheckUUID(UUID string) bool {
	for _, v := range s.Service {
		for _, v2 := range v.Unit {
			if UUID == v2.UUID {
				return true
			}
		}
	}
	return false
}

func (s *Service) RemapPayload() {
	if len(s.Unit) > 0 {
		limt := len(VtubersAgency) / len(s.Unit)
		chunk := [][]database.Group{}

		for i := 0; i < len(VtubersAgency); i += limt {
			nextLimit := limt + i
			if nextLimit+i > len(VtubersAgency) {
				nextLimit = len(VtubersAgency)
			}

			chunk = append(chunk, VtubersAgency[i:nextLimit])
			if VtubersAgency[nextLimit-1].ID == VtubersAgency[len(VtubersAgency)-1].ID {
				break
			}
		}

		for k, v := range chunk {
			s.Unit[k].Payload = v
		}
	}
}

func (k *UnitService) Marshal() []byte {

	dat, err := json.Marshal(k.Payload)
	if err != nil {
		log.Error(err)
		return nil
	}
	return dat
}

func (k *UnitService) ResetCounter() *UnitService {
	k.Counter = 1
	return k
}

func (s *Service) RemoveUnitFromDeadSvc(UUID string) {
	for k, v := range s.Unit {
		if v.UUID == UUID {
			s.Unit = append(s.Unit[0:k], s.Unit[k+1:]...)
		}
	}

	s.RemapPayload()
}

func (s *Server) RequestRunJobsOfService(ctx context.Context, in *ServiceMessage) (*RunJob, error) {
	for _, v := range s.Service {
		if v.Name == in.Service {
			//Check if request was new unit and if that unit was not scraping
			if !s.CheckUUID(in.ServiceUUID) {
				//New units detec
				log.WithFields(log.Fields{
					"Service": in.Service,
					"UUID":    in.ServiceUUID,
				}).Info("New Units detected,trying to register new unit")

				v.Unit = append(v.Unit, &UnitService{
					in.ServiceUUID,
					1,
					nil,
				})

				v.RemapPayload()

				payload := []string{}
				for _, v2 := range v.Unit {
					if v2.UUID == in.ServiceUUID {

						for _, v3 := range v2.Payload {
							payload = append(payload, v3.GroupName)
						}

						log.WithFields(log.Fields{
							"Agency Payload": payload,
							"UUID":           v2.UUID,
						})
					}
				}
				return &RunJob{
					Run:     false,
					Message: fmt.Sprintln("New units detected %s", payload),
				}, nil
			}

			for _, k := range v.Unit {
				if k.UUID == in.ServiceUUID {
					if in.Service == config.YoutubeCheckerService {
						v.SetNote(in.Message)
					}

					if in.Service == config.YoutubeCounterService && s.IsYtCheckerRunning() {
						log.WithFields(log.Fields{
							"Counter": v.CronJob,
							"UUID":    in.ServiceUUID,
							"Service": in.Service,
						}).Warn("Youtube Checker still running,skip Youtube Counter")

						return &RunJob{
							Run:     false,
							Message: "skip Youtube Counter",
						}, nil
					}

					if in.Message == "Done" {
						v.SetRun(false)
						k.ResetCounter()
						log.WithFields(log.Fields{
							"Counter": k.Counter,
							"UUID":    in.ServiceUUID,
							"Service": in.Service,
							"Running": v.Run,
							"Note":    v.Note,
						}).Info("Unit complete running job")
						return &RunJob{
							Run:     false,
							Message: fmt.Sprintf("Waiting schdule %d minutes remaining", v.CronJob-k.Counter),
						}, nil
					}

					if k.Counter%v.CronJob == 0 {
						log.WithFields(log.Fields{
							"Counter": v.CronJob,
							"UUID":    in.ServiceUUID,
							"Service": in.Service,
							"Running": v.Run,
							"Note":    in.Message,
						}).Info("Unit request for running job")

						log.WithFields(log.Fields{
							"Counter": v.CronJob,
							"UUID":    in.ServiceUUID,
							"Service": in.Service,
						}).Info("Approval request")

						v.SetRun(true)
						return &RunJob{
							Message:       "OK,Unit approved",
							Run:           true,
							VtuberPayload: k.Marshal(),
						}, nil

					} else {
						k.Counter++
						log.WithFields(log.Fields{
							"Counter": k.Counter,
							"Service": in.Service,
							"UUID":    in.ServiceUUID,
						}).Info("Waiting Request payload service")
						return &RunJob{
							Run:     false,
							Message: fmt.Sprintf("Waiting schdule %d minutes remaining", v.CronJob-k.Counter),
						}, nil
					}
				}
			}
		}
	}
	return &RunJob{}, nil
}

func (s *Server) ReportError(ctx context.Context, in *ServiceMessage) (*Empty, error) {
	ReportDeadService(in.Message, in.Service)
	return &Empty{}, nil
}

func (s *Server) MetricReport(ctx context.Context, in *Metric) (*Empty, error) {
	if in.State == config.FanartState {
		var FanArt database.DataFanart
		err := json.Unmarshal(in.MetricData, &FanArt)
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"Vtuber": FanArt.Member.EnName,
			"State":  in.State,
		}).Info("Update Fanart metric")

		metric.GetFanArt.WithLabelValues(
			FanArt.Member.Name,
			FanArt.Group.GroupName,
			FanArt.Author,
			BoolString(FanArt.Lewd),
			FanArt.State,
		).Inc()
	} else if in.State == config.SubsState {
		var Subs database.MemberSubs
		err := json.Unmarshal(in.MetricData, &Subs)
		if err != nil {
			log.Error(err)
		}
		log.WithFields(log.Fields{
			"Vtuber": Subs.Member.EnName,
			"State":  Subs.State,
		}).Info("Update Subs")
		if Subs.State == config.BiliLive && !Subs.Member.IsBiliNill() {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"BiliBili",
			).Set(float64(Subs.BiliFollow))

			metric.GetViews.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"BiliBili",
			).Set(float64(Subs.BiliViews))
		} else if Subs.State == config.YoutubeLive && !Subs.Member.IsYtNill() {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Youtube",
			).Set(float64(Subs.YtSubs))

			metric.GetViews.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Youtube",
			).Set(float64(Subs.YtViews))
		} else if Subs.State == config.TwitchLive && !Subs.Member.IsTwitchNill() {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Twitch",
			).Set(float64(Subs.TwitchFollow))

			metric.GetViews.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Twitch",
			).Set(float64(Subs.TwitchViews))
		} else {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Twitter",
			).Set(float64(Subs.TwFollow))
		}
	} else if in.State == config.LiveStatus {
		var LiveData database.LiveStream
		err := json.Unmarshal(in.MetricData, &LiveData)
		if err != nil {
			log.Error(err)
		}

		log.WithFields(log.Fields{
			"Vtuber": LiveData.Member.EnName,
			"State":  in.State,
		}).Info("Update Livestream metric")

		if LiveData.State == config.YoutubeLive && !LiveData.Member.IsYtNill() {
			metric.GetLive.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Youtube",
			).Inc()
		} else if LiveData.State == config.BiliLive && !LiveData.Member.IsBiliNill() {
			metric.GetLive.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"BiliBili",
			).Inc()
		} else if LiveData.State == config.TwitchLive {
			metric.GetLive.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Twitch",
			).Inc()
		}
	} else if in.State == config.PastStatus {
		var LiveData database.LiveStream
		err := json.Unmarshal(in.MetricData, &LiveData)
		if err != nil {
			log.Error(err)
		}

		if LiveData.End.IsZero() {
			return &Empty{}, nil
		}

		Time := LiveData.End.Sub(LiveData.Schedul).Minutes()

		log.WithFields(log.Fields{
			"Vtuber":  LiveData.Member.EnName,
			"State":   in.State,
			"Time":    int(Time),
			"VideoID": LiveData.VideoID,
		}).Info("Update past Livestream metric")

		if LiveData.State == config.YoutubeLive && !LiveData.Member.IsYtNill() {
			metric.GetLiveDuration.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Youtube",
			).Add(Time)
		} else if LiveData.State == config.BiliLive && !LiveData.Member.IsBiliNill() {
			metric.GetLiveDuration.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"BiliBili",
			).Add(Time)
		} else if LiveData.State == config.TwitchLive {
			metric.GetLiveDuration.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Twitch",
			).Add(Time)
		}
	}

	return &Empty{}, nil
}

func (s *Server) HeartBeat(in *ServiceMessage, stream PilotService_HeartBeatServer) error {
	for {
		for _, v := range s.Service {
			if v.Name == in.Service {
				for _, v2 := range v.Unit {
					if v2.UUID == in.ServiceUUID {
						log.WithFields(log.Fields{
							"Service":  in.Service,
							"Messsage": in.Message,
							"UUID":     in.ServiceUUID,
							"Status":   "Running",
						}).Debug("HeartBeat")

						err := stream.Send(&Empty{})
						if err != nil {
							log.WithFields(log.Fields{
								"Service":  in.Service,
								"Messsage": in.Message,
								"UUID":     in.ServiceUUID,
							}).Error(fmt.Sprintf("%s Removing unit if exist", err))
							v.RemoveUnitFromDeadSvc(in.ServiceUUID)
							ReportDeadService(err.Error(), in.Service)
							return err
						}
						time.Sleep(5 * time.Second)
					}
				}
			}
		}
	}
}

func RunHeartBeat(client PilotServiceClient, Service string, UUID string) {
	_, err := client.HeartBeat(context.Background(), &ServiceMessage{
		Service:     Service,
		Message:     "Bep Bob",
		ServiceUUID: UUID,
	})
	if err != nil {
		ReportDeadService("Pilot down", Service)
		log.Fatal(err)
	}
}

func ReportDeadService(message, module string) {
	PayloadBytes, err := json.Marshal(map[string]interface{}{
		"embeds": []interface{}{
			map[string]interface{}{
				"description": "Failed to send HeartBeat",
				"fields": []interface{}{
					map[string]interface{}{
						"name":   "Time down",
						"value":  time.Now().Format(time.RFC822),
						"inline": true,
					},
					map[string]interface{}{
						"name":   "Error message",
						"value":  message,
						"inline": true,
					},
					map[string]interface{}{
						"name":   "Module",
						"value":  module,
						"inline": true,
					},
				},
			},
		},
	})
	if err != nil {
		log.Error(err)
	}
	err = network.CurlPost(WeebHookURL, PayloadBytes)
	if err != nil {
		log.Error(err)
	}
}

func BoolString(a bool) string {
	if a {
		return "true"
	}
	return "false"
}
