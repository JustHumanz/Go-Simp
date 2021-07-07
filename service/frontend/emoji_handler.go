package main

import (
	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func UpdateChannel(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	for j, RegisterPayload := range Register.Payload {
		if RegisterPayload.AdminID == m.UserID+m.ChannelID {
			typ := func() {
				err := RegisterPayload.UpdateChannel(config.Type)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					log.Error(err)
				}
			}
			all := func() {
				err := RegisterPayload.UpdateChannel("all")
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					log.Error(err)
				}
			}

			if m.Emoji.MessageFormat() == config.One {
				RegisterPayload.Stop()
				RegisterPayload.ChoiceType()
				RegisterPayload.BreakPoint(2)

				if RegisterPayload.ChannelState.Group.GroupName == config.Indie {
					RegisterPayload.Stop()
					RegisterPayload.IndieNotif()
					RegisterPayload.BreakPoint(1)
				}

				if RegisterPayload.DisableState {
					typ()
					all()
				} else {
					RegisterPayload.ChangeLiveStream()
					all()
					typ()
				}

				if RegisterPayload.ChannelState.TypeTag == config.LewdType || RegisterPayload.ChannelState.TypeTag == config.LewdNArtType {
					if !RegisterPayload.CheckNSFW() {
						return
					}
				}

				_, err := s.ChannelMessageSend(m.ChannelID, "Done")
				if err != nil {
					log.Error(err)
				}
				CleanRegister(j)
				return

			} else if m.Emoji.MessageFormat() == config.Two {
				RegisterPayload.AddRegion()
			} else if m.Emoji.MessageFormat() == config.Three {
				RegisterPayload.DelRegion()
			} else if m.Emoji.MessageFormat() == config.Four {
				RegisterPayload.EmojiTrue().ChangeLiveStream().DisableChannelState(false)
				all()
				_, err := s.ChannelMessageSend(m.ChannelID, "Done")
				if err != nil {
					log.Error(err)
				}
				return
			}

			if RegisterPayload.State == AddRegion {
				RegisterPayload.Start()
				Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
				if Region != "" {
					RegisterPayload.AddNewRegion(Region)
				}
			} else if RegisterPayload.State == DelRegion {
				RegisterPayload.Start()
				Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
				if Region != "" {
					RegisterPayload.RemoveRegion(Region)
				}
			}
		}
	}
}

const (
	AddRegion = "addreg"
	DelRegion = "delreg"
)

func EmojiHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	for j, RegisterPayload := range Register.Payload {
		if RegisterPayload.AdminID == m.UserID+m.ChannelID && RegisterPayload.Emoji {
			RegisterPayload.Start()
			Message := func(msg string) {
				_, err := s.ChannelMessageSend(m.ChannelID, msg)
				if err != nil {
					log.Error(err)
				}
			}

			LewdLive := func(def int) {
				_, err := s.ChannelMessageSend(m.ChannelID, "[Error] you can't add livestream with lewd in same channel,canceling lewd")
				if err != nil {
					log.Error(err)
				}
				RegisterPayload.ChannelState.TypeTag = def
				CleanRegister(j)
			}

			NillType := func() {
				_, err := s.ChannelMessageSend(m.ChannelID, "[Error] you can't disable all type\nadios")
				if err != nil {
					log.Error(err)
				}
				CleanRegister(j)
			}

			if m.Emoji.MessageFormat() == config.Art {
				if RegisterPayload.ChannelState.TypeTag == config.LiveType {
					Message("[Info] <@" + m.UserID + "> enable Livestream and Art type on this channel")
					RegisterPayload.UpdateType(config.ArtNLiveType).DisableChannelState(false)
				} else if RegisterPayload.ChannelState.TypeTag == config.LewdType {
					Message("[Info] <@" + m.UserID + "> enable Lewd and Art type on this channel")
					RegisterPayload.UpdateType(config.LewdNArtType).DisableChannelState(false)
				} else if RegisterPayload.ChannelState.TypeTag == config.ArtType {
					NillType()
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.ArtNLiveType {
					Message("[Info] <@" + m.UserID + "> disable Art type on this channel")
					RegisterPayload.UpdateType(config.LiveType).DisableChannelState(true)
					return
				} else {
					Message("[Info] <@" + m.UserID + "> enable Art type on this channel")
					RegisterPayload.UpdateType(config.ArtType)
				}
			} else if m.Emoji.MessageFormat() == config.Live {
				if RegisterPayload.ChannelState.TypeTag == config.ArtType {
					Message("[Info] <@" + m.UserID + "> enable Livestream and Art type on this channel")
					RegisterPayload.UpdateType(config.ArtNLiveType).DisableChannelState(false)
				} else if RegisterPayload.ChannelState.TypeTag == config.LewdType {
					LewdLive(config.LewdType)
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.LiveType {
					NillType()
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.ArtNLiveType {
					Message("[Info] <@" + m.UserID + "> disable Live type on this channel")
					RegisterPayload.UpdateType(config.ArtType).DisableChannelState(true)
					return
				} else {
					Message("[Info] <@" + m.UserID + "> enable Livestream type on this channel")
					RegisterPayload.UpdateType(config.LiveType)
				}
			} else if m.Emoji.MessageFormat() == config.Lewd {
				if RegisterPayload.ChannelState.TypeTag == config.LiveType {
					LewdLive(config.LiveType)
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.ArtType {
					RegisterPayload.UpdateType(config.LewdNArtType)
				} else if RegisterPayload.ChannelState.TypeTag == config.ArtNLiveType {
					LewdLive(config.ArtNLiveType)
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.LewdType {
					NillType()
					return
				} else if RegisterPayload.ChannelState.TypeTag == config.LewdNArtType {
					Message("[Info] <@" + m.UserID + "> disable Lewd type on this channel")
					RegisterPayload.UpdateType(config.ArtType).DisableChannelState(true)
					return
				} else {
					RegisterPayload.UpdateType(config.LewdType)
				}
			}

			if m.Emoji.MessageFormat() == config.Ok {
				if RegisterPayload.State == config.LiveOnly {
					RegisterPayload.SetLiveOnly(true)
				} else if RegisterPayload.State == config.NewUpcoming {
					RegisterPayload.SetNewUpcoming(true)
				} else if RegisterPayload.State == config.Dynamic {
					RegisterPayload.SetDynamic(true)
				} else if RegisterPayload.State == config.LiteMode {
					RegisterPayload.SetLite(true)
				} else if RegisterPayload.State == config.IndieNotif {
					RegisterPayload.SetIndieNotif(true)
				}
			} else if m.Emoji.MessageFormat() == config.No {
				if RegisterPayload.State == config.LiveOnly {
					RegisterPayload.SetLiveOnly(false)
				} else if RegisterPayload.State == config.NewUpcoming {
					RegisterPayload.SetNewUpcoming(false)
				} else if RegisterPayload.State == config.Dynamic {
					RegisterPayload.SetDynamic(false)
				} else if RegisterPayload.State == config.LiteMode {
					RegisterPayload.SetLite(false)
				} else if RegisterPayload.State == config.IndieNotif {
					RegisterPayload.SetIndieNotif(false)
					_, err := s.ChannelMessageSend(m.ChannelID, "[Tips] create a dummy role and tag that role use `"+configfile.BotPrefix.General+"tag roles` command")
					if err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}
