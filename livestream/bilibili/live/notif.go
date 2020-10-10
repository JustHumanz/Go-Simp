package bilibili

import (
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data LiveBili) Tamod(MemberID int64) {
	id, DiscordChannelID := database.ChannelTag(MemberID, 2)
	var (
		wg sync.WaitGroup
	)
	for i, DiscordChannel := range DiscordChannelID {
		UserTagsList := database.GetUserList(id[i], MemberID)
		wg.Add(1)
		go func(DiscordChannel string, wg *sync.WaitGroup) {
			defer wg.Done()
			BotSession.ChannelMessageSendEmbed(DiscordChannel, Data.Embed)
			if UserTagsList != nil {
				msg, err := BotSession.ChannelMessageSend(DiscordChannel, "UserTags: "+strings.Join(UserTagsList, " "))
				if err != nil {
					log.Error(msg, err)
				}
			}
		}(DiscordChannel, &wg)
	}
	wg.Wait()
}

func (Data LiveBili) Crotttt(GroupIcon string) LiveBili {
	BiliBiliAccount := "https://space.bilibili.com/" + strconv.Itoa(Data.BiliBiliID)
	BiliBiliURL := "https://live.bilibili.com/" + strconv.Itoa(Data.RoomData.LiveRoomID)
	Online := strconv.Itoa(Data.RoomData.Online)
	expiresAt := time.Now().In(loc)

	Color, err := engine.GetColor("/tmp/bilThum", Data.RoomData.Thumbnail)
	if err != nil {
		log.Error(err)
	}

	var (
		msg  string
		msg1 string
		msg2 string
		msg3 string
	)
	DataRoom := Data.RoomData
	if DataRoom.Status == "Live" {
		msg = "Start live"
		msg1 = durafmt.Parse(expiresAt.Sub(DataRoom.ScheduledStart.In(loc))).LimitFirstN(2).String() + " Ago"
		msg2 = "Live right now"
		msg3 = "Online : " + Online
	} else if DataRoom.Status == "Upcoming" {
		msg = "Start live in"
		msg1 = durafmt.Parse(DataRoom.ScheduledStart.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2).String()
		msg2 = "New live schedule"
		msg3 = DataRoom.ScheduledStart.In(loc).Format(time.RFC822)
	} else if DataRoom.Status == "Reminder" {
		msg = "Start live in"
		msg1 = durafmt.Parse(DataRoom.ScheduledStart.In(loc).Sub(time.Now().In(loc))).LimitFirstN(2).String()
		msg2 = "Reminder"
		msg3 = "Online : " + Online
	}

	Data.Embed = engine.NewEmbed().
		SetAuthor(Data.VtuberName, Data.Face, BiliBiliAccount).
		SetTitle(msg2).
		SetThumbnail(GroupIcon).
		SetDescription(DataRoom.Description).
		SetImage(DataRoom.Thumbnail, "image").
		SetURL(BiliBiliURL).
		AddField(msg, msg1).
		SetFooter(msg3, config.BiliBiliIMG).
		InlineAllFields().
		SetColor(Color).MessageEmbed
	return Data

}
