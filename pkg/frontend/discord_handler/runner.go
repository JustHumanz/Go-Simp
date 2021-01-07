package discordhandler

import (
	"errors"
	"regexp"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	BotInfo *discordgo.User
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
	Info          = "info"
	Upcoming      = "upcoming"
	Past          = "past"
	Live          = "live"
	ModuleInfo    = "module"
)

//StartInit running the fe
func StartInit(path string) error {
	conf, err := config.ReadConfig(path)
	if err != nil {
		return err
	}
	db := conf.CheckSQL()

	Bot, _ := discordgo.New("Bot " + config.BotConf.Discord)
	err = Bot.Open()
	if err != nil {
		return err
	}
	BotInfo, err = Bot.User("@me")
	if err != nil {
		log.Error(err)
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
	Bot.AddHandler(Module)

	return nil
}

func Module(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.BotConf.BotPrefix.General
	if strings.HasPrefix(m.Content, Prefix) {
		if m.Content == Prefix+ModuleInfo {
			list := []string{}
			keys := make(map[string]bool)
			for _, Member := range database.GetModule() {
				if _, value := keys[Member]; !value {
					keys[Member] = true
					list = append(list, Member)
				}
			}
			_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(list, "\n"))
			if err != nil {
				log.Error(err)
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "Invalid command")
			if err != nil {
				log.Error(err)
			}
		}
	}
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
					GroupName:  Group.GroupName,
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
func FindGropName(GroupName string) (database.Group, error) {
	for _, Group := range engine.GroupData {
		if strings.ToLower(Group.GroupName) == strings.ToLower(GroupName) {
			return Group, nil
		}
	}
	return database.Group{}, errors.New(GroupName + " Name Vtuber not valid")
}

//RemovePic Remove twitter pic
func RemovePic(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)pic\.twitter.com\/.+`).ReplaceAllString(text, "${1}$2")
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
