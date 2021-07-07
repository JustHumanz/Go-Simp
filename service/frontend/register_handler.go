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
	Register = &NoRegister{}
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
	Out := func(j int) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Adios")
		if err != nil {
			log.Error(err)
		}
		CleanRegister(j)
	}

	if strings.HasPrefix(m.Content, Prefix) {
		Admin, err := MemberHasPermission(m.GuildID, m.Author.ID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, err.Error())
			if err != nil {
				log.Error(err)
			}
		}

		if m.Content == Prefix+Setup {
			if Admin {
				RegisterPayload := NewRegister(m.Author.ID, m.ChannelID)
				_, err := s.ChannelMessageSend(m.ChannelID, "Wellcome to setup mode\ntype `exit` to exit this mode")
				if err != nil {
					log.Error(err)
				}

				RegisterPayload.SetAdmin(m.Author.ID).SetChannelID(m.ChannelID)
				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID or Name of Vtuber group/agency you want to enable (Only one)")
				if err != nil {
					log.Error(err)
				}

				for _, v := range *GroupsPayload {
					table.Append([]string{strconv.Itoa(int(v.ID)), v.GroupName})
				}
				table.SetHeader([]string{"ID", "GroupName"})
				table.Render()
				_, err = s.ChannelMessageSend(m.ChannelID, "```"+tableString.String()+"```")
				if err != nil {
					log.Error(err)
				}
				RegisterPayload.UpdateState(FirstSetup)
				Register.Payload = append(Register.Payload, &RegisterPayload)
				return

			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Your roles don't have permission to enable/disable/update,make sure your roles have `Manage Channels` permission")
				if err != nil {
					log.Error(err)
				}
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

				RegisterPayload := NewRegister(m.Author.ID, m.ChannelID)

				RegisterPayload.SetAdmin(m.Author.ID).SetChannelID(m.ChannelID)
				RegisterPayload.SetChannels(ChannelData)

				_, err = s.ChannelMessageSend(m.ChannelID, "Select ID : ")
				if err != nil {
					log.Error(err)
				}
				RegisterPayload.UpdateState(UpdateState)
				Register.Payload = append(Register.Payload, &RegisterPayload)
				return
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}

	for j, RegisterPayload := range Register.Payload {
		if RegisterPayload.AdminID == m.Author.ID+m.ChannelID {
			if m.Content == "exit" {
				Out(j)
				return
			}
			if RegisterPayload.State == UpdateState {
				tmp, err := strconv.Atoi(m.Content)
				if err != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, "Worng input ID")
					if err != nil {
						log.Error(err)
					}
					Out(j)
					return
				} else {
					for _, ChannelState := range RegisterPayload.ChannelStates {
						if int(ChannelState.ID) == tmp {
							RegisterPayload.SetChannel(ChannelState)
						}
					}
					if RegisterPayload.ChannelState.ID != 0 {
						RegisterPayload.SetChannelID(m.ChannelID)
						_, err := s.ChannelMessageSend(m.ChannelID, "You selectd `"+RegisterPayload.ChannelState.Group.GroupName+"` with ID `"+strconv.Itoa(int(RegisterPayload.ChannelState.ID))+"`")
						if err != nil {
							log.Error(err)
						}
						table.SetHeader([]string{"Menu"})
						table.Append([]string{"Update Channel state"})
						table.Append([]string{"Add region in this channel"})
						table.Append([]string{"Delete region in this channel"})

						if RegisterPayload.ChannelState.TypeTag == 2 || RegisterPayload.ChannelState.TypeTag == 3 {
							table.Append([]string{"Change Livestream state"})
						}

						table.Render()
						MsgText, err := s.ChannelMessageSend(m.ChannelID, "`"+tableString.String()+"`")
						if err != nil {
							log.Error(err)
						}

						if RegisterPayload.ChannelState.TypeTag == 2 || RegisterPayload.ChannelState.TypeTag == 3 {
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

						RegisterPayload.UpdateMessageID(MsgText.ID)

					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "Channel ID not found")
						if err != nil {
							log.Error(err)
						}
						Out(j)
						return
					}
				}
			}

			//Fist Setup
			if RegisterPayload.State == FirstSetup {
				VTuberGroup, err := FindGropName(m.Content)
				if err != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, "`"+m.Content+"`,Name of Vtuber Group was not valid")
					if err != nil {
						log.Error(err)
					}
					return
				}
				RegisterPayload.SetGroup(VTuberGroup)

				if RegisterPayload.ChannelState.Group.IsNull() {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetDescription("Invalid ID,Group not found").
						SetImage(engine.NotFoundIMG()).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
					Out(j)
					return
				}

				RegisterPayload.Stop()
				for Key, Val := range RegList {
					if Key == RegisterPayload.ChannelState.Group.GroupName {
						if len(Val) > 3 {
							MsgTxt, err := s.ChannelMessageSend(m.ChannelID, "Select `"+Key+"` region")
							if err != nil {
								log.Error(err)
							}
							RegisterPayload.UpdateState(AddRegion)
							for _, v := range strings.Split(Val, ",") {
								err := s.MessageReactionAdd(m.ChannelID, MsgTxt.ID, engine.CountryCodetoUniCode(v))
								if err != nil {
									log.Error(err)
								}
							}
							RegisterPayload.UpdateMessageID(MsgTxt.ID)
							RegisterPayload.BreakPoint(5)

							RegisterPayload.FixRegion("add")
							if RegisterPayload.ChannelState.ChannelCheck() {
								_, err := s.ChannelMessageSend(m.ChannelID, "Already setup `"+RegisterPayload.ChannelState.Group.GroupName+"`,for add/del region use `Update` command")
								if err != nil {
									log.Error(err)
								}
								Out(j)
								return
							}
							RegisterPayload.Stop()
						} else {
							RegisterPayload.RegionTMP = strings.Split(Val, ",")
							RegisterPayload.FixRegion("add")
						}
					}
				}

				RegisterPayload.ChoiceType()
				RegisterPayload.BreakPoint(3)

				if RegisterPayload.ChannelState.TypeTag == 3 || RegisterPayload.ChannelState.TypeTag == 2 {
					RegisterPayload.Stop()
					RegisterPayload.LiveOnly()
					RegisterPayload.BreakPoint(1)

					if !RegisterPayload.ChannelState.LiveOnly {
						RegisterPayload.Stop()
						RegisterPayload.NewUpcoming()
						RegisterPayload.BreakPoint(1)
					}

					RegisterPayload.Stop()
					RegisterPayload.Dynamic()
					RegisterPayload.BreakPoint(1)

					RegisterPayload.Stop()
					RegisterPayload.Lite()
					RegisterPayload.BreakPoint(1)
				} else if RegisterPayload.ChannelState.TypeTag == 69 || RegisterPayload.ChannelState.TypeTag == 70 {
					if !RegisterPayload.CheckNSFW() {
						Out(j)
					}
				}

				if RegisterPayload.ChannelState.Group.GroupName == config.Indie {
					RegisterPayload.Stop()
					RegisterPayload.IndieNotif()
					RegisterPayload.BreakPoint(1)
				}

				if RegisterPayload.ChannelState.Group.GroupName != "" {
					err = RegisterPayload.ChannelState.AddChannel()
					if err != nil {
						log.Error(err)
					}

					_, err = s.ChannelMessageSend(m.ChannelID, "Done,you add `"+RegisterPayload.ChannelState.Group.GroupName+"` in this channel")
					if err != nil {
						log.Error(err)
					}
				}

				CleanRegister(j)
				return
			}
		}
	}
}
