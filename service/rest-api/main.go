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
	muxlogrus "github.com/pytimer/mux-logrus"
	"github.com/robfig/cron/v3"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	MembersData   []map[string]interface{}
	GroupsData    []map[string]interface{}
	VtuberMembers []database.Member
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	var (
		configfile config.ConfigFile
	)

	RequestPayload := func() {
		var (
			MembersDataTMP   []map[string]interface{}
			GroupsDataTMP    []map[string]interface{}
			VtuberMembersTMP []database.Member
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

		var Payload database.VtubersPayload
		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}

		configfile.InitConf()
		database.Start(configfile)

		for _, Group := range Payload.VtuberData {
			for _, Member := range Group.Members {
				VtuberMembersTMP = append(VtuberMembersTMP, Member)
				Subs, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}

				MemberData := FixMemberMap(Member)
				MemberData["GroupID"] = Group.ID

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
				MembersDataTMP = append(MembersDataTMP, MemberData)
			}
			GroupsDataTMP = append(GroupsDataTMP, map[string]interface{}{
				"ID":        Group.ID,
				"GroupIcon": Group.IconURL,
				"GroupName": Group.GroupName,
			})
		}
		VtuberMembers = VtuberMembersTMP
		MembersData = MembersDataTMP
		GroupsData = GroupsDataTMP
	}

	RequestPayload()
	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, RequestPayload)
	go pilot.RunHeartBeat(gRCPconn, "Rest_API")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "https://github.com/JustHumanz/Go-Simp/tree/master/service/rest-api/go-simp-api.yml", http.StatusTemporaryRedirect)
	})
	router.HandleFunc("/groups/", getGroup).Methods("GET")
	router.HandleFunc("/groups/{groupID}", getGroup).Methods("GET")

	router.HandleFunc("/members/", getMembers).Methods("GET")
	router.HandleFunc("/members/{memberID}", getMembers).Methods("GET")

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
	router.Use(muxlogrus.NewLogger().Middleware)
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
	vars := mux.Vars(r)["groupID"]
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []map[string]interface{}
		for _, Group := range GroupsData {
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
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(GroupsTMP)
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
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GroupsData)
		w.WriteHeader(http.StatusOK)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	grpID := r.FormValue("groupid")
	idstr := mux.Vars(r)["memberID"]
	var Members []map[string]interface{}

	if idstr != "" {
		key := strings.Split(idstr, ",")
		for _, Member := range MembersData {
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
				if grpID != "" {
					GroupID, err := strconv.Atoi(grpID)
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
						if GroupID == int(Member["GroupID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
							Members = append(Members, Member)
						}
					} else {
						if GroupID == int(Member["GroupID"].(int64)) {
							Members = append(Members, Member)
						}
					}
				} else {
					if region != "" {
						if MemberInt == int(Member["ID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
							Members = append(Members, Member)
						}
					} else {
						if MemberInt == int(Member["ID"].(int64)) {
							Members = append(Members, Member)
						}
					}
				}
			}
		}
	} else if grpID != "" {
		for _, Member := range MembersData {
			if grpID != "" {
				GroupID, err := strconv.Atoi(grpID)
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
					if GroupID == int(Member["GroupID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
						Members = append(Members, Member)
					}
				} else {
					if GroupID == int(Member["GroupID"].(int64)) {
						Members = append(Members, Member)
					}
				}
			}
		}
	} else if region != "" {
		for _, v := range MembersData {
			if strings.EqualFold(region, v["Region"].(string)) {
				Members = append(Members, v)
			}
		}
	} else {
		Members = MembersData
	}

	if Members != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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

func FixMemberMap(Member database.Member) map[string]interface{} {
	var EnName, JpName, Fanbase interface{}

	if Member.EnName == "" {
		EnName = nil
	} else {
		EnName = Member.EnName
	}

	if Member.JpName == "" {
		JpName = nil
	} else {
		JpName = Member.JpName
	}

	if Member.Fanbase == "" {
		Fanbase = nil
	} else {
		Fanbase = Member.Fanbase
	}

	return map[string]interface{}{
		"ID":       Member.ID,
		"NickName": Member.Name,
		"EnName":   EnName,
		"JpName":   JpName,
		"Fanbase":  Fanbase,
		"Region":   Member.Region,
	}
}

func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
