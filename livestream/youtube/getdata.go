package youtube

import (
	"encoding/json"
	"encoding/xml"
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
	funcvar := engine.GetFunctionName(GetRSS)
	engine.Debugging(funcvar, "In", YtID)
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
	engine.Debugging(funcvar, "Out", VideoID)
	return VideoID
}

func Filter(Name database.Name, Group database.GroupName, wg *sync.WaitGroup) {
	funcvar := engine.GetFunctionName(Filter)
	engine.Debugging(funcvar, "In", Name)
	defer wg.Done()
	for _, YoutubeID := range strings.Split(Name.YoutubeID, "\n") {
		VideoID := GetRSS(YoutubeID)
		Data, err := YtAPI(VideoID)
		if err != nil {
			log.Error(err)
		}
		MemberFixName := engine.FixName(Name.EnName, Name.JpName)
		for i := 0; i < len(Data.Items); i++ {
			var (
				yttype    string
				PushData  NotifStruct
				Viewers   string
				Starttime time.Time
				Thumb     string
			)
			duration := durafmt.Parse(ParseDuration(Data.Items[i].ContentDetails.Duration))
			if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|mv|covered|op)", strings.ToLower(Data.Items[i].Snippet.Title)); Cover {
				yttype = "Covering"
			} else if Chat, _ := regexp.MatchString("(?m)(free|chat|room)", Data.Items[i].Snippet.Title); Chat {
				yttype = "ChatRoom"
			} else {
				yttype = "Streaming"
			}
			DataDB := database.CheckVideoID(VideoID[i])

			if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				Viewers, err = GetWaiting(VideoID[i])
				if err != nil {
					log.Error(err)
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

			_, err := engine.Curl("http://i3.ytimg.com/vi/"+VideoID[i]+"/maxresdefault.jpg", nil)
			if err != nil {
				log.Error(err)
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/hqdefault.jpg"
			} else {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/maxresdefault.jpg"
			}

			PushData = NotifStruct{
				Group:  Group,
				Member: Name,
			}
			if DataDB != (database.YtDbData{}) {
				DataDB.Viewers = Viewers
				DataDB.End = Data.Items[i].LiveDetails.EndTime
				DataDB.Length = duration.String()

				PushData.Data = DataDB
				PushData.Data.VideoID = VideoID[i]
				if Data.Items[i].Snippet.VideoStatus == "none" && DataDB.Status == "live" {
					log.WithFields(log.Fields{
						"VideoData ID": VideoID[i],
						"Status":       "Past",
					}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past")
					DataDB.UpdateYt("past")
				} else if Data.Items[i].Snippet.VideoStatus == "live" && DataDB.Status == "upcoming" {
					log.WithFields(log.Fields{
						"VideoData ID": VideoID[i],
						"Status":       "Live",
					}).Info("Update video status from " + DataDB.Status + " to live")
					DataDB.UpdateYt("live")

					log.Info("Send to notify")
					PushData.GetEmbed("live").SendNude()
				} else if !Data.Items[i].LiveDetails.EndTime.IsZero() && DataDB.Status == "upcoming" {
					log.WithFields(log.Fields{
						"VideoData ID": VideoID[i],
						"Status":       "Past",
					}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past,probably member only")
					DataDB.UpdateYt("past")

				} else if Data.Items[i].Snippet.VideoStatus == "upcoming" && DataDB.Status == "past" {
					log.Info("maybe yt error or human error")
					DataDB.UpdateYt("upcoming")

					log.Info("Send to notify")
					PushData.GetEmbed("upcoming").SendNude()
				} else if Data.Items[i].Snippet.VideoStatus == "none" && DataDB.Viewers != Data.Items[i].Statistics.ViewCount {
					log.WithFields(log.Fields{
						"VideoData ID": VideoID[i],
						"Viwers past":  DataDB.Viewers,
						"Viwers now":   Data.Items[i].Statistics.ViewCount,
						"Status":       "Past",
					}).Info("Update viwers")
					DataDB.UpdateYt("past")
				} else if Data.Items[i].Snippet.VideoStatus == "live" {
					log.WithFields(log.Fields{
						"VideoData id": VideoID[i],
						"Viwers Live":  Data.Items[i].Statistics.ViewCount,
						"Status":       "Live",
					}).Info("Update viwers")
					DataDB.UpdateYt("live")

				} else if Data.Items[i].Snippet.VideoStatus == "upcoming" {
					if Data.Items[i].LiveDetails.StartTime != PushData.Data.Schedul {
						DataDB.Schedul = Data.Items[i].LiveDetails.StartTime
						log.Info("Livestream schdule changed")
						DataDB.UpdateYt("upcoming")

						log.Info("Send to notify")
						PushData.GetEmbed("upcoming").SendNude()
					}
					//send to reminder
					loc := Zawarudo(DataDB.Region)
					UpcominginMinutes := int(math.Round(PushData.Data.Schedul.In(loc).Sub(time.Now().In(loc)).Minutes()))
					if UpcominginMinutes > 60 && UpcominginMinutes < 66 || UpcominginMinutes > 30 && UpcominginMinutes < 36 {
						PushData.GetEmbed("reminder").SendNude()
					}
				} else {
					DataDB.UpdateYt(DataDB.Status)
				}
			} else {
				//verify
				PushData.Data = database.YtDbData{
					Status:    Data.Items[i].Snippet.VideoStatus,
					VideoID:   VideoID[i],
					Title:     Data.Items[i].Snippet.Title,
					Thumb:     Thumb,
					Desc:      Data.Items[i].Snippet.Description,
					Schedul:   Starttime,
					Published: Data.Items[i].Snippet.PublishedAt,
					Type:      yttype,
					Viewers:   Viewers,
				}
				if Data.Items[i].Snippet.VideoStatus == "upcoming" {
					PushData.Data.InputYt(Name.ID)
					log.WithFields(log.Fields{
						"YtID":       VideoID[i],
						"MemberName": MemberFixName,
					}).Info("New Upcoming live schedule")

					log.Info("Send to notify")
					PushData.GetEmbed("upcoming").SendNude()

				} else if Data.Items[i].Snippet.VideoStatus == "live" {
					log.WithFields(log.Fields{
						"YtID":       VideoID[i],
						"MemberName": MemberFixName,
					}).Info("New live stream right now")
					PushData.Data.InputYt(Name.ID)

					log.Info("Send to notify")
					PushData.GetEmbed("live").SendNude()

				} else if Data.Items[i].Snippet.VideoStatus == "none" && yttype == "Covering" {
					PushData.Data.Status = "past"
					PushData.Data.InputYt(Name.ID)

					log.WithFields(log.Fields{
						"YtID":       VideoID[i],
						"MemberName": MemberFixName,
					}).Info("New MV or Cover")
					PushData.GetEmbed("past").SendNude()

				} else if !Data.Items[i].Snippet.PublishedAt.IsZero() && Data.Items[i].Snippet.VideoStatus == "none" {
					PushData.Data.Status = "past"
					PushData.Data.InputYt(Name.ID)
					log.WithFields(log.Fields{
						"YtID":       VideoID[i],
						"MemberName": MemberFixName,
					}).Info("Suddenly upload new video")
					if PushData.Data.Schedul.IsZero() {
						PushData.Data.Schedul = PushData.Data.Published
					}
					PushData.GetEmbed("past").SendNude()

				} else {
					log.WithFields(log.Fields{
						"YtID":       VideoID[i],
						"MemberName": MemberFixName,
					}).Info("Past live stream")
					PushData.GetEmbed("past").SendNude()
				}
			}
		}
	}
}

func YtAPI(VideoID []string) (YtData, error) {
	funcvar := engine.GetFunctionName(YtAPI)
	engine.Debugging(funcvar, "In", VideoID)
	var (
		Data    YtData
		body    []byte
		curlerr error
		urls    = "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=" + strings.Join(VideoID, ",") + "&key=" + yttoken
	)

	body, curlerr = engine.Curl(urls, nil)
	if curlerr != nil {
		log.Error(curlerr, string(body))
		oldtoken := yttoken
		yttoken = engine.ChangeToken(yttoken)
		log.WithFields(log.Fields{
			"Old Token": oldtoken,
			"New Token": yttoken,
		}).Warn("Token out of limit,try to change")

	}
	err := json.Unmarshal(body, &Data)
	engine.BruhMoment(err, "", false)

	engine.Debugging(funcvar, "Out", fmt.Sprint(Data, nil))
	return Data, nil
}

func getXML(url string) ([]byte, error) {
	funcvar := engine.GetFunctionName(getXML)
	engine.Debugging(funcvar, "In", url)
	resp, err := http.Get(url)
	engine.BruhMoment(err, "", false)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	engine.BruhMoment(err, "", false)

	engine.Debugging(funcvar, "Out", fmt.Sprint(data, nil))
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
