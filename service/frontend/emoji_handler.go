package main

import (
	config "github.com/JustHumanz/Go-Simp/pkg/config"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func UpdateChannel(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == Register.Admin && m.MessageID == Register.MessageID {
		typ := func() {
			err := Register.UpdateChannel(config.Type)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				log.Error(err)
			}
		}
		all := func() {
			err := Register.UpdateChannel("all")
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				log.Error(err)
			}
		}

		if m.Emoji.MessageFormat() == config.One {
			Register.Stop()
			Register.ChoiceType()
			Register.BreakPoint(2)

			if Register.ChannelState.Group.GroupName == config.Indie {
				Register.Stop()
				Register.IndieNotif()
				Register.BreakPoint(1)
			}

			if Register.DisableState {
				typ()
				all()
			} else {
				Register.ChangeLiveStream()
				all()
				typ()
			}

			if Register.ChannelState.TypeTag == config.LewdType || Register.ChannelState.TypeTag == config.LewdNArtType {
				if !Register.CheckNSFW() {
					return
				}
			}

			_, err := s.ChannelMessageSend(m.ChannelID, "Done")
			if err != nil {
				log.Error(err)
			}
			CleanRegister()
			return

		} else if m.Emoji.MessageFormat() == config.Two {
			Register.AddRegion()
		} else if m.Emoji.MessageFormat() == config.Three {
			Register.DelRegion()
		} else if m.Emoji.MessageFormat() == config.Four {
			Register.EmojiTrue().ChangeLiveStream().DisableChannelState(false)
			all()
			_, err := s.ChannelMessageSend(m.ChannelID, "Done")
			if err != nil {
				log.Error(err)
			}
			return
		}

		if Register.State == AddRegion {
			Register.Start()
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				Register.AddNewRegion(Region)
			}
		} else if Register.State == DelRegion {
			Register.Start()
			Region := engine.UniCodetoCountryCode(m.Emoji.MessageFormat())
			if Region != "" {
				Register.RemoveRegion(Region)
			}
		}
	}
}

const (
	AddRegion = "addreg"
	DelRegion = "delreg"
)

func EmojiHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == Register.Admin && m.MessageID == Register.MessageID && Register.Emoji {
		Register.Start()
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
			Register.ChannelState.TypeTag = def
			CleanRegister()
		}

		NillType := func() {
			_, err := s.ChannelMessageSend(m.ChannelID, "[Error] you can't disable all type\nadios")
			if err != nil {
				log.Error(err)
			}
			CleanRegister()
		}

		if m.Emoji.MessageFormat() == config.Art {
			if Register.ChannelState.TypeTag == config.LiveType {
				Message("[Info] <@" + m.UserID + "> enable Livestream and Art type on this channel")
				Register.UpdateType(config.ArtNLiveType).DisableChannelState(false)
			} else if Register.ChannelState.TypeTag == config.LewdType {
				Message("[Info] <@" + m.UserID + "> enable Lewd and Art type on this channel")
				Register.UpdateType(config.LewdNArtType).DisableChannelState(false)
			} else if Register.ChannelState.TypeTag == config.ArtType {
				NillType()
				return
			} else if Register.ChannelState.TypeTag == config.ArtNLiveType {
				Message("[Info] <@" + m.UserID + "> disable Art type on this channel")
				Register.UpdateType(config.LiveType).DisableChannelState(true)
				return
			} else {
				Message("[Info] <@" + m.UserID + "> enable Art type on this channel")
				Register.UpdateType(config.ArtType)
			}
		} else if m.Emoji.MessageFormat() == config.Live {
			if Register.ChannelState.TypeTag == config.ArtType {
				Message("[Info] <@" + m.UserID + "> enable Livestream and Art type on this channel")
				Register.UpdateType(config.ArtNLiveType).DisableChannelState(false)
			} else if Register.ChannelState.TypeTag == config.LewdType {
				LewdLive(config.LewdType)
				return
			} else if Register.ChannelState.TypeTag == config.LiveType {
				NillType()
				return
			} else if Register.ChannelState.TypeTag == config.ArtNLiveType {
				Message("[Info] <@" + m.UserID + "> disable Live type on this channel")
				Register.UpdateType(config.ArtType).DisableChannelState(true)
				return
			} else {
				Message("[Info] <@" + m.UserID + "> enable Livestream type on this channel")
				Register.UpdateType(config.LiveType)
			}
		} else if m.Emoji.MessageFormat() == config.Lewd {
			if Register.ChannelState.TypeTag == config.LiveType {
				LewdLive(config.LiveType)
				return
			} else if Register.ChannelState.TypeTag == config.ArtType {
				Register.UpdateType(config.LewdNArtType)
			} else if Register.ChannelState.TypeTag == config.ArtNLiveType {
				LewdLive(config.ArtNLiveType)
				return
			} else if Register.ChannelState.TypeTag == config.LewdType {
				NillType()
				return
			} else if Register.ChannelState.TypeTag == config.LewdNArtType {
				Message("[Info] <@" + m.UserID + "> disable Lewd type on this channel")
				Register.UpdateType(config.ArtType).DisableChannelState(true)
				return
			} else {
				Register.UpdateType(config.LewdType)
			}
		}

		if m.Emoji.MessageFormat() == config.Ok {
			if Register.State == config.LiveOnly {
				Register.SetLiveOnly(true)
			} else if Register.State == config.NewUpcoming {
				Register.SetNewUpcoming(true)
			} else if Register.State == config.Dynamic {
				Register.SetDynamic(true)
			} else if Register.State == config.LiteMode {
				Register.SetLite(true)
			} else if Register.State == config.IndieNotif {
				Register.SetIndieNotif(true)
			}
		} else if m.Emoji.MessageFormat() == config.No {
			if Register.State == config.LiveOnly {
				Register.SetLiveOnly(false)
			} else if Register.State == config.NewUpcoming {
				Register.SetNewUpcoming(false)
			} else if Register.State == config.Dynamic {
				Register.SetDynamic(false)
			} else if Register.State == config.LiteMode {
				Register.SetLite(false)
			} else if Register.State == config.IndieNotif {
				Register.SetIndieNotif(false)
				_, err := s.ChannelMessageSend(m.ChannelID, "[Tips] create a dummy role and tag that role use `"+configfile.BotPrefix.General+"tag roles` command")
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
