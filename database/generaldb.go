package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	DB    *sql.DB
	debug bool
)

//Start Database session
func Start(dbsession *sql.DB) {
	DB = dbsession
	log.Info("Database module ready")
}

func BruhMoment(err error, msg string, exit bool) {
	if err != nil {
		log.Info(msg)
		log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Debugging(a, b, c interface{}) {
	if debug {
		f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err)
		}
		log.SetOutput(f)
		log.WithFields(log.Fields{
			"Func":   a,
			"Status": b,
			"Value":  c,
		}).Debug(c)
	}
}

//Get Twitter hashtag by group
func GetHashtag(GroupID int64) []MemberGroupID {
	funcvar := GetFunctionName(GetHashtag)
	Debugging(funcvar, "In", GroupID)
	rows, err := DB.Query(`SELECT VtuberMember.id,VtuberName,VtuberName_JP,VtuberGroup_id,Hashtag,VtuberGroupName,VtuberGroupIcon FROM VtuberMember INNER Join VtuberGroup ON VtuberGroup.id = VtuberMember.VtuberGroup_id WHERE VtuberGroup.id =?`, GroupID)
	BruhMoment(err, "", false)
	defer rows.Close()

	Data := []MemberGroupID{}
	for rows.Next() {
		var list MemberGroupID
		err = rows.Scan(&list.MemberID, &list.EnName, &list.JpName, &list.GroupID, &list.TwitterHashtags, &list.GroupName, &list.GroupIcon)
		BruhMoment(err, "", false)

		Data = append(Data, list)

	}
	Debugging(funcvar, "Out", Data)
	return Data
}

//Get all vtuber groupData
func GetGroup() []GroupName {
	funcvar := GetFunctionName(GetGroup)
	rows, err := DB.Query(`SELECT id,VtuberGroupName,VtuberGroupIcon FROM VtuberGroup`)
	BruhMoment(err, "", false)
	defer rows.Close()

	var Data []GroupName
	for rows.Next() {
		var list GroupName
		err = rows.Scan(&list.ID, &list.NameGroup, &list.IconURL)
		BruhMoment(err, "", false)

		Data = append(Data, list)

	}
	Debugging(funcvar, "Out", Data)
	return Data
}

//Get data of Vtuber member
func GetName(GroupID int64) []Name {
	funcvar := GetFunctionName(GetName)
	Debugging(funcvar, "In", GroupID)
	rows, err := DB.Query(`SELECT id,VtuberName,VtuberName_EN,VtuberName_JP,Youtube_ID,BiliBili_SpaceID,BiliBili_RoomID,Region,Hashtag,BiliBili_Hashtag,BiliBili_Avatar,Twitter_Username,Youtube_Avatar FROM VtuberMember WHERE VtuberGroup_id=?`, GroupID)
	BruhMoment(err, "", false)
	defer rows.Close()

	var Data []Name
	for rows.Next() {
		var list Name
		err = rows.Scan(&list.ID, &list.Name, &list.EnName, &list.JpName, &list.YoutubeID, &list.BiliBiliID, &list.BiliRoomID, &list.Region, &list.TwitterHashtags, &list.BiliBiliHashtags, &list.BiliBiliAvatar, &list.TwitterName, &list.YoutubeAvatar)
		BruhMoment(err, "", false)
		Data = append(Data, list)

	}
	Debugging(funcvar, "Out", Data)
	return Data
}

//gacha is gacha
func gacha() bool {
	return rand.Float32() < 0.5
}

//Get subs,follow,view,like data from Subscriber
func (Member Name) GetSubsCount() MemberSubs {
	var Data MemberSubs
	rows, err := DB.Query(`SELECT * FROM Subscriber WHERE VtuberMember_id=?`, Member.ID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.YtSubs, &Data.YtVideos, &Data.YtViews, &Data.BiliFollow, &Data.BiliVideos, &Data.BiliViews, &Data.TwFollow, &Data.MemberID)
		BruhMoment(err, "", false)
	}
	return Data

}

