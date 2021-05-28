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
	log "github.com/sirupsen/logrus"
)

func AddData(Data Vtuber) {
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
			DiscordChannel := GroupData.GetChannelByGroup(VtuberMember.Region)
			/*
				Add Member
			*/
			var MemberID int64
			row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND VtuberName_EN=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.ENName, VtuberMember.Youtube.YtID, VtuberMember.BiliBili.BiliBiliID, VtuberMember.BiliBili.BiliRoomID)
			err := row.Scan(&MemberID)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Twitter_Hashtag,Twitter_Lewd,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username,Twitch_Username,Twitch_Avatar,Fanbase,Status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
				if err != nil {
					log.Error(err)
				}
				BiliFace, err := VtuberMember.BiliBili.BliBiliFace()
				if err != nil {
					log.Error(err)
				}

				TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}
				res, err := stmt.Exec(VtuberMember.Name, VtuberMember.ENName,
					VtuberMember.JPName, VtuberMember.Twitter.TwitterFanart,
					VtuberMember.Twitter.TwitterLewd, VtuberMember.BiliBili.BiliBiliFanart,
					VtuberMember.Youtube.YtID, VtuberMember.Youtube.YtAvatar(),
					GroupData.ID, VtuberMember.Region, VtuberMember.BiliBili.BiliBiliID,
					VtuberMember.BiliBili.BiliRoomID, BiliFace,
					VtuberMember.Twitch.TwitchUsername, VtuberMember.Twitch.TwitchUsername,
					TwitchAvatar,
					VtuberMember.Fanbase, VtuberMember.Status)
				if err != nil {
					log.Error(err)
				}

				MemberID, err = res.LastInsertId()
				if err != nil {
					log.Error(err)
				}

				defer stmt.Close()

				for _, Channel := range DiscordChannel {
					msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewVtuber{
						Group:  GroupData,
						Member: VtuberMember,
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
								ID:      MemberID,
								Name:    VtuberMember.Name,
								EnName:  VtuberMember.ENName,
								JpName:  VtuberMember.JPName,
								GroupID: GroupData.ID,
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
				VtuberMember.InputSubs(MemberID)
				//New.SendNotif(Bot)
			} else if err != nil {
				log.Error(err)
			} else {
				TwitchAvatar, err := VtuberMember.Twitch.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}
				BiliFace, err := VtuberMember.BiliBili.BliBiliFace()
				if err != nil {
					log.Error(err)
				} else {
					log.WithFields(log.Fields{
						"VtuberGroup": "Independen",
						"Vtuber":      VtuberMember.ENName,
					}).Info("Update member")
					_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Twitter_Hashtag=?,Twitter_Lewd=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?,Twitch_Username=?,Twitch_Avatar=?,Fanbase=?,Status=?  Where id=?`,
						VtuberMember.Name, VtuberMember.ENName,
						VtuberMember.JPName, VtuberMember.Twitter.TwitterFanart,
						VtuberMember.Twitter.TwitterLewd,
						VtuberMember.BiliBili.BiliBiliFanart, VtuberMember.Region,
						VtuberMember.Youtube.YtID, VtuberMember.BiliBili.BiliBiliID,
						VtuberMember.BiliBili.BiliRoomID, BiliFace,
						VtuberMember.Youtube.YtAvatar(), VtuberMember.Twitter.TwitterUsername,
						VtuberMember.Twitch.TwitchUsername,
						TwitchAvatar, VtuberMember.Fanbase, VtuberMember.Status,
						MemberID)
					if err != nil {
						log.Error(err)
					}
				}
			}
			log.WithFields(log.Fields{
				"VtuberGroup": "Independen",
				"Vtuber":      VtuberMember.ENName,
			}).Info("Add subs info to database")

			/*
				Add subs info
			*/
			VtuberMember.InputSubs(MemberID)
			//time.Sleep(1 * time.Second)
		}

		if NewVtuberNamesIndependen != nil {
			Vtubers := strings.Join(NewVtuberNamesIndependen, ",")
			DiscordChannel := GroupData.GetChannelByGroup("")
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
				_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Push "+configfile.Emoji.Livestream[0]+" to add you in `"+Vtubers+"` ping list")
				if err != nil {
					log.Error(err)
				}

				_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Update your roles with `"+configfile.BotPrefix.General+"tag roles @somesimpsroles "+Vtubers+"`")
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

			DiscordChannel := GroupData.GetChannelByGroup("")
			for _, v := range GroupRaw.Members {
				/*
					Add Member
				*/
				var MemberID int64
				row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND (Youtube_ID=? OR BiliBili_SpaceID=? OR BiliBili_RoomID=?)", v.Name, v.Youtube.YtID, v.BiliBili.BiliBiliID, v.BiliBili.BiliRoomID)
				err := row.Scan(&MemberID)
				if err == sql.ErrNoRows {
					stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Twitter_Hashtag,Twitter_Lewd,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username,Twitch_Username,Twitch_Avatar,Fanbase,Status) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
					if err != nil {
						log.Error(err)
					}

					BiliFace, err := v.BiliBili.BliBiliFace()
					if err != nil {
						log.Error(err)
					}

					TwitchAvatar, err := v.Twitch.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}

					res, err := stmt.Exec(v.Name, v.ENName, v.JPName,
						v.Twitter.TwitterFanart, v.Twitter.TwitterLewd,
						v.BiliBili.BiliBiliFanart, v.Youtube.YtID, v.Youtube.YtAvatar(),
						GroupData.ID, v.Region, v.BiliBili.BiliBiliID,
						v.BiliBili.BiliRoomID, BiliFace,
						v.Twitter.TwitterUsername, v.Twitch.TwitchUsername,
						TwitchAvatar, v.Fanbase, v.Status)
					if err != nil {
						log.Error(err)
					}

					MemberID, err = res.LastInsertId()
					if err != nil {
						log.Error(err)
					}

					defer stmt.Close()

					for _, Channel := range DiscordChannel {
						msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, NewVtuber{
							Group:  GroupData,
							Member: v,
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
								ID:      MemberID,
								Name:    v.Name,
								EnName:  v.ENName,
								JpName:  v.JPName,
								GroupID: GroupData.ID,
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
					NewVtuberNames = append(NewVtuberNames, "`"+v.Name+"`")
					v.InputSubs(MemberID)

				} else if err != nil {
					log.Error(err)
				} else {
					log.WithFields(log.Fields{
						"VtuberGroup": GroupData.GroupName,
						"Vtuber":      v.Name,
					}).Info("Update member")
					TwitchAvatar, err := v.Twitch.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}

					BiliFace, err := v.BiliBili.BliBiliFace()
					if err != nil {
						log.Error(err)
					} else {
						_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Twitter_Hashtag=?,Twitter_Lewd=?,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?,Twitch_Username=?,Twitch_Avatar=?,Fanbase=?,Status=? Where id=?`,
							v.Name, v.ENName, v.JPName,
							v.Twitter.TwitterFanart, v.Twitter.TwitterLewd,
							v.BiliBili.BiliBiliFanart, v.Region,
							v.Youtube.YtID, v.BiliBili.BiliBiliID,
							v.BiliBili.BiliRoomID, BiliFace,
							v.Youtube.YtAvatar(), v.Twitter.TwitterUsername,
							v.Twitch.TwitchUsername, TwitchAvatar,
							v.Fanbase, v.Status, MemberID)
						if err != nil {
							log.Error(err)
						}
					}
				}
				log.WithFields(log.Fields{
					"VtuberGroup": GroupData.GroupName,
					"Vtuber":      v.Name,
				}).Info("Add subs info to database")

				/*
					Add subs info
				*/
				v.InputSubs(MemberID)
			}
			if NewVtuberNames != nil {
				Vtubers := strings.Join(NewVtuberNames, ",")
				DiscordChannel := GroupData.GetChannelByGroup("")
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
					_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Push "+configfile.Emoji.Livestream[0]+" to add you in `"+Vtubers+"` ping list")
					if err != nil {
						log.Error(err)
					}
					_, err = Bot.ChannelMessageSend(Channel.ChannelID, "Update your roles with `"+configfile.BotPrefix.General+"tag roles @somesimpsroles "+Vtubers+"`")
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

/*
//GetChannelByGroupGet DiscordChannelID from VtuberGroup
func GetChannelByGroup(VtuberGroupID int64) ([]string, []int64) {
	var (
		channellist []string
		ids         []int64
	)
	rows, err := db.Query(`SELECT id,DiscordChannelID FROM Channel WHERE VtuberGroup_id=?`, VtuberGroupID)
	if err != nil {
		log.Error(err)
	}

	defer rows.Close()
	for rows.Next() {
		var (
			list string
			id   int64
		)
		err = rows.Scan(&id, &list)
		if err != nil {
			log.Error(err)
		}
		channellist = append(channellist, list)
		ids = append(ids, id)
	}
	return channellist, ids
}
*/

func (Data Members) InputSubs(MemberID int64) {
	var tmp int64
	row := db.QueryRow("SELECT id FROM Subscriber WHERE VtuberMember_id=? ", MemberID)
	err := row.Scan(&tmp)
	Subs := Data.Youtube.GetYtSubs()
	Bili := Data.BiliBili.GetBiliFolow(Data.Name)

	ytsubs, _ := strconv.Atoi(Subs.Items[0].Statistics.SubscriberCount)
	ytvideos, _ := strconv.Atoi(Subs.Items[0].Statistics.VideoCount)
	ytviews, _ := strconv.Atoi(Subs.Items[0].Statistics.ViewCount)
	bilifoll := Bili.Follow.Data.Follower
	bilivideos := Bili.Video
	biliview := Bili.Like.Data.Archive.View
	twfollo, twerr := Data.Twitter.GetTwitterFollow()
	if twerr != nil {
		log.Error(err)
	}

	if err != nil || err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Subscriber (Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers,VtuberMember_id) values(?,?,?,?,?,?,?,?)")
		if err != nil {
			log.Error(err)
		}
		res, err := stmt.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, MemberID)
		if err != nil {
			log.Error(err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Error(err)
		}

		defer stmt.Close()
	} else {
		rows, err := db.Query(`SELECT Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers FROM Subscriber WHERE VtuberMember_id=?`, MemberID)
		if err != nil {
			log.Error(err)
		}
		var (
			ytsubstmp, ytvideostmp, ytviewstmp, bilifolltmp, bilivideostmp, biliviewtmp, twfollotmp int
		)

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&ytsubstmp, &ytvideostmp, &ytviewstmp, &bilifolltmp, &bilivideostmp, &biliviewtmp, &twfollotmp)
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

		Update, err := db.Prepare(`Update Subscriber set Youtube_Subscriber=?, Youtube_Videos=?, Youtube_Views=?, BiliBili_Followers=?, BiliBili_Videos=?, BiliBili_Views=?, Twitter_Followers=? Where id=?`)
		if err != nil {
			log.Error(err)
		}
		Update.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, tmp)
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
