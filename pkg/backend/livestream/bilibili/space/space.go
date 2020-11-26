package space

import (
	"math/rand"
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
		for _, Member := range database.GetName(Group.ID) {
			wg.Add(1)
			go func(Group database.GroupName, Member database.Name, wg *sync.WaitGroup) {
				defer wg.Done()
				if Member.BiliBiliID != 0 {
					log.WithFields(log.Fields{
						"Group":      Group.NameGroup,
						"Vtuber":     Member.EnName,
						"BiliBiliID": Member.BiliBiliID,
					}).Info("Check Space")
					Data := &CheckSctruct{
						SpaceID:    Member.BiliBiliID,
						MemberID:   Member.ID,
						GroupIcon:  Group.IconURL,
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
