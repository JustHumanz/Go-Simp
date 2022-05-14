package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/JustHumanz/Go-Simp/pkg/config"
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
	s := pilot.Server{
		Service: []*pilot.Service{
			//Fanart
			{
				Name: config.TBiliBiliService,
				//Counter: 1,
				CronJob: 7, //every 7 minutes
			},
			{
				Name: config.TwitterService,
				//Counter: 1,
				CronJob: 10, //every 10 minutes
			},
			{
				Name: config.PixivService,
				//Counter: 1,
				CronJob: 7, //every 7 minutes
			},

			//Live
			{
				Name: config.SpaceBiliBiliService,
				//Counter: 1,
				CronJob: 12, //every 12 minutes
			},
			{
				Name: config.LiveBiliBiliService,
				//Counter: 1,
				CronJob: 7, //every 7 minutes
			},
			{
				Name: config.TwitchService,
				//Counter: 1,
				CronJob: 10, //every 10 minutes
			},
			{
				Name: config.YoutubeCheckerService,
				//Counter: 1,
				CronJob: 5, //every 5 minutes
			},
			{
				Name: config.YoutubeCounterService,
				//Counter: 1,
				CronJob: 1, //every 1 minutes
			},
		},
	}

	grpcServer := grpc.NewServer()

	http.HandleFunc("/youtube/units/", func(w http.ResponseWriter, r *http.Request) {
		Data := GetUnitsMetadata(s, config.YoutubeCheckerService)
		if Data != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Data)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	http.HandleFunc("/twitter/units/", func(w http.ResponseWriter, r *http.Request) {
		Data := GetUnitsMetadata(s, config.TwitterService)
		if Data != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Data)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	go http.ListenAndServe(":8181", nil)

	pilot.RegisterPilotServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func GetUnitsMetadata(s pilot.Server, name string) []pilot.UnitMetadata {
	for _, v := range s.Service {
		if v.Name == name {
			if len(v.Unit) > 0 {
				Data := []pilot.UnitMetadata{}

				for _, v2 := range v.Unit {
					Data = append(Data, v2.Metadata)
				}
				return Data

			}
		}
	}

	return nil
}
