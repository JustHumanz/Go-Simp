package main

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"sync"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	_ "github.com/go-sql-driver/mysql"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

func AddData(Data Vtuber) {
	GetYtChannelInfo := func(Member Members) Items {
		YtItem := Items{}
		YoutubeChannelInfoRaw, err := Member.Youtube.GetYtInfo()
		if err != nil {
			log.Warn(err)
		}

		if len(YoutubeChannelInfoRaw.Items) > 0 {
			for _, Item := range YoutubeChannelInfoRaw.Items {
				YtItem = Item
				YtItem.Brandingsettings.Image.Bannerexternalurl += "=s1200"
			}
		}
		return YtItem
	}

	GetTwitterAccountInfo := func(Member Members) twitterscraper.Profile {
		profile, err := Member.Twitter.GetTwitterInfo()
		if err != nil {
			log.Error(err)
		}
		return profile
	}

	GetBiliBiliAccountInfo := func(Member Members) Avatar {
		BiliBili, err := Member.BiliBili.GetBiliBiliInfo()
		if err != nil {
			log.Error(err)
		}
		return BiliBili
	}

	Independen := func(wg *sync.WaitGroup) {
		var (
			GroupData = database.Group{
				ID:        0,
				GroupName: config.Indie,
				IconURL:   "https://cdn." + configfile.Domain + "/404.jpg",
			}
			NewVtuberNamesIndependen []string
			User                     = &database.UserStruct{
				Human:    true,
				Reminder: 0,
			}
		)
		defer wg.Done()
		row := db.QueryRow("SELECT id FROM VtuberGroup WHERE VtuberGroupName=?", GroupData.GroupName)
		err := row.Scan(&GroupData.ID)
		if err == sql.ErrNoRows {
			log.Error(err)
			stmt, err := db.Prepare("INSERT INTO VtuberGroup (VtuberGroupName,VtuberGroupIcon) values(?,?)")
			if err != nil {
				log.Error(err)
			}

			res, err := stmt.Exec(GroupData.GroupName, GroupData.IconURL)
			if err != nil {
				log.Error(err)
			}

			GroupData.ID, err = res.LastInsertId()
			if err != nil {
				log.Error(err)
			}

			defer stmt.Close()
		} else {
			log.WithFields(log.Fields{
				"VtuberGroup": GroupData.GroupName,
			}).Info("Update Vtuber Group Data")
			Update, err := db.Prepare(`Update VtuberGroup set VtuberGroupName=?, VtuberGroupIcon=? Where id=?`)
			if err != nil {
				log.Error(err)
			}
			Update.Exec(GroupData.GroupName, GroupData.IconURL, GroupData.ID)
		}

		for _, VtuberMember := range Data.VtuberData.Independent.Members {
			DiscordChannel, err := GroupData.GetChannelByGroup(VtuberMember.Region)
			if err != nil {
				log.Error(err)
				continue
			}
			/*
				Add Member
			*/
			YtItem := GetYtChannelInfo(VtuberMember)
			TwitterItem := GetTwitterAccountInfo(VtuberMember)
			BiliItem := GetBiliBiliAccountInfo(VtuberMember)

			Region := func() string {
				if YtItem.Brandingsettings.Channel.Country != "" {
					return YtItem.Brandingsettings.Channel.Country
				} else {
					return VtuberMember.Region
				}
			}()

			var MemberID int64
			row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND VtuberName_EN=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.EnName, VtuberMember.Youtube.YtID, VtuberMember.BiliBili.BiliBiliID, VtuberMember.BiliBili.BiliRoomID)
			err = row.Scan(&MemberID)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName, VtuberName_EN, VtuberName_JP, Twitter_Username, Twitter_Hashtag, Twitter_Lewd, Twitter_Avatar, Twitter_Banner, Youtube_ID, Youtube_Avatar, Youtube_Banner, BiliBili_SpaceID, BiliBili_RoomID, BiliBili_Avatar, BiliBili_Hashtag, BiliBili_Banner, Twitch_Username, Twitch_Avatar, Region, Fanbase, Status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
				if err != nil {
					log.Error(err)
				}

				TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}

				res, err := stmt.Exec(
					VtuberMember.Name,
					VtuberMember.EnName,
					VtuberMember.JpName,
					VtuberMember.Twitter.TwitterUsername,
					VtuberMember.Twitter.TwitterFanart,
					VtuberMember.Twitter.TwitterLewd,
					TwitterItem.Avatar,
					TwitterItem.Banner,
					VtuberMember.Youtube.YtID,
					YtItem.Brandingsettings.Image.Bannerexternalurl,
					VtuberMember.BiliBili.BiliBiliID,
					VtuberMember.BiliBili.BiliRoomID,
					BiliItem.Data.Face,
					VtuberMember.BiliBili.BiliBiliFanart,
					BiliItem.Data.TopPhoto,
					VtuberMember.Twitch.TwitchUsername,
					TwitchAvatar,
					Region,
				)
				if err != nil {
					log.Error(err)
				}

				MemberID, err = res.LastInsertId()
				if err != nil {
					log.Error(err)
				}

				defer stmt.Close()

				for _, Channel := range DiscordChannel {
					log.WithFields(log.Fields{
						"ChannelID": Channel.ChannelID,
					}).Info("Send notif")

					msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewVtuber{
						Group:      GroupData,
						Member:     VtuberMember,
						YtAvatar:   YtItem.Snippet.Thumbnails.High.URL,
						BiliAvatar: BiliItem.Data.Face,
					}.SendNotif())
					if err != nil {
						log.Error(msg, err)
						match, _ := regexp.MatchString("Unknown Channel", err.Error())
						if match {
							log.Info("Delete Discord Channel ", Channel)
							DeleteChannel(Channel.ID)
						}
					}

					if msg != nil {
						User.SetDiscordChannelID(Channel.ChannelID).
							SetGroup(GroupData).
							SetMember(database.Member{
								ID:     MemberID,
								Name:   VtuberMember.Name,
								EnName: VtuberMember.EnName,
								JpName: VtuberMember.JpName,
								Group: database.Group{
									ID: GroupData.ID,
								},
							})
						err = User.SendToCache(msg.ID)
						if err != nil {
							log.Error(err)
						}

						err = engine.Reacting(map[string]string{
							"ChannelID": Channel.ChannelID,
							"State":     "New Member",
							"MessageID": msg.ID,
						}, Bot)
						if err != nil {
							log.Error(err)
						}
					}
				}

				NewVtuberNamesIndependen = append(NewVtuberNamesIndependen, "`"+VtuberMember.Name+"`")
				VtuberMember.InputSubs(MemberID, YtItem, TwitterItem)
				//New.SendNotif(Bot)
			} else if err != nil {
				log.Error(err)
			} else {
				TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}

				log.WithFields(log.Fields{
					"VtuberGroup": "Independen",
					"Vtuber":      VtuberMember.EnName,
				}).Info("Update member")
				_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=?, Twitter_Username=?, Twitter_Hashtag=?, Twitter_Lewd=?, Twitter_Avatar=?, Twitter_Banner=?, Youtube_ID=?, Youtube_Avatar=?, Youtube_Banner=?, BiliBili_SpaceID=?, BiliBili_RoomID=?, BiliBili_Avatar=?, BiliBili_Hashtag=?, BiliBili_Banner=?, Twitch_Username=?, Twitch_Avatar=?, Region=?, Fanbase=?, Status=?  Where id=?`,
					VtuberMember.Name,
					VtuberMember.EnName,
					VtuberMember.JpName,
					VtuberMember.Twitter.TwitterUsername,
					VtuberMember.Twitter.TwitterFanart,
					VtuberMember.Twitter.TwitterLewd,
					TwitterItem.Avatar,
					TwitterItem.Banner,
					VtuberMember.Youtube.YtID,
					YtItem.Brandingsettings.Image.Bannerexternalurl,
					VtuberMember.BiliBili.BiliBiliID,
					VtuberMember.BiliBili.BiliRoomID,
					BiliItem.Data.Face,
					VtuberMember.BiliBili.BiliBiliFanart,
					BiliItem.Data.TopPhoto,
					VtuberMember.Twitch.TwitchUsername,
					TwitchAvatar,
					Region,
					VtuberMember.Fanbase,
					VtuberMember.Status,
					MemberID,
				)
				if err != nil {
					log.Error(err)
				}

			}
			log.WithFields(log.Fields{
				"VtuberGroup": "Independen",
				"Vtuber":      VtuberMember.EnName,
			}).Info("Add subs info to database")

			/*
				Add subs info
			*/
			VtuberMember.InputSubs(MemberID, YtItem, TwitterItem)
		}

		if NewVtuberNamesIndependen != nil {
			Vtubers := strings.Join(NewVtuberNamesIndependen, ",")
			DiscordChannel, err := GroupData.GetChannelByGroup("")
			if err != nil {
				log.Error(err)
			}
			for _, Channel := range DiscordChannel {
				msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "New Update!!!! "+Vtubers)
				if err != nil {
					log.Error(msg, err)
					match, _ := regexp.MatchString("Unknown Channel", err.Error())
					if match {
						log.Info("Delete Discord Channel ", Channel.ChannelID)
						DeleteChannel(Channel.ID)
					}
				}
				_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Push "+configfile.Emoji.Livestream[0]+" to add you in "+Vtubers+" ping list")
				if err != nil {
					log.Error(err)
				}

				_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Update your roles with `"+configfile.BotPrefix.General+"tag roles @somesimpsroles` "+Vtubers)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

	Agency := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for _, GroupRaw := range Data.VtuberData.Group {
			var (
				NewVtuberNames []string
				GroupData      = database.Group{
					ID:        0,
					GroupName: GroupRaw.GroupName,
					IconURL:   GroupRaw.GroupIcon,
				}
				User = &database.UserStruct{
					Human:    true,
					Reminder: 0,
				}
			)
			/*
				Add Group
			*/
			log.WithFields(log.Fields{
				"VtuberGroup":     GroupData.GroupName,
				"VtuberGroupIcon": GroupData.IconURL,
			}).Info("Add Group")

			row := db.QueryRow("SELECT id FROM VtuberGroup WHERE VtuberGroupName=? ", GroupData.GroupName)
			err := row.Scan(&GroupData.ID)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberGroup (VtuberGroupName,VtuberGroupIcon) values(?,?)")
				if err != nil {
					log.Error(err)
				}

				res, err := stmt.Exec(GroupData.GroupName, GroupData.IconURL)
				if err != nil {
					log.Error(err)
				}

				GroupData.ID, err = res.LastInsertId()
				if err != nil {
					log.Error(err)
				}

				defer stmt.Close()
			} else if err != nil {
				log.Error(err)
			} else {
				log.WithFields(log.Fields{
					"VtuberGroup":     GroupData.GroupName,
					"VtuberGroupIcon": GroupData.IconURL,
				}).Info("Update Vtuber Group Data")
				Update, err := db.Prepare(`Update VtuberGroup set VtuberGroupName=?, VtuberGroupIcon=? Where id=?`)
				if err != nil {
					log.Error(err)
				}
				Update.Exec(GroupData.GroupName, GroupData.IconURL, GroupData.ID)
			}

			for _, ytb := range GroupRaw.GroupChannel.Youtube {
				if ytb != nil {
					tmp := 0
					Youtube := ytb.(map[string]interface{})
					row := db.QueryRow("SELECT id FROM GroupYoutube WHERE YoutubeChannel=? ", Youtube["ChannelID"])
					err := row.Scan(&tmp)
					if err == sql.ErrNoRows {
						stmt, err := db.Prepare("INSERT INTO GroupYoutube (YoutubeChannel,Region,VtuberGroup_id) values(?,?,?)")
						if err != nil {
							log.Error(err)
						}

						res, err := stmt.Exec(Youtube["ChannelID"], Youtube["Region"], GroupData.ID)
						if err != nil {
							log.Error(err)
						}

						_, err = res.LastInsertId()
						if err != nil {
							log.Error(err)
						}

						defer stmt.Close()
					}
				}
			}

			for _, bil := range GroupRaw.GroupChannel.BiliBili {
				if bil != nil {
					tmp := 0
					BiliBili := bil.(map[string]interface{})
					row := db.QueryRow("SELECT id FROM GroupBiliBili WHERE AND BiliBili_SpaceID=? ", BiliBili["BiliBili_ID"])
					err := row.Scan(&tmp)
					if err == sql.ErrNoRows {
						stmt, err := db.Prepare("INSERT INTO GroupBiliBili (BiliBili_SpaceID,BiliBili_RoomID,Status,Region,VtuberGroup_id) values(?,?,?,?,?)")
						if err != nil {
							log.Error(err)
						}

						res, err := stmt.Exec(BiliBili["BiliBili_ID"].(float64), BiliBili["BiliRoom_ID"].(float64), config.UnknownStatus, BiliBili["Region"].(string), GroupData.ID)
						if err != nil {
							log.Error(err)
						}

						_, err = res.LastInsertId()
						if err != nil {
							log.Error(err)
						}

						defer stmt.Close()
					}
				}
			}

			DiscordChannel, err := GroupData.GetChannelByGroup("")
			if err != nil {
				log.Error(err)
			}
			for _, VtuberMember := range GroupRaw.Members {
				/*
					Add Member
				*/

				var MemberID int64
				YtItem := GetYtChannelInfo(VtuberMember)
				TwitterItem := GetTwitterAccountInfo(VtuberMember)
				BiliItem := GetBiliBiliAccountInfo(VtuberMember)

				Region := func() string {
					if YtItem.Brandingsettings.Channel.Country != "" {
						return YtItem.Brandingsettings.Channel.Country
					} else {
						return VtuberMember.Region
					}
				}()

				row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND VtuberName_EN=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.EnName, VtuberMember.Youtube.YtID, VtuberMember.BiliBili.BiliBiliID, VtuberMember.BiliBili.BiliRoomID)
				err = row.Scan(&MemberID)
				if err == sql.ErrNoRows {
					stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName, VtuberName_EN, VtuberName_JP, Twitter_Username, Twitter_Hashtag, Twitter_Lewd, Twitter_Avatar, Twitter_Banner, Youtube_ID, Youtube_Avatar, Youtube_Banner, BiliBili_SpaceID, BiliBili_RoomID, BiliBili_Avatar, BiliBili_Hashtag, BiliBili_Banner, Twitch_Username, Twitch_Avatar, Region, Fanbase, Status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
					if err != nil {
						log.Error(err)
					}

					TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}

					res, err := stmt.Exec(
						VtuberMember.Name,
						VtuberMember.EnName,
						VtuberMember.JpName,
						VtuberMember.Twitter.TwitterUsername,
						VtuberMember.Twitter.TwitterFanart,
						VtuberMember.Twitter.TwitterLewd,
						TwitterItem.Avatar,
						TwitterItem.Banner,
						VtuberMember.Youtube.YtID,
						YtItem.Brandingsettings.Image.Bannerexternalurl,
						VtuberMember.BiliBili.BiliBiliID,
						VtuberMember.BiliBili.BiliRoomID,
						BiliItem.Data.Face,
						VtuberMember.BiliBili.BiliBiliFanart,
						BiliItem.Data.TopPhoto,
						VtuberMember.Twitch.TwitchUsername,
						TwitchAvatar,
						Region,
						VtuberMember.Fanbase,
						VtuberMember.Status,
						MemberID,
					)
					if err != nil {
						log.Error(err)
					}

					MemberID, err = res.LastInsertId()
					if err != nil {
						log.Error(err)
					}

					defer stmt.Close()

					for _, Channel := range DiscordChannel {
						log.WithFields(log.Fields{
							"ChannelID": Channel.ChannelID,
						}).Info("Send notif")

						msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewVtuber{
							Group:      GroupData,
							Member:     VtuberMember,
							YtAvatar:   YtItem.Snippet.Thumbnails.High.URL,
							BiliAvatar: BiliItem.Data.Face,
						}.SendNotif())
						if err != nil {
							log.Error(msg, err)
							match, _ := regexp.MatchString("Unknown Channel", err.Error())
							if match {
								log.Info("Delete Discord Channel ", Channel)
								DeleteChannel(Channel.ID)
							}
						}
						User.SetDiscordChannelID(Channel.ChannelID).
							SetGroup(GroupData).
							SetMember(database.Member{
								ID:     MemberID,
								Name:   VtuberMember.Name,
								EnName: VtuberMember.EnName,
								JpName: VtuberMember.JpName,
								Group: database.Group{
									ID: GroupData.ID,
								},
							})
						err = User.SendToCache(msg.ID)
						if err != nil {
							log.Error(err)
						}
						err = engine.Reacting(map[string]string{
							"ChannelID": Channel.ChannelID,
							"State":     "New Member",
							"MessageID": msg.ID,
						}, Bot)
						if err != nil {
							log.Error(err)
						}
					}
					NewVtuberNames = append(NewVtuberNames, "`"+VtuberMember.Name+"`")
					VtuberMember.InputSubs(MemberID, YtItem, TwitterItem)

				} else if err != nil {
					log.Error(err)
				} else {
					log.WithFields(log.Fields{
						"VtuberGroup": GroupData.GroupName,
						"Vtuber":      VtuberMember.Name,
					}).Info("Update member")
					TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}
					_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=?, Twitter_Username=?, Twitter_Hashtag=?, Twitter_Lewd=?, Twitter_Avatar=?, Twitter_Banner=?, Youtube_ID=?, Youtube_Avatar=?, Youtube_Banner=?, BiliBili_SpaceID=?, BiliBili_RoomID=?, BiliBili_Avatar=?, BiliBili_Hashtag=?, BiliBili_Banner=?, Twitch_Username=?, Twitch_Avatar=?, Region=?, Fanbase=?, Status=?  Where id=?`,
						VtuberMember.Name,
						VtuberMember.EnName,
						VtuberMember.JpName,
						VtuberMember.Twitter.TwitterUsername,
						VtuberMember.Twitter.TwitterFanart,
						VtuberMember.Twitter.TwitterLewd,
						TwitterItem.Avatar,
						TwitterItem.Banner,
						VtuberMember.Youtube.YtID,
						YtItem.Brandingsettings.Image.Bannerexternalurl,
						VtuberMember.BiliBili.BiliBiliID,
						VtuberMember.BiliBili.BiliRoomID,
						BiliItem.Data.Face,
						VtuberMember.BiliBili.BiliBiliFanart,
						BiliItem.Data.TopPhoto,
						VtuberMember.Twitch.TwitchUsername,
						TwitchAvatar,
						Region,
						VtuberMember.Fanbase,
						VtuberMember.Status,
						MemberID,
					)
					if err != nil {
						log.Error(err)
					}

				}
				log.WithFields(log.Fields{
					"VtuberGroup": GroupData.GroupName,
					"Vtuber":      VtuberMember.Name,
				}).Info("Add subs info to database")

				/*
					Add subs info
				*/
				VtuberMember.InputSubs(MemberID, YtItem, TwitterItem)
			}
			if NewVtuberNames != nil {
				Vtubers := strings.Join(NewVtuberNames, ",")
				DiscordChannel, err := GroupData.GetChannelByGroup("")
				if err != nil {
					log.Error(err)
				}
				for _, Channel := range DiscordChannel {
					msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "New Update!!!! @here "+Vtubers)
					if err != nil {
						log.Error(msg, err)
						match, _ := regexp.MatchString("Unknown Channel", err.Error())
						if match {
							log.Info("Delete Discord Channel ", Channel)
							DeleteChannel(Channel.ID)
						}
					}
					_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Push "+configfile.Emoji.Livestream[0]+" to add you in "+Vtubers+" ping list")
					if err != nil {
						log.Error(err)
					}
					_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Update your roles with `"+configfile.BotPrefix.General+"tag roles @somesimpsroles` "+Vtubers)
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go Independen(&wg)
	go Agency(&wg)
	wg.Wait()
}

