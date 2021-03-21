package youtube

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/livestream/bilibili/live"

	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func SendNude(PushData database.LiveStream) error {
	Status := PushData.Status
	Avatar := PushData.Member.YoutubeAvatar
	YtChannel := "https://www.youtube.com/channel/" + PushData.Member.YoutubeID + "?sub_confirmation=1"
	YtURL := "https://www.youtube.com/watch?v=" + PushData.VideoID
	loc := engine.Zawarudo(PushData.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(PushData.Member.EnName, PushData.Member.JpName)

	FanBase := "simps"
	if PushData.Member.Fanbase != "" {
		FanBase = PushData.Member.Fanbase
	}

	var (
		Timestart time.Time
		User      = &database.UserStruct{
			Human:    true,
			Reminder: 0,
		}
		Viewers string
	)

	if match, _ := regexp.MatchString("404.jpg", PushData.Group.IconURL); match {
		PushData.Group.IconURL = ""
	}

	if !PushData.Schedul.IsZero() {
		Timestart = PushData.Schedul
	} else if PushData.Schedul.IsZero() && !PushData.Published.IsZero() {
		Timestart = PushData.Published
	} else if PushData.Schedul.IsZero() && PushData.Published.IsZero() {
		Timestart = time.Now()
	}

	if PushData.Viewers == "0" {
		PushData.Viewers = Ytwaiting
	} else {
		view, err := strconv.Atoi(PushData.Viewers)
		if err != nil {
			log.Error(err)
		}
		Viewers = engine.NearestThousandFormat(float64(view))
	}

	if Status == "upcoming" {
		Color, err := engine.GetColor(config.TmpDir, PushData.Thumb)
		if err != nil {
			return err
		}
		//id, DiscordChannelID
		var (
			wg sync.WaitGroup
		)

		ChannelData, err := database.ChannelTag(PushData.Member.ID, 2, config.NewUpcoming, PushData.Member.Region)
		if err != nil {
			log.Error(err)
		}

		for i, v := range ChannelData {
			v.SetMember(PushData.Member)

			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				ctx := context.Background()
				UserTagsList, err := Channel.GetUserList(ctx) //database.GetUserList(Channel.ID, PushData.Member.ID)
				if err != nil {
					log.Error(err)
				}
				if UserTagsList == nil && PushData.Group.GroupName != config.Indie {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && PushData.Group.GroupName == config.Indie && !Channel.IndieNotif {
					return nil
				}
				msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Avatar, YtChannel).
					SetTitle("New upcoming Livestream").
					SetDescription(PushData.Title).
					SetImage(PushData.Thumb).
					SetThumbnail(PushData.Group.IconURL).
					SetURL(YtURL).
					AddField("Type ", PushData.Type).
					AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()).
					InlineAllFields().
					AddField("Waiting", Viewers+" "+FanBase+" in ChatRoom").
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
				if !Channel.LiteMode {
					msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` New upcoming Livestream\nUserTags: "+strings.Join(UserTagsList, " "))
					if err != nil {
						return err
					}
				}
				return nil
			}(v, &wg)
			//Wait every ge 10 discord channel
			if i%config.Waiting == 0 && configfile.LowResources {
				log.WithFields(log.Fields{
					"Func":  "Youtube",
					"Value": config.Waiting,
				}).Warn("Waiting send message")
				wg.Wait()
				expiresAt = time.Now().In(loc)
			}
		}
		wg.Wait()

	} else if Status == "live" {
		Color, err := engine.GetColor(config.TmpDir, PushData.Thumb)
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
			wg sync.WaitGroup
		)

		ChannelData, err := database.ChannelTag(PushData.Member.ID, 2, config.Default, PushData.Member.Region)
		if err != nil {
			log.Error(err)
		}

		for i, v := range ChannelData {
			v.SetMember(PushData.Member)

			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				ctx := context.Background()
				UserTagsList, err := Channel.GetUserList(ctx)
				if err != nil {
					log.Error(err)
				}

				if UserTagsList == nil && PushData.Group.GroupName != config.Indie {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && PushData.Group.GroupName == config.Indie && !Channel.IndieNotif {
					return nil
				}

				MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Avatar, YtChannel).
					SetTitle("Live right now").
					SetDescription(PushData.Title).
					SetImage(PushData.Thumb).
					SetThumbnail(PushData.Group.IconURL).
					SetURL(YtURL).
					AddField("Type ", PushData.Type).
					AddField("Start live", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(1).String()+" Ago").
					InlineAllFields().
					AddField("Viewers", Viewers+" "+FanBase).
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
					return err
				}

				if Channel.Dynamic {
					log.WithFields(log.Fields{
						"DiscordChannel":  Channel.ChannelID,
						"VtuberGroupName": PushData.Group.GroupName,
						"YoutubeVideoID":  PushData.VideoID,
					}).Info("Set dynamic mode")
					Channel.SetVideoID(PushData.VideoID).
						SetMsgEmbedID(MsgEmbed.ID)
				}

				if !Channel.LiteMode {
					Msg := "Push " + configfile.Emoji.Livestream[0] + " to add you in `" + PushData.Member.Name + "` ping list\nPush " + configfile.Emoji.Livestream[1] + " to remove you from ping list"
					MsgFinal := ""

					if Bili {
						MsgFinal = "`" + PushData.Member.Name + "` Live right now at BiliBili And Youtube\nUserTags: " + strings.Join(UserTagsList, " ") + "\n" + Msg
					} else {
						MsgFinal = "`" + PushData.Member.Name + "` Live right now\nUserTags: " + strings.Join(UserTagsList, " ") + "\n" + Msg
					}

					msgText, err := Bot.ChannelMessageSend(Channel.ChannelID, MsgFinal)
					if err != nil {
						return err
					}

					User.SetDiscordChannelID(Channel.ChannelID).
						SetGroup(PushData.Group).
						SetMember(PushData.Member)
					err = User.SendToCache(msgText.ID)
					if err != nil {
						return err
					}

					Channel.SetMsgTextID(msgText.ID)

					err = engine.Reacting(map[string]string{
						"ChannelID": Channel.ChannelID,
						"State":     "Youtube",
						"MessageID": msgText.ID,
					}, Bot)
					if err != nil {
						return err
					}
				}

				Channel.PushReddis()

				return nil
			}(v, &wg)
			//Wait every ge 5 discord channel
			if i%config.Waiting == 0 && configfile.LowResources {
				log.WithFields(log.Fields{
					"Func":  "Youtube",
					"Value": config.Waiting,
				}).Warn("Waiting send message")
				wg.Wait()
				expiresAt = time.Now().In(loc)
			}
		}
		wg.Wait()

	} else if Status == "past" {
		Color, err := engine.GetColor(config.TmpDir, PushData.Thumb)
		if err != nil {
			return err
		}

		//id, DiscordChannelID
		var (
			wg sync.WaitGroup
		)

		ChannelData, err := database.ChannelTag(PushData.Member.ID, 2, config.NotLiveOnly, PushData.Member.Region)
		if err != nil {
			log.Error(err)
		}

		for i, v := range ChannelData {
			v.SetMember(PushData.Member)

			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				ctx := context.Background()
				UserTagsList, err := Channel.GetUserList(ctx)
				if err != nil {
					log.Error(err)
				}

				if UserTagsList == nil && PushData.Group.GroupName != config.Indie {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && PushData.Group.GroupName == config.Indie && !Channel.IndieNotif {
					return nil
				}

				msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Avatar, YtChannel).
					SetTitle("Uploaded a new video").
					SetDescription(PushData.Title).
					SetImage(PushData.Thumb).
					SetThumbnail(PushData.Group.IconURL).
					SetURL(YtURL).
					AddField("Type ", PushData.Type).
					AddField("Upload", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(1).String()+" Ago").
					AddField("Viewers", Viewers+" "+FanBase).
					AddField("Duration", PushData.Length).
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
				if !Channel.LiteMode {
					msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+PushData.Member.Name+"` Uploaded a new video\nUserTags: "+strings.Join(UserTagsList, " "))
					if err != nil {
						return err
					}
				}

				return nil
			}(v, &wg)
			//Wait every ge 5 discord channel
			if i%config.Waiting == 0 && configfile.LowResources {
				log.WithFields(log.Fields{
					"Func":  "Youtube",
					"Value": config.Waiting,
				}).Warn("Waiting send message")
				wg.Wait()
				expiresAt = time.Now().In(loc)
			}
		}
	} else if Status == "reminder" {
		UpcominginMinutes := int(Timestart.Sub(time.Now()).Minutes())
		if UpcominginMinutes > 10 && UpcominginMinutes < 70 {
			if database.CheckReminder(UpcominginMinutes) {
				ChanelData, err := database.ChannelTag(PushData.Member.ID, 2, config.Default, PushData.Member.Region)
				if err != nil {
					log.Error(err)
				}

				Color, err := engine.GetColor(config.TmpDir, PushData.Thumb)
				if err != nil {
					return err
				}
				LiveCount := durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()
				for _, Channel := range ChanelData {
					UserTagsList := database.GetUserReminderList(Channel.ID, PushData.Member.ID, UpcominginMinutes)
					if UserTagsList != nil {
						MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
							SetAuthor(VtuberName, Avatar, YtChannel).
							SetTitle(PushData.Member.EnName+" Live in "+LiveCount).
							SetDescription(PushData.Title).
							SetImage(PushData.Thumb).
							SetThumbnail(PushData.Group.IconURL).
							SetURL(YtURL).
							AddField("Type", PushData.Type).
							AddField("Start live in", LiveCount).
							InlineAllFields().
							AddField("Waiting", Viewers+" "+FanBase+" in ChatRoom").
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
					} else {
						break
					}
				}
			}
		}
	}
	return nil
}
