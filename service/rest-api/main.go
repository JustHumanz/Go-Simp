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

	router.HandleFunc("/youtube/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/youtube/group/{GroupID}/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/youtube/member/{MemberID}/{Status}", getYoutube).Methods("GET")

	router.HandleFunc("/bilibili/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/bilibili/Group/{GroupID}/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/bilibili/Member/{MemberID}/{Status}", getBilibili).Methods("GET")

	FanArt := router.PathPrefix("/fanart").Subrouter()
	FanArt.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Invalid request",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusBadRequest)
	}).Methods("GET")
	FanArt.HandleFunc("/random/member/{MemberID}", getRandomFanart).Methods("GET")
	FanArt.HandleFunc("/random/group/{groupID}", getRandomFanart).Methods("GET")

	FanArt.HandleFunc("/bilibili/random", getFanart).Methods("GET")
	FanArt.HandleFunc("/bilibili/member/{MemberID}", getFanart).Methods("GET")
	FanArt.HandleFunc("/bilibili/group/{groupID}", getFanart).Methods("GET")

	FanArt.HandleFunc("/pixiv/random", getFanart).Methods("GET")
	FanArt.HandleFunc("/pixiv/member/{MemberID}", getFanart).Methods("GET")
	FanArt.HandleFunc("/pixiv/group/{groupID}", getFanart).Methods("GET")
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

	http.ListenAndServe(":2525", router)
}

func getRandomFanart(w http.ResponseWriter, r *http.Request) {
	var (
		Vars = mux.Vars(r)
	)
	if Vars["groupID"] != "" {
		key := strings.Split(Vars["groupID"], ",")
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
					FanArt, err := database.GetFanart(GroupID, 0)
					if err != nil {
						log.Error(err)
						errstr := err.Error()
						if FanArt == nil {
							errstr = "Opps,fanart not found"
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: errstr,
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					FanArt.AddMember(GetMember(FanArt.Member.ID))
					var (
						PixivPostID, TwitterPostID, BiliBiliPostID, Video interface{}
					)

					if FanArt.State == config.PixivArt {
						PixivPostID = FanArt.PixivID
					} else {
						PixivPostID = nil
					}

					if FanArt.State == config.TwitterArt {
						TwitterPostID = FanArt.TweetID
					} else {
						TwitterPostID = nil
					}

					if FanArt.State == config.BiliBiliArt {
						BiliBiliPostID = FanArt.Dynamic_id
					} else {
						BiliBiliPostID = nil
					}

					if FanArt.Videos == "" {
						Video = nil
					} else {
						Video = FanArt.Videos
					}
					Art := map[string]interface{}{
						"Member": map[string]interface{}{
							"ID":       FanArt.Member.ID,
							"NickName": FanArt.Member.Name,
							"EnName":   FanArt.Member.EnName,
							"JpName":   FanArt.Member.JpName,
							"Region":   FanArt.Member.Region,
							"Hashtags": map[string]interface{}{
								"TwitterFanart":  FanArt.Member.TwitterHashtags,
								"BiliBiliFanart": FanArt.Member.BiliBiliHashtags,
							},
						},
						"Fanart": map[string]interface{}{
							"State":      FanArt.State,
							"URL":        FanArt.PermanentURL,
							"Photos":     FanArt.Photos,
							"Video":      Video,
							"Author":     FanArt.Author,
							"PixivID":    PixivPostID,
							"TwitterID":  TwitterPostID,
							"BiliBiliID": BiliBiliPostID,
						},
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(Art)
					w.WriteHeader(http.StatusOK)
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
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					MemberID := Member["ID"].(int64)
					if MemberIDint == int(MemberID) {
						FanArt, err := database.GetFanart(0, MemberID)
						if err != nil {
							log.Error(err)
							errstr := err.Error()
							if FanArt == nil {
								errstr = "Opps,fanart not found"
							}
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: errstr,
								Date:    time.Now(),
							})
						}
						FanArt.AddMember(GetMember(FanArt.Member.ID))
						var (
							PixivPostID, TwitterPostID, BiliBiliPostID, Video interface{}
						)

						if FanArt.State == config.PixivArt {
							PixivPostID = FanArt.PixivID
						} else {
							PixivPostID = nil
						}

						if FanArt.State == config.TwitterArt {
							TwitterPostID = FanArt.TweetID
						} else {
							TwitterPostID = nil
						}

						if FanArt.State == config.BiliBiliArt {
							BiliBiliPostID = FanArt.Dynamic_id
						} else {
							BiliBiliPostID = nil
						}

						if FanArt.Videos == "" {
							Video = nil
						} else {
							Video = FanArt.Videos
						}
						Art := map[string]interface{}{
							"Member": map[string]interface{}{
								"ID":       FanArt.Member.ID,
								"NickName": FanArt.Member.Name,
								"EnName":   FanArt.Member.EnName,
								"JpName":   FanArt.Member.JpName,
								"Region":   FanArt.Member.Region,
								"Hashtags": map[string]interface{}{
									"TwitterFanart":  FanArt.Member.TwitterHashtags,
									"BiliBiliFanart": FanArt.Member.BiliBiliHashtags,
								},
							},
							"Fanart": map[string]interface{}{
								"State":      FanArt.State,
								"URL":        FanArt.PermanentURL,
								"Photos":     FanArt.Photos,
								"Video":      Video,
								"Author":     FanArt.Author,
								"PixivID":    PixivPostID,
								"TwitterID":  TwitterPostID,
								"BiliBiliID": BiliBiliPostID,
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(Art)
						w.WriteHeader(http.StatusOK)
					}
				}
			}
		}
	}
}

