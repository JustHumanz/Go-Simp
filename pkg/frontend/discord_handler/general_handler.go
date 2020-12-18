package discordhandler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

//Fanart discord message handler
func Fanart(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PFanart
	var (
		Member      bool
		Group       bool
		Pic         = config.NotFound
		Msg         string
		embed       *discordgo.MessageEmbed
		DynamicData DynamicSvr
	)

	if strings.HasPrefix(m.Content, Prefix) {
		SendNude := func(Title, Author, Text, URL, Pic, Msg string, Color int, State, Dynamic string) bool {
			Msg = Msg + " *sometimes image not showing,because image oversize*"
			if State == "TBiliBili" {
				var (
					body    []byte
					errcurl error
					urls    = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=" + Dynamic
				)
				body, errcurl = network.Curl(urls, nil)
				if errcurl != nil {
					log.Error(errcurl, string(body))
					log.Info("Trying use tor")

					body, errcurl = network.CoolerCurl(urls, nil)
					if errcurl != nil {
						log.Error(errcurl)
					}
				}
				json.Unmarshal(body, &DynamicData)
				embed = engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Author).
					SetThumbnail(DynamicData.GetUserAvatar()).
					SetDescription(Text).
					SetURL(URL).
					SetImage(Pic).
					SetColor(Color).
					InlineAllFields().
					SetFooter(Msg, config.BiliBiliIMG).MessageEmbed
			} else {
				embed = engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Author).
					SetThumbnail(GetAuthorAvatar(Author)).
					SetDescription(Text).
					SetURL(URL).
					SetImage(Pic).
					SetColor(Color).
					InlineAllFields().
					SetFooter(Msg, config.TwitterIMG).MessageEmbed
			}
			msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				log.Error(err, msg)
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": m.ChannelID,
				"Content":   m.Content,
				"Prefix":    Prefix,
			}, s)
			if err != nil {
				log.Error(err)
			}
			return true
		}
		for _, GroupData := range engine.GroupData {
			if m.Content == strings.ToLower(Prefix+GroupData.GroupName) {
				Color, err := engine.GetColor("/tmp/mem.tmp", m.Author.AvatarURL("128"))
				if err != nil {
					log.Error(err)
				}

				DataFix := database.GetFanart(GroupData.ID, 0)
				_, err = network.Curl(DataFix.PermanentURL, nil)
				if err != nil {
					_, err = network.CoolerCurl(DataFix.PermanentURL, nil)
					if err != nil {
						log.Error(err)
					}
					log.Info("Delete fanart metadata ", DataFix.PermanentURL)
					err = DataFix.DeleteFanart()
					if err != nil {
						log.Error(err)
					}
				} else {
					if DataFix.Videos != "" {
						Msg = "Video type,check original post"
						Pic = config.NotFound
					} else if len(DataFix.Photos) > 0 {
						Pic = DataFix.Photos[0]
					}
					Group = SendNude(engine.FixName(DataFix.EnName, DataFix.JpName),
						DataFix.Author, RemovePic(DataFix.Text),
						DataFix.PermanentURL,
						Pic, Msg, Color,
						DataFix.State, DataFix.Dynamic_id)
					break
				}
			}
			for _, MemberData := range database.GetMembers(GroupData.ID) {
				if m.Content == strings.ToLower(Prefix+MemberData.Name) || m.Content == strings.ToLower(Prefix+MemberData.JpName) {
					Color, err := engine.GetColor("/tmp/mem.tmp", m.Author.AvatarURL("128"))
					if err != nil {
						log.Error(err)
					}

					DataFix := database.GetFanart(0, MemberData.ID)
					_, err = network.Curl(DataFix.PermanentURL, nil)
					if err != nil {
						_, err = network.CoolerCurl(DataFix.PermanentURL, nil)
						if err != nil {
							log.Error(err)
						}
						log.Info("Delete fanart metadata ", DataFix.PermanentURL)
						err = DataFix.DeleteFanart()
						if err != nil {
							log.Error(err)
						}
					} else {
						if DataFix.Videos != "" {
							Msg = "Video type,check original post"
							Pic = config.NotFound
						} else if len(DataFix.Photos) > 0 {
							Pic = DataFix.Photos[0]
						}
						Member = SendNude(engine.FixName(MemberData.EnName, MemberData.JpName),
							DataFix.Author, RemovePic(DataFix.Text),
							DataFix.PermanentURL,
							Pic, Msg, Color,
							DataFix.State, DataFix.Dynamic_id)
						break
					}
				}
			}
		}
		if Member || Group {
			return
		}
		if !Group && !Member {
			s.ChannelMessageSend(m.ChannelID, "`"+m.Content[len(Prefix):]+"` was invalid name")
		}
	}
}

