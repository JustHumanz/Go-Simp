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
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/prediction"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

var (
	BotInfo        *discordgo.User
	RegList        = make(map[string]string)
	GroupsName     []string
	GuildList      []string
	GroupsPayload  *[]database.Group
	configfile     config.ConfigFile
	Bot            *discordgo.Session
	PredictionConn prediction.PredictionClient
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
	Predick       = "prediction"
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
	PredictionConn = prediction.NewPredictionClient(network.InitgRPC(config.Prediction))
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
		GetPayload := func() {
			log.Info("Get Payload")
			res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
				Message: "Send me nude",
				Service: "Frontend",
			})
			if err != nil {
				if configfile.Discord != "" {
					pilot.ReportDeadService(err.Error(), "Frontend")
				}
				log.Error("Error when request payload: %s", err)
			}

			WaitMigrate = res.WaitMigrate
			err = json.Unmarshal(res.ConfigFile, &configfile)
			if err != nil {
				log.Error(err)
			}

			configfile.InitConf()

			err = json.Unmarshal(res.VtuberPayload, &GroupsPayload)
			if err != nil {
				log.Error(err)
			}
		}
		GetPayload()

		for _, Group := range *GroupsPayload {
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

		if !WaitMigrate || Counter == 6 {
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
			Bot.AddHandler(TwitchMessage)
			Bot.AddHandler(SubsMessage)
			Bot.AddHandler(Lewd)
			Bot.AddHandler(StartRegister)
			Bot.AddHandler(EmojiHandler)
			Bot.AddHandler(UpdateChannel)
			c.Stop()
			c2 := cron.New()
			c2.Start()
			c2.AddFunc(config.CheckPayload, GetPayload)

		} else {
			log.Info("Waiting migrate done")
			Counter++
		}
	}
	StartBot()

	if WaitMigrate {
		c.AddFunc("@every 0h5m0s", StartBot)
	} else if !WaitMigrate || Counter == 6 {
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
			_, err := Bot.ChannelMessageSend(m.ChannelID, strings.Join(list, "\n"))
			if err != nil {
				log.Error(err)
			}
		}
	}
}
*/

//FindName Find a valid Vtuber name from message handler
func FindVtuber(M interface{}) (database.Member, error) {
	MemberName, str := M.(string)
	if str {
		for _, Group := range *GroupsPayload {
			for _, Name := range Group.Members {
				if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName || MemberName == strconv.Itoa(int(Name.ID)) {
					return Name, nil
				}
			}
		}
	} else {
		MemberID := M.(int64)
		for _, Group := range *GroupsPayload {
			for _, Name := range Group.Members {
				if MemberID == Name.ID {
					return Name, nil
				}
			}
		}
	}

	return database.Member{}, errors.New("not found")
}

