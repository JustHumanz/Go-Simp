package space

import (
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"

	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
)

var (
	loc *time.Location
	Bot *discordgo.Session
)

//Start start twitter module
func Start(BotInit *discordgo.Session, cronInit *cron.Cron) {
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
	Bot = BotInit
	cronInit.AddFunc(config.BiliBiliSpace, CheckSpaceVideo)
	log.Info("Enable space bilibili module")
}

func CheckSpaceVideo() {
	for _, GroupData := range engine.GroupData {
		if GroupData.GroupName != "Hololive" {
			wg := new(sync.WaitGroup)
			for i, MemberData := range database.GetMembers(GroupData.ID) {
				wg.Add(1)
				go func(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
					defer wg.Done()
					if Member.BiliBiliID != 0 {
						log.WithFields(log.Fields{
							"Group":      Group.GroupName,
							"Vtuber":     Member.EnName,
							"BiliBiliID": Member.BiliBiliID,
						}).Info("Check Space")

						GroupIcon := ""
						if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
							GroupIcon = ""
						} else {
							GroupIcon = Group.IconURL
						}
						Data := &CheckSctruct{
							SpaceID:    Member.BiliBiliID,
							MemberID:   Member.ID,
							GroupIcon:  GroupIcon,
							MemberName: engine.FixName(Member.JpName, Member.EnName),
							MemberFace: Member.BiliBiliAvatar,
							MemberUrl:  "https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID),
						}
						Data.Check(strconv.Itoa(config.BotConf.LimitConf.SpaceBiliBili)).SendNude()

					}
				}(GroupData, MemberData, wg)
				if i%5 == 0 {
					time.Sleep(3 * time.Second)
				}
			}
			wg.Wait()
		}
	}
}
