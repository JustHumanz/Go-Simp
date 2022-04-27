package main

import (
	"fmt"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

//CacheChecker send upcoming livestream from db to cache
func CacheChecker() {
	log.Info("Start check upcoming from db")
	liveCache, err := database.GetUpcomingFromCache()
	if err != nil {
		log.Error(err)
	}

	for _, Agency := range agency {
		for _, Member := range Agency.Members {
			liveStreamDb, err := Member.GetYtLiveStream(config.UpcomingStatus)
			if err != nil {
				log.Error(err)
			}
			for _, live := range liveStreamDb {
				key := fmt.Sprintf("%s-%s", Member.Name, live.VideoID)
				_, ok := liveCache[key]
				if !ok {
					log.WithFields(log.Fields{
						"videoID": live.VideoID,
						"agency":  Agency.GroupName,
						"vtuber":  live.Member.Name,
					}).Info("Found upcoming on db but not in cache,send it to cache now")
					err := live.SendToCache(false)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}
