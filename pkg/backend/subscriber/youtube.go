package subscriber

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"
	log "github.com/sirupsen/logrus"
)

func CheckYtSubsCount() {
	var YTstate Subs
	for _, Group := range engine.GroupData {
		var VtubChannel []string
		Names := database.GetMembers(Group.ID)
		for i, Member := range Names {
			if Member.YoutubeID != "" {
				VtubChannel = append(VtubChannel, Member.YoutubeID)
			}

			if i == 24 || i == len(Names)-1 {
				body, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+strings.Join(VtubChannel, ",")+"&key="+engine.GetYtToken(), nil)
				if err != nil {
					log.Error(err, string(body))
					return
				}
				err = json.Unmarshal(body, &YTstate)
				if err != nil {
					log.Error(err)
				}
				for _, Name2 := range Names {
					for _, Item := range YTstate.Items {
						if Name2.YoutubeID == Item.ID && !Item.Statistics.HiddenSubscriberCount {
							YtSubsDB := Name2.GetSubsCount()
							YTSubscriberCount, err := strconv.Atoi(Item.Statistics.SubscriberCount)
							if err != nil {
								log.Error(err)
							}
							SendNotif := func(SubsCount string) {
								Color, err := engine.GetColor(config.TmpDir, Name2.YoutubeAvatar)
								if err != nil {
									log.Error(err)
								}
								SendNude(engine.NewEmbed().
									SetAuthor(Group.GroupName, Group.IconURL, "https://www.youtube.com/channel/"+Name2.YoutubeID+"?sub_confirmation=1").
									SetTitle(engine.FixName(Name2.EnName, Name2.JpName)).
									SetThumbnail(config.YoutubeIMG).
									SetDescription("Congratulation for "+SubsCount+" subscriber").
									SetImage(Name2.YoutubeAvatar).
									AddField("Viewers", Item.Statistics.ViewCount).
									AddField("Videos", Item.Statistics.VideoCount).
									InlineAllFields().
									SetURL("https://www.youtube.com/channel/"+Name2.YoutubeID+"?sub_confirmation=1").
									SetColor(Color).MessageEmbed, Group, Name2.ID)
							}
							if YtSubsDB.YtSubs != YTSubscriberCount {
								if YTSubscriberCount >= 1000000 {
									for i := 0; i < 10000001; i += 1000000 {
										if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
											SendNotif(strconv.Itoa(i))
										}
									}
								} else if YTSubscriberCount >= 100000 {
									for i := 0; i < 1000001; i += 100000 {
										if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
											SendNotif(strconv.Itoa(i))
										}
									}
								} else if YTSubscriberCount >= 10000 {
									for i := 0; i < 100001; i += 10000 {
										if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
											SendNotif(strconv.Itoa(i))
										}
									}
								} else if YTSubscriberCount >= 1000 {
									for i := 0; i < 10001; i += 1000 {
										if i == YTSubscriberCount && !Item.Statistics.HiddenSubscriberCount {
											SendNotif(strconv.Itoa(i))
										}
									}
								}

								log.WithFields(log.Fields{
									"Past Youtube subscriber":    YtSubsDB.YtSubs,
									"Current Youtube subscriber": YTSubscriberCount,
									"Vtuber":                     Name2.EnName,
								}).Info("Update Youtube subscriber")
								VideoCount, err := strconv.Atoi(Item.Statistics.VideoCount)
								if err != nil {
									log.Error(err)
								}
								ViewCount, err := strconv.Atoi(Item.Statistics.ViewCount)
								if err != nil {
									log.Error(err)
								}
								YtSubsDB.UpYtSubs(YTSubscriberCount).
									UpYtVideo(VideoCount).
									UpYtViews(ViewCount).
									UpdateSubs("yt")
							}
						}
					}
				}
			}
			if i%10 == 0 {
				time.Sleep(3 * time.Second)
			}
		}
	}
}
