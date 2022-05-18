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

	grpcServer := grpc.NewServer()
	router := mux.NewRouter()

	router.HandleFunc("/{service}/units/", func(w http.ResponseWriter, r *http.Request) {
		Service := mux.Vars(r)["service"]
		Data := []pilot.UnitMetadata{}

		if Service == "checker_youtube" {
			Data = GetUnitsMetadata(pilot.S, config.YoutubeCheckerService)
		} else if Service == "counter_youtube" {
			Data = GetUnitsMetadata(pilot.S, config.YoutubeCounterService)
		} else if Service == "space_bilibili" {
			Data = GetUnitsMetadata(pilot.S, config.SpaceBiliBiliService)
		} else if Service == "live_bilibili" {
			Data = GetUnitsMetadata(pilot.S, config.LiveBiliBiliService)
		} else if Service == "twitch" {
			Data = GetUnitsMetadata(pilot.S, config.TwitchService)
		} else if Service == "twitter" {
			Data = GetUnitsMetadata(pilot.S, config.TwitterService)
		} else if Service == "tbilibili" {
			Data = GetUnitsMetadata(pilot.S, config.TBiliBiliService)
		} else if Service == "pixiv" {
			Data = GetUnitsMetadata(pilot.S, config.PixivService)
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

	pilot.RegisterPilotServiceServer(grpcServer, &pilot.S)

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
