package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

//CheckGuild Check guild in database
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

//UpdateJoin update guild join
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

//InputGuild add new guild
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

//GetGuildsCount get count of guild
func GetGuildsCount() int {
	var count int
	err := DB.QueryRow(`SELECT Count(*) FROM Vtuber.GuildList`).Scan(&count)
	if err != nil {
		log.Error(err)
	}
	return count
}

//GetMemberCount get count of member
func GetMemberCount() int {
	var count int
	err := DB.QueryRow(`SELECT Count(*) from (Select COUNT(id) FROM Vtuber.User Group by DiscordID) t`).Scan(&count)
	if err != nil {
		log.Error(err)
	}
	return count
}
