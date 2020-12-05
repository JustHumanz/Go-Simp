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

	router.HandleFunc("/Twitter/", getTwitter).Methods("GET")

	http.ListenAndServe(":8000", router)
}

func getTwitter(w http.ResponseWriter, r *http.Request) {

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
	groupidstr := r.FormValue("groupid")
	region := r.FormValue("region")
	idstr := mux.Vars(r)["MemberID"]
	Members := []database.Member{}

	GetMembers := func(GroupID int64) {
		if idstr != "" {
			key := strings.Split(idstr, ",")
			for _, Member := range database.GetMembers(GroupID) {
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
		} else {
			for _, Member := range database.GetMembers(GroupID) {
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

	for _, Group := range database.GetGroups() {
		if groupidstr != "" {
			key := strings.Split(groupidstr, ",")
			for _, GroupID := range key {
				GroupIDInt, err := strconv.Atoi(GroupID)
				if err != nil {
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(MessageError{
						Message: err.Error(),
						Date:    time.Now(),
					})
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if GroupIDInt == int(Group.ID) {
					GetMembers(Group.ID)
				}
			}
		} else {
			GetMembers(Group.ID)
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

type MessageError struct {
	Message string
	Date    time.Time
}
