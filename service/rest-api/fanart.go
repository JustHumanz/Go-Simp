package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var FixFanart = func(FanArt database.DataFanart) map[string]interface{} {
	var (
		PixivPostID, TwitterPostID, BiliBiliPostID, Video interface{}
	)

	Data := make(map[string]interface{})
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

	Data["Member"] = FixMemberMap(FanArt.Member)
	Data["Fanart"] = map[string]interface{}{
		"State":       FanArt.State,
		"URL":         FanArt.PermanentURL,
		"Photos":      FanArt.Photos,
		"Video":       Video,
		"Author":      FanArt.Author,
		"PixivID":     PixivPostID,
		"TwitterID":   TwitterPostID,
		"BiliBiliID":  BiliBiliPostID,
		"Description": FanArt.Text,
	}
	return Data
}

func getRandomFanart(w http.ResponseWriter, r *http.Request) {
	var (
		Vars     = mux.Vars(r)
		GroupID  = Vars["groupID"]
		MemberID = Vars["memberID"]
	)
	if GroupID != "" {
		key := strings.Split(GroupID, ",")
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

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(FixFanart(*FanArt))
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	} else if MemberID != "" {
		key := strings.Split(MemberID, ",")
		for _, Member := range MembersData {
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
					FanArt.Member.ID = MemberID
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

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(FixFanart(*FanArt))
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	}
}
func getFanart(w http.ResponseWriter, r *http.Request) {
	var (
		Vars          = mux.Vars(r)
		GroupID       = Vars["groupID"]
		MemberID      = Vars["memberID"]
		FanArtDataFix []map[string]interface{}
	)

	if GroupID != "" {
		key := strings.Split(GroupID, ",")
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
					if strings.HasPrefix(r.URL.String(), "/fanart/twitter/") {
						FanArtData, err := GetFanartData(config.TwitterArt, GroupID, 0)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}

						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}

					} else if strings.HasPrefix(r.URL.String(), "/fanart/pixiv/") {
						FanArtData, err := GetFanartData(config.PixivArt, GroupID, 0)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}
						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					} else if strings.HasPrefix(r.URL.String(), "/fanart/bilibili/") {
						FanArtData, err := GetFanartData(config.BiliBiliArt, GroupID, 0)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}
						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					}
				}
			}
		}
	} else if MemberID != "" {
		key := strings.Split(MemberID, ",")
		for _, Member := range MembersData {
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
					if strings.HasPrefix(r.URL.String(), "/fanart/twitter/") {
						FanArtData, err := GetFanartData(config.TwitterArt, 0, MemberID)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}
						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					} else if strings.HasPrefix(r.URL.String(), "/fanart/pixiv/") {
						FanArtData, err := GetFanartData(config.PixivArt, 0, MemberID)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}
						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					} else if strings.HasPrefix(r.URL.String(), "/fanart/bilibili/") {
						FanArtData, err := GetFanartData(config.BiliBiliArt, 0, MemberID)
						if err != nil {
							log.Error(err)
							w.Header().Set("Content-Type", "application/json")
							json.NewEncoder(w).Encode(MessageError{
								Message: err.Error(),
								Date:    time.Now(),
							})
							return
						}
						for _, v := range FanArtData {
							v.AddMember(GetMember(v.Member.ID))
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					}
				}
			}
		}
	}

	if FanArtDataFix != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(FanArtDataFix)
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
