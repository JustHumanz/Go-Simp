package database

import (
	"encoding/json"
	"regexp"
	"time"
)

//Group group struct
type Group struct {
	ID              int64
	IconURL         string
	GroupName       string
	YoutubeChannels []GroupYtChannel
	Members         []Member
}

type GroupYtChannel struct {
	ID        int64
	YtChannel string
	Region    string
}

func (Data *Group) RemoveNillIconURL() *Group {
	if match, _ := regexp.MatchString("404.jpg", Data.IconURL); match {
		Data.IconURL = ""
	}
	return Data
}

//IsNull check if group struct is nil
func (Data Group) IsNull() bool {
	return Data.ID == 0
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
	TwitterLewd      string
	BiliBiliHashtags string
	BiliBiliAvatar   string
	TwitchName       string
	TwitchAvatar     string
	Region           string
	Fanbase          string
	Status           string
	GroupID          int64
}

func (Data Member) IsMemberNill() bool {
	return Data == (Member{})
}

func (Data Member) Active() bool {
	return Data.Status == "Active"
}

func (Data Member) IsYtNill() bool {
	return Data.YoutubeID != ""
}

func (Data Member) IsTwitchNill() bool {
	return Data.TwitchAvatar != ""
}

func (Data Member) IsBiliNill() bool {
	return Data.BiliRoomID != 0
}

func (Data Member) IsTwNill() bool {
	return Data.TwitterName != ""
}

type LiveStream struct {
	ID           int64
	Status       string
	VideoID      string
	Title        string
	Thumb        string
	Desc         string
	Schedul      time.Time
	End          time.Time
	Published    time.Time
	Game         string
	Type         string
	Viewers      string
	Length       string
	Member       Member
	Group        Group
	GroupYoutube GroupYtChannel
	State        string
	IsBiliLive   bool
}

func (Data *LiveStream) SetGroupYt(new GroupYtChannel) *LiveStream {
	Data.GroupYoutube = new
	return Data
}

func (Data *LiveStream) AddVideoID(new string) *LiveStream {
	Data.VideoID = new
	return Data
}

func (Data *LiveStream) SetType(new string) *LiveStream {
	Data.Type = new
	return Data
}

func (Data *LiveStream) SetState(new string) *LiveStream {
	Data.State = new
	return Data
}

func (Data *LiveStream) AddMember(new Member) *LiveStream {
	Data.Member = new
	return Data
}

func (Data *LiveStream) AddGroup(new Group) *LiveStream {
	Data.Group = new
	return Data
}

func (Data *LiveStream) UpdateStatus(new string) *LiveStream {
	Data.Status = new
	return Data
}

func (Data *LiveStream) UpdateSchdule(new time.Time) *LiveStream {
	Data.Schedul = new
	return Data
}

func (Data *LiveStream) UpdateViewers(new string) *LiveStream {
	Data.Viewers = new
	return Data
}

func (Data *LiveStream) UpdateThumbnail(new string) *LiveStream {
	Data.Thumb = new
	return Data
}

func (Data *LiveStream) UpdateTitle(new string) *LiveStream {
	Data.Title = new
	return Data
}

func (Data *LiveStream) UpdateEnd(new time.Time) *LiveStream {
	Data.End = new
	return Data
}

func (Data *LiveStream) UpdateLength(new string) *LiveStream {
	Data.Length = new
	return Data
}

func (Data *LiveStream) UpdatePublished(new time.Time) *LiveStream {
	Data.Published = new
	return Data
}

func (Data *LiveStream) UpdateGame(new string) *LiveStream {
	Data.Game = new
	return Data
}

func (Data *LiveStream) UpdateDesc(new string) *LiveStream {
	Data.Desc = new
	return Data
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
	YtSubs, YtVideos, YtViews               int
	BiliFollow, BiliVideos, BiliViews       int
	ID, TwFollow, TwitchFollow, TwitchViews int
	Member                                  Member
	Group                                   Group
	State                                   string
}

func (Data *MemberSubs) SetMember(new Member) *MemberSubs {
	Data.Member = new
	return Data
}

func (Data *MemberSubs) SetGroup(new Group) *MemberSubs {
	Data.Group = new
	return Data
}

func (Data *MemberSubs) UpdateState(new string) *MemberSubs {
	Data.State = new
	return Data
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
func (ac LiveStream) MarshalBinary() ([]byte, error) {
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

func (Data *DiscordChannel) IsFanart() bool {
	if Data.TypeTag == 1 || Data.TypeTag == 3 || Data.TypeTag == 70 {
		return true
	}
	return false
}

func (Data *DiscordChannel) IsLive() bool {
	if Data.TypeTag == 2 || Data.TypeTag == 3 {
		return true
	}
	return false
}

func (Data *DiscordChannel) IsLewd() bool {
	if Data.TypeTag == 69 || Data.TypeTag == 70 {
		return true
	}
	return false
}
