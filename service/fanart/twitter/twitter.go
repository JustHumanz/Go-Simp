package twitter

import (
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	configfile  config.ConfigFile
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	b.AddFunc(config.TwitterFanart, CheckNew)
	log.Info("Enable Twitter fanart module")
}

//CheckNew Check new fanart
func CheckNew() {
	wg := new(sync.WaitGroup)
	for _, GroupData := range VtubersData.VtuberData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			Fanarts, err := CreatePayload(Group, config.Scraper, config.GoSimpConf.LimitConf.TwitterFanart)
			if err != nil {
				log.WithFields(log.Fields{
					"Group": Group.GroupName,
				}).Error(err)
			} else {
				SendFanart(Fanarts, Group)
			}
		}(GroupData, wg)
	}
	wg.Wait()
}
