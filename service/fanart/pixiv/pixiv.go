package pixiv

import (
	"encoding/json"
	"sync"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
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

const (
	Artwork = "https://www.pixiv.net/ajax/search/artworks/Gawr_Gura?word=%s&order=date_d&mode=all&p=1&s_mode=s_tag_full&type=all&lang=en"
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	VtubersData = c
	configfile = d
	//b.AddFunc(config.TwitterFanart, CheckNew)
	log.Info("Enable Twitter fanart module")
}

//CheckNew Check new fanart
func CheckPixiv() {
	wg := new(sync.WaitGroup)
	var Art PixivArtworks
	for _, GroupData := range VtubersData.VtuberData {
		wg.Add(1)
		go func(Group database.Group, wg *sync.WaitGroup) {
			defer wg.Done()
			for _, Member := range Group.Members {
				if Member.Region == "JP" {
					URL := GetPixivURL(Member.JpName)
					body, err := network.Curl(URL, nil)
					if err != nil {
						log.Error(err)
					}
					err = json.Unmarshal(body, &Art)
					if err != nil {
						log.Error(err)
					}
				} else {

				}
			}
		}(GroupData, wg)
	}
	wg.Wait()
}

func GetPixivURL(str string) string {
	return "https://www.pixiv.net/ajax/search/artworks/" + str + "?word=" + str + "&order=date_d&mode=all&p=1&s_mode=s_tag_full&type=all&lang=en"
}
