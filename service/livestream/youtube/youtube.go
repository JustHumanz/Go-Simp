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
				YoutubeStatus, err := database.YtGetStatus(0, Member.ID, "upcoming", "")
				if err != nil {
					log.Error(err)
				}

				for _, v := range YoutubeStatus {
					YoutubeData := &NotifStruct{
						Member: Member,
						Group:  Group,
						YtData: &v,
					}
					if time.Now().Sub(v.Schedul) > v.Schedul.Sub(time.Now()) {
						log.WithFields(log.Fields{
							"Vtuber":  Member.EnName,
							"Group":   Group.GroupName,
							"VideoID": v.VideoID,
						}).Info("Vtuber upcoming schedule deadline,force change to live")

						yttoken = engine.GetYtToken()
						Data, err := YtAPI([]string{v.VideoID})
						if err != nil {
							log.Error(err)
						}
						if len(Data.Items) > 0 {
							if Data.Items[0].Snippet.VideoStatus != "none" {
								if Data.Items[0].Statistics.ViewCount != "" {
									YoutubeData.UpYtView(Data.Items[0].Statistics.ViewCount)
								} else if Data.Items[0].Statistics.ViewCount == "0" && v.Viewers == "0" {
									Viewers, err := GetWaiting(v.VideoID)
									if err != nil {
										log.Error(err)
									}
									YoutubeData.UpYtView(Viewers)
								}

								if !Data.Items[0].LiveDetails.ActualStartTime.IsZero() {
									YoutubeData.UpYtSchedul(Data.Items[0].LiveDetails.ActualStartTime)
								}
								Key := "0" + strconv.Itoa(int(Member.ID)) + "upcoming" + ""
								err = database.RemoveYtCache(Key, context.Background())
								if err != nil {
									log.Error(err)
								}

								YoutubeData.ChangeYtStatus("live").UpdateYtDB().SendNude()
							}
						}
					}
					YoutubeData.ChangeYtStatus("reminder").SendNude()
				}
			}
		}
	}
}

/*
func CheckYtScheduleTest(GroupName string) {
	yttoken = engine.GetYtToken()
	for _, Group := range engine.GroupData {
		if GroupName == Group.GroupName {
			var wg sync.WaitGroup
			for i, Member := range database.GetMembers(Group.ID) {
				if Member.YoutubeID != "" {
					wg.Add(1)
					log.WithFields(log.Fields{
						"Vtuber":        Member.EnName,
						"Group":         Group.GroupName,
						"Youtube ID":    Member.YoutubeID,
						"Vtuber Region": Member.Region,
					}).Info("Checking Youtube")
					go StartCheckYT(Member, Group, &wg)
				}
				if i%10 == 0 {
					wg.Wait()
				}
			}
			wg.Wait()
		}
	}
}
*/

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
	Check := func(Youtube database.YtDbData) {
		if Youtube.Status == "upcoming" && time.Now().Sub(Youtube.Schedul) > Youtube.Schedul.Sub(time.Now()) {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Member only video")
			Youtube.UpdateYt("past")
			err := engine.RemoveEmbed(Youtube.VideoID, Bot)
			if err != nil {
				log.Error(err)
			}
		} else if Youtube.Status == "live" && Youtube.Viewers == "" || Youtube.Status == "live" && int(math.Round(time.Now().Sub(Youtube.Schedul).Hours())) > 30 {
			log.WithFields(log.Fields{
				"VideoID": Youtube.VideoID,
			}).Info("Member only video")
			Youtube.UpdateYt("past")
			err := engine.RemoveEmbed(Youtube.VideoID, Bot)
			if err != nil {
				log.Error(err)
			}
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
	for _, Status := range []string{"upcoming", "past", "live", "private"} {
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
