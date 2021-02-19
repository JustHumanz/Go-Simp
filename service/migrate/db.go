package main

import (
	"database/sql"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
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
		LiveOnly TINYINT NOT NULL DEFAULT 0,
		NewUpcoming TINYINT NOT NULL DEFAULT 1,
		Dynamic TINYINT NOT NULL DEFAULT 0,
		Region VARCHAR(45) NOT NULL ,
		Lite TINYINT NOT NULL DEFAULT 0,
		VtuberGroup_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GuildList (
		id	int(11) NOT NULL AUTO_INCREMENT,
		GuildID	varchar(256),
		GuildName	varchar(256),
		JoinDate	timestamp,
		PRIMARY KEY(id)
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
		Twitch_Username varchar(24) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
		Twitch_Avatar varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Twitch (
		id int(11) NOT NULL AUTO_INCREMENT,
		Game varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Status varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
		Title varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		Thumbnails varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		ScheduledStart timestamp NOT NULL DEFAULT current_timestamp(),
		Viewers varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
		VtuberMember_id int(11) NOT NULL,
		PRIMARY KEY (id)
	  ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Subscriber (
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
		sts varchar(11),
		reg  varchar(11)
	)
	BEGIN
		IF reg != '' THEN
				IF sts = 'upcoming' THEN
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where VtuberGroup.id=grpid AND Status='upcoming' AND Region=reg Order by ScheduledStart DESC Limit 3;

			ELSEIF sts = 'live' OR sts = 'private' THEN 
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where VtuberGroup.id=grpid AND Status=sts AND Region=reg Limit 3;
			ELSE 
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where VtuberGroup.id=grpid AND Status='past' AND Region=reg AND EndStream !='' order by EndStream DESC Limit 3;
				
			END if;	
		ELSE
			IF sts = 'upcoming' THEN
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
				AND Status='upcoming' 
				Order by ScheduledStart DESC Limit lmt;
			ELSEIF sts = 'live' OR sts = 'private' THEN 
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
				AND Status=sts
				Limit lmt;
			ELSE 
				SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Type,
				Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers,VtuberMember.id,VtuberGroup.id 
				FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id 
				Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id 
				Where (VtuberGroup.id=grpid or VtuberMember.id=memid) 
				AND Status='past'
				AND EndStream !='' order by EndStream ASC Limit lmt;
				
			END if;	
		END if;
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
			Region,Hashtag,BiliBili_Hashtag,BiliBili_Avatar,Twitter_Username,Twitch_Username,Youtube_Avatar 
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
			AND Status=Sts Order by ScheduledStart DESC Limit lmt;
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
			SELECT Twitter.id,VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Vtuber.Twitter 
			Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id 
			Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id 
			where (VtuberGroup.id=GroupID OR VtuberMember.id=MemberID)  ORDER by RAND() LIMIT 1;
		else
			SELECT TBiliBili.id,VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Vtuber.TBiliBili  
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

func AddData(Data Vtuber) {
	Independen := func(wg *sync.WaitGroup) {
		var (
			GroupData = database.Group{
				ID:        0,
				GroupName: "Independen",
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

		for _, VtuberMember := range Data.Vtuber.Independen.Members {
			DiscordChannel := GroupData.GetChannelByGroup(VtuberMember.Region)
			/*
				Add Member
			*/
			var MemberID int64
			row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND VtuberName_EN=? AND (Youtube_ID=? OR  BiliBili_SpaceID=? OR BiliBili_RoomID=?)", VtuberMember.Name, VtuberMember.ENName, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID)
			err := row.Scan(&MemberID)
			if err == sql.ErrNoRows {
				stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username,Twitch_Username,Twitch_Avatar) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
				if err != nil {
					log.Error(err)
				}
				BiliFace, err := VtuberMember.BliBiliFace()
				if err != nil {
					log.Error(err)
				}

				TwitchAvatar, err := VtuberMember.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}
				res, err := stmt.Exec(VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.YtID, VtuberMember.YtAvatar(), GroupData.ID, VtuberMember.Region, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, BiliFace, VtuberMember.TwitterName, VtuberMember.TwitchName, TwitchAvatar)
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

				NewVtuberNamesIndependen = append(NewVtuberNamesIndependen, "`"+VtuberMember.Name+"`")
				VtuberMember.InputSubs(MemberID)
				//New.SendNotif(Bot)
			} else if err != nil {
				log.Error(err)
			} else {
				TwitchAvatar, err := VtuberMember.GetTwitchAvatar()
				if err != nil {
					log.Error(err)
				}
				BiliFace, err := VtuberMember.BliBiliFace()
				if err != nil {
					log.Error(err)
				} else {
					log.WithFields(log.Fields{
						"VtuberGroup": "Independen",
						"Vtuber":      VtuberMember.ENName,
					}).Info("Update member")
					_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?,Twitch_Username=?,Twitch_Avatar=?  Where id=?`, VtuberMember.Name, VtuberMember.ENName, VtuberMember.JPName, VtuberMember.Hashtag.Twitter, VtuberMember.Hashtag.BiliBili, VtuberMember.Region, VtuberMember.YtID, VtuberMember.BiliBiliID, VtuberMember.BiliRoomID, BiliFace, VtuberMember.YtAvatar(), VtuberMember.TwitterName, VtuberMember.TwitchName, TwitchAvatar, MemberID)
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
			}
		}
	}

	Agency := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for _, GroupRaw := range Data.Vtuber.Group {
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

			DiscordChannel := GroupData.GetChannelByGroup("")
			for _, v := range GroupRaw.Members {
				/*
					Add Member
				*/
				var MemberID int64
				row := db.QueryRow("SELECT id FROM VtuberMember WHERE VtuberName=? AND (Youtube_ID=? OR BiliBili_SpaceID=? OR BiliBili_RoomID=?)", v.Name, v.YtID, v.BiliBiliID, v.BiliRoomID)
				err := row.Scan(&MemberID)
				if err == sql.ErrNoRows {
					stmt, err := db.Prepare("INSERT INTO VtuberMember (VtuberName,VtuberName_EN,VtuberName_JP,Hashtag,BiliBili_Hashtag,Youtube_ID,Youtube_Avatar,VtuberGroup_id,Region,BiliBili_SpaceID,BiliBili_RoomID,BiliBili_Avatar,Twitter_Username,Twitch_Username,Twitch_Avatar) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
					if err != nil {
						log.Error(err)
					}

					BiliFace, err := v.BliBiliFace()
					if err != nil {
						log.Error(err)
					}

					TwitchAvatar, err := v.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}

					res, err := stmt.Exec(v.Name, v.ENName, v.JPName, v.Hashtag.Twitter, v.Hashtag.BiliBili, v.YtID, v.YtAvatar(), GroupData.ID, v.Region, v.BiliBiliID, v.BiliRoomID, BiliFace, v.TwitterName, v.TwitchName, TwitchAvatar)
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
					TwitchAvatar, err := v.GetTwitchAvatar()
					if err != nil {
						log.Error(err)
					}

					BiliFace, err := v.BliBiliFace()
					if err != nil {
						log.Error(err)
					} else {
						_, err = db.Exec(`Update VtuberMember set VtuberName=?, VtuberName_EN=?, VtuberName_JP=? ,Hashtag=? ,BiliBili_Hashtag=? ,Region=? ,Youtube_ID=? ,BiliBili_SpaceID=?,BiliBili_RoomID=?, BiliBili_Avatar=? ,Youtube_Avatar=?, Twitter_Username=?,Twitch_Username=?,Twitch_Avatar=? Where id=?`, v.Name, v.ENName, v.JPName, v.Hashtag.Twitter, v.Hashtag.BiliBili, v.Region, v.YtID, v.BiliBiliID, v.BiliRoomID, BiliFace, v.YtAvatar(), v.TwitterName, v.TwitchName, TwitchAvatar, MemberID)
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
	twfollo, twerr := Data.GetTwitterFollow()
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

func GetHastagMember(MemberID int64) string {
	rows, err := db.Query(`SELECT Hashtag FROM Vtuber.VtuberMember where id=?`, MemberID)
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

		stmt, err := db.Prepare("INSERT INTO Twitch (Game,Status,Title,Thumbnails,ScheduledStart,Viewers,VtuberMember_id) values(?,?,?,?,?,?,?)")
		if err != nil {
			return err
		}
		res, err := stmt.Exec(Data["Game"], Data["Status"], Data["Title"], Data["Thumbnails"], Data["ScheduledStart"], Data["Viewers"], Data["MemberID"])
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
