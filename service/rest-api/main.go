package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var (
	Groups []database.Group
)

func init() {
	config, err := config.ReadConfig("../../config.toml")
	if err != nil {
		log.Error(err)
	}
	database.Start(config.CheckSQL())
	Groups = database.GetGroups()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/Groups/", getGroup).Methods("GET")
	router.HandleFunc("/Groups/{GroupID}", getGroup).Methods("GET")

	router.HandleFunc("/Members/", getMembers).Methods("GET")
	router.HandleFunc("/Members/{MemberID}", getMembers).Methods("GET")

	router.HandleFunc("/Youtube/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/Youtube/Group/{GroupID}/{Status}", getYoutube).Methods("GET")
	router.HandleFunc("/Youtube/Member/{MemberID}/{Status}", getYoutube).Methods("GET")

	router.HandleFunc("/Bilibili/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/Bilibili/Group/{GroupID}/{Status}", getBilibili).Methods("GET")
	router.HandleFunc("/Bilibili/Member/{MemberID}/{Status}", getBilibili).Methods("GET")

	router.HandleFunc("/Twitter/", getFanart).Methods("GET")
	router.HandleFunc("/Twitter/Group/{GroupID}", getFanart).Methods("GET")
	router.HandleFunc("/Twitter/Member/{MemberID}", getFanart).Methods("GET")

	router.HandleFunc("/Tbilibili/", getFanart).Methods("GET")
	router.HandleFunc("/Tbilibili/Group/{GroupID}", getFanart).Methods("GET")
	router.HandleFunc("/Tbilibili/Member/{MemberID}", getFanart).Methods("GET")

	router.HandleFunc("/Subscribe/", getSubs).Methods("GET")
	router.HandleFunc("/Subscribe/Group/{GroupID}", getSubs).Methods("GET")
	router.HandleFunc("/Subscribe/Member/{MemberID}", getSubs).Methods("GET")

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
		for _, Group := range Groups {
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
		for _, Group := range Groups {
			for _, Member := range database.GetMembers(Group.ID) {
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
		for _, Group := range Groups {
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
		for _, Group := range Groups {
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
					YoutubeData = append(YoutubeData, database.YtGetStatus(Group.ID, 0, Status, Region)...)
				}
			}
		}
	} else if Vars["MemberID"] != "" {
		key := strings.Split(Vars["MemberID"], ",")
		for _, Group := range Groups {
			for _, Member := range database.GetMembers(Group.ID) {
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
						YoutubeData = append(YoutubeData, database.YtGetStatus(0, Member.ID, Status, Region)...)
					}
				}
			}
		}
	} else {
		for _, Group := range Groups {
			YoutubeData = append(YoutubeData, database.YtGetStatus(Group.ID, 0, Status, Region)...)
		}
	}

	if len(YoutubeData) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(YoutubeData)
	} else {
		w.Header().Set("Content-Type", "application/json")
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
		for _, Group := range Groups {
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
		for _, Group := range Groups {
			for _, Member := range database.GetMembers(Group.ID) {
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
		for _, Group := range Groups {
			BiliBiliData = append(BiliBiliData, database.BilGet(Group.ID, 0, Status)...)
		}
	}

	if len(BiliBiliData) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(BiliBiliData)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Reqest not found,404",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
	}
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["GroupID"]
	Groups = database.GetGroups()
	if vars != "" {
		key := strings.Split(vars, ",")
		var GroupsTMP []database.Group
		for _, Group := range Groups {
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
			json.NewEncoder(w).Encode(MessageError{
				Message: "Reqest not found,404",
				Date:    time.Now(),
			})
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Groups)
	}
}

func getMembers(w http.ResponseWriter, r *http.Request) {
	region := r.FormValue("region")
	idstr := mux.Vars(r)["MemberID"]
	Members := []database.Member{}

	for _, Group := range database.GetGroups() {
		if idstr != "" {
			key := strings.Split(idstr, ",")
			for _, Member := range database.GetMembers(Group.ID) {
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
							Member.GroupID = Group.ID
						}
					} else {
						if MemberInt == int(Member.ID) {
							Member.GroupID = Group.ID
						}
					}
				}
			}
		} else {
			for _, Member := range database.GetMembers(Group.ID) {
				Member.GroupID = Group.ID
				if region != "" {
					if strings.ToLower(region) == strings.ToLower(Member.Region) {
						Members = append(Members, Member)
					}
				} else {
					Members = append(Members, Member)
				}
			}
		}
	}
	if len(Members) > 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Members)
	} else {
		w.Header().Set("Content-Type", "application/json")
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
		for _, Group := range Groups {
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
					for _, Member := range database.GetMembers(Group.ID) {
						tmp := Member.GetSubsCount()
						SubsData = append(SubsData, SubsJson{
							MemberID:          tmp.MemberID,
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
		for _, Group := range Groups {
			for _, Member := range database.GetMembers(Group.ID) {
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
						tmp := Member.GetSubsCount()
						SubsData = append(SubsData, SubsJson{
							MemberID:          tmp.MemberID,
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
		for _, Group := range Groups {
			for _, Member := range database.GetMembers(Group.ID) {
				tmp := Member.GetSubsCount()
				SubsData = append(SubsData, SubsJson{
					MemberID:          tmp.MemberID,
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
