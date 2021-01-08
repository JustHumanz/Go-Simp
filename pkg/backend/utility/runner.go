package main

import (
	"os"

	"github.com/JustHumanz/Go-simp/pkg/backend/utility/runfunc"
	"github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	"github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/top-gg/go-dbl"
	"gopkg.in/robfig/cron.v2"
)

var (
	BotID string
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
	BotID = BotInfo.ID

	CheckServerCount()
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
				AddField("if you a broke gang,you can upvote "+BotInfo.Username, "[top.gg]("+config.BotConf.TopGG+")").
				AddField("or give some star on github", "[Github](https://github.com/JustHumanz/Go-Simp)").MessageEmbed)
		})
	}
	c.AddFunc("@every 0h1m0s", CheckServerCount)

	runfunc.Run()
}

func CheckServerCount() {
	dblClient, err := dbl.NewClient(os.Getenv("TOPGG"))
	if err != nil {
		log.Error(err)
	}

	err = dblClient.PostBotStats(BotID, &dbl.BotStatsPayload{
		Shards: []int{database.GetGuildsCount()},
	})
	if err != nil {
		log.Error(err)
	}
}
