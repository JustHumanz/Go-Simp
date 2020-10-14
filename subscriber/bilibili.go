package subscriber

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	"github.com/JustHumanz/Go-simp/livestream/bilibili/space"
	log "github.com/sirupsen/logrus"
)

func CheckBiliFollowCount() {
	for _, Group := range Data {
		Names := database.GetName(Group.ID)
		for _, Name := range Names {
			if Name.BiliBiliID != 0 {
				var (
					wg        sync.WaitGroup
					bilistate BiliBiliStat
				)
				wg.Add(3)
				go func() {
					var (
						body    []byte
						curlerr error
						urls    = "https://api.bilibili.com/x/relation/stat?vmid=" + strconv.Itoa(Name.BiliBiliID)
					)
					body, curlerr = engine.Curl(urls, nil)
					if curlerr != nil {
						log.Error(curlerr, string(body))

						log.Info("Trying use tor")
						body, curlerr = engine.CoolerCurl(urls, nil)
						if curlerr != nil {
							log.Error(curlerr)
						}
					}
					err := json.Unmarshal(body, &bilistate.Follow)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					body, err := engine.Curl("https://api.bilibili.com/x/space/upstat?mid="+strconv.Itoa(Name.BiliBiliID), []string{"Cookie", "SESSDATA=" + BiliSession})
					if err != nil {
						log.Error(err, string(body))
					}
					err = json.Unmarshal(body, &bilistate.LikeView)
					if err != nil {
						log.Error(err)
					}
					defer wg.Done()
				}()

				go func() {
					baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Name.BiliBiliID) + "&ps=100"
					url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
					for f := 0; f < len(url); f++ {
						body, err := engine.Curl(url[f], nil)
						if err != nil {
							log.Error(err, string(body))
						}
						var video space.SpaceVideo
						err = json.Unmarshal(body, &video)
						if err != nil {
							log.Error(err)
						}
						bilistate.Videos += video.Data.Page.Count
					}
					defer wg.Done()
				}()
				wg.Wait()

				if Name.BiliBiliID != 0 {
					if bilistate.Follow.Data.Follower <= 10000 {
						Avatar := Name.BiliBiliAvatar
						Color, err := engine.GetColor("/tmp/bili.tmp", Avatar)
						if err != nil {
							log.Error(err)
						}
						for i := 0; i < 1000001; i += 100000 {
							if i == bilistate.Follow.Data.Follower && bilistate.Follow.Data.Follower != 0 {
								SendNude(engine.NewEmbed().
									SetAuthor(Group.NameGroup, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
									SetTitle(engine.FixName(Name.EnName, Name.JpName)).
									SetThumbnail(config.BiliBiliIMG).
									SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
									SetImage(Avatar).
									AddField("Views", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
									AddField("Videos", strconv.Itoa(bilistate.Videos)).
									SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
									InlineAllFields().
									SetColor(Color).MessageEmbed, Group)
							}
						}
					} else {
						Avatar := Name.BiliBiliAvatar
						Color, err := engine.GetColor("/tmp/bili.tmp", Avatar)
						if err != nil {
							log.Error(err)
						}
						for i := 0; i < 10001; i += 1000 {
							if i == bilistate.Follow.Data.Follower && bilistate.Follow.Data.Follower != 0 {
								SendNude(engine.NewEmbed().
									SetAuthor(Group.NameGroup, Group.IconURL, "https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
									SetTitle(engine.FixName(Name.EnName, Name.JpName)).
									SetThumbnail(config.BiliBiliIMG).
									SetDescription("Congratulation for "+strconv.Itoa(i)+" followers").
									SetImage(Avatar).
									AddField("Views", strconv.Itoa(bilistate.LikeView.Data.Archive.View)).
									AddField("Videos", strconv.Itoa(bilistate.Videos)).
									SetURL("https://space.bilibili.com/"+strconv.Itoa(Name.BiliBiliID)).
									InlineAllFields().
									SetColor(Color).MessageEmbed, Group)
							}
						}
					}
				}
				BiliFollowDB := Name.GetSubsCount()
				log.WithFields(log.Fields{
					"Past BiliBili Follower":    BiliFollowDB.BiliFollow,
					"Current BiliBili Follower": bilistate.Follow.Data.Follower,
					"Vtuber":                    Name.EnName,
				}).Info("Update BiliBili Follower")
				BiliFollowDB = database.MemberSubs{
					BiliFollow: bilistate.Follow.Data.Follower,
					BiliVideos: bilistate.Videos,
					BiliViews:  bilistate.LikeView.Data.Archive.View,
				}
				BiliFollowDB.UpdateSubs("bili")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}
