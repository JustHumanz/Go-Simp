package youtube

import (
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	BotSession *discordgo.Session
	yttoken    string
	Ytwaiting  = "???"
)

func Start(Bot *discordgo.Session) {
	BotSession = Bot
	log.Info("Youtube module ready")
	//CheckSchedule()
}

func CheckSchedule() {
	yttoken = engine.GetYtToken()
	for _, Group := range engine.GroupData {
		var wg sync.WaitGroup
		for _, Member := range database.GetName(Group.ID) {
			if Member.YoutubeID != "" {
				wg.Add(1)
				log.WithFields(log.Fields{
					"Vtube":        Member.EnName,
					"Youtube ID":   Member.YoutubeID,
					"Vtube Region": Member.Region,
				}).Info("Checking Youtube")
				go Filter(Member, Group, &wg)
				time.Sleep(time.Duration(rand.Intn(config.RandomSleep-5)+5) * time.Millisecond)
			}
		}
		wg.Wait()
	}
}

func GetWaiting(VideoID string) (string, error) {
	var (
		bit     []byte
		curlerr error
		urls    = "https://www.youtube.com/watch?v=" + VideoID
	)
	bit, curlerr = engine.Curl(urls, nil)
	if curlerr != nil || bit == nil {
		bit, curlerr = engine.CoolerCurl(urls, nil)
		if curlerr != nil {
			return Ytwaiting, curlerr
		} else {
			log.WithFields(log.Fields{
				"Request": VideoID,
				"Func":    "YtGetWaiting",
			}).Info("Successfully use multi tor")
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
	log.Info("Start Check video")
	var (
		wg sync.WaitGroup
	)

	Check := func(Youtube database.YtDbData, wg *sync.WaitGroup) {
		defer wg.Done()

		var (
			tor bool
			err error
			url = "https://i3.ytimg.com/vi/" + Youtube.VideoID + "/hqdefault.jpg"
		)
		for {
			if tor {
				_, err = engine.CoolerCurl(url, nil)
			} else {
				_, err = engine.Curl(url, nil)
			}
			if Youtube.Status == "upcoming" && time.Now().Sub(Youtube.Schedul) > Youtube.Schedul.Sub(time.Now()) {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Member only video")
				Youtube.UpdateYt("past")
			} else if Youtube.Status == "live" && Youtube.Schedul.Minute() > time.Now().Add(2).Minute() {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Member only video")
				Youtube.UpdateYt("past")
			}

			if err != nil && strings.HasPrefix(err.Error(), "404") && Youtube.Status != "private" {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Private Video")
				Youtube.UpdateYt("private")
				break
			} else if err == nil && Youtube.Status == "private" {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("From Private Video to past")
				Youtube.UpdateYt("past")
				break
			} else if err != nil {
				log.Error(err)
				log.Info("Trying use tor")
				tor = true
				continue
			} else {
				log.WithFields(log.Fields{
					"VideoID": Youtube.VideoID,
				}).Info("Video was daijobu")
				break
			}
		}
	}

	log.Info("Start Check Private video")
	for _, Status := range []string{"upcoming", "past", "live", "private"} {
		for _, Group := range engine.GroupData {
			for i, Member := range database.GetName(Group.ID) {
				if i == 50 {
					break
				} else {
					YtData := database.YtGetStatus(0, Member.ID, Status, "")
					for j, Y := range YtData {
						Y.Status = Status
						wg.Add(1)
						go Check(Y, &wg)
						if j == 20 || j == len(YtData)-1 {
							wg.Wait()
						}
					}
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
	log.Info("Push to database")

	log.Info("Done")
}
