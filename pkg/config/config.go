package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/nicklaw5/helix"
	log "github.com/sirupsen/logrus"
)

// Public variables
var (
	GoSimpConf  ConfigFile
	TwitchToken string
	NotFound    string
	YoutubeIMG  string
	BiliBiliIMG string
	TwitchIMG   string
	TwitterIMG  string
	PixivIMG    string
	WorryIMG    string
	GoSimpIMG   string
	Longcatttt  = "https://cdn.ebaumsworld.com/2020/09/20/013235/86394200/longcat-pic.jpg"
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
	Scraper     *twitterscraper.Scraper
	Ytwaiting   = "???"
	CdnDomain   string

	//LewdBlacklist
	BlackList = []string{"loli", "gore"}
)

//Const config
const (
	TmpDir            = "/tmp/tmp.img"
	ChannelPermission = 8208
	BotPermission     = 452624
	GuildSupport      = "https://discord.com/invite/ydWC5knbJT"
	PixivIllustsEnd   = "https://www.pixiv.net/ajax/illust/"
	PixivUserEnd      = "https://www.pixiv.net/ajax/user/"
	PixivProxy        = "https://cdn.humanz.moe/pixiv/?pixivURL="
	Pilot             = "pilot:9000"
	Prediction        = "prediction:9001"
	Indie             = "independent"
	Fe                = "frontend"
	Sys               = "system"

	PixivModule          = "Pixiv Fanart"
	TwitterModule        = "Twitter Fanart"
	TBiliBiliModule      = "BiliBili Fanart"
	LiveBiliBiliModule   = "LiveBiliBili"
	SpaceBiliBiliModule  = "SpaceBiliBili"
	TwitchModule         = "Twitch"
	YoutubeCheckerModule = "Youtube_Checker"
	YoutubeCounterModule = "Youtube_Counter"
	SubscriberModule     = "Subscriber"

	//Crontab
	BiliBiliLive               = "@every 0h7m0s"
	BiliBiliSpace              = "@every 0h13m0s"
	Twitch                     = "@every 0h11m0s"
	TwitchFollowers            = "@every 0h30m0s"
	YoutubeCheckChannel        = "@every 0h5m0s"
	YoutubeCheckUpcomingByTime = "@every 0h1m0s"
	YoutubePrivateSlayer       = "@every 2h31m0s"
	YoutubeSubscriber          = "@every 1h30m0s"
	BiliBiliFollowers          = "@every 0h30m0s"
	TwitterFollowers           = "@every 0h17m0s"
	DonationMsg                = "@every 0h45m0s"
	CheckServerCount           = "@every 0h10m0s"
	CheckPayload               = "@every 1h30m0s"
	PilotGetGroups             = "@every 1h0m0s"

	//Time
	AddUserTTL      = 5 * time.Hour
	FanartSleep     = 5 * time.Second
	GetSubsCountTTL = 20 * time.Minute
	GetUserListTTL  = 30 * time.Minute
	ChannelTagTTL   = 35 * time.Minute
	YtGetStatusTTL  = 20 * time.Minute

	//Unicode
	Ok    = "✅"
	No    = "❎"
	One   = "1️⃣"
	Two   = "2️⃣"
	Three = "3️⃣"
	Four  = "4️⃣"
	Art   = "🎨"
	Live  = "🎥"
	Lewd  = "🔞"

	//GetChannel
	NotLiveOnly = "NotLiveOnly"
	NewUpcoming = "NewUpcoming"
	LewdChannel = "Lewd"
	Default     = "Default"
	Type        = "Type"
	LiveOnly    = "LiveOnly"
	Dynamic     = "Dynamic"
	Region      = "Region"
	LiteMode    = "LiteMode"
	IndieNotif  = "IndieNotif"

	//ChannelType
	ArtType      = 1
	LiveType     = 2
	ArtNLiveType = 3
	LewdType     = 69
	LewdNArtType = 70
	NillType     = 0

	//Live Status
	LiveStatus     = "live"
	PastStatus     = "past"
	UpcomingStatus = "upcoming"
	PrivateStatus  = "private"
	UnknownStatus  = "unknown"

	//FanartState
	TwitterArt  = "twitter"
	PixivArt    = "pixiv"
	BiliBiliArt = "bilibili"

	//LiveState
	YoutubeLive = "youtube"
	BiliLive    = "bilibili"
	SpaceBili   = "spacebili"
	TwitchLive  = "twitch"

	FanartState = "fanart"
	SubsState   = "subs"

	//Metric Query
	Get_Fanart        = "get_fanart"
	Get_Subscriber    = "get_subscriber"
	Get_Live          = "get_live"
	Get_Viewers       = "get_viewers"
	Get_Live_Duration = "get_live_duration"
)

//ConfigFile config file struct for config.toml
type ConfigFile struct {
	Discord        string `toml:"Discord"`
	BiliSess       string `toml:"BiliSess"`
	PixivSession   string `toml:"PixivSess"`
	SauceAPI       string `toml:"SauceAPI"`
	InviteLog      string `toml:"InviteLog"`
	PilotReporting string `toml:"PilotReporting"`
	MultiTOR       string `toml:"Multitor"`
	DonationLink   string `toml:"DonationLink"`
	TopGG          string `toml:"TOPGG"`
	Domain         string `toml:"Domain"`
	PrometheusURL  string `toml:"PrometheusURL"`
	LowResources   bool   `toml:"LowResources"` //Disable update like fanart & set wait every 5 counter
	Metric         bool   `toml:"Metric"`
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
		Twitch   string `toml:"Twitch"`
		General  string `toml:"General"`
		Lewd     string `toml:"Lewd"`
	} `toml:"BotPrefix"`
	Emoji struct {
		Fanart     []string `toml:"Fanart"`
		Livestream []string `toml:"Livestream"`
	} `toml:"Emoji"`
	YtToken []string `toml:"YoutubeToken"`
}

