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
	VtubersByte []byte
	BotByte     []byte
	dbByte      []byte
	confByte    []byte
	WeebHookURL string
	Groups      []database.Group
)

func Start() {
	configfile, err := config.ReadConfig("../../config.toml")
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
		for _, v := range database.GetGroups() {
			v.Members = database.GetMembers(v.ID)
			Grp = append(Grp, v)
		}
		Groups = Grp
	}
	GetGroups()

	c.AddFunc(config.PilotGetGroups, GetGroups)

	VtubersByte, err = json.Marshal(database.VtubersPayload{
		VtuberData: Groups,
	})
	if err != nil {
		log.Error(err)
	}

	WeebHookURL = configfile.PilotReporting
	confByte, err = json.Marshal(configfile)
	if err != nil {
		log.Error(err)
	}
}

type Server struct {
	ServiceList []ServiceMessage
	ModuleData  []ModuleData
	WaitMigrate *bool
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
		VtuberPayload: VtubersByte,
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

func (s *Server) MetricReport(ctx context.Context, in *Metric) (*Empty, error) {
	if in.State == config.FanartState {
		var FanArt database.DataFanart
		err := json.Unmarshal(in.MetricData, &FanArt)
		if err != nil {
			log.Error(err)
		}

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

		if Subs.State == config.BiliBiliLive {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"BiliBili",
			).Add(float64(Subs.NewSubs))

			metric.GetViews.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"BiliBili",
			).Add(float64(Subs.NewViews))
		} else if Subs.State == config.YoutubeLive {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Youtube",
			).Add(float64(Subs.NewSubs))

			metric.GetViews.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Youtube",
			).Add(float64(Subs.NewViews))
		} else {
			metric.GetSubs.WithLabelValues(
				Subs.Member.Name,
				Subs.Group.GroupName,
				"Twitter",
			).Add(float64(Subs.NewSubs))
		}
	} else if in.State == config.LiveStatus {
		var LiveData database.LiveStream
		err := json.Unmarshal(in.MetricData, &LiveData)
		if err != nil {
			log.Error(err)
		}
		if LiveData.State == config.YoutubeLive {
			metric.GetSubs.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Youtube",
			).Inc()
		} else if LiveData.State == config.BiliLive {
			metric.GetSubs.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"BiliBili",
			).Inc()
		} else if LiveData.State == config.TwitchLive {
			metric.GetSubs.WithLabelValues(
				LiveData.Member.Name,
				LiveData.Group.GroupName,
				"Twitch",
			).Inc()
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
			}).Info("HeartBeat")
		}

		err := stream.Send(&Empty{})
		if err != nil {
			ReportDeadService(err.Error())
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
		ReportDeadService("Pilot down")
		log.Fatal(err)
	}
}

func ReportDeadService(message string) {
	log.Error(message)
	PayloadBytes, err := json.Marshal(map[string]interface{}{
		"embeds": []interface{}{
			map[string]interface{}{
				"description": "Failed to send HeartBeat",
				"fields": []interface{}{
					map[string]interface{}{
						"name":   "Time down",
						"value":  time.Now(),
						"inline": true,
					},
					map[string]interface{}{
						"name":   "Error message",
						"value":  message,
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
