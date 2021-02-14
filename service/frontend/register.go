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

var (
	RegisterValue = &Regis{}
	OK            = false
)

func Answer(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == RegisterValue.Admin && m.MessageID == RegisterValue.MessageID {
		OK = true
		if m.Emoji.MessageFormat() == config.Ok {
			if RegisterValue.State == "LiveOnly" {
				RegisterValue.SetLiveOnly(true)
			} else if RegisterValue.State == "NewUpcoming" {
				RegisterValue.SetNewUpcoming(true)
			} else if RegisterValue.State == "Dynamic" {
				RegisterValue.SetDynamic(true)
			}
		} else if m.Emoji.MessageFormat() == config.No {
			if RegisterValue.State == "LiveOnly" {
				RegisterValue.SetLiveOnly(false)
			} else if RegisterValue.State == "NewUpcoming" {
				RegisterValue.SetNewUpcoming(false)
			} else if RegisterValue.State == "Dynamic" {
				RegisterValue.SetDynamic(false)
			}
		}

		if m.Emoji.MessageFormat() == config.Art {
			if RegisterValue.ChannelState.TypeTag == 2 {
				RegisterValue.UpdateType(3)
			} else {
				RegisterValue.UpdateType(1)
			}
		} else if m.Emoji.MessageFormat() == config.Live {
			if RegisterValue.ChannelState.TypeTag == 1 {
				RegisterValue.UpdateType(3)
			} else {
				RegisterValue.UpdateType(2)
			}
		}

		if m.Emoji.MessageFormat() == config.One {
			RegisterValue.LiveOnly(s)

			if !RegisterValue.ChannelState.LiveOnly {
				RegisterValue.NewUpcoming(s)
			}

			RegisterValue.Dynamic(s)

			RegisterValue.UpdateChannel()
			_, err := s.ChannelMessageSend(m.ChannelID, "Done")
			if err != nil {
				log.Error(err)
			}
			return

		} else if m.Emoji.MessageFormat() == config.Two {
			RegisterValue.AddRegion(s)

		} else if m.Emoji.MessageFormat() == config.Three {
			RegisterValue.DelRegion(s)
		}

		if RegisterValue.State == "AddReg" {
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				RegisterValue.AddNewRegion(Region)
			}
		} else if RegisterValue.State == "DelReg" {
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				RegisterValue.RemoveRegion(Region)
			}

		}
	}
}

func Register(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		RegisterValue.Clear()
	}

	if strings.HasPrefix(m.Content, Prefix) {
		if CheckPermission(m.Author.ID, m.ChannelID, s) {
			if m.Content == Prefix+"setup" {
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to setup mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}
				if database.CheckIfNewChannel(m.ChannelID) {
					RegisterValue.SetAdmin(m.Author.ID).SetChannel(m.ChannelID)
					_, err := s.ChannelMessageSend(m.ChannelID, "Select ID of Vtuber group/agency you want to enable (Only one)")
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
					RegisterValue.UpdateState("Group")

				}

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
					RegisterValue.SetAdmin(m.Author.ID).SetChannel(m.ChannelID)

					RegisterValue.ChannelStates = ChannelData
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
				RegisterValue.UpdateState("SelectChannel")
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
			RegisterValue.Clear()
			return
		}
	} else if m.Author.ID == RegisterValue.Admin && m.ChannelID == RegisterValue.ChannelState.ChannelID {
		if m.Content == "exit" {
			Out()
			return
		}
		if RegisterValue.State == "SelectChannel" {
			tmp, err := strconv.Atoi(m.Content)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "Worng input ID")
				if err != nil {
					log.Error(err)
				}
				Out()
				return
			} else {
				for _, ChannelState := range RegisterValue.ChannelStates {
					if int(ChannelState.ID) == tmp {
						RegisterValue.ChannelState = ChannelState
					}
				}
				if RegisterValue.ChannelState.ID != 0 {
					RegisterValue.SetChannel(m.ChannelID)
					_, err := s.ChannelMessageSend(m.ChannelID, "You select `"+RegisterValue.ChannelState.Group.GroupName+"` with ID `"+strconv.Itoa(int(RegisterValue.ChannelState.ID))+"`")
					if err != nil {
						log.Error(err)
					}
					table.SetHeader([]string{"ID", "Func"})
					table.Append([]string{"1", "Update Channel state"})
					table.Append([]string{"2", "Add region in this channel"})
					table.Append([]string{"3", "Delete region in this channel"})
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
					RegisterValue.UpdateMessageID(MsgText.ID)

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
		if RegisterValue.State == "Group" {
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
					RegisterValue.SetGroup(v)
					break
				}
			}
			//RegisterValue.ChannelState.ChannelCheck()
			if RegisterValue.ChannelState.Group.IsNull() {
				_, err := s.ChannelMessageSend(m.ChannelID, "Invalid ID,Group not found")
				if err != nil {
					log.Error(err)
				}
				Out()
				return
			}
			for Key, Val := range RegList {
				if Key == RegisterValue.ChannelState.Group.GroupName {
					MsgTxt, err := s.ChannelMessageSend(m.ChannelID, "Select `"+Key+"` region")
					if err != nil {
						log.Error(err)
					}
					for _, v := range strings.Split(Val, ",") {
						err := s.MessageReactionAdd(m.ChannelID, MsgTxt.ID, engine.CountryCodetoUniCode(v))
						if err != nil {
							log.Error(err)
						}
					}
					RegisterValue.UpdateMessageID(MsgTxt.ID)
					OK = false

					RegisterValue.FixRegion()
					if RegisterValue.ChannelState.ChannelCheck() {
						_, err := s.ChannelMessageSend(m.ChannelID, "Already setup `"+RegisterValue.ChannelState.Group.GroupName+"`,for add/del region use `Update` command")
						if err != nil {
							log.Error(err)
						}
						Out()
						return
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
					RegisterValue.UpdateMessageID(MsgText.ID)
					for i := 0; i < 100; i++ {
						if OK {
							break
						}
						time.Sleep(1 * time.Second)
					}
					if RegisterValue.ChannelState.TypeTag != 1 {
						RegisterValue.LiveOnly(s)
						for i := 0; i < 100; i++ {
							if OK {
								break
							}
							time.Sleep(1 * time.Second)
						}

						if !RegisterValue.ChannelState.LiveOnly {
							RegisterValue.NewUpcoming(s)
							for i := 0; i < 100; i++ {
								if OK {
									break
								}
								time.Sleep(1 * time.Second)
							}
						}

						RegisterValue.Dynamic(s)
						for i := 0; i < 100; i++ {
							if OK {
								break
							}
							time.Sleep(1 * time.Second)
						}
					}
					err = RegisterValue.ChannelState.AddChannel()
					if err != nil {
						log.Error(err)
					}
					_, err = s.ChannelMessageSend(m.ChannelID, "Done,you add `"+RegisterValue.ChannelState.Group.GroupName+"` in this channel")
					if err != nil {
						log.Error(err)
					}
					RegisterValue.Clear()
					return
				}
			}
		}
	}
}

