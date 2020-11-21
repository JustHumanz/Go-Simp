package main

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateDB(Data config.ConfigFile) error {
	log.Info("Create Database")

	db, err := sql.Open("mysql", Data.SQL.User+":"+Data.SQL.Pass+"@tcp("+Data.SQL.Host+":3306)/")
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}
	_, err = db.Exec("CREATE DATABASE Vtuber")

	_, err = db.Exec(`USE Vtuber`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Channel (
		id int(11) NOT NULL AUTO_INCREMENT,
		DiscordChannelID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Type int(11) NOT NULL,
		VtuberGroup_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS User (
		id int(11) NOT NULL AUTO_INCREMENT,
		DiscordID varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		DiscordUserName varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Reminder INT(2) DEFAULT 0,
		Human TINYINT DEFAULT 1,
		VtuberMember_id int(11) NOT NULL,
		Channel_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS VtuberGroup (
		id int(11) NOT NULL AUTO_INCREMENT,
		VtuberGroupName varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		VtuberGroupIcon varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return err
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
		return err
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
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS BiliBili (
		id int(11) NOT NULL AUTO_INCREMENT,
		VideoID varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Type varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Title varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Thumbnails varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Description text COLLATE utf8mb4_unicode_ci NOT NULL,
		UploadDate timestamp NOT NULL DEFAULT current_timestamp(),
		Viewers int(11) COLLATE utf8mb4_unicode_ci NOT NULL,
		Length varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return err
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
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Vtuber.Subscriber (
		id INT NOT NULL AUTO_INCREMENT,
		Youtube_Subscriber INT(11) NULL,
		Youtube_Videos INT(11) NULL,
		Youtube_Views INT(11) NULL,
		BiliBili_Followers INT(11) NULL,
		BiliBili_Videos INT(11) NULL,
		BiliBili_Views INT(11) NULL,
		Twitter_Followers INT(11) NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
		);`)

	log.Info("Create stored-procedure")

	log.Info("Create GetYt")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetYt;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetYt
		(
			memid int,
			grpid int,
			lmt int,
			sts varchar(11)
		)
		BEGIN
		IF sts = 'upcoming' THEN
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
			AND Status='upcoming' 
			Order by ScheduledStart ASC Limit lmt;
		ELSEIF sts = 'live' OR sts = 'private' THEN 
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
			AND Status=sts
			Limit lmt;
		ELSE 
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
			AND Status='past'
			AND EndStream !='' order by EndStream DESC Limit lmt;
			
		end if;	
		END`)
	if err != nil {
		return err
	}

	log.Info("Create GetYtByReg")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetYtByReg;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetYtByReg
		(
			grpid int,
			sts varchar(11),
			reg varchar(2)
		)
		BEGIN
		IF sts = 'upcoming' THEN
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where VtuberGroup.id=grpid AND Status='upcoming' AND Region=reg Order by ScheduledStart ASC Limit 3;

		ELSEIF sts = 'live' OR sts = 'private' THEN 
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where VtuberGroup.id=grpid AND Status=sts AND Region=reg Limit 3;
		ELSE 
			SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,
			Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers 
			FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where VtuberGroup.id=grpid AND Status='past' AND Region=reg AND EndStream !='' order by EndStream DESC Limit 3;
			
		end if;	
		END`)
	if err != nil {
		return err
	}

	log.Info("Create GetVtuberName")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetVtuberName;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetVtuberName
		(
			GroupID int
		)
		BEGIN
			SELECT id,VtuberName,VtuberName_EN,VtuberName_JP,Youtube_ID,BiliBili_SpaceID,BiliBili_RoomID,
			Region,Hashtag,BiliBili_Hashtag,BiliBili_Avatar,Twitter_Username,Youtube_Avatar 
			FROM Vtuber.VtuberMember WHERE VtuberGroup_id=GroupID 
			Order by Region,VtuberGroup_id;
		END`)
	if err != nil {
		return err
	}
	log.Info("Create GetLiveBiliBili")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetLiveBiliBili;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetLiveBiliBili
		(
			GroupID int,
			MemberID int,
			Sts varchar(11),
			lmt int
		)
		BEGIN
			SELECT RoomID,Status,Title,Thumbnails,Description,ScheduledStart,Viewers,VtuberName_EN,
			VtuberName_JP,BiliBili_Avatar FROM Vtuber.LiveBiliBili 
			Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where 
			(VtuberGroup.id=GroupID or VtuberMember.id=MemberID) 
			AND Status=Sts Order by ScheduledStart ASC Limit lmt;
		END`)
	if err != nil {
		return err
	}

	log.Info("Create GetSpaceBiliBili")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetSpaceBiliBili;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetSpaceBiliBili
		(
			GroupID int,
			MemberID int
		)
		BEGIN
		IF GroupID > 0 THEN
			SELECT VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberName_EN,VtuberName_JP,BiliBili_Avatar FROM Vtuber.BiliBili 
			Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where (VtuberGroup.id=GroupID or VtuberMember.id=MemberID) Order by UploadDate DESC limit 3;
		Else 
			SELECT VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberName_EN,VtuberName_JP,BiliBili_Avatar FROM Vtuber.BiliBili 
			Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
			Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
			Where (VtuberGroup.id=GroupID or VtuberMember.id=MemberID) Order by UploadDate DESC limit 5;		

		end if;						
		END`)
	if err != nil {
		return err
	}

	log.Info("Create GetArt")
	_, err = db.Exec(`DROP PROCEDURE IF EXISTS GetArt;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`CREATE PROCEDURE GetArt
		(
			GroupID int,
			MemberID int,
			State varchar(11)
		)
		BEGIN
		IF State = 'twitter' THEN
			SELECT VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Vtuber.Twitter 
			Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id 
			Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id 
			where (VtuberGroup.id=GroupID OR VtuberMember.id=MemberID)  ORDER by RAND() LIMIT 1;
		else
			SELECT VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Vtuber.TBiliBili  
			Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id 
			Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id 
			where (VtuberGroup.id=GroupID OR VtuberMember.id=MemberID)  ORDER by RAND() LIMIT 1;
			
		end if;
		END`)

	if err != nil {
		return err
	}
	log.Info("DB ok")
	db.Close()
	return nil
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

func AddData(Data Vtuber) {

	var (
		GroupIDIndependen        int64
		GroupName                = "Independen"
		GroupIcon                = "https://cdn.human-z.tech/404.jpg"
		NewVtuberNamesIndependen []string
	)
	row := db.QueryRow("SELECT id FROM VtuberGroup WHERE VtuberGroupName=?", GroupName)
	err := row.Scan(&GroupIDIndependen)
	if err == sql.ErrNoRows {
		log.Error(err)
		stmt, err := db.Prepare("INSERT INTO VtuberGroup (VtuberGroupName,VtuberGroupIcon) values(?,?)")
		engine.BruhMoment(err, "", false)

		res, err := stmt.Exec(GroupName, GroupIcon)
		engine.BruhMoment(err, "", false)

		GroupIDIndependen, err = res.LastInsertId()
		engine.BruhMoment(err, "", false)

		defer stmt.Close()
	} else {
		log.WithFields(log.Fields{
			"VtuberGroup": GroupName,
		}).Info("Update Vtuber Group Data")
		Update, err := db.Prepare(`Update VtuberGroup set VtuberGroupName=?, VtuberGroupIcon=? Where id=?`)
		engine.BruhMoment(err, "", false)
		Update.Exec(GroupName, GroupIcon, GroupIDIndependen)
	}

	Channels := GetChannelByGroup(GroupIDIndependen)
	wg2 := new(sync.WaitGroup)
	for _, VtuberMember := range Data.Vtuber.Independen.Members {
		wg2.Add(1)
		go func(VtuberMember Member, wg2 *sync.WaitGroup) {
			/*
				Add Member
			*/
			var tmp int64
			defer wg2.Done()
			row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID)
			err := row.Scan(&tmp)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
				engine.BruhMoment(err, "", false)

				res, err := stmt.Exec(VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.YtID, VtuberMember.YtAvatar(), GroupIDIndependen, VtuberMember.Region, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, VtuberMember.BliBiliFace(), VtuberMember.TwitterName)
				engine.BruhMoment(err, "", false)

				id, err := res.LastInsertId()
				engine.BruhMoment(err, "", false)

				defer stmt.Close()

				for _, Channel := range Channels {
					msg, err := Bot.ChannelMessageSendEmbed(Channel, NewVtuber{
						Group: database.GroupName{
							ID:        GroupIDIndependen,
							NameGroup: "Independen",
							IconURL:   "https://raw.githubusercontent.com/JustHumanz/Go-simp/main/Img/independen.png",
						},
						Member: VtuberMember,
					}.SendNotif())
					if err != nil {
						log.Error(msg, err)
					}
				}

				NewVtuberNamesIndependen = append(NewVtuberNamesIndependen, "`"+VtuberMember.Name+"`")
				VtuberMember.InputSubs(id)
				//New.SendNotif(Bot)
			} else {
				log.WithFields(log.Fields{
					"VtuberGroup": "Independen",
					"Vtuber":      VtuberMember.ENName,
				}).Info("already added...")
				_, err := db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?  Where id=?`, VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.Region, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, VtuberMember.BliBiliFace(), VtuberMember.YtAvatar(), VtuberMember.TwitterName, tmp)
				if err != nil {
					log.Error(err)
				}
				log.WithFields(log.Fields{
					"VtuberGroup": "Independen",
					"Vtuber":      VtuberMember.ENName,
				}).Info("update member")
			}
			log.WithFields(log.Fields{
				"VtuberGroup": "Independen",
				"Vtuber":      VtuberMember.ENName,
			}).Info("Add subs info to database")

			/*
				Add subs info
			*/
			VtuberMember.InputSubs(tmp)
			//time.Sleep(1 * time.Second)
		}(VtuberMember, wg2)
	}
	wg2.Wait()

	if NewVtuberNamesIndependen != nil {
		for _, Channel := range GetChannelByGroup(GroupIDIndependen) {
			msg, err := Bot.ChannelMessageSend(Channel, "New Update!!!! @here "+strings.Join(NewVtuberNamesIndependen, ","))
			if err != nil {
				log.Error(msg, err)
			}
		}
	}

	for i := 0; i < len(Data.Vtuber.Group); i++ {
		var (
			NewVtuberNames []string
			GroupID        int64
		)
		Group := Data.Vtuber.Group[i]
		/*
			Add Group
		*/
		log.WithFields(log.Fields{
			"VtuberGroup":     Group.GroupName,
			"VtuberGroupIcon": Group.GroupIcon,
		}).Info("Add Group")

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

		Channels := GetChannelByGroup(GroupID)
		wg2 := new(sync.WaitGroup)
		for j := 0; j < len(Group.Members); j++ {
			wg2.Add(1)
			go func(VtuberMember Member, wg2 *sync.WaitGroup) {
				defer wg2.Done()
				/*
					Add Member
				*/
				var tmp int64
				row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID)
				err := row.Scan(&tmp)
				if err != nil || err == sql.ErrNoRows {
					stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
					engine.BruhMoment(err, "", false)

					res, err := stmt.Exec(VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.YtID, VtuberMember.YtAvatar(), GroupID, VtuberMember.Region, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, VtuberMember.BliBiliFace(), VtuberMember.TwitterName)
					engine.BruhMoment(err, "", false)

					id, err := res.LastInsertId()
					engine.BruhMoment(err, "", false)

					defer stmt.Close()

					for _, Channel := range Channels {
						msg, err := Bot.ChannelMessageSendEmbed(Channel, NewVtuber{
							Group: database.GroupName{
								ID:        GroupID,
								NameGroup: Group.GroupName,
								IconURL:   Group.GroupIcon,
							},
							Member: VtuberMember,
						}.SendNotif())
						if err != nil {
							log.Error(msg, err)
						}
					}
					NewVtuberNames = append(NewVtuberNames, "`"+VtuberMember.Name+"`")
					VtuberMember.InputSubs(id)

				} else {
					log.WithFields(log.Fields{
						"VtuberGroup": Group.GroupName,
						"Vtuber":      VtuberMember.ENName,
					}).Info("already added...")
					_, err := db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?  Where id=?`, VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.Region, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, VtuberMember.BliBiliFace(), VtuberMember.YtAvatar(), VtuberMember.TwitterName, tmp)
					if err != nil {
						log.Error(err)
					}
					log.WithFields(log.Fields{
						"VtuberGroup": Group.GroupName,
						"Vtuber":      VtuberMember.ENName,
					}).Info("update member")
				}
				log.WithFields(log.Fields{
					"VtuberGroup": Group.GroupName,
					"Vtuber":      VtuberMember.ENName,
				}).Info("Add subs info to database")

				/*
					Add subs info
				*/
				VtuberMember.InputSubs(tmp)
			}(Group.Members[j], wg2)
		}
		wg2.Wait()
		if NewVtuberNames != nil {
			for _, Channel := range GetChannelByGroup(GroupID) {
				msg, err := Bot.ChannelMessageSend(Channel, "New Update!!!! @here "+strings.Join(NewVtuberNames, ","))
				if err != nil {
					log.Error(msg, err)
				}
			}
		}
	}
}

//Get discord channel id from VtuberGroup
func GetChannelByGroup(VtuberGroupID int64) []string {
	var channellist []string
	rows, err := db.Query(`SELECT DiscordChannelID FROM Channel WHERE VtuberGroup_id=?`, VtuberGroupID)
	if err != nil {
		log.Error(err)
	}

	defer rows.Close()
	for rows.Next() {
		var list string
		err = rows.Scan(&list)
		if err != nil {
			log.Error(err)
		}
		channellist = append(channellist, list)
	}
	return channellist
}

func (Data Member) InputSubs(MemberID int64) {
	var tmp int64
	row := db.QueryRow("SELECT id FROM Subscriber WHERE VtuberMember_id=? ", MemberID)
	err := row.Scan(&tmp)
	Subs := Data.GetYtSubs()
	Bili := Data.GetBiliFolow()

	ytsubs, _ := strconv.Atoi(Subs.Items[0].Statistics.SubscriberCount)
	ytvideos, _ := strconv.Atoi(Subs.Items[0].Statistics.VideoCount)
	ytviews, _ := strconv.Atoi(Subs.Items[0].Statistics.ViewCount)
	bilifoll := Bili.Follow.Data.Follower
	bilivideos := Bili.Video
	biliview := Bili.Like.Data.Archive.View
	twfollo := Data.GetTwitterFollow()
	if err != nil || err == sql.ErrNoRows {
		stmt, err := db.Prepare("INSERT INTO Subscriber (Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers,VtuberMember_id) values(?,?,?,?,?,?,?,?)")
		engine.BruhMoment(err, "", false)
		res, err := stmt.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, MemberID)
		engine.BruhMoment(err, "", false)
		_, err = res.LastInsertId()
		engine.BruhMoment(err, "", false)

		defer stmt.Close()
	} else {
		rows, err := db.Query(`SELECT Youtube_Subscriber,Youtube_Videos,Youtube_Views,BiliBili_Followers,BiliBili_Videos,BiliBili_Views,Twitter_Followers FROM Subscriber WHERE VtuberMember_id=?`, MemberID)
		engine.BruhMoment(err, "", false)
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
		engine.BruhMoment(err, "", false)
		Update.Exec(ytsubs, ytvideos, ytviews, bilifoll, bilivideos, biliview, twfollo, tmp)
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
		_, err := db.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, Viewers=? Where RoomID=? AND VtuberMember_id=?`, Data["Status"], Data["Title"], Data["Thumbnail"], Data["Description"], Data["PublishedAt"], Data["ScheduledStart"], Data["Online"], Data["LiveRoomID"], Data["MemberID"])
		engine.BruhMoment(err, "", false)
		return false
	}
}
