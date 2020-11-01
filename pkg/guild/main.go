package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var (
	H3llcome   = []string{config.Bonjour, config.Howdy, config.Guten, config.Koni, config.Selamat, config.Assalamu, config.Approaching}
	GuildList  []string
	PathLiteDB = "./guild.db"
	BotID      *discordgo.User
)

func main() {
	_, err := config.ReadConfig("../../config.toml")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Bot, _ := discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	BotID, err = Bot.User("@me")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for _, GuildID := range Bot.State.Guilds {
		GuildList = append(GuildList, GuildID.ID)
	}

	Bot.AddHandler(GuildJoin)
	log.Info("Guild handler ready.......")

	chain := make(chan os.Signal, 0)
	signal.Notify(chain, os.Interrupt)
	go func() {
		for sig := range chain {
			log.Warn("captured ", sig, ", stopping profiler and exiting..")
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	<-make(chan struct{})
	return
}

type Guild struct {
	ID     string
	Name   string
	Join   time.Time
	Dbconn *sql.DB
}

func CreateLite(path string) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
	}
	file.Close()
}

func OpenLiteDB(path string) *sql.DB {
	dblite, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Error(err)
	}
	dblite.Exec(`CREATE TABLE IF NOT EXISTS "GuildList" (
		"id"	INTEGER NOT NULL,
		"GuildID"	TEXT,
		"GuildName"	TEXT,
		"JoinDate"	TEXT,
		PRIMARY KEY("id" AUTOINCREMENT)
	);`)
	return dblite
}

func (Data Guild) CheckGuild() int {
	var (
		id int
	)
	err := Data.Dbconn.QueryRow(`SELECT id FROM GuildList WHERE GuildID=? `, Data.ID).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	} else {
		return id
	}
}

func (Data Guild) UpdateJoin(id int) error {
	stmt, err := Data.Dbconn.Prepare("UPDATE GuildList set JoinDate=? where id=?")
	if err != nil {
		log.Error(err)
	}
	_, err = stmt.Exec(Data.Join, id)
	if err != nil {
		log.Error(err)
	}
	return nil
}

func (Data Guild) InputGuild() error {
	stmt, err := Data.Dbconn.Prepare("INSERT INTO GuildList(GuildName, GuildID,JoinDate) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Data.Name, Data.ID, Data.Join)
	if err != nil {
		return err
	}
	return nil
}

func KillSqlConn(sql *sql.DB) {
	sql.Close()
}
