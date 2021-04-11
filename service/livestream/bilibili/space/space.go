package space

import (
	"strconv"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
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
	b.AddFunc(config.BiliBiliSpace, CheckSpaceVideo)
	log.Info("Enable space bilibili module")
}

//CheckSpaceVideo
func CheckSpaceVideo() {
	for _, GroupData := range VtubersData.VtuberData {
		if GroupData.GroupName != "Hololive" {
			wg := new(sync.WaitGroup)
			for i, MemberData := range GroupData.Members {
				wg.Add(1)
				go func(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
					defer wg.Done()
					if Member.BiliBiliID != 0 {
						log.WithFields(log.Fields{
							"Group":      Group.GroupName,
							"Vtuber":     Member.EnName,
							"BiliBiliID": Member.BiliBiliID,
						}).Info("Checking Space BiliBili")

						Group.RemoveNillIconURL()
						SpaceBiliLimit := strconv.Itoa(config.GoSimpConf.LimitConf.SpaceBiliBili)
						Data := &database.LiveStream{
							Member: Member,
							Group:  Group,
						}
						CheckSpace(Data, SpaceBiliLimit)

					}
				}(GroupData, MemberData, wg)
				if i%config.Waiting == 0 {
					wg.Wait()
				}
			}
			wg.Wait()
		}
	}
}