type Regis struct {
	Admin         string
	State         string
	MessageID     string
	RegionTMP     []string
	ChannelState  database.DiscordChannel
	ChannelStates []database.DiscordChannel
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
	GroupName := RegisterValue.ChannelState.Group.GroupName
	ChannelID := RegisterValue.ChannelState.ChannelID

	RegisterValue.UpdateState("AddReg")
	RegEmoji := []string{}
	for _, v := range strings.Split(RegisterValue.ChannelState.Region, ",") {
		RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v))
		RegisterValue.RegionTMP = append(RegisterValue.RegionTMP, v)
	}
	_, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region you already enabled in here "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgTxt2, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Select region you want to add it : ")
	if err != nil {
		log.Error(err)
	}

	RegisterValue.UpdateMessageID(MsgTxt2.ID)
	for Key, Val := range RegList {
		if Key == GroupName {
			if len(RegisterValue.RegionTMP) == len(strings.Split(Val, ",")) {
				_, err := s.ChannelMessageSend(ChannelID, "You already enable all region")
				if err != nil {
					log.Error(err)
				}
				return
			}

			for _, v := range RegisterValue.RegionTMP {
				for _, v2 := range strings.Split(Val, ",") {
					if v != v2 {
						err := s.MessageReactionAdd(ChannelID, MsgTxt2.ID, engine.CountryCodetoUniCode(v2))
						if err != nil {
							log.Error(err)
						}
					}
				}
			}
			OK = false
		}
	}

	for i := 100 - 1; i >= 0; i-- {
		if OK {
			break
		}
		time.Sleep(1 * time.Second)
	}
	RegisterValue.FixRegion()
	RegisterValue.ChannelState.UpdateChannel("Region")

	_, err = s.ChannelMessageSend(ChannelID, "Done,you adding "+RegisterValue.ChannelState.Region)
	if err != nil {
		log.Error(err)
	}
	return
}

func (Data *Regis) DelRegion(s *discordgo.Session) {
	GroupName := RegisterValue.ChannelState.Group.GroupName
	ChannelID := RegisterValue.ChannelState.ChannelID

	RegisterValue.UpdateState("DelReg")
	RegEmoji := []string{}
	for Key, Val := range RegList {
		if Key == GroupName {
			for _, v := range strings.Split(RegisterValue.ChannelState.Region, ",") {
				for _, v2 := range strings.Split(Val, ",") {
					if v == v2 {
						RegEmoji = append(RegEmoji, engine.CountryCodetoUniCode(v2))
						RegisterValue.RegionTMP = append(RegisterValue.RegionTMP, v2)
					}
				}
			}
		}
	}

	_, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Region you enabled in here "+strings.Join(RegEmoji, "  "))
	if err != nil {
		log.Error(err)
	}

	MsgTxt2, err := s.ChannelMessageSend(ChannelID, "`"+GroupName+"` Select region you want to delete it : ")
	if err != nil {
		log.Error(err)
	}
	for _, v := range RegisterValue.RegionTMP {
		err := s.MessageReactionAdd(ChannelID, MsgTxt2.ID, engine.CountryCodetoUniCode(v))
		if err != nil {
			log.Error(err)
		}
	}
	OK = false
	RegisterValue.UpdateMessageID(MsgTxt2.ID)

	for i := 100 - 1; i >= 0; i-- {
		if OK {
			break
		}
		time.Sleep(1 * time.Second)
	}
	RegisterValue.FixRegion()
	RegisterValue.ChannelState.UpdateChannel("Region")

	_, err = s.ChannelMessageSend(ChannelID, "Done,you remove "+RegisterValue.ChannelState.Region)
	if err != nil {
		log.Error(err)
	}
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

func (Data *Regis) FixRegion() {
	Data.ChannelState.Region = strings.Join(Data.RegionTMP, ",")
}

func (Data *Regis) AddNewRegion(new string) *Regis {
	Data.RegionTMP = append(Data.RegionTMP, new)
	return Data
}

func (Data *Regis) RemoveRegion(new string) *Regis {
	for _, v := range strings.Split(Data.ChannelState.Region, ",") {
		if v != new {
			Data.RegionTMP = append(Data.RegionTMP, new)
		}
	}
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
	RegisterValue = &Regis{}
	Data = &Regis{}
}
