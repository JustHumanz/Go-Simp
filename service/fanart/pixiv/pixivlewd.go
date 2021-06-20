package main

import (
	"context"
	"net/url"
	"sync"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	log "github.com/sirupsen/logrus"
)

func CheckPixivLewd() {
	for _, Group := range *GroupPayload {
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

				if Member.JpName != "" && Member.Region == "JP" {
					log.WithFields(log.Fields{
						"Member": Member.JpName,
						"Group":  Group.GroupName,
						"Lewd":   true,
					}).Info("Start curl pixiv")
					URLJP := GetPixivLewdURL(url.QueryEscape(Member.JpName))
					err := Pixiv(URLJP, FixFanArt, true)
					if err != nil {
						log.Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}
				} else if Member.EnName != "" && Member.Region != "JP" {
					log.WithFields(log.Fields{
						"Member": Member.EnName,
						"Group":  Group.GroupName,
						"Lewd":   true,
					}).Info("Start curl pixiv")
					URLEN := GetPixivLewdURL(engine.UnderScoreName(Member.EnName))
					err := Pixiv(URLEN, FixFanArt, true)
					if err != nil {
						log.Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}
				} else {
					log.WithFields(log.Fields{
						"Member": Member.EnName,
						"Group":  Group.GroupName,
						"Lewd":   true,
					}).Info("Start curl pixiv")
					URLEN := GetPixivLewdURL(engine.UnderScoreName(Member.Name))
					err := Pixiv(URLEN, FixFanArt, true)
					if err != nil {
						log.Error(err)
						gRCPconn.ReportError(context.Background(), &pilot.ServiceMessage{
							Message: err.Error(),
							Service: ModuleState,
						})
					}
				}

			}(&wg, Member)
			if i%4 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()
	}
}
