package main

import (
	"context"
	"encoding/json"
	"strconv"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckYoutube() {
	var YTstate Subs
	Token := engine.GetYtToken()
	for _, Group := range Payload.VtuberData {
		for _, Member := range Group.Members {
			if !Member.IsYtNill() && Member.Active() {
				body, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+Member.YoutubeID+"&key="+*Token, nil)
				if err != nil {
					log.Error(err, string(body))
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
						}
						YTSubscriberCount, err := strconv.Atoi(Item.Statistics.SubscriberCount)
						if err != nil {
							log.Error(err)
						}
						SendNotif := func(SubsCount string) {
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
							SendNude(engine.NewEmbed().
								SetAuthor(Group.GroupName, Group.IconURL, "https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
								SetTitle(engine.FixName(Member.EnName, Member.JpName)).
								SetThumbnail(config.YoutubeIMG).
								SetDescription("Congratulation for "+SubsCount+" subscriber").
								SetImage(Member.YoutubeAvatar).
								AddField("Viewers", engine.NearestThousandFormat(float64(ViewCount))).
								AddField("Videos", engine.NearestThousandFormat(float64(VideoCount))).
								InlineAllFields().
								SetURL("https://www.youtube.com/channel/"+Member.YoutubeID+"?sub_confirmation=1").
								SetColor(Color).MessageEmbed, Group, Member)
						}
						if YtSubsDB.YtSubs != YTSubscriberCount {
							if YTSubscriberCount >= 1000000 {
								for i := 0; i < 10000001; i += 1000000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(engine.NearestThousandFormat(float64(i)))
									}
								}
							} else if YTSubscriberCount >= 100000 {
								for i := 0; i < 1000001; i += 100000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(engine.NearestThousandFormat(float64(i)))
									}
								}
							} else if YTSubscriberCount >= 10000 {
								for i := 0; i < 100001; i += 10000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(engine.NearestThousandFormat(float64(i)))
									}
								}
							} else if YTSubscriberCount >= 1000 {
								for i := 0; i < 10001; i += 1000 {
									if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
										SendNotif(engine.NearestThousandFormat(float64(i)))
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
