package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// Public variables
var (
	ModuleList  = []string{"LiveBiliBili", "SpaceBiliBili", "Youtube", "TwitterFanart", "BiliBiliFanart", "YoutubeSubscriber", "BiliBiliFollowers", "TwitterFollowers"}
	BotConf     ConfigFile
	TwitchToken string
	NotFound    string
	YoutubeIMG  string
	BiliBiliIMG string
	TwitterIMG  string
	WorryIMG    string
	GoSimpIMG   string
	Longcatttt  = "https://cdn.ebaumsworld.com/2020/09/20/013235/86394200/longcat-pic.jpg"
	Dead        string
	BSD         string
	Sleep       string
	Bonjour     string
	Howdy       string
	Guten       string
	Koni        string
	Selamat     string
	Assalamu    string
	Approaching string
	CommandURL  string
	GuideURL    string
	VtubersData string
)

const (
	TmpDir            = "/tmp/tmp.img"
	ChannelPermission = 8208
	GuildSupport      = "https://discord.com/invite/ydWC5knbJT"

	//Crontab
	TwitterFanart              = "@every 0h3m0s"
	BiliBiliFanart             = "@every 0h6m0s"
	BiliBiliLive               = "@every 0h7m0s"
	BiliBiliSpace              = "@every 0h13m0s"
	Twitch                     = "@every 0h7m0s"
	YoutubeCheckChannel        = "@every 0h10m30s"
	YoutubeCheckUpcomingByTime = "@every 0h1m0s"
	YoutubePrivateSlayer       = "@every 2h31m0s"
	YoutubeSubscriber          = "@every 1h0m0s"
	BiliBiliFollowers          = "@every 0h30m0s"
	TwitterFollowers           = "@every 0h17m0s"
	DonationMsg                = "@every 0h30m0s"
	CheckServerCount           = "@every 0h1m0s"
)

type ConfigFile struct {
	Discord        string `toml:"Discord"`
	BiliSess       string `toml:"BiliSess"`
	SauceAPI       string `toml:"SauceAPI"`
	DiscordWebHook string `toml:"DiscordWebHook"`
	MultiTOR       string `toml:"Multitor"`
	DonationLink   string `toml:"DonationLink"`
	TopGG          string `toml:"TOPGG"`
	Domain         string `toml:"Domain"`
	Twitch         struct {
		ClientID     string `toml:"ClientID"`
		ClientSecret string `toml:"ClientSecret"`
	} `toml:"Twitch"`
	LimitConf struct {
		TwitterFanart int `toml:"TwitterLimit"`
		SpaceBiliBili int `toml:"SpaceBiliBili"`
		YoutubeLimit  int `toml:"YoutubeLimit"`
	} `toml:"Limit"`
	SQL struct {
		User         string `toml:"User"`
		Pass         string `toml:"Pass"`
		Host         string `toml:"Host"`
		Port         string `toml:"Port"`
		MaxOpenConns int    `toml:"MaxOpenConns"`
		MaxIdleConns int    `toml:"MaxIdleConns"`
	} `toml:"Sql"`
	Cached struct {
		Host string `toml:"Host"`
		Port string `toml:"Port"`
	} `toml:"Cached"`
	BotPrefix struct {
		Fanart   string `toml:"Fanart"`
		Youtube  string `toml:"Youtube"`
		Bilibili string `toml:"Bilibili"`
		General  string `toml:"General"`
	} `toml:"BotPrefix"`
	Emoji struct {
		Fanart     []string `toml:"Fanart"`
		Livestream []string `toml:"Livestream"`
	} `toml:"Emoji"`
	YtToken []string `toml:"YoutubeToken"`
}

//read from config file
func ReadConfig(path string) (ConfigFile, error) {
	fmt.Println("Reading config file...")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return ConfigFile{}, err
	}

	fmt.Println(string(file))

	_, err = toml.Decode(string(file), &BotConf)
	if err != nil {
		return ConfigFile{}, err
	}

	NotFound = fmt.Sprintf("https://cdn.%s/404.jpg", BotConf.Domain)
	YoutubeIMG = fmt.Sprintf("https://cdn.%s/youtube.png", BotConf.Domain)
	BiliBiliIMG = fmt.Sprintf("https://cdn.%s/bilibili.png", BotConf.Domain)
	TwitterIMG = fmt.Sprintf("https://cdn.%s/twitter.png", BotConf.Domain)
	WorryIMG = fmt.Sprintf("https://cdn.%s/parerunworry.png", BotConf.Domain)
	GoSimpIMG = fmt.Sprintf("https://cdn.%s/go-simp.png", BotConf.Domain)
	Dead = fmt.Sprintf("https://cdn.%s/dead.jpg", BotConf.Domain)
	BSD = fmt.Sprintf("https://cdn.%s/bsd.png", BotConf.Domain)
	Sleep = fmt.Sprintf("https://cdn.%s/sleep.png", BotConf.Domain)
	Bonjour = fmt.Sprintf("https://cdn.%s/bonjour.png", BotConf.Domain)
	Howdy = fmt.Sprintf("https://cdn.%s/howdy.png", BotConf.Domain)
	Guten = fmt.Sprintf("https://cdn.%s/guten.png", BotConf.Domain)
	Koni = fmt.Sprintf("https://cdn.%s/koni.png", BotConf.Domain)
	Selamat = fmt.Sprintf("https://cdn.%s/selamat.jpg", BotConf.Domain)
	Assalamu = fmt.Sprintf("https://cdn.%s/Assalamu.jpg", BotConf.Domain)
	Approaching = fmt.Sprintf("https://cdn.%s/approaching.jpg", BotConf.Domain)
	CommandURL = fmt.Sprintf("https://go-simp.%s/Exec/", BotConf.Domain)
	GuideURL = fmt.Sprintf("https://go-simp.%s/Guide/", BotConf.Domain)
	VtubersData = fmt.Sprintf("https://go-simp.%s", BotConf.Domain)

	if BotConf.LimitConf.YoutubeLimit >= 15 {
		BotConf.LimitConf.YoutubeLimit = 15
	}
	return BotConf, nil
}

func (Data ConfigFile) CheckSQL() *sql.DB {
	log.Info("Open DB")

	db, err := sql.Open("mysql", Data.SQL.User+":"+Data.SQL.Pass+"@tcp("+Data.SQL.Host+":"+Data.SQL.Port+")/Vtuber?parseTime=true")
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(Data.SQL.MaxOpenConns)
	db.SetMaxIdleConns(Data.SQL.MaxIdleConns)

	//make sure can access database
	_, err = db.Exec(`SELECT NOW()`)
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}
	return db
}

func (Data ConfigFile) GetTwitchAccessToken() string {
	if TwitchToken != "" {
		return TwitchToken
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		var (
			url  = "https://id.twitch.tv/oauth2/token?client_id=" + Data.Twitch.ClientID + "&client_secret=" + Data.Twitch.ClientSecret + "&grant_type=client_credentials"
			body []byte
		)
		spaceClient := http.Client{}
		request, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			log.Error(err)
		}
		request.Header.Set("cache-control", "no-cache")
		request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9")

		response, err := spaceClient.Do(request.WithContext(ctx))
		if err != nil {
			log.Error(err)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.WithFields(log.Fields{
				"Status": response.StatusCode,
				"Reason": response.Status,
			}).Error("Status code not daijobu")
		}

		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			log.Error(err)

		}
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Error(err)
		}
		TwitchToken = result["access_token"].(string)
		return TwitchToken
	}
}