//Update Subscriber data
func (Data MemberSubs) UpdateSubs(State string) {
	if State == "yt" {
		_, err := DB.Exec(`Update Subscriber set Youtube_Subs=?,Youtube_Videos=?,Youtube_Views=? Where id=? `, Data.YtSubs, Data.YtVideos, Data.YtViews, Data.ID)
		BruhMoment(err, "", false)
	} else if State == "bili" {
		_, err := DB.Exec(`Update Subscriber set BiliBili_Follows=?,BiliBili_Videos=?,BiliBili_Views=? Where id=? `, Data.BiliFollow, Data.BiliVideos, Data.BiliViews, Data.ID)
		BruhMoment(err, "", false)
	} else {
		_, err := DB.Exec(`Update Subscriber set Twitter_Follows=? Where id=? `, Data.TwFollow, Data.ID)
		BruhMoment(err, "", false)
	}
}

//Get Member fanart URL from TBiliBili and Twitter
func (Member Name) GetMemberURL() DataFanart {
	var (
		Data     DataFanart
		PhotoTmp string
	)

	if gacha() {
		rows, err := DB.Query(`SELECT PermanentURL,Author,Likes,Photos,Videos,Text FROM Twitter WHERE VtuberMember_id=? ORDER by RAND() LIMIT 1`, Member.ID)
		BruhMoment(err, "", false)
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Data.Videos, &Data.Text)
			BruhMoment(err, "", false)
		}
		Data.Photos = strings.Split(PhotoTmp, "\n")
		Data.State = "Twitter"
	} else {
		rows, err := DB.Query(`SELECT PermanentURL,Author,Likes,Photos,Text,Dynamic_id FROM TBiliBili WHERE VtuberMember_id=? ORDER by RAND() LIMIT 1`, Member.ID)
		BruhMoment(err, "", false)
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Data.Text, &Data.Dynamic_id)
			BruhMoment(err, "", false)

		}
		Data.Photos = strings.Split(PhotoTmp, "\n")
		Data.State = "TBiliBili"
	}
	return Data
}

//Get Group fanart URL from TBiliBili and Twitter
func (GroupData GroupName) GetGroupURL() DataFanart {
	rows, err := DB.Query(`SELECT VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Twitter Inner Join VtuberMember on VtuberMember.id = Twitter.VtuberMember_id Inner Join VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where VtuberGroup.id=? ORDER by RAND() LIMIT 1`, GroupData.ID)
	BruhMoment(err, "", false)
	defer rows.Close()

	var (
		Data     DataFanart
		PhotoTmp string
	)
	for rows.Next() {
		err = rows.Scan(&Data.EnName, &Data.JpName, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Videos, &Data.Text)
		if err != nil {
			log.Error(err)
		}
	}
	Data.Photos = strings.Fields(PhotoTmp)
	return Data
}

//Input new fanart from twitter
func (Data InputTW) InputTwitter() {

	stmt, err := DB.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	BruhMoment(err, "", false)
	defer stmt.Close()

	res, err := stmt.Exec(Data.Url, Data.Author, Data.Like, Data.Photos, Data.Video, Data.Text, Data.TweetID, Data.MemberID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)
}

//Update Likes data
func (Data InputTW) UpdateTwitter(id int) {
	_, err := DB.Exec(`Update Twitter set Likes=? Where id=?`, Data.Like, id)
	BruhMoment(err, "", false)
}

//Get Channel id from Discord ChannelID and VtuberGroupID
func GetChannelID(DiscordChannelID string, GroupID int64) int {
	funcvar := GetFunctionName(GetChannelID)
	Debugging(funcvar, "In", fmt.Sprint(DiscordChannelID, GroupID))
	var ChannelID int
	err := DB.QueryRow("SELECT id from Channel where DiscordChannelID=? AND VtuberGroup_id=?", DiscordChannelID, GroupID).Scan(&ChannelID)
	BruhMoment(err, "", false)
	Debugging(funcvar, "Out", ChannelID)
	return ChannelID
}

//Add user form `tag me command`
func (Data UserStruct) Adduser(MemberID int64) error {
	ChannelID := GetChannelID(Data.Channel_ID, Data.GroupID)
	tmp := CheckUser(Data.DiscordID, MemberID, ChannelID)
	if tmp {
		return errors.New("Already registered")
	} else {
		stmt, err := DB.Prepare(`INSERT INTO User (DiscordID,DiscordUserName,VtuberMember_id,Channel_id) values(?,?,?,?)`)
		BruhMoment(err, "", false)
		defer stmt.Close()
		res, err := stmt.Exec(Data.DiscordID, Data.DiscordUserName, MemberID, ChannelID)
		BruhMoment(err, "", false)

		_, err = res.LastInsertId()
		BruhMoment(err, "", false)

		return nil
	}
}

