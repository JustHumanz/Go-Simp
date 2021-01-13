package youtube

import (
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"

	log "github.com/sirupsen/logrus"
)

var (
	yttoken   string
	Ytwaiting = "???"
	Bot       *discordgo.Session
)

//Start start twitter module
func Start(BotInit *discordgo.Session, cronInit *cron.Cron) {
	Bot = BotInit
	cronInit.AddFunc(config.YoutubeCheckChannel, CheckYtSchedule)
	cronInit.AddFunc(config.YoutubePrivateSlayer, CheckPrivate)
	log.Info("Enable youtube module")
	//CheckYtScheduleTest("Hololive")
}

func CheckYtSchedule() {
	yttoken = engine.GetYtToken()
	for _, Group := range engine.GroupData {
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
	log.Info("Start Youtube checker")
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
		for _, Group := range engine.GroupData {
			for _, Member := range database.GetMembers(Group.ID) {
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
