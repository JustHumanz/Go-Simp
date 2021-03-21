-- MariaDB dump 10.19  Distrib 10.5.9-MariaDB, for Linux (x86_64)
--
-- Host: humanz-db.c1guxbt3kfqt.ap-southeast-1.rds.amazonaws.com    Database: Vtuber
-- ------------------------------------------------------
-- Server version	10.4.13-MariaDB-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `BiliBili`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS IF NOT EXISTS `BiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VideoID` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Type` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `UploadDate` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Viewers` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Length` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Channel`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `DiscordChannelID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Type` int(11) NOT NULL,
  `LiveOnly` tinyint(4) NOT NULL DEFAULT 0,
  `NewUpcoming` tinyint(4) NOT NULL DEFAULT 0,
  `Dynamic` tinyint(4) NOT NULL DEFAULT 0,
  `Region` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Lite` tinyint(4) NOT NULL DEFAULT 0,
  `IndieNotif` tinyint(4) NOT NULL DEFAULT 0,
  `VtuberGroup_id` tinyint(4) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `LiveBiliBili`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `LiveBiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `RoomID` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Status` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `Published` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `ScheduledStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Viewers` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `EndStream` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Subscriber`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Subscriber` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Youtube_Subscriber` int(11) DEFAULT NULL,
  `Youtube_Videos` int(11) DEFAULT NULL,
  `Youtube_Views` int(11) DEFAULT NULL,
  `BiliBili_Followers` int(11) DEFAULT NULL,
  `BiliBili_Videos` int(11) DEFAULT NULL,
  `BiliBili_Views` int(11) DEFAULT NULL,
  `Twitter_Followers` int(11) DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `TBiliBili`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `TBiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PermanentURL` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Likes` int(11) DEFAULT NULL,
  `Photos` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Videos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Text` text COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Dynamic_id` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Twitch`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Twitch` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `Game` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Status` varchar(24) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ScheduledStart` timestamp NOT NULL DEFAULT current_timestamp(),
  `Viewers` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Twitter`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Twitter` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PermanentURL` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Likes` int(11) DEFAULT NULL,
  `Photos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Videos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Text` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TweetID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Pixiv`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Pixiv` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PermanentURL` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Photos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Text` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PixivID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `Lewd`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Lewd` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PermanentURL` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Photos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Videos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Text` varchar(1024) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `TweetID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `PixivID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `User`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `User` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `DiscordID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `DiscordUserName` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Human` tinyint(4) DEFAULT 1,
  `Reminder` int(2) DEFAULT 0,
  `VtuberMember_id` int(11) NOT NULL,
  `Channel_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VtuberGroup`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `VtuberGroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VtuberGroupName` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberGroupIcon` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VtuberMember`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `VtuberMember` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VtuberName` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberName_EN` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberName_JP` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitter_Hashtag` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitter_Lewd` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,    
  `BiliBili_Hashtag` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Youtube_ID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Youtube_Avatar` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BiliBili_SpaceID` int(11) NOT NULL,
  `BiliBili_RoomID` int(11) NOT NULL,
  `BiliBili_Avatar` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitter_Username` varchar(16) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitch_Username` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitch_Avatar` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Region` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Fanbase` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberGroup_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Youtube`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Youtube` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VideoID` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Type` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Status` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `PublishedAt` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `ScheduledStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `EndStream` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Viewers` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Length` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping routines for database 'Vtuber'
--
/*!50003 DROP PROCEDURE IF EXISTS `GetArt` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE PROCEDURE `GetArt`(
			GroupID int,
			MemberID int,
			State varchar(11)
		)
BEGIN
		IF State = 'twitter' THEN
			SELECT Twitter.* FROM Vtuber.Twitter 
			Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id 
			Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id 
			where (VtuberGroup.id=GroupID OR VtuberMember.id=MemberID)  ORDER by RAND() LIMIT 1;
		else
			SELECT TBiliBili.* FROM Vtuber.TBiliBili  
			Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id 
			Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id 
			where (VtuberGroup.id=GroupID OR VtuberMember.id=MemberID)  ORDER by RAND() LIMIT 1;
			
		end if;
		END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetLiveBiliBili` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE PROCEDURE `GetLiveBiliBili`(
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
		END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetSpaceBiliBili` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE PROCEDURE `GetSpaceBiliBili`(
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
		END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetVtuberName` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE PROCEDURE `GetVtuberName`(
			GroupID int
		)
BEGIN
			SELECT VtuberMember.* FROM Vtuber.VtuberMember WHERE VtuberGroup_id=GroupID 
			Order by Region,VtuberGroup_id;
		END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetYt` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;
CREATE PROCEDURE `GetYt`(
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
	END ;;
DELIMITER ;
/*!50003 SET sql_mode              = @saved_sql_mode */ ;
/*!50003 SET character_set_client  = @saved_cs_client */ ;
/*!50003 SET character_set_results = @saved_cs_results */ ;
/*!50003 SET collation_connection  = @saved_col_connection */ ;
/*!50003 DROP PROCEDURE IF EXISTS `GetYtByReg` */;
/*!50003 SET @saved_cs_client      = @@character_set_client */ ;
/*!50003 SET @saved_cs_results     = @@character_set_results */ ;
/*!50003 SET @saved_col_connection = @@collation_connection */ ;
/*!50003 SET character_set_client  = utf8mb4 */ ;
/*!50003 SET character_set_results = utf8mb4 */ ;
/*!50003 SET collation_connection  = utf8mb4_general_ci */ ;
/*!50003 SET @saved_sql_mode       = @@sql_mode */ ;
/*!50003 SET sql_mode              = 'STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION' */ ;
DELIMITER ;;

-- Dump completed on 2021-03-02 14:31:03
