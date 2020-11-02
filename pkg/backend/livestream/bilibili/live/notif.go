package live

import (
	"strconv"
	"strings"
	"time"

	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

//push to discord channel
func (Data *LiveBili) Tamod() {
	MemberID := Data.Member.ID
	Bot := runner.Bot
	id, DiscordChannelID := database.ChannelTag(MemberID, 2)
	for i, DiscordChannel := range DiscordChannelID {
		UserTagsList := database.GetUserList(id[i], MemberID)
		if UserTagsList != nil {
			msg, err := Bot.ChannelMessageSendEmbed(DiscordChannel, Data.Embed)
			if err != nil {
				log.Error(msg, err)
			}
			msg, err = Bot.ChannelMessageSend(DiscordChannel, "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
			}
		}
	}
}

//make a embed
func (Data *LiveBili) Crotttt(GroupIcon string) *LiveBili {
	BiliBiliAccount := "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
	BiliBiliURL := "https://live.bilibili.com/" + strconv.Itoa(Data.RoomData.LiveRoomID)
	Online := strconv.Itoa(Data.RoomData.Online)
	Color, err := engine.GetColor("/tmp/bilThum", Data.RoomData.Thumbnail)
	if err != nil {
		log.Error(err)
	}

	if Data.RoomData.Status == "Live" {
		Data.Embed = engine.NewEmbed().
			SetAuthor(engine.FixName(Data.Member.EnName, Data.Member.JpName), Data.Member.BiliBiliAvatar, BiliBiliAccount).
			SetTitle("Live right now").
			SetThumbnail(GroupIcon).
			SetDescription(Data.RoomData.Description).
			SetImage(Data.RoomData.Thumbnail).
			SetURL(BiliBiliURL).
			AddField("Start live", durafmt.Parse(time.Now().In(loc).Sub(Data.RoomData.ScheduledStart.In(loc))).LimitFirstN(2).String()+" Ago").
			AddField("Online", Online).
			InlineAllFields().
			//AddField("Rank",Data).
			SetFooter(Data.RoomData.ScheduledStart.In(loc).Format(time.RFC822), config.BiliBiliIMG).
			SetColor(Color).MessageEmbed
	}

	/* else if DataRoom.Status == "Upcoming" {
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
	*/
	return Data

}
