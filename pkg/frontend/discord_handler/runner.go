package discordhandler

import (
	"errors"
	"regexp"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

//Prefix command
const (
	Enable        = "enable"
	Disable       = "disable"
	Update        = "update"
	TagMe         = "tag me"
	SetReminder   = "set reminder"
	DelTag        = "del tag"
	MyTags        = "my tags"
	TagRoles      = "tag roles"
	RolesTags     = "roles info"
	DelRoles      = "del roles"
	RolesReminder = "roles reminder"
	ChannelState  = "channel state"
	VtuberData    = "vtuber data"
	Subscriber    = "subscriber"
	Upcoming      = "upcoming"
	Past          = "past"
	Live          = "live"
)

//StartInit running the fe
func StartInit(path string) error {
	conf, err := config.ReadConfig(path)
	if err != nil {
		return err
	}
	db := conf.CheckSQL()

	Bot, _ := discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		return err
	}

	database.Start(db)
	engine.Start()

	Bot.AddHandler(Fanart)
	Bot.AddHandler(Tags)
	Bot.AddHandler(EnableState)
	Bot.AddHandler(Status)
	Bot.AddHandler(Help)
	Bot.AddHandler(BiliBiliMessage)
	Bot.AddHandler(BiliBiliSpace)
	Bot.AddHandler(YoutubeMessage)
	Bot.AddHandler(SubsMessage)

	return nil
}

//ValidName Find a valid name from user input
func ValidName(Name string) Memberst {
	for _, Group := range engine.GroupData {
		for _, Member := range database.GetMembers(Group.ID) {
			if Name == strings.ToLower(Member.Name) || Name == strings.ToLower(Member.JpName) {
				return Memberst{
					VTName:     engine.FixName(Member.EnName, Member.JpName),
					ID:         Member.ID,
					YtChannel:  Member.YoutubeID,
					SpaceID:    Member.BiliBiliID,
					BiliAvatar: Member.BiliBiliAvatar,
					YtAvatar:   Member.YoutubeAvatar,
				}
			}
		}
	}
	return Memberst{}
}

//FindName Find a valid Vtuber name from message handler
func FindName(MemberName string) NameStruct {
	for _, Group := range engine.GroupData {
		for _, Name := range database.GetMembers(Group.ID) {
			if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
				return NameStruct{
					GroupName:  Group.NameGroup,
					GroupID:    Group.ID,
					MemberName: Name.Name,
					MemberID:   Name.ID,
				}
			}
		}
	}
	return NameStruct{}

}

//NameStruct struct
type NameStruct struct {
	GroupName  string
	GroupID    int64
	MemberName string
	MemberID   int64
}

//FindGropName Find a valid Vtuber Group from message handler
func FindGropName(GroupName string) (database.GroupName, error) {
	for _, Group := range engine.GroupData {
		if strings.ToLower(Group.NameGroup) == strings.ToLower(GroupName) {
			return Group, nil
		}
	}
	return database.GroupName{}, errors.New(GroupName + " Name Vtuber not valid")
}

//RemovePic Remove twitter pic
func RemovePic(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)pic\.twitter.com\/.+`).ReplaceAllString(text, "${1}$2")
}

//GetAuthorAvatar Get twitter avatar
func GetAuthorAvatar(username string) string {
	profile, err := twitterscraper.GetProfile(username)
	if err != nil {
		log.Error(err)
	}
	return strings.Replace(profile.Avatar, "normal.jpg", "400x400.jpg", -1)
}

//GetUserAvatar Get bilibili user avatar
func (Data DynamicSvr) GetUserAvatar() string {
	return Data.Data.Card.Desc.UserProfile.Info.Face
}

//CheckReg Check available region
func CheckReg(GroupName, Reg string) bool {
	for Key, Val := range engine.RegList {
		if strings.ToLower(Key) == strings.ToLower(GroupName) {
			for _, Region := range strings.Split(strings.ToLower(Val), ",") {
				if Region == Reg {
					return true
				}
			}
		}
	}
	return false
}
