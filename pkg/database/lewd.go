package database

import (
	"database/sql"
	"strings"

	log "github.com/sirupsen/logrus"
)

func IsLewdNew(s, ID string) bool {
	var id int
	if s == "Pixiv" {
		err := DB.QueryRow(`SELECT id FROM Lewd WHERE PixivID=?`, ID).Scan(&id)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"State": "Pixiv",
			}).Info("New Lewd Fanart")
			return true
		}
	} else {
		Twitter := strings.Split(ID, "/")
		err := DB.QueryRow(`SELECT id FROM Lewd WHERE TweetID=?`, Twitter[len(Twitter)-1]).Scan(&id)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"State": "Twitter",
			}).Info("New Lewd Fanart")
			return true
		}
	}
	return false
}
func AddLewd(Data DataFanart) error {
	stmt, err := DB.Prepare(`INSERT INTO Lewd (PermanentURL,Author,Photos,Videos,Text,TweetID,PixivID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.PermanentURL, Data.Author, strings.Join(Data.Photos, "\n"), Data.Videos, Data.Text, Data.TweetID, Data.PixivID, Data.Member.ID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
