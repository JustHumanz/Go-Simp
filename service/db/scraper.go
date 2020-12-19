package main

import (
	"encoding/json"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/pkg/backend/fanart/twitter"
	bilibili "github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/live"
	youtube "github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	network "github.com/JustHumanz/Go-simp/tools/network"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func TwitterFanart() {
	scraper := twitterscraper.New()
	scraper.SetProxy(config.MultiTOR)
	for _, Group := range database.GetGroups() {
		var wg sync.WaitGroup
		for _, Member := range database.GetMembers(Group.ID) {
			wg.Add(1)
			go func(wg *sync.WaitGroup, Member database.Member, Group database.Group) {
				defer wg.Done()
				if Member.TwitterHashtags != "" || Member.Name != "Kaichou" {
					Newfanart := twitter.TwitterFanart{
						Member:  Member,
						Limit:   100,
						Group:   Group,
						Scraper: scraper,
					}
					Newfanart.CurlTwitter()
				} else {
					log.Info(Member.EnName + " don't have twitter hashtag")
				}
			}(&wg, Member, Group)
		}
		wg.Wait()
	}
}

func FilterYt(Dat database.Member, wg *sync.WaitGroup) {
	VideoID := youtube.GetRSS(Dat.YoutubeID)
	defer wg.Done()
	body, err := network.Curl("https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,actualEndTime),statistics(viewCount))&id="+strings.Join(VideoID, ",")+"&key="+YtToken, nil)
	if err != nil {
		log.Error(err, string(body))
	}
	var (
		Data    YtData
		Viewers string
		yttype  string
	)
	err = json.Unmarshal(body, &Data)
	if err != nil {
		log.Error(err)
	}
	for i := 0; i < len(Data.Items); i++ {
		if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|mv)", Data.Items[i].Snippet.Title); Cover {
			yttype = "Covering"
		} else if Chat, _ := regexp.MatchString("(?m)(free|chat|room)", Data.Items[i].Snippet.Title); Chat {
			yttype = "ChatRoom"
		} else {
			yttype = "Streaming"
		}
		if Dat.CheckYtVideo(VideoID[i]) != nil {
			continue
		} else {
			log.Info("New video")
			//verify
			if Data.Items[i].LiveDetails.Viewers == "" {
				Viewers = Data.Items[i].Statistics.ViewCount
			} else {
				Viewers = Data.Items[i].LiveDetails.Viewers
			}
			NewData := database.YtDbData{
				Status:    Data.Items[i].Snippet.VideoStatus,
				VideoID:   VideoID[i],
				Title:     Data.Items[i].Snippet.Title,
				Thumb:     "http://i3.ytimg.com/vi/" + VideoID[i] + "/maxresdefault.jpg",
				Desc:      Data.Items[i].Snippet.Description,
				Schedul:   Data.Items[i].LiveDetails.StartTime,
				Published: Data.Items[i].Snippet.PublishedAt,
				End:       Data.Items[i].LiveDetails.EndTime,
				Type:      yttype,
				Viewers:   Viewers,
			}

			if Data.Items[i].Snippet.VideoStatus != "upcoming" || Data.Items[i].Snippet.VideoStatus != "live" {
				NewData.Status = "past"
				NewData.InputYt(Dat.ID)
			} else {
				NewData.InputYt(Dat.ID)
			}
		}
	}
}

func (Data Member) YtAvatar() string {
	if Data.YtID != "" {
		var (
			avatar string
			bit    []byte
			err    error
			URL    = "https://www.youtube.com/channel/" + Data.YtID + "/about"
		)
		bit, err = network.Curl(URL, nil)
		if err != nil {
			bit, err = network.CoolerCurl(URL, nil)
			if err != nil {
				log.Error(err)
			}
		}
		submatchall := regexp.MustCompile(`(?ms)avatar.*?(http.*?)"`).FindAllStringSubmatch(string(bit), -1)
		for _, element := range submatchall {
			avatar = strings.Replace(element[1], "s48", "s800", -1)
			break
		}
		return avatar
	} else {
		return Data.BliBiliFace()
	}
}

