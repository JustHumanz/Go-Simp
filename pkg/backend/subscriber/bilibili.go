package main

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/space"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"
	log "github.com/sirupsen/logrus"
)

func CheckBiliBili() {
	BiliBiliSession := []string{"Cookie", "SESSDATA=" + config.BotConf.BiliSess}
	for _, Group := range engine.GroupData {
		Names := database.GetMembers(Group.ID)
		for i, Name := range Names {
			if Name.BiliBiliID != 0 {
				var (
					wg        sync.WaitGroup
					bilistate BiliBiliStat
				)
				wg.Add(3)
				go func() {
					body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/relation/stat?vmid="+strconv.Itoa(Name.BiliBiliID), BiliBiliSession)
					if curlerr != nil {
						log.Error(curlerr)
					}
					err := json.Unmarshal(body, &bilistate.Follow)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Name.BiliBiliID), BiliBiliSession)
					if curlerr != nil {
						log.Error(curlerr)
					}
					err := json.Unmarshal(body, &bilistate.LikeView)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Name.BiliBiliID) + "&ps=100"
					url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
					for f := 0; f < len(url); f++ {
						body, curlerr := network.CoolerCurl(url[f], BiliBiliSession)
						if curlerr != nil {
							log.Error(curlerr)
						}
						var video space.SpaceVideo
						err := json.Unmarshal(body, &video)
						if err != nil {
							log.Error(err)
						}
						bilistate.Videos += video.Data.Page.Count
					}
					defer wg.Done()
				}()
				wg.Wait()

				BiliFollowDB, err := Name.GetSubsCount()
				if err != nil {
					log.Error(err)
				}
				if Name.BiliBiliID != 0 {
					if BiliFollowDB.BiliFollow != bilistate.Follow.Data.Follower {
						if bilistate.Follow.Data.Follower <= 10000 {
							for i := 0; i < 1000001; i += 100000 {
								if i == bilistate.Follow.Data.Follower && bilistate.Follow.Data.Follower != 0 {
									Avatar := Name.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										SetTitle(engine.FixName(Name.EnName, Name.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
										SetImage(Avatar).
										AddField("Viewers", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
										AddField("Videos", strconv.Itoa(bilistate.Videos)).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										InlineAllFields().
										SetColor(Color).MessageEmbed, Group, Name.ID)
								}
							}
						} else {
							for i := 0; i < 10001; i += 1000 {
								if i == bilistate.Follow.Data.Follower && bilistate.Follow.Data.Follower != 0 {
									Avatar := Name.BiliBiliAvatar
									Color, err := engine.GetColor(config.TmpDir, Avatar)
									if err != nil {
										log.Error(err)
									}
									SendNude(engine.NewEmbed().
										SetAuthor(Group.GroupName, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										SetTitle(engine.FixName(Name.EnName, Name.JpName)).
										SetThumbnail(config.BiliBiliIMG).
										SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
										SetImage(Avatar).
										AddField("Views", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
										AddField("Videos", strconv.Itoa(bilistate.Videos)).
										SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
										InlineAllFields().
										SetColor(Color).MessageEmbed, Group, Name.ID)
								}
							}
						}
					}
				}
				log.WithFields(log.Fields{
					"Past BiliBili Follower":    BiliFollowDB.BiliFollow,
					"Current BiliBili Follower": bilistate.Follow.Data.Follower,
					"Vtuber":                    Name.EnName,
				}).Info("Update BiliBili Follower")
				BiliFollowDB.UpBiliFollow(bilistate.Follow.Data.Follower).
					UpBiliVideo(bilistate.Videos).
					UpBiliViews(bilistate.LikeView.Data.Archive.View).
					UpdateSubs("bili")
			}
			if i%10 == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}
}
