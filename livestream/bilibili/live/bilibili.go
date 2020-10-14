package bilibili

import (
	"sync"
	"time"

	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	loc        *time.Location
	BotSession *discordgo.Session
)

func Start(Bot *discordgo.Session) {
	BotSession = Bot
	loc, _ = time.LoadLocation("Asia/Shanghai")
	//go BotSession.AddHandler(Message)
	//go BotSession.AddHandler(Space)
}

func CheckSchedule() {
	log.Info("Start check Schedule")
	wg := new(sync.WaitGroup)
	for _, Group := range database.GetGroup() {
		for _, Member := range database.GetName(Group.ID) {
			wg.Add(1)
			go func(Group database.GroupName, Member database.Name, wg *sync.WaitGroup) {
				defer wg.Done()
				if Member.BiliBiliID != 0 {
					var (
						ScheduledStart time.Time
					)
					DataDB := database.GetRoomData(Member.ID, Member.BiliRoomID)
					Status := GetRoomStatus(Member.BiliRoomID)
					if Status.CheckScheduleLive() && DataDB.Status != "Live" {
						//Live
						if Status.Data.RoomInfo.LiveStartTime != 0 {
							ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0).In(loc)
						} else {
							ScheduledStart = time.Now().In(loc)
						}

						log.WithFields(log.Fields{
							"Group":  Group.NameGroup,
							"Vtuber": engine.FixName(Member.EnName, Member.JpName),
							"Start":  ScheduledStart,
						}).Info("Start live right now")
						Data := LiveBili{
							RoomData: database.LiveBiliDB{
								LiveRoomID:     Member.BiliRoomID,
								Status:         "Live",
								Title:          Status.Data.RoomInfo.Title,
								Thumbnail:      Status.Data.RoomInfo.Cover,
								Description:    Status.Data.NewsInfo.Content,
								PublishedAt:    time.Time{},
								ScheduledStart: ScheduledStart,
								Online:         Status.Data.RoomInfo.Online,
							},
							Face:       Status.Data.AnchorInfo.BaseInfo.Face,
							VtuberName: engine.FixName(Member.EnName, Member.JpName),
							BiliBiliID: Member.BiliBiliID,
						}
						Data.RoomData.UpdateLiveBili(Member.ID)
						Data.Crotttt(Group.IconURL).Tamod(Member.ID)

					} else if !Status.CheckScheduleLive() && DataDB.Status == "Live" {
						//prob past
						log.WithFields(log.Fields{
							"Group":  Group.NameGroup,
							"Vtuber": engine.FixName(Member.EnName, Member.JpName),
							"Start":  ScheduledStart,
						}).Info("Past live stream")
						DataDB.Status = "Past"
						DataDB.UpdateLiveBili(Member.ID)
					} else {
						//update online
						log.WithFields(log.Fields{
							"Group":  Group.NameGroup,
							"Vtuber": engine.FixName(Member.EnName, Member.JpName),
						}).Info("Update LiveBiliBili")
						DataDB.Online = Status.Data.RoomInfo.Online
						DataDB.UpdateLiveBili(Member.ID)
					}
				}
				//time.Sleep(time.Duration(int64(rand.Intn((20-8)+8))) * time.Second)
			}(Group, Member, wg)
		}
	}
	wg.Wait()
}
