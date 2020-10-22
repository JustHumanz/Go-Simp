package database

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

//Get Youtube data from status
func YtGetStatus(Group, Member int64, Status string) []YtDbData {
	funcvar := GetFunctionName(YtGetStatus)
	Debugging(funcvar, "In", fmt.Sprint(Group, Member, Status))
	var (
		rows  *sql.Rows
		err   error
		Data  []YtDbData
		list  YtDbData
		limit int
	)
	if Group != 0 {
		limit = 3
		if Status == "live" {
			limit = 5
		}
	} else {
		limit = 2525
	}
	if Status == "upcoming" {
		rows, err = DB.Query(`SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers FROM Youtube Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Status='upcoming' AND Type='Streaming' AND ScheduledStart !='' Order by ScheduledStart ASC Limit ? `, Group, Member, limit)
		BruhMoment(err, "", false)

	} else if Status == "live" {
		rows, err = DB.Query(`SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers FROM Youtube Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Status='live' limit ?`, Group, Member, limit)
		BruhMoment(err, "", false)

	} else if Status == "private" {
		rows, err = DB.Query(`SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers FROM Youtube Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Status='private'`, Group, Member)
		BruhMoment(err, "", false)
	} else {
		rows, err = DB.Query(`SELECT Youtube.id,VtuberGroupName,Youtube_ID,VtuberName_EN,VtuberName_JP,Youtube_Avatar,VideoID,Title,Thumbnails,Description,ScheduledStart,EndStream,Region,Viewers FROM Youtube Inner join VtuberMember on VtuberMember.id=VtuberMember_id Inner join VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Status='past' AND Type='Streaming' AND EndStream !='' order by EndStream DESC limit 3`, Group, Member)
		BruhMoment(err, "", false)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Group, &list.ChannelID, &list.NameEN, &list.NameJP, &list.YoutubeAvatar, &list.VideoID, &list.Title, &list.Thumb, &list.Desc, &list.Schedul, &list.End, &list.Region, &list.Viewers)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	Debugging(funcvar, "Out", Data)
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
func CheckVideoID(VideoID string) YtDbData {
	funcvar := GetFunctionName(CheckVideoID)
	Debugging(funcvar, "In", VideoID)
	var Data YtDbData
	rows, err := DB.Query(`SELECT id,VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers FROM Vtuber.Youtube Where VideoID=?`, VideoID)
	BruhMoment(err, "", false)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers)
		BruhMoment(err, "", false)
	}
	if Data.ID == 0 {
		return YtDbData{}
	} else {
		return Data
	}
}

//Update youtube data
func (Data YtDbData) UpdateYt(Status string) {
	_, err := DB.Exec(`Update Youtube set Type=?,Status=?,Title=?,Thumbnails=?,Description=?,PublishedAt=?,ScheduledStart=?,EndStream=?,Viewers=?,Length=? where id=?`, Data.Type, Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.ID)
	BruhMoment(err, "", false)
}
