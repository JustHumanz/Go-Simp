package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"

	log "github.com/sirupsen/logrus"
)

//StartCheckYT Youtube rss and API
func StartCheckYT(Group database.Group, Update bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if Group.YoutubeChannels != nil {
		for _, YtChan := range Group.YoutubeChannels {
			log.WithFields(log.Fields{
				"Group": Group.GroupName,
			}).Info("Checking Group channel")

			VideoID := engine.GetRSS(YtChan.YtChannel)
			for _, ID := range VideoID {
				YoutubeData, err := YtChan.CheckYoutubeVideo(ID)
				if err != nil {
					log.Error(err)
				}

				if YoutubeData == nil {
					var Thumb string
					Data, err := YtAPI([]string{ID})
					if err != nil {
						log.Error(err)
					}
					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
					}
					Items := Data.Items[0]

					_, err = network.Curl("http://i3.ytimg.com/vi/"+ID+"/maxresdefault.jpg", nil)
					if err != nil {
						Thumb = "http://i3.ytimg.com/vi/" + ID + "/hqdefault.jpg"
					} else {
						Thumb = "http://i3.ytimg.com/vi/" + ID + "/maxresdefault.jpg"
					}

					YtType := engine.YtFindType(Items.Snippet.Title)
					if YtType == "Streaming" && Items.ContentDetails.Duration != "P0D" && Items.LiveDetails.StartTime.IsZero() {
						YtType = "Regular video"
					}

					NewYoutubeData := &database.LiveStream{
						Status:       Items.Snippet.VideoStatus,
						VideoID:      ID,
						Title:        Items.Snippet.Title,
						Thumb:        Thumb,
						Desc:         Items.Snippet.Description,
						Schedul:      Items.LiveDetails.StartTime,
						Published:    Items.Snippet.PublishedAt,
						Type:         YtType,
						Viewers:      Items.Statistics.ViewCount,
						Length:       durafmt.Parse(ParseDuration(Items.ContentDetails.Duration)).String(),
						Group:        Group,
						GroupYoutube: YtChan,
						State:        config.YoutubeLive,
					}
					if Items.Snippet.VideoStatus == "none" {
						if YtType == "Covering" {
							log.WithFields(log.Fields{
								"YtID":      ID,
								"GroupName": Group.GroupName,
							}).Info("New MV or Cover")

							NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
							engine.SendLiveNotif(NewYoutubeData, Bot)

						} else if !Items.Snippet.PublishedAt.IsZero() {
							log.WithFields(log.Fields{
								"YtID":      ID,
								"GroupName": Group.GroupName,
							}).Info("Suddenly upload new video")
							if NewYoutubeData.Schedul.IsZero() {
								NewYoutubeData.UpdateSchdule(NewYoutubeData.Published)
							}

							NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
							engine.SendLiveNotif(NewYoutubeData, Bot)

						} else {
							log.WithFields(log.Fields{
								"YtID":      ID,
								"GroupName": Group.GroupName,
							}).Info("Past live stream")
							NewYoutubeData.UpdateStatus(config.PastStatus)
							engine.SendLiveNotif(NewYoutubeData, Bot)
						}
					}
				}
			}
		}
	}

	for _, Member := range Group.Members {
		if !Member.IsYtNill() && Member.Active() {
			log.WithFields(log.Fields{
				"Vtuber": Member.EnName,
				"Group":  Group.GroupName,
			}).Info("Checking Vtuber channel")

			VideoID := engine.GetRSS(Member.YoutubeID)
			for _, ID := range VideoID {
				YoutubeData, err := Member.CheckYoutubeVideo(ID)
				if err != nil {
					log.Error(err)
				}

				if YoutubeData == nil {
					var (
						Viewers string
						Thumb   string
					)

					Data, err := YtAPI([]string{ID})
					if err != nil {
						log.Error(err)
					}

					log.WithFields(log.Fields{
						"Group":  Group.GroupName,
						"Member": Member.Name,
					}).Info("Checking New VideoID")

					Items := Data.Items[0]

					if Items.Snippet.VideoStatus == config.UpcomingStatus {
						if YoutubeData == nil {
							Viewers, err = GetWaiting(ID)
							if err != nil {
								log.Error(err)
							}
						} else if YoutubeData.Viewers != config.Ytwaiting {
							Viewers = YoutubeData.Viewers
						} else {
							Viewers, err = GetWaiting(ID)
							if err != nil {
								log.Error(err)
							}
						}
					} else if Items.LiveDetails.Viewers == "" {
						Viewers = Items.Statistics.ViewCount
					} else {
						Viewers = Items.LiveDetails.Viewers
					}

					_, err = network.Curl("http://i3.ytimg.com/vi/"+ID+"/maxresdefault.jpg", nil)
					if err != nil {
						Thumb = "http://i3.ytimg.com/vi/" + ID + "/hqdefault.jpg"
					} else {
						Thumb = "http://i3.ytimg.com/vi/" + ID + "/maxresdefault.jpg"
					}

					YtType := engine.YtFindType(Items.Snippet.Title)
					if YtType == "Streaming" && Items.ContentDetails.Duration != "P0D" && Items.LiveDetails.StartTime.IsZero() {
						YtType = "Regular video"
					}

					NewYoutubeData := &database.LiveStream{
						Status:    Items.Snippet.VideoStatus,
						VideoID:   ID,
						Title:     Items.Snippet.Title,
						Thumb:     Thumb,
						Desc:      Items.Snippet.Description,
						Schedul:   Items.LiveDetails.StartTime,
						Published: Items.Snippet.PublishedAt,
						Type:      YtType,
						Viewers:   Viewers,
						Length:    durafmt.Parse(ParseDuration(Items.ContentDetails.Duration)).String(),
						Member:    Member,
						Group:     Group,
						State:     config.YoutubeLive,
					}

					if Items.Snippet.VideoStatus == config.UpcomingStatus {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.EnName,
							"Message":    "Send to notify",
						}).Info("New Upcoming live schedule")

						NewYoutubeData.UpdateStatus(config.UpcomingStatus)
						_, err := NewYoutubeData.InputYt()
						if err != nil {
							log.Error(err)
						}

						UpcominginHours := int(time.Until(NewYoutubeData.Schedul).Hours())
						if UpcominginHours > 6 {
							engine.SendLiveNotif(NewYoutubeData, Bot)
						}

					} else if Items.Snippet.VideoStatus == config.LiveStatus {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.EnName,
							"Message":    "Send to notify",
						}).Info("New live stream right now")

						NewYoutubeData.UpdateStatus(config.LiveStatus)
						_, err := NewYoutubeData.InputYt()
						if err != nil {
							log.Error(err)
						}

						if Member.BiliRoomID != 0 {
							LiveBili, err := engine.GetRoomStatus(Member.BiliRoomID)
							if err != nil {
								log.Error(err)
							}
							if LiveBili.CheckScheduleLive() {
								NewYoutubeData.SetBiliLive(true).UpdateBiliToLive()
							}
						}

						if config.GoSimpConf.Metric {
							bit, err := NewYoutubeData.MarshalBinary()
							if err != nil {
								log.Error(err)
							}
							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bit,
								State:      config.LiveStatus,
							})
						}

						if !Items.LiveDetails.ActualStartTime.IsZero() {
							NewYoutubeData.UpdateSchdule(Items.LiveDetails.ActualStartTime)
							engine.SendLiveNotif(NewYoutubeData, Bot)

						} else {
							engine.SendLiveNotif(NewYoutubeData, Bot)
						}

					} else if Items.Snippet.VideoStatus == "none" && YtType == "Covering" {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.EnName,
						}).Info("New MV or Cover")

						NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
						engine.SendLiveNotif(NewYoutubeData, Bot)

					} else if !Items.Snippet.PublishedAt.IsZero() && Items.Snippet.VideoStatus == "none" {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.EnName,
						}).Info("Suddenly upload new video")
						if NewYoutubeData.Schedul.IsZero() {
							NewYoutubeData.UpdateSchdule(NewYoutubeData.Published)
						}

						NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
						engine.SendLiveNotif(NewYoutubeData, Bot)

					} else {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.EnName,
						}).Info("Past live stream")
						NewYoutubeData.UpdateStatus(config.PastStatus)
						engine.SendLiveNotif(NewYoutubeData, Bot)
					}
				} else if Update {
					log.WithFields(log.Fields{
						"Group":   Group.GroupName,
						"Member":  Member.Name,
						"VideoID": ID,
					}).Info("Update VideoID")

					Data, err := YtAPI([]string{ID})
					if err != nil {
						log.Error(err)
					}

					Items := Data.Items[0]

					YoutubeData.UpdateEnd(Items.LiveDetails.EndTime).
						UpdateViewers(Items.Statistics.ViewCount).
						UpdateLength(durafmt.Parse(ParseDuration(Items.ContentDetails.Duration)).String()).
						SetState(config.YoutubeLive).
						AddMember(Member).
						AddGroup(Group)

					if Items.Snippet.VideoStatus == "none" && YoutubeData.Status == config.LiveStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       config.PastStatus,
						}).Info("Update video status from " + Items.Snippet.VideoStatus + " to past")
						YoutubeData.UpdateYt(config.PastStatus)

						engine.RemoveEmbed(ID, Bot)

						if config.GoSimpConf.Metric {
							bit, err := YoutubeData.MarshalBinary()
							if err != nil {
								log.Error(err)
							}
							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bit,
								State:      config.PastStatus,
							})
						}

					} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Status == config.UpcomingStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       config.LiveStatus,
						}).Info("Update video status from " + YoutubeData.Status + " to live")
						YoutubeData.UpdateStatus(config.LiveStatus)

						log.Info("Update database")
						if !Items.LiveDetails.ActualStartTime.IsZero() {
							YoutubeData.UpdateSchdule(Items.LiveDetails.ActualStartTime)
						}

						YoutubeData.UpdateYt(YoutubeData.Status)

					} else if (!Items.LiveDetails.EndTime.IsZero() && YoutubeData.Status == config.UpcomingStatus) || (YoutubeData.Status == config.UpcomingStatus && Items.Snippet.VideoStatus == "none") {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       config.PastStatus,
						}).Info("Update video status from " + Items.Snippet.VideoStatus + " to past,probably member only")
						YoutubeData.UpdateYt(config.PastStatus)

						engine.RemoveEmbed(ID, Bot)

						if config.GoSimpConf.Metric {
							bit, err := YoutubeData.MarshalBinary()
							if err != nil {
								log.Error(err)
							}
							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bit,
								State:      config.PastStatus,
							})
						}

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus && YoutubeData.Status == config.PastStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       Items.Snippet.VideoStatus,
						}).Info("maybe yt error or human error")

						YoutubeData.UpdateStatus(config.UpcomingStatus)
					} else if (Items.Snippet.VideoStatus == "none" && YoutubeData.Viewers != Items.Statistics.ViewCount) || (Items.Snippet.VideoStatus == config.LiveStatus) {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.PastStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateYt(config.PastStatus)

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus {
						if Items.LiveDetails.StartTime != YoutubeData.Schedul {
							log.WithFields(log.Fields{
								"VideoData ID": ID,
								"old schdule":  YoutubeData.Schedul,
								"new schdule":  Items.LiveDetails.StartTime,
								"Status":       config.UpcomingStatus,
							}).Info("Livestream schdule changed")

							YoutubeData.UpdateSchdule(Items.LiveDetails.StartTime)
							YoutubeData.UpdateYt(config.UpcomingStatus)
						}
					} else {
						YoutubeData.UpdateYt(YoutubeData.Status)
					}
				}
			}
		}
	}
}

