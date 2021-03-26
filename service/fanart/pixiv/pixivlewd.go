package pixiv

import (
	"net/url"
	"sync"

	"github.com/JustHumanz/Go-Simp/pkg/database"
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
				if Member.JpName != "" {
					log.WithFields(log.Fields{
						"Member": Member.JpName,
						"Group":  Group.GroupName,
						"Lewd":   true,
					}).Info("Start curl lewd pixiv")
					URLJP := GetPixivLewdURL(url.QueryEscape(Member.JpName))
					err := Pixiv(URLJP, FixFanArt, true)
					if err != nil {
						log.Error(err)
					}
				}

				if Member.EnName == Member.Name {
					if Member.EnName != "" {
						log.WithFields(log.Fields{
							"Member": Member.EnName,
							"Group":  Group.GroupName,
							"Lewd":   true,
						}).Info("Start curl lewd pixiv")
						URLEN := GetPixivLewdURL(url.QueryEscape(Member.EnName))
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
						}).Info("Start curl lewd pixiv")
						URLEN := GetPixivLewdURL(url.QueryEscape(Member.EnName))
						err := Pixiv(URLEN, FixFanArt, true)
						if err != nil {
							log.Error(err)
						}

					}
					if Member.Name != "" {
						log.WithFields(log.Fields{
							"Member": Member.Name,
							"Group":  Group.GroupName,
							"Lewd":   true,
						}).Info("Start curl lewd pixiv")
						URL := GetPixivLewdURL(url.QueryEscape(Member.Name))
						err := Pixiv(URL, FixFanArt, true)
						if err != nil {
							log.Error(err)
						}
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
