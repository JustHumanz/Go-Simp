package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	Payload database.VtubersPayload
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	var (
		configfile config.ConfigFile
	)

	res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
		Message: "Send me nude",
		Service: "Rest_API",
	})
	if err != nil {
		log.Fatalf("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(res.VtuberPayload, &Payload)
	if err != nil {
		log.Panic(err)
	}

	config.GoSimpConf = configfile
	database.Start(configfile)
	go pilot.RunHeartBeat(gRCPconn, "Rest_API")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/All", getGroup).Methods("GET")
	router.HandleFunc("/Groups/{GroupID}", getGroup).Methods("GET")
	router.HandleFunc("/Members/{MemberID}", getMembers).Methods("GET")

	router.HandleFunc("/Youtube/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/Youtube/Group/{GroupID}/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/Youtube/Member/{MemberID}/{Status}", getYoutube).Methods("GET")

	router.HandleFunc("/Bilibili/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/Bilibili/Group/{GroupID}/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/Bilibili/Member/{MemberID}/{Status}", getBilibili).Methods("GET")

	router.HandleFunc("/Twitter", getFanart).Methods("GET")
	router.HandleFunc("/Twitter/Group/{GroupID}", getFanart).Methods("GET")
	router.HandleFunc("/Twitter/Member/{MemberID}", getFanart).Methods("GET")

	router.HandleFunc("/Tbilibili", getFanart).Methods("GET")
	router.HandleFunc("/Tbilibili/Group/{GroupID}", getFanart).Methods("GET")
	router.HandleFunc("/Tbilibili/Member/{MemberID}", getFanart).Methods("GET")

	router.HandleFunc("/Subscriber", getSubs).Methods("GET")
	router.HandleFunc("/Subscriber/Group/{GroupID}", getSubs).Methods("GET")
	router.HandleFunc("/Subscriber/Member/{MemberID}", getSubs).Methods("GET")

	http.ListenAndServe(":2525", router)
}

func getFanart(w http.ResponseWriter, r *http.Request) {
	var (
		Vars    = mux.Vars(r)
		Fanart  []database.DataFanart
		Twitter = strings.HasPrefix(r.URL.String(), "/Twitter")
	)

	if Vars["GroupID"] != "" {
		key := strings.Split(Vars["GroupID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, GroupIDstr := range key {
				GroupIDint, err := strconv.Atoi(GroupIDstr)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDint == int(Group.ID) {
					State := ""
					if Twitter {
						State = "Twitter"
					} else {
						State = "Tbilibili"
					}
					Fanart = append(Fanart, GetFanartData(State, Group.ID, 0)...)
				}
			}
		}
	} else if Vars["MemberID"] != "" {
		key := strings.Split(Vars["MemberID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if MemberIDint == int(Member.ID) {
						State := ""
						if Twitter {
							State = "Twitter"
						} else {
							State = "Tbilibili"
						}
						Fanart = append(Fanart, GetFanartData(State, 0, Member.ID)...)
					}
				}
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			State := ""
			if Twitter {
				State = "Twitter"
			} else {
				State = "Tbilibili"
			}
			Fanart = append(Fanart, GetFanartData(State, Group.ID, 0)...)
		}
	}

	if len(Fanart) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Fanart)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

func getYoutube(w http.ResponseWriter, r *http.Request) {
	var (
		YoutubeData []database.YtDbData
		Vars        = mux.Vars(r)
		Region      = strings.ToLower(r.FormValue("region"))
		Status      = strings.ToLower(Vars["Status"])
	)
	if Status == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Status notfound",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if Vars["GroupID"] != "" {
		key := strings.Split(Vars["GroupID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, GroupIDstr := range key {
				GroupIDint, err := strconv.Atoi(GroupIDstr)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDint == int(Group.ID) {
					YTData, err := database.YtGetStatus(Group.ID, 0, Status, Region)
					if err != nil {
						log.Error(err)
					}
					YoutubeData = append(YoutubeData, YTData...)
				}
			}
		}
	} else if Vars["MemberID"] != "" {
		key := strings.Split(Vars["MemberID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if MemberIDint == int(Member.ID) {
						YTData, err := database.YtGetStatus(0, Member.ID, Status, Region)
						if err != nil {
							log.Error(err)
						}
						YoutubeData = append(YoutubeData, YTData...)
					}
				}
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			YTData, err := database.YtGetStatus(Group.ID, 0, Status, Region)
			if err != nil {
				log.Error(err)
			}
			YoutubeData = append(YoutubeData, YTData...)
		}
	}

	if len(YoutubeData) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(YoutubeData)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

