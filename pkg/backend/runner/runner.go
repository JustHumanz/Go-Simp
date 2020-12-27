package runner

import (
	"flag"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
)

//Public variable for Discord bot session
var (
	Bot     *discordgo.Session
	BotInfo *discordgo.User
	Flags   map[string]interface{}
)

//StartInit Start running BE
func StartInit(path string) error {
	conf, err := config.ReadConfig(path)
	if err != nil {
		return err
	}
	Flags = make(map[string]interface{})

	Flags["twitterfanart"] = flag.Bool("TwitterFanart", false, "Enable twitter fanart module")
	Flags["bilibilifanart"] = flag.Bool("BiliBiliFanart", false, "Enable bilibili fanart module")

	Flags["youtube"] = flag.Bool("Youtube", false, "Enable youtube module")
	Flags["bilibili"] = flag.Bool("Bilibili", false, "Enable bilibili module")
	Flags["space.bilibili"] = flag.Bool("Space.Bilibili", false, "Enable space.bilibili module")

	Flags["subscriber"] = flag.Bool("Subscriber", false, "Enable subscriber module")

	flag.Parse()

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
