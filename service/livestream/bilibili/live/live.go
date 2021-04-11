package live

import (
	"regexp"
	"strconv"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/livestream/notif"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	loc         *time.Location
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload) {
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
	Bot = a
	VtubersData = c
	b.AddFunc(config.BiliBiliLive, CheckLiveSchedule)
	log.Info("Enable Live BiliBili module")
}

func CheckLiveSchedule() {
	for _, GroupData := range VtubersData.VtuberData {
		var wg sync.WaitGroup
		if GroupData.GroupName != "Hololive" {
			for i, MemberData := range GroupData.Members {
				if MemberData.BiliRoomID != 0 {
					wg.Add(1)
					log.WithFields(log.Fields{
						"Group":  GroupData.GroupName,
						"Vtuber": MemberData.EnName,
					}).Info("Checking LiveBiliBili")
					go CheckBili(GroupData, MemberData, &wg)
					if i%10 == 0 {
						wg.Wait()
					}

				}
			}

		}
	}
}

func CheckBili(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
	defer wg.Done()
	if Member.BiliBiliID != 0 {
		var (
			ScheduledStart time.Time
		)

		LiveBiliDB, err := database.GetRoomData(Member.ID, Member.BiliRoomID)
		if err != nil {
			log.Error(err)
		}

		Status, err := GetRoomStatus(Member.BiliRoomID)
		if err != nil {
			log.Error(err)
		}

		if LiveBiliDB != nil {
			LiveBiliDB.AddMember(Member).AddGroup(Group)
			if Status.CheckScheduleLive() && LiveBiliDB.Status != config.LiveStatus {
				//Live
				if Status.Data.RoomInfo.LiveStartTime != 0 {
					ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0).In(loc)
				} else {
					ScheduledStart = time.Now().In(loc)
				}

				if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
					Group.IconURL = ""
				}

				log.WithFields(log.Fields{
					"Group":  Group.GroupName,
					"Vtuber": Member.EnName,
					"Start":  ScheduledStart,
				}).Info("Start live right now")

				LiveBiliDB.UpdateStatus(config.LiveStatus).
					UpdateSchdule(ScheduledStart).
					UpdateViewers(strconv.Itoa(Status.Data.RoomInfo.Online)).
					UpdateThumbnail(Status.Data.RoomInfo.Cover).
					UpdateTitle(Status.Data.RoomInfo.Title).
					AddMember(Member).
					AddGroup(Group).
					SetState(config.BiliLive)

				err := LiveBiliDB.UpdateLiveBili()
				if err != nil {
					log.Error(err)
				}

				notif.SendDude(LiveBiliDB, Bot)

			} else if !Status.CheckScheduleLive() && LiveBiliDB.Status == config.LiveStatus {
				log.WithFields(log.Fields{
					"Group":  Group.GroupName,
					"Vtuber": Member.EnName,
					"Start":  LiveBiliDB.Schedul,
				}).Info("Past live stream")
				engine.RemoveEmbed(strconv.Itoa(LiveBiliDB.Member.BiliRoomID), Bot)
				LiveBiliDB.UpdateStatus(config.PastStatus)

				err = LiveBiliDB.UpdateLiveBili()
				if err != nil {
					log.Error(err)
				}
			} else {
				LiveBiliDB.UpdateViewers(strconv.Itoa(Status.Data.RoomInfo.Online))
				err := LiveBiliDB.UpdateLiveBili()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
