package main

import (
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

var (
	Register = &ChannelRegister{}
)

const (
	UpdateState = "SelectChannel"
	FirstSetup  = "Setup"
)

func StartRegister(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		CleanRegister()
	}

	if strings.HasPrefix(m.Content, Prefix) {
		Admin, err := MemberHasPermission(m.GuildID, m.Author.ID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
			if err != nil {
				log.Error(err)
			}
		}
		Register = NewRegister()

		if m.Content == Prefix+Setup {
			if Admin {
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to setup mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}

				Register.SetAdmin(m.Author.ID).SetChannelID(m.ChannelID)
				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID or Name of Vtuber group/agency you want to enable (Only one)")
				if err != nil {
					log.Error(err)
				}

				for _, v := range *GroupsPayload {
					table.Append([]string{strconv.Itoa(int(v.ID)), v.GroupName})
				}
				table.SetHeader([]string{"ID", "GroupName"})
				table.Render()
				_, err = s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
				if err != nil {
					log.Error(err)
				}
				Register.UpdateState(FirstSetup)
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Your roles don't have permission to enable/disable/update,make sure your roles have `Manage Channels` permission")
				if err != nil {
					log.Error(err)
				}
				CleanRegister()
				return
			}

			//Update state
		} else if m.Content == Prefix+Update2 {
			if Admin {
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to update mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}

				var (
					Typestr     string
					LiveOnly    = config.No
					NewUpcoming = config.No
					Dynamic     = config.No
					LiteMode    = config.No
					Indie       = ""
					Region      = "All"
				)
				ChannelData, err := database.ChannelStatus(m.ChannelID)
				if err != nil {
					SendError(map[string]string{
						"ChannelID": m.ChannelID,
						"Username":  m.Author.Username,
						"AvatarURL": m.Author.AvatarURL("128"),
					})
				}
				if len(ChannelData) > 0 {
					for _, Channel := range ChannelData {
						if Channel.Region != "" {
							Region = Channel.Region
						}

						if Channel.IsFanart() {
							Typestr = "FanArt"
						}

						if Channel.IsLive() {
							Typestr = "Live"
						}

						if Channel.IsFanart() && Channel.IsLive() {
							Typestr = "FanArt & Livestream"
						}

						if Channel.IsLewd() {
							Typestr = "Lewd"
						}

						if Channel.IsLewd() && Channel.IsFanart() {
							Typestr = "FanArt & Lewd"
						}

						if Channel.LiveOnly {
							LiveOnly = config.Ok
						}

						if Channel.NewUpcoming {
							NewUpcoming = config.Ok
						}

						if Channel.Dynamic {
							Dynamic = config.Ok
						}

						if Channel.LiteMode {
							LiteMode = config.Ok
						}

						if Channel.Group.GroupName == config.Indie {
							if Channel.IndieNotif {
								Indie = config.Ok
							} else if Channel.Group.GroupName != config.Indie {
								Indie = "-"
							} else {
								Indie = config.No
							}
							Channel.Group.RemoveNillIconURL()

							_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetThumbnail(Channel.Group.IconURL).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle("ID "+strconv.Itoa(int(Channel.ID))).
								AddField("Type", Typestr).
								AddField("LiveOnly", LiveOnly).
								AddField("Dynamic", Dynamic).
								AddField("Upcoming", NewUpcoming).
								AddField("Lite", LiteMode).
								AddField("Regions", Region).
								AddField("Independent notif", Indie).
								InlineAllFields().MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						} else {
							_, err = s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetThumbnail(Channel.Group.IconURL).
								SetDescription("Channel States of "+Channel.Group.GroupName).
								SetTitle("ID "+strconv.Itoa(int(Channel.ID))).
								AddField("Type", Typestr).
								AddField("LiveOnly", LiveOnly).
								AddField("Dynamic", Dynamic).
								AddField("Upcoming", NewUpcoming).
								AddField("Lite", LiteMode).
								AddField("Regions", Region).
								InlineAllFields().MessageEmbed)
							if err != nil {
								log.Error(err)
							}
						}
					}
				} else {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetTitle("404 Not found,use `"+Prefix+Setup+"` first").
						SetThumbnail(config.GoSimpIMG).
						SetImage(engine.NotFoundIMG()).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
					return
				}

				Register.SetAdmin(m.Author.ID).SetChannelID(m.ChannelID)
				Register.SetChannels(ChannelData)

				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID : ")
				if err != nil {
					log.Error(err)
				}

				Register.UpdateState(UpdateState)
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
				if err != nil {
					log.Error(err)
				}
				CleanRegister()
				return
			}
		}
	} else if m.Author.ID == Register.Admin && m.ChannelID == Register.ChannelState.ChannelID {
		if m.Content == "exit" {
			Out()
			return
		}
		if Register.State == UpdateState {
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
						Register.SetChannel(ChannelState)
					}
				}
				if Register.ChannelState.ID != 0 {
					Register.SetChannelID(m.ChannelID)
					_, err := s.ChannelMessageSend(m.ChannelID, "You selectd `"+Register.ChannelState.Group.GroupName+"` with ID `"+strconv.Itoa(int(Register.ChannelState.ID))+"`")
					if err != nil {
						log.Error(err)
					}
					table.SetHeader([]string{"Menu"})
					table.Append([]string{"Update Channel state"})
					table.Append([]string{"Add region in this channel"})
					table.Append([]string{"Delete region in this channel"})

					if Register.ChannelState.TypeTag == 2 || Register.ChannelState.TypeTag == 3 {
						table.Append([]string{"Change Livestream state"})
					}

					table.Render()
					MsgText, err := s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
					if err != nil {
						log.Error(err)
					}

					if Register.ChannelState.TypeTag == 2 || Register.ChannelState.TypeTag == 3 {
						err = engine.Reacting(map[string]string{
							"ChannelID": m.ChannelID,
							"State":     "Menu2",
							"MessageID": MsgText.ID,
						}, s)
						if err != nil {
							log.Error(err)
						}
					} else {
						err = engine.Reacting(map[string]string{
							"ChannelID": m.ChannelID,
							"State":     "Menu",
							"MessageID": MsgText.ID,
						}, s)
						if err != nil {
							log.Error(err)
						}
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

		//Fist Setup
		if Register.State == FirstSetup {
			VTuberGroup, err := FindGropName(m.Content)
			if err != nil {
				_, err := s.ChannelMessageSend(m.ChannelID, "`"+m.Content+"`,Name of Vtuber Group was not valid")
				if err != nil {
					log.Error(err)
				}
				return
			}
			Register.SetGroup(VTuberGroup)

			if Register.ChannelState.Group.IsNull() {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetDescription("Invalid ID,Group not found").
					SetImage(engine.NotFoundIMG()).MessageEmbed)
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
						Register.UpdateState(AddRegion)
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

			Register.ChoiceType()
			Register.BreakPoint(3)

			if Register.ChannelState.TypeTag == 3 || Register.ChannelState.TypeTag == 2 {
				Register.Stop()
				Register.LiveOnly()
				Register.BreakPoint(1)

				if !Register.ChannelState.LiveOnly {
					Register.Stop()
					Register.NewUpcoming()
					Register.BreakPoint(1)
				}

				Register.Stop()
				Register.Dynamic()
				Register.BreakPoint(1)

				Register.Stop()
				Register.Lite()
				Register.BreakPoint(1)
			} else if Register.ChannelState.TypeTag == 69 || Register.ChannelState.TypeTag == 70 {
				if !Register.CheckNSFW() {
					Out()
				}
			}

			if Register.ChannelState.Group.GroupName == config.Indie {
				Register.Stop()
				Register.IndieNotif()
				Register.BreakPoint(1)
			}

			if Register.ChannelState.Group.GroupName != "" {
				err = Register.ChannelState.AddChannel()
				if err != nil {
					log.Error(err)
				}

				_, err = s.ChannelMessageSend(m.ChannelID, "Done,you add `"+Register.ChannelState.Group.GroupName+"` in this channel")
				if err != nil {
					log.Error(err)
				}
			}

			CleanRegister()
			return
		}
	}
}