//Delete user from `del` command
func (Data UserStruct) Deluser(MemberID int64) error {
	ChannelID := GetChannelID(Data.Channel_ID, Data.GroupID)
	tmp := CheckUser(Data.DiscordID, MemberID, ChannelID)
	if tmp {
		stmt, err := DB.Prepare(`DELETE From User WHERE DiscordID=? AND VtuberMember_id=? AND Channel_id=?`)
		BruhMoment(err, "", false)
		defer stmt.Close()

		stmt.Exec(Data.DiscordID, MemberID, ChannelID)
		return nil
	} else {
		return errors.New("Already removed")
	}
}

//Check user if already tagged
func CheckUser(DiscordID string, MemberID int64, Channel_ChannelID int) bool {
	funcvar := GetFunctionName(CheckUser)
	Debugging(funcvar, "In", fmt.Sprint(DiscordID, MemberID, Channel_ChannelID))
	var tmp int
	row := DB.QueryRow("SELECT id FROM User WHERE DiscordID=? AND VtuberMember_id=? AND Channel_id=?", DiscordID, MemberID, Channel_ChannelID)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		Debugging(funcvar, "Out", false)
		return false
	} else {
		Debugging(funcvar, "Out", true)
		return true
	}
}

//Add new discord channel from `enable` command
func AddChannel(ChannelID string, typetag int, VtuberGroupID int64) error {
	funcvar := GetFunctionName(AddChannel)
	Debugging(funcvar, "In", fmt.Sprint(ChannelID, VtuberGroupID))
	stmt, err := DB.Prepare(`INSERT INTO Channel (DiscordChannelID,Type,VtuberGroup_id) values(?,?,?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	res, err := stmt.Exec(ChannelID, typetag, VtuberGroupID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}

	return nil
}

//delete discord channel from `disable` command
func DelChannel(ChannelID string, VtuberGroupID int64) error {
	stmt, err := DB.Prepare(`DELETE From Channel WHERE DiscordChannelID=? AND VtuberGroup_id=?`)
	if err != nil {
		return err
	}

	stmt.Exec(ChannelID, VtuberGroupID)
	defer stmt.Close()
	return nil

}

//delete discord channel type from `update` command
func UpdateChannel(ChannelID string, typetag int, VtuberGroupID int64) error {
	if ChannelCheck(VtuberGroupID, ChannelID) {
		var tmp int
		row := DB.QueryRow(`SELECT Type FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?`, VtuberGroupID, ChannelID)
		row.Scan(&tmp)
		if tmp == typetag {
			return errors.New("Already enable type on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set Type=? Where VtuberGroup_id=? AND DiscordChannelID=?`, typetag, VtuberGroupID, ChannelID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//Get discord channel id from VtuberGroup
func (Data GroupName) GetChannelByGroup() []string {
	var channellist []string
	rows, err := DB.Query(`SELECT DiscordChannelID FROM Channel WHERE VtuberGroup_id=? group by DiscordChannelID`, Data.ID)
	BruhMoment(err, "", false)

	defer rows.Close()
	for rows.Next() {
		var list string
		err = rows.Scan(&list)
		BruhMoment(err, "", false)

		channellist = append(channellist, list)
	}
	return channellist
}

//Check Discord Channel from VtuberGroup
func ChannelCheck(VtuberGroupID int64, ChannelID string) bool {
	funcvar := GetFunctionName(ChannelCheck)
	Debugging(funcvar, "In", fmt.Sprint(VtuberGroupID, ChannelID))
	var tmp int
	row := DB.QueryRow("SELECT id FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?", VtuberGroupID, ChannelID)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		Debugging(funcvar, "Out", false)
		return false
	} else {
		Debugging(funcvar, "Out", true)
		return true
	}
}

//Check enable or disable discord channel from `tag,del` command
func CheckChannelEnable(ChannelID, VtuberName string, GroupID int64) bool {
	funcvar := GetFunctionName(CheckChannelEnable)
	Debugging(funcvar, "In", fmt.Sprint(ChannelID, VtuberName))
	var DiscordChannelID string
	row := DB.QueryRow("Select DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where Channel.DiscordChannelID=? AND (VtuberMember.VtuberName_EN=? OR VtuberMember.VtuberName_JP=? OR VtuberGroup.id=?) group by Channel.id", ChannelID, VtuberName, VtuberName, GroupID)
	err := row.Scan(&DiscordChannelID)
	if err != nil || err == sql.ErrNoRows {
		Debugging(funcvar, "Out", false)
		return false
	} else {
		Debugging(funcvar, "Out", true)
		return true
	}
}

//Get userinfo(tags) from discord channel
func UserStatus(UserID, Channel string) []string {
	funcvar := GetFunctionName(UserStatus)
	Debugging(funcvar, "In", fmt.Sprint(UserID, Channel))
	var taglist []string
	rows, err := DB.Query(`SELECT VtuberName FROM User INNER JOIN VtuberMember ON User.VtuberMember_id= VtuberMember.id Inner Join Channel on Channel.id=User.Channel_id WHERE DiscordChannelID=? And DiscordID=?`, Channel, UserID)
	BruhMoment(err, "", false)

	defer rows.Close()
	for rows.Next() {
		var list string
		err = rows.Scan(&list)
		BruhMoment(err, "", false)

		taglist = append(taglist, list)
	}
	Debugging(funcvar, "Out", taglist)
	return taglist
}

//Get Discord channel status
func ChannelStatus(ChannelID string) ([]string, []int) {
	funcvar := GetFunctionName(UserStatus)
	Debugging(funcvar, "In", ChannelID)
	var (
		Taglist []string
		Type    []int
	)
	rows, err := DB.Query(`SELECT VtuberGroupName,Channel.Type FROM Channel INNER JOIn VtuberGroup on VtuberGroup.id=Channel.VtuberGroup_id WHERE DiscordChannelID=?`, ChannelID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		var (
			tmp  string
			tmp2 int
		)
		err = rows.Scan(&tmp, &tmp2)
		BruhMoment(err, "", false)

		Taglist = append(Taglist, tmp)
		Type = append(Type, tmp2)
	}
	Debugging(funcvar, "Out", fmt.Sprint(Taglist, Type))
	return Taglist, Type
}

//get channel tags from `channel tags` command
func ChannelTag(MemberID int64, typetag int) ([]int, []string) {
	funcvar := GetFunctionName(ChannelTag)
	Debugging(funcvar, "In", MemberID)
	var (
		IdDiscordChannelID []int
		DiscordChannelID   []string
	)
	rows, err := DB.Query(`Select Channel.id,DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=? OR Channel.type=3)`, MemberID, typetag)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		var (
			tmp  int
			tmp2 string
		)
		err = rows.Scan(&tmp, &tmp2)
		BruhMoment(err, "", false)

		IdDiscordChannelID = append(IdDiscordChannelID, tmp)
		DiscordChannelID = append(DiscordChannelID, tmp2)
	}
	Debugging(funcvar, "Out", fmt.Sprint(IdDiscordChannelID, DiscordChannelID))
	return IdDiscordChannelID, DiscordChannelID
}

