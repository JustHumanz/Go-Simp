package main

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"

	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

type Video struct {
	ID      string
	Preview string
}

func Tweet(Group string, NameID int64, Limit int) {
	var (
		a []string
	)
	if NameID > 0 {
		tmp := GetHastagMember(NameID)
		a = append(a, tmp)
	} else {
		for _, hashtag := range GetHashtag(Group) {
			if hashtag.TwitterHashtags != "" && hashtag.JpName == "桐生ココ" {
				a = append(a, hashtag.TwitterHashtags)
			}
		}
	}
	Hashtag := strings.Join(a, " OR ")
	log.Info(Hashtag)

	for tweet := range twitterscraper.SearchTweets(context.Background(), Hashtag+"  filter:links -filter:replies filter:media", Limit) {
		engine.BruhMoment(tweet.Error, "", false)
		Data := &InputTwitter{
			TwitterData: tweet.Tweet,
			Group:       Group,
			MemberID:    NameID,
		}
		Data.FiltterTweet().InputData()
	}
}

func (Data *InputTwitter) FiltterTweet() *InputTwitter {
	for _, Hashtag := range GetHashtag(Data.Group) {
		matched, _ := regexp.MatchString(Hashtag.TwitterHashtags, strings.Join(Data.TwitterData.Hashtags, " "))
		if matched {
			Data.MemberID = Hashtag.MemberID
		}
	}
	return Data
}

type InputTwitter struct {
	TwitterData twitterscraper.Tweet
	Group       string
	MemberID    int64
}

func FilterYt(Dat database.Name, wg *sync.WaitGroup) {
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
