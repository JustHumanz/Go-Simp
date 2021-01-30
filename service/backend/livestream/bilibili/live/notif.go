package live

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

//push to discord channel
func (Data *LiveBili) Crotttt() error {
	BiliBiliAccount := "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
	BiliBiliURL := "https://live.bilibili.com/" + strconv.Itoa(Data.RoomData.LiveRoomID)
	Color, err := engine.GetColor(config.TmpDir, Data.RoomData.Thumbnail)
	if err != nil {
		return err
	}
	BiliBiliRoomID := strconv.Itoa(Data.RoomData.LiveRoomID)

	User := &database.UserStruct{
		Human:    true,
		Reminder: 0,
	}
	if Data.RoomData.Status == "Live" {
		Data.Embed = engine.NewEmbed().
			SetAuthor(engine.FixName(Data.Member.EnName, Data.Member.JpName), Data.Member.BiliBiliAvatar, BiliBiliAccount).
			SetTitle("Live right now").
			SetThumbnail(Data.Group.IconURL).
			SetDescription(Data.RoomData.Description).
			SetImage(Data.RoomData.Thumbnail).
			SetURL(BiliBiliURL).
			AddField("Start live", durafmt.Parse(time.Now().In(loc).Sub(Data.RoomData.ScheduledStart.In(loc))).LimitFirstN(1).String()+" Ago").
			AddField("Online", engine.NearestThousandFormat(float64(Data.RoomData.Online))).
			InlineAllFields().
			//AddField("Rank",Data).
			SetFooter(Data.RoomData.ScheduledStart.In(loc).Format(time.RFC822), config.BiliBiliIMG).
			SetColor(Color).MessageEmbed
	} else {
		return errors.New("it's not live")
	}

	MemberID := Data.Member.ID
	//id, DiscordChannelID
	var (
		wg          sync.WaitGroup
		ChannelData = database.ChannelTag(MemberID, 2, "")
	)
	for i, v := range ChannelData {
		wg.Add(1)
		go func(Channel database.DiscordChannel, wg *sync.WaitGroup) {
			defer wg.Done()
			UserTagsList := Channel.GetUserList(context.Background())
			if UserTagsList == nil && Data.Group.GroupName != "Independen" {
				UserTagsList = []string{"_"}
			} else if UserTagsList == nil && Data.Group.GroupName == "Independen" {
				return
			}

			MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, Data.Embed)
			if err != nil {
				log.Error(err)
			}
			Msg := "Push " + config.BotConf.Emoji.Livestream[0] + " to add you in `" + Data.Member.Name + "` ping list\nPush " + config.BotConf.Emoji.Livestream[1] + " to remove you from ping list"
			MsgTxt, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` Live right now\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
			if err != nil {
				log.Error(err)
			} else {
				if Channel.Dynamic {
					log.WithFields(log.Fields{
						"DiscordChannel": Channel.ChannelID,
						"VtuberGroupID":  Data.Group.ID,
						"BiliBiliRoomID": BiliBiliRoomID,
					}).Info("Set dynamic mode")
					Channel.SetVideoID(BiliBiliRoomID).
						SetMsgEmbedID(MsgEmbed.ID).
						SetMsgTextID(MsgTxt.ID).
						PushReddis()
				}

				User.SetDiscordChannelID(Channel.ChannelID).
					SetGroup(Data.Group).
					SetMember(Data.Member)
				err = User.SendToCache(MsgTxt.ID)
				if err != nil {
					log.Error(err)
				}
				err = engine.Reacting(map[string]string{
					"ChannelID": Channel.ChannelID,
					"State":     "Youtube",
					"MessageID": MsgTxt.ID,
				}, Bot)
				if err != nil {
					log.Error(err)
				}
			}

		}(v, &wg)
		//Wait every ge 10 discord channel
		if i%10 == 0 {
			wg.Wait()
		}
	}
	wg.Wait()

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