//FindGropName Find a valid Vtuber Group from message handler
func FindGropName(g interface{}) (database.Group, error) {
	Grp, str := g.(string)
	if str {
		for _, Group := range *GroupsPayload {
			if strings.EqualFold(Group.GroupName, Grp) || strconv.Itoa(int(Group.ID)) == Grp {
				return Group, nil
			}
		}
	} else {
		GrpID := g.(int64)
		for _, Group := range *GroupsPayload {
			if Group.ID == GrpID {
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
		if strings.EqualFold(Key, GroupName) {
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
		return false, errors.New("you not enabled by any roles with administrator permissions")
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

func (Data *ChannelRegister) ChoiceType() *ChannelRegister {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	_, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Select Channel Type: ")
	if err != nil {
		log.Error(err)
	}
	Fanart := "Disabled"
	Live := "Disabled"
	Lewd := "Disabled"
	if Data.ChannelState.IsFanart() {
		Fanart = "Enabled"
	}

	if Data.ChannelState.IsLive() {
		Live = "Enabled"
	}

	if Data.ChannelState.IsLewd() {
		Lewd = "Enabled"
	}

	table.SetHeader([]string{"Type", "Select", "Status"})
	table.Append([]string{"Fanart", config.Art, Fanart})
	table.Append([]string{"Livestream", config.Live, Live})
	table.Append([]string{"Lewd Fanart", config.Lewd, Lewd})
	table.Render()

	MsgText, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "`"+tableString.String()+"`")
	if err != nil {
		Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Something error, "+err.Error())
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "TypeChannel",
		"MessageID": MsgText.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}

	Data.UpdateMessageID(MsgText.ID).EmojiTrue()
	return Data
}

func (Data *ChannelRegister) LiveOnly() *ChannelRegister {
	MsgTxt, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable LiveOnly? **Set livestreams in strict mode(ignoring covering or regular video) notification**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState(config.LiveOnly)
	return Data
}

func (Data *ChannelRegister) NewUpcoming() *ChannelRegister {
	MsgTxt, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable NewUpcoming? **Bot will send new upcoming livestream**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState(config.NewUpcoming)
	return Data
}

func (Data *ChannelRegister) Dynamic() *ChannelRegister {
	MsgTxt, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable Dynamic mode? **Livestream message will disappear after livestream ended**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState(config.Dynamic)
	return Data
}

func (Data *ChannelRegister) Lite() *ChannelRegister {
	MsgTxt, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable Lite mode? **Disabling ping user/role function**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState(config.LiteMode)
	return Data
}

func (Data *ChannelRegister) IndieNotif() *ChannelRegister {
	MsgTxt, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "Send all independent vtubers notification?")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, Bot)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState(config.IndieNotif)
	return Data
}

func (Data *ChannelRegister) UpdateChannel(s string) error {
	if s == config.Type {
		err := Data.ChannelState.UpdateChannel(config.Type)
		if err != nil {
			return err
		}
	} else {
		err := Data.ChannelState.UpdateChannel(config.LiveOnly)
		if err != nil {
			return err
		}

		err = Data.ChannelState.UpdateChannel(config.Dynamic)
		if err != nil {
			return err
		}

		if !Register.Payload[Data.Index].ChannelState.LiveOnly {
			err = Data.ChannelState.UpdateChannel(config.NewUpcoming)
			if err != nil {
				return err
			}
		}

		err = Data.ChannelState.UpdateChannel(config.LiveOnly)
		if err != nil {
			return err
		}

		if Register.Payload[Data.Index].ChannelState.Group.GroupName == config.Indie {
			err = Data.ChannelState.UpdateChannel(config.IndieNotif)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (Data *ChannelRegister) AddRegion() {
	GroupName := Register.Payload[Data.Index].ChannelState.Group.GroupName
	ChannelID := Register.Payload[Data.Index].ChannelState.ChannelID

	Register.Payload[Data.Index].UpdateState(AddRegion)
	RegEmoji := []string{}
	ChannelRegion := strings.Split(Register.Payload[Data.Index].ChannelState.Region, ",")
	for _, v := range ChannelRegion {
		if v != "" {
			RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v))
			Register.Payload[Data.Index].RegionTMP = append(Register.Payload[Data.Index].RegionTMP, v)
		}
	}
	_, err := Bot.ChannelMessageSend(ChannelID, "`"+GroupName+"` Regions you already enabled: "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgTxt2, err := Bot.ChannelMessageSend(ChannelID, "`"+GroupName+" `Select the region you want to add: ")
	if err != nil {
		log.Error(err)
	}

	Register.Payload[Data.Index].UpdateMessageID(MsgTxt2.ID)
	for Key, Val := range RegList {
		GroupRegList := strings.Split(Val, ",")

		if Key == GroupName {
			if len(ChannelRegion) == len(GroupRegList) {
				_, err := Bot.ChannelMessageSend(ChannelID, "You already enable all region")
				if err != nil {
					log.Error(err)
				}
				return
			}

			Register.Payload[Data.Index].Stop()
			for _, v2 := range GroupRegList {
				skip := false
				for _, v := range ChannelRegion {
					if v == v2 {
						skip = true
						break
					}
				}
				if !skip {
					err := Bot.MessageReactionAdd(ChannelID, MsgTxt2.ID, engine.CountryCodetoUniCode(v2))
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

	Register.Payload[Data.Index].BreakPoint(4)
	Register.Payload[Data.Index].FixRegion("add")
	Register.Payload[Data.Index].ChannelState.UpdateChannel(config.Region)

	_, err = Bot.ChannelMessageSend(ChannelID, "Done,you added "+strings.Join(Register.Payload[Data.Index].AddRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
	CleanRegister(Data.Index)
}

func (Data *ChannelRegister) DelRegion() {
	GroupName := Register.Payload[Data.Index].ChannelState.Group.GroupName
	ChannelID := Register.Payload[Data.Index].ChannelState.ChannelID

	Register.Payload[Data.Index].UpdateState(DelRegion)
	RegEmoji := []string{}
	for Key, Val := range RegList {
		if Key == GroupName {
			for _, v2 := range strings.Split(Val, ",") {
				for _, v := range strings.Split(Register.Payload[Data.Index].ChannelState.Region, ",") {
					if v == v2 {
						RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v2))
						Register.Payload[Data.Index].RegionTMP = append(Register.Payload[Data.Index].RegionTMP, v2)
					}
				}
			}
		}
	}

	_, err := Bot.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region you already enabled in this channel "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgID := ""
	if len(RegEmoji) == 0 {
		_, err := Bot.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region 404,add first")
		if err != nil {
			log.Error(err)
		}
		return
	} else {
		MsgTxt2, err := Bot.ChannelMessageSend(ChannelID, "`"+GroupName+"` Select region you want to delete : ")
		if err != nil {
			log.Error(err)
		}
		MsgID = MsgTxt2.ID
	}

	Register.Payload[Data.Index].Stop()
	for _, v := range RegEmoji {
		err := Bot.MessageReactionAdd(ChannelID, MsgID, v)
		if err != nil {
			log.Error(err)
		}
	}
	Register.Payload[Data.Index].UpdateMessageID(MsgID)
	Register.Payload[Data.Index].BreakPoint(4)
	Register.Payload[Data.Index].FixRegion("del")
	err = Register.Payload[Data.Index].ChannelState.UpdateChannel(config.Region)
	if err != nil {
		log.Error(err)
	}

	_, err = Bot.ChannelMessageSend(ChannelID, "Done,you remove "+strings.Join(Data.DelRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
	CleanRegister(Data.Index)
}

func (Data *ChannelRegister) CheckNSFW() bool {
	ChannelRaw, err := Bot.Channel(Data.ChannelState.ChannelID)
	if err != nil {
		log.Error(err)
	}

	if !ChannelRaw.NSFW {
		if Register.Payload[Data.Index].ChannelState.TypeTag == 69 {
			_, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "This Channel was not a NSFW channel")
			if err != nil {
				log.Error(err)
			}
			return false
		} else {
			_, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "This Channel was not a NSFW channel,change channel type to fanart")
			if err != nil {
				log.Error(err)
			}
			Register.Payload[Data.Index].ChannelState.TypeTag = 1
		}

	}
	return true
}
