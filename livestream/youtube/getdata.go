package youtube

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	log "github.com/sirupsen/logrus"
)

func GetRSS(YtID string) []string {
	var DataXml YtXML

	Data, err := engine.Curl("https://www.youtube.com/feeds/videos.xml?channel_id="+YtID+"&q=searchterms", nil)
	if err != nil {
		log.Error(err, string(Data))
	}

	xml.Unmarshal(Data, &DataXml)

	var VideoID []string
	for i := 0; i < len(DataXml.Entry); i++ {
		VideoID = append(VideoID, DataXml.Entry[i].VideoId)
	}
	return VideoID
}

func Filter(Name database.Name, Group database.GroupName, wg *sync.WaitGroup) error {
	VideoID := GetRSS(Name.YoutubeID)
	Data, err := YtAPI(VideoID)
	if err != nil {
		return err
	}
	defer wg.Done()

	for i := 0; i < len(Data.Items); i++ {
		var (
			yttype    string
			Viewers   string
			Starttime time.Time
			Thumb     string
		)
		duration := durafmt.Parse(ParseDuration(Data.Items[i].ContentDetails.Duration))
		if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|mv|covered|op|ed)", strings.ToLower(Data.Items[i].Snippet.Title)); Cover {
			yttype = "Covering"
		} else if Chat, _ := regexp.MatchString("(?m)(chat|room)", Data.Items[i].Snippet.Title); Chat {
			yttype = "ChatRoom"
		} else {
			yttype = "Streaming"
		}
		YoutubeData := &NotifStruct{
			YtData: Name.CheckYtVideo(VideoID[i]),
			Group:  Group,
			Member: Name,
		}

		if Data.Items[i].Snippet.VideoStatus == "upcoming" {
			if YoutubeData.YtData == nil {
				Viewers, err = GetWaiting(VideoID[i])
				if err != nil {
					log.Error(err)
				}
			} else if YoutubeData.YtData.Viewers != Ytwaiting {
				Viewers = YoutubeData.YtData.Viewers
			} else {
				Viewers, err = GetWaiting(VideoID[i])
				if err != nil {
					log.Error(err)
				}
			}
		} else if Data.Items[i].LiveDetails.Viewers == "" {
			Viewers = Data.Items[i].Statistics.ViewCount
		} else {
			Viewers = Data.Items[i].LiveDetails.Viewers
		}

		if Data.Items[i].LiveDetails.StartTime.IsZero() {
			Starttime = Data.Items[i].LiveDetails.ActualStartTime
		} else if !Data.Items[i].LiveDetails.StartTime.IsZero() {
			Starttime = Data.Items[i].LiveDetails.StartTime
		} else {
			Starttime = Data.Items[i].Snippet.PublishedAt
		}

		if YoutubeData.YtData != nil {
			YoutubeData.
				UpYtView(Viewers).
				UpYtEnd(Data.Items[i].LiveDetails.EndTime).
				UpYtLen(duration.String())

			if Data.Items[i].Snippet.VideoStatus == "none" && YoutubeData.YtData.Status == "live" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Past",
				}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past")
				YoutubeData.ChangeYtStatus("past").UpdateYtDB()

			} else if Data.Items[i].Snippet.VideoStatus == "live" && YoutubeData.YtData.Status == "upcoming" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Live",
				}).Info("Update video status from " + YoutubeData.YtData.Status + " to live")
				log.Info("Send to notify")
				YoutubeData.ChangeYtStatus("live").SendtoDB().SendNude()

			} else if !Data.Items[i].LiveDetails.EndTime.IsZero() && YoutubeData.YtData.Status == "upcoming" || YoutubeData.YtData.Status == "upcoming" && Data.Items[i].Snippet.VideoStatus == "none" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Past",
				}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past,probably member only")
				YoutubeData.YtData.UpdateYt("past")

			} else if Data.Items[i].Snippet.VideoStatus == "upcoming" && YoutubeData.YtData.Status == "past" {
				log.Info("maybe yt error or human error")
				log.Info("Send to notify")
				YoutubeData.ChangeYtStatus("upcoming").SendtoDB().SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "none" && YoutubeData.YtData.Viewers != Data.Items[i].Statistics.ViewCount {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Viwers past":  YoutubeData.YtData.Viewers,
					"Viwers now":   Data.Items[i].Statistics.ViewCount,
					"Status":       "Past",
				}).Info("Update viwers")
				YoutubeData.YtData.UpdateYt("live")

			} else if Data.Items[i].Snippet.VideoStatus == "live" {
				log.WithFields(log.Fields{
					"VideoData id": VideoID[i],
					"Viwers Live":  Data.Items[i].Statistics.ViewCount,
					"Status":       "Live",
				}).Info("Update viwers")
				YoutubeData.ChangeYtStatus("live").UpdateYtDB()

			} else if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				if Data.Items[i].LiveDetails.StartTime != YoutubeData.YtData.Schedul {
					log.Info("Livestream schdule changed")
					log.Info("Send to notify")

					YoutubeData.ChangeYtStatus("upcoming").
						UpYtSchedul(Data.Items[i].LiveDetails.StartTime).SendNude()

					YoutubeData.UpdateYtDB()
				}
				//send to reminder
				loc := engine.Zawarudo(YoutubeData.YtData.Region)
				UpcominginMinutes := int(math.Round(YoutubeData.YtData.Schedul.In(loc).Sub(time.Now().In(loc)).Minutes()))
				if UpcominginMinutes > 60 && UpcominginMinutes < 66 || UpcominginMinutes > 30 && UpcominginMinutes < 36 {
					YoutubeData.ChangeYtStatus("reminder").SendNude()
				}
			} else {
				YoutubeData.YtData.UpdateYt(YoutubeData.YtData.Status)
			}
		} else {
			MemberFixName := engine.FixName(Name.EnName, Name.JpName)
			_, err := engine.Curl("http://i3.ytimg.com/vi/"+VideoID[i]+"/maxresdefault.jpg", nil)
			if err != nil {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/hqdefault.jpg"
			} else {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/maxresdefault.jpg"
			}

			//verify
			YoutubeData.AddData(&database.YtDbData{
				Status:    Data.Items[i].Snippet.VideoStatus,
				VideoID:   VideoID[i],
				Title:     Data.Items[i].Snippet.Title,
				Thumb:     Thumb,
				Desc:      Data.Items[i].Snippet.Description,
				Schedul:   Starttime,
				Published: Data.Items[i].Snippet.PublishedAt,
				Type:      yttype,
				Viewers:   Viewers,
			})

			if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": MemberFixName,
					"Message":    "Send to notify",
				}).Info("New Upcoming live schedule")

				YoutubeData.ChangeYtStatus("upcoming").SendtoDB().SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "live" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": MemberFixName,
					"Message":    "Send to notify",
				}).Info("New live stream right now")
				YoutubeData.ChangeYtStatus("live").SendtoDB().SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "none" && yttype == "Covering" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": MemberFixName,
				}).Info("New MV or Cover")

				YoutubeData.ChangeYtStatus("past").SendtoDB().SendNude()

			} else if !Data.Items[i].Snippet.PublishedAt.IsZero() && Data.Items[i].Snippet.VideoStatus == "none" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": MemberFixName,
				}).Info("Suddenly upload new video")
				if YoutubeData.YtData.Schedul.IsZero() {
					YoutubeData.UpYtSchedul(YoutubeData.YtData.Published)
				}
				YoutubeData.ChangeYtStatus("past").SendtoDB().SendNude()

			} else {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": MemberFixName,
				}).Info("Past live stream")
				YoutubeData.ChangeYtStatus("past").SendNude()
			}
		}
	}
	return nil
}

func YtAPI(VideoID []string) (YtData, error) {
	var (
		Data    YtData
		body    []byte
		curlerr error
		urls    = "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=" + strings.Join(VideoID, ",") + "&key=" + yttoken
	)

	body, curlerr = engine.Curl(urls, nil)
	if curlerr != nil {
		return YtData{}, errors.New("Token out of limit")
	}
	err := json.Unmarshal(body, &Data)
	engine.BruhMoment(err, "", false)

	return Data, nil
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	engine.BruhMoment(err, "", false)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	engine.BruhMoment(err, "", false)

	return data, nil
}

func ParseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := ParseInt64(matches[1])
	months := ParseInt64(matches[2])
	days := ParseInt64(matches[3])
	hours := ParseInt64(matches[4])
	minutes := ParseInt64(matches[5])
	seconds := ParseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func ParseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}
