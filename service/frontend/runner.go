package main

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
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
	GuildList  []string
	Payload    database.VtubersPayload
	configfile config.ConfigFile
	Bot        *discordgo.Session
	FanBase    = "simps"
)

//Prefix command
const (
	Enable        = "enable"
	Disable       = "disable"
	Update        = "update_v1"
	Update2       = "update"
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
	Setup         = "setup"
	Kings         = "kings"
	Upvote        = "upvote"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

//StartInit running the fe
func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	var (
		WaitMigrate = true
		Counter     int
	)
	c := cron.New()
	c.Start()

	StartBot := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Frontend",
		})
		if err != nil {
			if configfile.Discord != "" {
				pilot.ReportDeadService(err.Error())
			}
			log.Fatalf("Error when request payload: %s", err)
		}

		WaitMigrate = res.WaitMigrate
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

		if !res.WaitMigrate || Counter == 3 {
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

			for _, GuildID := range Bot.State.Guilds {
				GuildList = append(GuildList, GuildID.ID)
			}

			configfile.InitConf()
			database.Start(configfile)

			err = Bot.UpdateStreamingStatus(0, config.GoSimpConf.BotPrefix.General+"help", config.VtubersData)
			if err != nil {
				log.Error(err)
			}

			Bot.AddHandler(Fanart)
			Bot.AddHandler(Tags)
			Bot.AddHandler(EnableState)
			Bot.AddHandler(Status)
			Bot.AddHandler(Help)
			Bot.AddHandler(BiliBiliMessage)
			Bot.AddHandler(BiliBiliSpace)
			Bot.AddHandler(YoutubeMessage)
			Bot.AddHandler(SubsMessage)
			//Bot.AddHandler(Module)
			Bot.AddHandler(RegisterFunc)
			Bot.AddHandler(Answer)
			c.Stop()

		} else {
			log.Info("Waiting migrate done")
			Counter++
		}
	}
	StartBot()

	if WaitMigrate {
		c.AddFunc("@every 0h5m0s", StartBot)
	} else if !WaitMigrate || Counter == 10 {
		c.Stop()
	}

	go pilot.RunHeartBeat(gRCPconn, "Frontend")
	runfunc.Run(Bot)
}

/*
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
*/

//FindName Find a valid Vtuber name from message handler
func FindVtuber(MemberName string, ID int64) database.Member {
	if ID == 0 {
		for _, Group := range Payload.VtuberData {
			for _, Name := range Group.Members {
				if Name.ID == ID {
					return Name
				}
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			for _, Name := range Group.Members {
				if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
					return Name
				}
			}
		}
	}

	return database.Member{}
}

//FindGropName Find a valid Vtuber Group from message handler
func FindGropName(g interface{}) (database.Group, error) {
	ID, err := strconv.Atoi(g.(string))
	if err != nil {
		for _, Group := range Payload.VtuberData {
			if strings.ToLower(Group.GroupName) == strings.ToLower(g.(string)) {
				return Group, nil
			}
		}
	} else {
		for _, Group := range Payload.VtuberData {
			if Group.ID == int64(ID) {
				return Group, nil
			}
		}
	}
	return database.Group{}, errors.New(g.(string) + " Name Vtuber not valid")
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

func MemberHasPermission(guildID string, userID string) (bool, error) {
	member, err := Bot.State.Member(guildID, userID)
	if err != nil {
		if member, err = Bot.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	if len(member.Roles) == 0 {
		return false, errors.New("You not enabled by any roles with administrator permissions")
	}

	for _, roleID := range member.Roles {
		role, err := Bot.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}

		if role.Permissions&config.ChannelPermission != 0 {
			return true, nil
		}
	}

	return false, nil
}
