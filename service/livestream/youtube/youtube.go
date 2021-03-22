package youtube

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	log "github.com/sirupsen/logrus"
)

var (
	yttoken     *string
	Ytwaiting   = "???"
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	configfile  config.ConfigFile
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile) {
	Bot = a
	configfile = d
	b.AddFunc(config.YoutubeCheckChannel, CheckYtSchedule)
	b.AddFunc(config.YoutubeCheckUpcomingByTime, CheckYtByTime)
	b.AddFunc(config.YoutubePrivateSlayer, CheckPrivate)
	VtubersData = c
	log.Info("Enable youtube module")
	//CheckYtScheduleTest("Hololive")
}

func CheckYtSchedule() {
	yttoken = engine.GetYtToken()
	var (
		wg sync.WaitGroup
	)
	for _, Group := range VtubersData.VtuberData {
		wg.Add(1)
		go StartCheckYT(Group, &wg)
	}
	wg.Wait()
}

func CheckYtByTime() {
	for _, Group := range VtubersData.VtuberData {
		for _, Member := range Group.Members {
			if Member.YoutubeID != "" {
				log.WithFields(log.Fields{
					"Vtuber": Member.EnName,
					"Group":  Group.GroupName,
				}).Info("Checking Upcoming schedule")
				YoutubeStatus, err := database.YtGetStatus(0, Member.ID, config.UpcomingStatus, "")
				if err != nil {
					log.Error(err)
				}

				for _, Youtube := range YoutubeStatus {
					if time.Now().Sub(Youtube.Schedul) > Youtube.Schedul.Sub(time.Now()) {
						log.WithFields(log.Fields{
							"Vtuber":  Member.EnName,
							"Group":   Group.GroupName,
							"VideoID": Youtube.VideoID,
						}).Info("Vtuber upcoming schedule deadline,force change to live")

						yttoken = engine.GetYtToken()
						Data, err := YtAPI([]string{Youtube.VideoID})
						if err != nil {
							log.Error(err)
						}
						if len(Data.Items) > 0 {
							if Data.Items[0].Snippet.VideoStatus != "none" {
								if Data.Items[0].Statistics.ViewCount != "" {
									Youtube.UpdateViewers(Data.Items[0].Statistics.ViewCount)
								} else if Data.Items[0].Statistics.ViewCount == "0" && Youtube.Viewers == "0" {
									Viewers, err := GetWaiting(Youtube.VideoID)
									if err != nil {
										log.Error(err)
									}
									Youtube.UpdateViewers(Viewers)
								}

								if !Data.Items[0].LiveDetails.ActualStartTime.IsZero() {
									Youtube.UpdateSchdule(Data.Items[0].LiveDetails.ActualStartTime)
								}
								Key := "0" + strconv.Itoa(int(Member.ID)) + config.UpcomingStatus + ""
								err = database.RemoveYtCache(Key, context.Background())
								if err != nil {
									log.Error(err)
								}

								Youtube.UpdateYt(config.LiveStatus)
								err = SendNude(Youtube)
								if err != nil {
									log.Error(err)
								}
							}
						}
					}
					Youtube.UpdateStatus("reminder")
					err = SendNude(Youtube)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

func GetWaiting(VideoID string) (string, error) {
	var (
		bit     []byte
		curlerr error
		urls    = "https://www.youtube.com/watch?v=" + VideoID
	)
	bit, curlerr = network.Curl(urls, nil)
	if curlerr != nil || bit == nil {
		bit, curlerr = network.CoolerCurl(urls, nil)
		if curlerr != nil {
			return Ytwaiting, curlerr
		}
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return Ytwaiting, err
	}
	for _, element := range regexp.MustCompile(`(?m)videoViewCountRenderer.*?text([0-9\s]+).+(isLive\strue)`).FindAllStringSubmatch(reg.ReplaceAllString(string(bit), " "), -1) {
		tmp := strings.Replace(element[1], " ", "", -1)
		if tmp != "" {
			Ytwaiting = tmp
		}
	}
	return Ytwaiting, nil
}

func CheckPrivate() {
	log.Info("Start Video private slayer")
	Check := func(Youtube database.LiveStream) {
		if Youtube.Status == "upcoming" && time.Now().Sub(Youtube.Schedul) > Youtube.Schedul.Sub(time.Now()) {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Member only video")
			Youtube.UpdateYt("past")
			engine.RemoveEmbed(Youtube.VideoID, Bot)
		} else if Youtube.Status == "live" && Youtube.Viewers == "0" || Youtube.Status == "live" && int(math.Round(time.Now().Sub(Youtube.Schedul).Hours())) > 30 {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Member only video")
			Youtube.UpdateYt("past")
			engine.RemoveEmbed(Youtube.VideoID, Bot)
		}

		_, err := network.Curl("https://i3.ytimg.com/vi/"+Youtube.VideoID+"/hqdefault.jpg", nil)
		if err != nil && strings.HasPrefix(err.Error(), "404") && Youtube.Status != "private" {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Private Video")
			Youtube.UpdateYt("private")
		} else if err == nil && Youtube.Status == "private" {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("From Private Video to past")
			Youtube.UpdateYt("past")
		} else {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Video was daijobu")
		}
	}

	log.Info("Start Check Private video")
	for _, Status := range []string{config.UpcomingStatus, config.PastStatus, config.LiveStatus, config.PrivateStatus} {
		for _, Group := range VtubersData.VtuberData {
			for _, Member := range Group.Members {
				YtData, err := database.YtGetStatus(0, Member.ID, Status, "")
				if err != nil {
					log.Error(err)
				}
				for _, Y := range YtData {
					Y.Status = Status
					Check(Y)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	log.Info("Done")
}
