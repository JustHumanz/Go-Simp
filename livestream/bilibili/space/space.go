package space

import (
	"strconv"
	"sync"
	"time"

	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

var (
	loc *time.Location
)

func Start() {
	loc, _ = time.LoadLocation("Asia/Shanghai") /*Use CST*/
}

func CheckVideo() {
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
					Data := CheckSctruct{
						SpaceID:    Member.BiliBiliID,
						MemberID:   Member.ID,
						GroupIcon:  Group.IconURL,
						MemberName: engine.FixName(Member.JpName, Member.EnName),
						MemberFace: Member.BiliBiliAvatar,
						MemberUrl:  "https://space.bilibili.com/" + strconv.Itoa(Member.BiliBiliID),
					}
					Data.Check("10").SendNude()
					//time.Sleep(time.Duration(int64(rand.Intn((10-4)+4))) * time.Second)

				}
			}(Group, Member, wg)
		}
	}
	wg.Wait()
}