//get tags
func GetUserList(IdDiscordChannelID int, Member int64) []string {
	funcvar := GetFunctionName(GetUserList)
	Debugging(funcvar, "In", fmt.Sprint(IdDiscordChannelID, Member))
	var UserTagsList []string
	rows, err := DB.Query(`SELECT DiscordID From User WHERE Channel_id=? And VtuberMember_id =?`, IdDiscordChannelID, Member)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		BruhMoment(err, "", false)
		UserTagsList = append(UserTagsList, "<@"+tmp+">")
	}
	Debugging(funcvar, "Out", UserStatus)
	return UserTagsList
}

//Scrapping twitter followers
func (Data Name) GetTwitterFollow() TwitterUser {
	Client := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}
	request, err := http.NewRequest(http.MethodGet, "https://api.allorigins.win/raw?url=https://socialbearing.com/scripts/get-user.php?user="+Data.TwitterName, nil)
	if err != nil {
		log.Error(err)
	}
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; MacOS x86_64; rv:81.0) Gecko/20100101 Firefox/81.0")

	result, err := Client.Do(request)
	if err != nil {
		log.Error(err)
	}

	if result.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"Status": result.StatusCode,
			"Reason": result.Status,
		}).Warn("Status code not daijobu")
	}

	defer result.Body.Close()
	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Error(err)

	}

	var Profile TwitterUser
	err = json.Unmarshal(body, &Profile)
	if err != nil {
		log.Error(err)
	}
	return Profile
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
