package discordhandler

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
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
		DynamicData Dynamic_svr
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
				body, errcurl = engine.Curl(urls, nil)
				if errcurl != nil {
					log.Error(errcurl, string(body))
					log.Info("Trying use tor")

					body, errcurl = engine.CoolerCurl(urls, nil)
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
					SetFooter(Msg, config.TwitterIMG).MessageEmbed
			} else {
				embed = engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetTitle(Author).
					SetThumbnail(GetUserAvatar(Author)).
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
			if m.Content == strings.ToLower(Prefix+GroupData.NameGroup) {
				Color, err := engine.GetColor("/tmp/mem.tmp", m.Author.AvatarURL("80"))
				if err != nil {
					log.Error(err)
				}

				DataFix := database.GetFanart(GroupData.ID, 0)
				if DataFix.Videos != "" {
					Msg = "Video type,check original post"
					Pic = config.NotFound
				} else if len(DataFix.Photos) > 0 {
					Pic = DataFix.Photos[0]
					Color, err = engine.GetColor("/tmp/mem.tmp", DataFix.Photos[0])
					if err != nil {
						log.Error(err)
					}
				} else {
					Msg = "Original post was deleted"
					Pic = config.NotFound

					log.WithFields(log.Fields{
						"PermanentURL": DataFix.PermanentURL,
					}).Warn("Original post was deleted")
				}
				Group = SendNude(engine.FixName(DataFix.EnName, DataFix.JpName),
					DataFix.Author, RemovePic(DataFix.Text),
					DataFix.PermanentURL,
					Pic, Msg, Color,
					DataFix.State, DataFix.Dynamic_id)
				break
			}
			for _, MemberData := range database.GetName(GroupData.ID) {
				if m.Content == strings.ToLower(Prefix+MemberData.Name) || m.Content == strings.ToLower(Prefix+MemberData.JpName) {
					Color, err := engine.GetColor("/tmp/mem.tmp", m.Author.AvatarURL("80"))
					if err != nil {
						log.Error(err)
					}

					DataFix := database.GetFanart(0, MemberData.ID)
					if DataFix.Videos != "" {
						Msg = "Video type,check original post"
						Pic = config.NotFound
					} else if len(DataFix.Photos) > 0 {
						Pic = DataFix.Photos[0]
						Color, err = engine.GetColor("/tmp/mem.tmp", DataFix.Photos[0])
						if err != nil {
							log.Error(err)
						}
					} else {
						Msg = "Original post was deleted"

						log.WithFields(log.Fields{
							"PermanentURL": DataFix.PermanentURL,
						}).Warn("Original post was deleted")
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
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `vtuber data` command to see vtubers name or see at my github https://github.com/JustHumanz/Go-Simp")
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID).
								SetReminder(ReminderUser)

							for _, Member := range database.GetName(VTuberGroup.ID) {
								err := User.Adduser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
								}
							}
							if Already != nil || Done != nil {
								if Already != nil {
									s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your tag list").
										AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
										SetImage(VTuberGroup.IconURL).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
										SetColor(Color).MessageEmbed)
								}
								if Done != nil {
									s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Add\n"+strings.Join(Done, " ")+" to your tag list").
										AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
										SetImage(VTuberGroup.IconURL).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
										SetColor(Color).MessageEmbed)
								}
							}
						} else {
							s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.NameGroup+"`")
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
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						return
					}
				}

				if Already != nil || Done != nil {
					if Already != nil {
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your list").
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)

					}
					if Done != nil {
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Add\n"+strings.Join(Done, " ")+" to your list").
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)
					}
				}

			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `"+TagMe+"` command")
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
							s.ChannelMessageSend(m.ChannelID, "Can't set Reminder over than 60 Minutes")
							return
						}

					} else {
						//Invaild
						s.ChannelMessageSend(m.ChannelID, "Invaild number")
						return
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Number not found")
					return
				}

				tmp := strings.Split(FindInt[1], ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data.GroupID == 0 {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `vtuber data` command to see vtubers name or see at my github https://github.com/JustHumanz/Go-Simp")
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID).
								SetReminder(ReminderUser)
							for _, Member := range database.GetName(VTuberGroup.ID) {
								err := User.UpdateReminder(Member.ID)
								if err != nil {
									log.Error(err)
								}
								Done = append(Done, "`"+Member.Name+"`")
							}
							if Done != nil {
								s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You Update reminder time\n"+strings.Join(Done, " ")+" to your list").
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									SetColor(Color).MessageEmbed)
							}
							Done = nil
						} else {
							s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.NameGroup+"`")
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
						}
						Done = append(Done, "`"+tmp[i]+"`")
					} else {
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						return
					}
				}
				if Done != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You Update reminder time\n"+strings.Join(Done, " ")+" to your list").
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
						SetColor(Color).MessageEmbed)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `"+SetReminder+"` command")
				return
			}
		} else if strings.HasPrefix(m.Content, Prefix+DelTag) {
			Already = nil
			Done = nil
			VtuberName := strings.TrimSpace(strings.Replace(m.Content, Prefix+DelTag, "", -1))
			if (VtuberName) != "" {
				tmp := strings.Split(VtuberName, ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data == (NameStruct{}) {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `vtuber data` command to see vtubers name or see at my github https://github.com/JustHumanz/Go-Simp")
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.SetGroupID(VTuberGroup.ID)
							for _, Member := range database.GetName(VTuberGroup.ID) {
								err := User.Deluser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
								}
							}
							if Already != nil {
								s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
									AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
									SetImage(VTuberGroup.IconURL).
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									SetColor(Color).MessageEmbed)
							} else if Done != nil {
								s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
									AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
									SetThumbnail(config.GoSimpIMG).
									SetImage(VTuberGroup.IconURL).
									SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
									SetColor(Color).MessageEmbed)
							}
						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+VTuberGroup.NameGroup+"`").
								SetImage(VTuberGroup.IconURL).
								SetThumbnail(config.GoSimpIMG).
								SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
								SetColor(Color).MessageEmbed)
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
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						return
					}
				}

				if Already != nil || Done != nil {
					if Already != nil {
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)
					}

					if Done != nil {
						//return
						s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+MyTags+"\" to show you tags list").
							SetColor(Color).MessageEmbed)
					}
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `"+DelTag+"` command")
			}
		} else if strings.HasPrefix(m.Content, Prefix+TagRoles) {
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
								s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `vtuber data` command to see vtubers name or see at my github https://github.com/JustHumanz/Go-Simp")
								return
							}

							if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
								for _, Role := range guild.Roles {
									for _, UserRole := range VtuberName {
										if UserRole == Role.Mention() {
											for _, Member := range database.GetName(VTuberGroup.ID) {
												User := database.UserStruct{
													DiscordID:       Role.ID,
													DiscordUserName: Role.Name,
													Channel_ID:      m.ChannelID,
													GroupID:         VTuberGroup.ID,
													Human:           false,
												}
												err := User.Adduser(Member.ID)
												if err != nil {
													Already = append(Already, "`"+Member.Name+"`")
												} else {
													Done = append(Done, "`"+Member.Name+"`")

												}
											}
											if Already != nil || Done != nil {
												if Already != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
														AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to role tags list").
														SetColor(Color).MessageEmbed)
													Already = nil
												}
												if Done != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+"Add\n"+strings.Join(Done, " ")).
														AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to show role tags list").
														SetColor(Color).MessageEmbed)
													Done = nil
												}
											}
										}
									}
								}
							} else {
								s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.NameGroup+"`")
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
										}
										err := User.Adduser(Member.MemberID)
										if err != nil {
											Already = append(Already, "`"+tmp[i]+"`")
										} else {
											Done = append(Done, "`"+tmp[i]+"`")
										}

										if Already != nil || Done != nil {
											if Already != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+"\" to show role tags list").
													SetColor(Color).MessageEmbed)
												Already = nil
											}
											if Done != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Add\n"+strings.Join(Done, " ")).
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+"\" to show you tags list").
													SetColor(Color).MessageEmbed)
												Done = nil
											}
										}
									}
								}
							}

						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+Member.GroupName+"`").
								SetThumbnail(config.GoSimpIMG).
								SetColor(Color).MessageEmbed)
							return
						}
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Incomplete `tag role` command")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You don't have enough permission to use this command")
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
								s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid,use `vtuber data` command to see vtubers name or see at my github https://github.com/JustHumanz/Go-Simp")
								return
							}

							if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
								for _, Role := range guild.Roles {
									for _, UserRole := range VtuberName {
										if UserRole == Role.Mention() {
											for _, Member := range database.GetName(VTuberGroup.ID) {
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
											if Already != nil || Done != nil {
												if Already != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Already Removed from tags list or "+Role.Mention()+" never add them \n"+strings.Join(Already, " ")).
														AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to role tags list").
														SetColor(Color).MessageEmbed)
													Already = nil
												}
												if Done != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Remove\n"+strings.Join(Done, " ")+"\n from tag list").
														AddField("Group Name", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to show role tags list").
														SetColor(Color).MessageEmbed)
													Done = nil
												}
											}
										}
									}
								}
							} else {
								s.ChannelMessageSend(m.ChannelID, "look like this channel not enable `"+VTuberGroup.NameGroup+"`")
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
										}
										err := User.Deluser(Member.MemberID)
										if err != nil {
											Already = append(Already, "`"+tmp[i]+"`")
										} else {
											Done = append(Done, "`"+tmp[i]+"`")

										}

										if Already != nil || Done != nil {
											if Already != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Already Removed from tags list or "+Role.Mention()+" never add them \n"+strings.Join(Already, " ")).
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to show role tags list").
													SetColor(Color).MessageEmbed)
												Already = nil
											}
											if Done != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Remove\n"+strings.Join(Done, " ")+"\n from tag list").
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+RolesTags+Role.Mention()+"\" to show you tags list").
													SetColor(Color).MessageEmbed)
												Done = nil
											}
										}
									}
								}
							}
						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+Member.GroupName+"`").
								SetThumbnail(config.GoSimpIMG).
								SetColor(Color).MessageEmbed)
							return
						}
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Incomplete `tag role` command")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You don't have enough permission to use this command,Only user with permission `Manage Channel or Higher` can use this command ")
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
		var (
			tagtype int
			already []string
			done    []string
			msg1    = "Remember boys,always respect the author,**do not save the fanart without permission from the author**"
			msg2    = "Every livestream have some **rule**,follow the **rule** and don't be asshole"
		)
		CommandArray := strings.Split(m.Content, " ")
		if CommandArray[0] == Prefix+Enable {
			if len(CommandArray) > 1 {
				if CommandArray[1] == "art" {
					tagtype = 1
				} else if CommandArray[1] == "live" {
					tagtype = 2
				} else {
					tagtype = 3
				}
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[len(CommandArray)-1]), ",")

				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							already = append(already, "`"+VTuberGroup.NameGroup+"`")
						} else {
							err := database.AddChannel(m.ChannelID, tagtype, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
							}
							done = append(done, "`"+VTuberGroup.NameGroup+"`")

						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if done != nil {
					s.ChannelMessageSend(m.ChannelID, "done,@here <@"+m.Author.ID+"> is enable "+strings.Join(done, ",")+" on this channel")
					if tagtype == 1 {
						s.ChannelMessageSend(m.ChannelID, msg1+"\n@here")
					} else if tagtype == 2 {
						s.ChannelMessageSend(m.ChannelID, msg2+"\n@here")
					} else {
						s.ChannelMessageSend(m.ChannelID, msg1+"\n"+msg2+"\n@here")
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+", already added")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete enable command")
			}
		} else if CommandArray[0] == Prefix+Disable {
			if len(CommandArray) > 1 {
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[1]), ",")
				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was not valid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							err := database.DelChannel(m.ChannelID, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
								return
							}
							done = append(done, "`"+VTuberGroup.NameGroup+"`")

						} else {
							already = append(already, "`"+VTuberGroup.NameGroup+"`")
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if done != nil {
					s.ChannelMessageSend(m.ChannelID, "done,@here <@"+m.Author.ID+"> is disable "+strings.Join(done, ",")+" from this channel")
				} else {
					s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+", already removed or never enable on this channel")
				}
			}
		} else if CommandArray[0] == Prefix+Update {
			if len(CommandArray) > 1 {
				if CommandArray[1] == "art" {
					tagtype = 1
				} else if CommandArray[1] == "live" {
					tagtype = 2
				} else {
					tagtype = 3
				}
				FindGroupArry := strings.Split(strings.TrimSpace(CommandArray[len(CommandArray)-1]), ",")

				for i := 0; i < len(FindGroupArry); i++ {
					VTuberGroup, err := FindGropName(FindGroupArry[i])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "`"+FindGroupArry[i]+"`,Name of Vtuber Group was invalid")
						return
					}
					if CheckPermission(m.Author.ID, m.ChannelID, s) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							err := database.UpdateChannel(m.ChannelID, tagtype, VTuberGroup.ID)
							if err != nil {
								already = append(already, "`"+VTuberGroup.NameGroup+"`")
							} else {
								done = append(done, "`"+VTuberGroup.NameGroup+"`")

							}
						} else {
							s.ChannelMessageSend(m.ChannelID, "this channel not enable `"+VTuberGroup.NameGroup+"`")
							return
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if done != nil {
					s.ChannelMessageSend(m.ChannelID, "done,<@"+m.Author.ID+"> update channel state "+strings.Join(done, ","))
					if tagtype == 1 {
						s.ChannelMessageSend(m.ChannelID, msg1+"\n@here")
					} else if tagtype == 2 {
						s.ChannelMessageSend(m.ChannelID, msg2+"\n@here")
					} else {
						s.ChannelMessageSend(m.ChannelID, msg1+"\n"+msg2+"\n@here")
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, strings.Join(already, ",")+" Same type")
				}

			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `update` command")
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
			s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetURL(config.CommandURL).
				SetDescription("See at my github repository\n"+config.CommandURL).
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
			return
		} else if m.Content == Prefix+"help jp" { //i'm just joking lol
			s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetDescription("日本語が話せるようになってヘルプメニューを作りたい\n~Dev").
				SetImage("https://i.imgur.com/f0no1r7.png").
				SetFooter("More like,help me").
				SetColor(Color).MessageEmbed)
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

			for _, UserRoles := range strings.Split(strings.TrimSpace(strings.Replace(m.Content, Prefix+RolesTags, "", -1)), " ") {
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

							s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetThumbnail(m.Author.AvatarURL("128")).
								SetDescription("Role "+Role.Mention()+"\n```"+tableString.String()+"```").
								SetColor(Color).MessageEmbed)

						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
								SetTitle("404 Not found").
								SetImage(config.NotFound).
								SetColor(Color).MessageEmbed)
						}
					}
				}
			}

		} else if m.Content == Prefix+MyTags {
			list := database.UserStatus(m.Author.ID, m.ChannelID)

			if list != nil {
				table.SetHeader([]string{"Vtuber Group", "Vtuber Name", "Reminder"})
				table.AppendBulk(list)
				table.Render()

				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(m.Author.AvatarURL("128")).
					SetDescription("```"+tableString.String()+"```").
					SetColor(Color).MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
			}
		} else if m.Content == Prefix+ChannelState {
			list, Type := database.ChannelStatus(m.ChannelID)
			if list != nil {
				var (
					Typestr string
				)
				for i := 0; i < len(list); i++ {
					if Type[i] == 1 {
						Typestr = "Art"
					} else if Type[i] == 2 {
						Typestr = "Live"
					} else {
						Typestr = "All"
					}
					table.Append([]string{list[i], Typestr})
				}
				table.SetHeader([]string{"Group", "Type"})
				table.Render()

				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
					SetDescription("```"+tableString.String()+"```").
					SetThumbnail(config.GoSimpIMG).
					SetColor(Color).
					SetFooter("Use `update` command to change type of channel").MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
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
						if Grp == strings.ToLower(Group.NameGroup) {
							for _, Member := range database.GetName(Group.ID) {
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
											table.Append([]string{Member.Name, Member.Region, yt, bl, Group.NameGroup})
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
					s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetThumbnail(config.GoSimpIMG).
						SetURL(config.VtubersData).
						SetDescription("Data too longgggggg\nsee Vtubers Data at my github\n"+config.VtubersData).
						SetImage(config.Longcatttt).
						SetColor(Color).MessageEmbed)
				} else if len(tableString.String()) > 1500 {
					s.ChannelMessageSend(m.ChannelID, "```"+tableString.String()+"```")
				} else if tableString.String() == "" {
					s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetThumbnail(config.GoSimpIMG).
						SetTitle("List of Vtuber Groups").
						SetURL(config.VtubersData).
						SetDescription("```"+strings.Join(engine.GroupsName, "\t")+"```For more detail see at "+config.VtubersData).
						SetColor(Color).
						SetFooter("Use Name of group to show vtuber members").MessageEmbed)

				} else {
					s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("```"+tableString.String()+"```").
						SetColor(Color).
						SetFooter("Use \"Nickname\" as parameter").MessageEmbed)
				}

				if NiggList != nil {
					for key, val := range NiggList {
						s.ChannelMessageSend(m.ChannelID, "`"+strings.Title(key)+"` don't have member in `"+strings.ToUpper(val)+"`")
					}
				}
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(config.GoSimpIMG).
					SetURL(config.CommandURL).
					SetDescription("Invalid command,see command at my github\n"+config.CommandURL).
					SetColor(Color).MessageEmbed)
				return
			}
		}
	}
}

