package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/nicklaw5/helix"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	youtube "github.com/JustHumanz/Go-Simp/service/livestream/youtube"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	JsonData        Vtuber
	Limit           int
	db              *sql.DB
	YoutubeToken    string
	Publish         time.Time
	Roomstatus      string
	BiliBiliSession []string
	Bot             *discordgo.Session
	TwitchClient    *helix.Client
	TwitchToken     string
	configfile      config.ConfigFile
	gRCPconn        pilot.PilotServiceClient
)

type NewVtuber struct {
	Member Member
	Group  database.Group
}

func RequestPay(Message string) {
	res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
		Message: Message,
		Service: "Migrate",
	})
	if err != nil {
		log.Fatalf("Error when request payload: %s", err)
	}
	err = json.Unmarshal(res.ConfigFile, &configfile)
	if err != nil {
		log.Panic(err)
	}
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

	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	RequestPay("Start migrate new vtuber")

	Bot, err = discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	db = configfile.CheckSQL()
	database.DB = db
	config.GoSimpConf = configfile

	YoutubeToken = engine.GetYtToken()
	BiliBiliSession = []string{"Cookie", "SESSDATA=" + configfile.BiliSess}
	Limit = 100
	TwitchToken = configfile.GetTwitchAccessToken()
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = CreateDB(configfile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	TwitchClient, err = helix.NewClient(&helix.Options{
		ClientID:     configfile.Twitch.ClientID,
		ClientSecret: configfile.Twitch.ClientSecret,
	})
	if err != nil {
		log.Error(err)
	}
	TwitchClient.SetUserAccessToken(TwitchToken)

	Bot.AddHandler(Dead)
}

func main() {
	AddData(JsonData)
	go CheckYoutube()
	go CheckLiveBiliBili()
	go CheckTwitch()
	time.Sleep(5 * time.Minute)

	go CheckSpaceBiliBili()
	go CheckTBiliBili()
	time.Sleep(5 * time.Minute)

	go youtube.CheckPrivate()
	go TwitterFanart()
	time.Sleep(5 * time.Minute)
	log.Info("Done")
	RequestPay("Done migrate new vtuber")
	os.Exit(0)
}

func (Data NewVtuber) SendNotif() *discordgo.MessageEmbed {
	var (
		Twitterfanart  string
		Bilibilifanart string
		Bilibili       string
		Twitch         string
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
		Color, err = engine.GetColor(config.TmpDir, Avatar)
		if err != nil {
			log.Error(err)
		}

	} else {
		Youtube = "✘"
		URL = "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBiliID)
		Avatar, err = Data.Member.BliBiliFace()
		if err != nil {
			log.Error(err)
		}
		Color, err = engine.GetColor(config.TmpDir, Avatar)
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

	if Data.Member.TwitchName != "" {
		Twitch = "✓"
	} else {
		Twitch = "✘"
	}

	return engine.NewEmbed().
		SetAuthor(Data.Group.GroupName, Data.Group.IconURL).
		SetTitle(engine.FixName(Data.Member.ENName, Data.Member.JPName)).
		SetImage(Avatar).
		SetDescription("New Vtuber has been added to list").
		AddField("Nickname", Data.Member.Name).
		AddField("Region", Data.Member.Region).
		AddField("Twitter Fanart", Twitterfanart).
		AddField("BiliBili Fanart", Bilibilifanart).
		AddField("Youtube Notif", Youtube).
		AddField("BiliBili Notif", Bilibili).
		AddField("Twitch Notif", Twitch).
		InlineAllFields().
		SetURL(URL).
		SetColor(Color).MessageEmbed
}

func Dead(s *discordgo.Session, m *discordgo.MessageCreate) {
	General := configfile.BotPrefix.General
	Fanart := configfile.BotPrefix.Fanart
	BiliBili := configfile.BotPrefix.Bilibili
	Youtube := configfile.BotPrefix.Youtube
	m.Content = strings.ToLower(m.Content)
	Color, err := engine.GetColor(config.TmpDir, m.Author.AvatarURL("128"))
	if err != nil {
		log.Error(err)
	}
	if m.Content != "" {
		if len(regexp.MustCompile("(?m)("+General+"|"+Fanart+"|"+BiliBili+"|"+Youtube+")").FindAllString(m.Content, -1)) > 0 {
			s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Bot update new Vtubers").
				SetURL("https://github.com/JustHumanz/Go-Simp/blob/master/CHANGELOG.md").
				SetDescription("Still Processing new data,Comeback when i ready to bang you (around 10-20 minutes or more,~~idk i don't fvcking count~~)").
				AddField("See update at", "[Changelog](https://github.com/JustHumanz/Go-Simp/blob/master/CHANGELOG.md)").
				SetThumbnail(config.Sleep).
				SetImage(config.Dead).
				SetColor(Color).
				SetFooter("Adios~").MessageEmbed)
		}
	}
}