func (Data Member) GetYtSubs() Subs {
	var (
		datasubs Subs
	)
	if Data.YtID != "" {
		body, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+Data.YtID+"&key="+YtToken, nil)
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(body, &datasubs)
		if err != nil {
			log.Error(err)
		}
		return datasubs
	} else {
		return datasubs.Default()
	}
}

func (Data Member) GetBiliFolow() BiliStat {
	var (
		wg      sync.WaitGroup
		stat    BiliStat
		body    []byte
		curlerr error
	)
	if Data.BiliRoomID != 0 {
		wg.Add(3)
		go func() {
			urls := "https://api.bilibili.com/x/relation/stat?vmid=" + strconv.Itoa(Data.BiliBiliID)
			body, curlerr = network.Curl(urls, nil)
			if curlerr != nil {
				log.Warn("Trying use tor")
				body, curlerr = network.CoolerCurl(urls, nil)
				if curlerr != nil {
					log.Error(curlerr)
				}
			}
			err := json.Unmarshal(body, &stat.Follow)
			if err != nil {
				log.Error(err)
			}
			defer wg.Done()
		}()

		go func() {
			urls := "https://api.bilibili.com/x/space/upstat?mid=" + strconv.Itoa(Data.BiliBiliID)
			body, curlerr = network.Curl(urls, []string{"Cookie", "SESSDATA=" + BiliSession})
			if curlerr != nil {
				log.Warn("Trying use tor")
				body, curlerr = network.CoolerCurl(urls, []string{"Cookie", "SESSDATA=" + BiliSession})
				if curlerr != nil {
					log.Error(curlerr)
				}
			}
			err := json.Unmarshal(body, &stat.Like)
			if err != nil {
				log.Error(err)
			}
			defer wg.Done()
		}()

		go func() {
			baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Data.BiliBiliID) + "&ps=100"
			url := []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
			for f := 0; f < len(url); f++ {
				body, curlerr = network.Curl(url[f], nil)
				if curlerr != nil {
					log.Warn("Trying use tor")
					body, curlerr = network.CoolerCurl(url[f], nil)
					if curlerr != nil {
						log.Error(curlerr)
					}
				}
				var video SpaceVideo
				err := json.Unmarshal(body, &video)
				if err != nil {
					log.Error(err)
				}
				stat.Video += video.Data.Page.Count
			}
			defer wg.Done()
		}()
		wg.Wait()
		return stat
	} else {
		log.WithFields(log.Fields{
			"Vtuber": Data.ENName,
		}).Info("BiliBili Space nill")
		return stat
	}
}

func (Data Member) GetTwitterFollow() int {
	if Data.TwitterName != "" {
		profile, err := twitterscraper.GetProfile(Data.TwitterName)
		if err != nil {
			log.Error(err)
		}
		return profile.FollowersCount
	} else {
		return 0
	}
}

func (Data Member) BliBiliFace() string {
	if Data.BiliBiliID == 0 {
		return ""
	} else {
		var (
			Info    Avatar
			body    []byte
			errcurl error
			url     = "https://api.bilibili.com/x/space/acc/info?mid=" + strconv.Itoa(Data.BiliBiliID)
		)
		body, errcurl = network.Curl(url, nil)
		if body == nil {
			log.Info("Not daijobu,trying use multitor")
			body, errcurl = network.CoolerCurl(url, nil)

			if errcurl != nil {
				log.Error(errcurl)
				return ""
			}
		} else if errcurl != nil {
			log.Error(errcurl)
			return ""
		}
		err := json.Unmarshal(body, &Info)
		if err != nil {
			log.Error(err)
			return ""
		}

		return strings.Replace(Info.Data.Face, "http", "https", -1)
	}
}

func CheckYT() {
	Data := database.GetGroups()
	for i := 0; i < len(Data); i++ {
		var wg sync.WaitGroup
		for _, Name := range database.GetMembers(Data[i].ID) {
			wg.Add(1)
			go func(Name database.Member) {
				if Name.YoutubeID != "" {
					log.WithFields(log.Fields{
						"Vtube":        Name.EnName,
						"Youtube ID":   Name.YoutubeID,
						"Vtube Region": Name.Region,
					}).Info("Checking yt")
					FilterYt(Name, &wg)
				}

			}(Name)
		}
		wg.Wait()
	}
}

