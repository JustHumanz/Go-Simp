package engine

import (
	"context"
	"regexp"
	"strings"
	"time"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func SubsSendEmbed(Embed *discordgo.MessageEmbed, Group database.Group, Member database.Member, Bot *discordgo.Session) {
	if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
		Embed.Author.IconURL = ""
	}
	ChannelData, err := Group.GetChannelByGroup(Member.Region)
	if err != nil {
		log.Error(err)
	}
	for i, Channel := range ChannelData {

		Channel.SetMember(Member)
		Tmp := &Channel
		ctx := context.Background()
		UserTagsList, err := Tmp.SetMember(Member).SetGroup(Group).GetUserList(ctx)
		if err != nil {
			log.Error(err)
		}
		msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, Embed)
		if err != nil {
			log.Error(msg, err)
		}
		if UserTagsList != nil {
			msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
			}
		}

		Wait := GetMaxSqlConn()
		if i%Wait == 0 && i != 0 {
			log.WithFields(log.Fields{
				"Func":  "Subscriber",
				"Value": Wait,
			}).Warn("Waiting send message")
			time.Sleep(100 * time.Millisecond)
		}
	}
}
