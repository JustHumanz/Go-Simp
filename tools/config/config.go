package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

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
	BotConf ConfigFile
	TmpDir  = "/var/tmp.img"
)

const (
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
	LimitConf      struct {
		TwitterFanart int `toml:"TwitterLimit"`
		SpaceBiliBili int `toml:"SpaceBiliBili"`
		YoutubeLimit  int `toml:"YoutubeLimit"`
	} `toml:"Limit"`
	SQL struct {
		User string `toml:"User"`
		Pass string `toml:"Pass"`
		Host string `toml:"Host"`
	} `toml:"Sql"`
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
	/*
		TwitterToken = config.TwitterBearer
		ImgurClient = config.ImgurClinet
		BiliBiliSes = config.BiliSess
		SauceAPI = config.SauceAPI

		DiscordWebHook = config.DiscordWebHook
		MultiTOR = config.MultiTOR
		KoFiLink = config.DonationLink

		Token = config.Discord
		YtToken = config.YtToken
		EmojiFanart = config.Emoji.Fanart

		PGeneral = config.BotPrefix.General
		PFanart = config.BotPrefix.Fanart
		PYoutube = config.BotPrefix.Youtube
		PBilibili = config.BotPrefix.Bilibili
	*/

	return BotConf, nil
}

func (Data ConfigFile) CheckSQL() *sql.DB {
	log.Info("Open DB")

	db, err := sql.Open("mysql", Data.SQL.User+":"+Data.SQL.Pass+"@tcp("+Data.SQL.Host+":3306)/Vtuber?parseTime=true")
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}

	//make sure can access database
	_, err = db.Exec(`SELECT NOW()`)
	if err != nil {
		log.Error(err, " Something worng with database,make sure you create Vtuber database first")
		os.Exit(1)
	}
	return db
}
