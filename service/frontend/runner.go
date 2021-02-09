package main

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	BotInfo    *discordgo.User
	RegList    = make(map[string]string)
	GroupsName []string
	Payload    database.VtubersPayload
	configfile config.ConfigFile
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

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	time.Sleep(5 * time.Minute)
}

//StartInit running the fe
func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		Bot         *discordgo.Session
		WaitMigrate = true
	)
	c := cron.New()
	c.Start()

	StartBot := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Frontend",
		})

		WaitMigrate = res.WaitMigrate
		if err != nil {
			log.Fatalf("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}

		for _, Group := range Payload.VtuberData {
			GroupsName = append(GroupsName, Group.GroupName)
			list := []string{}
			keys := make(map[string]bool)
			for _, Member := range Group.Members {
				if _, value := keys[Member.Region]; !value {
					keys[Member.Region] = true
					list = append(list, Member.Region)
				}
			}
			RegList[Group.GroupName] = strings.Join(list, ",")
		}

		if !res.WaitMigrate {
			log.Info("Start Frontend")
			var err error
			Bot, err = discordgo.New("Bot " + configfile.Discord)
			if err != nil {
				log.Error(err)
			}

			err = Bot.Open()
			if err != nil {
				log.Error(err)
			}

			BotInfo, err = Bot.User("@me")
			if err != nil {
				log.Error(err)
			}

			config.GoSimpConf = configfile
			database.Start(configfile)

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
		} else {
			log.Info("Waiting migrate done")
		}
	}
	StartBot()

	if WaitMigrate {
		c.AddFunc("@every 0h10m0s", StartBot)
	} else {
		c.Stop()
	}

	go pilot.RunHeartBeat(gRCPconn, "Frontend")
	runfunc.Run(Bot)
}

func Module(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := configfile.BotPrefix.General
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
		}
	}
}

//ValidName Find a valid name from user input
func ValidName(Name string) Memberst {
	for _, Group := range Payload.VtuberData {
		for _, Member := range Group.Members {
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
func FindName(MemberName string) *NameStruct {
	for _, Group := range Payload.VtuberData {
		for _, Name := range Group.Members {
			if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
				return &NameStruct{
					Group:  Group,
					Member: Name,
				}
			}
		}
	}
	return nil

}

//NameStruct struct
type NameStruct struct {
	Group  database.Group
	Member database.Member
}

//FindGropName Find a valid Vtuber Group from message handler
func FindGropName(GroupName string) (database.Group, error) {
	for _, Group := range Payload.VtuberData {
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
	for Key, Val := range RegList {
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
