package main

import (
	"context"
	"encoding/json"
	"fmt"
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
	VtuberMembers []MembersPayload
	VtuberAgency  []GroupPayload
	Payload       []*database.Group
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	var (
		configfile config.ConfigFile
	)

	RequestPayload := func() {
		log.Info("Request payload")

		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Rest_API",
		})
		if err != nil {
			log.Fatalf("Error when request payload: %s", err)
		}
		var (
			VtuberMembersTMP []MembersPayload
			VtuberAgencyTMP  []GroupPayload
		)

		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}

		configfile.InitConf()
		database.Start(configfile)

		for _, Agency := range Payload {
			VtuberAgencyTMP = append(VtuberAgencyTMP, GroupPayload{
				ID:        Agency.ID,
				GroupName: Agency.GroupName,
				GroupIcon: Agency.IconURL,
				Youtube: func() interface{} {
					if Agency.YoutubeChannels != nil {
						return Agency.YoutubeChannels
					} else {
						return nil
					}
				}(),
			})
			for _, Member := range Agency.Members {
				Subs, err := Member.GetSubsCount()
				if err != nil {
					log.Error(err)
				}

				isYtLive, err := Member.GetYtLiveStream(config.LiveStatus)
				if err != nil {
					log.Error(err)
				}

				isBlLive, err := Member.GetBlLiveStream(config.LiveStatus)
				if err != nil {
					log.Error(err)
				}

				isTwitchLive, err := Member.GetTwitchLiveStream(config.LiveStatus)
				if err != nil {
					log.Error(err)
				}

				log.WithFields(log.Fields{
					"Group":  Agency.GroupName,
					"Vtuber": Member.EnName,
					"isYtLive": func() bool {
						if isYtLive != nil {
							return true
						} else {
							return false
						}
					}(),
					"isBiliLive":   isBlLive.ID,
					"isTwitchLive": isTwitchLive.ID,
				}).Info("Processing data")

				VtuberMembersTMP = append(VtuberMembersTMP, MembersPayload{
					ID:       Member.ID,
					NickName: Member.Name,
					EnName:   Member.EnName,
					JpName:   Member.JpName,
					Region:   Member.Region,
					Status:   Member.Status,
					Fanbase:  Member.Fanbase,
					Group: database.Group{
						ID:              Agency.ID,
						IconURL:         Agency.IconURL,
						GroupName:       Agency.GroupName,
						YoutubeChannels: Agency.YoutubeChannels,
					},
					BiliBili: func() interface{} {
						if Member.BiliBiliID != 0 {
							return map[string]interface{}{
								"Avatar":      Member.BiliBiliAvatar,
								"Fanart":      Member.BiliBiliHashtags,
								"SpaceID":     Member.BiliBiliID,
								"LiveID":      Member.BiliRoomID,
								"TotalVideos": Subs.BiliVideos,
								"ViwersCount": Subs.BiliViews,
								"Followers":   Subs.BiliFollow,
							}
						} else {
							return nil
						}
					}(),
					Youtube: func() interface{} {
						if Member.YoutubeID != "" {
							return map[string]interface{}{
								"Avatar":      Member.YoutubeAvatar,
								"YoutubeID":   Member.YoutubeID,
								"Subscriber":  Subs.YtSubs,
								"TotalVideos": Subs.YtVideos,
								"ViwersCount": Subs.YtViews,
							}
						} else {
							return nil
						}
					}(),
					Twitter: func() interface{} {
						if Member.TwitterName != "" {
							return map[string]interface{}{
								"Username":   Member.TwitterName,
								"Fanart":     Member.TwitterHashtags,
								"LewdFanart": Member.TwitterLewd,
								"Followers":  Subs.TwFollow,
							}
						} else {
							return nil
						}
					}(),
					Twitch: func() interface{} {
						if Member.TwitchName != "" {
							return map[string]interface{}{
								"Avatar":      Member.TwitchAvatar,
								"Username":    Member.TwitchName,
								"Followers":   Subs.TwitchFollow,
								"ViwersCount": Subs.TwitchViews,
							}
						} else {
							return nil
						}
					}(),
					IsLive: func() interface{} {
						tmp := make(map[string]interface{})

						if len(isYtLive) > 0 {
							for _, v := range isYtLive {
								tmp["Youtube"] = map[string]interface{}{
									"URL": "https://www.youtube.com/watch?v=" + v.VideoID,
								}
								break
							}
						} else {
							tmp["Youtube"] = nil
						}

						if isBlLive.ID != 0 {
							tmp["BiliBili"] = map[string]interface{}{
								"URL": fmt.Sprintf("https://live.bilibili.com/%d", Member.BiliRoomID),
							}
						} else {
							tmp["BiliBili"] = nil
						}

						if isTwitchLive.ID != 0 {
							tmp["Twitch"] = map[string]interface{}{
								"URL": fmt.Sprintf("https://www.twitch.tv/%s", Member.TwitchName),
							}
						} else {
							tmp["Twitch"] = nil
						}

						return tmp
					}(),
				})
			}
		}

		VtuberAgency = VtuberAgencyTMP
		VtuberMembers = VtuberMembersTMP
	}

	RequestPayload()
	c := cron.New()
	c.Start()
	c.AddFunc("@every 0h20m0s", RequestPayload)
	go pilot.RunHeartBeat(gRCPconn, "Rest_API")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, "https://github.com/JustHumanz/Go-Simp/tree/master/service/rest-api/go-simp-api.yml", http.StatusTemporaryRedirect)
	})
	router.HandleFunc("/groups/", getGroup).Methods("GET")
	router.HandleFunc("/groups/{groupID}", getGroup).Methods("GET")

	//router.HandleFunc("/prediction/{memberID}", getPrediction).Methods("GET")

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

