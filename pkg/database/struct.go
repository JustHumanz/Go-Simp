package database

import (
	"encoding/json"
	"time"
)

/*
type MemberGroupID struct {
	EnName          string
	JpName          string
	TwitterHashtags string
	MemberID        int64
	GroupID         int64
	GroupName       string
	GroupIcon       string
}
*/

//DataFanart fanart struct
type DataFanart struct {
	ID           int64
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

//Group group struct
type Group struct {
	ID        int64
	IconURL   string
	GroupName string
	Members   []Member
}

//IsNull check if group struct is nil
func (Data Group) IsNull() bool {
	if Data.ID == 0 {
		return true
	}
	return false
}

//VtubersPayload payload with multiple group
type VtubersPayload struct {
	VtuberData []Group
}

//Member Member struct
type Member struct {
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
	TwitchName       string
	Region           string
	GroupID          int64
}

//YtDbData Youtube database struct
type YtDbData struct {
	ID            int64
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
	MemberID      int64
	GroupID       int64
}

//UserStruct user struct
type UserStruct struct {
	DiscordID, DiscordUserName, Channel_ID string
	Group                                  Group
	Member                                 Member
	Reminder                               int
	Human                                  bool
}

//SetDiscordID set UserStruct ID
func (Data *UserStruct) SetDiscordID(new string) *UserStruct {
	Data.DiscordID = new
	return Data
}

//SetDiscordUserName set UserStruct Discord user name
func (Data *UserStruct) SetDiscordUserName(new string) *UserStruct {
	Data.DiscordUserName = new
	return Data
}

//SetDiscordChannelID set UserStruct Set DiscordChannel ID
func (Data *UserStruct) SetDiscordChannelID(new string) *UserStruct {
	Data.Channel_ID = new
	return Data
}

//SetGroup set UserStruct group
func (Data *UserStruct) SetGroup(new Group) *UserStruct {
	Data.Group = new
	return Data
}

//SetHuman set UserStruct human or roles
func (Data *UserStruct) SetHuman(new bool) *UserStruct {
	Data.Human = new
	return Data
}

//SetReminder set UserStruct reminder time
func (Data *UserStruct) SetReminder(new int) *UserStruct {
	Data.Reminder = new
	return Data
}

//SetMember set UserStruct member
func (Data *UserStruct) SetMember(new Member) *UserStruct {
	Data.Member = new
	return Data
}

//LiveBiliDB live bilibili database struct
type LiveBiliDB struct {
	LiveRoomID, Online, ID                                        int
	Status, Title, Thumbnail, Description, EnName, JpName, Avatar string
	ScheduledStart, PublishedAt                                   time.Time
}

//SpaceBiliDB spacebilibili database struct
type SpaceBiliDB struct {
	Viewers                                                                      int
	VideoID, Title, Thumbnail, Description, EnName, JpName, Avatar, Type, Length string
	UploadDate                                                                   time.Time
}

//InputTW twitter fanart struct
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

//MemberSubs subscribe struct
type MemberSubs struct {
	YtSubs, YtVideos, YtViews         int
	BiliFollow, BiliVideos, BiliViews int
	ID, TwFollow                      int
	MemberID                          int64
}

//InputBiliBili input bilibili struct
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

//TBiliBili tbilibili struct
type TBiliBili struct {
	URL        string
	Author     string
	Avatar     string
	Like       int
	Photos     []string
	Videos     string
	Text       string
	Dynamic_id string
	Member     Member
	Group      Group
}

//Guild guild struct
type Guild struct {
	ID   string
	Name string
	Join time.Time
}

//MarshalBinary change struct to binary
func (ac MemberSubs) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

//MarshalBinary change struct to binary
func (ac YtDbData) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

//MarshalBinary change struct to binary
func (ac Member) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

//MarshalBinary change struct to binary
func (ac DiscordChannel) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

//MarshalBinary change struct to binary
func (ac UserStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

//DiscordChannel discord channel struct
type DiscordChannel struct {
	ID             int64
	ChannelID      string
	TypeTag        int
	LiveOnly       bool
	NewUpcoming    bool
	Dynamic        bool
	LiteMode       bool
	IndieNotif     bool
	Group          Group
	Member         Member
	VideoID        string
	EmbedMessageID string
	TextMessageID  string
	Region         string
}

func (Data *DiscordChannel) SetChannel(new string) *DiscordChannel {
	Data.ChannelID = new
	return Data
}

func (Data *DiscordChannel) SetTypeTag(new int) *DiscordChannel {
	Data.TypeTag = new
	return Data
}

func (Data *DiscordChannel) SetLiveOnly(new bool) *DiscordChannel {
	Data.LiveOnly = new
	return Data
}

func (Data *DiscordChannel) SetNewUpcoming(new bool) *DiscordChannel {
	Data.NewUpcoming = new
	return Data
}

func (Data *DiscordChannel) SetLite(new bool) *DiscordChannel {
	Data.LiteMode = new
	return Data
}

func (Data *DiscordChannel) SetIndieNotif(new bool) *DiscordChannel {
	Data.IndieNotif = new
	return Data
}

func (Data *DiscordChannel) SetVtuberGroupID(new int64) *DiscordChannel {
	Data.Group.ID = new
	return Data
}

func (Data *DiscordChannel) SetDynamic(new bool) *DiscordChannel {
	Data.Dynamic = new
	return Data
}

func (Data *DiscordChannel) SetVideoID(new string) *DiscordChannel {
	Data.VideoID = new
	return Data
}

func (Data *DiscordChannel) SetMsgEmbedID(new string) *DiscordChannel {
	Data.EmbedMessageID = new
	return Data
}

func (Data *DiscordChannel) SetMsgTextID(new string) *DiscordChannel {
	Data.TextMessageID = new
	return Data
}

func (Data *DiscordChannel) SetMember(new Member) *DiscordChannel {
	Data.Member = new
	return Data
}

func (Data *DiscordChannel) SetGroup(new Group) *DiscordChannel {
	Data.Group = new
	return Data
}

type TwitchDB struct {
	ID             int64
	Game           string
	Status         string
	Title          string
	Thumbnails     string
	ScheduledStart time.Time
	Viewers        int
}

func (Data *TwitchDB) UpdateViewers(new int) *TwitchDB {
	Data.Viewers = new
	return Data
}

func (Data *TwitchDB) UpdateStatus(new string) *TwitchDB {
	Data.Status = new
	return Data
}

func (Data *TwitchDB) UpdateTitle(new string) *TwitchDB {
	Data.Title = new
	return Data
}

func (Data *TwitchDB) UpdateThumbnails(new string) *TwitchDB {
	Data.Thumbnails = new
	return Data
}

func (Data *TwitchDB) UpdateSchedule(new time.Time) *TwitchDB {
	Data.ScheduledStart = new
	return Data
}

func (Data *TwitchDB) UpdateGame(new string) *TwitchDB {
	Data.Game = new
	return Data
}
