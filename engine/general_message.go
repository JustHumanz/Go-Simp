package engine

import (
	"encoding/json"
	"errors"
	"math/rand"
	"regexp"
	"strings"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
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

	Color, err := GetColor("/tmp/mem.tmp", m.Author.AvatarURL("80"))

	if strings.HasPrefix(m.Content, Prefix) {
		SendNude := func(Title, Author, Text, URL, Pic, Msg string, Color int, State, Dynamic string) bool {
			Msg = Msg + " *sometimes image not showing,because image oversize*"
			if State == "TBiliBili" {
				var (
					body    []byte
					errcurl error
					urls    = "https://api.vc.bilibili.com/dynamic_svr/v1/dynamic_svr/get_dynamic_detail?dynamic_id=" + Dynamic
				)
				body, errcurl = Curl(urls, nil)
				if errcurl != nil {
					log.Error(errcurl, string(body))
					log.Info("Trying use tor")

					body, errcurl = CoolerCurl(urls, nil)
					if errcurl != nil {
						log.Error(errcurl)
					}
				}
				json.Unmarshal(body, &DynamicData)
				embed = NewEmbed().
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
				embed = NewEmbed().
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
			err = Reacting(map[string]string{
				"ChannelID": m.ChannelID,
				"Content":   m.Content,
				"Prefix":    Prefix,
			})
			if err != nil {
				log.Error(err)
			}
			return true
		}
		for _, GroupData := range GroupData {
			if m.Content == strings.ToLower(Prefix+GroupData.NameGroup) {
				DataFix := database.GetFanart(GroupData.ID, 0)
				if DataFix.Videos != "" {
					Msg = "Video type,check original post"
					Pic = config.NotFound
				} else if len(DataFix.Photos) > 0 {
					Pic = DataFix.Photos[0]
					Color, err = GetColor("/tmp/mem.tmp", DataFix.Photos[0])
					if err != nil {
						log.Error(err)
					}
				} else {
					Msg = "Video type,check original post"
					Pic = config.NotFound
				}
				Group = SendNude(FixName(DataFix.EnName, DataFix.JpName),
					DataFix.Author, RemovePic(DataFix.Text),
					DataFix.PermanentURL,
					Pic, Msg, Color,
					DataFix.State, DataFix.Dynamic_id)
				break
			}
			for _, MemberData := range database.GetName(GroupData.ID) {
				if m.Content == strings.ToLower(Prefix+MemberData.Name) || m.Content == strings.ToLower(Prefix+MemberData.JpName) {
					DataFix := database.GetFanart(0, MemberData.ID)
					if DataFix.Videos != "" {
						Msg = "Video type,check original post"
						Pic = config.NotFound
					} else if len(DataFix.Photos) > 0 {
						Pic = DataFix.Photos[0]
						Color, err = GetColor("/tmp/mem.tmp", DataFix.Photos[0])
						if err != nil {
							log.Error(err)
						}
					} else {
						Msg = "Video type,check original post"
					}
					Member = SendNude(FixName(MemberData.EnName, MemberData.JpName),
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
			counter   bool
			Already   []string
			Done      []string
			MemberTag []NameStruct
		)
		User := database.UserStruct{
			DiscordID:       m.Author.ID,
			DiscordUserName: m.Author.Username,
			Channel_ID:      m.ChannelID,
		}
		Color, err := GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}
		if strings.HasPrefix(m.Content, Prefix+TagMe) {
			Already = nil
			Done = nil
			VtuberName := strings.TrimSpace(strings.Replace(m.Content, Prefix+TagMe, "", -1))
			if VtuberName != "" {
				tmp := strings.Split(VtuberName, ",")
				for _, Name := range tmp {
					Data := FindName(Name)
					if Data.GroupID == 0 {
						VTuberGroup, err := FindGropName(Name)
						if err != nil {
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid")
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.GroupID = VTuberGroup.ID
							for _, Member := range database.GetName(VTuberGroup.ID) {
								err := User.Adduser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
									counter = true
								}
							}
							if Already != nil || Done != nil {
								if Already != nil {
									s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your tag list").
										AddField("GroupName", "**"+VTuberGroup.NameGroup+"**").
										SetImage(VTuberGroup.IconURL).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
										SetColor(Color).MessageEmbed)
								}
								if Done != nil {
									s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
										SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
										SetDescription("You Add\n"+strings.Join(Done, " ")+" to your tag list").
										AddField("GroupName", "**"+VTuberGroup.NameGroup+"**").
										SetImage(VTuberGroup.IconURL).
										SetThumbnail(config.GoSimpIMG).
										SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
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
						User.GroupID = Member.GroupID
						err := User.Adduser(Member.MemberID)
						if err != nil {
							Already = append(Already, "`"+tmp[i]+"`")
						} else {
							Done = append(Done, "`"+tmp[i]+"`")
							counter = true
						}
					} else {
						s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						return
					}
				}

				if Already != nil || Done != nil {
					if Already != nil {
						s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Already Added\n"+strings.Join(Already, " ")+" from your list").
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
							SetColor(Color).MessageEmbed)

					}
					if Done != nil {
						s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("You Add\n"+strings.Join(Done, " ")+" to your list").
							SetThumbnail(config.GoSimpIMG).
							SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
							SetColor(Color).MessageEmbed)
					}
				}

			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete `tag me` command")
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
							s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid")
							return
						}
						if database.CheckChannelEnable(m.ChannelID, Name, VTuberGroup.ID) {
							User.GroupID = VTuberGroup.ID
							for _, Member := range database.GetName(VTuberGroup.ID) {
								err := User.Deluser(Member.ID)
								if err != nil {
									Already = append(Already, "`"+Member.Name+"`")
								} else {
									Done = append(Done, "`"+Member.Name+"`")
									counter = true
								}
							}
							if Already != nil {
								s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
									AddField("GroupName", "**"+VTuberGroup.NameGroup+"**").
									SetImage(VTuberGroup.IconURL).
									SetThumbnail(config.GoSimpIMG).
									SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
									SetColor(Color).MessageEmbed)
								return
							} else if Done != nil {
								s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
									SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
									SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
									SetThumbnail(config.GoSimpIMG).
									SetImage(VTuberGroup.IconURL).
									SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
									SetColor(Color).MessageEmbed)
								return
							}
						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetDescription("look like this channel not enable `"+VTuberGroup.NameGroup+"`").
								SetImage(VTuberGroup.IconURL).
								SetThumbnail(config.GoSimpIMG).
								SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
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
							counter = true
						}
					} else {
						s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
							SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
							SetDescription("look like this channel not enable `"+Member.GroupName+"`").
							SetThumbnail(config.GoSimpIMG).
							SetColor(Color).MessageEmbed)
						return
					}
				}
				if Already != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("Already Removed from your tags or You never tag them\n"+strings.Join(Already, " ")).
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
						SetColor(Color).MessageEmbed)
					return
				}

				if counter {
					//return
					s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("You remove "+strings.Join(Done, " ")+" from your tag list").
						SetThumbnail(config.GoSimpIMG).
						SetFooter("Use \""+config.PGeneral+"my tags\" to show you tags list").
						SetColor(Color).MessageEmbed)
					return
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Incomplete del tag command")
			}
		} else if strings.HasPrefix(m.Content, Prefix+TagRoles) {
			if CheckPermission(m.Author.ID, m.ChannelID) {
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
								s.ChannelMessageSend(m.ChannelID, "`"+Name+"` was invalid")
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
												}
												err := User.Adduser(Member.ID)
												if err != nil {
													Already = append(Already, "`"+Member.Name+"`")
												} else {
													Done = append(Done, "`"+Member.Name+"`")
													counter = true
												}
											}
											if Already != nil || Done != nil {
												if Already != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
														AddField("GroupName", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+"role tags\" to role tags list").
														SetColor(Color).MessageEmbed)
													Already = nil
												}
												if Done != nil {
													s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
														SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
														SetDescription(Role.Mention()+"Add\n"+strings.Join(Done, " ")).
														AddField("GroupName", "**"+VTuberGroup.NameGroup+"**").
														SetImage(VTuberGroup.IconURL).
														SetThumbnail(config.GoSimpIMG).
														SetFooter("Use \""+config.PGeneral+"role tags\" to show role tags list").
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
											counter = true
										}

										if Already != nil || Done != nil {
											if Already != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Already Added\n"+strings.Join(Already, " ")).
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+"role tags\" to show role tags list").
													SetColor(Color).MessageEmbed)
												Already = nil
											}
											if Done != nil {
												s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
													SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
													SetDescription(Role.Mention()+" Add\n"+strings.Join(Done, " ")).
													SetThumbnail(config.GoSimpIMG).
													SetFooter("Use \""+config.PGeneral+"role tags\" to show you tags list").
													SetColor(Color).MessageEmbed)
												Done = nil
											}
										}
									}
								}
							}

						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
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
		}
	}
}

