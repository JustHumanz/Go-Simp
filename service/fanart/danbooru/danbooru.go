package danbooru

import (
	"github.com/JustHumanz/Go-Simp/pkg/config"
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
	FistRunning *bool
	DataStore   = make(map[string][]int)
	LewdPics    []int
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	tmp := true
	FistRunning = &tmp
	log.Info("Enable danbooru fanart module")
	CheckDan()
	b.AddFunc(config.DanbooruFanart, CheckDan)
}

func CheckDan() {
	log.Info("Start Checking Danbooru lewd")
	for i, Group := range VtubersData.VtuberData {
		GetDan(Group)
		if i == len(VtubersData.VtuberData)-1 && *FistRunning {
			tmp := false
			FistRunning = &tmp
			log.Info("Set FistRunning to false")
		}
	}
}
