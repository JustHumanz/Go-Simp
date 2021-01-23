package main

import (
	"database/sql"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

func GetFanartData(State string, GroupID, MemberID int64) []database.DataFanart {
	var (
		Datafanart []database.DataFanart
		db         = database.DB
		PhotoTmp   sql.NullString
	)
	if State == "Twitter" {
		rows, err := db.Query(`SELECT Twitter.id,VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Videos,Text FROM Vtuber.Twitter Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?) Limit 15`, GroupID, MemberID)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()
		for rows.Next() {
			var Data database.DataFanart
			err = rows.Scan(&Data.ID, &Data.EnName, &Data.JpName, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Videos, &Data.Text)
			if err != nil {
				log.Error(err)
			}
			Data.Photos = strings.Fields(PhotoTmp.String)
			Datafanart = append(Datafanart, Data)
		}
	} else {
		rows, err := db.Query(`SELECT TBiliBili.id,VtuberName_EN,VtuberName_JP,PermanentURL,Author,Photos,Text FROM Vtuber.TBiliBili  Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?) Limit 15`, GroupID, MemberID)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()
		for rows.Next() {
			var Data database.DataFanart
			err = rows.Scan(&Data.ID, &Data.EnName, &Data.JpName, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Text)
			if err != nil {
				log.Error(err)
			}
			Data.Photos = strings.Fields(PhotoTmp.String)
			Datafanart = append(Datafanart, Data)
		}
	}
	return Datafanart
}
