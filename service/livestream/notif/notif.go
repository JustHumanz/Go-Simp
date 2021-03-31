package notif

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
)

func SendDude(Data *database.LiveStream, Bot *discordgo.Session) {
	loc := engine.Zawarudo(Data.Member.Region)
	expiresAt := time.Now().In(loc)
	VtuberName := engine.FixName(Data.Member.EnName, Data.Member.JpName)
	FanBase := "simps"
	if Data.Member.Fanbase != "" {
		FanBase = Data.Member.Fanbase
	}

	Data.Group.RemoveNillIconURL()

	Color := func() int {
		clr, err := engine.GetColor(config.TmpDir, Data.Thumb)
		if err != nil {
			log.Error(err)
		}
		return clr
	}
	if Data.State == config.YoutubeLive {
		Status := Data.Status
		Avatar := Data.Member.YoutubeAvatar
		YtChannel := "https://www.youtube.com/channel/" + Data.Member.YoutubeID + "?sub_confirmation=1"
		YtURL := "https://www.youtube.com/watch?v=" + Data.VideoID

		var (
			Timestart time.Time
			User      = &database.UserStruct{
				Human:    true,
				Reminder: 0,
			}
			Viewers string
		)

		if !Data.Schedul.IsZero() {
			Timestart = Data.Schedul
		} else if Data.Schedul.IsZero() && !Data.Published.IsZero() {
			Timestart = Data.Published
		} else if Data.Schedul.IsZero() && Data.Published.IsZero() {
			Timestart = time.Now()
		}

		if Data.Viewers == "0" {
			Data.Viewers = config.Ytwaiting
		} else {
			view, err := strconv.Atoi(Data.Viewers)
			if err != nil {
				log.Error(err)
				Viewers = config.Ytwaiting
			} else {
				Viewers = engine.NearestThousandFormat(float64(view))
			}
		}

		if Status == config.UpcomingStatus {
			var (
				wg sync.WaitGroup
			)

			ChannelData, err := database.ChannelTag(Data.Member.ID, 2, config.NewUpcoming, Data.Member.Region)
			if err != nil {
				log.Error(err)
			}

			for i, v := range ChannelData {
				v.SetMember(Data.Member)

				wg.Add(1)
				go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
					defer wg.Done()
					ctx := context.Background()
					UserTagsList, err := Channel.GetUserList(ctx) //database.GetUserList(Channel.ID, Data.Member.ID)
					if err != nil {
						log.Error(err)
					}
					if UserTagsList == nil && Data.Group.GroupName != config.Indie {
						UserTagsList = []string{"_"}
					} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
						return nil
					}
					msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
						SetAuthor(VtuberName, Avatar, YtChannel).
						SetTitle("New upcoming Livestream").
						SetDescription(Data.Title).
						SetImage(Data.Thumb).
						SetThumbnail(Data.Group.IconURL).
						SetURL(YtURL).
						AddField("Type ", Data.Type).
						AddField("Start live in", durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()).
						InlineAllFields().
						AddField("Waiting", Viewers+" "+FanBase+" in ChatRoom").
						SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
						SetColor(Color()).MessageEmbed)
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
						msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` New upcoming Livestream\nUserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							return err
						}
					}
					return nil
				}(v, &wg)
				//Wait every ge 10 discord channel
				if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
					log.WithFields(log.Fields{
						"Func":  "Youtube",
						"Value": config.Waiting,
					}).Warn("Waiting send message")
					wg.Wait()
					expiresAt = time.Now().In(loc)
				}
			}
			wg.Wait()

		} else if Status == config.LiveStatus {
			var (
				wg sync.WaitGroup
			)

			ChannelData, err := database.ChannelTag(Data.Member.ID, 2, config.Default, Data.Member.Region)
			if err != nil {
				log.Error(err)
			}

			for i, v := range ChannelData {
				v.SetMember(Data.Member)

				wg.Add(1)
				go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
					defer wg.Done()
					ctx := context.Background()
					UserTagsList, err := Channel.GetUserList(ctx)
					if err != nil {
						log.Error(err)
					}

					if UserTagsList == nil && Data.Group.GroupName != config.Indie {
						UserTagsList = []string{"_"}
					} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
						return nil
					}

					MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
						SetAuthor(VtuberName, Avatar, YtChannel).
						SetTitle("Live right now").
						SetDescription(Data.Title).
						SetImage(Data.Thumb).
						SetThumbnail(Data.Group.IconURL).
						SetURL(YtURL).
						AddField("Type ", Data.Type).
						AddField("Start live", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(1).String()+" Ago").
						InlineAllFields().
						AddField("Viewers", Viewers+" "+FanBase).
						SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
						SetColor(Color()).MessageEmbed)
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
							"VtuberGroupName": Data.Group.GroupName,
							"YoutubeVideoID":  Data.VideoID,
						}).Info("Set dynamic mode")
						Channel.SetVideoID(Data.VideoID).
							SetMsgEmbedID(MsgEmbed.ID)
					}

					if !Channel.LiteMode {
						Msg := "Push " + config.GoSimpConf.Emoji.Livestream[0] + " to add you in `" + Data.Member.Name + "` ping list\nPush " + config.GoSimpConf.Emoji.Livestream[1] + " to remove you from ping list"
						MsgFinal := ""

						if Data.IsBiliLive {
							MsgFinal = "`" + Data.Member.Name + "` Live right now at BiliBili And Youtube\nUserTags: " + strings.Join(UserTagsList, " ") + "\n" + Msg
						} else {
							MsgFinal = "`" + Data.Member.Name + "` Live right now\nUserTags: " + strings.Join(UserTagsList, " ") + "\n" + Msg
						}

						msgText, err := Bot.ChannelMessageSend(Channel.ChannelID, MsgFinal)
						if err != nil {
							return err
						}

						User.SetDiscordChannelID(Channel.ChannelID).
							SetGroup(Data.Group).
							SetMember(Data.Member)
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
				if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
					log.WithFields(log.Fields{
						"Func":  "Youtube",
						"Value": config.Waiting,
					}).Warn("Waiting send message")
					wg.Wait()
					expiresAt = time.Now().In(loc)
				}
			}
			wg.Wait()

		} else if Status == config.PastStatus {
			//id, DiscordChannelID
			var (
				wg sync.WaitGroup
			)

			ChannelData, err := database.ChannelTag(Data.Member.ID, 2, config.NotLiveOnly, Data.Member.Region)
			if err != nil {
				log.Error(err)
			}

			for i, v := range ChannelData {
				v.SetMember(Data.Member)

				wg.Add(1)
				go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
					defer wg.Done()
					ctx := context.Background()
					UserTagsList, err := Channel.GetUserList(ctx)
					if err != nil {
						log.Error(err)
					}

					if UserTagsList == nil && Data.Group.GroupName != config.Indie {
						UserTagsList = []string{"_"}
					} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
						return nil
					}

					msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
						SetAuthor(VtuberName, Avatar, YtChannel).
						SetTitle("Uploaded a new video").
						SetDescription(Data.Title).
						SetImage(Data.Thumb).
						SetThumbnail(Data.Group.IconURL).
						SetURL(YtURL).
						AddField("Type ", Data.Type).
						AddField("Upload", durafmt.Parse(expiresAt.Sub(Timestart.In(loc))).LimitFirstN(1).String()+" Ago").
						AddField("Viewers", Viewers+" "+FanBase).
						AddField("Duration", Data.Length).
						InlineAllFields().
						SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
						SetColor(Color()).MessageEmbed)
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
						msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` Uploaded a new video\nUserTags: "+strings.Join(UserTagsList, " "))
						if err != nil {
							return err
						}
					}

					return nil
				}(v, &wg)
				//Wait every ge 5 discord channel
				if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
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
					ChanelData, err := database.ChannelTag(Data.Member.ID, 2, config.Default, Data.Member.Region)
					if err != nil {
						log.Error(err)
					}
					LiveCount := durafmt.Parse(Timestart.In(loc).Sub(expiresAt)).LimitFirstN(1).String()
					for _, Channel := range ChanelData {
						UserTagsList := database.GetUserReminderList(Channel.ID, Data.Member.ID, UpcominginMinutes)
						if UserTagsList != nil {
							MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
								SetAuthor(VtuberName, Avatar, YtChannel).
								SetTitle(Data.Member.EnName+" Live in "+LiveCount).
								SetDescription(Data.Title).
								SetImage(Data.Thumb).
								SetThumbnail(Data.Group.IconURL).
								SetURL(YtURL).
								AddField("Type", Data.Type).
								AddField("Start live in", LiveCount).
								InlineAllFields().
								AddField("Waiting", Viewers+" "+FanBase+" in ChatRoom").
								SetFooter(Timestart.In(loc).Format(time.RFC822), config.YoutubeIMG).
								SetColor(Color()).MessageEmbed)
							if err != nil {
								log.WithFields(log.Fields{
									"Message":          MsgEmbed,
									"ChannelID":        Channel.ID,
									"DiscordChannelID": Channel.ChannelID,
								}).Error(err)
								err = Channel.DelChannel(err.Error())
								log.Error(err)
							}
						} else {
							break
						}
					}
				}
			}
		}
	} else if Data.State == config.BiliLive {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		BiliBiliAccount := "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
		BiliBiliURL := "https://live.bilibili.com/" + strconv.Itoa(Data.Member.BiliRoomID)

		BiliBiliRoomID := strconv.Itoa(Data.Member.BiliRoomID)
		User := &database.UserStruct{
			Human:    true,
			Reminder: 0,
		}
		if Data.Status == config.LiveStatus {

			MemberID := Data.Member.ID
			//id, DiscordChannelID
			var (
				wg sync.WaitGroup
			)

			ChannelData, err := database.ChannelTag(MemberID, 2, config.Default, Data.Member.Region)
			if err != nil {
				log.Error(err)
			}
			for i, v := range ChannelData {
				v.SetMember(Data.Member)

				wg.Add(1)
				go func(Channel database.DiscordChannel, wg *sync.WaitGroup) {
					defer wg.Done()
					ctx := context.Background()
					UserTagsList, err := Channel.GetUserList(ctx)
					if err != nil {
						log.Error(err)
					}
					if UserTagsList == nil && Data.Group.GroupName != config.Indie {
						UserTagsList = []string{"_"}
					} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
						return
					}
					view, err := strconv.Atoi(Data.Viewers)
					if err != nil {
						log.Error(err)
					}

					Start := durafmt.Parse(expiresAt.Sub(Data.Schedul.In(loc))).LimitFirstN(1)
					MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
						SetAuthor(VtuberName, Data.Member.BiliBiliAvatar, BiliBiliAccount).
						SetTitle(Data.Title).
						SetThumbnail(Data.Group.IconURL).
						SetDescription(Data.Desc).
						SetImage(Data.Thumb).
						SetURL(BiliBiliURL).
						AddField("Start live", Start.String()+" Ago").
						AddField("Online", engine.NearestThousandFormat(float64(view))+" "+FanBase).
						InlineAllFields().
						SetFooter(Data.Schedul.In(loc).Format(time.RFC822), config.BiliBiliIMG).
						SetColor(Color()).MessageEmbed)
					if err != nil {
						log.Error(err)
					}

					if Channel.Dynamic {
						log.WithFields(log.Fields{
							"DiscordChannel": Channel.ChannelID,
							"VtuberGroupID":  Data.Group.ID,
							"BiliBiliRoomID": BiliBiliRoomID,
						}).Info("Set dynamic mode")
						Channel.SetVideoID(BiliBiliRoomID).
							SetMsgEmbedID(MsgEmbed.ID)
					}

					if !Channel.LiteMode {
						Msg := "Push " + config.GoSimpConf.Emoji.Livestream[0] + " to add you in `" + Data.Member.Name + "` ping list\nPush " + config.GoSimpConf.Emoji.Livestream[1] + " to remove you from ping list"
						MsgTxt, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` Live right now\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
						if err != nil {
							log.Error(err)
							return
						}
						User.SetDiscordChannelID(Channel.ChannelID).
							SetGroup(Data.Group).
							SetMember(Data.Member)
						err = User.SendToCache(MsgTxt.ID)
						if err != nil {
							log.Error(err)
						}

						Channel.SetMsgTextID(MsgTxt.ID)
						err = engine.Reacting(map[string]string{
							"ChannelID": Channel.ChannelID,
							"State":     "Youtube",
							"MessageID": MsgTxt.ID,
						}, Bot)
						if err != nil {
							log.Error(err)
						}
					}

					Channel.PushReddis()

				}(v, &wg)
				//Wait every ge 5 discord channel
				if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
					log.WithFields(log.Fields{
						"Func":  "BiliBili Live",
						"Value": config.Waiting,
					}).Warn("Waiting send message")
					wg.Wait()
					expiresAt = time.Now().In(loc)
				}
			}
			wg.Wait()
		} else {
			log.Warn("it's not live")
		}
	} else if Data.State == config.TwitchLive {

		var (
			wg        sync.WaitGroup
			ImgURL    = "https://www.twitch.tv/" + Data.Member.TwitchName
			expiresAt = time.Now().In(loc)
			User      = &database.UserStruct{
				Human:    true,
				Reminder: 0,
			}
		)

		ChannelData, err := database.ChannelTag(Data.Member.ID, 2, config.Default, Data.Member.Region)
		if err != nil {
			log.Error(err)
		}

		for i, v := range ChannelData {
			v.SetMember(Data.Member)

			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
				defer wg.Done()
				ctx := context.Background()
				UserTagsList, err := Channel.GetUserList(ctx)
				if err != nil {
					log.Error(err)
				}
				if UserTagsList == nil && Data.Group.GroupName != config.Indie {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
					return nil
				}

				View, err := strconv.Atoi(Data.Viewers)
				if err != nil {
					log.Error()
				}
				MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Data.Member.YoutubeAvatar, ImgURL).
					SetTitle("Live right now").
					SetDescription(Data.Title).
					SetImage(Data.Thumb).
					SetThumbnail(Data.Group.IconURL).
					SetURL(ImgURL).
					AddField("Start live", durafmt.Parse(expiresAt.Sub(Data.Schedul.In(loc))).LimitFirstN(1).String()+" Ago").
					AddField("Viewers", engine.NearestThousandFormat(float64(View))+" "+FanBase).
					InlineAllFields().
					AddField("Game", Data.Game).
					SetFooter(Data.Schedul.In(loc).Format(time.RFC822), config.TwitchIMG).
					SetColor(Color()).MessageEmbed)
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
						"DiscordChannel": Channel.ChannelID,
						"VtuberGroupID":  Data.Group.ID,
						"TwitchID":       "Twitch" + Data.Member.TwitchName,
					}).Info("Set dynamic mode")
					Channel.SetVideoID("Twitch" + Data.Member.TwitchName).
						SetMsgEmbedID(MsgEmbed.ID)
				}

				if !Channel.LiteMode {
					Msg := "Push " + config.GoSimpConf.Emoji.Livestream[0] + " to add you in `" + Data.Member.Name + "` ping list\nPush " + config.GoSimpConf.Emoji.Livestream[1] + " to remove you from ping list"
					msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` Live right now\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
					if err != nil {
						return err
					}
					User.SetDiscordChannelID(Channel.ChannelID).
						SetGroup(Data.Group).
						SetMember(Data.Member)

					err = User.SendToCache(msg.ID)
					if err != nil {
						return err
					}

					Channel.SetMsgTextID(msg.ID)

					err = engine.Reacting(map[string]string{
						"ChannelID": Channel.ChannelID,
						"State":     "Youtube",
						"MessageID": msg.ID,
					}, Bot)
					if err != nil {
						return err
					}
				}

				Channel.PushReddis()

				return nil
			}(v, &wg)
			//Wait every ge 5 discord channel
			if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
				log.WithFields(log.Fields{
					"Func":  "Twitch",
					"Value": config.Waiting,
				}).Warn("Waiting send message")
				wg.Wait()
				expiresAt = time.Now().In(loc)
			}
		}
		wg.Wait()
	} else {
		var (
			wg sync.WaitGroup
		)

		ChannelData, err := database.ChannelTag(Data.Member.ID, 2, config.NotLiveOnly, Data.Member.Region)
		if err != nil {
			log.Error(err)
		}
		for i, v := range ChannelData {
			v.SetMember(Data.Member)

			wg.Add(1)
			go func(Channel database.DiscordChannel, wg *sync.WaitGroup) {
				defer wg.Done()
				ctx := context.Background()
				UserTagsList, err := Channel.GetUserList(ctx)
				if err != nil {
					log.Error(err)
				}

				if UserTagsList == nil && Data.Group.GroupName != config.Indie {
					UserTagsList = []string{"_"}
				} else if UserTagsList == nil && Data.Group.GroupName == config.Indie && !Channel.IndieNotif {
					return
				}
				viwe, err := strconv.Atoi(Data.Viewers)
				if err != nil {
					log.Error(err)
				}

				msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(VtuberName, Data.Member.BiliBiliAvatar, "https://space.bilibili/"+strconv.Itoa(Data.Member.BiliBiliID)).
					SetTitle("Uploaded new video").
					SetDescription(Data.Title).
					SetImage(Data.Thumb).
					SetThumbnail(Data.Group.IconURL).
					SetURL("https://www.bilibili.com/video/"+Data.VideoID).
					AddField("Type ", Data.Type).
					AddField("Duration ", Data.Length).
					InlineAllFields().
					AddField("Viwers ", engine.NearestThousandFormat(float64(viwe))+" "+FanBase).
					SetFooter(Data.Schedul.In(loc).Format(time.RFC822), config.BiliBiliIMG).
					SetColor(Color()).MessageEmbed)
				if err != nil {
					log.Error(msg, err)
				} else {
					msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "UserTags: "+strings.Join(UserTagsList, " "))
					if err != nil {
						log.Error(msg, err)
					}
				}
			}(v, &wg)
			//Wait every ge 5 discord channel
			if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
				log.WithFields(log.Fields{
					"Func":  "BiliBili space",
					"Value": config.Waiting,
				}).Warn("Waiting send message")
				wg.Wait()
				expiresAt = time.Now().In(loc)
			}
		}
		wg.Wait()

	}
}
