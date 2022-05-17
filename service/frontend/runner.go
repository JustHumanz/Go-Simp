package main

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

var TESTINGGUILD = "721009835889393705"

var (
	BotInfo            *discordgo.User
	RegList            = make(map[string]string)
	GroupsName         []string
	GuildList          []string
	GroupsPayload      []database.Group
	configfile         config.ConfigFile
	Bot                *discordgo.Session
	ServiceUUID        = uuid.New().String()
	VtuberGroupChoices []*discordgo.ApplicationCommandOptionChoice
)

const (
	Kings       = "kings"
	Upvote      = "upvote"
	UpdateState = "SelectChannel"
	ServiceName = config.FrontendService
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

//StartInit running the fe
func main() {
	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	log.Info("Get Payload")
	res, err := gRCPconn.GetBotPayload(context.Background(), &pilot.ServiceMessage{
		Message:     "Init " + ServiceName + " service",
		Service:     ServiceName,
		ServiceUUID: ServiceUUID,
	})
	if err != nil {
		if configfile.Discord != "" {
			pilot.ReportDeadService(err.Error(), "Frontend")
		}
		log.Error("Error when request payload: %s", err)
	}

	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Fatalln(err)
	}

	configfile.InitConf()

	hostname := engine.GetHostname()

	go func() {
		for {
			log.WithFields(log.Fields{
				"Running": false,
				"UUID":    ServiceUUID,
			}).Info("request for running job")

			res2, err := gRCPconn.GetAgencyPayload(context.Background(), &pilot.ServiceMessage{
				Service:     ServiceName,
				Message:     "Request",
				ServiceUUID: ServiceUUID,
				Hostname:    hostname,
			})
			if err != nil {
				log.Error(err)
			}

			err = json.Unmarshal(res2.AgencyVtubers, &GroupsPayload)
			if err != nil {
				log.Fatalln(err)
			}

			for _, Group := range GroupsPayload {
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

			time.Sleep(30 * time.Minute)
		}
	}()

	log.Info("Start Frontend")
	Bot = engine.StartBot(false)

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

	for _, v := range GroupsPayload {
		VtuberGroupChoices = append(VtuberGroupChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  v.GroupName,
			Value: v.ID,
		})
	}

	database.Start(configfile)

	err = Bot.UpdateStreamingStatus(0, config.GoSimpConf.BotPrefix.General+"help", config.VtubersData)
	if err != nil {
		log.Error(err)
	}

	Bot.AddHandler(Help)

	go pilot.RunHeartBeat(gRCPconn, ServiceName, ServiceUUID)
	go engine.InitSlash(Bot, GroupsPayload, nil)

	Bot.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	guildID := func() []string {
		ids := []string{}
		for _, v := range Bot.State.Guilds {
			ids = append(ids, v.ID)
		}
		return ids
	}()

	Bot.AddHandler(func(s *discordgo.Session, g *discordgo.GuildCreate) {
		if g.Unavailable {
			log.Error("joined unavailable guild", g.Guild.ID)
			return
		}

		for _, v := range guildID {
			if v == g.ID {
				return
			}
		}

		log.WithFields(log.Fields{
			"GuildName": g.Name,
			"OwnerID":   g.OwnerID,
			"JoinDate":  g.JoinedAt.Format(time.RFC822),
		}).Info("New invite")
		engine.InitSlash(Bot, GroupsPayload, g.Guild)

	})

	Bot.AddHandler(func(s *discordgo.Session, g *discordgo.GuildDelete) {
		if g.Unavailable {
			log.Error("joined unavailable guild", g.Guild.ID)
			return
		}

		for k, v := range guildID {
			if v == g.ID {
				guildID[k] = ""
			}
		}

		log.WithFields(log.Fields{
			"GuildName": g.Name,
			"OwnerID":   g.OwnerID,
			"JoinDate":  g.JoinedAt.Format(time.RFC822),
		}).Info("Bot got kicked")

	})

	runfunc.Run(Bot)
}