//Tags command message handler
func Tags(s *discordgo.Session, m *discordgo.MessageCreate) {
	Prefix := config.PGeneral
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, Prefix) {
		var (
			Already     []string
			Done        []string
			MemberTag   []NameStruct
			ReminderInt = 0
		)
		User := &database.UserStruct{
			DiscordID:       m.Author.ID,
			DiscordUserName: m.Author.Username,
			Channel_ID:      m.ChannelID,
			Human:           true,
			Reminder:        ReminderInt,
		}
		Color, err := engine.GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}
		if strings.HasPrefix(m.Content, Prefix+TagMe) {
			Already = nil
			Done = nil
			UserInput := strings.Replace(m.Content, Prefix+TagMe, "", -1)
			var (
				VtuberName   string
				ReminderUser int
				re           = regexp.MustCompile(`(?m)-setreminder\s[1-9].`)
			)

			if len(re.FindAllString(UserInput, -1)) > 0 {
				tmpvar := re.FindAllString(UserInput, -1)[0]
				if tmpvar[len(tmpvar)-2:] != "" {
					tmpvar2 := tmpvar[len(tmpvar)-2:]
					tmpvar3, err := strconv.Atoi(tmpvar2[:len(tmpvar2)-1] + "6")
					if err != nil {
						log.Error(err)
					} else {
						ReminderUser = tmpvar3
					}
				} else {
					tmpvar2, err := strconv.Atoi(tmpvar + "6")
					if err != nil {
						log.Error(err)
					} else {
						ReminderUser = tmpvar2
					}
				}

				fmt.Println(ReminderUser)
				if ReminderUser > 67 {
					s.ChannelMessageSend(m.ChannelID, "Can't set Reminder over than 60 Minutes")
					return
				}
				VtuberName = strings.TrimSpace(strings.Replace(UserInput, tmpvar, "", -1))

			} else {
				VtuberName = strings.TrimSpace(strings.Replace(UserInput, "-setreminder", "", -1))
			}
			if VtuberName != "" {
				tmp := strings.Split(VtuberName, ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data.GroupID == 0 {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID).
								SetReminder(ReminderUser)

							for _, Member := range database.GetMembers(VTuberGroup.ID) {
								err := User.Adduser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
								}
							}
							if Already != nil {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your tag list").
									AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
							if Done != nil {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You Add\n"+strings.Join(Done, " ")+" to your tag list").
									AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.GroupName+"`")
							if err != nil {
								log.Error(err)
							}
							return
						}
					} else {
						MemberTag = append(MemberTag, Data)
					}
					Already = nil
					Done = nil
				}
				for i, Member := range MemberTag {
					if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
						User.SetGroupID(Member.GroupID).
							SetReminder(ReminderUser)

						err := User.Adduser(Member.MemberID)
						if err != nil {
							Already = append(Already, "`"+tmp[i]+"`")
						} else {
							Done = append(Done, "`"+tmp[i]+"`")
						}
					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
				if Already != nil {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your list").
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}

				}
				if Done != nil {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You Add\n"+strings.Join(Done, " ")+" to your list").
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+TagMe+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.HasPrefix(m.Content, Prefix+SetReminder) {
			var (
				ReminderUser int
				UserInput    = strings.Replace(m.Content, Prefix+SetReminder, "", -1)
				FindInt      = strings.Split(UserInput, " ")
			)

			if UserInput != "" {
				if len(FindInt) > 2 {
					tmpvar := FindInt[2]
					if lenstr, _ := regexp.MatchString("^[0-9]{2}[:.,-]?$", tmpvar); lenstr {
						tmpvar3, err := strconv.Atoi(tmpvar[:len(tmpvar)-1] + "6")
						if err != nil {
							log.Error(err)
						}
						ReminderUser = tmpvar3
						if ReminderUser > 67 {
							_, err := s.ChannelMessageSend(m.ChannelID, "Can't set Reminder over than 60 Minutes")
							if err != nil {
								log.Error(err)
							}
							return
						}

					} else if FindInt[2] == "0" {
						_, err := s.ChannelMessageSend(m.ChannelID, "You disable reminder time")
						if err != nil {
							log.Error(err)
						}
						ReminderUser = 0
					} else {
						//Invaild
						_, err := s.ChannelMessageSend(m.ChannelID, "Invaild number\nmin 10,max 60")
						if err != nil {
							log.Error(err)
						}
						return
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Invalid `"+SetReminder+"` command")
					if err != nil {
						log.Error(err)
					}
					return
				}

				tmp := strings.Split(FindInt[1], ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data.GroupID == 0 {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
							if err != nil {
								log.Error(err)
							}
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID).
								SetReminder(ReminderUser)
							for _, Member := range database.GetMembers(VTuberGroup.ID) {
								err = User.UpdateReminder(Member.ID)
								if err != nil {
									log.Error(err)
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
								}
							}
							if Done != nil {
								if ReminderUser == 0 {
									_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Disable reminder time\n"+strings.Join(Done, " ")).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
										SetColor(Color).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								} else {
									_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Update reminder time\n"+strings.Join(Done, " ")).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
										SetColor(Color).MessageEmbed)
									if err != nil {
										log.Error(err)
									}
								}
								Done = nil
							} else if Already != nil {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You not tag \n"+strings.Join(Already, " ")).
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
								Already = nil
							}
						} else {
							_, err := s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.GroupName+"`")
							if err != nil {
								log.Error(err)
							}
							return
						}
					} else {
						MemberTag = append(MemberTag, Data)
					}
				}
				for i, Member := range MemberTag {
					if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
						User.SetGroupID(Member.GroupID).
							SetReminder(ReminderUser)
						err := User.UpdateReminder(Member.MemberID)
						if err != nil {
							log.Error(err)
							Already = append(Already, "`"+tmp[i]+"`")
						} else {
							Done = append(Done, "`"+tmp[i]+"`")
						}

					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
				if Done != nil {
					if ReminderUser == 0 {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Disable reminder time\n"+strings.Join(Done, " ")).
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Update reminder time\n"+strings.Join(Done, " ")).
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					}
				} else if Already != nil {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You not tag \n"+strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}

				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+SetReminder+"` command")
				if err != nil {
					log.Error(err)
				}
				return
			}
		} else if strings.HasPrefix(m.Content, Prefix+DelTag) {
			Already = nil
			Done = nil
			VtuberName := strings.TrimSpace(strings.Replace(m.Content, Prefix+DelTag, "", -1))
			if VtuberName != "" {
				tmp := strings.Split(VtuberName, ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data == (NameStruct{}) {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							_, err := s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
							if err != nil {
								log.Error(err)
							}
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID)
							for _, Member := range database.GetMembers(VTuberGroup.ID) {
								err := User.Deluser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
								}
							}
							if Already != nil {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
									AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
							if Done != nil {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
									AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									InlineAllFields().
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+VTuberGroup.GroupName+"`").
								SetImage(VTuberGroup.IconURL).
								SetThumbnail(config.GoSimpIMG).
								SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}

					} else {
						MemberTag = append(MemberTag, Data)
					}
				}
				Already = nil
				Done = nil
				for i, Member := range MemberTag {
					if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
						User.GroupID = Member.GroupID
						err := User.Deluser(Member.MemberID)
						if err != nil {
							Already = append(Already, "`"+tmp[i]+"`")
						} else {
							Done = append(Done, "`"+tmp[i]+"`")
						}
					} else {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						return
					}
				}

				if Already != nil {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}

				if Done != nil {
					//return
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+DelTag+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.HasPrefix(m.Content, Prefix+TagRoles) {
			if CheckPermission(m.Author.ID, m.ChannelID, s) {
				Already = nil
				Done = nil
				UserInput := strings.Replace(m.Content, Prefix+TagRoles, "", -1)
				ReminderUser := 0
				re := regexp.MustCompile(`(?m)-setreminder\s[1-9].`)

				if len(re.FindAllString(UserInput, -1)) > 0 {
					tmpvar := re.FindAllString(UserInput, -1)[0]
					if tmpvar[len(tmpvar)-2:] != "" {
						tmpvar2 := tmpvar[len(tmpvar)-2:]
						tmpvar3, err := strconv.Atoi(tmpvar2[:len(tmpvar2)-1] + "6")
						if err != nil {
							log.Error(err)
						} else {
							ReminderUser = tmpvar3
						}
					} else {
						tmpvar2, err := strconv.Atoi(tmpvar + "6")
						if err != nil {
							log.Error(err)
						} else {
							ReminderUser = tmpvar2
						}
					}

					if ReminderUser > 67 {
						s.ChannelMessageSend(m.ChannelID, "Can't set Reminder over than 60 Minutes")
						return
					}
					UserInput = strings.TrimSpace(strings.Replace(UserInput, tmpvar, "", -1))

				} else {
					UserInput = strings.TrimSpace(strings.Replace(UserInput, "-setreminder", "", -1))
				}

				VtuberName := strings.Split(strings.TrimSpace(UserInput), " ")
				guild, err := s.Guild(m.GuildID)
				if err != nil {
					log.Error(err)
				}

				if len(VtuberName[len(VtuberName)-1:]) > 0 {
					tmp := strings.Split(VtuberName[len(VtuberName)-1:][0], ",")
					for _, Name := range tmp {
						Data := FindName(Name)
						if Data.GroupID == 0 {
							VTuberGroup, err := FindGropName(Name)
							if err != nil {
								log.Error(err)
								_, err := s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
								if err != nil {
									log.Error(err)
								}
								return
							}

							if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
								for _, Role := range guild.Roles {
									for _, UserRole := range VtuberName {
										if UserRole == Role.Mention() {
											for _, Member := range database.GetMembers(VTuberGroup.ID) {
												User := database.UserStruct{
													DiscordID:       Role.ID,
													DiscordUserName: Role.Name,
													Channel_ID:      m.ChannelID,
													GroupID:         VTuberGroup.ID,
													Human:           false,
													Reminder:        ReminderUser,
												}
												err := User.Adduser(Member.ID)
												if err != nil {
													Already = append(Already, "`"+Member.Name+"`")
												} else {
													Done = append(Done, "`"+Member.Name+"`")

												}
											}
											if Already != nil {
												_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
													AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to role tags list").
													InlineAllFields().
													SetColor(Color).MessageEmbed)
												if err != nil {
													log.Error(err)
												}
												Already = nil
											}
											if Done != nil {
												_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+"Add\n"+strings.Join(Done, " ")).
													AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
													InlineAllFields().
													SetColor(Color).MessageEmbed)
												if err != nil {
													log.Error(err)
												}
												Done = nil
											}
										}
									}
								}
							} else {
								_, err := s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.GroupName+"`")
								if err != nil {
									log.Error(err)
								}
								return
							}
						} else {
							MemberTag = append(MemberTag, Data)
						}
						Already = nil
						Done = nil
					}
					for i, Member := range MemberTag {
						if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
							for _, Role := range guild.Roles {
								for _, UserRole := range VtuberName {
									if UserRole == Role.Mention() {
										User := database.UserStruct{
											DiscordID:       Role.ID,
											DiscordUserName: Role.Name,
											Channel_ID:      m.ChannelID,
											GroupID:         Member.GroupID,
											Human:           false,
											Reminder:        ReminderUser,
										}
										err := User.Adduser(Member.MemberID)
										if err != nil {
											Already = append(Already, "`"+tmp[i]+"`")
										} else {
											Done = append(Done, "`"+tmp[i]+"`")
										}

										if Already != nil {
											_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
												SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
												SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
												SetThumbnail(config.GoSimpIMG).
												SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
												SetColor(Color).MessageEmbed)
											if err != nil {
												log.Error(err)
											}
											Already = nil
										}
										if Done != nil {
											_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
												SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
												SetDescription(Role.Mention()+" Add\n"+strings.Join(Done, " ")).
												SetThumbnail(config.GoSimpIMG).
												SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show you tags list").
												SetColor(Color).MessageEmbed)
											if err != nil {
												log.Error(err)
											}
											Done = nil
										}
									}
								}
							}

						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+Member.GroupName+"`").
								SetThumbnail(config.GoSimpIMG).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+TagRoles+"` command")
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "You don't have enough permission to use this command,Only user with permission `Manage Channel or Higher` can use this commandYou don't have enough permission to use this command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.HasPrefix(m.Content, Prefix+DelRoles) {
			if CheckPermission(m.Author.ID, m.ChannelID, s) {
				Already = nil
				Done = nil
				VtuberName := strings.Split(strings.TrimSpace(strings.Replace(m.Content, Prefix+TagRoles, "", -1)), " ")

				guild, err := s.Guild(m.GuildID)
				if err != nil {
					log.Error(err)
				}

				if len(VtuberName[len(VtuberName)-1:]) > 0 {
					tmp := strings.Split(VtuberName[len(VtuberName)-1:][0], ",")

					for _, Name := range tmp {
						Data := FindName(Name)
						if Data.GroupID == 0 {
							VTuberGroup, err := FindGropName(Name)
							if err != nil {
								log.Error(err)
								_, err := s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
								if err != nil {
									log.Error(err)
								}
								return
							}

							if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
								for _, Role := range guild.Roles {
									for _, UserRole := range VtuberName {
										if UserRole == Role.Mention() {
											for _, Member := range database.GetMembers(VTuberGroup.ID) {
												User := database.UserStruct{
													DiscordID:       Role.ID,
													DiscordUserName: Role.Name,
													Channel_ID:      m.ChannelID,
													GroupID:         VTuberGroup.ID,
													Human:           false,
												}
												err := User.Deluser(Member.ID)
												if err != nil {
													Already = append(Already, "`"+Member.Name+"`")
												} else {
													Done = append(Done, "`"+Member.Name+"`")
												}
											}

											if Already != nil {
												_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Already Remove "+strings.Join(Already, " ")+" from tags list or "+Role.Mention()+" never add them \n").
													AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to role tags list").
													InlineAllFields().
													SetColor(Color).MessageEmbed)
												if err != nil {
													log.Error(err)
												}
												Already = nil
											} else if Done != nil {
												_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Remove\n"+strings.Join(Done, " ")+"\n from tag list").
													AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
													InlineAllFields().
													SetColor(Color).MessageEmbed)
												if err != nil {
													log.Error(err)
												}
												Done = nil
											}
										}
									}
								}
							} else {
								_, err := s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.GroupName+"`")
								if err != nil {
									log.Error(err)
								}
								return
							}
						} else {
							MemberTag = append(MemberTag, Data)
						}
					}
					Already = nil
					Done = nil
					for i, Member := range MemberTag {
						if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
							for _, Role := range guild.Roles {
								for _, UserRole := range VtuberName {
									if UserRole == Role.Mention() {
										User := database.UserStruct{
											DiscordID:       Role.ID,
											DiscordUserName: Role.Name,
											Channel_ID:      m.ChannelID,
											GroupID:         Member.GroupID,
											Human:           false,
										}
										err := User.Deluser(Member.MemberID)
										if err != nil {
											Already = append(Already, "`"+tmp[i]+"`")
										} else {
											Done = append(Done, "`"+tmp[i]+"`")
										}

										if Already != nil {
											_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
												SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
												SetDescription(Role.Mention()+" Already Remove "+strings.Join(Already, " ")+" from tags list or "+Role.Mention()+" never add them \n").
												SetThumbnail(config.GoSimpIMG).
												SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
												SetColor(Color).MessageEmbed)
											if err != nil {
												log.Error(err)
											}
											Already = nil
										}
										if Done != nil {
											_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
												SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
												SetDescription(Role.Mention()+" Remove\n"+strings.Join(Done, " ")+"\n from tag list").
												SetThumbnail(config.GoSimpIMG).
												SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show you tags list").
												SetColor(Color).MessageEmbed)
											if err != nil {
												log.Error(err)
											}
											Done = nil
										}
									}
								}
							}
						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+Member.GroupName+"`").
								SetThumbnail(config.GoSimpIMG).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+DelRoles+"` command")
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "You don't have enough permission to use this command,Only user with permission `Manage Channel or Higher` can use this command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.HasPrefix(m.Content, Prefix+RolesReminder) {
			if CheckPermission(m.Author.ID, m.ChannelID, s) {
				Already = nil
				Done = nil
				UserInput := strings.Replace(m.Content, Prefix+RolesReminder, "", -1)
				VtuberName := strings.Split(strings.TrimSpace(UserInput), " ")
				tmpvar := VtuberName[len(VtuberName)-1]

				ReminderUser, err := strconv.Atoi(tmpvar[:len(tmpvar)-1] + "6")
				if err != nil {
					log.Error(err)
					_, err := s.ChannelMessageSend(m.ChannelID, "Invalid number")
					if err != nil {
						log.Error(err)
					}
					return
				} else if ReminderUser > 67 {
					s.ChannelMessageSend(m.ChannelID, "Can't set Reminder over than 60 Minutes\nMin 10 Max 60")
					return
				} else if ReminderUser == 6 {
					ReminderUser = 0
					s.ChannelMessageSend(m.ChannelID, "You disable reminder time")
				}

				guild, err := s.Guild(m.GuildID)
				if err != nil {
					log.Error(err)
				}

				if len(VtuberName[len(VtuberName)-2:]) > 0 {
					tmp := strings.Split(VtuberName[len(VtuberName)-2:][0], ",")
					for _, Name := range tmp {
						Data := FindName(Name)
						if Data.GroupID == 0 {
							VTuberGroup, err := FindGropName(Name)
							if err != nil {
								log.Error(err)
								_, err := s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `"+VtuberData+"` command to see vtubers name or see at web site \n "+config.VtubersData)
								if err != nil {
									log.Error(err)
								}
								return
							}

							if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
								for _, Role := range guild.Roles {
									for _, UserRole := range VtuberName {
										if UserRole == Role.Mention() {
											for _, Member := range database.GetMembers(VTuberGroup.ID) {
												User := database.UserStruct{
													DiscordID:       Role.ID,
													DiscordUserName: Role.Name,
													Channel_ID:      m.ChannelID,
													GroupID:         VTuberGroup.ID,
													Human:           false,
													Reminder:        ReminderUser,
												}
												err := User.UpdateReminder(Member.ID)
												if err != nil {
													log.Error(err)
													_, err := s.ChannelMessageSend(m.ChannelID, Role.Mention()+" not tag `"+VTuberGroup.GroupName+"`")
													if err != nil {
														log.Error(err)
													}
													return
												} else {
													Done = append(Done, "`"+Member.Name+"`")
												}
											}
											if Done != nil {
												if ReminderUser == 0 {
													_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Disable reminder time\n"+strings.Join(Done, " ")).
														AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
														InlineAllFields().
														SetColor(Color).MessageEmbed)
													if err != nil {
														log.Error(err)
													}
												} else {
													_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Change reminder to "+strconv.Itoa(ReminderUser)+"\n"+strings.Join(Done, " ")).
														AddField("Group Name", "**"+VTuberGroup.GroupName+"**").
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show role tags list").
														InlineAllFields().
														SetColor(Color).MessageEmbed)
													if err != nil {
														log.Error(err)
													}

												}
												Done = nil
											}
										}
									}
								}
							} else {
								_, err := s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.GroupName+"`")
								if err != nil {
									log.Error(err)
								}
								return
							}
						} else {
							MemberTag = append(MemberTag, Data)
						}
						Done = nil
					}
					for i, Member := range MemberTag {
						if database.CheckChannelEnable(m.ChannelID, tmp[i], Member.GroupID) {
							for _, Role := range guild.Roles {
								for _, UserRole := range VtuberName {
									if UserRole == Role.Mention() {
										User := database.UserStruct{
											DiscordID:       Role.ID,
											DiscordUserName: Role.Name,
											Channel_ID:      m.ChannelID,
											GroupID:         Member.GroupID,
											Human:           false,
											Reminder:        ReminderUser,
										}
										err := User.UpdateReminder(Member.MemberID)
										if err != nil {
											log.Error(err)
											_, err := s.ChannelMessageSend(m.ChannelID, Role.Mention()+" not tag `"+Member.MemberName+"`")
											if err != nil {
												log.Error(err)
											}
											return
										} else {
											Done = append(Done, "`"+tmp[i]+"`")
										}
									}

									if Done != nil {
										_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
											SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
											SetDescription(Role.Mention()+"Change reminder to "+strconv.Itoa(ReminderUser)+"\n"+strings.Join(Done, " ")).
											SetThumbnail(config.GoSimpIMG).
											SetFooter("Use \""+config.PGeneral+RolesTags+" @"+Role.Name+"\" to show you tags list").
											SetColor(Color).MessageEmbed)
										if err != nil {
											log.Error(err)
										}
										Done = nil
									}
								}
							}

						} else {
							_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+Member.GroupName+"`").
								SetThumbnail(config.GoSimpIMG).
								SetColor(Color).MessageEmbed)
							if err != nil {
								log.Error(err)
							}
							return
						}
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+RolesReminder+"` command")
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "You don't have enough permission to use this command,Only user with permission `Manage Channel or Higher` can use this commandYou don't have enough permission to use this command")
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

