package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/nicklaw5/helix"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	JsonData        Vtuber
	Limit           int
	db              *sql.DB
	YoutubeToken    *string
	Publish         time.Time
	Roomstatus      string
	BiliBiliSession map[string]string
	Bot             *discordgo.Session
	TwitchClient    *helix.Client
	configfile      config.ConfigFile
	gRCPconn        pilot.PilotServiceClient
)

type NewVtuber struct {
	Member     Members
	Group      database.Group
	YtAvatar   string
	BiliAvatar string
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
	fmt.Println("Reading json files...")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	files, err := ioutil.ReadDir("json/")
	if err != nil {
		log.Fatal(err)
	}

	var (
		Indie  Independent
		Groups []Group
	)

	for _, f := range files {
		file, err := ioutil.ReadFile("json/" + f.Name())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(file))

		if f.Name() == "Independent.json" {
			err = json.Unmarshal(file, &Indie)
			if err != nil {
				fmt.Println(err)
			}

		} else {
			var Gtmp Group
			err = json.Unmarshal(file, &Gtmp)
			if err != nil {
				fmt.Println(err)
			}
			Groups = append(Groups, Gtmp)
		}
	}
	JsonData = Vtuber{
		VtuberData: Data{
			Independent: Indie,
			Group:       Groups,
		},
	}

	gRCPconn = pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))
	RequestPay("Start migrate new vtuber")

	db = configfile.CheckSQL()
	database.Start(configfile)
	configfile.InitConf()

	YoutubeToken = engine.GetYtToken()
	BiliBiliSession = map[string]string{
		"Cookie": "SESSDATA=" + configfile.BiliSess,
	} //[]string{"Cookie", "SESSDATA=" + configfile.BiliSess}
	Limit = 100

	Bot = engine.StartBot(false)
	err = Bot.Open()
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

	resp, err := TwitchClient.RequestAppAccessToken([]string{"user:read:email"})
	if err != nil {
		log.Error(err)
	}

	TwitchClient.SetAppAccessToken(resp.Data.AccessToken)

	Bot.AddHandler(Dead)
	err = Bot.UpdateStreamingStatus(0, "Maintenance!!!!", config.VtubersData)
	if err != nil {
		log.Error(err)
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Shutting down...")
		if Bot != nil {
			log.Info("Close bot wss session")
			err := Bot.Close()
			if err != nil {
				log.Error(err)
			}

			log.Info("Migrate killed")
			RequestPay("Done migrate new vtuber")
			os.Exit(1)
		}
	}()
	AddData(JsonData)
	time.Sleep(1 * time.Minute)

	go CheckYoutube()
	go CheckLiveBiliBili()
	go CheckTwitch()
	time.Sleep(5 * time.Minute)

	go CheckSpaceBiliBili()
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

	Data.Group.RemoveNillIconURL()

	if Data.Member.Youtube.YtID != "" {
		Youtube = "✓"
		URL = "https://www.youtube.com/channel/" + Data.Member.Youtube.YtID + "?sub_confirmation=1"

		Color, err = engine.GetColor(config.TmpDir, Data.YtAvatar)
		if err != nil {
			log.Error(err)
		}

	} else {
		Youtube = "✘"
		URL = "https://space.bilibili.com/" + strconv.Itoa(Data.Member.BiliBili.BiliBiliID)
		Color, err = engine.GetColor(config.TmpDir, Data.BiliAvatar)
		if err != nil {
			log.Error(err)
		}
	}

	if Data.Member.Twitter.TwitterFanart != "" {
		Twitterfanart = "✓"
	} else {
		Twitterfanart = "✘"
	}

	if Data.Member.BiliBili.BiliBiliFanart != "" {
		Bilibilifanart = "✓"
	} else {
		Bilibilifanart = "✘"
	}

	if Data.Member.BiliBili.BiliRoomID != 0 {
		Bilibili = "✓"
	} else {
		Bilibili = "✘"
	}

	if Data.Member.Twitch.TwitchUsername != "" {
		Twitch = "✓"
	} else {
		Twitch = "✘"
	}

	return engine.NewEmbed().
		SetAuthor(Data.Group.GroupName, Data.Group.IconURL).
		SetTitle(engine.FixName(Data.Member.EnName, Data.Member.JpName)).
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
		if len(regexp.MustCompile("(?m)("+General+"|"+Fanart+"|"+BiliBili+"|"+Youtube+")").FindAllString(m.Content, -1)) > 0 && !m.Author.Bot {
			s.ChannelMessageSendEmbed(m.ChannelID, engine.NewEmbed().
				SetAuthor(m.Author.Username, m.Author.AvatarURL("128")).
				SetTitle("Bot update new Vtubers").
				SetURL("https://github.com/JustHumanz/Go-Simp/blob/master/CHANGELOG.md").
				SetDescription("Still Processing new data,Comeback when i ready to bang you (around 10-20 minutes or more)").
				AddField("See update at", "[Changelog](https://github.com/JustHumanz/Go-Simp/blob/master/CHANGELOG.md)").
				SetThumbnail(config.Sleep).
				SetImage(engine.MaintenanceIMG()).
				SetColor(Color).
				SetFooter("Adios~").MessageEmbed)
		}
	}
}