func DeleteChannel(id int64) error {
	_, err := db.Exec(`DELETE From Channel WHERE id=?`, id)
	if err != nil {
		return err
	}
	return nil
}

func (Data Members) InputSubs(MemberID int64, YtItem Items, TwitterItem twitterscraper.Profile) {
	var tmp int64
	Bili := Data.BiliBili.GetBiliFolow(Data.Name)
	TwitchFollow, TwitchViwers, err := Data.Twitch.GetTwitchFollowers()
	if err != nil {
		log.Error(err)
	}

	ytsubs, err := strconv.Atoi(YtItem.Statistics.SubscriberCount)
	if err != nil {
		log.Error(err)
	}

	ytvideos, err := strconv.Atoi(YtItem.Statistics.VideoCount)
	if err != nil {
		log.Error(err)
	}

	ytviews, err := strconv.Atoi(YtItem.Statistics.ViewCount)
	if err != nil {
		log.Error(err)
	}

	bilifoll := Bili.Follow.Data.Follower
	bilivideos := Bili.Video
	biliview := Bili.Like.Data.Archive.View
	twfollo := TwitterItem.FollowersCount

	row := db.QueryRow("SELECT id FROM Subscriber WHERE VtuberMember_id=? ", MemberID)
	err = row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Subscriber (Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers,Twitch_Followers,Twitch_Views,VtuberMember_id) values(?,?,?,?,?,?,?,?,?,?)")
		if err != nil {
			log.Error(err)
		}
		res, err := stmt.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, TwitchFollow, TwitchViwers, MemberID)
		if err != nil {
			log.Error(err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Error(err)
		}

		defer stmt.Close()
	} else {
		rows, err := db.Query(`SELECT Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers,Twitch_Followers,Twitch_Views FROM Subscriber WHERE VtuberMember_id=?`, MemberID)
		if err != nil {
			log.Error(err)
		}
		var (
			ytsubstmp, ytvideostmp, ytviewstmp, bilifolltmp, bilivideostmp, biliviewtmp, twfollotmp, twitchfollotmp, twitchviwerstmp int
		)

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&ytsubstmp, &ytvideostmp, &ytviewstmp, &bilifolltmp, &bilivideostmp, &biliviewtmp, &twfollotmp, &twitchfollotmp, &twitchviwerstmp)
			if err != nil {
				log.Error(err)
			}
		}

		if (ytsubs == 0 && ytsubstmp != 0) || (ytsubs == 0 && ytsubstmp == 0) {
			ytsubs = ytsubstmp
		}
		if (ytvideos == 0 && ytvideostmp != 0) || (ytvideos == 0 && ytvideostmp == 0) {
			ytvideos = ytvideostmp
		}
		if (ytviews == 0 && ytviewstmp != 0) || (ytviews == 0 && ytviewstmp == 0) {
			ytviews = ytviewstmp
		}
		if (bilifoll == 0 && bilifolltmp != 0) || (bilifoll == 0 && bilifolltmp == 0) {
			bilifoll = bilifolltmp
		}
		if (bilivideos == 0 && bilivideostmp != 0) || (bilivideos == 0 && bilivideostmp == 0) {
			bilivideos = bilivideostmp
		}
		if (biliview == 0 && biliviewtmp != 0) || (biliview == 0 && biliviewtmp == 0) {
			biliview = biliviewtmp
		}
		if (twfollo == 0 && twfollotmp != 0) || (twfollo == 0 && twfollotmp == 0) {
			twfollo = twfollotmp
		}
		if (TwitchFollow == 0 && twfollotmp != 0) || (TwitchFollow == 0 && twitchfollotmp == 0) {
			TwitchFollow = twitchfollotmp
		}
		if (TwitchViwers == 0 && twitchviwerstmp != 0) || (TwitchViwers == 0 && twitchviwerstmp == 0) {
			TwitchViwers = twitchviwerstmp
		}

		Update, err := db.Prepare(`Update Subscriber set Youtube_Subscriber=?, Youtube_Videos=?, Youtube_Views=?, BiliBili_Followers=?, BiliBili_Videos=?, BiliBili_Views=?, Twitter_Followers=?,Twitch_Followers=?,Twitch_Views=? Where id=?`)
		if err != nil {
			log.Error(err)
		}
		Update.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, TwitchFollow, TwitchViwers, tmp)
	}
}