//Find a valid Vtuber name from message handler
func FindName(MemberName string) NameStruct {
	for _, Group := range engine.GroupData {
		for _, Name := range database.GetName(Group.ID) {
			if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
				return NameStruct{
					GroupName: Group.NameGroup,
					GroupID:   Group.ID,
					MemberID:  Name.ID,
				}
			}
		}
	}
	return NameStruct{}

}

type NameStruct struct {
	GroupName string
	GroupID   int64
	MemberID  int64
}

//Find a valid Vtuber Group from message handler
func FindGropName(GroupName string) (database.GroupName, error) {
	for _, Group := range engine.GroupData {
		if strings.ToLower(Group.NameGroup) == strings.ToLower(GroupName) {
			return Group, nil
		}
	}
	return database.GroupName{}, errors.New(GroupName + " Name Vtuber not valid")
}

//Remove twitter pic
func RemovePic(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)pic\.twitter.com\/.+`).ReplaceAllString(text, "${1}$2")
}

//Get twitter avatar
func GetUserAvatar(username string) string {
	var (
		bit     []byte
		curlerr error
		avatar  string
		url     = "https://mobile.twitter.com/" + regexp.MustCompile("[[:^ascii:]]").ReplaceAllLiteralString(username, "")
	)
	bit, curlerr = engine.Curl(url, nil)
	if curlerr != nil {
		bit, curlerr = engine.CoolerCurl(url, nil)
		if curlerr != nil {
			log.Error(curlerr)
		}
	}

	re := regexp.MustCompile(`(?ms)avatar.*?(http.*?)"`)
	if len(re.FindStringIndex(string(bit))) > 0 {
		re2 := regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
		submatchall := re2.FindAllStringSubmatch(re.FindString(string(bit)), -1)
		for _, element := range submatchall {
			avatar = strings.Replace(element[1], "normal.jpg", "400x400.jpg", -1)
		}
	}
	return avatar
}

//Get bilibili user avatar
func (Data Dynamic_svr) GetUserAvatar() string {
	return Data.Data.Card.Desc.UserProfile.Info.Face
}
