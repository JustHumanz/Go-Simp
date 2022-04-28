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

	Data["Member"] = map[string]interface{}{
		"ID":       FanArt.Member.ID,
		"NickName": FanArt.Member.Name,
		"EnName":   FanArt.Member.EnName,
		"JpName":   FanArt.Member.JpName,
		"Region":   FanArt.Member.Region,
	}
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
		for _, Group := range Payload {
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
					FanArt, err := Group.GetRandomFanart()
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

					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(FixFanart(*FanArt))
					w.WriteHeader(http.StatusOK)
				}
			}
		}
	} else if MemberID != "" {
		key := strings.Split(MemberID, ",")
		for _, Group := range Payload {
			for _, Member := range Group.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					if MemberIDint == int(Member.ID) {
						FanArt, err := Member.GetRandomFanart()
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

						w.Header().Set("Access-Control-Allow-Origin", "*")
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(FixFanart(*FanArt))
						w.WriteHeader(http.StatusOK)
					}
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
		for _, Group := range Payload {
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
					if strings.HasPrefix(r.URL.String(), "/fanart/twitter/") {
						FanArtData, err := Group.GetFanartData(config.TwitterArt, 10)
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
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}

					} else if strings.HasPrefix(r.URL.String(), "/fanart/pixiv/") {
						FanArtData, err := Group.GetFanartData(config.PixivArt, 10) //GetFanartData(config.PixivArt, GroupID, 0)
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
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					} else if strings.HasPrefix(r.URL.String(), "/fanart/bilibili/") {
						FanArtData, err := Group.GetFanartData(config.BiliBiliArt, 10) //GetFanartData(config.BiliBiliArt, GroupID, 0)
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
							FanArtDataFix = append(FanArtDataFix, FixFanart(v))
						}
					}
				}
			}
		}
	} else if MemberID != "" {
		key := strings.Split(MemberID, ",")
		for _, Agency := range Payload {
			for _, Member := range Agency.Members {
				for _, MemberIDstr := range key {
					MemberIDint, err := strconv.Atoi(MemberIDstr)
					if err != nil {
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					if MemberIDint == int(Member.ID) {
						if strings.HasPrefix(r.URL.String(), "/fanart/twitter/") {
							FanArtData, err := Member.GetFanartData(config.TwitterArt, 10) //GetFanartData(config.TwitterArt, 0, MemberID)
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
								FanArtDataFix = append(FanArtDataFix, FixFanart(v))
							}
						} else if strings.HasPrefix(r.URL.String(), "/fanart/pixiv/") {
							FanArtData, err := Member.GetFanartData(config.PixivArt, 10) //GetFanartData(config.PixivArt, 0, MemberID)
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
								FanArtDataFix = append(FanArtDataFix, FixFanart(v))
							}
						} else if strings.HasPrefix(r.URL.String(), "/fanart/bilibili/") {
							FanArtData, err := Member.GetFanartData(config.BiliBiliArt, 10) //GetFanartData(config.BiliBiliArt, 0, MemberID)
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
								FanArtDataFix = append(FanArtDataFix, FixFanart(v))
							}
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
