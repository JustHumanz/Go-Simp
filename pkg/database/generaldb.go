package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	DB            *sql.DB
	UserTagCache  *redis.Client
	LiveCache     *redis.Client
	GeneralCache  *redis.Client
	UpcomingCache *redis.Client
)

//Start Database session
func Start(configfile config.ConfigFile) {
	DB = configfile.CheckSQL()
	RedisHost := configfile.Cached.Host + ":" + configfile.Cached.Port
	UserTagCache = redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: "",
		DB:       0,
	})

	LiveCache = redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: "",
		DB:       1,
	})

	GeneralCache = redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: "",
		DB:       2,
	})

	UpcomingCache = redis.NewClient(&redis.Options{
		Addr:     RedisHost,
		Password: "",
		DB:       3,
	})
	log.Info("Database module ready")
}

//GetGroups Get all vtuber groupData
func GetGroups() ([]Group, error) {
	rows, err := DB.Query(`SELECT id,VtuberGroupName,VtuberGroupIcon FROM VtuberGroup`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Data []Group
	for rows.Next() {
		var list Group
		err = rows.Scan(&list.ID, &list.GroupName, &list.IconURL)
		if err != nil {
			return nil, err
		}

		list.Members, err = GetMembers(list.ID)
		if err != nil {
			return nil, err
		}
		list.YoutubeChannels, err = GetGroupsYtChannel(list.ID)
		if err != nil {
			return nil, err
		}

		Data = append(Data, list)

	}
	return Data, nil
}

//Get vtuber agency yt channel
func GetGroupsYtChannel(i int64) ([]GroupYtChannel, error) {
	rows, err := DB.Query(`SELECT YoutubeChannel,Region FROM GroupYoutube where VtuberGroup_id=?`, i)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Data []GroupYtChannel
	for rows.Next() {
		var (
			list GroupYtChannel
		)
		err = rows.Scan(&list.YtChannel, &list.Region)
		if err != nil {
			return nil, err
		}
		list.ID = i
		Data = append(Data, list)
	}
	return Data, nil
}

//GetMembers Get data of Vtuber member
func GetMembers(GroupID int64) ([]Member, error) {
	var (
		list Member
		Data []Member
	)

	rows, err := DB.Query(`SELECT id, VtuberName, VtuberName_EN, VtuberName_JP, Twitter_Username, Twitter_Hashtag, Twitter_Avatar, Twitter_Banner, Twitter_Lewd, Youtube_ID, Youtube_Avatar, Youtube_Banner, BiliBili_SpaceID, BiliBili_RoomID, BiliBili_Avatar, BiliBili_Hashtag, BiliBili_Banner, Twitch_Username, Twitch_Avatar, Region, Fanbase, Status, VtuberGroup_id FROM Vtuber.VtuberMember WHERE VtuberGroup_id=? Order by Region,VtuberGroup_id;`, GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(
			&list.ID,
			&list.Name,
			&list.EnName,
			&list.JpName,
			&list.TwitterName,
			&list.TwitterHashtag,
			&list.TwitterAvatar,
			&list.TwitterBanner,
			&list.TwitterLewd,
			&list.YoutubeID,
			&list.YoutubeAvatar,
			&list.YoutubeBanner,
			&list.BiliBiliID,
			&list.BiliBiliRoomID,
			&list.BiliBiliAvatar,
			&list.BiliBiliHashtag,
			&list.BiliBiliBanner,
			&list.TwitchName,
			&list.TwitchAvatar,
			&list.Region,
			&list.Fanbase,
			&list.Status,
			&list.Group.ID,
		)
		if err != nil {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
		Data = append(Data, list)
	}
	return Data, nil
}

//GetSubsCount Get subs,follow,view,like data from Subscriber
func (Member Member) GetSubsCount() (*MemberSubs, error) {
	var (
		Data MemberSubs
		Key  = strconv.Itoa(int(Member.ID)) + Member.Name
	)

	val, err := GeneralCache.Get(context.Background(), Key).Result()
	if val == "" || err == redis.Nil {
		rows, err := DB.Query(`SELECT * FROM Subscriber WHERE VtuberMember_id=?`, Member.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.YtSubs, &Data.YtVideos, &Data.YtViews, &Data.BiliFollow, &Data.BiliVideos, &Data.BiliViews, &Data.TwFollow, &Data.TwitchFollow, &Data.TwitchViews, &Data.Member.ID)
			if err != nil {
				return nil, err
			}
		}
		err = GeneralCache.Set(context.Background(), Key, Data, config.GetSubsCountTTL).Err()
		if err != nil {
			return nil, err
		}

	} else {
		err := json.Unmarshal([]byte(val), &Data)
		if err != nil {
			return nil, err
		}
	}
	return &Data, nil
}

//RemoveSubsCache
func (Member Member) RemoveSubsCache() error {
	var (
		Key = strconv.Itoa(int(Member.ID)) + Member.Name
	)

	err := GeneralCache.Del(context.Background(), Key).Err()
	if err != nil {
		return err
	}
	return nil
}

//UpdateBiliBiliFollowers update bilibili state
func (Member *MemberSubs) UpdateBiliBiliFollowers(new int) *MemberSubs {
	Member.BiliFollow = new
	return Member
}

//UpdateBiliBiliVideos Add bilibili Videos
func (Member *MemberSubs) UpdateBiliBiliVideos(new int) *MemberSubs {
	Member.BiliVideos = new
	return Member
}

//UpdateBiliBiliViewers Add views
func (Member *MemberSubs) UpdateBiliBiliViewers(new int) *MemberSubs {
	Member.BiliViews = new
	return Member
}

//UpdateYoutubeSubs update youtube state
func (Member *MemberSubs) UpdateYoutubeSubs(new int) *MemberSubs {
	Member.YtSubs = new
	return Member
}

//UpdateYoutubeVideos Update youtube videos
func (Member *MemberSubs) UpdateYoutubeVideos(new int) *MemberSubs {
	Member.YtVideos = new
	return Member
}

//UpdateYoutubeViewers Update youtube views
func (Member *MemberSubs) UpdateYoutubeViewers(new int) *MemberSubs {
	Member.YtViews = new
	return Member
}

//UpdateTwitterFollowes Update twitter state
func (Member *MemberSubs) UpdateTwitterFollowes(new int) *MemberSubs {
	Member.TwFollow = new
	return Member
}

//UpdateTwitchFollowes Update Twitch state
func (Member *MemberSubs) UpdateTwitchFollowes(new int) *MemberSubs {
	Member.TwitchFollow = new
	return Member
}

//UpdateTwitchViewers Update Twitch state
func (Member *MemberSubs) UpdateTwitchViewers(new int) *MemberSubs {
	Member.TwitchViews = new
	return Member
}

//UpdateSubs Update Subscriber data
func (Member *MemberSubs) UpdateSubs() error {
	if Member.State == config.YoutubeLive {
		_, err := DB.Exec(`Update Subscriber set Youtube_Subscriber=?,Youtube_Videos=?,Youtube_Views=? Where id=? `, Member.YtSubs, Member.YtVideos, Member.YtViews, Member.ID)
		if err != nil {
			return err
		}
	} else if Member.State == config.BiliLive {
		_, err := DB.Exec(`Update Subscriber set BiliBili_Followers=?,BiliBili_Videos=?,BiliBili_Views=? Where id=? `, Member.BiliFollow, Member.BiliVideos, Member.BiliViews, Member.ID)
		if err != nil {
			return err
		}
	} else if Member.State == config.TwitchLive {
		_, err := DB.Exec(`Update Subscriber set Twitch_Followers=?,Twitch_Views=?, Where id=? `, Member.TwitchFollow, Member.TwitchViews, Member.ID)
		if err != nil {
			return err
		}
	} else {
		_, err := DB.Exec(`Update Subscriber set Twitter_Followers=? Where id=? `, Member.TwFollow, Member.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetChannelID Get Channel id from Discord ChannelID and VtuberGroupID
func GetChannelID(DiscordChannelID string, GroupID int64) (int, error) {
	var ChannelID int
	err := DB.QueryRow("SELECT id from Channel where DiscordChannelID=? AND VtuberGroup_id=?", DiscordChannelID, GroupID).Scan(&ChannelID)
	if err != nil {
		return 0, err
	}
	return ChannelID, nil
}

//Adduser form `tag me command`
func (Data *UserStruct) Adduser() error {
	ChannelID, err := GetChannelID(Data.Channel_ID, Data.Group.ID)
	if err != nil {
		return err
	}
	tmp := CheckUser(Data.DiscordID, Data.Member.ID, ChannelID)
	if tmp {
		return errors.New("already registered")
	} else {
		stmt, err := DB.Prepare(`INSERT INTO User (DiscordID,DiscordUserName,Human,Reminder,VtuberMember_id,Channel_id) values(?,?,?,?,?,?)`)
		if err != nil {
			return err
		}
		defer stmt.Close()
		res, err := stmt.Exec(Data.DiscordID, Data.DiscordUserName, Data.Human, Data.Reminder, Data.Member.ID, ChannelID)
		if err != nil {
			return err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return err
		}

		return nil
	}
}

//SendToCache send messageID to reddis
func (Data *UserStruct) SendToCache(MessageID string) error {
	err := GeneralCache.Set(context.Background(), MessageID, Data, config.AddUserTTL).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetChannelMessage get messageID from redis
func GetChannelMessage(MessageID string) (*UserStruct, error) {
	var (
		data UserStruct
	)
	val := GeneralCache.Get(context.Background(), MessageID).Val()
	err := json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

//UpdateReminder Update reminder time
func (Data UserStruct) UpdateReminder() error {
	ChannelID, err := GetChannelID(Data.Channel_ID, Data.Group.ID)
	if err != nil {
		return nil
	}
	tmp := CheckUser(Data.DiscordID, Data.Member.ID, ChannelID)
	if tmp {
		_, err := DB.Exec(`Update User set Reminder=? where DiscordID=? And VtuberMember_id=? And Channel_id=?`, Data.Reminder, Data.DiscordID, Data.Member.ID, ChannelID)
		if err != nil {
			return err
		}
	} else {
		return errors.New("User not tag " + strconv.Itoa(int(Data.Member.ID)))
	}
	return nil
}

//Deluser Delete user from `del` command
func (Data UserStruct) Deluser() error {
	ChannelID, err := GetChannelID(Data.Channel_ID, Data.Group.ID)
	if err != nil {
		return err
	}
	tmp := CheckUser(Data.DiscordID, Data.Member.ID, ChannelID)
	if tmp {
		stmt, err := DB.Prepare(`DELETE From User WHERE DiscordID=? AND VtuberMember_id=? AND Channel_id=?`)
		if err != nil {
			log.Error(err)
		}
		defer stmt.Close()

		stmt.Exec(Data.DiscordID, Data.Member.ID, ChannelID)
		return nil
	} else {
		return errors.New("already removed or you not in tag list")
	}
}

//CheckUser Check user if already tagged
func CheckUser(DiscordID string, MemberID int64, ChannelChannelID int) bool {
	var tmp int
	row := DB.QueryRow("SELECT id FROM User WHERE DiscordID=? AND VtuberMember_id=? AND Channel_id=?", DiscordID, MemberID, ChannelChannelID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Error(err)
		return true
	} else {
		return true
	}
}

//AddChannel Add new discord channel from `enable` command
func (Data *DiscordChannel) AddChannel() error {
	if Data.Dynamic {
		Data.SetNewUpcoming(false).SetLiveOnly(true)
	}
	stmt, err := DB.Prepare(`INSERT INTO Channel (DiscordChannelID,Type,LiveOnly,NewUpcoming,Dynamic,Region,IndieNotif,Lite,VtuberGroup_id) values(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}

	defer stmt.Close()
	res, err := stmt.Exec(Data.ChannelID, Data.TypeTag, Data.LiveOnly, Data.NewUpcoming, Data.Dynamic, Data.Region, Data.IndieNotif, Data.LiteMode, Data.Group.ID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}

	if Data.Dynamic {
		return errors.New("force to set Dynamic and live only true & disable new upcoming")
	}

	return nil
}

//DelChannel delete discord channel from `disable` command
func (Data *DiscordChannel) DelChannel() error {
	/*
		match, err := regexp.MatchString("Unknown|403||Delete|Missing|non-text", errmsg)
		if err != nil {
			return err
		}
	*/

	log.Warn("Delete Discord Channel ", Data.ChannelID)
	var (
		ID int64
	)

	row := DB.QueryRow("SELECT id FROM Channel WHERE DiscordChannelID=? AND VtuberGroup_id=?", Data.ChannelID, Data.Group.ID)
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
	stmt.Exec(Data.ChannelID, Data.Group.ID)
	return nil
}

//UpdateChannel update discord channel type from `update` command
func (Data *DiscordChannel) UpdateChannel(UpdateType string) error {
	var (
		channeltype int
		liveonly    bool
		newupcoming bool
		dynamic     bool
		lite        bool
		indienotif  bool
	)
	rows, err := DB.Query(`SELECT Type,LiveOnly,NewUpcoming,Dynamic,Lite,IndieNotif FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?`, Data.Group.ID, Data.ChannelID)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&channeltype, &liveonly, &newupcoming, &dynamic, &lite, &indienotif)
		if err != nil {
			log.Error(err)
		}
	}
	if channeltype == 1 {
		liveonly = false
		newupcoming = false
		dynamic = false
		lite = false
	}

	if UpdateType == config.Type {
		if channeltype == Data.TypeTag {
			if channeltype == 1 {
				return errors.New("already set fanart type on this channel")
			} else if channeltype == 2 {
				return errors.New("already set livestream type on this channel")
			} else if channeltype == 3 {
				return errors.New("already set fanart & livestream type on this channel")
			} else if channeltype == 69 {
				return errors.New("already set lewd type on this channel")
			} else if channeltype == 70 {
				return errors.New("already set lewd & fanart type on this channel")
			}

		} else {
			_, err := DB.Exec(`Update Channel set Type=?,LiveOnly=?,NewUpcoming=?,Dynamic=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.TypeTag, liveonly, newupcoming, dynamic, Data.Group.ID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	} else if UpdateType == config.LiveOnly {
		if liveonly == Data.LiveOnly && liveonly {
			return errors.New("already set LiveOnly on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set Type=?,LiveOnly=?,NewUpcoming=?,Dynamic=? Where VtuberGroup_id=? AND DiscordChannelID=?`, channeltype, Data.LiveOnly, newupcoming, dynamic, Data.Group.ID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	} else if UpdateType == config.Dynamic {
		if dynamic == Data.Dynamic && dynamic {
			return errors.New("already set Dynamic on this channel")
		} else if newupcoming && dynamic {
			newupcoming = false
			liveonly = true
		} else {
			_, err := DB.Exec(`Update Channel set Type=?,LiveOnly=?,NewUpcoming=?,Dynamic=? Where VtuberGroup_id=? AND DiscordChannelID=?`, channeltype, liveonly, newupcoming, Data.Dynamic, Data.Group.ID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	} else if UpdateType == config.Region {
		_, err := DB.Exec(`Update Channel set Region=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.Region, Data.Group.ID, Data.ChannelID)
		if err != nil {
			return err
		}
	} else if UpdateType == config.LiteMode {
		_, err := DB.Exec(`Update Channel set Lite=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.LiteMode, Data.Group.ID, Data.ChannelID)
		if err != nil {
			return err
		}
	} else if UpdateType == config.IndieNotif {
		_, err := DB.Exec(`Update Channel set IndieNotif=? Where VtuberGroup_id=? AND DiscordChannelID=?`, Data.IndieNotif, Data.Group.ID, Data.ChannelID)
		if err != nil {
			return err
		}
	} else {
		if newupcoming == Data.NewUpcoming && newupcoming {
			return errors.New("already set NewUpcoming on this channel")
		} else {
			_, err := DB.Exec(`Update Channel set Type=?,LiveOnly=?,NewUpcoming=?,Dynamic=? Where VtuberGroup_id=? AND DiscordChannelID=?`, channeltype, liveonly, Data.NewUpcoming, dynamic, Data.Group.ID, Data.ChannelID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//GetChannelByGroup Get DiscordChannelID from VtuberGroup
func (Data Group) GetChannelByGroup(Region string) ([]DiscordChannel, error) {
	var (
		list        string
		id          int64
		channeltype int
		ChannelData []DiscordChannel
	)
	rows, err := DB.Query(`SELECT id,DiscordChannelID,Type FROM Channel WHERE VtuberGroup_id=? AND (Type=1 OR Type=2 OR Type=3) AND (Channel.Region like ? OR Channel.Region='') group by DiscordChannelID`, Data.ID, "%"+Region+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &list, &channeltype)
		if err != nil {
			return nil, err
		}
		ChannelData = append(ChannelData, DiscordChannel{
			ID:        id,
			ChannelID: list,
			TypeTag:   channeltype,
		})
	}
	return ChannelData, nil
}

//ChannelCheck Check Discord Channel from VtuberGroup
func (Data *DiscordChannel) ChannelCheck() bool {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Channel WHERE VtuberGroup_id=? AND DiscordChannelID=?", Data.Group.ID, Data.ChannelID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Error(err)
		return true
	} else {
		return true
	}
}

//CheckIfNewChannel Check Discord Channel from VtuberGroup
func CheckIfNewChannel(ChannelID string) bool {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Channel WHERE DiscordChannelID=?", ChannelID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Error(err)
		return false
	} else {
		return true
	}
}

//CheckChannelEnable Check enable or disable discord channel from `tag,del` command
func CheckChannelEnable(ChannelID, VtuberName string, GroupID int64) bool {
	var DiscordChannelID string
	row := DB.QueryRow("Select DiscordChannelID FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where Channel.DiscordChannelID=? AND (VtuberMember.VtuberName_EN=? OR VtuberMember.VtuberName_JP=? OR VtuberGroup.id=?) group by Channel.id", ChannelID, VtuberName, VtuberName, GroupID)
	err := row.Scan(&DiscordChannelID)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Error(err)
		return false
	} else {
		return true
	}
}

//UserStatus Get userinfo(tags) from discord channel
func UserStatus(UserID, Channel string) ([][]string, error) {
	var (
		GroupName  string
		VtuberName string
		Reminder   string
		taglist    [][]string
	)
	rows, err := DB.Query(`SELECT VtuberGroupName,VtuberName,User.Reminder FROM User INNER JOIN VtuberMember ON User.VtuberMember_id=VtuberMember.id Join VtuberGroup ON VtuberGroup.id = VtuberMember.VtuberGroup_id Inner Join Channel on Channel.id=User.Channel_id WHERE DiscordChannelID=? And DiscordID=?`, Channel, UserID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		tmpReminder := 0
		err = rows.Scan(&GroupName, &VtuberName, &tmpReminder)
		if err != nil {
			return nil, err
		}

		if tmpReminder == 0 {
			Reminder = "None"
		} else {
			Reminder = strconv.Itoa(tmpReminder) + " Minutes"
		}

		taglist = append(taglist, []string{GroupName, VtuberName, Reminder})
	}
	return taglist, nil
}

//ChannelStatus Get Discord channel status
func ChannelStatus(ChannelID string) ([]DiscordChannel, error) {
	var (
		Data []DiscordChannel
	)

	rows, err := DB.Query(`SELECT Channel.*,VtuberGroup.VtuberGroupName,VtuberGroupIcon FROM Channel INNER JOIn VtuberGroup on VtuberGroup.id=Channel.VtuberGroup_id WHERE DiscordChannelID=?`, ChannelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			Region      string
			GroupData   Group
			ChannelData DiscordChannel
		)
		err = rows.Scan(&ChannelData.ID, &ChannelData.ChannelID, &ChannelData.TypeTag, &ChannelData.LiveOnly, &ChannelData.NewUpcoming, &ChannelData.Dynamic, &Region, &ChannelData.LiteMode, &ChannelData.IndieNotif, &GroupData.ID, &GroupData.GroupName, &GroupData.IconURL)
		if err != nil {
			return nil, err
		}
		ChannelData.Group = GroupData
		ChannelData.Region = strings.ToUpper(Region)
		Data = append(Data, ChannelData)
	}
	return Data, nil
}

//ChannelTag get channel tags from `channel tags` command
func ChannelTag(MemberID int64, typetag int, Options string, Reg string) ([]DiscordChannel, error) {
	var (
		Data        []DiscordChannel
		rows        *sql.Rows
		err         error
		Key         = strconv.Itoa(int(MemberID)) + strconv.Itoa(typetag) + Options + Reg
		rds         = UserTagCache
		DiscordChan = DiscordChannel{
			Member: Member{
				ID: MemberID,
			},
			Region: Reg,
		}
	)
	ctx := context.Background()
	val, err := rds.LRange(ctx, Key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	if err == redis.Nil || len(val) == 0 {
		if Options == "NotLiveOnly" {
			rows, err = DB.Query(`Select Channel.id,DiscordChannelID,Dynamic,Lite,IndieNotif,VtuberGroup.id FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=2 OR Channel.type=3) AND LiveOnly=0 AND (Channel.Region like ? OR Channel.Region='')`, MemberID, "%"+Reg+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()

		} else if Options == "NewUpcoming" {
			rows, err = DB.Query(`Select Channel.id,DiscordChannelID,Dynamic,Lite,IndieNotif,VtuberGroup.id FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=2 OR Channel.type=3) AND NewUpcoming=1 AND (Channel.Region like ? OR Channel.Region='')`, MemberID, "%"+Reg+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		} else if Options == "Lewd" {
			rows, err = DB.Query(`Select Channel.id,DiscordChannelID,Dynamic,Lite,IndieNotif,VtuberGroup.id FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=69 OR Channel.type=70) AND (Channel.Region like ? OR Channel.Region='')`, MemberID, "%"+Reg+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		} else if Options == "Default" {
			rows, err = DB.Query(`Select Channel.id,DiscordChannelID,Dynamic,Lite,IndieNotif,VtuberGroup.id FROM Channel Inner join VtuberGroup on VtuberGroup.id = Channel.VtuberGroup_id inner Join VtuberMember on VtuberMember.VtuberGroup_id = VtuberGroup.id Where VtuberMember.id=? AND (Channel.type=? OR Channel.type=3) AND (Channel.Region like ? OR Channel.Region='')`, MemberID, typetag, "%"+Reg+"%")
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		}
		for rows.Next() {
			err = rows.Scan(&DiscordChan.ID, &DiscordChan.ChannelID, &DiscordChan.Dynamic, &DiscordChan.LiteMode, &DiscordChan.IndieNotif, &DiscordChan.Group.ID)
			if err != nil {
				return nil, err
			}
			Data = append(Data, DiscordChan)
			err = rds.LPush(context.Background(), Key, DiscordChan).Err()
			if err != nil {
				return nil, err
			}
		}
		err = rds.Expire(context.Background(), Key, config.ChannelTagTTL).Err()
		if err != nil {
			return nil, err
		}
	} else {
		for _, result := range unique(val) {
			err := json.Unmarshal([]byte(result), &DiscordChan)
			if err != nil {
				return nil, err
			}
			Data = append(Data, DiscordChan)
		}
	}
	return Data, nil
}

//unique Remove dupicate string
var unique = func(Slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

//PushReddis Push DiscordChannel state to reddis
func (Data *DiscordChannel) PushReddis() error {
	err := GeneralCache.LPush(context.Background(), Data.VideoID, Data).Err()
	if err != nil {
		return err
	}
	return nil
}

//GetLiveNotifMsg get MessageID with live status
func GetLiveNotifMsg(Key string) ([]DiscordChannel, error) {
	var (
		Data        []DiscordChannel
		ChannelData DiscordChannel
	)
	val := GeneralCache.LRange(context.Background(), Key, 0, -1).Val()
	for _, v := range val {
		err := json.Unmarshal([]byte(v), &ChannelData)
		if err != nil {
			return nil, err
		}
		Data = append(Data, ChannelData)
	}
	rederr := GeneralCache.Del(context.Background(), Key).Err()
	if rederr != nil {
		return nil, rederr
	}

	return Data, nil
}

//GetUserList GetUser tags
func (Data *DiscordChannel) GetUserList(ctx context.Context) ([]string, error) {
	var (
		DataUser      []string
		DiscordUserID string
		Type          bool
		Key           = Data.Member.Name + strconv.Itoa(int(Data.ID))
		rds           = UserTagCache
	)
	val2, err := rds.LRange(ctx, Key, 0, -1).Result()
	if err == redis.Nil || len(val2) == 0 {
		rows, err := DB.Query(`SELECT DiscordID,Human From User WHERE Channel_id=? And VtuberMember_id=?`, Data.ID, Data.Member.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := ""
			err = rows.Scan(&DiscordUserID, &Type)
			if err != nil {
				return nil, err
			}

			if Type {
				tmp = "<@" + DiscordUserID + ">"
			} else {
				tmp = "<@&" + DiscordUserID + ">"
			}

			DataUser = append(DataUser, tmp)
			err := rds.LPush(ctx, Key, tmp).Err()
			if err != nil {
				log.Error(err)
			}
		}

		err = rds.Expire(ctx, Key, config.GetUserListTTL).Err()
		if err != nil {
			log.Error(err)
		}

	} else if err != nil {
		return nil, err
	} else {
		DataUser = val2
	}
	return DataUser, nil
}

//GetUserReminderList get Reminder tags
func GetUserReminderList(ChannelIDDiscord int64, Member int64, Reminder int) ([]string, error) {
	var (
		UserTagsList  []string
		DiscordUserID string
		Type          bool
		Key           = strconv.Itoa(int(ChannelIDDiscord)) + strconv.Itoa(int(Member)) + strconv.Itoa(int(Reminder))
		rds           = UserTagCache
		ctx           = context.Background()
	)
	val2, err := rds.LRange(ctx, Key, 0, -1).Result()
	if err == redis.Nil || len(val2) == 0 {
		rows, err := DB.Query(`SELECT DiscordID,Human From User WHERE Channel_id=? And VtuberMember_id=? And Reminder=?`, ChannelIDDiscord, Member, Reminder)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			tmp := ""
			err = rows.Scan(&DiscordUserID, &Type)
			if err != nil {
				return nil, err
			}
			if Type {
				tmp = "<@" + DiscordUserID + ">"
			} else {
				tmp = "<@&" + DiscordUserID + ">"
			}

			UserTagsList = append(UserTagsList, tmp)
			err = rds.LPush(ctx, Key, tmp).Err()
			if err != nil {
				log.Error(err)
			}
		}

		err = rds.Expire(ctx, Key, config.GetUserListTTL).Err()
		if err != nil {
			log.Error(err)
		}
	} else {
		UserTagsList = val2
	}

	return UserTagsList, nil
}

//GetRanChannel get random id channel
func GetRanChannel() string {
	var tmp string
	row := DB.QueryRow("SELECT DiscordChannelID FROM Vtuber.Channel ORDER BY RAND() LIMIT 1")
	err := row.Scan(&tmp)
	if err != nil {
		log.Error(err)
	}
	return tmp
}

//GetMemberCount get count of member
func GetMemberCount() int {
	var count int
	err := DB.QueryRow(`SELECT Count(*) from (Select COUNT(id) FROM Vtuber.User Group by DiscordID) t`).Scan(&count)
	if err != nil {
		log.Error(err)
	}
	return count
}

func DbStop() {
	err := DB.Close()
	if err != nil {
		log.Error(err)
	}
}

func (i *LiveStream) RemoveCache(Key string) error {
	ctx := context.Background()
	log.WithFields(log.Fields{
		"Key": Key,
	}).Info("Drop cache")

	err := LiveCache.Del(ctx, Key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (i *LiveStream) RemoveUpcomingCache(Key string) error {
	ctx := context.Background()
	log.WithFields(log.Fields{
		"Key": Key,
	}).Info("Drop cache")

	err := UpcomingCache.Del(ctx, Key).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetAllUser() []string {
	var (
		UserID []string
	)
	rows, err := DB.Query(`SELECT DiscordID FROM Vtuber.User where Human=1 group by DiscordID;`)
	if err != nil {
		log.Error(err)
	}

	defer rows.Close()
	for rows.Next() {
		id := ""
		err = rows.Scan(&id)
		if err != nil {
			log.Error(err)
		}
		UserID = append(UserID, id)
	}
	return UserID
}

func DeleteDeletedUser(users []string) {
	sqlq := fmt.Sprintf(`Delete FROM User where DiscordID in ('%s')`, strings.Join(users, "','"))
	_, err := DB.Query(sqlq)
	if err != nil {
		log.Error(err)
	}
}