//FindName Find a valid Vtuber name from message handler
func FindVtuber(M interface{}) (database.Member, error) {
	MemberName, str := M.(string)
	if str {
		for _, Group := range GroupsPayload {
			for _, Name := range Group.Members {
				if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName || MemberName == strconv.Itoa(int(Name.ID)) {
					return Name, nil
				}
			}
		}
	} else {
		MemberID := M.(int64)
		for _, Group := range GroupsPayload {
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
		for _, Group := range GroupsPayload {
			if strings.EqualFold(Group.GroupName, Grp) || strconv.Itoa(int(Group.ID)) == Grp {
				return Group, nil
			}
		}
	} else {
		GrpID := g.(int64)
		for _, Group := range GroupsPayload {
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

	MsgText, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "```"+tableString.String()+"```")
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

		if !Data.ChannelState.LiveOnly {
			err = Data.ChannelState.UpdateChannel(config.NewUpcoming)
			if err != nil {
				return err
			}
		}

		err = Data.ChannelState.UpdateChannel(config.LiveOnly)
		if err != nil {
			return err
		}

		if Data.ChannelState.Group.GroupName == config.Indie {
			err = Data.ChannelState.UpdateChannel(config.IndieNotif)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (Data *ChannelRegister) AddRegion() {
	GroupName := Data.ChannelState.Group.GroupName
	ChannelID := Data.ChannelState.ChannelID

	Data.UpdateState(AddRegion)
	RegEmoji := []string{}
	ChannelRegion := strings.Split(Data.ChannelState.Region, ",")
	for _, v := range ChannelRegion {
		if v != "" {
			RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v))
			Data.RegionTMP = append(Data.RegionTMP, v)
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

	Data.UpdateMessageID(MsgTxt2.ID)
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

			Data.Stop()
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

	Data.BreakPoint(time.Duration(len(RegEmoji)) * time.Second)
	Data.FixRegion("add")
	Data.ChannelState.UpdateChannel(config.Region)

	_, err = Bot.ChannelMessageSend(ChannelID, "Done,you added "+strings.Join(Data.AddRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
}

func (Data *ChannelRegister) DelRegion() {
	GroupName := Data.ChannelState.Group.GroupName
	ChannelID := Data.ChannelState.ChannelID

	Data.UpdateState(DelRegion)
	RegEmoji := []string{}
	for Key, Val := range RegList {
		if Key == GroupName {
			for _, v2 := range strings.Split(Val, ",") {
				for _, v := range strings.Split(Data.ChannelState.Region, ",") {
					if v == v2 {
						RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v2))
						Data.RegionTMP = append(Data.RegionTMP, v2)
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

	Data.Stop()
	for _, v := range RegEmoji {
		err := Bot.MessageReactionAdd(ChannelID, MsgID, v)
		if err != nil {
			log.Error(err)
		}
	}
	Data.UpdateMessageID(MsgID)
	Data.BreakPoint(time.Duration(len(RegEmoji)) * time.Second)
	Data.FixRegion("del")
	err = Data.ChannelState.UpdateChannel(config.Region)
	if err != nil {
		log.Error(err)
	}

	_, err = Bot.ChannelMessageSend(ChannelID, "Done,you remove "+strings.Join(Data.DelRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
}

func (Data *ChannelRegister) CheckNSFW() bool {
	ChannelRaw, err := Bot.Channel(Data.ChannelState.ChannelID)
	if err != nil {
		log.Error(err)
	}

	if !ChannelRaw.NSFW {
		if Data.ChannelState.TypeTag == 69 {
			_, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "[error]This Channel was not a NSFW channel")
			if err != nil {
				log.Error(err)
			}
			return false
		} else {
			_, err := Bot.ChannelMessageSend(Data.ChannelState.ChannelID, "[error]This Channel was not a NSFW channel,change channel type to fanart")
			if err != nil {
				log.Error(err)
			}
			return false
		}

	}
	return true
}

func Clear(v interface{}) {
	p := reflect.ValueOf(v).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func gacha() bool {
	return rand.Float32() < 0.5
}
