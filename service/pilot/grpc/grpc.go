package pilot

import (
	context "context"
	"encoding/json"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/metric"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	VtubersByte *[]byte
	confByte    []byte
	WeebHookURL string
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
		var Grp []database.Group
		GroupData, err := database.GetGroups()
		if err != nil {
			log.Error(err)
		}
		for _, v := range GroupData {
			v.Members, err = database.GetMembers(v.ID)
			if err != nil {
				log.Error(err)
			}
			Grp = append(Grp, v)
		}

		tmp, err := json.Marshal(Grp)
		VtubersByte = &tmp
		if err != nil {
			log.Error(err)
		}
	}
	GetGroups()

	c.AddFunc(config.PilotGetGroups, GetGroups)

	WeebHookURL = configfile.PilotReporting
	confByte, err = json.Marshal(configfile)
	if err != nil {
		log.Error(err)
	}
}

type ModuleList struct {
	Name    string
	Counter int
	CronJob int
	Note    string
	Run     bool
}

type Server struct {
	ServiceList []ServiceMessage
	ModuleData  []ModuleData
	Modules     []*ModuleList
	WaitMigrate *bool
	UnimplementedPilotServiceServer
}

func (s *Server) ReqData(ctx context.Context, in *ServiceMessage) (*VtubersData, error) {
	log.WithFields(log.Fields{
		"Service": in.Service,
		"Message": in.Message,
	}).Info("Request payload")
	for _, v := range s.ServiceList {
		if v.Service != in.Service {
			v.Alive = false
		} else {
			v.Alive = true
		}
	}

	if in.Service == "Migrate" && in.Message == "Start migrate new vtuber" {
		log.Info("Migrate process detect")
		migrate := true
		s.WaitMigrate = &migrate
	} else if in.Message == "Done migrate new vtuber" {
		log.Info("Not detect any vtuber Migrate,Change value to false")
		migrate := false
		s.WaitMigrate = &migrate
	}

	return &VtubersData{
		VtuberPayload: *VtubersByte,
		ConfigFile:    confByte,
		WaitMigrate:   *s.WaitMigrate,
	}, nil
}

func (s *Server) ModuleList(ctx context.Context, in *ModuleData) (*Empty, error) {
	log.WithFields(log.Fields{
		"Module":  in.Module,
		"Enabled": in.Enabled,
	}).Info("Report module list")
	for _, v := range s.ModuleData {
		if v.Module == in.Module {
			v.Enabled = in.Enabled
		}
	}
	return &Empty{}, nil
}

func (i *ModuleList) SetRun(a bool) *ModuleList {
	i.Run = a
	return i
}

func (i *ModuleList) SetNote(n string) *ModuleList {
	i.Note = n
	return i
}

func (s *Server) isYtCheckerRunning() bool {
	for _, v := range s.Modules {
		if v.Name == config.YoutubeCheckerModule {
			if v.Run && v.Note == "Update" {
				return true
			}
		}
	}
	return false
}

func (s *Server) RunModuleJob(ctx context.Context, in *ServiceMessage) (*RunJob, error) {
	for _, v := range s.Modules {
		if v.Name == in.Service {
			if in.Service == config.YoutubeCheckerModule {
				v.SetNote(in.Message)
			}

			if in.Service == config.YoutubeCounterModule && s.isYtCheckerRunning() {
				log.WithFields(log.Fields{
					"Counter": v.CronJob,
					"Service": in.Service,
				}).Warn("Youtube Checker still running,skip Youtube Counter")

				return &RunJob{
					Run:     false,
					Message: "skip Youtube Counter",
				}, nil
			}

			if in.Alive && v.Counter%v.CronJob == 0 {
				log.WithFields(log.Fields{
					"Counter": v.CronJob,
					"Service": in.Service,
					"Running": v.Run,
					"Note":    in.Message,
				}).Info("Module request for running job")

				if v.Run {
					log.WithFields(log.Fields{
						"Counter": v.CronJob,
						"Service": in.Service,
					}).Warn("Job still running")

					return &RunJob{
						Run:     false,
						Message: "job still running",
					}, nil
				}

				log.WithFields(log.Fields{
					"Counter": v.CronJob,
					"Service": in.Service,
				}).Info("Approval request")

				v.SetRun(true)
				return &RunJob{
					Message: "OK,module approved",
					Run:     true,
				}, nil

			} else if !in.Alive && in.Message == "Done" {
				v.SetRun(false)
				log.WithFields(log.Fields{
					"Counter": v.Counter,
					"Service": in.Service,
					"Running": v.Run,
					"Note":    v.Note,
				}).Info("Module complete running job")
			}
			v.Counter++
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
		if in.Alive {
			log.WithFields(log.Fields{
				"Service":  in.Service,
				"Messsage": in.Message,
				"Status":   "Running",
			}).Debug("HeartBeat")
		}

		err := stream.Send(&Empty{})
		if err != nil {
			log.Error(in.Service + " Dead")
			ReportDeadService(err.Error(), in.Service)
			return err
		}
		time.Sleep(5 * time.Second)
	}
}

func RunHeartBeat(client PilotServiceClient, Service string) {
	_, err := client.HeartBeat(context.Background(), &ServiceMessage{
		Service: Service,
		Message: "Service 200 daijoubu",
		Alive:   true,
	})
	if err != nil {
		ReportDeadService("Pilot down", Service)
		log.Fatal(err)
	}
}

func ReportDeadService(message, module string) {
	log.Error(message, module)
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
