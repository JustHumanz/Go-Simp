package database

import (
	"time"
)

type MemberGroupID struct {
	EnName          string
	JpName          string
	TwitterHashtags string
	MemberID        int64
	GroupID         int64
	GroupName       string
	GroupIcon       string
}

type DataFanart struct {
	EnName       string
	JpName       string
	PermanentURL string
	Author       string
	Photos       []string
	Videos       string
	Text         string
	Likes        int
	Dynamic_id   string
	State        string
}

type GroupName struct {
	ID        int64
	IconURL   string
	NameGroup string
}

type Name struct {
	ID               int64
	Name             string
	EnName           string
	JpName           string
	YoutubeID        string
	YoutubeAvatar    string
	BiliBiliID       int
	BiliRoomID       int
	TwitterHashtags  string
	TwitterName      string
	BiliBiliHashtags string
	BiliBiliAvatar   string
	Region           string
}

type YtDbData struct {
	ID            int
	ChannelID     string
	Group         string
	Status        string
	NameEN        string
	NameJP        string
	VideoID       string
	Title         string
	Thumb         string
	Desc          string
	YoutubeAvatar string
	Schedul       time.Time
	End           time.Time
	Published     time.Time
	Type          string
	Region        string
	Viewers       string
	Length        string
}

type UserStruct struct {
	DiscordID, DiscordUserName, Channel_ID string
	GroupID                                int64
}

type LiveBiliDB struct {
	LiveRoomID, Online, ID                                        int
	Status, Title, Thumbnail, Description, EnName, JpName, Avatar string
	ScheduledStart, PublishedAt                                   time.Time
}

type SpaceBiliDB struct {
	Viewers                                                                      int
	VideoID, Title, Thumbnail, Description, EnName, JpName, Avatar, Type, Length string
	UploadDate                                                                   time.Time
}

type InputTW struct {
	Url      string
	Author   string
	Like     int
	Photos   string
	Video    string
	Text     string
	TweetID  string
	MemberID int64
}

type MemberSubs struct {
	YtSubs, YtVideos, YtViews         int
	BiliFollow, BiliVideos, BiliViews int
	id, TwFollow                      int
	MemberID                          int64
}

type InputBiliBili struct {
	VideoID  string
	Type     string
	Title    string
	Thum     string
	Desc     string
	Update   time.Time
	Viewers  int
	MemberID int64
	Length   string
}

type InputTBiliBili struct {
	URL        string
	Author     string
	Avatar     string
	Like       int
	Photos     string
	Videos     string
	Text       string
	Dynamic_id string
}
