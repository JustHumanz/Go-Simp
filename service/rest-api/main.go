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

		configfile.InitConf()
		database.Start(configfile)

		for _, Group := range Payload.VtuberData {
			var Mem []map[string]interface{}
			VtuberGroups = append(VtuberGroups, Group)
			for _, Member := range Group.Members {
				Subs, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}

				MemberData := map[string]interface{}{
					"ID":       Member.ID,
					"NickName": Member.Name,
					"EnName":   Member.EnName,
					"JpName":   Member.JpName,
					"Fanbase":  Member.Fanbase,
					"Region":   Member.Region,
				}

				if Member.YoutubeID != "" {
					MemberData["Youtube"] = map[string]interface{}{
						"ID":          Member.YoutubeID,
						"Avatar":      Member.YoutubeAvatar,
						"Subscriber":  Subs.YtSubs,
						"ViwersCount": Subs.YtViews,
						"TotalVideos": Subs.YtVideos,
					}
				} else {
					MemberData["Youtube"] = nil
				}

				if Member.BiliBiliID != 0 {
					var fanart interface{}
					if Member.BiliBiliHashtags != "" {
						fanart = Member.BiliBiliHashtags
					} else {
						fanart = nil
					}
					MemberData["BiliBili"] = map[string]interface{}{
						"ID":          Member.BiliBiliID,
						"RoomID":      Member.BiliRoomID,
						"Avatar":      Member.BiliBiliAvatar,
						"Fanart":      fanart,
						"Followers":   Subs.BiliFollow,
						"TotalVideos": Subs.BiliVideos,
						"ViwersCount": Subs.BiliViews,
					}
				} else {
					MemberData["BiliBili"] = nil
				}

				if Member.TwitterName != "" {
					var (
						Lewd   interface{}
						Fanart interface{}
					)

					if Member.TwitterLewd != "" {
						Lewd = Member.TwitterLewd
					} else {
						Lewd = nil
					}

					if Member.TwitterHashtags != "" {
						Fanart = Member.TwitterHashtags
					} else {
						Fanart = nil
					}

					MemberData["Twitter"] = map[string]interface{}{
						"UserName":  Member.TwitterName,
						"Fanart":    Fanart,
						"Lewd":      Lewd,
						"Followers": Subs.TwFollow,
					}
				} else {
					MemberData["Twitter"] = nil
				}

				if Member.TwitchName != "" {
					MemberData["Twitch"] = map[string]interface{}{
						"UserName": Member.TwitchName,
						"Avatar":   Member.TwitchAvatar,
					}
				} else {
					MemberData["Twitch"] = nil
				}
				Mem = append(Mem, MemberData)
				VtuberMembers = append(VtuberMembers, Member)
			}
			Data = append(Data, map[string]interface{}{
				"ID":        Group.ID,
				"GroupIcon": Group.IconURL,
				"GroupName": Group.GroupName,
				"Members":   Mem,
			})
		}
	}

	RequestPayload()
	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, RequestPayload)
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
	http.ListenAndServe(":2525", LowerCaseURI(router))
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
				Member["GroupID"] = Group["ID"].(int64)
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
						if MemberInt == int(Member["GroupID"].(int64)) && strings.ToLower(region) == strings.ToLower(Member["Region"].(string)) {
							Members = append(Members, Member)
						}
					} else {
						if MemberInt == int(Member["GroupID"].(int64)) {
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

func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
