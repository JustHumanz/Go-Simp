package danbooru

import (
	"encoding/json"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	log "github.com/sirupsen/logrus"
)

func GetDan(Data database.Group) {
	wg := new(sync.WaitGroup)
	for i, Mem := range Data.Members {
		wg.Add(1)
		if i%10 == 0 {
			time.Sleep(5 * time.Second)
		}
		go func(Member database.Member, w *sync.WaitGroup) {
			defer w.Done()
			log.WithFields(log.Fields{
				"Group":   Data.GroupName,
				"Vtubers": Member.Name,
			}).Info("Check lewd pic")

			var (
				DanPayload []Danbooru
			)
			databyte, err := network.Curl(config.DanbooruEndPoint+strings.Replace(Member.EnName, " ", "_", -1)+"&limit=10", nil)
			if err != nil {
				log.Error(err)
			}
			err = json.Unmarshal(databyte, &DanPayload)
			if err != nil {
				log.Error(err)
			}
			for _, Dan := range DanPayload {
				if Dan.CheckLewd() {
					if *FistRunning {
						DataStore[Member.Name] = append(DataStore[Member.Name], Dan.Source)
					} else {
						if Dan.IsNew(Member) && Dan.ParentID == nil {
							log.WithFields(log.Fields{
								"Group":      Data.GroupName,
								"Vtubers":    Member.Name,
								"DanbooruID": Dan.ID,
							}).Info("New Lewd pic")
							DataStore[Member.Name] = append(DataStore[Member.Name], Dan.Source)
							Dan.SendNotif(Data, Member)
						}
					}
				}
			}
		}(Mem, wg)
	}
	wg.Wait()
}

func (Data Danbooru) IsNew(Member database.Member) bool {
	if DataStore != nil {
		for _, v := range DataStore[Member.Name] {
			if v == Data.Source {
				return false
			}
		}
	} else {
		return false
	}
	return true
}

func (Data Danbooru) CheckLewd() bool {
	safebutcrott, _ := regexp.MatchString("(swimsuits|lingerie|pantyshot)", Data.TagString)
	if Data.Rating == "e" || Data.Rating == "q" || safebutcrott {
		return true
	}
	return false
}
