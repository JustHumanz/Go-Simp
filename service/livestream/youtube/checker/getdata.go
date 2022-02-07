package main

import (
	"context"
	"fmt"
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

	//check vtuber agency youtube channel
	if Group.YoutubeChannels != nil {
		for _, YtChan := range Group.YoutubeChannels {
			log.WithFields(log.Fields{
				"agency":    Group.GroupName,
				"channelID": YtChan.YtChannel,
				"region":    YtChan.Region,
			}).Info("Checking agency channel")

			VideoID, err := engine.GetRSS(YtChan.YtChannel, *proxy)
			if err != nil {
				log.Error(err)
				gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
					Message: err.Error(),
					Service: ModuleState,
				})
			}
			for _, ID := range VideoID {
				YoutubeData, err := YtChan.CheckYoutubeVideo(ID)
				if err != nil {
					log.Warn(err)
				}

				if YoutubeData == nil {
					Data, err := engine.YtAPI([]string{ID})
					if err != nil {
						log.Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}

					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
						continue
					}
					Items := Data.Items[0]

					YtType := engine.YtFindType(Items.Snippet.Title)
					if YtType == "Streaming" && Items.ContentDetails.Duration != "P0D" && Items.LiveDetails.StartTime.IsZero() {
						YtType = "Regular video"
					}

					NewYoutubeData := &database.LiveStream{
						Status:  Items.Snippet.VideoStatus,
						VideoID: ID,
						Title:   Items.Snippet.Title,
						Thumb: func() string {
							_, err = network.Curl("http://i3.ytimg.com/vi/"+ID+"/maxresdefault.jpg", nil)
							if err != nil {
								return "http://i3.ytimg.com/vi/" + ID + "/hqdefault.jpg"
							} else {
								return "http://i3.ytimg.com/vi/" + ID + "/maxresdefault.jpg"
							}

						}(),
						Desc:         Items.Snippet.Description,
						Schedul:      Items.LiveDetails.StartTime,
						Published:    Items.Snippet.PublishedAt,
						Type:         YtType,
						Viewers:      Items.Statistics.ViewCount,
						Length:       durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String(),
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
					} else if Items.Snippet.VideoStatus == config.UpcomingStatus {
						log.WithFields(log.Fields{
							"YtID":      ID,
							"GroupName": Group.GroupName,
							"Message":   "Send to notify",
						}).Info("New Upcoming live schedule")

						NewYoutubeData.UpdateStatus(config.UpcomingStatus)
						_, err := NewYoutubeData.InputYt()
						if err != nil {
							log.Error(err)
						}

						err = NewYoutubeData.SendToCache(true)
						if err != nil {
							log.Error(err)
						}

						engine.SendLiveNotif(NewYoutubeData, Bot)
					}
				} else if Update {
					log.WithFields(log.Fields{
						"Group":   Group.GroupName,
						"VideoID": ID,
					}).Info("Update VideoID")

					Data, err := engine.YtAPI([]string{ID})
					if err != nil {
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
						log.Error(err)
					}

					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
						continue
					}

					Items := Data.Items[0]

					YoutubeData.UpdateEnd(Items.LiveDetails.EndTime).
						UpdateViewers(Items.Statistics.ViewCount).
						UpdateLength(durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String()).
						SetState(config.YoutubeLive).
						AddGroup(Group)

					if Items.Snippet.VideoStatus == "none" && YoutubeData.Status == config.LiveStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       config.PastStatus,
						}).Info("Update video status from " + Items.Snippet.VideoStatus + " to past")
						YoutubeData.UpdateGroupYt(config.PastStatus)

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

						YoutubeData.UpdateGroupYt(YoutubeData.Status)
						engine.SendLiveNotif(YoutubeData, Bot)

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus && YoutubeData.Status == config.PastStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       Items.Snippet.VideoStatus,
						}).Info("maybe yt error or human error")

						YoutubeData.UpdateStatus(config.UpcomingStatus)
					} else if Items.Snippet.VideoStatus == "none" && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.PastStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateGroupYt(config.PastStatus)

					} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.LiveStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateGroupYt(config.LiveStatus)

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus {
						if Items.LiveDetails.StartTime != YoutubeData.Schedul {
							log.WithFields(log.Fields{
								"VideoData ID": ID,
								"old schdule":  YoutubeData.Schedul,
								"new schdule":  Items.LiveDetails.StartTime,
								"Status":       config.UpcomingStatus,
							}).Info("Livestream schdule changed")

							YoutubeData.UpdateSchdule(Items.LiveDetails.StartTime)
							YoutubeData.UpdateGroupYt(config.UpcomingStatus)
						}
					}
				}
			}
		}
	}

	//check vtuber agency members youtube channel
	for _, Member := range Group.Members {
		if !Member.IsYtNill() && Member.Active() {
			log.WithFields(log.Fields{
				"Vtuber": Member.Name,
				"Group":  Group.GroupName,
			}).Info("Checking Vtuber channel")

			VideoID, err := engine.GetRSS(Member.YoutubeID, *proxy)
			if err != nil {
				log.Error(err)
				gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
					Message: err.Error(),
					Service: ModuleState,
				})
			}
			for _, ID := range VideoID {
				YoutubeData, err := Member.CheckYoutubeVideo(ID)
				if err != nil {
					log.Warn(err)
				}

				if YoutubeData == nil {
					var (
						Viewers string
					)

					Data, err := engine.YtAPI([]string{ID})
					if err != nil {
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
						log.Error(err)
					}

					log.WithFields(log.Fields{
						"Group":  Group.GroupName,
						"Member": Member.Name,
					}).Info("Checking New VideoID")

					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
						pilot.ReportDeadService("Yt Item Nill", ModuleState)
						continue
					}

					Items := Data.Items[0]

					if Items.Snippet.VideoStatus == config.UpcomingStatus {
						if YoutubeData == nil {
							Viewers, err = engine.GetWaiting(ID)
							if err != nil {
								log.Error(err)
							}
						} else if YoutubeData.Viewers != config.Ytwaiting {
							Viewers = YoutubeData.Viewers
						} else {
							Viewers, err = engine.GetWaiting(ID)
							if err != nil {
								log.Error(err)
							}
						}
					} else if Items.LiveDetails.Viewers == "" {
						Viewers = Items.Statistics.ViewCount
					} else {
						Viewers = Items.LiveDetails.Viewers
					}

					YtType := engine.YtFindType(Items.Snippet.Title)
					if YtType == "Streaming" && Items.ContentDetails.Duration != "P0D" && Items.LiveDetails.StartTime.IsZero() {
						YtType = "Regular video"
					}

					NewYoutubeData := &database.LiveStream{
						Status:  Items.Snippet.VideoStatus,
						VideoID: ID,
						Title:   Items.Snippet.Title,
						Thumb: func() string {
							_, err = network.Curl("http://i3.ytimg.com/vi/"+ID+"/maxresdefault.jpg", nil)
							if err != nil {
								return "http://i3.ytimg.com/vi/" + ID + "/hqdefault.jpg"
							} else {
								return "http://i3.ytimg.com/vi/" + ID + "/maxresdefault.jpg"
							}

						}(),
						Desc:      Items.Snippet.Description,
						Schedul:   Items.LiveDetails.StartTime,
						Published: Items.Snippet.PublishedAt,
						Type:      YtType,
						Viewers:   Viewers,
						Length:    durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String(),
						Member:    Member,
						Group:     Group,
						State:     config.YoutubeLive,
					}

					if Items.Snippet.VideoStatus == config.UpcomingStatus {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.Name,
							"Message":    "Send to notify",
						}).Info("New Upcoming live schedule")

						NewYoutubeData.UpdateStatus(config.UpcomingStatus)
						_, err := NewYoutubeData.InputYt()
						if err != nil {
							log.Error(err)
						}

						err = NewYoutubeData.SendToCache(false)
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
							"MemberName": Member.Name,
							"Message":    "Send to notify",
						}).Info("Suddenly live stream")

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
							"MemberName": Member.Name,
						}).Info("New MV or Cover")

						NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
						engine.SendLiveNotif(NewYoutubeData, Bot)

					} else if !Items.Snippet.PublishedAt.IsZero() && Items.Snippet.VideoStatus == "none" {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.Name,
						}).Info("Suddenly upload new video")
						if NewYoutubeData.Schedul.IsZero() {
							NewYoutubeData.UpdateSchdule(NewYoutubeData.Published)
						}

						NewYoutubeData.UpdateStatus(config.PastStatus).InputYt()
						engine.SendLiveNotif(NewYoutubeData, Bot)

					} else {
						log.WithFields(log.Fields{
							"YtID":       ID,
							"MemberName": Member.Name,
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

					Data, err := engine.YtAPI([]string{ID})
					if err != nil {
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
						log.Error(err)
					}

					if len(Data.Items) == 0 {
						fmt.Println("Opps something error\n", Data)
						continue
					}

					Items := Data.Items[0]

					YoutubeData.UpdateEnd(Items.LiveDetails.EndTime).
						UpdateViewers(Items.Statistics.ViewCount).
						UpdateLength(durafmt.Parse(engine.ParseDuration(Items.ContentDetails.Duration)).String()).
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
						engine.SendLiveNotif(YoutubeData, Bot)

					} else if Items.Snippet.VideoStatus == config.UpcomingStatus && YoutubeData.Status == config.PastStatus {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Status":       Items.Snippet.VideoStatus,
						}).Info("maybe yt error or human error")

						YoutubeData.UpdateStatus(config.UpcomingStatus)
					} else if Items.Snippet.VideoStatus == "none" && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.PastStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateYt(config.PastStatus)

					} else if Items.Snippet.VideoStatus == config.LiveStatus && YoutubeData.Viewers != Items.Statistics.ViewCount {
						log.WithFields(log.Fields{
							"VideoData ID": ID,
							"Viwers past":  YoutubeData.Viewers,
							"Viwers now":   Items.Statistics.ViewCount,
							"Status":       config.LiveStatus,
						}).Info("Update Viwers")
						YoutubeData.UpdateYt(config.LiveStatus)

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
					}
				}
			}
		}
	}
}
