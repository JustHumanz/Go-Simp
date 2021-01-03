package main

import (
	"github.com/JustHumanz/Go-simp/pkg/backend/utility/runfunc"
	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	conf, err := config.ReadConfig("../../../config.toml")
	if err != nil {
		log.Panic(err)
	}

	Bot, _ := discordgo.New("Bot " + config.BotConf.Discord)
	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}
	BotInfo, err := Bot.User("@me")
	if err != nil {
		log.Panic(err)
	}

	Donation := config.BotConf.DonationLink
	database.Start(conf.CheckSQL())
	engine.Start()

	c := cron.New()
	c.Start()
	if config.BotConf.DonationLink != "" {
		c.AddFunc("@every 0h30m0s", func() {
			Bot.ChannelMessageSendEmbed(database.GetRanChannel(), engine.NewEmbed().
				SetTitle("Donate").
				SetURL(Donation).
				SetThumbnail(BotInfo.AvatarURL("128")).
				SetImage(config.GoSimpIMG).
				SetColor(14807034).
				SetDescription("Enjoy the bot?\ndon't forget to support this bot and dev").
				AddField("Ko-Fi", "[Link]("+Donation+")").
				AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[top.gg]("+config.BotConf.TopGG+")").MessageEmbed)
		})
	}

	runfunc.Run()
}
