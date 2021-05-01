package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	twitterscraper "github.com/n0madic/twitter-scraper"
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
	Waiting           = 5
	Pilot             = "pilot"
	Indie             = "independent"
	Fe                = "frontend"
	Sys               = "system"

	//Crontab
	TwitterFanart              = "@every 0h3m0s"
	BiliBiliFanart             = "@every 0h6m0s"
	PixivFanart                = "@every 0h15m0s"
	PixivFanartLewd            = "@every 1h0m0s"
	BiliBiliLive               = "@every 0h7m0s"
	BiliBiliSpace              = "@every 0h13m0s"
	Twitch                     = "@every 0h11m0s"
	YoutubeCheckChannel        = "@every 0h5m0s"
	YoutubeCheckUpcomingByTime = "@every 0h1m0s"
	YoutubePrivateSlayer       = "@every 2h31m0s"
	YoutubeSubscriber          = "@every 1h0m0s"
	BiliBiliFollowers          = "@every 0h30m0s"
	TwitterFollowers           = "@every 0h17m0s"
	DonationMsg                = "@every 0h45m0s"
	CheckServerCount           = "@every 0h10m0s"
	CheckPayload               = "@every 2h0m0s"
	PilotGetGroups             = "@every 1h30m0s"

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
	Get_Fanart       = "get_fanart"
	Get_Subscriber   = "get_subscriber"
	Get_Live         = "get_live"
	Get_Viewers      = "get_viewers"
	Get_Flying_Hours = "get_flying_hours"
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

/*
//GetTwitchAccessToken get twitch access token
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
*/

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
		err := Scraper.SetProxy(GoSimpConf.MultiTOR)
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
