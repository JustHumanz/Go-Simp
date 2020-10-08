package database

import (
	"database/sql"

	log "github.com/sirupsen/logrus"
)

//get RoomData from LiveBiliBili
func GetRoomData(MemberID int64, RoomID int) LiveBiliDB {
	rows, err := DB.Query(`SELECT id,RoomID,Status,Title,Thumbnails,Description,ScheduledStart,Published,Viewers FROM LiveBiliBili Where VtuberMember_id=? AND RoomID=? order by ScheduledStart ASC`, MemberID, RoomID)
	BruhMoment(err, "", false)

	defer rows.Close()

	var (
		Data LiveBiliDB
	)
	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.LiveRoomID, &Data.Status, &Data.Title, &Data.Thumbnail, &Data.Description, &Data.ScheduledStart, &Data.PublishedAt, &Data.Online)
		BruhMoment(err, "", false)
	}
	return Data
}

//Update LiveBiliBili Data
func (Data LiveBiliDB) UpdateLiveBili(MemberID int64) {
	_, err := DB.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, Viewers=? Where RoomID=? AND VtuberMember_id=?`, Data.Status, Data.Title, Data.Thumbnail, Data.Description, Data.PublishedAt, Data.ScheduledStart, Data.Online, Data.LiveRoomID, MemberID)
	BruhMoment(err, "", false)
}

//Check new post on TBiliBili
func GetTBiliBili(dynamic_id string) bool {
	var tmp int64
	row := DB.QueryRow("SELECT id FROM Vtuber.TBiliBili where Dynamic_id=?", dynamic_id)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		return true
	} else {
		return false
	}
}

//Get LiveBiliBili by Status (live,past)
func BilGet(GroupID int64, MemberID int64, Status string) []LiveBiliDB {
	funcvar := GetFunctionName(BilGet)
	Debugging(funcvar, "In", GroupID)
	rows, err := DB.Query(`SELECT RoomID,Status,Title,Thumbnails,Description,ScheduledStart,Viewers,VtuberName_EN,VtuberName_JP,BiliBili_Avatar FROM LiveBiliBili Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Status=? Order by ScheduledStart ASC`, GroupID, MemberID, Status)
	BruhMoment(err, "", false)

	defer rows.Close()

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
	Debugging(funcvar, "Out", Data)
	return Data
}

//Get SpaceBiliBili Data
func SpaceGet(GroupID int64, MemberID int64) []SpaceBiliDB {
	funcvar := GetFunctionName(BilGet)
	Debugging(funcvar, "In", GroupID)
	rows, err := DB.Query(`SELECT VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberName_EN,VtuberName_JP,BiliBili_Avatar FROM BiliBili Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) Order by UploadDate DESC limit 3`, GroupID, MemberID)
	BruhMoment(err, "", false)

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
	Debugging(funcvar, "Out", Data)
	return Data
}

//Input data to SpaceBiliBili
func InputSpaceVideo(Data InputBiliBili) {
	funcvar := GetFunctionName(InputSpaceVideo)
	Debugging(funcvar, "In", Data)

	stmt, err := DB.Prepare(`INSERT INTO BiliBili (VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?)`)
	BruhMoment(err, "", false)
	defer stmt.Close()

	res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Title, Data.Thum, Data.Desc, Data.Update, Data.Viewers, Data.Length, Data.MemberID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)
}

//Check New video from SpaceBiliBili
func (Data InputBiliBili) CheckVideo() (bool, int) {
	var tmp int
	row := DB.QueryRow("SELECT id FROM BiliBili WHERE VideoID=? AND VtuberMember_id=?", Data.VideoID, Data.MemberID)
	err := row.Scan(&tmp)
	if err != nil || err == sql.ErrNoRows {
		log.Error(err)
		return true, 0
	} else {
		return false, tmp
	}
}

//Update SpaceBiliBili data
func (Data InputBiliBili) UpdateView(id int) {
	log.WithFields(log.Fields{
		"VideoData ID": Data.VideoID,
		"Viwers":       Data.Viewers,
	}).Info("Update Space.Bilibili")
	_, err := DB.Exec(`Update BiliBili set Viewers=? Where id=?`, Data.Viewers, id)
	BruhMoment(err, "", false)

}

//Input TBiliBili data
func (Data InputTBiliBili) InputTBiliBili(MemberID int64) {
	stmt, err := DB.Prepare(`INSERT INTO TBiliBili (PermanentURL,Author,Likes,Photos,Videos,Text,Dynamic_id,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Error(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.URL, Data.Author, Data.Like, Data.Photos, Data.Videos, Data.Text, Data.Dynamic_id, MemberID)
	if err != nil {
		log.Error(err)
	}

	_, err = res.LastInsertId()
	if err != nil {
		log.Error(err)
	}

}
