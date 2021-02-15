package main

import (
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

type Regis struct {
	Admin         string
	State         string
	MessageID     string
	RegionTMP     []string
	AddRegionVal  []string
	DelRegionVal  []string
	Gass          bool
	ChannelState  database.DiscordChannel
	ChannelStates []database.DiscordChannel
}

var (
	Register = &Regis{}
)

func Answer(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == Register.Admin && m.MessageID == Register.MessageID {
		Register.Start()
		if m.Emoji.MessageFormat() == config.Ok {
			if Register.State == "LiveOnly" {
				Register.SetLiveOnly(true)
			} else if Register.State == "NewUpcoming" {
				Register.SetNewUpcoming(true)
			} else if Register.State == "Dynamic" {
				Register.SetDynamic(true)
			}
		} else if m.Emoji.MessageFormat() == config.No {
			if Register.State == "LiveOnly" {
				Register.SetLiveOnly(false)
			} else if Register.State == "NewUpcoming" {
				Register.SetNewUpcoming(false)
			} else if Register.State == "Dynamic" {
				Register.SetDynamic(false)
			}
		}

		if m.Emoji.MessageFormat() == config.Art {
			if Register.ChannelState.TypeTag == 2 {
				Register.UpdateType(3)
			} else {
				Register.UpdateType(1)
			}
		} else if m.Emoji.MessageFormat() == config.Live {
			if Register.ChannelState.TypeTag == 1 {
				Register.UpdateType(3)
			} else {
				Register.UpdateType(2)
			}
		}

		if m.Emoji.MessageFormat() == config.One {
			Register.LiveOnly(s)

			if !Register.ChannelState.LiveOnly {
				Register.NewUpcoming(s)
			}

			Register.Dynamic(s)

			Register.UpdateChannel()
			_, err := s.ChannelMessageSend(m.ChannelID, "Done")
			if err != nil {
				log.Error(err)
			}
			return

		} else if m.Emoji.MessageFormat() == config.Two {
			Register.AddRegion(s)

		} else if m.Emoji.MessageFormat() == config.Three {
			Register.DelRegion(s)
		}

		if Register.State == "AddReg" {
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				Register.AddNewRegion(Region)
			}
		} else if Register.State == "DelReg" {
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				Register.RemoveRegion(Region)
			}

		}
	}
}

func RegisterFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	Prefix := configfile.BotPrefix.General
	Out := func() {
		_, err := s.ChannelMessageSend(m.ChannelID, "Adios")
		if err != nil {
			log.Error(err)
		}
		Register.Clear()
	}

	if strings.HasPrefix(m.Content, Prefix) {
		if CheckPermission(m.Author.ID, m.ChannelID, s) {
			if m.Content == Prefix+"setup" {
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to setup mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}

				Register.SetAdmin(m.Author.ID).SetChannel(m.ChannelID)
				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID of Vtuber group/agency you want to enable (Only one)")
				if err != nil {
					log.Error(err)
				}

				for _, v := range Payload.VtuberData {
					table.Append([]string{strconv.Itoa(int(v.ID)), v.GroupName})
				}
				table.SetHeader([]string{"ID", "GroupName"})
				table.Render()
				_, err = s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
				if err != nil {
					log.Error(err)
				}
				Register.UpdateState("Group")

			} else if m.Content == Prefix+"update_v2" {
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to update mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}
				ChannelData := database.ChannelStatus(m.ChannelID)
				if len(ChannelData) > 0 {
					var (
						Typestr string
					)
					Register.SetAdmin(m.Author.ID).SetChannel(m.ChannelID)

					Register.ChannelStates = ChannelData
					table.SetHeader([]string{"ID", "Group", "Type", "LiveOnly", "Dynamic", "NewUpcoming", "Region"})
					for i := 0; i < len(ChannelData); i++ {
						if ChannelData[i].TypeTag == 1 {
							Typestr = "Art"
						} else if ChannelData[i].TypeTag == 2 {
							Typestr = "Live"
						} else {
							Typestr = "All"
						}
						LiveOnly := config.No
						NewUpcoming := config.No
						Dynamic := config.No

						if ChannelData[i].LiveOnly {
							LiveOnly = config.Ok
						}

						if ChannelData[i].NewUpcoming {
							NewUpcoming = config.Ok
						}

						if ChannelData[i].Dynamic {
							Dynamic = config.Ok
						}
						table.Append([]string{strconv.Itoa(int(ChannelData[i].ID)), ChannelData[i].Group.GroupName, Typestr, LiveOnly, Dynamic, NewUpcoming, ChannelData[i].Region})
					}
					table.Render()
					_, err := s.ChannelMessageSend(m.ChannelID, "```"+tableString.String()+"```")
					if err != nil {
						log.Error(err)
					}
				} else {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetTitle("404 Not found").
						SetThumbnail(config.GoSimpIMG).
						SetImage(config.NotFound).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}
				Register.UpdateState("SelectChannel")
				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID of Channel state: ")
				if err != nil {
					log.Error(err)
				}
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
			if err != nil {
				log.Error(err)
			}
			Register.Clear()
			return
		}
	} else if m.Author.ID == Register.Admin && m.ChannelID == Register.ChannelState.ChannelID {
		if m.Content == "exit" {
			Out()
			return
		}
		if Register.State == "SelectChannel" {
			tmp, err := strconv.Atoi(m.Content)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "Worng input ID")
				if err != nil {
					log.Error(err)
				}
				Out()
				return
			} else {
				for _, ChannelState := range Register.ChannelStates {
					if int(ChannelState.ID) == tmp {
						Register.ChannelState = ChannelState
					}
				}
				if Register.ChannelState.ID != 0 {
					Register.SetChannel(m.ChannelID)
					_, err := s.ChannelMessageSend(m.ChannelID, "You select `"+Register.ChannelState.Group.GroupName+"` with ID `"+strconv.Itoa(int(Register.ChannelState.ID))+"`")
					if err != nil {
						log.Error(err)
					}
					table.SetHeader([]string{"Func"})
					table.Append([]string{"Update Channel state"})
					table.Append([]string{"Add region in this channel"})
					table.Append([]string{"Delete region in this channel"})
					table.Render()
					MsgText, err := s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
					if err != nil {
						log.Error(err)
					}
					err = engine.Reacting(map[string]string{
						"ChannelID": m.ChannelID,
						"State":     "Menu",
						"MessageID": MsgText.ID,
					}, s)
					if err != nil {
						log.Error(err)
					}
					Register.UpdateMessageID(MsgText.ID)

				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Channel ID not found")
					if err != nil {
						log.Error(err)
					}
					Out()
					return
				}
			}
		}
		if Register.State == "Group" {
			GroupID := 0
			tmp, err := strconv.Atoi(m.Content)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "Worng input ID")
				if err != nil {
					log.Error(err)
				}
				Out()
				return
			} else {
				GroupID = tmp
			}
			for _, v := range Payload.VtuberData {
				if int(v.ID) == GroupID {
					_, err := s.ChannelMessageSend(m.ChannelID, "Vtuber Group `"+v.GroupName+"`")
					if err != nil {
						log.Error(err)
					}
					Register.SetGroup(v)
					break
				}
			}
			//Register.ChannelState.ChannelCheck()
			if Register.ChannelState.Group.IsNull() {
				_, err := s.ChannelMessageSend(m.ChannelID, "Invalid ID,Group not found")
				if err != nil {
					log.Error(err)
				}
				Out()
				return
			}

			Register.Stop()
			for Key, Val := range RegList {
				if Key == Register.ChannelState.Group.GroupName {
					if len(Val) > 3 {
						MsgTxt, err := s.ChannelMessageSend(m.ChannelID, "Select `"+Key+"` region")
						if err != nil {
							log.Error(err)
						}
						Register.UpdateState("AddReg")
						for _, v := range strings.Split(Val, ",") {
							err := s.MessageReactionAdd(m.ChannelID, MsgTxt.ID, engine.CountryCodetoUniCode(v))
							if err != nil {
								log.Error(err)
							}
						}
						Register.UpdateMessageID(MsgTxt.ID)
						Register.BreakPoint(5)

						Register.FixRegion("add")
						if Register.ChannelState.ChannelCheck() {
							_, err := s.ChannelMessageSend(m.ChannelID, "Already setup `"+Register.ChannelState.Group.GroupName+"`,for add/del region use `Update` command")
							if err != nil {
								log.Error(err)
							}
							Out()
							return
						}
						Register.Stop()
					} else {
						Register.RegionTMP = strings.Split(Val, ",")
						Register.FixRegion("add")
					}
				}
			}

			_, err = s.ChannelMessageSend(m.ChannelID, "Select Channel Type: ")
			if err != nil {
				log.Error(err)
			}
			table.SetHeader([]string{"Type"})
			table.Append([]string{"Fanart " + config.Art})
			table.Append([]string{"Livestream " + config.Live})
			table.Render()

			MsgText, err := s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
			if err != nil {
				log.Error(err)
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": m.ChannelID,
				"State":     "TypeChannel",
				"MessageID": MsgText.ID,
			}, s)
			if err != nil {
				log.Error(err)
			}

			Register.UpdateMessageID(MsgText.ID)
			Register.BreakPoint(3)

			if Register.ChannelState.TypeTag != 1 {
				Register.Stop()
				Register.LiveOnly(s)
				Register.BreakPoint(1)

				if !Register.ChannelState.LiveOnly {
					Register.Stop()
					Register.NewUpcoming(s)
					Register.BreakPoint(1)
				}

				Register.Stop()
				Register.Dynamic(s)
				Register.BreakPoint(1)
			}
			err = Register.ChannelState.AddChannel()
			if err != nil {
				log.Error(err)
			}
			_, err = s.ChannelMessageSend(m.ChannelID, "Done,you add `"+Register.ChannelState.Group.GroupName+"` in this channel")
			if err != nil {
				log.Error(err)
			}
			Register.Clear()
			return
		}
	}
}