func getBilibili(w http.ResponseWriter, r *http.Request) {
	var (
		BiliBiliData []database.LiveBiliDB
		Vars         = mux.Vars(r)
		Status       = strings.ToLower(Vars["Status"])
	)
	if Status == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Status notfound",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if Vars["GroupID"] != "" {
		key := strings.Split(Vars["GroupID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, GroupIDstr := range key {
				GroupIDint, err := strconv.Atoi(GroupIDstr)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDint == int(Group.ID) {
					BiliBiliData = append(BiliBiliData, database.BilGet(Group.ID, 0, Status)...)
				}
			}
		}
	} else if Vars["MemberID"] != "" {
		key := strings.Split(Vars["MemberID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if MemberIDint == int(Member.ID) {
						BiliBiliData = append(BiliBiliData, database.BilGet(0, Member.ID, Status)...)
					}
				}
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			BiliBiliData = append(BiliBiliData, database.BilGet(Group.ID, 0, Status)...)
		}
	}

	if len(BiliBiliData) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BiliBiliData)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["GroupID"]
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []database.Group
		for _, Group := range Payload.VtuberData {
			for _, GroupIDstr := range key {
				GroupIDint, err := strconv.Atoi(GroupIDstr)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDint == int(Group.ID) {
					GroupsTMP = append(GroupsTMP, Group)
				}
			}
		}
		if len(GroupsTMP) > 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(GroupsTMP)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(MessageError{
				Message: "Reqest not found,404",
				Date:    time.Now(),
			})
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Payload.VtuberData)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	idstr := mux.Vars(r)["MemberID"]
	Members := []database.Member{}

	for _, Group := range Payload.VtuberData {
		if idstr != "" {
			key := strings.Split(idstr, ",")
			for _, Member := range Group.Members {
				Member.GroupID = Group.ID
				for _, MemberStr := range key {
					MemberInt, err := strconv.Atoi(MemberStr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if region != "" {
						if MemberInt == int(Member.ID) && strings.ToLower(region) == strings.ToLower(Member.Region) {
							Members = append(Members, Member)
						}
					} else {
						if MemberInt == int(Member.ID) {
							Members = append(Members, Member)
						}
					}
				}
			}
		}
	}
	if len(Members) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Members)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

func getSubs(w http.ResponseWriter, r *http.Request) {
	type SubsJson struct {
		MemberID          int64
		Name              string
		EnName            string
		JpName            string
		YoutubeSubscribe  int
		YoutubeVideos     int
		YoutubeViews      int
		BiliBiliFollowers int
		BiliBiliVideos    int
		BiliBiliViews     int
		TwitterFollowers  int
	}
	var (
		SubsData []SubsJson
		Vars     = mux.Vars(r)
	)
	if Vars["GroupID"] != "" {
		key := strings.Split(Vars["GroupID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, GroupIDstr := range key {
				GroupIDint, err := strconv.Atoi(GroupIDstr)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDint == int(Group.ID) {
					for _, Member := range Group.Members {
						tmp, err := Member.GetSubsCount()
						if err != nil {
							log.Error(err)
						}
						SubsData = append(SubsData, SubsJson{
							MemberID:          tmp.MemberID,
							Name:              Member.Name,
							EnName:            Member.EnName,
							JpName:            Member.JpName,
							YoutubeSubscribe:  tmp.YtSubs,
							YoutubeVideos:     tmp.YtVideos,
							YoutubeViews:      tmp.YtViews,
							BiliBiliFollowers: tmp.BiliFollow,
							BiliBiliVideos:    tmp.BiliVideos,
							BiliBiliViews:     tmp.BiliViews,
							TwitterFollowers:  tmp.TwFollow,
						})
					}
				}
			}
		}
	} else if Vars["MemberID"] != "" {
		key := strings.Split(Vars["MemberID"], ",")
		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if MemberIDint == int(Member.ID) {
						tmp, err := Member.GetSubsCount()
						if err != nil {
							log.Error(err)
						}
						SubsData = append(SubsData, SubsJson{
							MemberID:          tmp.MemberID,
							Name:              Member.Name,
							EnName:            Member.EnName,
							JpName:            Member.JpName,
							YoutubeSubscribe:  tmp.YtSubs,
							YoutubeVideos:     tmp.YtVideos,
							YoutubeViews:      tmp.YtViews,
							BiliBiliFollowers: tmp.BiliFollow,
							BiliBiliVideos:    tmp.BiliVideos,
							BiliBiliViews:     tmp.BiliViews,
							TwitterFollowers:  tmp.TwFollow,
						})
					}
				}
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				tmp, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}
				SubsData = append(SubsData, SubsJson{
					MemberID:          tmp.MemberID,
					Name:              Member.Name,
					EnName:            Member.EnName,
					JpName:            Member.JpName,
					YoutubeSubscribe:  tmp.YtSubs,
					YoutubeVideos:     tmp.YtVideos,
					YoutubeViews:      tmp.YtViews,
					BiliBiliFollowers: tmp.BiliFollow,
					BiliBiliVideos:    tmp.BiliVideos,
					BiliBiliViews:     tmp.BiliViews,
					TwitterFollowers:  tmp.TwFollow,
				})
			}
		}
	}

	if len(SubsData) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SubsData)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

type MessageError struct {
	Message string
	Date    time.Time
}
