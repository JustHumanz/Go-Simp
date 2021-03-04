package youtube

import (
	"encoding/json"
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"

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
		if i == configfile.LimitConf.YoutubeLimit {
			break
		}
	}
	return VideoID
}

//StartCheckYT Youtube rss and API
func StartCheckYT(Group database.Group, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, Member := range Group.Members {
		if Member.YoutubeID != "" {
			VideoID := GetRSS(Member.YoutubeID)
			Data, err := YtAPI(VideoID)
			if err != nil {
				log.Error(err)
			}

			log.WithFields(log.Fields{
				"Group":  Group.GroupName,
				"Member": Member.Name,
			}).Info("Checking Youtube Channels")

			for i, Items := range Data.Items {
				var (
					Viewers   string
					Thumb     string
					YtVideoID = VideoID[i]
				)

				YtData, err := Member.CheckYtVideo(YtVideoID)
				if err != nil {
					log.Error(err)
				}

				YoutubeData := &NotifStruct{
					YtData: YtData,
					Group:  Group,
					Member: Member,
				}

				if Items.Snippet.VideoStatus == "upcoming" {
					if YoutubeData.YtData == nil {
						Viewers, err = GetWaiting(YtVideoID)
						if err != nil {
							log.Error(err)
						}
					} else if YoutubeData.YtData.Viewers != Ytwaiting {
						Viewers = YoutubeData.YtData.Viewers
					} else {
						Viewers, err = GetWaiting(YtVideoID)
						if err != nil {
							log.Error(err)
						}
					}
				} else if Items.LiveDetails.Viewers == "" {
					Viewers = Items.Statistics.ViewCount
				} else {
					Viewers = Items.LiveDetails.Viewers
				}

				if YoutubeData.YtData != nil {
					YoutubeData.
						UpYtView(Viewers).
						UpYtEnd(Items.LiveDetails.EndTime).
						UpYtLen(durafmt.Parse(ParseDuration(Items.ContentDetails.Duration)).String())

					if Items.Snippet.VideoStatus == "none" && YoutubeData.YtData.Status == "live" {
						log.WithFields(log.Fields{
							"VideoData ID": YtVideoID,
							"Status":       "Past",
						}).Info("Update video status from " + Items.Snippet.VideoStatus + " to past")
						YoutubeData.ChangeYtStatus("past").UpdateYtDB()
						engine.RemoveEmbed(YtVideoID, Bot)

					} else if Items.Snippet.VideoStatus == "live" && YoutubeData.YtData.Status == "upcoming" {
						log.WithFields(log.Fields{
							"VideoData ID": YtVideoID,
							"Status":       "Live",
						}).Info("Update video status from " + YoutubeData.YtData.Status + " to live")
						YoutubeData.ChangeYtStatus("live")

						log.Info("Send to notify")
						if !Items.LiveDetails.ActualStartTime.IsZero() {
							YoutubeData.SetActuallyStart(Items.LiveDetails.ActualStartTime)
						} else {
							//YoutubeData.SendNude()
						}

						YoutubeData.UpdateYtDB()

					} else if (!Items.LiveDetails.EndTime.IsZero() && YoutubeData.YtData.Status == "upcoming") || (YoutubeData.YtData.Status == "upcoming" && Items.Snippet.VideoStatus == "none") {
						log.WithFields(log.Fields{
							"VideoData ID": YtVideoID,
							"Status":       "Past",
						}).Info("Update video status from " + Items.Snippet.VideoStatus + " to past,probably member only")
						YoutubeData.ChangeYtStatus("past").UpdateYtDB()
						engine.RemoveEmbed(YtVideoID, Bot)

					} else if Items.Snippet.VideoStatus == "upcoming" && YoutubeData.YtData.Status == "past" {
						log.WithFields(log.Fields{
							"VideoData ID": YtVideoID,
							"Status":       Items.Snippet.VideoStatus,
						}).Info("maybe yt error or human error")

						YoutubeData.ChangeYtStatus("upcoming").UpdateYtDB().SendNude()

					} else if Items.Snippet.VideoStatus == "none" && YoutubeData.YtData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": YtVideoID,
							"Viwers past":  YoutubeData.YtData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       "past",
						}).Info("Update Viwers")
						YoutubeData.YtData.UpdateYt("past")

					} else if Items.Snippet.VideoStatus == "live" {
						log.WithFields(log.Fields{
							"VideoData id": YtVideoID,
							"Viwers Live":  Items.Statistics.ViewCount,
							"Status":       "Live",
						}).Info("Update Viwers")
						YoutubeData.ChangeYtStatus("live").UpdateYtDB()

					} else if Items.Snippet.VideoStatus == "upcoming" {
						if Items.LiveDetails.StartTime != YoutubeData.YtData.Schedul {
							log.WithFields(log.Fields{
								"VideoData ID": YtVideoID,
								"old schdule":  YoutubeData.YtData.Schedul,
								"new schdule":  Items.LiveDetails.StartTime,
								"Status":       "upcoming",
							}).Info("Livestream schdule changed")

							YoutubeData.ChangeYtStatus("upcoming").
								UpYtSchedul(Items.LiveDetails.StartTime).
								UpdateYtDB()
						}
					} else {
						YoutubeData.YtData.UpdateYt(YoutubeData.YtData.Status)
					}
				} else {
					_, err := network.Curl("http://i3.ytimg.com/vi/"+YtVideoID+"/maxresdefault.jpg", nil)
					if err != nil {
						Thumb = "http://i3.ytimg.com/vi/" + YtVideoID + "/hqdefault.jpg"
					} else {
						Thumb = "http://i3.ytimg.com/vi/" + YtVideoID + "/maxresdefault.jpg"
					}

					YtType := engine.YtFindType(Items.Snippet.Title)
					if YtType == "Streaming" && Items.ContentDetails.Duration != "P0D" && Items.LiveDetails.StartTime.IsZero() {
						YtType = "Regular video"
					}

					YoutubeData.AddData(&database.YtDbData{
						Status:    Items.Snippet.VideoStatus,
						VideoID:   YtVideoID,
						Title:     Items.Snippet.Title,
						Thumb:     Thumb,
						Desc:      Items.Snippet.Description,
						Schedul:   Items.LiveDetails.StartTime,
						Published: Items.Snippet.PublishedAt,
						Type:      YtType,
						Viewers:   Viewers,
						Length:    durafmt.Parse(ParseDuration(Items.ContentDetails.Duration)).String(),
						MemberID:  Member.ID,
					})

					if Items.Snippet.VideoStatus == "upcoming" {
						log.WithFields(log.Fields{
							"YtID":       YtVideoID,
							"MemberName": Member.EnName,
							"Message":    "Send to notify",
						}).Info("New Upcoming live schedule")

						YoutubeData.ChangeYtStatus("upcoming")
						ID, err := YoutubeData.YtData.InputYt()
						if err != nil {
							log.Error(err)
						}
						YoutubeData.SetYoutubeID(ID).SendNude()

					} else if Items.Snippet.VideoStatus == "live" {
						log.WithFields(log.Fields{
							"YtID":       YtVideoID,
							"MemberName": Member.EnName,
							"Message":    "Send to notify",
						}).Info("New live stream right now")

						YoutubeData.ChangeYtStatus("live")
						ID, err := YoutubeData.YtData.InputYt()
						if err != nil {
							log.Error(err)
						}
						if !Items.LiveDetails.ActualStartTime.IsZero() {
							YoutubeData.SetActuallyStart(Items.LiveDetails.ActualStartTime).
								SetYoutubeID(ID).
								SendNude()
						} else {
							YoutubeData.
								SetYoutubeID(ID).SendNude()
						}

					} else if Items.Snippet.VideoStatus == "none" && YtType == "Covering" {
						log.WithFields(log.Fields{
							"YtID":       YtVideoID,
							"MemberName": Member.EnName,
						}).Info("New MV or Cover")

						YoutubeData.ChangeYtStatus("past").SendNude()
						YoutubeData.YtData.InputYt()

					} else if !Items.Snippet.PublishedAt.IsZero() && Items.Snippet.VideoStatus == "none" {
						log.WithFields(log.Fields{
							"YtID":       YtVideoID,
							"MemberName": Member.EnName,
						}).Info("Suddenly upload new video")
						if YoutubeData.YtData.Schedul.IsZero() {
							YoutubeData.UpYtSchedul(YoutubeData.YtData.Published)
						}

						YoutubeData.ChangeYtStatus("past").SendNude()
						YoutubeData.YtData.InputYt()

					} else {
						log.WithFields(log.Fields{
							"YtID":       YtVideoID,
							"MemberName": Member.EnName,
						}).Info("Past live stream")
						YoutubeData.ChangeYtStatus("past").SendNude()
					}
				}
			}

		}
		time.Sleep(1 * time.Second)
	}
}

//YtAPI Get data from youtube api
func YtAPI(VideoID []string) (YtData, error) {
	var (
		Data YtData
	)

	body, curlerr := network.Curl("https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id="+strings.Join(VideoID, ",")+"&key="+*yttoken, nil)
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
