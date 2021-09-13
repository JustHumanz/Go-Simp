package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckYoutube() {
	var YTstate Subs
	Token := engine.GetYtToken()
	for _, Group := range *Payload {
		for _, Member := range Group.Members {
			if !Member.IsYtNill() && Member.Active() {
				body, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+Member.YoutubeID+"&key="+*Token, nil)
				if err != nil {
					if err.Error() == "403 Forbidden" {
						log.Error(err, string(body))
					} else {
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}
				}
				err = json.Unmarshal(body, &YTstate)
				if err != nil {
					log.Error(err)
				}
				for _, Item := range YTstate.Items {
					if Member.YoutubeID == Item.ID && !Item.Statistics.HiddenSubscriberCount {
						YtSubsDB, err := Member.GetSubsCount()
						if err != nil {
							log.Error(err)
							gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
								Message: err.Error(),
								Service: ModuleState,
							})
						}
						YTSubscriberCount, err := strconv.Atoi(Item.Statistics.SubscriberCount)
						if err != nil {
							log.Error(err)
						}
						SendNotif := func(SubsCount, NextCount string, nextdt, score int64) {
							err = Member.RemoveSubsCache()
							if err != nil {
								log.Error(err)
							}
							log.WithFields(log.Fields{
								"Vtuber":             Member.Name,
								"Group":              Group.GroupName,
								"TimeSampPrediction": nextdt,
								"Score":              score,
							}).Info("Congratulation for " + SubsCount + " subscriber")

							Color, err := engine.GetColor(config.TmpDir, Member.YoutubeAvatar)
							if err != nil {
								log.Error(err)
							}
							VideoCount, err := strconv.Atoi(Item.Statistics.VideoCount)
							if err != nil {
								log.Error(err)
							}
							ViewCount, err := strconv.Atoi(Item.Statistics.ViewCount)
							if err != nil {
								log.Error(err)
							}

							Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22Youtube%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"

							if nextdt != 0 && score != 0 {
								datePredic := time.Unix(nextdt, 0)
								datePredicStr := fmt.Sprintf("%s/%d", datePredic.Month().String(), datePredic.Day())
								SendNude(engine.NewEmbed().
									SetAuthor(Group.GroupName, Group.IconURL, "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
									SetTitle(engine.FixName(Member.EnName, Member.JpName)).
									SetThumbnail(config.YoutubeIMG).
									SetDescription("Congratulation for "+SubsCount+" subscriber").
									SetImage(Member.YoutubeAvatar).
									AddField("Viewers", engine.NearestThousandFormat(float64(ViewCount))).
									AddField("Videos", engine.NearestThousandFormat(float64(VideoCount))).
									InlineAllFields().
									AddField("Graph", Graph).
									SetURL("https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
									AddField("Next Milestone Prediction", datePredicStr+" - "+NextCount).
									SetFooter("Prediction Score "+strconv.Itoa(int(score))+"%").
									SetColor(Color).MessageEmbed, Group, Member)
							} else {
								SendNude(engine.NewEmbed().
									SetAuthor(Group.GroupName, Group.IconURL, "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
									SetTitle(engine.FixName(Member.EnName, Member.JpName)).
									SetThumbnail(config.YoutubeIMG).
									SetDescription("Congratulation for "+SubsCount+" subscriber").
									SetImage(Member.YoutubeAvatar).
									AddField("Viewers", engine.NearestThousandFormat(float64(ViewCount))).
									AddField("Videos", engine.NearestThousandFormat(float64(VideoCount))).
									InlineAllFields().
									AddField("Graph", Graph).
									SetURL("https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
									SetColor(Color).MessageEmbed, Group, Member)
							}

						}
						if YtSubsDB.YtSubs != YTSubscriberCount {
							if YTSubscriberCount >= 1000000 {
								for i := 0; i < 10000001; i += 1000000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										NextCount := YTSubscriberCount + 1000000
										dt, sc, err := SubsPreDick(NextCount, "Youtube", Member.Name)
										if err != nil {
											log.Error(err)
											gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
												Message: err.Error(),
												Service: ModuleState,
											})
										}
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
											engine.NearestThousandFormat(float64(NextCount)),
											dt, sc)
									}
								}
							} else if YTSubscriberCount >= 100000 {
								for i := 0; i < 1000001; i += 100000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										NextCount := YTSubscriberCount + 100000
										dt, sc, err := SubsPreDick(NextCount, "Youtube", Member.Name)
										if err != nil {
											log.Error(err)
											gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
												Message: err.Error(),
												Service: ModuleState,
											})
										}
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
											engine.NearestThousandFormat(float64(NextCount)),
											dt, sc)
									}
								}
							} else if YTSubscriberCount >= 10000 {
								for i := 0; i < 100001; i += 10000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										NextCount := YTSubscriberCount + 10000
										dt, sc, err := SubsPreDick(NextCount, "Youtube", Member.Name)
										if err != nil {
											log.Error(err)
											gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
												Message: err.Error(),
												Service: ModuleState,
											})
										}
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
											engine.NearestThousandFormat(float64(NextCount)),
											dt, sc)
									}
								}
							} else if YTSubscriberCount >= 1000 {
								for i := 0; i < 10001; i += 1000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										NextCount := YTSubscriberCount + 1000
										dt, sc, err := SubsPreDick(NextCount, "Youtube", Member.Name)
										if err != nil {
											log.Error(err)
											gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
												Message: err.Error(),
												Service: ModuleState,
											})
										}
										SendNotif(
											engine.NearestThousandFormat(float64(i)),
											engine.NearestThousandFormat(float64(NextCount)),
											dt, sc)
									}
								}
							}
						}

						log.WithFields(log.Fields{
							"Past Youtube subscriber":    YtSubsDB.YtSubs,
							"Current Youtube subscriber": YTSubscriberCount,
							"Vtuber":                     Member.EnName,
						}).Info("Update Youtube subscriber")
						VideoCount, err := strconv.Atoi(Item.Statistics.VideoCount)
						if err != nil {
							log.Error(err)
						}
						ViewCount, err := strconv.Atoi(Item.Statistics.ViewCount)
						if err != nil {
							log.Error(err)
						}

						YtSubsDB.SetMember(Member).SetGroup(Group).
							UpYtSubs(YTSubscriberCount).
							UpYtVideo(VideoCount).
							UpYtViews(ViewCount).
							UpdateState(config.YoutubeLive).
							UpdateSubs()

						bin, err := YtSubsDB.MarshalBinary()
						if err != nil {
							log.Error(err)
						}
						if config.GoSimpConf.Metric {
							gRCPconn.MetricReport(context.Background(), &pilot.Metric{
								MetricData: bin,
								State:      config.SubsState,
							})
						}
					}
				}
			}
		}
	}
}
