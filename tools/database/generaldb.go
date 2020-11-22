package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	DB    *sql.DB
	debug bool
)

//Start Database session
func Start(dbsession *sql.DB) {
	DB = dbsession
	log.Info("Database module ready")
}

//BruhMoment Bruh moment
func BruhMoment(err error, msg string, exit bool) {
	if err != nil {
		log.Info(msg)
		log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

//GetHashtag Get Twitter hashtag by group
func GetHashtag(GroupID int64) []MemberGroupID {
	rows, err := DB.Query(`SELECT VtuberMember.id,VtuberName,VtuberName_JP,VtuberGroup_id,Hashtag,VtuberGroupName,VtuberGroupIcon FROM VtuberMember INNER Join VtuberGroup ON VtuberGroup.id = VtuberMember.VtuberGroup_id WHERE Hashtag !="" AND VtuberGroup.id =?`, GroupID)
	BruhMoment(err, "", false)
	defer rows.Close()

	var (
		Data []MemberGroupID
		list MemberGroupID
	)
	for rows.Next() {
		err = rows.Scan(&list.MemberID, &list.EnName, &list.JpName, &list.GroupID, &list.TwitterHashtags, &list.GroupName, &list.GroupIcon)
		BruhMoment(err, "", false)

		Data = append(Data, list)

	}
	return Data
}

//GetGroup Get all vtuber groupData
func GetGroup() []GroupName {
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
	return Data
}

//GetName Get data of Vtuber member
func GetName(GroupID int64) []Name {
	rows, err := DB.Query(`call GetVtuberName(?)`, GroupID)
	BruhMoment(err, "", false)
	defer rows.Close()

	var Data []Name
	for rows.Next() {
		var list Name
		err = rows.Scan(&list.ID, &list.Name, &list.EnName, &list.JpName, &list.YoutubeID, &list.BiliBiliID, &list.BiliRoomID, &list.Region, &list.TwitterHashtags, &list.BiliBiliHashtags, &list.BiliBiliAvatar, &list.TwitterName, &list.YoutubeAvatar)
		BruhMoment(err, "", false)
		Data = append(Data, list)

	}
	return Data
}

//gacha is gacha
func gacha() bool {
	return rand.Float32() < 0.5
}

//GetSubsCount Get subs,follow,view,like data from Subscriber
func (Member Name) GetSubsCount() *MemberSubs {
	var Data MemberSubs
	rows, err := DB.Query(`SELECT * FROM Subscriber WHERE VtuberMember_id=?`, Member.ID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.YtSubs, &Data.YtVideos, &Data.YtViews, &Data.BiliFollow, &Data.BiliVideos, &Data.BiliViews, &Data.TwFollow, &Data.MemberID)
		BruhMoment(err, "", false)
	}
	return &Data
}

//UpBiliFollow update bilibili state
func (Member *MemberSubs) UpBiliFollow(new int) *MemberSubs {
	Member.BiliFollow = new
	return Member
}

//UpBiliVideo Add bilibili Videos
func (Member *MemberSubs) UpBiliVideo(new int) *MemberSubs {
	Member.BiliVideos = new
	return Member
}

//UpBiliViews Add views
func (Member *MemberSubs) UpBiliViews(new int) *MemberSubs {
	Member.BiliViews = new
	return Member
}

//UpYtSubs update youtube state
func (Member *MemberSubs) UpYtSubs(new int) *MemberSubs {
	Member.YtSubs = new
	return Member
}

//UpYtVideo Update youtube videos
func (Member *MemberSubs) UpYtVideo(new int) *MemberSubs {
	Member.YtVideos = new
	return Member
}

//UpYtViews Update youtube views
func (Member *MemberSubs) UpYtViews(new int) *MemberSubs {
	Member.YtViews = new
	return Member
}

//UptwFollow Update twitter state
func (Member *MemberSubs) UptwFollow(new int) *MemberSubs {
	Member.TwFollow = new
	return Member
}

//UpdateSubs Update Subscriber data
func (Data *MemberSubs) UpdateSubs(State string) {
	if State == "yt" {
		_, err := DB.Exec(`Update Subscriber set Youtube_Subscriber=?,Youtube_Videos=?,Youtube_Views=? Where id=? `, Data.YtSubs, Data.YtVideos, Data.YtViews, Data.ID)
		BruhMoment(err, "", false)
	} else if State == "bili" {
		_, err := DB.Exec(`Update Subscriber set BiliBili_Followers=?,BiliBili_Videos=?,BiliBili_Views=? Where id=? `, Data.BiliFollow, Data.BiliVideos, Data.BiliViews, Data.ID)
		BruhMoment(err, "", false)
	} else {
		_, err := DB.Exec(`Update Subscriber set Twitter_Followers=? Where id=? `, Data.TwFollow, Data.ID)
		BruhMoment(err, "", false)
	}
}

//GetFanart Get Member fanart URL from TBiliBili and Twitter
func GetFanart(GroupID, MemberID int64) DataFanart {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
		rows     *sql.Rows
		err      error
	)

	if gacha() {
		rows, err = DB.Query(`Call GetArt(?,?,'twitter')`, GroupID, MemberID)
		BruhMoment(err, "", false)
	} else {
		rows, err = DB.Query(`Call GetArt(?,?,'westtaiwan')`, GroupID, MemberID)
		BruhMoment(err, "", false)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.EnName, &Data.JpName, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Video, &Data.Text)
		if err != nil {
			log.Error(err)
		}
	}
	Data.Videos = Video.String
	Data.Photos = strings.Fields(PhotoTmp.String)
	return Data

}

