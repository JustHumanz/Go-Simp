package subscriber

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func CheckYtSubsCount() {
	for k := 0; k < len(Data); k++ {
		var (
			YTstate     Subs
			Counter     int
			YtSubsDB    []database.MemberSubs
			ChannelList []string
			Names       []database.Name
		)
		tmp := database.GetName(Data[k].ID)
		for _, Name := range tmp {
			Counter++
			YtSubsDB = append(YtSubsDB, Name.GetSubsCount())
			ChannelList = append(ChannelList, strings.Split(Name.YoutubeID, "\n")[0])
			Names = append(Names, Name)
			if Counter == 24 || Counter == len(tmp)-1 {
				body, err := engine.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id="+strings.Join(ChannelList, ",")+"&key="+yttoken, nil)
				if err != nil {
					log.Error(err, string(body))
					oldtoken := yttoken
					yttoken = engine.ChangeToken(yttoken)
					log.WithFields(log.Fields{
						"Old Token": oldtoken,
						"New Token": yttoken,
					}).Warn("Token out of limit,try to change")
				}
				err = json.Unmarshal(body, &YTstate)
				if err != nil {
					log.Error(err)
				} else {
					for z, ChannelInfo := range YTstate.Items {
						if !ChannelInfo.Statistics.HiddenSubscriberCount {
							YTSubscriberCount, err := strconv.Atoi(ChannelInfo.Statistics.SubscriberCount)
							if err != nil {
								log.Error(err)
							}
							if YTSubscriberCount != YtSubsDB[z].YtSubs {
								if YTSubscriberCount <= 100000 {
									Avatar := Names[z].YoutubeAvatar
									Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
									if err != nil {
										log.Error(err)
									}
									for i := 0; i < 1000001; i += 100000 {
										if i == YTSubscriberCount && YTSubscriberCount != 0 {
											SendNude(engine.NewEmbed().
												SetAuthor(Data[k].NameGroup, Data[k].IconURL, "https://www.youtube.com/channel/"+Names[z].YoutubeID+"?sub_confirmation=1").
												SetTitle(engine.FixName(Names[z].EnName, Names[z].JpName)).
												SetThumbnail(config.YoutubeIMG).
												SetDescription("Congratulation for "+strconv.Itoa(i)+" subs").
												SetImage(Avatar).
												AddField("Views", ChannelInfo.Statistics.ViewCount).
												AddField("Videos", ChannelInfo.Statistics.VideoCount).
												InlineAllFields().
												SetURL("https://www.youtube.com/channel/"+Names[z].YoutubeID).
												SetColor(Color).MessageEmbed, Data[k])
										}
									}
								} else {
									Avatar := Names[z].YoutubeAvatar
									Color, err := engine.GetColor("/tmp/yt.tmp", Avatar)
									if err != nil {
										log.Error(err)
									}
									for i := 0; i < 10001; i += 1000 {
										if i == YTSubscriberCount && YTSubscriberCount != 0 {
											SendNude(engine.NewEmbed().
												SetAuthor(Data[k].NameGroup, Data[k].IconURL, "https://www.youtube.com/channel/"+Names[z].YoutubeID+"?sub_confirmation=1").
												SetTitle(engine.FixName(Names[z].EnName, Names[z].JpName)).
												SetThumbnail(config.YoutubeIMG).
												SetDescription("Congratulation for "+strconv.Itoa(i)+" subs").
												SetImage(Avatar).
												AddField("Views", ChannelInfo.Statistics.ViewCount).
												AddField("Videos", ChannelInfo.Statistics.VideoCount).
												InlineAllFields().
												SetURL("https://www.youtube.com/channel/"+Names[z].YoutubeID).
												SetColor(Color).MessageEmbed, Data[k])
										}
									}
								}
							}
							log.WithFields(log.Fields{
								"Past Youtube subscriber":    YtSubsDB[z].YtSubs,
								"Current Youtube subscriber": YTSubscriberCount,
								"Vtuber":                     Names[z].EnName,
							}).Info("Update Youtube subscriber")
							VideoCount, err := strconv.Atoi(ChannelInfo.Statistics.VideoCount)
							if err != nil {
								log.Error(err)
							}
							ViewCount, err := strconv.Atoi(ChannelInfo.Statistics.ViewCount)
							if err != nil {
								log.Error(err)
							}
							YtSubsDB[z].YtSubs = YTSubscriberCount
							YtSubsDB[z].YtVideos = VideoCount
							YtSubsDB[z].YtViews = ViewCount

							YtSubsDB[z].UpdateSubs("yt")
						}
					}
				}
				Counter = 0
				YtSubsDB = nil
				ChannelList = nil
				Names = nil
			}
		}
	}
}
