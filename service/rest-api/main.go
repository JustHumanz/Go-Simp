package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/prediction"
	muxlogrus "github.com/pytimer/mux-logrus"
	"github.com/robfig/cron/v3"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	MembersData    []map[string]interface{}
	GroupsData     []map[string]interface{}
	VtuberMembers  []database.Member
	PredictionConn prediction.PredictionClient
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	PredictionConn = prediction.NewPredictionClient(network.InitgRPC(config.Prediction))

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

		var Payload []*database.Group
		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}

		configfile.InitConf()
		database.Start(configfile)

		for _, Group := range Payload {
			for _, Member := range Group.Members {
				VtuberMembersTMP = append(VtuberMembersTMP, Member)
				Subs, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}

				MemberData := FixMemberMap(Member)
				MemberData["GroupID"] = Group.ID
				MemberData["GroupName"] = Group.GroupName

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
						"UserName":    Member.TwitchName,
						"Avatar":      Member.TwitchAvatar,
						"Followers":   Subs.TwitchFollow,
						"ViwersCount": Subs.TwitchViews,
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

	router.HandleFunc("/prediction/{memberID}", getPrediction).Methods("GET")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageError{
		Message: "Invalid request.check your request and path",
		Date:    time.Now(),
	})
	w.WriteHeader(http.StatusBadRequest)
}