//CheckPermission Check user permission
func CheckPermission(User, Channel string, bot *discordgo.Session) bool {
	a, err := bot.UserChannelPermissions(User, Channel)
	if err != nil {
		log.Error(err)
	}
	if a&config.ChannelPermission != 0 {
		return true
	}
	return false
}

//EnableState Enable command message handler
func EnableState(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	if strings.HasPrefix(m.Content, Prefix) {
		Color, err := engine.GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}

		var (
			ChannelState = database.DiscordChannel{
				TypeTag:     0,
				LiveOnly:    false,
				NewUpcoming: true,
				ChannelID:   m.ChannelID,
			}
			already []string
			done    []string
			msg1    = engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetDescription("Remember boys,always respect the author,**do not save the fanart without permission from the author**").
				SetThumbnail(config.GoSimpIMG).
				SetColor(Color)

			msg2 = engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetDescription("Every livestream have some **rule**,follow the **rule** and don't be asshole").
				SetThumbnail(config.GoSimpIMG).
				SetColor(Color)

			msg3 = engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetDescription("Remember boys,always respect the author,**do not save the fanart without permission from the author**\nEvery livestream have some **rule**,follow the **rule** and don't be asshole").
				SetThumbnail(config.GoSimpIMG).
				SetColor(Color)

			msg4 = engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Donate").
				SetDescription("Buy dev coffie to improve bot performance and make dev happy").
				SetURL(config.KoFiLink).
				SetThumbnail(config.GoSimpIMG).
				SetColor(Color)
		)
		CommandArray := strings.Split(m.Content, " ")
		if CommandArray[0] == Prefix+Enable {
			if len(CommandArray) > 2 {
				if CommandArray[1] == "art" {
					ChannelState.SetTypeTag(1)
				} else if CommandArray[1] == "live" {
					ChannelState.SetTypeTag(2)
				} else if CommandArray[1] == "all" {
					ChannelState.SetTypeTag(3)
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[1]+"`,Invalid type")
					if err != nil {
						log.Error(err)
					}
					return
				}

				if (ChannelState.TypeTag == 2 || ChannelState.TypeTag == 3) && len(CommandArray) >= 4 {
					if CommandArray[3] == "-liveonly" {
						ChannelState.SetLiveOnly(true)
					} else if CommandArray[3] == "-newupcoming" {
						ChannelState.SetNewUpcoming(true)
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[3]+"`,Invalid options")
						if err != nil {
							log.Error(err)
						}
						return
					}
				} else if ChannelState.TypeTag == 1 && len(CommandArray) >= 4 {
					_, err := s.ChannelMessageSend(m.ChannelID, "You enable `Art` state,Ignoring `"+CommandArray[3]+"` options")
					if err != nil {
						log.Error(err)
					}
					ChannelState.SetLiveOnly(false)
					ChannelState.SetNewUpcoming(false)
				}

				if (ChannelState.TypeTag == 2 || ChannelState.TypeTag == 3) && len(CommandArray) >= 5 {
					if CommandArray[4] == "-newupcoming" {
						ChannelState.SetNewUpcoming(true)
					} else if CommandArray[4] == "-liveonly" {
						ChannelState.SetLiveOnly(true)
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[4]+"`,Invalid options")
						if err != nil {
							log.Error(err)
						}
						return
					}
				}

				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[2]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						if err != nil {
							log.Error(err)
						}
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						ChannelState.SetVtuberGroupID(VTuberGroup.ID)
						if ChannelState.ChannelCheck() {
							already = append(already, "`"+VTuberGroup.GroupName+"`")
						} else {
							err := ChannelState.AddChannel()
							if err != nil {
								log.Error(err)
								_, err := s.ChannelMessageSend(m.ChannelID, "Something error XD")
								if err != nil {
									log.Error(err)
								}
							}
							done = append(done, "`"+VTuberGroup.GroupName+"`")

						}
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
				if done != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, "done, <@"+m.Author.ID+"> is enable "+strings.Join(done, ",")+" on this channel")
					if err != nil {
						log.Error(err)
					}
					msgID := ""
					if ChannelState.TypeTag == 1 {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg1.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					} else if ChannelState.TypeTag == 2 {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg2.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					} else {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg3.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					}
					MessagePinned, err := s.ChannelMessagesPinned(m.ChannelID)

					for _, Message := range MessagePinned {
						if Message.Author.ID == BotInfo.ID {
							err := s.ChannelMessageUnpin(m.ChannelID, Message.ID)
							if err != nil {
								log.Error(err)
							}
						}
					}
					err = s.ChannelMessagePin(m.ChannelID, msgID)
					if err != nil {
						log.Error(err)
					}
					if rand.Float32() < 0.5 && config.KoFiLink != "" {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg4.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+", already enabled")
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Enable+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if CommandArray[0] == Prefix+Disable {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not valid")
						if err != nil {
							log.Error(err)
						}
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						ChannelState.SetVtuberGroupID(VTuberGroup.ID)
						if ChannelState.ChannelCheck() {
							err := ChannelState.DelChannel("Delete")
							if err != nil {
								log.Error(err)
								_, err := s.ChannelMessageSend(m.ChannelID, "Something error XD")
								if err != nil {
									log.Error(err)
								}
								return
							}
							done = append(done, "`"+VTuberGroup.GroupName+"`")
						} else {
							already = append(already, "`"+VTuberGroup.GroupName+"`")
						}
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
				if done != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, "done, <@"+m.Author.ID+"> is disable "+strings.Join(done, ",")+" from this channel")
					if err != nil {
						log.Error(err)
					}
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+", already removed or never enable on this channel")
					if err != nil {
						log.Error(err)
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Disable+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		} else if CommandArray[0] == Prefix+Update {
			if len(CommandArray) > 2 {
				if CommandArray[1] == "art" {
					ChannelState.SetTypeTag(1)
				} else if CommandArray[1] == "live" {
					ChannelState.SetTypeTag(2)
				} else if CommandArray[1] == "all" {
					ChannelState.SetTypeTag(3)
				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[1]+"`,Invalid type")
					if err != nil {
						log.Error(err)
					}
					return
				}

				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[2]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						if err != nil {
							log.Error(err)
						}
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						ChannelState.SetVtuberGroupID(VTuberGroup.ID)
						if ChannelState.ChannelCheck() {
							if (ChannelState.TypeTag == 2 || ChannelState.TypeTag == 3) && len(CommandArray) >= 4 {
								if CommandArray[3] == "-liveonly" {
									err := ChannelState.SetLiveOnly(true).UpdateChannel("LiveOnly")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[3] == "-newupcoming" {
									err := ChannelState.SetNewUpcoming(true).UpdateChannel("NewUpcoming")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[3] == "-rm_liveonly" {
									err := ChannelState.SetLiveOnly(false).UpdateChannel("LiveOnly")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[3] == "-rm_newupcoming" {
									err := ChannelState.SetNewUpcoming(false).UpdateChannel("NewUpcoming")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else {
									_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[3]+"`,Invalid options")
									if err != nil {
										log.Error(err)
									}
									return
								}
							}

							if (ChannelState.TypeTag == 2 || ChannelState.TypeTag == 3) && len(CommandArray) >= 5 {
								if CommandArray[4] == "-newupcoming" {
									err := ChannelState.SetNewUpcoming(true).UpdateChannel("NewUpcoming")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[4] == "-liveonly" {
									err := ChannelState.SetLiveOnly(true).UpdateChannel("LiveOnly")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[4] == "-rm_liveonly" {
									err := ChannelState.SetLiveOnly(false).UpdateChannel("LiveOnly")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else if CommandArray[4] == "-rm_newupcoming" {
									err := ChannelState.SetNewUpcoming(false).UpdateChannel("NewUpcoming")
									if err != nil {
										already = append(already, "`"+VTuberGroup.GroupName+"`")
									} else {
										done = append(done, "`"+VTuberGroup.GroupName+"`")
									}
								} else {
									_, err := s.ChannelMessageSend(m.ChannelID, "`"+CommandArray[4]+"`,Invalid options")
									if err != nil {
										log.Error(err)
									}
									return
								}
							}

							err := ChannelState.UpdateChannel("Type")
							if err != nil {
								already = append(already, "`"+VTuberGroup.GroupName+"`")
							} else {
								done = append(done, "`"+VTuberGroup.GroupName+"`")
							}
						} else {
							_, err := s.ChannelMessageSend(m.ChannelID, "this channel not enable `"+VTuberGroup.GroupName+"`")
							if err != nil {
								log.Error(err)
							}
							return
						}
					} else {
						_, err := s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						if err != nil {
							log.Error(err)
						}
						return
					}
				}
				if done != nil {
					_, err := s.ChannelMessageSend(m.ChannelID, "done,<@"+m.Author.ID+"> update channel state "+strings.Join(done, ","))
					if err != nil {
						log.Error(err)
					}

					msgID := ""
					if ChannelState.TypeTag == 1 {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg1.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					} else if ChannelState.TypeTag == 2 {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg2.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					} else {
						tmp, err := s.ChannelMessageSendEmbed(m.ChannelID, msg3.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
						msgID = tmp.ID
					}
					MessagePinned, err := s.ChannelMessagesPinned(m.ChannelID)

					for _, Message := range MessagePinned {
						if Message.Author.ID == BotInfo.ID {
							err := s.ChannelMessageUnpin(m.ChannelID, Message.ID)
							if err != nil {
								log.Error(err)
							}
						}
					}
					err = s.ChannelMessagePin(m.ChannelID, msgID)
					if err != nil {
						log.Error(err)
					}
					if rand.Float32() < 0.5 && config.KoFiLink != "" {
						_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg4.MessageEmbed)
						if err != nil {
							log.Error(err)
						}
					}

				} else {
					_, err := s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+" Same state")
					if err != nil {
						log.Error(err)
					}
				}

			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+Update+"` command")
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}

//Help helmp command message handler
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	if strings.HasPrefix(m.Content, Prefix) {
		Color, err := engine.GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}
		if m.Content == Prefix+"help en" || m.Content == Prefix+"help" {
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetURL(config.CommandURL).
				SetDescription("See at web site\n"+config.CommandURL).
				/*
					AddField(Prefix+Enable+" {art/live/all} [Vtuber Group]", "This command will declare if [Vtuber Group] enable in this channel\nExample:\n`"+config.PGeneral+Enable+" all hanayori,hololive` so other users can use `"+config.PGeneral+TagMe+" kanochi` or "+"`"+config.PGeneral+TagMe+" gura`").
					AddField(Prefix+Update+" {art/live/all} [Vtuber Group]", "Use this command if you want to change enable state").
					AddField(Prefix+Disable+" [Vtuber Group]", "Just like enable but this disable command :3 ").
					AddField(config.PFanart+"{Group/Member name}", "Show fanart with randomly with their fanart hashtag\nExample: \n`"+config.PFanart+"Kanochi` or `"+config.PFanart+"hololive`").
					AddField(Prefix+TagMe+" [Group/Member name]", "This command will add you to tags list if any new fanart and livestream schedule\nExample: \n`"+config.PGeneral+TagMe+" Kanochi`,then you will get tagged when there is a new fanart and livestream schedule of kano").
					AddField(Prefix+DelTag+" [Group/Member name]", "This command will remove you from tags list").
					AddField(Prefix+MyTags, "Show all lists that you are subscribed on this channel").
					AddField(Prefix+TagRoles+" [Roles name]", "Same like `tag me` but this will tag roles").
					AddField(Prefix+DelRoles+" [Roles name]", "Remove roles from tags list").
					AddField(Prefix+RolesTags+" [Roles name]", "Show all tags list that roles subscribed on this channel").
					AddField(Prefix+ChannelState, "Show what is enable in this channel").
					AddField(Prefix+VtuberData+" [Group] -Region {region}", "Show available Vtuber data ").
					AddField(Prefix+Subscriber+" {Member name}", "Show Vtuber count of subscriber and followers ").
					AddField(config.PYoutube+Upcoming+" [Vtuber Group/Member] -Region {region}", "This command will show Upcoming live streams on Youtube (*only 3 if use Group name*)").
					AddField(config.PYoutube+Live+" [Vtuber Group/Member] -Region {region}", "This command will show all live streams right now on Youtube").
					AddField(config.PYoutube+Past+" [Vtuber Group/Member] -Region {region}", "This command will show past streams on Youtube (*only 3 if use Group name*)").
					AddField("~~"+config.PBilibili+Upcoming+" [Vtuber Group/Member]~~", "~~This command will show all Upcoming live streams on BiliBili~~").
					AddField(config.PBilibili+Live+" [Vtuber Group/Member]", "This command will show all live streams right now on BiliBili (*only 3 if use Group name*)").
					AddField(config.PBilibili+Past+" [Vtuber Group/Member]", "This command will show all past streams on BiliBili").
					AddField("sp_"+config.PBilibili+"[Vtuber Group/Member]", "This command will show latest video on bilibili").
					AddField(Prefix+"Help EN", "Well,you using it right now").
					AddField(Prefix+"Help JP", "Like this but in Japanese").
				*/
				SetThumbnail(config.BSD).
				//SetFooter("Only user with permission \"Manage Channel or Higher\" can Enable/Disable/Update Vtuber Group").
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.Error(err)
			}
			return
		} else if m.Content == Prefix+"help jp" { //i'm just joking lol
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetDescription("日本語が話せるようになってヘルプメニューを作りたい\n~Dev").
				SetImage("https://i.imgur.com/f0no1r7.png").
				SetFooter("More like,help me").
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.Error(err)
			}
			return
		}
	}
}

