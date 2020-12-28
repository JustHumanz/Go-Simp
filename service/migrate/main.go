package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	youtube "github.com/JustHumanz/Go-simp/pkg/backend/livestream/youtube"
	"github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	JsonData    Vtuber
	Limit       int
	db          *sql.DB
	YtToken     string
	Publish     time.Time
	Roomstatus  string
	BiliSession string
	Bot         *discordgo.Session
)

type NewVtuber struct {
	Member Member
	Group  database.Group
}

func init() {
	fmt.Println("Reading hashtag file...")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	file, err := ioutil.ReadFile("./vtuber.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &JsonData)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	config, err := config.ReadConfig("../../config.toml")
	if err != nil {
		log.Error(err)
	}
	YtToken = config.YtToken[0]
	BiliSession = config.BiliSess
	Limit = 100
	Bot, _ = discordgo.New("Bot " + config.Discord)
	db = config.CheckSQL()
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = CreateDB(config)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	Bot.AddHandler(Dead)
}

func main() {
	database.Start(db)
	AddData(JsonData)
	go CheckYT()
	go CheckSchedule()
	go CheckVideoSpace()
	go CheckTBili()
	go youtube.CheckPrivate()
	go TwitterFanart()

	log.Info("Done")
	time.Sleep(6 * time.Minute)
	os.Exit(0)
}

func (Data NewVtuber) SendNotif() *discordgo.MessageEmbed {
	var (
		Twitterfanart  string
		Bilibilifanart string
		Bilibili       string
		Youtube        string
		URL            string
		Color          int
		Avatar         string
		err            error
	)

	if Data.Member.YtID != "" {
		Youtube = "✓"
		URL = "https://www.youtube.com/channel/" + Data.Member.YtID + "?sub_confirmation=1"

		Avatar = Data.Member.YtAvatar()
		Color, err = engine.GetColor("/tmp/notf.gg", Avatar)
		if err != nil {
			log.Error(err)
		}

	} else {
		Youtube = "✘"
		URL = "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
		Avatar = Data.Member.BliBiliFace()
		Color, err = engine.GetColor("/tmp/notf.gg", Avatar)
		if err != nil {
			log.Error(err)
		}
	}

	if Data.Member.Hashtag.Twitter != "" {
		Twitterfanart = "✓"
	} else {
		Twitterfanart = "✘"
	}

	if Data.Member.Hashtag.BiliBili != "" {
		Bilibilifanart = "✓"
	} else {
		Bilibilifanart = "✘"
	}

	if Data.Member.BiliRoomID != 0 {
		Bilibili = "✓"
	} else {
		Bilibili = "✘"
	}

	return engine.NewEmbed().
		SetAuthor(Data.Group.GroupName, Data.Group.IconURL).
		SetTitle(engine.FixName(Data.Member.ENName, Data.Member.JPName)).
		SetImage(Avatar).
		SetThumbnail("https://justhumanz.me/update.png").
		SetDescription("New Vtuber has been added to list").
		AddField("Nickname", Data.Member.Name).
		AddField("Region", Data.Member.Region).
		AddField("Twitter Fanart", Twitterfanart).
		AddField("BiliBili Fanart", Bilibilifanart).
		AddField("Youtube Notif", Youtube).
		AddField("BiliBili Notif", Bilibili).
		InlineAllFields().
		SetURL(URL).
		SetColor(Color).MessageEmbed
}

func Dead(s *discordgo.Session, m *discordgo.MessageCreate) {
	General := config.BotConf.BotPrefix.General
	Fanart := config.BotConf.BotPrefix.Fanart
	BiliBili := config.BotConf.BotPrefix.Bilibili
	Youtube := config.BotConf.BotPrefix.Youtube
	m.Content = strings.ToLower(m.Content)
	Color, err := engine.GetColor("/tmp/discordpp.tmp", m.Author.AvatarURL("128"))
	if err != nil {
		log.Error(err)
	}
	if m.Content != "" {
		if len(regexp.MustCompile("(?m)("+General+"|"+Fanart+"|"+BiliBili+"|"+Youtube+")").FindAllString(m.Content, -1)) > 0 {
			s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Bot update new Vtubers").
				SetURL("https://github.com/JustHumanz/Go-Simp/blob/main/CHANGELOG.md").
				SetDescription("Still Processing new data,Comeback when i ready to bang you (around 10-20 minutes or more,~~idk i don't fvcking count~~)").
				AddField("See update at", "[Changelog](https://github.com/JustHumanz/Go-Simp/blob/main/CHANGELOG.md)").
				SetThumbnail(config.Sleep).
				SetImage(config.Dead).
				SetColor(Color).
				SetFooter("Adios~").MessageEmbed)
		}
	}
}
