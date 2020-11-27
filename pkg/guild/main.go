package guild

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	H3llcome  = []string{config.Bonjour, config.Howdy, config.Guten, config.Koni, config.Selamat, config.Assalamu, config.Approaching}
	GuildList []string
	BotID     *discordgo.User
)

func main() {
	conf, err := config.ReadConfig("../../config.toml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Bot, _ := discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	BotID, err = Bot.User("@me")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	database.Start(conf.CheckSQL())
	engine.Start()

	for _, GuildID := range Bot.State.Guilds {
		GuildList = append(GuildList, GuildID.ID)
	}

	Bot.AddHandler(GuildJoin)
	log.Info("Guild handler ready.......")

	chain := make(chan os.Signal, 0)
	signal.Notify(chain, os.Interrupt)
	go func() {
		for sig := range chain {
			log.Warn("captured ", sig, ", stopping profiler and exiting..")
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	<-make(chan struct{})
	return
}
