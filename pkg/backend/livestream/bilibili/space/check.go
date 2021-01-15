package space

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "github.com/JustHumanz/Go-simp/tools/database"
	network "github.com/JustHumanz/Go-simp/tools/network"
	"github.com/go-redis/redis/v8"

	log "github.com/sirupsen/logrus"
)

func (Space *CheckSctruct) Check(limit string) *CheckSctruct {
	var (
		Videotype string
		PushVideo SpaceVideo
		NewVideo  Vlist
		ctx       = context.Background()
		rdb       = database.LiveCache
	)
	body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/arc/search?mid="+strconv.Itoa(Space.SpaceID)+"&ps="+limit, nil)
	if curlerr != nil {
		log.Error(curlerr)
	}

	err := json.Unmarshal(body, &PushVideo)
	if err != nil {
		log.Error(err)
	}

	for _, video := range PushVideo.Data.List.Vlist {
		_, err := rdb.Get(ctx, video.Bvid).Result()
		if err == redis.Nil {
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
				MemberID: Space.Member.ID,
			}
			new, id := Data.CheckVideo()
			if new {
				Data.InputSpaceVideo()
				video.Pic = "https:" + video.Pic
				video.VideoType = Videotype
				NewVideo = append(NewVideo, video)
			} else {
				Data.UpdateView(id)
			}
			err := rdb.Set(ctx, video.Bvid, video, 6*time.Hour).Err()
			if err != nil {
				log.Error(err)
			}
		} else if err != nil {
			log.Error(err)
		}
	}
	Space.VideoList = NewVideo
	return Space
}
