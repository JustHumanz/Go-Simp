package lewd

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
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	log.Info("Enable lewd fanart module")
	b.AddFunc(config.DanbooruFanart, CheckLewd)
}

func CheckLewd() {
	log.Info("Start Checking Danbooru lewd")
	for _, Group := range VtubersData.VtuberData {
		GetDan(Group)
	}
}
