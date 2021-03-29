package space

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/livestream/notif"

	log "github.com/sirupsen/logrus"
)

func CheckSpace(Data *database.LiveStream, limit string) {
	var (
		Videotype string
		PushVideo SpaceVideo
	)

	body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/arc/search?mid="+strconv.Itoa(Data.Member.BiliBiliID)+"&ps="+limit, nil)
	if curlerr != nil {
		log.Error(curlerr)
	}

	err := json.Unmarshal(body, &PushVideo)
	if err != nil {
		log.Error(err)
	}

	for _, video := range PushVideo.Data.List.Vlist {
		if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|翻唱|mv|歌曲)", strings.ToLower(video.Title)); Cover || video.Typeid == 31 {
			Videotype = "Covering"
		} else {
			Videotype = "Streaming"
		}

		Data.AddVideoID(video.Bvid).SetType(Videotype).
			UpdateTitle(video.Title).
			UpdateThumbnail("https:" + video.Pic).UpdateSchdule(time.Unix(int64(video.Created), 0).In(loc)).
			UpdateViewers(strconv.Itoa(video.Play)).UpdateLength(video.Length)
		/*
			Data = database.LiveStream{
				VideoID: video.Bvid,
				Type:    Videotype,
				Title:   video.Title,
				Thumb:   "https:" + video.Pic,
				Desc:    video.Description,
				Schedul: time.Unix(int64(video.Created), 0).In(loc),
				Viewers: strconv.Itoa(video.Play),
				Length:  video.Length,
				Member:  Data.Member,
				Group:   Data.Group,
			}
		*/
		new, id := Data.CheckVideo()
		if new {
			log.WithFields(log.Fields{
				"Vtuber": Data.Member.Name,
			}).Info("New video uploaded")

			Data.InputSpaceVideo()
			video.Pic = "https:" + video.Pic
			video.VideoType = Videotype
			notif.SendDude(Data, Bot)
		} else {
			if !config.GoSimpConf.LowResources {
				Data.UpdateSpaceViews(id)
			}
		}
	}
}