//ReadConfig read from config file
func ReadConfig(path string) (ConfigFile, error) {
	fmt.Println("Reading config file...")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return ConfigFile{}, err
	}

	fmt.Println(string(file))

	_, err = toml.Decode(string(file), &GoSimpConf)
	if err != nil {
		return ConfigFile{}, err
	}

	return GoSimpConf, nil
}

//CheckSQL check if database conn is daijobou
func (Data ConfigFile) CheckSQL() *sql.DB {
	log.Info("Open DB")

	db, err := sql.Open("mysql", Data.SQL.User+":"+Data.SQL.Pass+"@tcp("+Data.SQL.Host+":"+Data.SQL.Port+")/Vtuber?parseTime=true")
	if err != nil {
		log.Panic(err, " Something worng with database,make sure you create Vtuber database first")
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(Data.SQL.MaxOpenConns)
	db.SetMaxIdleConns(Data.SQL.MaxIdleConns)

	//make sure can access database
	_, err = db.Exec(`SELECT NOW()`)
	if err != nil {
		log.Panic(err, " Something worng with database,make sure you create Vtuber database first")
	}
	return db
}

//InitConf initializing config file
func (Data ConfigFile) InitConf() {
	GoSimpConf = Data
	CdnDomain = fmt.Sprintf("https://cdn.%s/", Data.Domain)
	NotFound = fmt.Sprintf("https://cdn.%s/404.jpg", Data.Domain)
	YoutubeIMG = fmt.Sprintf("https://cdn.%s/youtube.png", Data.Domain)
	BiliBiliIMG = fmt.Sprintf("https://cdn.%s/bilibili.png", Data.Domain)
	TwitterIMG = fmt.Sprintf("https://cdn.%s/twitter.png", Data.Domain)
	PixivIMG = fmt.Sprintf("https://cdn.%s/pixiv.png", Data.Domain)
	TwitchIMG = fmt.Sprintf("https://cdn.%s/twitch.png", Data.Domain)
	WorryIMG = fmt.Sprintf("https://cdn.%s/parerunworry.png", Data.Domain)
	GoSimpIMG = fmt.Sprintf("https://cdn.%s/go-simp.png", Data.Domain)
	BSD = fmt.Sprintf("https://cdn.%s/bsd.png", Data.Domain)
	Sleep = fmt.Sprintf("https://cdn.%s/sleep.png", Data.Domain)
	Bonjour = fmt.Sprintf("https://cdn.%s/bonjour.png", Data.Domain)
	Howdy = fmt.Sprintf("https://cdn.%s/howdy.png", Data.Domain)
	Guten = fmt.Sprintf("https://cdn.%s/guten.png", Data.Domain)
	Koni = fmt.Sprintf("https://cdn.%s/koni.png", Data.Domain)
	Selamat = fmt.Sprintf("https://cdn.%s/selamat.jpg", Data.Domain)
	Assalamu = fmt.Sprintf("https://cdn.%s/Assalamu.jpg", Data.Domain)
	Approaching = fmt.Sprintf("https://cdn.%s/approaching.jpg", Data.Domain)
	CommandURL = fmt.Sprintf("https://go-simp.%s/Exec/", Data.Domain)
	GuideURL = fmt.Sprintf("https://go-simp.%s/Guide/", Data.Domain)
	VtubersData = fmt.Sprintf("https://go-simp.%s", Data.Domain)

	if Data.LimitConf.YoutubeLimit >= 15 {
		GoSimpConf.LimitConf.YoutubeLimit = 15
	}

	Scraper = twitterscraper.New()
	Scraper.SetSearchMode(twitterscraper.SearchLatest)
	if Data.MultiTOR != "" {
		err := Scraper.SetProxy("socks5://" + GoSimpConf.MultiTOR)
		if err != nil {
			log.Error(err)
		}
	}

	if Data.BotPrefix.Bilibili == "" {
		log.Fatal("Bilibili Prefix not found")
	}

	if Data.BotPrefix.Fanart == "" {
		log.Fatal("Fanart Prefix not found")
	}

	if Data.BotPrefix.General == "" {
		log.Fatal("General Prefix not found")
	}

	if Data.BotPrefix.Twitch == "" {
		log.Fatal("Twitch Prefix not found")
	}

	if Data.BotPrefix.Youtube == "" {
		log.Fatal("Youtube Prefix not found")
	}
}

func (i *ConfigFile) StartBot() *discordgo.Session {
	tmp, err := discordgo.New("Bot " + i.Discord)
	if err != nil {
		log.Panic(err)
	}
	return tmp
}

func (i *ConfigFile) GetTwitchTkn() *helix.Client {
	TwitchClient, err := helix.NewClient(&helix.Options{
		ClientID:     i.Twitch.ClientID,
		ClientSecret: i.Twitch.ClientSecret,
	})
	if err != nil {
		log.Panic(err)
	}
	return TwitchClient
}
