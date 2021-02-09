package pilot

import (
	context "context"
	"encoding/json"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	log "github.com/sirupsen/logrus"
)

var (
	VtubersByte []byte
	BotByte     []byte
	dbByte      []byte
	confByte    []byte
	WeebHookURL string
)

func Start() {
	configfile, err := config.ReadConfig("../../config.toml")
	if err != nil {
		log.Panic(err)
	}

	configfile.InitConf()
	database.Start(configfile)

	Groups := database.GetGroups()
	for i, _ := range Groups {
		Groups[i].Members = database.GetMembers(Groups[i].ID)
	}

	/*
			Bot, err := discordgo.New("Bot " + config.GoSimpConf.Discord)
			if err != nil {
				log.Error(err)
			}


		BotByte, err = json.Marshal(Bot)
		if err != nil {
			log.Error(err)
			BotByte = nil
		}
	*/

	VtubersByte, err = json.Marshal(database.VtubersPayload{
		VtuberData: Groups,
	})
	if err != nil {
		log.Error(err)
	}

	WeebHookURL = configfile.DiscordWebHook
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

func (s *Server) HeartBeat(stream PilotService_HeartBeatServer) error {
	for {
		err := stream.Send(&Empty{})
		if err != nil {
			/*
				Service die,Send notif via discord webhook
			*/
			ReportDeadService(err.Error())
			return err
		}
		client, err := stream.Recv()
		if err != nil {
			ReportDeadService(err.Error())
			return err
		}

		if client.Alive {
			log.WithFields(log.Fields{
				"Service":  client.Service,
				"Messsage": client.Message,
				"Status":   "Running",
			}).Info("HeartBeat")
		}
		time.Sleep(2 * time.Second)
	}
	return nil
}

func RunHeartBeat(client PilotServiceClient, Service string) {
	for {
		stream, err := client.HeartBeat(context.Background())
		if err != nil {
			log.Fatalf("%v.RecordRoute(_) = _, %v", client, err)
		}
		err = stream.Send(&ServiceMessage{
			Service: Service,
			Message: "Service 200 daijoubu",
			Alive:   true,
		})
		if err != nil {
			log.Error(err)
		}
		time.Sleep(5 * time.Second)
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
