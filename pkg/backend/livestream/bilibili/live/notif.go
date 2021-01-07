package live

import (
	"errors"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	"github.com/hako/durafmt"
)

//push to discord channel
func (Data *LiveBili) Crotttt(GroupIcon string) error {
	BiliBiliAccount := "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
	BiliBiliURL := "https://live.bilibili.com/" + strconv.Itoa(Data.RoomData.LiveRoomID)
	Color, err := engine.GetColor(config.TmpDir, Data.RoomData.Thumbnail)
	if err != nil {
		return err
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
			AddField("Online", engine.NearestThousandFormat(float64(Data.RoomData.Online))).
			InlineAllFields().
			//AddField("Rank",Data).
			SetFooter(Data.RoomData.ScheduledStart.In(loc).Format(time.RFC822), config.BiliBiliIMG).
			SetColor(Color).MessageEmbed
	} else {
		return errors.New("it's not live")
	}

	MemberID := Data.Member.ID
	id, DiscordChannelID := database.ChannelTag(MemberID, 2, "")
	for i, DiscordChannel := range DiscordChannelID {
		UserTagsList := database.GetUserList(id[i], MemberID)
		if UserTagsList != nil {
			_, err := Bot.ChannelMessageSendEmbed(DiscordChannel, Data.Embed)
			if err != nil {
				return err
			}
			_, err = Bot.ChannelMessageSend(DiscordChannel, Data.Member.Name+"Live right now\nUserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				return err
			}
		}
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

	return nil
}
