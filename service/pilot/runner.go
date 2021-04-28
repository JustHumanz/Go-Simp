package main

import (
	"flag"
	"net/http"

	"github.com/JustHumanz/Go-Simp/pkg/metric"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	pilot.Start()

	var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(*addr, nil)

	metric.Init()
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