func getPrediction(w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["memberID"]
	days := r.FormValue("days")
	target := r.FormValue("target")
	state := r.FormValue("state")
	var daysint = 7
	var targetint int
	if days != "" {
		var err error
		daysint, err = strconv.Atoi(days)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageError{
				Message: err.Error(),
				Date:    time.Now(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if target != "" {
		var err error
		targetint, err = strconv.Atoi(target)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageError{
				Message: err.Error(),
				Date:    time.Now(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if days != "" && target != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Not support multiple prediction",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if state == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Nill state",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if daysint > 10 || daysint < 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Can't predic more than 10 days",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if idstr != "" {
		idint, err := strconv.Atoi(idstr)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageError{
				Message: err.Error(),
				Date:    time.Now(),
			})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, Member := range MembersData {
			if Member["ID"].(int64) == int64(idint) {
				FinalData := make(map[string]interface{})
				isError := false
				var Fetch = func(wg *sync.WaitGroup, State string, Days bool) {
					defer wg.Done()
					var RawData *prediction.MessageResponse
					var err error
					if Days {
						RawData, err = PredictionConn.GetSubscriberPrediction(context.Background(), &prediction.Message{
							State: State,
							Name:  Member["NickName"].(string),
							Limit: int64(daysint),
						})
					} else {
						RawData, err = PredictionConn.GetReverseSubscriberPrediction(context.Background(), &prediction.Message{
							State: State,
							Name:  Member["NickName"].(string),
							Limit: int64(targetint),
						})
					}
					if err != nil {
						log.Error(err)
						isError = true
					}

					if RawData.Code == 0 {
						Data := map[string]interface{}{}
						if State == "Twitter" {
							Tw := Member["Twitter"].(map[string]interface{})
							Data["Current_followes/subscriber"] = Tw["Followers"]
						} else if State == "Youtube" {
							Yt := Member["Youtube"].(map[string]interface{})
							Data["Current_followes/subscriber"] = Yt["Subscriber"]
						} else if State == "BiliBili" {
							Bl := Member["BiliBili"].(map[string]interface{})
							Data["Current_followes/subscriber"] = Bl["Followers"]
						}
						if Days {
							Data["Prediction"] = RawData.Prediction
						} else {
							if Data["Current_followes/subscriber"].(int) < targetint {
								Data["Prediction"] = targetint
							} else {
								isError = true
							}
						}

						Data["Score"] = RawData.Score
						if targetint == 0 {
							FinalData["Prediction_Days"] = time.Now().AddDate(0, 0, daysint)
						} else {
							unxtoko := time.Unix(RawData.Prediction, 0)
							FinalData["Prediction_Days"] = unxtoko
						}

						FinalData[State] = Data
					} else {
						isError = true
					}
				}

				var wg sync.WaitGroup

				for _, st := range strings.Split(state, ",") {
					if st == "bilibili" {
						st = "BiliBili"
					}
					if Member[strings.Title(st)] != nil {
						wg.Add(1)
						if targetint == 0 {
							go Fetch(&wg, strings.Title(st), true)
						} else {
							go Fetch(&wg, strings.Title(st), false)
							break
						}
					} else {
						FinalData[strings.Title(st)] = nil
					}
				}

				wg.Wait()

				if !isError {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(FinalData)
					w.WriteHeader(http.StatusOK)
				} else {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: "oops,something goes wrong",
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusBadRequest)
				}
				return
			}
		}
	}
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
					w.WriteHeader(http.StatusBadRequest)
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
	cekLive := r.FormValue("live")
	idstr := mux.Vars(r)["memberID"]
	var Members []map[string]interface{}
	var ww sync.WaitGroup

	var CekLiveMember = func(Member map[string]interface{}) map[string]interface{} {
		if cekLive == "true" {
			MemberID := Member["ID"].(int64)
			MemberName := Member["NickName"].(string)
			if Member["Youtube"] != nil {
				YTData, _, err := database.YtGetStatus(map[string]interface{}{
					"MemberID":   MemberID,
					"MemberName": MemberName,
					"Status":     config.LiveStatus,
					"State":      config.Sys,
				})
				if err != nil {
					log.Error(err)
				}
				if YTData != nil {
					Member["IsYtLive"] = true
					Member["LiveURL"] = "https://www.youtube.com/watch?v=" + YTData[0].VideoID
				} else {
					Member["IsYtLive"] = false
				}
			}

			if Member["BiliBili"] != nil {
				BiliData, _, err := database.BilGet(map[string]interface{}{
					"MemberID": Member["ID"].(int64),
					"Status":   config.LiveStatus,
				})
				if err != nil {
					log.Error(err)
				}

				if BiliData != nil {
					Bili := Member["BiliBili"].(map[string]interface{})
					Member["IsBiliLive"] = true
					Member["LiveURL"] = "https://live.bilibili.com/" + strconv.Itoa(Bili["RoomID"].(int))
				} else {
					Member["IsBiliLive"] = false
				}
			}

			if Member["Twitch"] != nil {
				TwitchData, err := database.GetTwitch(Member["ID"].(int64))
				if err != nil {
					log.Error(err)
				}
				if TwitchData != nil {
					Twitch := Member["Twitch"].(map[string]interface{})
					Member["IsTwitchLive"] = true
					Member["LiveURL"] = "https://www.twitch.tv/" + Twitch["UserName"].(string)
				} else {
					Member["IsTwitchLive"] = false
				}
			}
		}
		return Member
	}
	if idstr != "" {
		key := strings.Split(idstr, ",")
		for _, M := range MembersData {
			ww.Add(1)
			go func(Member map[string]interface{}, wg *sync.WaitGroup) {
				defer wg.Done()
				for _, MemberStr := range key {
					MemberInt, err := strconv.Atoi(MemberStr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusBadRequest)
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
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						if region != "" {
							if GroupID == int(Member["GroupID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
								Members = append(Members, CekLiveMember(Member))
							}
						} else {
							if GroupID == int(Member["GroupID"].(int64)) {
								Members = append(Members, CekLiveMember(Member))
							}
						}
					} else {
						if region != "" {
							if MemberInt == int(Member["ID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
								Members = append(Members, CekLiveMember(Member))
							}
						} else {
							if MemberInt == int(Member["ID"].(int64)) {
								Members = append(Members, CekLiveMember(Member))
							}
						}
					}
				}
			}(M, &ww)
		}
		ww.Wait()
	} else if grpID != "" {
		for _, M := range MembersData {
			ww.Add(1)
			go func(Member map[string]interface{}, wg *sync.WaitGroup) {
				defer wg.Done()
				if grpID != "" {
					GroupID, err := strconv.Atoi(grpID)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					if region != "" {
						if GroupID == int(Member["GroupID"].(int64)) && strings.EqualFold(region, Member["Region"].(string)) {
							Members = append(Members, CekLiveMember(Member))
						}
					} else {
						if GroupID == int(Member["GroupID"].(int64)) {
							Members = append(Members, CekLiveMember(Member))
						}
					}
				}
			}(M, &ww)
		}
		ww.Wait()
	} else if region != "" {
		for _, v := range MembersData {
			if strings.EqualFold(region, v["Region"].(string)) {
				Members = append(Members, CekLiveMember(v))
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
