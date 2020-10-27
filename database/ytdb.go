package database

import (
	"database/sql"
	"time"

	log "github.com/sirupsen/logrus"
)

//Get Youtube data from status
func YtGetStatus(Group, Member int64, Status, Region string) []YtDbData {
	var (
		Data  []YtDbData
		list  YtDbData
		rows  *sql.Rows
		limit int
		err   error
	)
	if (Group != 0 && Status != "live") || (Member != 0 && Status == "past") {
		limit = 3
	} else {
		limit = 2525
	}
	if Region == "" {
		rows, err = DB.Query(`call GetYt(?,?,?,?)`, Member, Group, limit, Status)
	} else {
		rows, err = DB.Query(`call GetYtByReg(?,?,?)`, Group, Status, Region)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Group, &list.ChannelID, &list.NameEN, &list.NameJP, &list.YoutubeAvatar, &list.VideoID, &list.Title, &list.Thumb, &list.Desc, &list.Schedul, &list.End, &list.Region, &list.Viewers)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	return Data

}

//Input youtube new video
func (Data YtDbData) InputYt(MemberID int64) {
	stmt, err := DB.Prepare(`INSERT INTO Youtube (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?,?,?,?)`)
	BruhMoment(err, "", false)
	defer stmt.Close()

	res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, MemberID)
	BruhMoment(err, "", false)

	_, err = res.LastInsertId()
	BruhMoment(err, "", false)

}

//Check new video or not
func (Member Name) CheckYtVideo(VideoID string) *YtDbData {
	var Data YtDbData
	rows, err := DB.Query(`SELECT id,VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers FROM Vtuber.Youtube Where VideoID=? AND VtuberMember_id=?`, VideoID, Member.ID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers)
		BruhMoment(err, "", false)
	}
	if Data.ID == 0 {
		return nil
	} else {
		return &Data
	}
}

func (Data *YtDbData) UpView(new string) *YtDbData {
	Data.Viewers = new
	return Data
}

func (Data *YtDbData) UpEnd(new time.Time) *YtDbData {
	Data.End = new
	return Data
}

func (Data *YtDbData) UpLength(new string) *YtDbData {
	Data.Length = new
	return Data
}

//Update youtube data
func (Data YtDbData) UpdateYt(Status string) {
	_, err := DB.Exec(`Update Youtube set Type=?,Status=?,Title=?,Thumbnails=?,Description=?,PublishedAt=?,ScheduledStart=?,EndStream=?,Viewers=?,Length=? where id=?`, Data.Type, Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.ID)
	BruhMoment(err, "", false)
}
