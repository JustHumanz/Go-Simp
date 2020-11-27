package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

func (Data Guild) CheckGuild() int {
	var (
		id int
	)
	err := DB.QueryRow(`SELECT id FROM GuildList WHERE GuildID=? `, Data.ID).Scan(&id)
	if err == sql.ErrNoRows {
		return 0
	} else {
		return id
	}
}

func (Data Guild) UpdateJoin(id int) error {
	stmt, err := DB.Prepare("UPDATE GuildList set JoinDate=? where id=?")
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
	stmt, err := DB.Prepare("INSERT INTO GuildList(GuildName, GuildID,JoinDate) values(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(Data.Name, Data.ID, Data.Join)
	if err != nil {
		return err
	}
	return nil
}