//Need rework
/*
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
		for _, Member := range VtuberMembers {
			if Member.ID == int64(idint) {
				FinalData := make(map[string]interface{})
				isError := false
				var Fetch = func(wg *sync.WaitGroup, State string, Days bool) {
					defer wg.Done()
					var DataPredic int
					var err error
					if Days {
						DataPredic, err = engine.Prediction(database.Member{Name: Member.NickName}, State, daysint)
						if err != nil {
							log.Error(err)
							isError = true
						}
					}
					if err != nil {
						log.Error(err)
						isError = true
					}

					Data := map[string]interface{}{}
					if State == "Twitter" {
						Data["Current_followes/subscriber"] = Member.Twitter["Followers"]
					} else if State == "Youtube" {
						Yt := Member["Youtube"].(map[string]interface{})
						Data["Current_followes/subscriber"] = Yt["Subscriber"]
					} else if State == "BiliBili" {
						Bl := Member["BiliBili"].(map[string]interface{})
						Data["Current_followes/subscriber"] = Bl["Followers"]
					}
					if Days {
						Data["Prediction"] = DataPredic
					} else {
						if Data["Current_followes/subscriber"].(int) < targetint {
							Data["Prediction"] = targetint
						} else {
							isError = true
						}
					}

					if targetint == 0 {
						FinalData["Prediction_Days"] = time.Now().AddDate(0, 0, daysint)
					}

					FinalData[State] = Data
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
*/

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["groupID"]
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []GroupPayload
		for _, Group := range VtuberAgency {
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

				if GroupIDint == int(Group.ID) {
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
		json.NewEncoder(w).Encode(VtuberAgency)
		w.WriteHeader(http.StatusOK)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	grpID := r.FormValue("groupid")
	idstr := mux.Vars(r)["memberID"]
	var Members []MembersPayload

	if idstr != "" {
		key := strings.Split(idstr, ",")

		for _, Member := range VtuberMembers {
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
						if GroupID == int(Member.Group.ID) && strings.EqualFold(region, Member.Region) {
							Members = append(Members, Member)
						}
					} else {
						if GroupID == int(Member.Group.ID) {
							Members = append(Members, Member)
						}
					}
				} else {
					if region != "" {
						if MemberInt == int(Member.ID) && strings.EqualFold(region, Member.Region) {
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
	} else if grpID != "" {
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

		for _, Member := range VtuberMembers {
			if region != "" {
				if GroupID == int(Member.Group.ID) && strings.EqualFold(region, Member.Region) {
					Members = append(Members, Member)
				}
			} else {
				if GroupID == int(Member.Group.ID) {
					Members = append(Members, Member)
				}
			}
		}
	} else if region != "" {
		for _, v := range VtuberMembers {
			if strings.EqualFold(region, v.Region) {
				Members = append(Members, v)
			}
		}
	} else {
		Members = VtuberMembers
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

func GetMember(i int64) MembersPayload {
	for _, v := range VtuberMembers {
		if v.ID == i {
			return v
		}
	}
	return MembersPayload{}
}

func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
