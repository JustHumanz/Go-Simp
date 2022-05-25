package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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

	ServiceName := func(Name string) string {
		if Name == "checker_youtube" {
			return config.YoutubeCheckerService
		} else if Name == "counter_youtube" {
			return config.YoutubeCounterService
		} else if Name == "live_youtube" {
			return config.YoutubeLiveTrackerService
		} else if Name == "past_youtube" {
			return config.YoutubePastTrackerService
		} else if Name == "space_bilibili" {
			return config.SpaceBiliBiliService
		} else if Name == "live_bilibili" {
			return config.LiveBiliBiliService
		} else if Name == "twitch" {
			return config.TwitchService
		} else if Name == "twitter" {
			return config.TwitterService
		} else if Name == "tbilibili" {
			return config.TBiliBiliService
		} else if Name == "pixiv" {
			return config.PixivService
		}
		return ""
	}

	router.HandleFunc("/get/{service}/units/", func(w http.ResponseWriter, r *http.Request) {
		Service := mux.Vars(r)["service"]
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		Data := GetUnitsMetadata(pilot.S, ServiceName(Service))
		if Data != nil {
			json.NewEncoder(w).Encode(Data)
			w.WriteHeader(http.StatusOK)
		} else {
			json.NewEncoder(w).Encode(nil)
			w.WriteHeader(http.StatusNotFound)
		}
	}).Methods(http.MethodGet)

	//Admin can delete the unit payload in case the unit is killed or dead but the payload still remaining
	router.HandleFunc("/delete/{service}/{uuid}/", func(w http.ResponseWriter, r *http.Request) {
		Service := mux.Vars(r)["service"]
		UUID := mux.Vars(r)["uuid"]
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		err := DeleteUnitsMetadata(pilot.S, ServiceName(Service), UUID)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Error":   true,
				"Message": err.Error(),
			})
			w.WriteHeader(http.StatusBadRequest)
		} else {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"Error":   false,
				"Message": nil,
			})
			w.WriteHeader(http.StatusOK)
		}
	}).Methods(http.MethodDelete)

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

func DeleteUnitsMetadata(s pilot.Server, Service, UUID string) error {
	log.WithFields(log.Fields{
		"Service": Service,
		"UUID":    UUID,
	}).Warn("Deleting unit payload")

	for _, v := range s.Service {
		if v.Name == Service {
			if len(v.Unit) > 0 {
				for _, v2 := range v.Unit {
					if v2.UUID == UUID {
						v.RemoveUnitFromDeadNode(v2.UUID)
						fmt.Println(v)
						return nil
					}
				}
			}
		}
	}
	return errors.New("invalid uuid")
}