func GetTwitterFanart(MemberID int64) string {
	rows, err := db.Query(`SELECT Twitter_Hashtag FROM Vtuber.VtuberMember where id=?`, MemberID)
	if err != nil {
		log.Error(err)
	}

	var Data string
	for rows.Next() {
		err = rows.Scan(&Data)
		if err != nil {
			log.Error(err)
		}
	}
	defer rows.Close()
	return Data
}

func LiveBiliBili(Data map[string]interface{}) error {
	var tmp int
	row := db.QueryRow("SELECT RoomID FROM LiveBiliBili WHERE RoomID=? AND VtuberMember_id=?", Data["LiveRoomID"], Data["MemberID"])
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO LiveBiliBili (RoomID,Status,Title,Thumbnails,Description,Published,ScheduledStart,Viewers,VtuberMember_id) values(?,?,?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
		res, err := stmt.Exec(Data["LiveRoomID"], Data["Status"], Data["Title"], Data["Thumbnail"], Data["Description"], Data["PublishedAt"], Data["ScheduledStart"], Data["Online"], Data["MemberID"])
		if err != nil {
			return err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return err
		}
		defer stmt.Close()
	} else if err != nil {
		return err
	} else {
		_, err := db.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, Viewers=? Where RoomID=? AND VtuberMember_id=?`, Data["Status"], Data["Title"], Data["Thumbnail"], Data["Description"], Data["PublishedAt"], Data["ScheduledStart"], Data["Online"], Data["LiveRoomID"], Data["MemberID"])
		if err != nil {
			return err
		}
	}
	return nil
}

func AddTwitchInfo(Data map[string]interface{}) error {
	var tmp int
	row := db.QueryRow("SELECT id FROM Twitch WHERE VtuberMember_id=?", Data["MemberID"])
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"Group":      Data["GroupName"],
			"VtuberName": Data["MemberName"],
		}).Info("New Member Twitch")

		stmt, err := db.Prepare("INSERT INTO Twitch (Game,Status,Title,Thumbnails,ScheduledStart,EndStream,Viewers,VtuberMember_id) values(?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
		res, err := stmt.Exec(Data["Game"], Data["Status"], Data["Title"], Data["Thumbnails"], Data["ScheduledStart"], Data["EndStream"], Data["Viewers"], Data["MemberID"])
		if err != nil {
			return err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return err
		}
		defer stmt.Close()
	} else if err != nil {
		return err
	} else {
		log.WithFields(log.Fields{
			"Group":      Data["GroupName"],
			"VtuberName": Data["MemberName"],
		}).Info("Update Member Twitch")
		_, err := db.Exec(`Update Twitch set Game=?,Status=?,Thumbnails=?,ScheduledStart=?,Viewers=? Where id=? AND VtuberMember_id=?`, tmp, Data["Game"], Data["Status"], Data["Title"], Data["Thumbnails"], Data["ScheduledStart"], Data["Viewers"], Data["MemberID"])
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