//InputTwitter Input new fanart from twitter
func (Data InputTW) InputTwitter() {

	stmt, err := DB.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	BruhMoment(err, "", false)
	defer stmt.Close()

	res, err := stmt.Exec(Data.Url, Data.Author, Data.Like, Data.Photos, Data.Video, Data.Text, Data.TweetID, Data.MemberID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)
}

//UpdateTwitter Update Likes data
func (Data InputTW) UpdateTwitter(id int) {
	_, err := DB.Exec(`Update Twitter set Likes=? Where id=?`, Data.Like, id)
	BruhMoment(err, "", false)
}

//GetChannelID Get Channel id from Discord ChannelID and VtuberGroupID
func GetChannelID(DiscordChannelID string, GroupID int64) int {
	var ChannelID int
	err := DB.QueryRow("SELECT id from Channel where DiscordChannelID=? AND VtuberGroup_id=?", DiscordChannelID, GroupID).Scan(&ChannelID)
	BruhMoment(err, "", false)
	return ChannelID
}

//Adduser form `tag me command`
func (Data UserStruct) Adduser(MemberID int64) error {
	ChannelID := GetChannelID(Data.Channel_ID, Data.GroupID)
	tmp := CheckUser(Data.DiscordID, MemberID, ChannelID)
	if tmp {
		return errors.New("Already registered")
	} else {
		stmt, err := DB.Prepare(`INSERT INTO User (DiscordID,DiscordUserName,Human,Reminder,VtuberMember_id,Channel_id) values(?,?,?,?,?,?)`)
		BruhMoment(err, "", false)
		defer stmt.Close()
		res, err := stmt.Exec(Data.DiscordID, Data.DiscordUserName, Data.Human, Data.Reminder, MemberID, ChannelID)
		BruhMoment(err, "", false)

		_, err = res.LastInsertId()
		BruhMoment(err, "", false)

		return nil
	}
}

