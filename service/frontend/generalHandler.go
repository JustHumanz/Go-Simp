package main

import (
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//Help helmp command message handler
func Help(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	Prefix := configfile.BotPrefix.General
	if strings.HasPrefix(m.Content, Prefix) {
		Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}
		if m.Content == Prefix+"help en" || m.Content == Prefix+"help" {
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Help").
				SetURL(config.CommandURL).
				SetDescription("A simple VTuber bot which pings you or your roles if any new Videos, Fanarts, or Livestreams and Upcoming streams of VTubers are posted!").
				AddField("Command list", "[Exec]("+config.CommandURL+")").
				AddField("Guide", "[Guide]("+config.GuideURL+")").
				AddField("Vtuber list", "[Vtubers]("+config.VtubersData+")").
				AddField("Made by Golang", "[Go-Simp](https://github.com/JustHumanz/Go-Simp)").
				AddField("Bot WebApp", "[Web](web-admin.humanz.moe/Login)").
				AddField("Server count", strconv.Itoa(len(GuildList))).
				AddField("Member count", strconv.Itoa(database.GetMemberCount())).
				InlineAllFields().
				AddField("Join Dev server", "[Invite](https://discord.com/invite/ydWC5knbJT)").
				SetThumbnail(config.BSD).
				SetFooter("Prefix command already deprecated,use slash command").
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
		} else if m.Content == Prefix+Kings {
			s.ChannelMessageSend(m.ChannelID, "https://github.com/JustHumanz/Go-Simp/blob/master/King.md")
		} else if m.Content == Prefix+Upvote {
			s.ChannelMessageSend(m.ChannelID, config.GoSimpConf.TopGG)
		}
	}
}

func SendError(Messsage map[string]string) {
	_, err := Bot.ChannelMessageSendEmbed(Messsage["ChannelID"], engine.NewEmbed().
		SetAuthor(Messsage["Username"], Messsage["AvatarURL"]).
		SetTitle(Messsage["GroupName"]).
		SetDescription("Internal error XD").
		SetImage(engine.NotFoundIMG()).MessageEmbed)
	if err != nil {
		log.Error(err)
	}
}
