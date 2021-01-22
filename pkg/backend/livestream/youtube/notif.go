package youtube

import (
	"context"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/pkg/backend/livestream/bilibili/live"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (PushData *NotifStruct) SendNude() error {
	Status := PushData.YtData.Status
	Avatar := PushData.Member.YoutubeAvatar
	YtChannel := "https://www.youtube.com/channel/" + PushData.Member.YoutubeID + "?sub_confirmation=1"
	YtURL := "https://www.youtube.com/watch?v=" + PushData.YtData.VideoID
	loc := engine.Zawarudo(PushData.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(PushData.Member.EnName, PushData.Member.JpName)

	var (
		Timestart time.Time
		User      = &database.UserStruct{
			Human:    true,
			Reminder: 0,
		}
	)

	if match, _ := regexp.MatchString("404.jpg", PushData.Group.IconURL); match {
		PushData.Group.IconURL = ""
	}

	if !PushData.YtData.Schedul.IsZero() {
		Timestart = PushData.YtData.Schedul
	} else if PushData.YtData.Schedul.IsZero() && !PushData.YtData.Published.IsZero() {
		Timestart = PushData.YtData.Published
	} else if PushData.YtData.Schedul.IsZero() && PushData.YtData.Published.IsZero() {
		Timestart = time.Now()
	}

	Views, err := strconv.Atoi(PushData.YtData.Viewers)
	if err != nil {
		log.Error(err)
	}
	PushData.YtData.Viewers = engine.NearestThousandFormat(float64(Views))
	if Status == "upcoming" {
		Color, err := engine.GetColor(config.TmpDir, PushData.YtData.Thumb)
		if err != nil {
			return err
		}
		//id, DiscordChannelID
		var (
			wg          sync.WaitGroup
			ChannelData = database.ChannelTag(PushData.Member.ID, 2, "NewUpcoming")
		)
		for i, v := range ChannelData {
			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				UserTagsList := Channel.GetUserList(context.Background()) //database.GetUserList(Channel.ID, PushData.Member.ID)
				if UserTagsList == nil && PushData.Group.GroupName != "Independen" {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && PushData.Group.GroupName == "Independen" {
					return nil
				}
				msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Avatar, YtChannel).
					SetTitle("New upcoming Livestream").
					SetDescription(PushData.YtData.Title).
					SetImage(PushData.YtData.Thumb).
					SetThumbnail(PushData.Group.IconURL).
					SetURL(YtURL).
					AddField("Type ", PushData.YtData.Type).
					AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(2).String()).
					InlineAllFields().
					AddField("Waiting", PushData.YtData.Viewers+" Simps in ChatRoom").
					SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.WithFields(log.Fields{
						"Message":          msg,
						"ChannelID":        Channel.ID,
						"DiscordChannelID": Channel.ChannelID,
					}).Error(err)
					err = Channel.DelChannel(err.Error())
					if err != nil {
						return err
					}
				}
				msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` New upcoming Livestream\nUserTags: "+strings.Join(UserTagsList, " "))
				if err != nil {
					return err
				}
				return nil
			}(v, &wg)
			//Wait every ge 10 discord channel
			if i%10 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()

	} else if Status == "live" {
		Color, err := engine.GetColor(config.TmpDir, PushData.YtData.Thumb)
		if err != nil {
			return err
		}
		Bili := false
		if PushData.Member.BiliRoomID != 0 {
			LiveBili, err := live.GetRoomStatus(PushData.Member.BiliRoomID)
			if err != nil {
				return err
			}
			if LiveBili.CheckScheduleLive() {
				Bili = true
				database.SetRoomToLive(PushData.Member.ID)
			}
		}
		//id, DiscordChannelID
		var (
			wg          sync.WaitGroup
			ChannelData = database.ChannelTag(PushData.Member.ID, 2, "")
		)
		for i, v := range ChannelData {
			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				UserTagsList := Channel.GetUserList(context.Background()) //database.GetUserList(Channel.ID, PushData.Member.ID)
				if UserTagsList == nil && PushData.Group.GroupName != "Independen" {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && PushData.Group.GroupName == "Independen" {
					return nil
				}
				MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Avatar, YtChannel).
					SetTitle("Live right now").
					SetDescription(PushData.YtData.Title).
					SetImage(PushData.YtData.Thumb).
					SetThumbnail(PushData.Group.IconURL).
					SetURL(YtURL).
					AddField("Type ", PushData.YtData.Type).
					AddField("Start live", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(2).String()+" Ago").
					InlineAllFields().
					AddField("Viewers", PushData.YtData.Viewers+" simps").
					SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.WithFields(log.Fields{
						"Message":          MsgEmbed,
						"ChannelID":        Channel.ID,
						"DiscordChannelID": Channel.ChannelID,
					}).Error(err)
					err = Channel.DelChannel(err.Error())
					if err != nil {
						return err
					}
				}
				if Channel.Dynamic {
					log.WithFields(log.Fields{
						"DiscordChannel": Channel.ChannelID,
						"VtuberGroupID":  PushData.Group.ID,
						"YoutubeID":      PushData.YtData.ID,
					}).Info("Set dynamic mode")
					Channel.SetVideoID(PushData.YtData.VideoID).
						SetMsgEmbedID(MsgEmbed.ID)
				}
				Msg := "Push " + config.BotConf.Emoji.Livestream[0] + " to add you in `" + PushData.Member.Name + "` ping list\nPush " + config.BotConf.Emoji.Livestream[1] + " to remove you from ping list"
				MsgID := ""
				if Bili {
					msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` Live right now at BiliBili And Youtube\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
					if err != nil {
						return err
					}
					if Channel.Dynamic {
						Channel.SetMsgTextID(msg.ID).PushReddis()
					}
					MsgID = msg.ID

				} else {
					msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` Live right now\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
					if err != nil {
						return err
					}
					if Channel.Dynamic {
						Channel.SetMsgTextID(msg.ID).PushReddis()
					}
					MsgID = msg.ID
				}
				User.SetDiscordChannelID(Channel.ChannelID).
					SetGroup(PushData.Group).
					SetMember(PushData.Member)
				err = User.SendToCache(MsgID)
				if err != nil {
					return err
				}

				err = engine.Reacting(map[string]string{
					"ChannelID": Channel.ChannelID,
					"State":     "Youtube",
					"MessageID": MsgID,
				}, Bot)
				if err != nil {
					return err
				}
				return nil
			}(v, &wg)
			//Wait every ge 15 discord channel
			if i%15 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()

	} else if Status == "past" {
		Color, err := engine.GetColor(config.TmpDir, PushData.YtData.Thumb)
		if err != nil {
			return err
		}

		//id, DiscordChannelID
		var (
			wg          sync.WaitGroup
			ChannelData = database.ChannelTag(PushData.Member.ID, 2, "NotLiveOnly")
		)
		for i, v := range ChannelData {
			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				UserTagsList := Channel.GetUserList(context.Background()) //database.GetUserList(Channel.ID, PushData.Member.ID)
				if UserTagsList != nil {
					msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
						SetAuthor(VtuberName, Avatar, YtChannel).
						SetTitle("Uploaded a new video").
						SetDescription(PushData.YtData.Title).
						SetImage(PushData.YtData.Thumb).
						SetThumbnail(PushData.Group.IconURL).
						SetURL(YtURL).
						AddField("Type ", PushData.YtData.Type).
						AddField("Upload", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(2).String()+" Ago").
						AddField("Viewers", PushData.YtData.Viewers).
						AddField("Duration", PushData.YtData.Length).
						InlineAllFields().
						SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.WithFields(log.Fields{
							"Message":          msg,
							"ChannelID":        Channel.ID,
							"DiscordChannelID": Channel.ChannelID,
						}).Error(err)
						err = Channel.DelChannel(err.Error())
						if err != nil {
							return err
						}
					}
					msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` Uploaded a new video\nUserTags: "+strings.Join(UserTagsList, " "))
					if err != nil {
						return err
					}
				}
				return nil
			}(v, &wg)
			//Wait every ge 10 discord channel
			if i%10 == 0 {
				wg.Wait()
			}
		}
	} else if Status == "reminder" {
		UpcominginMinutes := int(math.Round(Timestart.In(loc).Sub(time.Now().In(loc)).Minutes()))
		//id, DiscordChannelID
		var Color int
		ChanelData := database.ChannelTag(PushData.Member.ID, 2, "")
		for _, Channel := range ChanelData {
			for ii := 0; ii < 70; ii++ {
				if UpcominginMinutes == ii {
					if Color != 0 {
						Color, err = engine.GetColor(config.TmpDir, PushData.YtData.Thumb)
						if err != nil {
							log.Error(err)
						}
					}

					UserTagsList := database.GetUserReminderList(int(Channel.ID), PushData.Member.ID, ii)
					LiveCount := durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()
					if UserTagsList != nil {
						MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
							SetAuthor(VtuberName, Avatar, YtChannel).
							SetTitle(PushData.Member.EnName+" Live in "+LiveCount).
							SetDescription(PushData.YtData.Title).
							SetImage(PushData.YtData.Thumb).
							SetThumbnail(PushData.Group.IconURL).
							SetURL(YtURL).
							AddField("Type ", PushData.YtData.Type).
							AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(2).String()).
							InlineAllFields().
							AddField("Waiting", PushData.YtData.Viewers+" Simps in ChatRoom").
							SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.WithFields(log.Fields{
								"Message":          MsgEmbed,
								"ChannelID":        Channel.ID,
								"DiscordChannelID": Channel.ChannelID,
							}).Error(err)
							err = Channel.DelChannel(err.Error())
							return err
						}
						MsgText, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` Live in "+LiveCount+"\nUserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							return err
						}
						if Channel.Dynamic {
							log.WithFields(log.Fields{
								"DiscordChannel": Channel.ChannelID,
								"VtuberGroupID":  PushData.Group.ID,
								"YoutubeID":      PushData.YtData.ID,
							}).Info("Set dynamic mode")
							Channel.SetVideoID(PushData.YtData.VideoID).
								SetMsgEmbedID(MsgEmbed.ID).
								SetMsgTextID(MsgText.ID)
						}
					}
				}
			}
		}
	}
	return nil
}
