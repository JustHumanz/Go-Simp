package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

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

func (Data *Member) GetBlLiveStream(status string) (LiveStream, error) {
	var (
		Query string
		ctx   = context.Background()
		Live  LiveStream
		Key   = fmt.Sprintf("%d-%s-%s-%d", Data.ID, Data.Name, status, Data.BiliBiliID)
	)

	val, err := LiveCache.Get(ctx, Key).Result()
	if err == redis.Nil {
		Limit := func() int {
			if status == config.PastStatus {
				return 3
			} else {
				return 4
			}
		}()

		if status == config.LiveStatus {
			Query = "SELECT LiveBiliBili.* FROM Vtuber.LiveBiliBili Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Where  VtuberMember.id=? AND LiveBiliBili.Status=? Order by ScheduledStart ASC Limit ?"
		} else {
			Query = "SELECT LiveBiliBili.* FROM Vtuber.LiveBiliBili Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Where  VtuberMember.id=? AND LiveBiliBili.Status=? Order by ScheduledStart DESC Limit ?"
		}

		rows, err := DB.Query(Query, Data.ID, status, Limit)
		if err != nil {
			return LiveStream{}, err
		}
		defer rows.Close()

		for rows.Next() {
			var list LiveStream
			err = rows.Scan(&list.ID, &list.Member.BiliBiliRoomID, &list.Status, &list.Title, &list.Thumb, &list.Desc, &list.Published, &list.Schedul, &list.Viewers, &list.End, &list.Member.ID)
			if err != nil {
				return LiveStream{}, err
			}
		}

		err = LiveCache.Set(ctx, Key, Live, 0).Err()
		if err != nil {
			return LiveStream{}, err
		}

		err = LiveCache.Expire(ctx, Key, config.YtGetStatusTTL).Err()
		if err != nil {
			return LiveStream{}, err
		}

	} else if err != nil {
		return LiveStream{}, err
	} else {
		err := json.Unmarshal([]byte(val), &Live)
		if err != nil {
			return LiveStream{}, err
		}
	}

	Live.AddMember(*Data)

	return Live, nil
}

//BilGet Get LiveBiliBili by Status (live,past)
func (Data *Group) GetBlLiveStream(status string) ([]LiveStream, error) {
	var LiveStreamData []LiveStream
	for _, Member := range Data.Members {
		Live, err := Member.GetBlLiveStream(status)
		if err != nil {
			return nil, err
		}

		if Live.ID == 0 {
			continue
		}

		LiveStreamData = append(LiveStreamData, Live)
	}

	return LiveStreamData, nil
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

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	Data.AddBiliBiliSpaceToCache(id)

	return nil
}

//CheckVideo Check New video from SpaceBiliBili
func (Data LiveStream) SpaceCheckVideo() error {
	var tmp int
	row := DB.QueryRow("SELECT id FROM Vtuber.BiliBili WHERE VideoID=? AND VtuberMember_id=?", Data.VideoID, Data.Member.ID)
	err := row.Scan(&tmp)
	if err == sql.ErrNoRows {
		log.WithFields(log.Fields{
			"Agency": Data.Group.GroupName,
			"Vtuber": Data.Member.Name,
		}).Info("New video uploaded")
		err := Data.InputSpaceVideo()
		if err != nil {
			return err
		}

	} else if err != nil {
		return err
	}
	return nil
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

//
func (Data LiveStream) AddBiliBiliSpaceToCache(id int64) {
	key := fmt.Sprintf("cache-layer-%s", Data.VideoID)
	log.WithFields(log.Fields{
		"Key": key,
	}).Info("Add fanart into cache")

	Data.ID = id
	bytelive, err := json.Marshal(Data)
	if err != nil {
		log.Fatal(err)
	}
	err = LiveCache.Set(context.Background(), key, bytelive, 10*time.Minute).Err()
	if err != nil {
		log.Fatal(err)
	}
}
