package main

import (
	"net/url"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	log "github.com/sirupsen/logrus"
)

func CheckPixivLewd() {
	for _, Group := range VtubersData.VtuberData {
		var wg sync.WaitGroup
		for i, Member := range Group.Members {
			wg.Add(1)
			go func(wg *sync.WaitGroup, Member database.Member) {
				defer wg.Done()
				FixFanArt := &database.DataFanart{
					Member: Member,
					Group:  Group,
					Lewd:   true,
				}
				var wg2 sync.WaitGroup
				wg2.Add(3)

				go func(w *sync.WaitGroup) {
					defer w.Done()
					if Member.JpName != "" {
						log.WithFields(log.Fields{
							"Member": Member.JpName,
							"Group":  Group.GroupName,
							"Lewd":   true,
						}).Info("Start curl pixiv")
						URLJP := GetPixivLewdURL(url.QueryEscape(Member.JpName))
						err := Pixiv(URLJP, FixFanArt, true)
						if err != nil {
							log.Error(err)
						}
					}
				}(&wg2)
				time.Sleep(1 * time.Second)

				go func(w *sync.WaitGroup) {
					defer w.Done()
					if Member.EnName == Member.Name {
						if Member.EnName != "" {
							log.WithFields(log.Fields{
								"Member": Member.EnName,
								"Group":  Group.GroupName,
								"Lewd":   true,
							}).Info("Start curl pixiv")
							URLEN := GetPixivLewdURL(engine.UnderScoreName(Member.EnName))
							err := Pixiv(URLEN, FixFanArt, true)
							if err != nil {
								log.Error(err)
							}

						}
					} else {
						if Member.EnName != "" {
							log.WithFields(log.Fields{
								"Member": Member.EnName,
								"Group":  Group.GroupName,
								"Lewd":   true,
							}).Info("Start curl pixiv")
							URLEN := GetPixivLewdURL(engine.UnderScoreName(Member.EnName))
							err := Pixiv(URLEN, FixFanArt, true)
							if err != nil {
								log.Error(err)
							}

						}
					}
				}(&wg2)
				time.Sleep(1 * time.Second)

				go func(w *sync.WaitGroup) {
					defer w.Done()
					if Member.TwitterLewd != "" {
						log.WithFields(log.Fields{
							"Member": Member.Name,
							"Group":  Group.GroupName,
							"Lewd":   true,
						}).Info("Start curl pixiv")
						URL := GetPixivLewdURL(engine.UnderScoreName(Member.TwitterLewd[1:]))
						err := Pixiv(URL, FixFanArt, true)
						if err != nil {
							log.Error(err)
						}
					}
				}(&wg2)
				wg2.Wait()

			}(&wg, Member)
			if i%4 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()
	}
}
