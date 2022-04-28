package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckBiliBili() {
	BiliBiliSession := map[string]string{
		"Cookie": "SESSDATA=" + configfile.BiliSess,
	}
	for _, Group := range Agency {
		for _, Member := range Group.Members {
			if Member.BiliBiliID != 0 && Member.Active() {
				var (
					bilistate BiliBiliStat
				)

				body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/relation/stat?vmid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
				if curlerr != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					})
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message:     curlerr.Error(),
						Service:     ServiceName,
						ServiceUUID: ServiceUUID,
					})
				}
				err := json.Unmarshal(body, &bilistate.Follow)
				if err != nil {
					log.Error(err)
				}

				body, curlerr = network.CoolerCurl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Member.BiliBiliID), BiliBiliSession)
				if curlerr != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					})
					gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
						Message:     curlerr.Error(),
						Service:     ServiceName,
						ServiceUUID: ServiceUUID,
					})
				}
				err = json.Unmarshal(body, &bilistate.LikeView)
				if err != nil {
					log.Error(err)
				}

				baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Member.BiliBiliID) + "&ps=100"
				url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
				for f := 0; f < len(url); f++ {
					body, curlerr := network.CoolerCurl(url[f], BiliBiliSession)
					if curlerr != nil {
						log.WithFields(log.Fields{
							"Agency": Group.GroupName,
							"Vtuber": Member.Name,
						})
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message:     curlerr.Error(),
							Service:     ServiceName,
							ServiceUUID: ServiceUUID,
						})
					}
					var video engine.SpaceVideo
					err := json.Unmarshal(body, &video)
					if err != nil {
						log.Error(err)
					}
					bilistate.Videos += video.Data.Page.Count
				}

				BiliFollowDB, err := Member.GetSubsCount()
				if err != nil {
					log.WithFields(log.Fields{
						"Agency": Group.GroupName,
						"Vtuber": Member.Name,
					})
				}

				if bilistate.Follow.Data.Follower != 0 {
					log.WithFields(log.Fields{
						"Past BiliBili Follower":    BiliFollowDB.BiliFollow,
						"Current BiliBili Follower": bilistate.Follow.Data.Follower,
						"Vtuber":                    Member.Name,
					}).Info("Update BiliBili Follower")

					err := BiliFollowDB.SetMember(Member).SetGroup(Group).
						UpdateBiliBiliFollowers(bilistate.Follow.Data.Follower).
						UpdateBiliBiliVideos(bilistate.Videos).
						UpdateBiliBiliViewers(bilistate.LikeView.Data.Archive.View).
						UpdateState(config.BiliLive).UpdateSubs()
					if err != nil {
						log.Error(err)
					}

					if BiliFollowDB.BiliFollow != bilistate.Follow.Data.Follower {
						if bilistate.Follow.Data.Follower <= 10000 {
							for i := 0; i < 1000001; i += 100000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Member.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.WithFields(log.Fields{
											"Agency": Group.GroupName,
											"Vtuber": Member.Name,
										})
									}

									err = Member.RemoveSubsCache()
									if err != nil {
										log.WithFields(log.Fields{
											"Agency": Group.GroupName,
											"Vtuber": Member.Name,
										})
									}
									Graph := "[View as Graph](" + os.Getenv("PrometheusURL") + "/graph?g0.expr=get_subscriber%7Bstate%3D%22BiliBili%22%2C%20vtuber%3D%22" + Member.Name + "%22%7D&g0.tab=0&g0.stacked=0&g0.range_input=4w)"
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										SetTitle(engine.FixName(Member.EnName, Member.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Viewers", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
										AddField("Videos", strconv.Itoa(bilistate.Videos)).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										InlineAllFields().
										AddField("Graph", Graph).
										SetColor(Color).MessageEmbed, Group, Member)
								}
							}
						} else {
							for i := 0; i < 10001; i += 1000 {
								if i == bilistate.Follow.Data.Follower {
									Avatar := Member.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										SetTitle(engine.FixName(Member.EnName, Member.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+engine.NearestThousandFormat(float64(i))+" followers").
										SetImage(Avatar).
										AddField("Views", engine.NearestThousandFormat(float64(bilistate.LikeView.Data.Archive.View))).
										AddField("Videos", engine.NearestThousandFormat(float64(bilistate.Videos))).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Member.BiliBiliID)).
										InlineAllFields().
										SetColor(Color).MessageEmbed, Group, Member)
								}
							}
						}
					}

					bin, err := BiliFollowDB.MarshalBinary()
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