//Status command message handler
func Status(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral

	if strings.HasPrefix(m.Content, Prefix) {
		Color, err := engine.GetColor("/tmp/discordpp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(true)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		if strings.HasPrefix(m.Content, Prefix+RolesTags) {
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				log.Error(err)
			}
			RolesInput := strings.Split(strings.TrimSpace(strings.Replace(m.Content, Prefix+RolesTags, "", -1)), " ")
			if len(RolesInput) > 0 {
				for _, UserRoles := range RolesInput {
					for _, Role := range guild.Roles {
						if UserRoles == Role.Mention() {
							list := database.UserStatus(Role.ID, m.ChannelID)
							if list != nil {
								tableString := &strings.Builder{}
								table := tablewriter.NewWriter(tableString)
								table.SetHeader([]string{"Vtuber Group", "Vtuber Name", "Reminder"})
								table.SetAutoWrapText(false)
								table.SetAutoFormatHeaders(true)
								table.SetCenterSeparator("")
								table.SetColumnSeparator("")
								table.SetRowSeparator("")
								table.SetHeaderLine(true)
								table.SetBorder(false)
								table.SetTablePadding("\t")
								table.SetNoWhiteSpace(true)
								table.AppendBulk(list)
								table.Render()

								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetThumbnail(m.Author.AvatarURL("128")).
									SetDescription("Role "+Role.Mention()+"\n```"+tableString.String()+"```").
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}

							} else {
								_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetTitle("404 Not found").
									SetImage(config.NotFound).
									SetColor(Color).MessageEmbed)
								if err != nil {
									log.Error(err)
								}
							}
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSend(m.ChannelID, "Incomplete `"+RolesTags+"` command")
				if err != nil {
					log.Error(err)
				}
			}

		} else if m.Content == Prefix+MyTags {
			list := database.UserStatus(m.Author.ID, m.ChannelID)

			if list != nil {
				table.SetHeader([]string{"Vtuber Group", "Vtuber Name", "Reminder"})
				table.AppendBulk(list)
				table.Render()

				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(m.Author.AvatarURL("128")).
					SetDescription("```"+tableString.String()+"```").
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			}
		} else if m.Content == Prefix+ChannelState {
			list, Type, LiveOnly, NewUpcoming := database.ChannelStatus(m.ChannelID)
			if list != nil {
				var (
					Typestr string
				)
				table.SetHeader([]string{"Group", "Type", "LiveOnly", "NewUpcoming"})
				for i := 0; i < len(list); i++ {
					if Type[i] == 1 {
						Typestr = "Art"
					} else if Type[i] == 2 {
						Typestr = "Live"
					} else {
						Typestr = "All"
					}
					table.Append([]string{list[i], Typestr, LiveOnly[i], NewUpcoming[i]})
				}
				table.Render()

				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetDescription("```"+tableString.String()+"```").
					SetThumbnail(config.GoSimpIMG).
					SetColor(Color).
					SetFooter("Use `update` command to change type of channel").MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
			}
		} else if strings.HasPrefix(m.Content, Prefix+VtuberData) {
			var (
				re        = regexp.MustCompile(`(?m)-region\s.+`)
				tmpvar    = re.FindAllString(m.Content, -1)
				UserInput = strings.Replace(m.Content, Prefix+VtuberData, "", -1)
				RegInput  []string
			)

			if len(tmpvar) > 0 {
				vartmp2 := strings.TrimSpace(strings.Replace(tmpvar[0], "-region", "", -1))
				RegInput = strings.Split(vartmp2, ",")

				UserInput = strings.Replace(UserInput, tmpvar[0], "", -1)
			} else {
				UserInput = strings.Replace(UserInput, "-region", "", -1)
			}
			GroupInput := strings.Split(strings.TrimSpace(UserInput), ",")
			if len(GroupInput) > 0 {
				var (
					GroupsByReg = RegInput
					NiggList    = make(map[string]string)
				)
				if len(RegInput) > 0 {
					for _, Group := range GroupInput {
						var (
							black []string
						)
						for _, Reg := range GroupsByReg {
							Counter := CheckReg(Group, Reg)
							if !Counter {
								black = append(black, Reg)
							}
						}
						if black != nil {
							NiggList[Group] = strings.Join(black, ",")
						}
					}
				}
				for _, Group := range engine.GroupData {
					for _, Grp := range GroupInput {
						if Grp == strings.ToLower(Group.GroupName) {
							for _, Member := range database.GetMembers(Group.ID) {
								yt := ""
								bl := ""
								if Member.YoutubeID != "" {
									yt = "✓"
								} else {
									yt = "✗"
								}

								if Member.BiliBiliID != 0 {
									bl = "✓"
								} else {
									bl = "✗"
								}

								if GroupsByReg != nil {
									table.SetHeader([]string{"Nickname", "Region", "Youtube", "BiliBili", "Group"})
									for _, Reg := range GroupsByReg {
										if Reg == strings.ToLower(Member.Region) {
											table.Append([]string{Member.Name, Member.Region, yt, bl, Group.GroupName})
										}
									}
								} else {
									table.SetHeader([]string{"Nickname", "Region", "Youtube", "BiliBili"})
									table.Append([]string{Member.Name, Member.Region, yt, bl})
								}
							}
						}
					}
				}
				table.Render()
				if len(tableString.String()) > engine.EmbedLimitDescription {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetThumbnail(config.GoSimpIMG).
						SetURL(config.VtubersData).
						SetDescription("Data too longgggggg\nsee Vtubers Data at web site\n"+config.VtubersData).
						SetImage(config.Longcatttt).
						SetColor(Color).MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				} else if len(tableString.String()) > 1500 {
					_, err := s.ChannelMessageSend(m.ChannelID, "```"+tableString.String()+"```")
					if err != nil {
						log.Error(err)
					}
				} else if tableString.String() == "" {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetThumbnail(config.GoSimpIMG).
						SetTitle("List of Vtuber Groups").
						SetURL(config.VtubersData).
						SetDescription("```"+strings.Join(engine.GroupsName, "\t")+"```For more detail see at "+config.VtubersData).
						SetColor(Color).
						SetFooter("Use Name of group to show vtuber members").MessageEmbed)
					if err != nil {
						log.Error(err)
					}

				} else if len(tableString.String()) > 42 {
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("```"+tableString.String()+"```").
						SetColor(Color).
						SetFooter("Use \"Nickname\" as parameter").MessageEmbed)
					if err != nil {
						log.Error(err)
					}
				}

				if NiggList != nil {
					for key, val := range NiggList {
						_, err := s.ChannelMessageSend(m.ChannelID, "`"+strings.Title(key)+"` don't have member in `"+strings.ToUpper(val)+"`")
						if err != nil {
							log.Error(err)
						}
					}
				}
			} else {
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(config.GoSimpIMG).
					SetURL(config.CommandURL).
					SetDescription("Invalid command,see command at my github\n"+config.CommandURL).
					SetColor(Color).MessageEmbed)
				if err != nil {
					log.Error(err)
				}
				return
			}
		}
	}
}
