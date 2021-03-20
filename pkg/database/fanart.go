package database

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

//GetFanart Get Member fanart URL from TBiliBili and Twitter
func GetFanart(GroupID, MemberID int64) (*DataFanart, error) {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
		rows     *sql.Rows
		err      error
	)

	Twitter := func(State string) error {
		rows, err = DB.Query(`Call GetArt(?,?,'twitter')`, GroupID, MemberID)
		if err != nil {
			return err
		} else if err == sql.ErrNoRows {
			return errors.New("Vtuber don't have any fanart")
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		Data.State = State
		return nil
	}
	Tbilibili := func(State string) error {
		rows, err = DB.Query(`Call GetArt(?,?,'westtaiwan')`, GroupID, MemberID)
		if err != nil {
			return err
		} else if err == sql.ErrNoRows {
			return errors.New("Vtuber don't have any fanart")
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.Dynamic_id, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		Data.State = State
		return nil
	}

	if gacha() {
		err := Twitter("Twitter")
		if err != nil {
			return nil, err
		}
	} else {
		err := Tbilibili("TBiliBili")
		if err != nil {
			return nil, err
		}

		if Data.ID == 0 {
			log.Warn("Tbilibili nill")
			err := Twitter("Twitter")
			if err != nil {
				return nil, err
			}
		}
	}

	Data.Videos = Video.String
	Data.Photos = strings.Fields(PhotoTmp.String)
	return &Data, nil

}

//DeleteFanart Delete fanart when get 404 error status
func (Data DataFanart) DeleteFanart(e string) error {
	if notfound, _ := regexp.MatchString("404", e); notfound {
		log.Info("Delete fanart metadata ", Data.PermanentURL)
		if Data.State == "Twitter" {
			stmt, err := DB.Prepare(`DELETE From Twitter WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		} else {
			stmt, err := DB.Prepare(`DELETE From TBiliBili WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		}
	} else {
		return nil
	}
}

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

//CheckMemberFanart Check if `that` was a new fanart
func (FanArt DataFanart) CheckTweetFanArt() (bool, error) {
	var (
		id int
	)
	err := DB.QueryRow(`SELECT id FROM Twitter WHERE TweetID=?`, FanArt.TweetID).Scan(&id)
	if err == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"Name":    FanArt.Member.EnName,
			"Hashtag": FanArt.Member.TwitterHashtags,
		}).Info("New Fanart")

		stmt, err := DB.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
		if err != nil {
			return false, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.TweetID, FanArt.Member.ID)
		if err != nil {
			return false, err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		return false, err
	} else {
		if !config.GoSimpConf.LowResources {
			//update like
			log.WithFields(log.Fields{
				"Name":    FanArt.Member.EnName,
				"Hashtag": FanArt.Member.TwitterHashtags,
				"Likes":   FanArt.Likes,
			}).Info("Update like")
			_, err := DB.Exec(`Update Twitter set Likes=? Where id=? `, FanArt.Likes, id)
			if err != nil {
				return false, err
			}
		}
	}
	return false, nil
}

func (FanArt DataFanart) CheckTBiliBiliFanArt() (bool, error) {
	var tmp int64
	row := DB.QueryRow("SELECT id FROM Vtuber.TBiliBili where Dynamic_id=?", FanArt.Dynamic_id)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"Vtuber": FanArt.Member.EnName,
			"Img":    FanArt.Photos,
		}).Info("New Fanart")
		stmt, err := DB.Prepare(`INSERT INTO TBiliBili (PermanentURL,Author,Likes,Photos,Videos,Text,Dynamic_id,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
		if err != nil {
			return false, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.Dynamic_id, FanArt.Member.ID)
		if err != nil {
			return false, err
		}

		_, err = res.LastInsertId()
		if err != nil {
			return false, err
		}
		return true, nil
	} else if err != nil {
		log.Error(err)
	}
	return false, nil

}