//YtAPI Get data from youtube api
func YtAPI(VideoID []string) (engine.YtData, error) {
	var (
		Data engine.YtData
	)
	log.WithFields(log.Fields{
		"VideoID": VideoID,
	}).Info("Checking from youtubeAPI")

	for i, Token := range config.GoSimpConf.YtToken {
		if exTknList != nil {
			isExhaustion := false
			for _, v := range exTknList {
				if v == Token {
					isExhaustion = true
					break
				}
			}

			if isExhaustion {
				continue
			}
		}
		url := "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=" + strings.Join(VideoID, ",") + "&key=" + Token

		var bdy []byte
		var curlerr error
		if *Proxy {
			bdy, curlerr = network.CoolerCurl(url, nil)
			if curlerr != nil {
				if curlerr.Error() == "403 Forbidden" {
					exTknList = append(exTknList, Token)
				}
				log.Error(curlerr)
				if i == len(config.GoSimpConf.YtToken)-1 {
					break
				}
				continue
			}

			err := json.Unmarshal(bdy, &Data)
			if err != nil {
				return Data, err
			}
			return Data, nil

		} else {
			bdy, curlerr = network.Curl(url, nil)
			if curlerr != nil {
				if curlerr.Error() == "403 Forbidden" {
					exTknList = append(exTknList, Token)
				}
				log.Error(curlerr)
				if i == len(config.GoSimpConf.YtToken)-1 {
					break
				}
				continue
			}

			err := json.Unmarshal(bdy, &Data)
			if err != nil {
				return Data, err
			}
			return Data, nil
		}
	}
	return engine.YtData{}, errors.New("exhaustion Token")
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
