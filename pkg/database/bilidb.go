package database

import (
	"database/sql"
	"strings"

	log "github.com/sirupsen/logrus"
)

//GetRoomData get RoomData from LiveBiliBili
func GetRoomData(MemberID int64, RoomID int) (*LiveBiliDB, error) {
	rows, err := DB.Query(`SELECT id,RoomID,Status,Title,Thumbnails,Description,ScheduledStart,Published,Viewers FROM LiveBiliBili Where VtuberMember_id=? AND RoomID=? order by ScheduledStart ASC`, MemberID, RoomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		Data LiveBiliDB
	)
	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.LiveRoomID, &Data.Status, &Data.Title, &Data.Thumbnail, &Data.Description, &Data.ScheduledStart, &Data.PublishedAt, &Data.Online)
		if err != nil {
			return nil, err
		}
	}
	return &Data, nil
}

//UpdateLiveBili Update LiveBiliBili Data
func (Data *LiveBiliDB) UpdateLiveBili(MemberID int64) {
	_, err := DB.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, Viewers=? Where id=? AND VtuberMember_id=?`, Data.Status, Data.Title, Data.Thumbnail, Data.Description, Data.PublishedAt, Data.ScheduledStart, Data.Online, Data.ID, MemberID)
	if err != nil {
		log.Error(err)
	}
}

//SetRoomToLive force bilibili room status to live
func SetRoomToLive(MemberID int64) {
	_, err := DB.Exec(`Update LiveBiliBili set Status='Live' Where VtuberMember_id=?`, MemberID)
	if err != nil {
		log.Error(err)
	}
}

//GetTBiliBili Check new post on TBiliBili
func GetTBiliBili(DynamicID string) bool {
	var tmp int64
	row := DB.QueryRow("SELECT id FROM Vtuber.TBiliBili where Dynamic_id=?", DynamicID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows || tmp == 0 {
		return true
	} else if err != nil {
		log.Error(err)
	}
	return false
}

//BilGet Get LiveBiliBili by Status (live,past)
func BilGet(GroupID int64, MemberID int64, Status string) []LiveBiliDB {
	var (
		rows *sql.Rows
		err  error
	)

	if GroupID > 0 && Status != "Live" {
		rows, err = DB.Query(`call GetLiveBiliBili(?,?,?,?)`, GroupID, MemberID, Status, 3)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()
	} else {
		rows, err = DB.Query(`call GetLiveBiliBili(?,?,?,?)`, GroupID, MemberID, Status, 2525)
		if err != nil {
			log.Error(err)
		}
		defer rows.Close()
	}

	var (
		Data []LiveBiliDB
		list LiveBiliDB
	)
	for rows.Next() {
		err = rows.Scan(&list.LiveRoomID, &list.Status, &list.Title, &list.Thumbnail, &list.Description, &list.ScheduledStart, &list.Online, &list.EnName, &list.JpName, &list.Avatar)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	return Data
}

//SpaceGet Get SpaceBiliBili Data
func SpaceGet(GroupID int64, MemberID int64) []SpaceBiliDB {
	rows, err := DB.Query(`Call GetSpaceBiliBili(?,?)`, GroupID, MemberID)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	var (
		Data []SpaceBiliDB
		list SpaceBiliDB
	)
	for rows.Next() {
		err = rows.Scan(&list.VideoID, &list.Type, &list.Title, &list.Thumbnail, &list.Description, &list.UploadDate, &list.Viewers, &list.Length, &list.EnName, &list.JpName, &list.Avatar)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	return Data
}

//InputSpaceVideo Input data to SpaceBiliBili
func (Data InputBiliBili) InputSpaceVideo() {
	stmt, err := DB.Prepare(`INSERT INTO BiliBili (VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Error(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Title, Data.Thum, Data.Desc, Data.Update, Data.Viewers, Data.Length, Data.MemberID)
	if err != nil {
		log.Error(err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Error(err)
	}
}

//CheckVideo Check New video from SpaceBiliBili
func (Data InputBiliBili) CheckVideo() (bool, int) {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Vtuber.BiliBili WHERE VideoID=? AND VtuberMember_id=?;", Data.VideoID, Data.MemberID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return true, 0
	} else if err != nil {
		log.Error(err)
		return false, 0
	} else {
		return false, tmp
	}
}

//UpdateView Update SpaceBiliBili data
func (Data InputBiliBili) UpdateView(id int) {
	log.WithFields(log.Fields{
		"VideoData ID": Data.VideoID,
		"Viwers":       Data.Viewers,
	}).Info("Update Space.Bilibili")
	_, err := DB.Exec(`Update BiliBili set Viewers=? Where id=?`, Data.Viewers, id)
	if err != nil {
		log.Error(err)
	}
}

//InputTBiliBili Input TBiliBili data
func (Data TBiliBili) InputTBiliBili() error {
	stmt, err := DB.Prepare(`INSERT INTO TBiliBili (PermanentURL,Author,Likes,Photos,Videos,Text,Dynamic_id,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.URL, Data.Author, Data.Like, strings.Join(Data.Photos, "\n"), Data.Videos, Data.Text, Data.Dynamic_id, Data.Member.ID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}
