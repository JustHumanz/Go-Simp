package runner

import (
	"flag"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
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
	BotInfo, err = Bot.User("@me")
	if err != nil {
		return err
	}

	Bot.AddHandler(ModuleInfo)

	database.Start(conf.CheckSQL())
	engine.Start()

	return nil
	//cronjob.InitCron()
}

//ModuleInfo send user about module info
func ModuleInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	m.Content = strings.ToLower(m.Content)
	if strings.HasPrefix(m.Content, config.PGeneral+"module") {
		Color, err := engine.GetColor("/tmp/mem.tmp", m.Author.AvatarURL("128"))
		if err != nil {
			log.Error(err)
		}

		var (
			TwitterFanart  string
			BiliBiliFanart string
			LiveBiliBili   string
			SpaceBiliBili  string
			Youtube        string
			Subscriber     string
		)

		if *Flags["twitterfanart"].(*bool) {
			TwitterFanart = "✓"
		} else {
			TwitterFanart = "✘"
		}

		if *Flags["bilibilifanart"].(*bool) {
			BiliBiliFanart = "✓"
		} else {
			BiliBiliFanart = "✘"
		}

		if *Flags["youtube"].(*bool) {
			Youtube = "✓"
		} else {
			Youtube = "✘"
		}

		if *Flags["bilibili"].(*bool) {
			LiveBiliBili = "✓"
		} else {
			Youtube = "✘"
		}

		if *Flags["space.bilibili"].(*bool) {
			SpaceBiliBili = "✓"
		} else {
			SpaceBiliBili = "✘"
		}

		if *Flags["subscriber"].(*bool) {
			Subscriber = "✓"
		} else {
			Subscriber = "✘"
		}

		s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
			SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
			SetTitle(BotInfo.Username).
			SetImage(config.GoSimpIMG).
			SetThumbnail(BotInfo.AvatarURL("128")).
			SetDescription("Bot Module").
			AddField("Youtube notif", Youtube).
			AddField("Live BiliBili notif", LiveBiliBili).
			AddField("Space BiliBili notif", SpaceBiliBili).
			AddField("Twitter fanart", TwitterFanart).
			AddField("BiliBili fanart", BiliBiliFanart).
			AddField("Subscriber count notif", Subscriber).
			InlineAllFields().
			SetColor(Color).MessageEmbed)
	}
}
