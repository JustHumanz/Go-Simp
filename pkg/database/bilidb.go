package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

//GetRoomData get RoomData from LiveBiliBili
func GetRoomData(MemberID int64, RoomID int) (*LiveStream, string, error) {
	var (
		Data LiveStream
		Key  []string
		ctx  = context.Background()
	)
	Key = append(Key, strconv.Itoa(RoomID), config.Sys)

	Key2 := strings.Join(Key, "-")
	val, err := LiveCache.LRange(ctx, Key2, 0, -1).Result()
	if err != nil {
		return nil, Key2, err
	}

	if len(val) == 0 {
		err := DB.QueryRow("SELECT * FROM LiveBiliBili Where VtuberMember_id=? AND RoomID=?", MemberID, RoomID).Scan(&Data.ID, &Data.Member.BiliRoomID, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.Viewers, &Data.End, &Data.Member.ID)
		if err != nil {
			return nil, Key2, err
		}

		err = LiveCache.LPush(ctx, Key2, Data).Err()
		if err != nil {
			return nil, Key2, err
		}

		err = LiveCache.Expire(ctx, Key2, config.YtGetStatusTTL).Err()
		if err != nil {
			return nil, Key2, err
		}

	} else {
		for _, result := range unique(val) {
			err := json.Unmarshal([]byte(result), &Data)
			if err != nil {
				return nil, Key2, err
			}
		}
	}

	return &Data, Key2, nil
}

//UpdateLiveBili Update LiveBiliBili Data
func (Data *LiveStream) UpdateLiveBili() error {
	_, err := DB.Exec(`Update LiveBiliBili set Status=? , Title=? ,Thumbnails=?, Description=?, Published=?, ScheduledStart=?, EndStream=?,Viewers=? Where id=? AND VtuberMember_id=?`, Data.Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.ID, Data.Member.ID)
	if err != nil {
		return err
	}
	return nil
}

//UpdateBiliToLive force bilibili room status to live
func (Data *LiveStream) UpdateBiliToLive() error {
	_, err := DB.Exec(`Update LiveBiliBili set Status=? Where VtuberMember_id=?`, Data.Status, Data.Member.ID)
	if err != nil {
		return err
	}
	return nil
}

//UpdateBiliToLive force bilibili room status to live
func (Data *LiveStream) SetBiliLive(new bool) *LiveStream {
	Data.IsBiliLive = new
	return Data
}

//BilGet Get LiveBiliBili by Status (live,past)
func BilGet(Payload map[string]interface{}) ([]LiveStream, error) {
	var Group, Member int64
	Status := Payload["Status"].(string)

	if Payload["GroupID"] != nil {
		Group = Payload["GroupID"].(int64)
	} else {
		Member = Payload["MemberID"].(int64)
	}

	var (
		Limit int
	)

	if Group > 0 && Status != "Live" {
		Limit = 3
	} else {
		Limit = 2525
	}

	Query := ""
	if Status == config.LiveStatus {
		Query = "SELECT LiveBiliBili.* FROM Vtuber.LiveBiliBili Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND LiveBiliBili.Status=? Order by ScheduledStart ASC Limit ?"
	} else {
		Query = "SELECT LiveBiliBili.* FROM Vtuber.LiveBiliBili Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND LiveBiliBili.Status=? Order by ScheduledStart DESC Limit ?"
	}

	rows, err := DB.Query(Query, Group, Member, Status, Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		Data []LiveStream
		list LiveStream
	)

	if !rows.Next() {
		return nil, errors.New("not found any schdule")
	}

	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Member.BiliRoomID, &list.Status, &list.Title, &list.Thumb, &list.Desc, &list.Published, &list.Schedul, &list.Viewers, &list.End, &list.Member.ID)
		if err != nil {
			return nil, err
		}
		Data = append(Data, list)
	}
	return Data, nil
}

//SpaceGet Get SpaceBiliBili Data
func SpaceGet(GroupID int64, MemberID int64) ([]LiveStream, error) {
	var (
		Limit int
	)
	if GroupID == 0 {
		Limit = 3
	} else {
		Limit = 5
	}

	rows, err := DB.Query(`SELECT * FROM Vtuber.BiliBili Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) Order by UploadDate DESC limit ?`, GroupID, MemberID, Limit)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	var (
		Data []LiveStream
		list LiveStream
	)

	if !rows.Next() {
		return nil, errors.New("not found any new video")
	}

	for rows.Next() {
		err = rows.Scan(&list.ID, &list.VideoID, &list.Type, &list.Title, &list.Thumb, &list.Desc, &list.Schedul, &list.Viewers, &list.Length, &list.Member.ID)
		if err != nil {
			return nil, err
		}
		Data = append(Data, list)
	}
	return Data, err
}

//InputSpaceVideo Input data to SpaceBiliBili
func (Data LiveStream) InputSpaceVideo() error {
	stmt, err := DB.Prepare(`INSERT INTO BiliBili (VideoID,Type,Title,Thumbnails,Description,UploadDate,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Title, Data.Thumb, Data.Desc, Data.Schedul, Data.Viewers, Data.Length, Data.Member.ID)
	if err != nil {
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

//CheckVideo Check New video from SpaceBiliBili
func (Data LiveStream) CheckVideo() (bool, int) {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Vtuber.BiliBili WHERE VideoID=? AND VtuberMember_id=?", Data.VideoID, Data.Member.ID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		return true, 0
	} else if err != nil {
		return false, 0
	} else {
		return false, tmp
	}
}

//UpdateView Update SpaceBiliBili data
func (Data LiveStream) UpdateSpaceViews(id int) error {
	log.WithFields(log.Fields{
		"VideoData ID": Data.VideoID,
		"Viwers":       Data.Viewers,
	}).Info("Update Space.Bilibili")
	_, err := DB.Exec(`Update BiliBili set Viewers=? Where id=?`, Data.Viewers, id)
	if err != nil {
		return err
	}
	return nil
}