func (Data *Regis) LiveOnly(s *discordgo.Session) *Regis {
	MsgTxt, err := s.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable LiveOnly? **Set livestreams in strict mode(ignoring covering or regular video) notification**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, s)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState("LiveOnly")
	return Data
}

func (Data *Regis) NewUpcoming(s *discordgo.Session) *Regis {
	MsgTxt, err := s.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable NewUpcoming? **Bot will send new upcoming livestream**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, s)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState("NewUpcoming")
	return Data
}

func (Data *Regis) Dynamic(s *discordgo.Session) *Regis {
	MsgTxt, err := s.ChannelMessageSend(Data.ChannelState.ChannelID, "Enable Dynamic mode? **Livestream message will disappear after livestream ended**")
	if err != nil {
		log.Error(err)
	}
	err = engine.Reacting(map[string]string{
		"ChannelID": Data.ChannelState.ChannelID,
		"State":     "SelectType",
		"MessageID": MsgTxt.ID,
	}, s)
	if err != nil {
		log.Error(err)
	}
	Data.UpdateMessageID(MsgTxt.ID)
	Data.UpdateState("Dynamic")
	return Data
}

func (Data *Regis) UpdateChannel() error {
	err := Data.ChannelState.UpdateChannel("LiveOnly")
	ChannelID := Data.ChannelState.ChannelID
	if err != nil {
		_, err := Bot.ChannelMessageSend(ChannelID, err.Error())
		if err != nil {
			log.Error(err)
		}
	}
	err = Data.ChannelState.UpdateChannel("Dynamic")
	if err != nil {
		_, err := Bot.ChannelMessageSend(ChannelID, err.Error())
		if err != nil {
			log.Error(err)
		}
	}
	err = Data.ChannelState.UpdateChannel("NewUpcoming")
	if err != nil {
		_, err := Bot.ChannelMessageSend(ChannelID, err.Error())
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (Data *Regis) AddRegion(s *discordgo.Session) {
	GroupName := Register.ChannelState.Group.GroupName
	ChannelID := Register.ChannelState.ChannelID

	Register.UpdateState("AddReg")
	RegEmoji := []string{}
	ChannelRegion := strings.Split(Register.ChannelState.Region, ",")
	for _, v := range ChannelRegion {
		RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v))
		Register.RegionTMP = append(Register.RegionTMP, v)
	}
	_, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region you already enabled in here "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgTxt2, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Select region you want to add it : ")
	if err != nil {
		log.Error(err)
	}

	Register.UpdateMessageID(MsgTxt2.ID)
	for Key, Val := range RegList {
		if Key == GroupName {
			if len(ChannelRegion) == len(strings.Split(Val, ",")) {
				_, err := s.ChannelMessageSend(ChannelID, "You already enable all region")
				if err != nil {
					log.Error(err)
				}
				return
			}

			Register.Stop()
			for _, v2 := range strings.Split(Val, ",") {
				skip := false
				for _, v := range ChannelRegion {
					if v == v2 {
						skip = true
						break
					}
				}
				if !skip {
					err := s.MessageReactionAdd(ChannelID, MsgTxt2.ID, engine.CountryCodetoUniCode(v2))
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}

	Register.BreakPoint(3)
	Register.FixRegion("add")
	Register.ChannelState.UpdateChannel("Region")

	_, err = s.ChannelMessageSend(ChannelID, "Done,you add "+strings.Join(Register.AddRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
	Register.Clear()
	return
}

func (Data *Regis) DelRegion(s *discordgo.Session) {
	GroupName := Register.ChannelState.Group.GroupName
	ChannelID := Register.ChannelState.ChannelID

	Register.UpdateState("DelReg")
	RegEmoji := []string{}
	for Key, Val := range RegList {
		if Key == GroupName {
			for _, v2 := range strings.Split(Val, ",") {
				for _, v := range strings.Split(Register.ChannelState.Region, ",") {
					if v == v2 {
						RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v2))
						Register.RegionTMP = append(Register.RegionTMP, v2)
					}
				}
			}
		}
	}

	_, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region you already enabled in here "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgTxt2, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Select region you want to delete : ")
	if err != nil {
		log.Error(err)
	}

	Register.Stop()
	for _, v := range RegEmoji {
		err := s.MessageReactionAdd(ChannelID, MsgTxt2.ID, v)
		if err != nil {
			log.Error(err)
		}
	}
	Register.UpdateMessageID(MsgTxt2.ID)
	Register.BreakPoint(4)
	Register.FixRegion("del")
	Register.ChannelState.UpdateChannel("Region")

	_, err = s.ChannelMessageSend(ChannelID, "Done,you remove "+strings.Join(Data.DelRegionVal, ","))
	if err != nil {
		log.Error(err)
	}
	Register.Clear()
	return
}

func (Data *Regis) SetLiveOnly(new bool) *Regis {
	Data.ChannelState.LiveOnly = new
	return Data
}

func (Data *Regis) SetNewUpcoming(new bool) *Regis {
	Data.ChannelState.NewUpcoming = new
	return Data
}

func (Data *Regis) SetDynamic(new bool) *Regis {
	Data.ChannelState.Dynamic = new
	return Data
}

func (Data *Regis) SetChannel(new string) *Regis {
	Data.ChannelState.ChannelID = new
	return Data
}

func (Data *Regis) SetAdmin(new string) *Regis {
	Data.Admin = new
	return Data
}

func (Data *Regis) UpdateState(new string) *Regis {
	Data.State = new
	return Data
}

func (Data *Regis) SetGroup(new database.Group) *Regis {
	Data.ChannelState.Group = new
	return Data
}

func (Data *Regis) FixRegion(s string) {
	list := []string{}
	keys := make(map[string]bool)
	for _, Reg := range Data.RegionTMP {
		if _, value := keys[Reg]; !value {
			keys[Reg] = true
			list = append(list, Reg)
		}
	}
	if s == "add" {
		Data.ChannelState.Region = strings.Join(list, ",")
	} else {
		tmp := []string{}
		for _, v := range list {
			skip := false
			for _, v2 := range Data.DelRegionVal {
				if v2 == v {
					skip = true
					break
				}
			}
			if !skip {
				tmp = append(tmp, v)
			}
		}
		Data.ChannelState.Region = strings.Join(tmp, ",")
	}
}

func (Data *Regis) AddNewRegion(new string) *Regis {
	Data.AddRegionVal = append(Data.AddRegionVal, new)
	Data.RegionTMP = append(Data.RegionTMP, new)
	return Data
}

func (Data *Regis) RemoveRegion(new string) *Regis {
	Data.DelRegionVal = append(Data.DelRegionVal, new)
	return Data
}

func (Data *Regis) UpdateType(new int) *Regis {
	Data.ChannelState.TypeTag = new
	return Data
}

func (Data *Regis) UpdateMessageID(new string) *Regis {
	Data.MessageID = new
	return Data
}

func (Data *Regis) Clear() {
	Register = &Regis{}
	Data = &Regis{}
}

func (Data *Regis) Stop() {
	Data.Gass = false
}

func (Data *Regis) Start() {
	Data.Gass = true
}

func (Data *Regis) BreakPoint(num time.Duration) {
	for i := 0; i < 100; i++ {
		if Data.Gass {
			break
		}
		time.Sleep(num * time.Second)
	}
}
