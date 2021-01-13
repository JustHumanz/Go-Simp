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

type Group struct {
	ID        int64
	IconURL   string
	GroupName string
}

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
	TwitchUserName   string
	Region           string
	GroupID          int64
}

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

type UserStruct struct {
	DiscordID, DiscordUserName, Channel_ID string
	Group                                  Group
	Member                                 Member
	Reminder                               int
	Human                                  bool
}

func (Data *UserStruct) SetDiscordID(new string) *UserStruct {
	Data.DiscordID = new
	return Data
}

func (Data *UserStruct) SetDiscordUserName(new string) *UserStruct {
	Data.DiscordUserName = new
	return Data
}

func (Data *UserStruct) SetDiscordChannelID(new string) *UserStruct {
	Data.Channel_ID = new
	return Data
}

func (Data *UserStruct) SetGroup(new Group) *UserStruct {
	Data.Group = new
	return Data
}

func (Data *UserStruct) SetHuman(new bool) *UserStruct {
	Data.Human = new
	return Data
}

func (Data *UserStruct) SetReminder(new int) *UserStruct {
	Data.Reminder = new
	return Data
}

func (Data *UserStruct) SetMember(new Member) *UserStruct {
	Data.Member = new
	return Data
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
	ID, TwFollow                      int
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

type TwitterUser struct {
	ID              int         `json:"id"`
	IDStr           string      `json:"id_str"`
	Name            string      `json:"name"`
	ScreenName      string      `json:"screen_name"`
	Location        string      `json:"location"`
	ProfileLocation interface{} `json:"profile_location"`
	Description     string      `json:"description"`
	URL             string      `json:"url"`
	Entities        struct {
		URL struct {
			Urls []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"url"`
		Description struct {
			Urls []struct {
				URL         string `json:"url"`
				ExpandedURL string `json:"expanded_url"`
				DisplayURL  string `json:"display_url"`
				Indices     []int  `json:"indices"`
			} `json:"urls"`
		} `json:"description"`
	} `json:"entities"`
	Protected       bool        `json:"protected"`
	FollowersCount  int         `json:"followers_count"`
	FriendsCount    int         `json:"friends_count"`
	ListedCount     int         `json:"listed_count"`
	CreatedAt       string      `json:"created_at"`
	FavouritesCount int         `json:"favourites_count"`
	UtcOffset       interface{} `json:"utc_offset"`
	TimeZone        interface{} `json:"time_zone"`
	GeoEnabled      bool        `json:"geo_enabled"`
	Verified        bool        `json:"verified"`
	StatusesCount   int         `json:"statuses_count"`
	Lang            interface{} `json:"lang"`
	Status          struct {
		CreatedAt string `json:"created_at"`
		ID        int64  `json:"id"`
		IDStr     string `json:"id_str"`
		Text      string `json:"text"`
		Truncated bool   `json:"truncated"`
		Entities  struct {
			Hashtags     []interface{} `json:"hashtags"`
			Symbols      []interface{} `json:"symbols"`
			UserMentions []struct {
				ScreenName string `json:"screen_name"`
				Name       string `json:"name"`
				ID         int    `json:"id"`
				IDStr      string `json:"id_str"`
				Indices    []int  `json:"indices"`
			} `json:"user_mentions"`
			Urls []interface{} `json:"urls"`
		} `json:"entities"`
		Source               string      `json:"source"`
		InReplyToStatusID    int64       `json:"in_reply_to_status_id"`
		InReplyToStatusIDStr string      `json:"in_reply_to_status_id_str"`
		InReplyToUserID      int         `json:"in_reply_to_user_id"`
		InReplyToUserIDStr   string      `json:"in_reply_to_user_id_str"`
		InReplyToScreenName  string      `json:"in_reply_to_screen_name"`
		Geo                  interface{} `json:"geo"`
		Coordinates          interface{} `json:"coordinates"`
		Place                interface{} `json:"place"`
		Contributors         interface{} `json:"contributors"`
		IsQuoteStatus        bool        `json:"is_quote_status"`
		RetweetCount         int         `json:"retweet_count"`
		FavoriteCount        int         `json:"favorite_count"`
		Favorited            bool        `json:"favorited"`
		Retweeted            bool        `json:"retweeted"`
		Lang                 string      `json:"lang"`
	} `json:"status"`
	ContributorsEnabled            bool   `json:"contributors_enabled"`
	IsTranslator                   bool   `json:"is_translator"`
	IsTranslationEnabled           bool   `json:"is_translation_enabled"`
	ProfileBackgroundColor         string `json:"profile_background_color"`
	ProfileBackgroundImageURL      string `json:"profile_background_image_url"`
	ProfileBackgroundImageURLHTTPS string `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool   `json:"profile_background_tile"`
	ProfileImageURL                string `json:"profile_image_url"`
	ProfileImageURLHTTPS           string `json:"profile_image_url_https"`
	ProfileBannerURL               string `json:"profile_banner_url"`
	ProfileLinkColor               string `json:"profile_link_color"`
	ProfileSidebarBorderColor      string `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool   `json:"profile_use_background_image"`
	HasExtendedProfile             bool   `json:"has_extended_profile"`
	DefaultProfile                 bool   `json:"default_profile"`
	DefaultProfileImage            bool   `json:"default_profile_image"`
	Following                      bool   `json:"following"`
	FollowRequestSent              bool   `json:"follow_request_sent"`
	Notifications                  bool   `json:"notifications"`
	TranslatorType                 string `json:"translator_type"`
}

type Guild struct {
	ID   string
	Name string
	Join time.Time
}

func (ac MemberSubs) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac YtDbData) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac Member) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac DiscordChannel) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

func (ac UserStruct) MarshalBinary() ([]byte, error) {
	return json.Marshal(ac)
}

type DiscordChannel struct {
	ID             int64
	ChannelID      string
	TypeTag        int
	LiveOnly       bool
	NewUpcoming    bool
	Dynamic        bool
	Group          Group
	VideoID        string
	EmbedMessageID string
	TextMessageID  string
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
