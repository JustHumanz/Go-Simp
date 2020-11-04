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

//Get Twitter hashtag by group
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

//Get all vtuber groupData
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

//Get data of Vtuber member
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

//Get subs,follow,view,like data from Subscriber
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

//update bilibili state
func (Member *MemberSubs) UpBiliFollow(new int) *MemberSubs {
	Member.BiliFollow = new
	return Member
}

func (Member *MemberSubs) UpBiliVideo(new int) *MemberSubs {
	Member.BiliVideos = new
	return Member
}

func (Member *MemberSubs) UpBiliViews(new int) *MemberSubs {
	Member.BiliViews = new
	return Member
}

//update youtube state
func (Member *MemberSubs) UpYtSubs(new int) *MemberSubs {
	Member.YtSubs = new
	return Member
}

func (Member *MemberSubs) UpYtVideo(new int) *MemberSubs {
	Member.YtVideos = new
	return Member
}

func (Member *MemberSubs) UpYtViews(new int) *MemberSubs {
	Member.YtViews = new
	return Member
}

//update twitter state
func (Member *MemberSubs) UptwFollow(new int) *MemberSubs {
	Member.TwFollow = new
	return Member
}

//Update Subscriber data
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

//Get Member fanart URL from TBiliBili and Twitter
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
	var ChannelID int
	err := DB.QueryRow("SELECT id from Channel where DiscordChannelID=? AND VtuberGroup_id=?", DiscordChannelID, GroupID).Scan(&ChannelID)
	BruhMoment(err, "", false)
	return ChannelID
}

//Add user form `tag me command`
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

//Add new discord channel from `enable` command
func AddChannel(ChannelID string, typetag int, VtuberGroupID int64) error {
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
	defer stmt.Close()
	stmt.Exec(ChannelID, VtuberGroupID)
	return nil

}

//update discord channel type from `update` command
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
func ChannelCheck(VtuberGroupID int64, ChannelID string) bool {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?", VtuberGroupID, ChannelID)
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
func ChannelStatus(ChannelID string) ([]string, []int) {
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
	return Taglist, Type
}

//get channel tags from `channel tags` command
func ChannelTag(MemberID int64, typetag int) ([]int, []string) {
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