func CheckTBili() {
	DataGroup := database.GetGroups()
	for k := 0; k < len(DataGroup); k++ {
		DataMember := database.GetMembers(DataGroup[k].ID)
		for z := 0; z < len(DataMember); z++ {
			if DataMember[z].BiliBiliHashtags != "" {
				log.WithFields(log.Fields{
					"Group":  DataGroup[k].NameGroup,
					"Vtuber": DataMember[z].EnName,
				}).Info("Start crawler T.bilibili")
				body, err := network.Curl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(DataMember[z].BiliBiliHashtags), nil)
				if err != nil {
					log.Error(err)
				}
				var (
					TB              TBiliBili
					DynamicIDStrTmp string
				)
				_ = json.Unmarshal(body, &TB)
				if (len(TB.Data.Cards) > 0) && TB.Data.Cards[0].Desc.DynamicIDStr != DynamicIDStrTmp {
					DynamicIDStrTmp = TB.Data.Cards[0].Desc.DynamicIDStr
					for i := 0; i < len(TB.Data.Cards); i++ {
						var (
							STB  SubTbili
							img  []string
							nope bool
						)
						_ = json.Unmarshal([]byte(TB.Data.Cards[i].Card), &STB)
						if STB.Item.Pictures != nil && TB.Data.Cards[i].Desc.Type == 2 { //type 2 is picture post (prob,heheheh)
							niggerlist := []string{"解锁专属粉丝卡片", "Official", "twitter.com", "咖啡厅", "CD", "专辑", "PIXIV", "遇", "marshmallow-qa.com"}
							for _, Nworld := range niggerlist {
								nope, _ = regexp.MatchString(Nworld, STB.Item.Description)
								if nope {
									break
								}
							}
							New := database.GetTBiliBili(TB.Data.Cards[i].Desc.DynamicIDStr)

							if New && !nope {
								log.WithFields(log.Fields{
									"Group":  DataGroup[k].NameGroup,
									"Vtuber": DataMember[z].EnName,
								}).Info("New Fanart")
								for l := 0; l < len(STB.Item.Pictures); l++ {
									img = append(img, STB.Item.Pictures[l].ImgSrc)
								}

								Data := database.InputTBiliBili{
									URL:        "https://t.bilibili.com/" + TB.Data.Cards[i].Desc.DynamicIDStr + "?tab=2",
									Author:     TB.Data.Cards[i].Desc.UserProfile.Info.Uname,
									Avatar:     TB.Data.Cards[i].Desc.UserProfile.Info.Face,
									Like:       TB.Data.Cards[i].Desc.Like,
									Photos:     strings.Join(img, "\n"),
									Dynamic_id: TB.Data.Cards[i].Desc.DynamicIDStr,
									Text:       STB.Item.Description,
								}
								log.Info("Send to database")
								Data.InputTBiliBili(DataMember[z].ID)
							} else {
								log.WithFields(log.Fields{
									"Group":  DataGroup[k].NameGroup,
									"Vtuber": DataMember[z].EnName,
								}).Info("Still same")
							}
						}
					}
				} else {
					log.WithFields(log.Fields{
						"Group":  DataGroup[k].NameGroup,
						"Vtuber": DataMember[z].EnName,
					}).Info("Still same")
				}
				time.Sleep(time.Duration(int64(rand.Intn((7-1)+1))) * time.Second)
			}
		}
	}
}

