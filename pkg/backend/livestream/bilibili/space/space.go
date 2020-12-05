package space

import (
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
)

var (
	loc *time.Location
)

func CheckVideo() {
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
	wg := new(sync.WaitGroup)
	for _, Group := range engine.GroupData {
		for _, Member := range database.GetMembers(Group.ID) {
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
					Data.Check("3").SendNude()

				}
			}(Group, Member, wg)
			time.Sleep(time.Duration(rand.Intn(config.RandomSleep-400)+400) * time.Millisecond)
		}
	}
	wg.Wait()
}
