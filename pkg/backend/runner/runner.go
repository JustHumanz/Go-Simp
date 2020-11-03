package runner

import (
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
)

var (
	Bot *discordgo.Session
)

func StartInit(path string) error {
	conf, err := config.ReadConfig(path)
	if err != nil {
		return err
	}

	Bot, _ = discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		return err
	}

	database.Start(conf.CheckSQL())
	engine.Start()

	return nil
	//cronjob.InitCron()
}