//Check user permission
func CheckPermission(User, Channel string) bool {
	a, err := BotSession.UserChannelPermissions(User, Channel)
	BruhMoment(err, "", false)
	if a&config.ChannelPermission != 0 {
		return true
	} else {
		return false
	}
}

//Enable command message handler
func EnableState(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	if strings.HasPrefix(m.Content, Prefix) {
		var (
			counter bool
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
					if CheckPermission(m.Author.ID, m.ChannelID) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							already = append(already, "`"+VTuberGroup.NameGroup+"`")
							counter = false
						} else {
							err := database.AddChannel(m.ChannelID, tagtype, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
							}
							done = append(done, "`"+VTuberGroup.NameGroup+"`")
							counter = true
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if counter {
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
					if CheckPermission(m.Author.ID, m.ChannelID) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							err := database.DelChannel(m.ChannelID, VTuberGroup.ID)
							if err != nil {
								log.Error(err)
								s.ChannelMessageSend(m.ChannelID, "Something error XD")
								return
							}
							done = append(done, "`"+VTuberGroup.NameGroup+"`")
							counter = true
						} else {
							already = append(already, "`"+VTuberGroup.NameGroup+"`")
							counter = false
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "You don't have permission to enable/disable/update")
						return
					}
				}
				if counter {
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
					if CheckPermission(m.Author.ID, m.ChannelID) {
						if database.ChannelCheck(VTuberGroup.ID, m.ChannelID) {
							err := database.UpdateChannel(m.ChannelID, tagtype, VTuberGroup.ID)
							if err != nil {
								already = append(already, "`"+VTuberGroup.NameGroup+"`")
								counter = false
							} else {
								done = append(done, "`"+VTuberGroup.NameGroup+"`")
								counter = true
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
				if counter {
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

//helmp command message handler
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := config.PGeneral
	Color, err := GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
	if err != nil {
		log.Error(err)
	}
	if strings.HasPrefix(m.Content, Prefix) {
		if m.Content == Prefix+"help en" || m.Content == Prefix+"help" {
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetDescription("`[]` => array,support multiple VtuberName/VtuberGroup *separated by commas*\n`{}` => only support single VtuberName/VtuberGroup").
				AddField(Prefix+Enable+" {art/live/all} [Vtuber Group]", "This command will declare if [Vtuber Group] enable in this channel\nExample:\n`"+config.PGeneral+Enable+" all hanayori,hololive` so other users can use `"+config.PGeneral+TagMe+" kanochi` or "+"`"+config.PGeneral+TagMe+" gura`").
				AddField(Prefix+Update+" {art/live/all} [Vtuber Group]", "Use this command if you want to change enable state").
				AddField(Prefix+Disable+" [Vtuber Group]", "Just like enable but this disable command :3 ").
				AddField(config.PFanart+"{Group/Member name}", "Show fanart with randomly with their fanart hashtag\nExample: \n`"+config.PFanart+"Kanochi` or `"+config.PFanart+"hololive`").
				AddField(Prefix+TagMe+" [Group/Member name]", "This command will add you to tags list if any new fanart\nExample: \n`"+config.PGeneral+TagMe+" Kanochi`,then you will get tagged when there is a new fanart and livestream schedule of kano").
				AddField(Prefix+DelTag+" [Group/Member name]", "This command will remove you from tags list").
				AddField(Prefix+MyTags, "Show all lists that you are subscribed").
				AddField(Prefix+ChannelState, "Show what is enable in this channel").
				AddField(Prefix+VtuberData+" [Group] [Region]", "Show available Vtuber data ").
				AddField(Prefix+Subscriber+" {Member name}", "Show Vtuber count of subscriber and followers ").
				AddField(config.PYoutube+Upcoming+" [Vtuber Group/Member] {Region}", "This command will show Upcoming live streams on Youtube  *only 3 if use Vtuber Group*").
				AddField(config.PYoutube+Live+" [Vtuber Group/Member] {Region}", "This command will show all live streams right now on Youtube").
				AddField(config.PYoutube+Past+" [Vtuber Group/Member] {Region}", "This command will show past streams on Youtube *only 3 if use Vtuber Group*").
				AddField("~~"+config.PBilibili+Upcoming+" [Vtuber Group/Member]~~", "~~This command will show all Upcoming live streams on BiliBili~~").
				AddField(config.PBilibili+Live+" [Vtuber Group/Member]", "This command will show all live streams right now on BiliBili").
				AddField(config.PBilibili+Past+" [Vtuber Group/Member]", "This command will show all past streams on BiliBili").
				AddField("sp_"+config.PBilibili+"[Vtuber Group/Member]", "This command will show latest video on bilibili  *only 3 if use Vtuber Group*").
				AddField(Prefix+"Help EN", "Well,you using it right now").
				AddField(Prefix+"Help JP", "Like this but in Japanese").
				SetThumbnail(config.BSD).
				SetFooter("Only user with permission \"Manage Channel or higher\" can Enable/Disable/Update Vtuber Group").
				SetColor(Color).MessageEmbed)
			return
		} else if m.Content == Prefix+"help jp" { //i'm just joking lol
			s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
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
	Color, err := GetColor("/tmp/discordpp", m.Author.AvatarURL("128"))
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

	if strings.HasPrefix(m.Content, Prefix) {
		if strings.HasPrefix(m.Content, Prefix+RoleTags) {
			guild, err := s.Guild(m.GuildID)
			if err != nil {
				log.Error(err)
			}

			for _, UserRoles := range strings.Split(strings.TrimSpace(strings.Replace(m.Content, Prefix+RoleTags, "", -1)), " ") {
				for _, Role := range guild.Roles {
					if UserRoles == Role.Mention() {
						list := database.UserStatus(Role.ID, m.ChannelID)
						if list != nil {
							tableString := &strings.Builder{}
							table := tablewriter.NewWriter(tableString)
							table.SetHeader([]string{"Vtuber Group", "Vtuber Name"})
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

							s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
								SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
								SetThumbnail(m.Author.AvatarURL("128")).
								SetDescription("Role "+Role.Mention()+"\n```"+tableString.String()+"```").
								SetColor(Color).MessageEmbed)

						} else {
							s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
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
				table.SetHeader([]string{"Vtuber Group", "Vtuber Name"})
				table.AppendBulk(list)
				table.Render()

				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(m.Author.AvatarURL("128")).
					SetDescription("```"+tableString.String()+"```").
					SetColor(Color).MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
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

				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("80")).
					SetDescription("```"+tableString.String()+"```").
					SetThumbnail(config.GoSimpIMG).
					SetColor(Color).
					SetFooter("Use `update` command to change type of channel").MessageEmbed)
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetTitle("404 Not found").
					SetImage(config.NotFound).
					SetColor(Color).MessageEmbed)
			}
		} else if strings.HasPrefix(m.Content, Prefix+VtuberData) {
			Parameter := strings.Split(m.Content, " ")
			if len(Parameter) > 2 {
				Groups := strings.Split(Parameter[2], ",")
				var (
					GroupsByReg []string
					NiggList    = make(map[string]string)
				)
				if len(Parameter) > 3 {
					GroupsByReg = strings.Split(Parameter[3], ",")
					for _, Group := range Groups {
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
				for _, Group := range GroupData {
					for _, Grp := range Groups {
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
				if len(tableString.String()) > EmbedLimitDescription {
					s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetThumbnail(config.GoSimpIMG).
						SetDescription("Data too longgggggg").
						SetImage(config.Longcatttt).
						SetColor(Color).MessageEmbed)
				} else {
					s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
						SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
						SetDescription("```"+tableString.String()+"```").
						SetColor(Color).
						SetFooter("Use `Nickname` as parameter").MessageEmbed)
				}

				if NiggList != nil {
					for key, val := range NiggList {
						s.ChannelMessageSend(m.ChannelID, "`"+strings.Title(key)+"` don't have member in `"+strings.ToUpper(val)+"`")
					}
				}
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, NewEmbed().
					SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
					SetThumbnail(config.GoSimpIMG).
					SetDescription("List of Vtuber Groups\n```"+strings.Join(GroupsName, "\t")+"```").
					SetColor(Color).
					SetFooter("Use Name of group to show vtuber members").MessageEmbed)
			}
		}
	}
}

//Find a valid Vtuber name from message handler
func FindName(MemberName string) NameStruct {
	for i := 0; i < len(GroupData); i++ {
		Names := database.GetName(GroupData[i].ID)
		for _, Name := range Names {
			if strings.ToLower(Name.Name) == MemberName || strings.ToLower(Name.JpName) == MemberName {
				return NameStruct{
					GroupName: GroupData[i].NameGroup,
					GroupID:   GroupData[i].ID,
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
	for i := 0; i < len(GroupData); i++ {
		if strings.ToLower(GroupData[i].NameGroup) == strings.ToLower(GroupName) {
			return GroupData[i], nil
		}
	}
	return database.GroupName{}, errors.New(GroupName + " Name Vtuber not valid")
}

//Remove twitter pic
func RemovePic(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)pic\.twitter.com\/.+`).ReplaceAllString(text, "${1}$2")
}

func Reacting(Data map[string]string) error {
	EmojiList := config.EmojiFanart
	ChannelID := Data["ChannelID"]
	MessID, err := BotSession.Channel(ChannelID)
	if err != nil {
		return errors.New(err.Error() + " ChannelID: " + ChannelID)
	}
	for l := 0; l < len(EmojiList); l++ {
		if Data["Content"][len(Data["Prefix"]):] == "kanochi" {
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[0])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}
			break
		} else if Data["Content"][len(Data["Prefix"]):] == "cleaire" {
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}
			if l == len(EmojiList)-1 {
				err = BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":latom:767810745860751391")
				if err != nil {
					return errors.New(err.Error() + " ChannelID: " + ChannelID)
					//log.Error(err, ChannelID)
				}
			}
		} else if Data["Content"][len(Data["Prefix"]):] == "senchou" {
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}

			if l == len(EmojiList)-1 {
				err = BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":hormny:768700671750176790")
				if err != nil {
					return errors.New(err.Error() + " ChannelID: " + ChannelID)
					//log.Error(err, ChannelID)
				}
			}
		} else {
			err := BotSession.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
				//break
			}
		}
	}
	return nil
}

//Get twitter avatar
func GetUserAvatar(username string) string {
	var (
		bit     []byte
		curlerr error
		avatar  string
		url     = "https://mobile.twitter.com/" + regexp.MustCompile("[[:^ascii:]]").ReplaceAllLiteralString(username, "")
	)
	bit, curlerr = Curl(url, nil)
	if curlerr != nil {
		bit, curlerr = CoolerCurl(url, nil)
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

//Guild join handler
func GuildJoin(s *discordgo.Session, g *discordgo.GuildCreate) {
	if g.Unavailable {
		log.Info("joined unavailable guild", g.Guild.ID)
		return
	}
	sqlite := OpenLiteDB(PathLiteDB)
	timejoin, err := g.Guild.JoinedAt.Parse()
	if err != nil {
		log.Error(err)
		return
	}
	DataGuild := Guild{
		ID:     g.Guild.ID,
		Name:   g.Guild.Name,
		Join:   timejoin,
		Dbconn: sqlite,
	}
	Info := DataGuild.CheckGuild()
	SendInvite, err := s.UserChannelCreate(config.OwnerDiscordID)
	if err != nil {
		log.Error(err)
	}

	if Info == 0 {
		for _, Channel := range g.Guild.Channels {
			BotPermission, err := s.UserChannelPermissions(BotID, Channel.ID)
			if err != nil {
				log.Error(err)
				return
			}
			if Channel.Type == 0 && BotPermission&2048 != 0 {
				s.ChannelMessageSendEmbed(Channel.ID, NewEmbed().
					SetTitle("Thx for invite me to this server <3 ").
					SetThumbnail(config.GoSimpIMG).
					SetImage(H3llcome[rand.Intn(len(H3llcome))]).
					SetColor(14807034).
					SetDescription("Type `"+config.PGeneral+"help` to show options").MessageEmbed)

				//send server name to my discord
				err := DataGuild.InputGuild()
				if err != nil {
					log.Error(err)
					return
				}
				s.ChannelMessageSend(SendInvite.ID, g.Guild.Name+" invited me")
				return
			}
		}
	} else {
		s.ChannelMessageSend(SendInvite.ID, g.Guild.Name+" reinvite me")
		err := DataGuild.UpdateJoin(Info)
		if err != nil {
			log.Error(err)
			return
		}
	}
	KillSqlite(sqlite)
}