func CheckSchedule() {
	log.Info("Start check BiliBili room")
	for _, Group := range database.GetGroups() {
		for _, Member := range database.GetMembers(Group.ID) {
			if Member.BiliBiliID != 0 {
				log.WithFields(log.Fields{
					"Group":   Group.NameGroup,
					"SpaceID": Member.EnName,
				}).Info("Check Room")
				var (
					ScheduledStart time.Time
				)
				DataDB := database.GetRoomData(Member.ID, Member.BiliRoomID)
				Status, err := bilibili.GetRoomStatus(Member.BiliRoomID)
				if err != nil {
					log.Error(err)
				}
				loc, _ := time.LoadLocation("Asia/Shanghai")
				if Status.Data.RoomInfo.LiveStartTime != 0 {
					ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0).In(loc)
				} else {
					ScheduledStart = time.Time{}
				}
				Data := map[string]interface{}{
					"LiveRoomID":     Member.BiliRoomID,
					"Status":         "",
					"Title":          Status.Data.RoomInfo.Title,
					"Thumbnail":      Status.Data.RoomInfo.Cover,
					"Description":    Status.Data.NewsInfo.Content,
					"PublishedAt":    time.Time{},
					"ScheduledStart": ScheduledStart,
					"Face":           Status.Data.AnchorInfo.BaseInfo.Face,
					"Online":         Status.Data.RoomInfo.Online,
					"BiliBiliID":     Member.BiliBiliID,
					"MemberID":       Member.ID,
				}
				if Status.CheckScheduleLive() {
					//Live
					log.WithFields(log.Fields{
						"Group":      Group.NameGroup,
						"VtuberName": Member.Name,
					}).Info("Status Live")
					Data["Status"] = "Live"
					LiveBiliBili(Data)
				} else if !Status.CheckScheduleLive() && DataDB.Status == "Live" {
					//prob past
					log.WithFields(log.Fields{
						"Group":      Group.NameGroup,
						"VtuberName": Member.Name,
					}).Info("Status Past")
					Data["Status"] = "Past"
					LiveBiliBili(Data)
				} else if DataDB.LiveRoomID == 0 {
					log.WithFields(log.Fields{
						"Group":      Group.NameGroup,
						"VtuberName": Member.Name,
					}).Info("Status Unknown")
					Data["Status"] = "Unknown"
					LiveBiliBili(Data)
				}
			}
		}
	}
}

func CheckVideoSpace() {
	Group := database.GetGroups()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for z := 0; z < len(Group); z++ {
		Name := database.GetMembers(Group[z].ID)
		for k := 0; k < len(Name); k++ {
			if Name[k].BiliBiliID != 0 {
				log.WithFields(log.Fields{
					"Group":   Group[z].NameGroup,
					"SpaceID": Name[k].EnName,
				}).Info("Check Space")
				var (
					PushVideo SpaceVideo
					videotype string
					url       []string
				)
				baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Name[k].BiliBiliID) + "&ps=100"
				url = []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
				for f := 0; f < len(url); f++ {
					body, err := network.Curl(url[f], nil)
					if err != nil {
						log.Error(err, string(body))
					}
					var tmp SpaceVideo
					err = json.Unmarshal(body, &tmp)
					if err != nil {
						log.Error(err)
					}
					for _, Vlist := range tmp.Data.List.Vlist {
						PushVideo.Data.List.Vlist = append(PushVideo.Data.List.Vlist, Vlist)
					}
				}

				for _, video := range PushVideo.Data.List.Vlist {
					if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|翻唱|mv)", strings.ToLower(video.Title)); Cover {
						videotype = "Covering"
					} else {
						videotype = "Streaming"
					}
					tmp := database.InputBiliBili{
						VideoID:  video.Bvid,
						Type:     videotype,
						Title:    video.Title,
						Thum:     "https:" + video.Pic,
						Desc:     video.Description,
						Update:   time.Unix(int64(video.Created), 0).In(loc),
						Viewers:  video.Play,
						MemberID: Name[k].ID,
					}
					tmp.InputSpaceVideo()
				}
			}
		}
	}
}

type BiliStat struct {
	Follow BiliFollow
	Like   LikeView
	Video  int
}

type LikeView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Archive struct {
			View int `json:"view"`
		} `json:"archive"`
		Article struct {
			View int `json:"view"`
		} `json:"article"`
		Likes int `json:"likes"`
	} `json:"data"`
}

type BiliFollow struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int `json:"mid"`
		Following int `json:"following"`
		Whisper   int `json:"whisper"`
		Black     int `json:"black"`
		Follower  int `json:"follower"`
	} `json:"data"`
}
