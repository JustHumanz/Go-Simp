package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateDB() (*sql.DB, error) {
	log.Info("Create Database")

	db, err := sql.Open("mysql", os.Getenv("DBUSER")+":"+os.Getenv("DBPASS")+"@tcp("+os.Getenv("DBHOST")+":3306)/")
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}

	_, err = db.Exec("CREATE DATABASE Vtuber")

	_, err = db.Exec(`USE Vtuber`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Channel (
		id int(11) NOT NULL AUTO_INCREMENT,
		DiscordChannelID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Type int(11) NOT NULL,
		VtuberGroup_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS User (
		id int(11) NOT NULL AUTO_INCREMENT,
		DiscordID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		DiscordUserName varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberMember_id int(11) NOT NULL,
		Channel_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Twitter (
		id int(11) NOT NULL AUTO_INCREMENT,
		PermanentURL varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Author varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Likes int(11) DEFAULT NULL,
		Photos varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Videos varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Text varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		TweetID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TBiliBili (
		id int(11) NOT NULL AUTO_INCREMENT,
		PermanentURL varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Author varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Likes int(11) DEFAULT NULL,
		Photos TEXT COLLATE utf8mb4_unicode_ci DEFAULT NULL,  /*i'm not joking,they use sha1 hash for image identify,so the url very fucking long*/
		Videos varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Text TEXT COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Dynamic_id varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS VtuberGroup (
		id int(11) NOT NULL AUTO_INCREMENT,
		VtuberGroupName varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberGroupIcon varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS VtuberMember (
		id int(11) NOT NULL AUTO_INCREMENT,
		VtuberName varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberName_EN varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberName_JP varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Hashtag varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		BiliBili_Hashtag varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Youtube_ID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Youtube_Avatar varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		BiliBili_SpaceID INT(11) DEFAULT NULL,
		BiliBili_RoomID INT(11) DEFAULT NULL,
		BiliBili_Avatar varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Twitter_Username varchar(24) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Region varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberGroup_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Youtube (
		id int(11) NOT NULL AUTO_INCREMENT,
		VideoID varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Type varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Status varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Title varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Thumbnails varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Description text COLLATE utf8mb4_unicode_ci NOT NULL,
		PublishedAt timestamp NOT NULL DEFAULT current_timestamp(),
		ScheduledStart timestamp NOT NULL DEFAULT current_timestamp(),
		EndStream timestamp NOT NULL DEFAULT current_timestamp(),
		Viewers varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Length varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS BiliBili (
		id int(11) NOT NULL AUTO_INCREMENT,
		VideoID varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Type varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Title varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Thumbnails varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Description text COLLATE utf8mb4_unicode_ci NOT NULL,
		UploadDate timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
		Viewers int(11) COLLATE utf8mb4_unicode_ci NOT NULL,
		Length varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS LiveBiliBili (
		id int(11) NOT NULL AUTO_INCREMENT,
		RoomID varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Status varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Title varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Thumbnails varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Description text COLLATE utf8mb4_unicode_ci NOT NULL,
		Published timestamp NOT NULL DEFAULT current_timestamp(),
		ScheduledStart timestamp NOT NULL DEFAULT current_timestamp(),
		Viewers varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Vtuber.Subscriber (
		id INT NOT NULL AUTO_INCREMENT,
		Youtube_Subs INT(11) NULL,
		Youtube_Videos INT(11) NULL,
		Youtube_Views INT(11) NULL,
		BiliBili_Follows INT(11) NULL,
		BiliBili_Videos INT(11) NULL,
		BiliBili_Views INT(11) NULL,
		Twitter_Follows INT(11) NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
		);`)
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(30)
	db.SetMaxOpenConns(50)
	log.Info("DB ok")
	return db, nil
}

func GetHashtag(Group string) []database.MemberGroupID {
	rows, err := db.Query(`SELECT VtuberMember.id,VtuberName,VtuberName_JP,VtuberGroup_id,Hashtag,VtuberGroupName,VtuberGroupIcon FROM VtuberMember INNER Join VtuberGroup ON VtuberGroup.id = VtuberMember.VtuberGroup_id WHERE VtuberGroup.VtuberGroupName =?`, Group)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	Data := []database.MemberGroupID{}
	for rows.Next() {
		var list database.MemberGroupID
		err = rows.Scan(&list.MemberID, &list.EnName, &list.JpName, &list.GroupID, &list.TwitterHashtags, &list.GroupName, &list.GroupIcon)
		if err != nil {
			log.Error(err)
		}

		Data = append(Data, list)

	}
	return Data
}

func (Data Member) BliBiliFace() string {
	if Data.BiliBiliID == 0 {
		return ""
	} else {
		var (
			Info    Avatar
			body    []byte
			errcurl error
			url     = "https://api.bilibili.com/x/space/acc/info?mid=" + strconv.Itoa(Data.BiliBiliID)
		)
		body, errcurl = engine.Curl(url, nil)
		if body == nil {
			log.Info("Not daijobu,trying use multitor")
			body, errcurl = engine.CoolerCurl(url)

			if errcurl != nil {
				log.Error(errcurl)
				return ""
			}
		} else if errcurl != nil {
			log.Error(errcurl)
			return ""
		}
		err := json.Unmarshal(body, &Info)
		if err != nil {
			log.Error(err)
			return ""
		}

		return strings.Replace(Info.Data.Face, "http", "https", -1)
	}
}

func AddData(Data Vtuber) {
	var (
		wg sync.WaitGroup
	)
	wg.Add(2)
	go func() {
		defer wg.Done()

		var (
			GroupID   int64
			GroupName = "Independen"
			GroupIcon = "https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/independen.png"
		)
		row := db.QueryRow("SELECT id FROM VtuberGroup WHERE VtuberGroupName=?", GroupName)
		err := row.Scan(&GroupID)
		if err == sql.ErrNoRows {
			log.Error(err)
			stmt, err := db.Prepare("INSERT INTO VtuberGroup (VtuberGroupName,VtuberGroupIcon) values(?,?)")
			engine.BruhMoment(err, "", false)

			res, err := stmt.Exec(GroupName, GroupIcon)
			engine.BruhMoment(err, "", false)

			GroupID, err = res.LastInsertId()
			engine.BruhMoment(err, "", false)

			defer stmt.Close()
		} else {
			log.WithFields(log.Fields{
				"VtuberGroup": GroupName,
			}).Info("Update Vtuber Group Data")
			Update, err := db.Prepare(`Update VtuberGroup set VtuberGroupName=?, VtuberGroupIcon=? Where id=?`)
			engine.BruhMoment(err, "", false)
			Update.Exec(GroupName, GroupIcon, GroupID)
		}
		for _, Member := range Data.Vtuber.Independen.Members {
			/*
				Add Member
			*/
			var tmp int64
			row := db.QueryRow("SELECT id FROM VtuberMember WHERE Youtube_ID=? ", strings.Join(Member.YtID, "\n"))
			err := row.Scan(&tmp)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
				engine.BruhMoment(err, "", false)

				res, err := stmt.Exec(Member.Name, Member.ENName, Member.JPName, Member.Hashtag.Twitter, Member.Hashtag.BiliBili, strings.Join(Member.YtID, "\n"), Member.YtAvatar(), GroupID, Member.Region, Member.BiliBiliID, Member.BiliRoomID, Member.BliBiliFace(), Member.TwitterName)
				engine.BruhMoment(err, "", false)

				_, err = res.LastInsertId()
				engine.BruhMoment(err, "", false)

				defer stmt.Close()
				New = append(New, NewVtuber{
					Group: database.GroupName{
						ID:        GroupID,
						NameGroup: "Independen",
						IconURL:   "https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/independen.png",
					},
					Member: Member,
				})
			} else {
				log.WithFields(log.Fields{
					"VtuberGroup": "Independen",
					"Vtuber":      Member.ENName,
				}).Info("already added...")
				_, err := db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?  Where id=?`, Member.Name, Member.ENName, Member.JPName, Member.Hashtag.Twitter, Member.Hashtag.BiliBili, Member.Region, strings.Join(Member.YtID, "\n"), Member.BiliBiliID, Member.BiliRoomID, Member.BliBiliFace(), Member.YtAvatar(), Member.TwitterName, tmp)
				if err != nil {
					log.Error(err)
				}
				log.WithFields(log.Fields{
					"VtuberGroup": "Independen",
					"Vtuber":      Member.ENName,
				}).Info("update member")
			}
			log.WithFields(log.Fields{
				"VtuberGroup": "Independen",
				"Vtuber":      Member.ENName,
			}).Info("Add subs info to database")

			/*
				Add subs info
			*/
			Member.InputSubs(tmp)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()

		wg2 := new(sync.WaitGroup)
		for i := 0; i < len(Data.Vtuber.Group); i++ {
			wg2.Add(1)
			Group := Data.Vtuber.Group[i]
			go func() {
				/*
					Add Group
				*/
				log.WithFields(log.Fields{
					"VtuberGroup":     Group.GroupName,
					"VtuberGroupIcon": Group.GroupIcon,
				}).Info("Add Group")
				defer wg2.Done()

				var GroupID int64
				row := db.QueryRow("SELECT id FROM VtuberGroup WHERE VtuberGroupName=? ", Group.GroupName)
				err := row.Scan(&GroupID)
				if err != nil || err == sql.ErrNoRows {
					stmt, err := db.Prepare("INSERT INTO VtuberGroup (VtuberGroupName,VtuberGroupIcon) values(?,?)")
					engine.BruhMoment(err, "", false)

					res, err := stmt.Exec(Group.GroupName, Group.GroupIcon)
					engine.BruhMoment(err, "", false)

					GroupID, err = res.LastInsertId()
					engine.BruhMoment(err, "", false)

					defer stmt.Close()
				} else {
					log.WithFields(log.Fields{
						"VtuberGroup":     Group.GroupName,
						"VtuberGroupIcon": Group.GroupIcon,
					}).Info("Update Vtuber Group Data")
					Update, err := db.Prepare(`Update VtuberGroup set VtuberGroupName=?, VtuberGroupIcon=? Where id=?`)
					engine.BruhMoment(err, "", false)
					Update.Exec(Group.GroupName, Group.GroupIcon, GroupID)
				}

				for j := 0; j < len(Group.Members); j++ {
					Member := Group.Members[j]
					/*
						Add Member
					*/
					var tmp int64
					row := db.QueryRow("SELECT id FROM VtuberMember WHERE Youtube_ID=? ", strings.Join(Member.YtID, "\n"))
					err := row.Scan(&tmp)
					if err != nil || err == sql.ErrNoRows {
						stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
						engine.BruhMoment(err, "", false)

						res, err := stmt.Exec(Member.Name, Member.ENName, Member.JPName, Member.Hashtag.Twitter, Member.Hashtag.BiliBili, strings.Join(Member.YtID, "\n"), Member.YtAvatar(), GroupID, Member.Region, Member.BiliBiliID, Member.BiliRoomID, Member.BliBiliFace(), Member.TwitterName)
						engine.BruhMoment(err, "", false)

						_, err = res.LastInsertId()
						engine.BruhMoment(err, "", false)

						defer stmt.Close()
						New = append(New, NewVtuber{
							Group: database.GroupName{
								ID:        GroupID,
								NameGroup: Group.GroupName,
								IconURL:   Group.GroupIcon,
							},
							Member: Member,
						})
					} else {
						log.WithFields(log.Fields{
							"VtuberGroup": Group.GroupName,
							"Vtuber":      Member.ENName,
						}).Info("already added...")
						_, err := db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?  Where id=?`, Member.Name, Member.ENName, Member.JPName, Member.Hashtag.Twitter, Member.Hashtag.BiliBili, Member.Region, strings.Join(Member.YtID, "\n"), Member.BiliBiliID, Member.BiliRoomID, Member.BliBiliFace(), Member.YtAvatar(), Member.TwitterName, tmp)
						if err != nil {
							log.Error(err)
						}
						log.WithFields(log.Fields{
							"VtuberGroup": Group.GroupName,
							"Vtuber":      Member.ENName,
						}).Info("update member")
					}
					log.WithFields(log.Fields{
						"VtuberGroup": Group.GroupName,
						"Vtuber":      Member.ENName,
					}).Info("Add subs info to database")

					/*
						Add subs info
					*/
					Member.InputSubs(tmp)
					time.Sleep(1 * time.Second)
				}
			}()
		}
		wg2.Wait()
	}()
	wg.Wait()
}

func (Data Member) InputSubs(MemberID int64) {
	var tmp int64
	row := db.QueryRow("SELECT id FROM Subscriber WHERE VtuberMember_id=? ", MemberID)
	err := row.Scan(&tmp)
	Subs := Data.GetYtSubs()
	Bili := Data.GetBiliFolow()

	if err != nil || err == sql.ErrNoRows {
		log.Error(err)
		stmt, err := db.Prepare("INSERT INTO Subscriber (Youtube_Subs,Youtube_Videos,Youtube_Views,BiliBili_Follows,BiliBili_Videos,BiliBili_Views,Twitter_Follows,VtuberMember_id) values(?,?,?,?,?,?,?,?)")
		engine.BruhMoment(err, "", false)
		res, err := stmt.Exec(Subs[0].Data.Subscribers, Subs[0].Data.Videos, Subs[0].Data.Views, Bili.Follow.Data.Follower, Bili.Video, Bili.Like.Data.Archive.View, Data.GetTwitterFollow(), MemberID)
		engine.BruhMoment(err, "", false)

		_, err = res.LastInsertId()
		engine.BruhMoment(err, "", false)

		defer stmt.Close()
	} else {
		Update, err := db.Prepare(`Update Subscriber set Youtube_Subs=?, Youtube_Videos=?, Youtube_Views=?, BiliBili_Follows=?, BiliBili_Videos=?, BiliBili_Views=?, Twitter_Follows=? Where id=?`)
		engine.BruhMoment(err, "", false)
		Update.Exec(Subs[0].Data.Subscribers, Subs[0].Data.Videos, Subs[0].Data.Views, Bili.Follow.Data.Follower, Bili.Video, Bili.Like.Data.Archive.View, Data.GetTwitterFollow(), tmp)

	}
}

func GetHastagMember(MemberID int64) string {
	rows, err := db.Query(`SELECT Hashtag FROM Vtuber.VtuberMember where id=?`, MemberID)
	engine.BruhMoment(err, "", false)

	var Data string
	for rows.Next() {
		err = rows.Scan(&Data)
		engine.BruhMoment(err, "", false)
	}
	defer rows.Close()
	return Data
}

func (Data InputTwitter) InputData() {
	if Data.MemberID != 0 {
		var (
			tmp   string
			Video string
		)

		Photos := strings.Join(Data.TwitterData.Photos, "\n")
		if Data.TwitterData.Videos != nil {
			Video = "https://pbs.twimg.com/tweet_video/" + Data.TwitterData.Videos[0].ID + ".mp4"
		}

		row := db.QueryRow("SELECT PermanentURL FROM Twitter WHERE PermanentURL=? AND VtuberMember_id=?", Data.TwitterData.PermanentURL, Data.MemberID)
		err := row.Scan(&tmp)
		if err != nil || err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Username": Data.TwitterData.Username,
				"Like":     Data.TwitterData.Likes,
				"TweetID":  Data.TwitterData.ID,
			}).Info("New Tweet")
			stmt, err := db.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
			engine.BruhMoment(err, "", false)

			res, err := stmt.Exec(Data.TwitterData.PermanentURL, Data.TwitterData.Username, Data.TwitterData.Likes, Photos, Video, Data.TwitterData.Text, Data.TwitterData.ID, Data.MemberID)
			engine.BruhMoment(err, "", false)

			_, err = res.LastInsertId()
			engine.BruhMoment(err, "", false)

			defer stmt.Close()
		} else {
			log.Info("already added...")
		}
	}
}

func LiveBiliBili(Data map[string]interface{}) bool {
	var tmp int
	row := db.QueryRow("SELECT RoomID FROM LiveBiliBili WHERE RoomID=? AND VtuberMember_id=?", Data["LiveRoomID"], Data["MemberID"])
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO LiveBiliBili (RoomID,Status,Title,Thumbnails,Description,Published,ScheduledStart,Viewers,VtuberMember_id) values(?,?,?,?,?,?,?,?,?)")
		if err != nil {
			log.Error(err)
		}
		res, err := stmt.Exec(Data["LiveRoomID"], Data["Status"], Data["Title"], Data["Thumbnail"], Data["Description"], Data["PublishedAt"], Data["ScheduledStart"], Data["Online"], Data["MemberID"])
		if err != nil {
			log.Error(err)
		}

		_, err = res.LastInsertId()
		if err != nil {
			log.Error(err)
		}
		defer stmt.Close()
		return true
	} else {
		log.Info("already added...")
		log.Info("Update LiveBiliBili")
		_, err := db.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, Viewers=? Where RoomID=? AND VtuberMember_id=?`, Data["Status"], Data["Title"], Data["Thumbnail"], Data["Description"], Data["PublishedAt"], Data["ScheduledStart"], Data["Online"], Data["LiveRoomID"], Data["MemberID"])
		engine.BruhMoment(err, "", false)
		return false
	}

}
