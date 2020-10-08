package subscriber

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func CheckYtSubsCount() {
	for k := 0; k < len(Data); k++ {
		Names := database.GetName(Data[k].ID)
		for _, Name := range Names {
			var (
				ytstate Subs
			)
			head := []string{"Referer", "https://akshatmittal.com/youtube-realtime/"}
			body, err := engine.Curl("https://counts.live/api/youtube-subscriber-count/"+Name.YoutubeID+"/live", head)
			if err != nil {
				log.Error(err, string(body))
			}
			err = json.Unmarshal(body, &ytstate)
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
								SetThumbnail("https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/youtube.png").
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
								SetThumbnail("https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/youtube.png").
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
