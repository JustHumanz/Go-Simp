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
	/*
		Token          string
		YtToken        []string
		EmojiFanart    []string
		EmojiStream    []string
		PFanart        string
		PYoutube       string
		PBilibili      string
		PGeneral       string
		TwitterToken   []string
		ImgurClient    string
		BiliBiliSes    string
		SauceAPI       string
		Logging        string
		DiscordWebHook string
		MultiTOR       string
		KoFiLink       string
	*/
	ModuleList  = []string{"LiveBiliBili", "SpaceBiliBili", "Youtube", "TwitterFanart", "BiliBiliFanart", "YoutubeSubscriber", "BiliBiliFollowers", "TwitterFollowers"}
	BotConf     ConfigFile
	TwitchToken string
)

const (
	TmpDir            = "/tmp/tmp.img"
	NotFound          = "https://cdn.human-z.tech/404.jpg"
	YoutubeIMG        = "https://cdn.human-z.tech/youtube.png"
	BiliBiliIMG       = "https://cdn.human-z.tech/bilibili.png"
	TwitterIMG        = "https://cdn.human-z.tech/twitter.png"
	WorryIMG          = "https://cdn.human-z.tech/parerunworry.png"
	GoSimpIMG         = "https://cdn.human-z.tech/go-simp.png"
	Longcatttt        = "https://cdn.ebaumsworld.com/2020/09/20/013235/86394200/longcat-pic.jpg"
	Dead              = "https://cdn.human-z.tech/dead.jpg"
	BSD               = "https://cdn.human-z.tech/bsd.png"
	Sleep             = "https://cdn.human-z.tech/sleep.png"
	Bonjour           = "https://cdn.human-z.tech/bonjour.png"
	Howdy             = "https://cdn.human-z.tech/howdy.png"
	Guten             = "https://cdn.human-z.tech/guten.png"
	Koni              = "https://cdn.human-z.tech/koni.png"
	Selamat           = "https://cdn.human-z.tech/selamat.jpg"
	Assalamu          = "https://cdn.human-z.tech/Assalamu.jpg"
	Approaching       = "https://cdn.human-z.tech/approaching.jpg"
	CommandURL        = "https://go-simp.human-z.tech/Exec/"
	VtubersData       = "https://go-simp.human-z.tech"
	ChannelPermission = 16

	//Crontab
	TwitterFanart        = "@every 0h3m0s"
	BiliBiliFanart       = "@every 0h6m0s"
	BiliBiliLive         = "@every 0h7m0s"
	BiliBiliSpace        = "@every 0h13m0s"
	YoutubeCheckChannel  = "@every 0h5m0s"
	YoutubePrivateSlayer = "@every 2h31m0s"
	YoutubeSubscriber    = "@every 1h0m0s"
	BiliBiliFollowers    = "@every 0h30m0s"
	TwitterFollowers     = "@every 0h17m0s"
	DonationMsg          = "@every 0h30m0s"
	CheckServerCount     = "@every 0h1m0s"
)

type ConfigFile struct {
	Discord        string   `toml:"Discord"`
	TwitterBearer  []string `toml:"TwitterBearer"`
	ImgurClinet    string   `toml:"ImgurClinet"`
	BiliSess       string   `toml:"BiliSess"`
	SauceAPI       string   `toml:"SauceAPI"`
	DiscordWebHook string   `toml:"DiscordWebHook"`
	MultiTOR       string   `toml:"Multitor"`
	DonationLink   string   `toml:"DonationLink"`
	TopGG          string   `toml:"TOPGG"`
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
