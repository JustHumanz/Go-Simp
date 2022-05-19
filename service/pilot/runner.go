package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/metric"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	muxlogrus "github.com/pytimer/mux-logrus"
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
				CronJob: 15, //every 15 minutes
			},
			{
				Name: config.YoutubeCounterService,
				//Counter: 1,
				CronJob: 1, //every 1 minutes
			},
		},
	}

	grpcServer := grpc.NewServer()
	router := mux.NewRouter()

	router.HandleFunc("/{service}/units/", func(w http.ResponseWriter, r *http.Request) {
		Service := mux.Vars(r)["service"]
		Data := []pilot.UnitMetadata{}

		if Service == "checker_youtube" {
			Data = GetUnitsMetadata(s, config.YoutubeCheckerService)
		} else if Service == "counter_youtube" {
			Data = GetUnitsMetadata(s, config.YoutubeCounterService)
		} else if Service == "space_bilibili" {
			Data = GetUnitsMetadata(s, config.SpaceBiliBiliService)
		} else if Service == "live_bilibili" {
			Data = GetUnitsMetadata(s, config.LiveBiliBiliService)
		} else if Service == "twitch" {
			Data = GetUnitsMetadata(s, config.TwitchService)
		} else if Service == "twitter" {
			Data = GetUnitsMetadata(s, config.TwitterService)
		} else if Service == "tbilibili" {
			Data = GetUnitsMetadata(s, config.TBiliBiliService)
		} else if Service == "pixiv" {
			Data = GetUnitsMetadata(s, config.PixivService)
		}

		if Data != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Data)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})

	router.Use(muxlogrus.NewLogger().Middleware)
	go http.ListenAndServe(":8181", engine.LowerCaseURI(router))

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
