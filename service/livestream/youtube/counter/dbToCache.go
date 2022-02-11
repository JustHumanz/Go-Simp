package main

import (
	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

//CacheChcker send upcoming livestream from db to cache
func CacheChcker() {
	log.Info("Start check upcoming from db")
	liveCache, err := database.GetUpcomingFromCache()
	if err != nil {
		log.Error(err)
	}
	for _, v := range *GroupPayload {
		liveStream, err := v.GetYtLiveStream(config.UpcomingStatus, nil)
		if err != nil {
			log.Error(err)
		}
		for _, k := range liveStream {
			for _, j := range liveCache {
				j := j.(database.LiveStream)
				if k.VideoID == j.VideoID {
					continue
				} else {
					log.WithFields(log.Fields{
						"videoID": j.VideoID,
						"agency":  j.Group.GroupName,
						"vtuber":  j.Member.Name,
					}).Info("Found upcoming on db but not in cache,send it to cache now")
					err := k.SendToCache(false)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}
