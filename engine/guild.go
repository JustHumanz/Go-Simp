package engine

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Guild struct {
	ID   string
	Name string
	Join time.Time
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

func (Data Guild) CheckGuild(dblite *sql.DB) string {
	var (
		id int
	)
	err := dblite.QueryRow(`SELECT id FROM GuildList WHERE GuildID=? `, Data.ID).Scan(&id)
	if err == sql.ErrNoRows {
		return "New"
	} else {
		if dblite.QueryRow(`SELECT id FROM GuildList WHERE JoinDate=? `, Data.Join).Scan(&id) == sql.ErrNoRows {
			stmt, err := dblite.Prepare("UPDATE GuildList set JoinDate=? where id=?")
			if err != nil {
				log.Error(err)
			}
			_, err = stmt.Exec(Data.Join, id)
			if err != nil {
				log.Error(err)
			}
			return "Rejoin"
		} else {
			return ""
		}
	}
}

func (Data Guild) InputGuild(dblite *sql.DB) error {
	stmt, err := dblite.Prepare("INSERT INTO GuildList(GuildName, GuildID,JoinDate) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Data.Name, Data.ID, Data.Join)
	if err != nil {
		return err
	}
	return nil
}

func KillSqlite(sql *sql.DB) {
	sql.Close()
}
