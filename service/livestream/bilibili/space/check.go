package main

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"

	log "github.com/sirupsen/logrus"
)

func CheckSpace(Data *database.LiveStream, limit string) {
	var (
		Videotype string
		PushVideo engine.SpaceVideo
	)

	body, curlerr := network.CoolerCurl("https://api.bilibili.com/x/space/arc/search?mid="+strconv.Itoa(Data.Member.BiliBiliID)+"&ps="+limit, nil)
	if curlerr != nil {
		log.Error(curlerr)
		gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
			Message: curlerr.Error(),
			Service: ModuleState,
		})
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
			UpdateThumbnail(video.Pic).UpdateSchdule(time.Unix(int64(video.Created), 0).In(loc)).
			UpdateViewers(strconv.Itoa(video.Play)).UpdateLength(video.Length).SetState(config.SpaceBili)
		new, id := Data.CheckVideo()
		if new {
			log.WithFields(log.Fields{
				"Vtuber": Data.Member.Name,
			}).Info("New video uploaded")

			Data.InputSpaceVideo()
			video.VideoType = Videotype
			engine.SendLiveNotif(Data, Bot)
		} else {
			if !config.GoSimpConf.LowResources {
				Data.UpdateSpaceViews(id)
			}
		}
	}
}
