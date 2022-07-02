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
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/google/uuid"
	muxlogrus "github.com/pytimer/mux-logrus"
	"github.com/robfig/cron/v3"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	VtuberMembers *[]MembersPayload
	VtuberAgency  *[]GroupPayload
	Payload       []*database.Group
	ServiceUUID   = uuid.New().String()
)

func init() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	var (
		configfile config.ConfigFile
	)

	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + config.ResetApiService + " service",
		Service:     config.ResetApiService,
		ServiceUUID: ServiceUUID,
	})
	if err != nil {
		log.Fatalf("Error when request payload: %s", err)
	}

	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}

	hostname := engine.GetHostname()

	RequestPayload := func() {
		var (
			VtuberMembersTMP []MembersPayload
			VtuberAgencyTMP  []GroupPayload
		)
		res2, err := gRCPconn.GetAgencyPayload(context.Background(), &pilot.ServiceMessage{
			Service:     config.ResetApiService,
			Message:     "Request",
			ServiceUUID: ServiceUUID,
			Hostname:    hostname,
		})
		if err != nil {
			log.Fatalf("Error when request payload: %s", err)
		}

		var Payload []*database.Group
		err = json.Unmarshal(res2.AgencyVtubers, &Payload)
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
					"Vtuber": Member.Name,
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
					Group: map[string]interface{}{
						"ID":        Agency.ID,
						"IconURL":   Agency.IconURL,
						"GroupName": Agency.GroupName,
						"YoutubeChannels": func() []map[string]interface{} {
							if Agency.YoutubeChannels != nil {
								tmp := []map[string]interface{}{}
								for _, v := range Agency.YoutubeChannels {
									tmp = append(tmp, map[string]interface{}{
										"YtChannel": v.YtChannel,
										"Region":    v.Region,
									})
								}
								return tmp
							}
							return nil
						}(),
					},
					BiliBili: func() interface{} {
						if Member.BiliBiliID != 0 {
							return map[string]interface{}{
								"Avatar":      Member.BiliBiliAvatar,
								"Banner":      Member.BiliBiliBanner,
								"Fanart":      Member.BiliBiliHashtag,
								"SpaceID":     Member.BiliBiliID,
								"LiveID":      Member.BiliBiliRoomID,
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
								"Banner":      Member.YoutubeBanner,
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
								"Avatar":     Member.TwitterAvatar,
								"Banner":     Member.TwitterBanner,
								"Username":   Member.TwitterName,
								"Fanart":     Member.TwitterHashtag,
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
								"URL": fmt.Sprintf("https://live.bilibili.com/%d", Member.BiliBiliRoomID),
							}
						} else {
							tmp["BiliBili"] = nil
						}

						if isTwitchLive.ID != 0 {
							tmp["Twitch"] = map[string]interface{}{
								"URL": "https://www.twitch.tv/" + Member.TwitchName,
							}
						} else {
							tmp["Twitch"] = nil
						}

						return tmp
					}(),
				})
			}
		}

		VtuberAgency = &VtuberAgencyTMP
		VtuberMembers = &VtuberMembersTMP
	}

	RequestPayload()
	c := cron.New()
	c.Start()
	c.AddFunc(config.CheckPayload, RequestPayload)
	go pilot.RunHeartBeat(gRCPconn, config.ResetApiService, ServiceUUID)
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
	http.ListenAndServe(":2525", engine.LowerCaseURI(router))
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

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["groupID"]
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []GroupPayload
		for _, Group := range *VtuberAgency {
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

	type simpleMember struct {
		ID       int64
		NickName string
		EnName   string
		JpName   string
		Region   string
		Fanbase  string
		Status   string
		BiliBili interface{}
		Youtube  interface{}
		Twitter  interface{}
		Twitch   interface{}
		Group    interface{}
		IsLive   interface{}
	}

	var FixMember = func(Member MembersPayload) simpleMember {
		tmp := simpleMember{
			ID:       Member.ID,
			NickName: Member.NickName,
			EnName:   Member.EnName,
			JpName:   Member.JpName,
			Region:   Member.Region,
			Fanbase:  Member.Fanbase,
			Status:   Member.Status,
			BiliBili: func() interface{} {
				if Member.BiliBili != nil {
					return map[string]interface{}{
						"Avatar":      Member.BiliBili.(map[string]interface{})["Avatar"],
						"SpaceID":     Member.BiliBili.(map[string]interface{})["SpaceID"],
						"LiveID":      Member.BiliBili.(map[string]interface{})["LiveID"],
						"ViwersCount": Member.BiliBili.(map[string]interface{})["ViwersCount"],
						"Followers":   Member.BiliBili.(map[string]interface{})["Followers"],
					}
				}
				return nil
			}(),
			Youtube: func() interface{} {
				if Member.Youtube != nil {
					return map[string]interface{}{
						"Avatar":      Member.Youtube.(map[string]interface{})["Avatar"],
						"YoutubeID":   Member.Youtube.(map[string]interface{})["YoutubeID"],
						"Subscriber":  Member.Youtube.(map[string]interface{})["Subscriber"],
						"ViwersCount": Member.Youtube.(map[string]interface{})["ViwersCount"],
					}
				}
				return nil
			}(),
			Twitter: func() interface{} {
				if Member.Twitter != nil {
					return map[string]interface{}{
						"Avatar":    Member.Twitter.(map[string]interface{})["Avatar"],
						"Username":  Member.Twitter.(map[string]interface{})["Username"],
						"Followers": Member.Twitter.(map[string]interface{})["Followers"],
					}
				}
				return nil
			}(),
			Twitch: Member.Twitch,
			Group: func() interface{} {
				return map[string]interface{}{
					"ID":        Member.Group.(map[string]interface{})["ID"],
					"IconURL":   Member.Group.(map[string]interface{})["IconURL"],
					"GroupName": Member.Group.(map[string]interface{})["GroupName"],
				}
			}(),
			IsLive: Member.IsLive,
		}
		return tmp
	}

	if idstr != "" {
		var Members []MembersPayload
		key := strings.Split(idstr, ",")
		for _, Member := range *VtuberMembers {
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

					memberAgencyID := Member.Group.(map[string]interface{})
					if region != "" {
						if GroupID == int(memberAgencyID["ID"].(int64)) && strings.EqualFold(region, Member.Region) {
							Members = append(Members, Member)
						}
					} else {
						if GroupID == int(memberAgencyID["ID"].(int64)) {
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

		if Members != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(Members)
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
		return
	} else if grpID != "" {
		var Members []simpleMember
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

		for _, Member := range *VtuberMembers {
			memberAgencyID := Member.Group.(map[string]interface{})

			if region != "" {
				if GroupID == int(memberAgencyID["ID"].(int64)) && strings.EqualFold(region, Member.Region) {
					Members = append(Members, FixMember(Member))
				}
			} else {
				if GroupID == int(memberAgencyID["ID"].(int64)) {
					Members = append(Members, FixMember(Member))
				}
			}
		}

		if Members != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(Members)
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
		return

	} else if region != "" {
		var Members []MembersPayload
		for _, v := range *VtuberMembers {
			if strings.EqualFold(region, v.Region) {
				Members = append(Members, v)
			}
		}

		if Members != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(Members)
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
		return

	} else {
		var Members []simpleMember
		for _, Member := range *VtuberMembers {
			Members = append(Members, FixMember(Member))
		}
		if Members != nil {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(Members)
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
		return
	}
}

type MessageError struct {
	Message string
	Date    time.Time
}
