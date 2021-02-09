package main

import (
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	pilot.Start()
	lis := network.InitNet()
	defmigrate := false
	s := pilot.Server{
		WaitMigrate: &defmigrate,
		ServiceList: []pilot.ServiceMessage{
			pilot.ServiceMessage{
				Service: "Fanart",
			},
			pilot.ServiceMessage{
				Service: "Livestream",
			},
			pilot.ServiceMessage{
				Service: "Guild",
			},
			pilot.ServiceMessage{
				Service: "Subscriber",
			},
			pilot.ServiceMessage{
				Service: "Utility",
			},
			pilot.ServiceMessage{
				Service: "Rest_API",
			},
		},
		ModuleData: []pilot.ModuleData{
			pilot.ModuleData{
				Module: "TwitterFanart",
			},
			pilot.ModuleData{
				Module: "BiliBiliFanart",
			},
			pilot.ModuleData{
				Module: "Youtube",
			},
			pilot.ModuleData{
				Module: "SpaceBiliBili",
			},
			pilot.ModuleData{
				Module: "LiveBiliBili",
			},
			pilot.ModuleData{
				Module: "Twitch",
			},
			pilot.ModuleData{
				Module: "YoutubeSubscriber",
			},
			pilot.ModuleData{
				Module: "BiliBiliFollowers",
			},
			pilot.ModuleData{
				Module: "TwitterFollowers",
			},
		},
	}

	grpcServer := grpc.NewServer()

	pilot.RegisterPilotServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
