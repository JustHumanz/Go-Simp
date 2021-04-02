package main

import (
	"database/sql"
	"errors"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
)

func GetFanartData(State string, GroupID, MemberID int64) ([]database.DataFanart, error) {
	var (
		Datafanart []database.DataFanart
		DB         = database.DB
		PhotoTmp   sql.NullString
		Video      sql.NullString
	)

	Twitter := func() error {
		rows, err := DB.Query(`SELECT Twitter.* FROM Vtuber.Twitter Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by id desc LIMIT 10`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data database.DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.Member.ID)
			if err != nil {
				return err
			}

			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}
			Data.State = config.TwitterArt
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}
	Tbilibili := func() error {
		rows, err := DB.Query(`SELECT TBiliBili.* FROM Vtuber.TBiliBili Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by id desc LIMIT 10`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data database.DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.Dynamic_id, &Data.Member.ID)
			if err != nil {
				return err
			}
			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}

			Data.State = config.BiliBiliArt
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}

	Pixiv := func() error {
		rows, err := DB.Query(`SELECT Pixiv.* FROM Vtuber.Pixiv Inner Join Vtuber.VtuberMember on VtuberMember.id = Pixiv.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by id desc LIMIT 10`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data database.DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Text, &Data.PixivID, &Data.Member.ID)
			if err != nil {
				return err
			}
			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}

			Data.State = config.PixivArt
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}

	if State == config.PixivArt {
		err := Pixiv()
		if err != nil {
			return nil, err
		}
	} else if State == config.BiliBiliArt {
		err := Tbilibili()
		if err != nil {
			return nil, err
		}
	} else {
		err := Twitter()
		if err != nil {
			return nil, err
		}
	}
	return Datafanart, nil
}
