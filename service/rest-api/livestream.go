package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func getYoutube(w http.ResponseWriter, r *http.Request) {
	var (
		YoutubeData []map[string]interface{}
		Vars        = mux.Vars(r)
		Region      = strings.ToLower(r.FormValue("region"))
		Status      = strings.ToLower(Vars["status"])
		GroupIDs    = Vars["groupID"]
		MemberIDs   = Vars["memberID"]
		Rgx         = "(" + config.LiveStatus + "|" + config.PastStatus + "|" + config.UpcomingStatus + "|" + config.PrivateStatus + ")"
		ww          sync.WaitGroup
	)

	if match, _ := regexp.MatchString(Rgx, Status); !match {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Status not found",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if GroupIDs != "" {
		key := strings.Split(GroupIDs, ",")
		for _, GroupData := range Payload {
			ww.Add(1)
			go func(Agency database.Group, wg *sync.WaitGroup) {
				defer wg.Done()
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

					if GroupIDint == int(Agency.ID) {
						YTData, err := Agency.GetYtLiveStream(
							Status,
							Region,
						)
						if err != nil {
							log.Error(err)
						}
						for _, v := range YTData {
							v.SetState(config.YoutubeLive)
							YoutubeData = append(YoutubeData, FixLive(v))
						}
					}
				}
			}(*GroupData, &ww)
		}
		ww.Wait()
	} else if MemberIDs != "" {
		key := strings.Split(MemberIDs, ",")
		for _, Agency := range Payload {
			for _, M := range Agency.Members {
				ww.Add(1)
				go func(Member database.Member, wg *sync.WaitGroup) {
					defer wg.Done()
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
							YTData, err := Member.GetYtLiveStream(Status)
							if err != nil {
								log.Error(err)
							}
							for _, v := range YTData {
								v.SetState(config.YoutubeLive)
								YoutubeData = append(YoutubeData, FixLive(v))
							}
						}
					}
				}(M, &ww)
			}
		}
		ww.Wait()
	}

	if YoutubeData != nil {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(YoutubeData)
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

func getBilibili(w http.ResponseWriter, r *http.Request) {
	var (
		BiliBiliData []map[string]interface{}
		Vars         = mux.Vars(r)
		Status       = strings.ToLower(Vars["status"])
		GroupID      = Vars["groupID"]
		MemberID     = Vars["memberID"]
		Rgx          = "(" + config.LiveStatus + "|" + config.PastStatus + ")"
	)

	if match, _ := regexp.MatchString(Rgx, Status); !match {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(MessageError{
			Message: "Status not found",
			Date:    time.Now(),
		})
		w.WriteHeader(http.StatusNotFound)
		return
	}
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
					BiliData, err := Group.GetBlLiveStream(Status)
					if err != nil {
						log.Error(err)
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(MessageError{
							Message: err.Error(),
							Date:    time.Now(),
						})
						w.WriteHeader(http.StatusBadRequest)
						return

					}

					for _, v := range BiliData {
						if v.ID == 0 {
							break
						}

						v.SetState(config.BiliLive)
						BiliBiliData = append(BiliBiliData, FixLive(v))
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
						BiliData, err := Member.GetBlLiveStream(Status)
						if err != nil {
							log.Error(err)
						}

						if BiliData.ID == 0 {
							break
						}

						BiliData.SetState(config.BiliLive)
						BiliBiliData = append(BiliBiliData, FixLive(BiliData))
					}
				}
			}
		}
	}

	if BiliBiliData != nil {
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

var FixLive = func(Live database.LiveStream) map[string]interface{} {
	if Live.State == config.YoutubeLive {
		var EndStream interface{}
		if Live.End.IsZero() {
			EndStream = nil
		} else {
			EndStream = Live.End
		}
		return map[string]interface{}{
			"Member": map[string]interface{}{
				"ID":            Live.Member.ID,
				"NickName":      Live.Member.Name,
				"EnName":        Live.Member.EnName,
				"JpName":        Live.Member.JpName,
				"Region":        Live.Member.Region,
				"YoutubeID":     Live.Member.YoutubeID,
				"YoutubeAvatar": Live.Member.YoutubeAvatar,
			},
			"Youtube": map[string]interface{}{
				"VideoID":            Live.VideoID,
				"URL":                "https://www.youtube.com/watch?v=" + Live.VideoID,
				"Status":             Live.Status,
				"Title":              Live.Title,
				"Description":        Live.Desc,
				"Thumbnail":          Live.Thumb,
				"StartStreamSchedul": Live.Schedul,
				"EndStream":          EndStream,
				"SchedulPublished":   Live.Published,
				"Viewers":            Live.Viewers,
				"Length":             Live.Length,
			},
		}
	} else if Live.State == config.BiliLive {
		return map[string]interface{}{
			"Member": map[string]interface{}{
				"ID":             Live.Member.ID,
				"NickName":       Live.Member.Name,
				"EnName":         Live.Member.EnName,
				"JpName":         Live.Member.JpName,
				"Region":         Live.Member.Region,
				"BiliBiliRoomID": Live.Member.BiliBiliRoomID,
				"BiliBiliAvatar": Live.Member.BiliBiliAvatar,
			},
			"BiliBili": map[string]interface{}{
				"BiliBiliRoomID":     Live.Member.BiliBiliRoomID,
				"URL":                "https://live.bilibili.com/" + strconv.Itoa(Live.Member.BiliBiliRoomID),
				"Status":             Live.Status,
				"Title":              Live.Title,
				"Description":        Live.Desc,
				"Thumbnail":          Live.Thumb,
				"StartStreamSchedul": Live.Schedul,
				"Viewers":            Live.Viewers,
			},
		}
	} else {
		var StartStream interface{}
		if Live.Schedul.IsZero() {
			StartStream = nil
		} else {
			StartStream = Live.Schedul
		}
		return map[string]interface{}{
			"Member": map[string]interface{}{
				"ID":             Live.Member.ID,
				"NickName":       Live.Member.Name,
				"EnName":         Live.Member.EnName,
				"JpName":         Live.Member.JpName,
				"Region":         Live.Member.Region,
				"TwitchUserName": Live.Member.TwitchName,
				"TwitchAvatar":   Live.Member.TwitchAvatar,
			},
			"Twitch": map[string]interface{}{
				"UserName":           Live.Member.TwitchName,
				"URL":                "twitch.tv/" + Live.Member.TwitchName,
				"Status":             Live.Status,
				"Title":              Live.Title,
				"Thumbnail":          Live.Thumb,
				"StartStreamSchedul": StartStream,
				"Viewers":            Live.Viewers,
				"Game":               Live.Game,
			},
		}
	}
}
