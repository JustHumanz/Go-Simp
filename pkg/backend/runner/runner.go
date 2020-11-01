package runner

import (
	"fmt"
	"os"

	config "github.com/JustHumanz/Go-simp/config"
	//cronjob "github.com/JustHumanz/Go-simp/tools/cronjob"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	Bot *discordgo.Session
)

func StartInit(path string) {
	conf, err := config.ReadConfig(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Bot, _ = discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	database.Start(conf.CheckSQL())
	engine.Start()
	//cronjob.InitCron()
}