func getFanart(w http.ResponseWriter, r *http.Request) {
	var (
		Vars   = mux.Vars(r)
		Fanart []database.DataFanart
	)

	if strings.HasPrefix(r.URL.String(), "/fanart/twitter/") {

	} else if r.URL.String() == "/fanart/bilibili/random" {

	} else if r.URL.String() == "/fanart/pixiv/random" {

	}

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
				GroupID := Group["ID"].(int64)
				if GroupIDint == int(GroupID) {
					//Fanart = append(Fanart, GetFanartData(State, GroupID, 0)...)
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
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					MemberID := Member["ID"].(int64)
					if MemberIDint == int(MemberID) {
						//	Fanart = append(Fanart, GetFanartData(State, 0, MemberID)...)
					}
				}
			}
		}
	} else {
		/*
			for _, Group := range Data {
				//Fanart = append(Fanart, GetFanartData(State, Group.ID, 0)...)
			}
		*/
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
		YoutubeData []database.LiveStream
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
					YTData, err := database.YtGetStatus(GroupID, 0, Status, Region)
					if err != nil {
						log.Error(err)
					}
					YoutubeData = append(YoutubeData, YTData...)
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
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					MemberID := Member["ID"].(int64)
					if MemberIDint == int(MemberID) {
						YTData, err := database.YtGetStatus(0, MemberID, Status, Region)
						if err != nil {
							log.Error(err)
						}
						YoutubeData = append(YoutubeData, YTData...)
					}
				}
			}
		}
	} else {
		for _, Group := range Data {
			YTData, err := database.YtGetStatus(Group["ID"].(int64), 0, Status, Region)
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
		BiliBiliData []database.LiveStream
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
					BiliBiliData = append(BiliBiliData, database.BilGet(GroupID, 0, Status)...)
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
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					MemberID := Member["ID"].(int64)
					if MemberIDint == int(MemberID) {
						BiliBiliData = append(BiliBiliData, database.BilGet(0, MemberID, Status)...)
					}
				}
			}
		}
	} else {
		for _, Group := range Data {
			BiliBiliData = append(BiliBiliData, database.BilGet(Group["ID"].(int64), 0, Status)...)
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
		json.NewEncoder(w).Encode(Data)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	idstr := mux.Vars(r)["MemberID"]
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
