-- MariaDB dump 10.17  Distrib 10.4.13-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: Vtuber
-- ------------------------------------------------------
-- Server version	10.4.13-MariaDB

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
-- Table structure for table `Channel`
--
CREATE DATABASE IF NOT EXISTS `Vtuber`;

USE `Vtuber`;

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `Channel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `DiscordChannelID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Type` int(11) NOT NULL,
  `VtuberGroup_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `Human` TINYINT DEFAULT 1,
  `VtuberMember_id` int(11) NOT NULL,
  `Channel_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `T.Bilibili`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `TBiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `PermanentURL` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Author` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Likes` int(11) DEFAULT NULL,
  `Photos` TEXT COLLATE utf8mb4_unicode_ci DEFAULT NULL,  /*i'm not joking,they use sha1 hash for image identify,so the url very fucking long*/
  `Videos` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Text` TEXT COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Dynamic_id` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
--
-- Table structure for table `VtuberGroup`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `VtuberGroup` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VtuberGroupName` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberGroupIcon` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VtuberMember`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE IF NOT EXISTS `VtuberMember` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VtuberName` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberName_EN` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberName_JP` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Hashtag` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BiliBili_Hashtag` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Youtube_ID` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Youtube_Avatar` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `BiliBili_SpaceID` INT(11) DEFAULT NULL,
  `BiliBili_RoomID` INT(11) DEFAULT NULL,
  `BiliBili_Avatar` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Twitter_Username` varchar(24) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `Region` varchar(5) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `VtuberGroup_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
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
  `Length` varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

CREATE TABLE IF NOT EXISTS `BiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `VideoID` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Type` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `UploadDate` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Viewers` int(11) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Length` varchar(11) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


CREATE TABLE IF NOT EXISTS `LiveBiliBili` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `RoomID` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Status` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Title` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Thumbnails` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `Description` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `Published` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `ScheduledStart` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `Viewers` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

CREATE TABLE IF NOT EXISTS `Vtuber`.`Subscriber` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `Youtube_Subs` INT(11) NULL,
  `Youtube_Videos` INT(11) NULL,
  `Youtube_Views` INT(11) NULL,
  `BiliBili_Follows` INT(11) NULL,
  `BiliBili_Videos` INT(11) NULL,
  `BiliBili_Views` INT(11) NULL,
  `Twitter_Follows` INT(11) NULL,
  `VtuberMember_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
  );

-- Dump completed on 2020-07-12 18:51:30