package space

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	log "github.com/sirupsen/logrus"
)

func (Space *CheckSctruct) Check(limit string) *CheckSctruct {
	var (
		Videotype string
		PushVideo SpaceVideo
		NewVideo  Vlist
		body      []byte
		curlerr   error
		urls      = "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Space.SpaceID) + "&ps=" + limit
	)
	body, curlerr = engine.Curl(urls, nil)
	if curlerr != nil {
		log.Info("Trying use tor")
		body, curlerr = engine.CoolerCurl(urls, nil)
		if curlerr != nil {
			log.Error(curlerr)
		} else {
			log.Info("Successfully use tor")
		}
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
		Data := database.InputBiliBili{
			VideoID:  video.Bvid,
			Type:     Videotype,
			Title:    video.Title,
			Thum:     "https:" + video.Pic,
			Desc:     video.Description,
			Update:   time.Unix(int64(video.Created), 0).In(loc),
			Viewers:  video.Play,
			Length:   video.Length,
			MemberID: Space.MemberID,
		}
		if new, id := Data.CheckVideo(); new {
			Data.InputSpaceVideo()
			video.Pic = "https:" + video.Pic
			video.VideoType = Videotype
			NewVideo = append(NewVideo, video)
		} else {
			Data.UpdateView(id)
		}
	}
	Space.VideoList = NewVideo
	return Space
}
