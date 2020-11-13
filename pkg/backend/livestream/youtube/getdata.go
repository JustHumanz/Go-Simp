package youtube

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	log "github.com/sirupsen/logrus"
)

//Get RSS from Channel
func GetRSS(YtID string) []string {
	var (
		DataXml YtXML
		VideoID []string
	)

	Data, err := engine.Curl("https://www.youtube.com/feeds/videos.xml?channel_id="+YtID+"&q=searchterms", nil)
	if err != nil {
		log.Error(err, string(Data))
	}

	xml.Unmarshal(Data, &DataXml)

	for i := 0; i < len(DataXml.Entry); i++ {
		VideoID = append(VideoID, DataXml.Entry[i].VideoId)
	}
	return VideoID
}

//StartCheckYT Youtube rss and API
func StartCheckYT(Member database.Name, Group database.GroupName, wg *sync.WaitGroup) error {
	VideoID := GetRSS(Member.YoutubeID)
	Data, err := YtAPI(VideoID)
	if err != nil {
		return err
	}
	defer wg.Done()

	for i := 0; i < len(Data.Items); i++ {
		var (
			Viewers string
			Thumb   string
		)

		YoutubeData := &NotifStruct{
			YtData: Member.CheckYtVideo(VideoID[i]),
			Group:  Group,
			Member: Member,
		}

		if Data.Items[i].Snippet.VideoStatus == "upcoming" {
			if YoutubeData.YtData == nil {
				Viewers, err = GetWaiting(VideoID[i])
				if err != nil {
					return err
				}
			} else if YoutubeData.YtData.Viewers != Ytwaiting {
				Viewers = YoutubeData.YtData.Viewers
			} else {
				Viewers, err = GetWaiting(VideoID[i])
				if err != nil {
					return err
				}
			}
		} else if Data.Items[i].LiveDetails.Viewers == "" {
			Viewers = Data.Items[i].Statistics.ViewCount
		} else {
			Viewers = Data.Items[i].LiveDetails.Viewers
		}

		if YoutubeData.YtData != nil {
			YoutubeData.
				UpYtView(Viewers).
				UpYtEnd(Data.Items[i].LiveDetails.EndTime).
				UpYtLen(durafmt.Parse(ParseDuration(Data.Items[i].ContentDetails.Duration)).String())

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
				YoutubeData.ChangeYtStatus("live").UpdateYtDB()
				YoutubeData.SendNude()

			} else if !Data.Items[i].LiveDetails.EndTime.IsZero() && YoutubeData.YtData.Status == "upcoming" || YoutubeData.YtData.Status == "upcoming" && Data.Items[i].Snippet.VideoStatus == "none" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Past",
				}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past,probably member only")
				YoutubeData.ChangeYtStatus("past").UpdateYtDB()

			} else if Data.Items[i].Snippet.VideoStatus == "upcoming" && YoutubeData.YtData.Status == "past" {
				log.Info("maybe yt error or human error")
				log.Info("Send to notify")
				YoutubeData.ChangeYtStatus("upcoming").UpdateYtDB()

				YoutubeData.SendNude()

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
				YoutubeData.ChangeYtStatus("reminder").SendNude()

			} else {
				YoutubeData.YtData.UpdateYt(YoutubeData.YtData.Status)
			}
		} else {
			_, err := engine.Curl("http://i3.ytimg.com/vi/"+VideoID[i]+"/maxresdefault.jpg", nil)
			if err != nil {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/hqdefault.jpg"
			} else {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/maxresdefault.jpg"
			}

			yttype := engine.YtFindType(Data.Items[i].Snippet.Title)
			if yttype == "Streaming" && Data.Items[i].ContentDetails.Duration != "P0D" {
				yttype = "Regular video"
			}

			var (
				timestart time.Time
			)

			if !Data.Items[i].LiveDetails.StartTime.IsZero() {
				timestart = Data.Items[i].LiveDetails.StartTime
			} else if !Data.Items[i].Snippet.PublishedAt.IsZero() && Data.Items[i].LiveDetails.StartTime.IsZero() {
				timestart = Data.Items[i].Snippet.PublishedAt
			} else if Data.Items[i].LiveDetails.StartTime.IsZero() && Data.Items[i].Snippet.PublishedAt.IsZero() {
				timestart = time.Now()
			}

			//verify
			YoutubeData.AddData(&database.YtDbData{
				Status:    Data.Items[i].Snippet.VideoStatus,
				VideoID:   VideoID[i],
				Title:     Data.Items[i].Snippet.Title,
				Thumb:     Thumb,
				Desc:      Data.Items[i].Snippet.Description,
				Schedul:   timestart,
				Published: Data.Items[i].Snippet.PublishedAt,
				Type:      yttype,
				Viewers:   Viewers,
				Length:    durafmt.Parse(ParseDuration(Data.Items[i].ContentDetails.Duration)).String(),
			})

			if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
					"Message":    "Send to notify",
				}).Info("New Upcoming live schedule")

				err := YoutubeData.ChangeYtStatus("upcoming").SendtoDB()
				if err != nil {
					return err
				}
				YoutubeData.SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "live" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
					"Message":    "Send to notify",
				}).Info("New live stream right now")

				err := YoutubeData.ChangeYtStatus("live").SendtoDB()
				if err != nil {
					return err
				}
				YoutubeData.SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "none" && yttype == "Covering" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
				}).Info("New MV or Cover")

				err := YoutubeData.ChangeYtStatus("past").SendtoDB()
				if err != nil {
					return err
				}
				YoutubeData.SendNude()

			} else if !Data.Items[i].Snippet.PublishedAt.IsZero() && Data.Items[i].Snippet.VideoStatus == "none" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
				}).Info("Suddenly upload new video")
				if YoutubeData.YtData.Schedul.IsZero() {
					YoutubeData.UpYtSchedul(YoutubeData.YtData.Published)
				}
				err := YoutubeData.ChangeYtStatus("past").SendtoDB()
				if err != nil {
					return err
				}
				YoutubeData.SendNude()

			} else {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
				}).Info("Past live stream")
				YoutubeData.ChangeYtStatus("past").SendNude()
			}
		}
	}
	return nil
}

//Get data from youtube api
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

//Parse video duration
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