func (Data UserStruct) UpdateReminder(MemberID int64) error {
	ChannelID := GetChannelID(Data.Channel_ID, Data.GroupID)

	stmt, err := DB.Prepare(`Update User set Reminder=? where DiscordID=? And VtuberMember_id=? And Channel_id=?`)
	BruhMoment(err, "", false)
	defer stmt.Close()
	res, err := stmt.Exec(Data.Reminder, Data.DiscordID, MemberID, ChannelID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)

	return nil
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
	var tmp int
	row := DB.QueryRow("SELECT id FROM User WHERE DiscordID=? AND VtuberMember_id=? AND Channel_id=?", DiscordID, MemberID, Channel_ChannelID)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

type DiscordChannel struct {
	ChannelID     string
	TypeTag       int
	LiveOnly      bool
	NewUpcoming   bool
	VtuberGroupID int64
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
	Data.VtuberGroupID = new
	return Data
}

//Add new discord channel from `enable` command
func (Data *DiscordChannel) AddChannel() error {
	stmt, err := DB.Prepare(`INSERT INTO Channel (DiscordChannelID,Type,LiveOnly,NewUpcoming,VtuberGroup_id) values(?,?,?,?,?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	res, err := stmt.Exec(Data.ChannelID, Data.TypeTag, Data.LiveOnly, Data.NewUpcoming, Data.VtuberGroupID)
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
func (Data *DiscordChannel) DelChannel() error {
	var (
		ID int64
	)

	row := DB.QueryRow("SELECT id FROM Channel WHERE DiscordChannelID=? AND VtuberGroup_id=?", Data.ChannelID, Data.VtuberGroupID)
	err := row.Scan(&ID)
	if err != nil {
		return err
	}

	rows, err := DB.Query(`SELECT id FROM User Where Channel_id=?`, ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var tmp int64
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}

		stmt, err := DB.Prepare(`DELETE From User WHERE id=?`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		stmt.Exec(tmp)
	}

	stmt, err := DB.Prepare(`DELETE From Channel WHERE DiscordChannelID=? AND VtuberGroup_id=?`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	stmt.Exec(Data.ChannelID, Data.VtuberGroupID)
	return nil
}

//update discord channel type from `update` command
func (Data *DiscordChannel) UpdateChannel(UpdateType string) error {
	var (
		typ   int
		live  bool
		newup bool
	)
	rows, err := DB.Query(`SELECT Type,LiveOnly,NewUpcoming FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?`, Data.VtuberGroupID, Data.ChannelID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&typ, &live, &newup)
		BruhMoment(err, "", false)
	}
	if UpdateType == "Type" {
		if typ == Data.TypeTag {
			return errors.New("Already enable type on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set Type=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.TypeTag, Data.VtuberGroupID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	} else if UpdateType == "LiveOnly" {
		if live == Data.LiveOnly {
			return errors.New("Already enable LiveOnly on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set LiveOnly=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.LiveOnly, Data.VtuberGroupID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	} else {
		if newup == Data.NewUpcoming {
			return errors.New("Already enable LiveOnly on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set NewUpcoming=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.NewUpcoming, Data.VtuberGroupID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//Get DiscordChannelID from VtuberGroup
func (Data GroupName) GetChannelByGroup() ([]int, []string) {
	var (
		channellist []string
		idlist      []int
		list        string
		id          int
	)
	rows, err := DB.Query(`SELECT id,DiscordChannelID FROM Channel WHERE VtuberGroup_id=? group by DiscordChannelID`, Data.ID)
	BruhMoment(err, "", false)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &list)
		BruhMoment(err, "", false)
		channellist = append(channellist, list)
		idlist = append(idlist, id)
	}
	return idlist, channellist
}

//Check Discord Channel from VtuberGroup
func (Data *DiscordChannel) ChannelCheck() bool {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?", Data.VtuberGroupID, Data.ChannelID)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

//Check enable or disable discord channel from `tag,del` command
func CheckChannelEnable(ChannelID, VtuberName string, GroupID int64) bool {
	var DiscordChannelID string
	row := DB.QueryRow("Select DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where Channel.DiscordChannelID=? AND (VtuberMember.VtuberName_EN=? OR VtuberMember.VtuberName_JP=? OR VtuberGroup.id=?) group by Channel.id", ChannelID, VtuberName, VtuberName, GroupID)
	err := row.Scan(&DiscordChannelID)
	if err != nil || err == sql.ErrNoRows {
		return false
	} else {
		return true
	}
}

//Get userinfo(tags) from discord channel
func UserStatus(UserID, Channel string) [][]string {
	var (
		GroupName  string
		VtuberName string
		Reminder   string
		taglist    [][]string
	)
	rows, err := DB.Query(`SELECT VtuberGroupName,VtuberName,User.Reminder FROM User INNER JOIN VtuberMember ON User.VtuberMember_id=VtuberMember.id Join VtuberGroup ON VtuberGroup.id = VtuberMember.VtuberGroup_id Inner Join Channel on Channel.id=User.Channel_id WHERE DiscordChannelID=? And DiscordID=?`, Channel, UserID)
	BruhMoment(err, "", false)

	defer rows.Close()
	for rows.Next() {
		tmpReminder := 0
		err = rows.Scan(&GroupName, &VtuberName, &tmpReminder)
		BruhMoment(err, "", false)

		if tmpReminder == 0 {
			Reminder = "None"
		} else {
			Reminder = strconv.Itoa(tmpReminder) + " Minutes"
		}

		taglist = append(taglist, []string{GroupName, VtuberName, Reminder})
	}
	return taglist
}

//Get Discord channel status
func ChannelStatus(ChannelID string) ([]string, []int, []string, []string) {
	var (
		Taglist     []string
		Type        []int
		LiveOnly    []string
		NewUpcoming []string
	)
	rows, err := DB.Query(`SELECT VtuberGroupName,Channel.Type,Channel.LiveOnly,Channel.NewUpcoming FROM Channel INNER JOIn VtuberGroup on VtuberGroup.id=Channel.VtuberGroup_id WHERE DiscordChannelID=?`, ChannelID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		var (
			tmp  string
			tmp2 int
			tmp3 bool
			tmp4 bool
		)
		err = rows.Scan(&tmp, &tmp2, &tmp3, &tmp4)
		BruhMoment(err, "", false)
		if tmp3 {
			LiveOnly = append(LiveOnly, "Enable")
		} else {
			LiveOnly = append(LiveOnly, "Disable")
		}
		if tmp4 {
			NewUpcoming = append(NewUpcoming, "Enable")
		} else {
			NewUpcoming = append(NewUpcoming, "Disable")
		}
		Taglist = append(Taglist, tmp)
		Type = append(Type, tmp2)
	}
	return Taglist, Type, LiveOnly, NewUpcoming
}

//ChannelTag get channel tags from `channel tags` command
func ChannelTag(MemberID int64, typetag int, Options string) ([]int, []string) {
	var (
		IdDiscordChannelID []int
		DiscordChannelID   []string
		rows               *sql.Rows
		err                error
	)
	if Options == "LiveOnly" {
		rows, err = DB.Query(`Select Channel.id,DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=2 OR Channel.type=3) AND LiveOnly=0`, MemberID)
		BruhMoment(err, "", false)
	} else if Options == "NewUpcoming" {
		rows, err = DB.Query(`Select Channel.id,DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=2 OR Channel.type=3) AND NewUpcoming=1`, MemberID)
		BruhMoment(err, "", false)
	} else {
		rows, err = DB.Query(`Select Channel.id,DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=? OR Channel.type=3)`, MemberID, typetag)
		BruhMoment(err, "", false)
	}
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
	return IdDiscordChannelID, DiscordChannelID
}

//get tags
func GetUserList(ChannelIDDiscord int, Member int64) []string {
	var (
		UserTagsList  []string
		DiscordUserID string
		Type          bool
	)
	rows, err := DB.Query(`SELECT DiscordID,Human From User WHERE Channel_id=? And VtuberMember_id =?`, ChannelIDDiscord, Member)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&DiscordUserID, &Type)
		BruhMoment(err, "", false)
		if Type {
			UserTagsList = append(UserTagsList, "<@"+DiscordUserID+">")
		} else {
			UserTagsList = append(UserTagsList, "<@&"+DiscordUserID+">")
		}
	}
	return UserTagsList
}

//get Reminder tags
func GetUserReminderList(ChannelIDDiscord int, Member int64, Reminder int) []string {
	var (
		UserTagsList  []string
		DiscordUserID string
		Type          bool
	)
	rows, err := DB.Query(`SELECT DiscordID,Human From User WHERE Channel_id=? And VtuberMember_id=? And Reminder=?`, ChannelIDDiscord, Member, Reminder)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&DiscordUserID, &Type)
		BruhMoment(err, "", false)
		if Type {
			UserTagsList = append(UserTagsList, "<@"+DiscordUserID+">")
		} else {
			UserTagsList = append(UserTagsList, "<@&"+DiscordUserID+">")
		}
	}
	return UserTagsList
}

//Scrapping twitter followers
func (Data Name) GetTwitterFollow() TwitterUser {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	Client := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}
	request, err := http.NewRequest(http.MethodGet, "https://api.allorigins.win/raw?url=https://socialbearing.com/scripts/get-user.php?user="+Data.TwitterName, nil)
	if err != nil {
		log.Error(err)
	}
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; MacOS x86_64; rv:81.0) Gecko/20100101 Firefox/81.0")

	result, err := Client.Do(request.WithContext(ctx))
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
