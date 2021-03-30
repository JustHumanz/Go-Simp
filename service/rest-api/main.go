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
	"github.com/robfig/cron/v3"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	Data          []map[string]interface{}
	VtuberMembers []database.Member
	VtuberGroups  []database.Group
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	var (
		configfile config.ConfigFile
	)

	RequestPayload := func() {
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

		var Payload database.VtubersPayload
		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}

		for _, Group := range Payload.VtuberData {
			var Mem []map[string]interface{}
			VtuberGroups = append(VtuberGroups, Group)
			for _, Member := range Group.Members {
				VtuberMembers = append(VtuberMembers, Member)
				tmp := map[string]interface{}{
					"ID":       Member.ID,
					"NickName": Member.Name,
					"EnName":   Member.EnName,
					"JpName":   Member.JpName,
					"Region":   Member.Region,
				}
				if Member.BiliBiliID != 0 {
					tmp["BiliBili"] = map[string]interface{}{
						"ID":       Member.BiliBiliID,
						"LiveRoom": Member.BiliRoomID,
						"Avatar":   Member.BiliBiliAvatar,
						"FanArt":   Member.BiliBiliHashtags,
					}
				} else {
					tmp["BiliBili"] = nil
				}

				if Member.YoutubeID != "" {
					tmp["Youtube"] = map[string]interface{}{
						"ID":     Member.YoutubeID,
						"Avatar": Member.YoutubeAvatar,
					}
				} else {
					tmp["Youtube"] = nil
				}

				if Member.TwitterName != "" {
					tmp["Twiter"] = map[string]interface{}{
						"UserName": Member.TwitterName,
						"FanArt":   Member.TwitterHashtags,
						"Lewd":     Member.TwitterLewd,
					}
				} else {
					tmp["Twiter"] = nil
				}

				if Member.TwitchName != "" {
					tmp["Twitch"] = map[string]interface{}{
						"UserName": Member.TwitchName,
						"Avatar":   Member.TwitchAvatar,
					}
				} else {
					tmp["Twitch"] = nil
				}
				Mem = append(Mem, tmp)
			}
			Data = append(Data, map[string]interface{}{
				"ID":        Group.ID,
				"GroupIcon": Group.IconURL,
				"GroupName": Group.GroupName,
				"Members":   Mem,
			})
		}
		configfile.InitConf()
	}
	RequestPayload()

	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, RequestPayload)
	database.Start(configfile)
	go pilot.RunHeartBeat(gRCPconn, "Rest_API")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/all", getGroup).Methods("GET")
	router.HandleFunc("/groups/{GroupID}", getGroup).Methods("GET")
	router.HandleFunc("/members/{MemberID}", getMembers).Methods("GET")

	FanArt := router.PathPrefix("/fanart").Subrouter()
	FanArt.HandleFunc("/", invalidPath).Methods("GET")

	RandomArt := FanArt.PathPrefix("/random").Subrouter()
	RandomArt.HandleFunc("/", invalidPath).Methods("GET")
	RandomArt.HandleFunc("/member/{memberID}", getRandomFanart).Methods("GET")
	RandomArt.HandleFunc("/group/{groupID}", getRandomFanart).Methods("GET")

	BiliArt := FanArt.PathPrefix("/bilibili").Subrouter()
	BiliArt.HandleFunc("/", invalidPath).Methods("GET")
	BiliArt.HandleFunc("/member/{memberID}", getFanart).Methods("GET")
	BiliArt.HandleFunc("/group/{groupID}", getFanart).Methods("GET")

	TwitterArt := FanArt.PathPrefix("/twitter").Subrouter()
	TwitterArt.HandleFunc("/", invalidPath).Methods("GET")
	TwitterArt.HandleFunc("/member/{memberID}", getFanart).Methods("GET")
	TwitterArt.HandleFunc("/group/{groupID}", getFanart).Methods("GET")

	PixivArt := FanArt.PathPrefix("/pixiv").Subrouter()
	PixivArt.HandleFunc("/", invalidPath).Methods("GET")
	PixivArt.HandleFunc("/member/{memberID}", getFanart).Methods("GET")
	PixivArt.HandleFunc("/group/{groupID}", getFanart).Methods("GET")
	/*
		router.HandleFunc("/Twitter", getFanart).Methods("GET")
		router.HandleFunc("/Twitter/Group/{GroupID}", getFanart).Methods("GET")
		router.HandleFunc("/Twitter/Member/{MemberID}", getFanart).Methods("GET")

		router.HandleFunc("/Tbilibili", getFanart).Methods("GET")
		router.HandleFunc("/Tbilibili/Group/{GroupID}", getFanart).Methods("GET")
		router.HandleFunc("/Tbilibili/Member/{MemberID}", getFanart).Methods("GET")

		router.HandleFunc("/Subscriber", getSubs).Methods("GET")
		router.HandleFunc("/Subscriber/Group/{GroupID}", getSubs).Methods("GET")
		router.HandleFunc("/Subscriber/Member/{MemberID}", getSubs).Methods("GET")
	*/

	Live := router.PathPrefix("/livestream").Subrouter()
	Live.HandleFunc("/", invalidPath).Methods("GET")

	LiveYoutube := Live.PathPrefix("/youtube").Subrouter()
	LiveYoutube.HandleFunc("/", invalidPath).Methods("GET")
	LiveYoutube.HandleFunc("/group/{groupID}/{status}", getYoutube).Methods("GET")
	LiveYoutube.HandleFunc("/member/{memberID}/{status}", getYoutube).Methods("GET")

	LiveBili := Live.PathPrefix("/bilibili").Subrouter()
	LiveBili.HandleFunc("/", invalidPath).Methods("GET")
	LiveBili.HandleFunc("/group/{groupID}/{status}", getBilibili).Methods("GET")
	LiveBili.HandleFunc("/member/{memberID}/{status}", getBilibili).Methods("GET")
	http.ListenAndServe(":2525", router)
}

func invalidPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageError{
		Message: "Invalid request.check your request and path",
		Date:    time.Now(),
	})
	w.WriteHeader(http.StatusBadRequest)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["GroupID"]
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []map[string]interface{}
		for _, Group := range Data {
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
				GroupID := Group["ID"].(int64)
				if GroupIDint == int(GroupID) {
					GroupsTMP = append(GroupsTMP, Group)
				}
			}
		}

		if GroupsTMP != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(GroupsTMP)
			w.WriteHeader(http.StatusOK)
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
		json.NewEncoder(w).Encode(Data)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	idstr := mux.Vars(r)["memberID"]
	var Members []map[string]interface{}

	for _, Group := range Data {
		if idstr != "" {
			key := strings.Split(idstr, ",")
			for _, Member := range Group["Members"].([]map[string]interface{}) {
				Member["GroupID"] = Group["ID"]
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
					MemberID := Member["ID"].(int64)
					if region != "" {
						if MemberInt == int(MemberID) && strings.ToLower(region) == strings.ToLower(Member["Region"].(string)) {
							Members = append(Members, Member)
						}
					} else {
						if MemberInt == int(MemberID) {
							Members = append(Members, Member)
						}
					}
				}
			}
		}
	}

	if Members != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Members)
		w.WriteHeader(http.StatusOK)
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

/*

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
		for _, Group := range Data {
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
		for _, Group := range Data {
			for _, Member := range Group["Members"].([]map[string]interface{}) {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						log.Error(err)
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
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							w.WriteHeader(http.StatusInternalServerError)
							return
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
*/

type MessageError struct {
	Message string
	Date    time.Time
}

func GetMember(i int64) database.Member {
	for _, v := range VtuberMembers {
		if v.ID == i {
			return v
		}
	}
	return database.Member{}
}
