package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

var (
	// Public variables
	Token   string
	YtToken []string
	//helmp me to pick emoji
	EmojiFanart    []string
	EmojiStream    []string
	PFanart        string
	PYoutube       string
	PBilibili      string
	PGeneral       string
	TwitterToken   string
	ImgurClient    string
	BiliBiliSes    string
	NotFound       string
	SauceAPI       string
	Logging        string
	OwnerDiscordID string
	// Private variables
	config *configStruct
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	TwitterToken = os.Getenv("TWBEARER")
	ImgurClient = os.Getenv("IMGUR")
	BiliBiliSes = os.Getenv("BILISESS")
	NotFound = "https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/404.jpg"
	SauceAPI = os.Getenv("SAUCEAPI")
	Logging = os.Getenv("LOG")
	OwnerDiscordID = os.Getenv("OWNER")
}

type configStruct struct {
	Token     string   `json:"Token"`
	BotPrefix Prefix   `json:"BotPrefix"`
	Emoji     Emoji    `json:"Emoji"`
	YtToken   []string `json:"YtToken"`
}

type Emoji struct {
	Fanart     []string `json:"Fanart"`
	Livestream []string `json:"Livestream"`
}

type Prefix struct {
	Fanart, Youtube, Bilibili, General string
}

//read from config file
func ReadConfig() (*sql.DB, error) {
	fmt.Println("Reading config file...")

	log.Info("Open DB")
	user := os.Getenv("SQLUSER")
	pass := os.Getenv("SQLPASS")
	host := os.Getenv("DBHOST")
	db, err := sql.Open("mysql", ""+user+":"+pass+"@tcp("+host+":3306)/Vtuber?parseTime=true")
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

	file, err := ioutil.ReadFile("config/config.json")

	if err != nil {
		fmt.Println(err.Error())
		return db, err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		return db, err
	}
	Token = config.Token
	YtToken = config.YtToken
	EmojiFanart = config.Emoji.Fanart

	PGeneral = config.BotPrefix.General
	PFanart = config.BotPrefix.Fanart
	PYoutube = config.BotPrefix.Youtube
	PBilibili = config.BotPrefix.Bilibili

	return db, nil
}
