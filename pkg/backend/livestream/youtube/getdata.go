package youtube

import (
	"encoding/json"
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/tools/config"

	"github.com/hako/durafmt"

	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"

	log "github.com/sirupsen/logrus"
)

//GetRSS GetRSS from Channel
func GetRSS(YtID string) []string {
	var (
		DataXML YtXML
		VideoID []string
	)

	Data, err := network.Curl("https://www.youtube.com/feeds/videos.xml?channel_id="+YtID+"&q=searchterms", nil)
	if err != nil {
		log.Error(err, string(Data))
	}

	xml.Unmarshal(Data, &DataXML)

	for i := 0; i < len(DataXML.Entry); i++ {
		VideoID = append(VideoID, DataXML.Entry[i].VideoId)
		if i == config.BotConf.LimitConf.YoutubeLimit {
			break
		}
	}
	return VideoID
}

//StartCheckYT Youtube rss and API
func StartCheckYT(Member database.Member, Group database.Group, wg *sync.WaitGroup) error {
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

		YtData, err := Member.CheckYtVideo(VideoID[i])
		if err != nil {
			log.Error(err)
		}

		YoutubeData := &NotifStruct{
			YtData: YtData,
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
				err := engine.RemoveEmbed(VideoID[i], Bot)
				if err != nil {
					log.Error(err)
				}

			} else if Data.Items[i].Snippet.VideoStatus == "live" && YoutubeData.YtData.Status == "upcoming" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Live",
				}).Info("Update video status from " + YoutubeData.YtData.Status + " to live")
				YoutubeData.ChangeYtStatus("live").UpdateYtDB()

				log.Info("Send to notify")
				if !Data.Items[i].LiveDetails.ActualStartTime.IsZero() {
					YoutubeData.SetActuallyStart(Data.Items[i].LiveDetails.ActualStartTime).SendNude()
				} else {
					YoutubeData.SendNude()
				}

			} else if (!Data.Items[i].LiveDetails.EndTime.IsZero() && YoutubeData.YtData.Status == "upcoming") || (YoutubeData.YtData.Status == "upcoming" && Data.Items[i].Snippet.VideoStatus == "none") {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       "Past",
				}).Info("Update video status from " + Data.Items[i].Snippet.VideoStatus + " to past,probably member only")
				YoutubeData.ChangeYtStatus("past").UpdateYtDB()
				err := engine.RemoveEmbed(VideoID[i], Bot)
				if err != nil {
					log.Error(err)
				}

			} else if Data.Items[i].Snippet.VideoStatus == "upcoming" && YoutubeData.YtData.Status == "past" {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Status":       Data.Items[i].Snippet.VideoStatus,
				}).Info("maybe yt error or human error")

				YoutubeData.ChangeYtStatus("upcoming").UpdateYtDB().SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "none" && YoutubeData.YtData.Viewers != Data.Items[i].Statistics.ViewCount {
				log.WithFields(log.Fields{
					"VideoData ID": VideoID[i],
					"Viwers past":  YoutubeData.YtData.Viewers,
					"Viwers now":   Data.Items[i].Statistics.ViewCount,
					"Status":       "past",
				}).Info("Update viwers")
				YoutubeData.YtData.UpdateYt("past")

			} else if Data.Items[i].Snippet.VideoStatus == "live" {
				log.WithFields(log.Fields{
					"VideoData id": VideoID[i],
					"Viwers Live":  Data.Items[i].Statistics.ViewCount,
					"Status":       "Live",
				}).Info("Update viwers")
				YoutubeData.ChangeYtStatus("live").UpdateYtDB()

			} else if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				if Data.Items[i].LiveDetails.StartTime != YoutubeData.YtData.Schedul {
					log.WithFields(log.Fields{
						"VideoData ID": VideoID[i],
						"old schdule":  YoutubeData.YtData.Schedul,
						"new schdule":  Data.Items[i].LiveDetails.StartTime,
						"Status":       "upcoming",
					}).Info("Livestream schdule changed")

					YoutubeData.ChangeYtStatus("upcoming").
						UpYtSchedul(Data.Items[i].LiveDetails.StartTime).
						UpdateYtDB()
				}

				if time.Now().Sub(Data.Items[i].LiveDetails.StartTime) > Data.Items[i].LiveDetails.StartTime.Sub(time.Now()) && Data.Items[i].LiveDetails.ActualStartTime.IsZero() {
					if YoutubeData.YtData.Status == "upcoming" {
						YoutubeData.ChangeYtStatus("live").UpdateYtDB()

						log.Info("Send to notify")
						YoutubeData.SendNude()
						log.WithFields(log.Fields{
							"VideoData ID": VideoID[i],
							"Status":       YoutubeData.YtData.Status,
						}).Info("Livestream schedule late,change video status to live")
					}
				}
				//send to reminder
				//YoutubeData.ChangeYtStatus("reminder").SendNude()

			} else {
				YoutubeData.YtData.UpdateYt(YoutubeData.YtData.Status)
			}
		} else {
			_, err := network.Curl("http://i3.ytimg.com/vi/"+VideoID[i]+"/maxresdefault.jpg", nil)
			if err != nil {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/hqdefault.jpg"
			} else {
				Thumb = "http://i3.ytimg.com/vi/" + VideoID[i] + "/maxresdefault.jpg"
			}

			YtType := engine.YtFindType(Data.Items[i].Snippet.Title)
			if YtType == "Streaming" && Data.Items[i].ContentDetails.Duration != "P0D" && Data.Items[i].LiveDetails.StartTime.IsZero() {
				YtType = "Regular video"
			}

			YoutubeData.AddData(&database.YtDbData{
				Status:    Data.Items[i].Snippet.VideoStatus,
				VideoID:   VideoID[i],
				Title:     Data.Items[i].Snippet.Title,
				Thumb:     Thumb,
				Desc:      Data.Items[i].Snippet.Description,
				Schedul:   Data.Items[i].LiveDetails.StartTime,
				Published: Data.Items[i].Snippet.PublishedAt,
				Type:      YtType,
				Viewers:   Viewers,
				Length:    durafmt.Parse(ParseDuration(Data.Items[i].ContentDetails.Duration)).String(),
				MemberID:  Member.ID,
			})

			if Data.Items[i].Snippet.VideoStatus == "upcoming" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
					"Message":    "Send to notify",
				}).Info("New Upcoming live schedule")

				YoutubeData.ChangeYtStatus("upcoming")
				ID, err := YoutubeData.YtData.InputYt()
				if err != nil {
					log.Error(err)
				}
				YoutubeData.SetYoutubeID(ID).SendNude()

			} else if Data.Items[i].Snippet.VideoStatus == "live" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
					"Message":    "Send to notify",
				}).Info("New live stream right now")

				YoutubeData.ChangeYtStatus("live")
				ID, err := YoutubeData.YtData.InputYt()
				if err != nil {
					log.Error(err)
				}
				if !Data.Items[i].LiveDetails.ActualStartTime.IsZero() {
					YoutubeData.SetActuallyStart(Data.Items[i].LiveDetails.ActualStartTime).
						SetYoutubeID(ID).
						SendNude()
				} else {
					YoutubeData.
						SetYoutubeID(ID).SendNude()
				}

			} else if Data.Items[i].Snippet.VideoStatus == "none" && YtType == "Covering" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
				}).Info("New MV or Cover")

				YoutubeData.ChangeYtStatus("past").SendNude()
				YoutubeData.YtData.InputYt()

			} else if !Data.Items[i].Snippet.PublishedAt.IsZero() && Data.Items[i].Snippet.VideoStatus == "none" {
				log.WithFields(log.Fields{
					"YtID":       VideoID[i],
					"MemberName": Member.EnName,
				}).Info("Suddenly upload new video")
				if YoutubeData.YtData.Schedul.IsZero() {
					YoutubeData.UpYtSchedul(YoutubeData.YtData.Published)
				}

				YoutubeData.ChangeYtStatus("past").SendNude()
				YoutubeData.YtData.InputYt()

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

//YtAPI Get data from youtube api
func YtAPI(VideoID []string) (YtData, error) {
	var (
		Data YtData
	)

	body, curlerr := network.Curl("https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id="+strings.Join(VideoID, ",")+"&key="+yttoken, nil)
	if curlerr != nil {
		log.Error(curlerr)
	}
	err := json.Unmarshal(body, &Data)
	if err != nil {
		return Data, err
	}

	return Data, nil
}

//ParseDuration Parse video duration
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
