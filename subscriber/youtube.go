package subscriber

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func CheckYtSubsCount() {
	for k := 0; k < len(Data); k++ {
		for _, Name := range database.GetName(Data[k].ID) {
			var (
				ytstate Subs
			)
			head := []string{"Referer", "https://akshatmittal.com/youtube-realtime/"}
			body, err := engine.Curl("https://counts.live/api/youtube-subscriber-count/"+strings.Split(Name.YoutubeID, "\n")[0]+"/live", head)
			if err != nil {
				log.Error(err, string(body))
			}
			err = json.Unmarshal(body, &ytstate)
			if err != nil {
				log.Error(err)
			}
			//Check Subs count
			YtSubsDB := Name.GetSubsCount()

			if ytstate.Data.Subscribers != YtSubsDB.YtSubs {
				if ytstate.Data.Subscribers <= 100000 {
					Avatar := Name.YoutubeAvatar
					Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
					if err != nil {
						log.Error(err)
					}
					for i := 0; i < 1000001; i += 100000 {
						if i == ytstate.Data.Subscribers && ytstate.Data.Subscribers != 0 {
							SendNude(engine.NewEmbed().
								SetAuthor(Data[k].NameGroup, Data[k].IconURL, "https://www.youtube.com/channel/"+Name.YoutubeID+"?sub_confirmation=1").
								SetTitle(engine.FixName(Name.EnName, Name.JpName)).
								SetThumbnail(config.YoutubeIMG).
								SetDescription("Congratulation for "+strconv.Itoa(i)+" subs").
								SetImage(Avatar).
								AddField("Views", strconv.Itoa(ytstate.Data.Views)).
								AddField("Videos", strconv.Itoa(ytstate.Data.Videos)).
								InlineAllFields().
								SetURL("https://www.youtube.com/channel/"+Name.YoutubeID).
								SetColor(Color).MessageEmbed, Data[k])
						}
					}
				} else {
					Avatar := Name.YoutubeAvatar
					Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
					if err != nil {
						log.Error(err)
					}
					for i := 0; i < 10001; i += 1000 {
						if i == ytstate.Data.Subscribers && ytstate.Data.Subscribers != 0 {
							SendNude(engine.NewEmbed().
								SetAuthor(Data[k].NameGroup, Data[k].IconURL, "https://www.youtube.com/channel/"+Name.YoutubeID+"?sub_confirmation=1").
								SetTitle(engine.FixName(Name.EnName, Name.JpName)).
								SetThumbnail(config.YoutubeIMG).
								SetDescription("Congratulation for "+strconv.Itoa(i)+" subs").
								SetImage(Avatar).
								AddField("Views", strconv.Itoa(ytstate.Data.Views)).
								AddField("Videos", strconv.Itoa(ytstate.Data.Videos)).
								InlineAllFields().
								SetURL("https://www.youtube.com/channel/"+Name.YoutubeID).
								SetColor(Color).MessageEmbed, Data[k])
						}
					}
				}
			}
			log.WithFields(log.Fields{
				"Past Youtube subscriber":    YtSubsDB.YtSubs,
				"Current Youtube subscriber": ytstate.Data.Subscribers,
				"Vtuber":                     Name.EnName,
			}).Info("Update Youtube subscriber")
			YtSubsDB = database.MemberSubs{
				YtSubs:   ytstate.Data.Subscribers,
				YtVideos: ytstate.Data.Videos,
				YtViews:  ytstate.Data.Views,
			}
			YtSubsDB.UpdateSubs("yt")
			time.Sleep(1 * time.Second)
		}
	}
}
