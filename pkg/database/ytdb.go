package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

/*
//YtGetStatusMap Get Youtube data from status
func YtGetStatus(Payload map[string]interface{}) ([]LiveStream, string, error) {
	var (
		Data          []LiveStream
		list          LiveStream
		Key           []string //= strconv.Itoa(int(Member)) + Status + Region + Uniq
		rows          *sql.Rows
		err           error
		ctx           = context.Background()
		Group, Member int64
		Status        = Payload["Status"].(string)
		Region        string
	)

	if Payload["GroupID"] != nil {
		Group = Payload["GroupID"].(int64)
		Key = append(Key, strconv.Itoa(int(Group)), Payload["GroupName"].(string))
		if Payload["Region"] != nil {
			Region = Payload["Region"].(string)
			if Region != "" {
				Key = append(Key, Region)
			}
		}
	} else {
		Member = Payload["MemberID"].(int64)
		Key = append(Key, strconv.Itoa(int(Member)), Payload["MemberName"].(string))
	}

	Key2 := strings.Join(append(Key, Status, Payload["State"].(string)), "-")

	val, err := LiveCache.LRange(ctx, Key2, 0, -1).Result()
	if err != nil {
		return nil, Key2, err
	}
	if len(val) == 0 {
		limit := func() int {
			if Payload["State"].(string) == config.Sys {
				return 40
			} else {
				if (Group != 0 && Status != "live") || (Member != 0 && Status == "past") {
					return 3
				} else {
					return 2525
				}
			}

		}()
		Query := ""
		if Status == config.PastStatus {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Youtube.Status=? Order by EndStream DESC Limit ?"
		} else if Status == config.UpcomingStatus {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Youtube.Status=? Order by PublishedAt DESC Limit ?"
		} else {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Youtube.Status=? Order by ScheduledStart DESC Limit ?"
		}

		if Region != "" {
			rows, err = DB.Query(`SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=? AND Youtube.Status=? AND Region=? Order by ScheduledStart DESC Limit ?`, Group, Status, Region, limit)
			if err != nil {
				return nil, Key2, err
			}
			defer rows.Close()

		} else {
			rows, err = DB.Query(Query, Group, Member, Status, limit)
			if err != nil {
				return nil, Key2, err
			}
			defer rows.Close()
		}

		for rows.Next() {
			err = rows.Scan(&list.ID, &list.VideoID, &list.Type, &list.Status, &list.Title, &list.Thumb, &list.Desc, &list.Published, &list.Schedul, &list.End, &list.Viewers, &list.Length, &list.Member.ID)
			if err != nil {
				return nil, Key2, err
			}
			UpcominginHours := int(time.Until(list.Schedul).Hours())
			if UpcominginHours > 730 && Status != config.PastStatus || Status == config.Live && UpcominginHours > 168 {
				continue
			}

			Data = append(Data, list)
			err = LiveCache.LPush(ctx, Key2, list).Err()
			if err != nil {
				return nil, Key2, err
			}
		}

		if len(Data) > 0 {
			log.WithFields(log.Fields{
				"State": Payload["State"],
				"Key":   Key2,
			}).Info("Append new cache")
			if Payload["State"].(string) == config.Sys {
				err = LiveCache.Expire(ctx, Key2, config.YtGetStatusTTL).Err()
				if err != nil {
					return nil, Key2, err
				}
			} else {
				err = LiveCache.Expire(ctx, Key2, 5*time.Minute).Err()
				if err != nil {
					return nil, Key2, err
				}
			}
		}

	} else {
		for _, result := range unique(val) {
			err := json.Unmarshal([]byte(result), &list)
			if err != nil {
				return nil, Key2, err
			}

			//Skip empty data
			if list.YtIsEmpty() {
				continue
			}

			Data = append(Data, list)
		}
	}

	return Data, Key2, nil

}
*/
func (Data *Group) GetYtLiveStream(status string, reg string) ([]LiveStream, error) {
	var (
		Query          string
		LiveStreamData []LiveStream
		ctx            = context.Background()
		Live           LiveStream
		rows           *sql.Rows
		err            error
		Key            = func() string {
			if reg != "" {
				return fmt.Sprintf("%d-%s-%s-%s", Data.ID, Data.GroupName, status, reg)
			} else {
				return fmt.Sprintf("%d-%s-%s", Data.ID, Data.GroupName, status)
			}
		}()
	)

	val, err := LiveCache.LRange(ctx, Key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(val) > 0 {
		for _, result := range val {
			err := json.Unmarshal([]byte(result), &Live)
			if err != nil {
				return nil, err
			}

			//Skip empty data
			if Live.YtIsEmpty() {
				continue
			}

			LiveStreamData = append(LiveStreamData, Live)
		}
	} else {
		limit := func() int {
			if status == config.PastStatus {
				return 3
			} else {
				return len(Data.Members)
			}
		}()

		if reg == "" {
			if status == config.PastStatus {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? Order by EndStream DESC Limit ?"
			} else if status == config.UpcomingStatus {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? Order by PublishedAt DESC Limit ?"
			} else {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? Order by ScheduledStart DESC Limit ?"
			}

			rows, err = DB.Query(Query, Data.ID, status, limit)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

		} else {
			if status == config.PastStatus {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? AND Region=? Order by EndStream DESC Limit ?"
			} else if status == config.UpcomingStatus {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? AND Region=? Order by PublishedAt DESC Limit ?"
			} else {
				Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where VtuberGroup.id=?  AND Youtube.Status=? AND Region=? Order by ScheduledStart DESC Limit ?"
			}

			rows, err = DB.Query(Query, Data.ID, status, reg, limit)
			if err != nil {
				return nil, err
			}
			defer rows.Close()
		}

		for rows.Next() {
			err = rows.Scan(&Live.ID, &Live.VideoID, &Live.Type, &Live.Status, &Live.Title, &Live.Thumb, &Live.Desc, &Live.Published, &Live.Schedul, &Live.End, &Live.Viewers, &Live.Length, &Live.Member.ID)
			if err != nil {
				return nil, err
			}

			UpcominginHours := int(time.Until(Live.Schedul).Hours())
			if (UpcominginHours > 730 && status != config.PastStatus) || (status == config.Live && UpcominginHours > 168) {
				continue
			}

			for _, j := range Data.Members {
				if j.ID == Live.Member.ID {
					Live.Member = j
				}
			}

			LiveStreamData = append(LiveStreamData, Live)

			err = LiveCache.LPush(ctx, Key, Live).Err()
			if err != nil {
				return nil, err
			}
		}

		err = LiveCache.Expire(ctx, Key, config.YtGetStatusTTL).Err()
		if err != nil {
			return nil, err
		}

	}

	return LiveStreamData, nil

}

func (Data *Member) GetYtLiveStream(status string) ([]LiveStream, error) {
	var (
		Query          string
		LiveStreamData []LiveStream
		ctx            = context.Background()
		Live           LiveStream
		Key            = fmt.Sprintf("%d-%s-%s", Data.ID, Data.Name, status)
	)

	val, err := LiveCache.LRange(ctx, Key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(val) > 0 {
		for _, result := range val {
			err := json.Unmarshal([]byte(result), &Live)
			if err != nil {
				return nil, err
			}

			//Skip empty data
			if Live.YtIsEmpty() {
				continue
			}

			LiveStreamData = append(LiveStreamData, Live)
		}
	} else {
		if status == config.PastStatus {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id where VtuberMember.id=?  AND Youtube.Status=? Order by EndStream DESC Limit ?"
		} else if status == config.UpcomingStatus {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id where VtuberMember.id=? AND Youtube.Status=? Order by PublishedAt DESC Limit ?"
		} else {
			Query = "SELECT Youtube.* FROM Vtuber.Youtube Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id where VtuberMember.id=?  AND Youtube.Status=? Order by ScheduledStart DESC Limit ?"
		}

		rows, err := DB.Query(Query, Data.ID, status, 5)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&Live.ID, &Live.VideoID, &Live.Type, &Live.Status, &Live.Title, &Live.Thumb, &Live.Desc, &Live.Published, &Live.Schedul, &Live.End, &Live.Viewers, &Live.Length, &Live.Member.ID)
			if err != nil {
				return nil, err
			}
			UpcominginHours := int(time.Until(Live.Schedul).Hours())
			if UpcominginHours > 730 && status != config.PastStatus || status == config.Live && UpcominginHours > 168 {
				continue
			}

			Live.Member = *Data

			LiveStreamData = append(LiveStreamData, Live)

			err = LiveCache.LPush(ctx, Key, Live).Err()
			if err != nil {
				return nil, err
			}
		}

		err = LiveCache.Expire(ctx, Key, config.YtGetStatusTTL).Err()
		if err != nil {
			return nil, err
		}

	}

	return LiveStreamData, nil

}

func (Data *LiveStream) SendToUpcomingCache(isAgency bool) error {
	key := []string{}
	if isAgency {
		key = append(key, Data.Group.GroupName, Data.VideoID)
	} else {
		key = append(key, Data.Member.Name, Data.VideoID)
	}

	err := UpcomingCache.Set(context.Background(), strings.Join(key, "-"), Data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetUpcomingFromCache() (map[string]interface{}, error) {
	ctx := context.Background()
	liveData := make(map[string]interface{})

	iter := UpcomingCache.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		val, err := UpcomingCache.Get(ctx, iter.Val()).Result()
		if err != nil {
			return nil, err
		}

		var Data LiveStream
		err = json.Unmarshal([]byte(val), &Data)
		if err != nil {
			return nil, err
		}

		liveData[iter.Val()] = Data
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return liveData, nil
}

//Input youtube new video
func (Data *LiveStream) InputYt() (int64, error) {
	if !Data.Member.IsMemberNill() {
		stmt, err := DB.Prepare(`INSERT INTO Youtube (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,Length,VtuberMember_id) values(?,?,?,?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.Member.ID)
		if err != nil {
			return 0, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		Data.AddYoutubeToCache(id)

		return id, nil
	} else {
		stmt, err := DB.Prepare(`INSERT INTO GroupVideos (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,Length,LiveBili,VtuberGroup_id) values(?,?,?,?,?,?,?,?,?,?,?,?,?)`)
		if err != nil {
			return 0, err
		}
		defer stmt.Close()

		res, err := stmt.Exec(Data.VideoID, Data.Type, Data.Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.IsBiliLive, Data.Group.ID)
		if err != nil {
			return 0, err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		Data.AddYoutubeToCache(id)

		return id, nil
	}
}

//Check new video or not
func (Member Member) CheckYoutubeVideo(VideoID string) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT id FROM Vtuber.Youtube Where VideoID=? AND VtuberMember_id=?`, VideoID, Member.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID)
		if err != nil {
			return nil, err
		}
	}
	if Data.ID == 0 {
		return nil, errors.New("VideoID not found in database")
	} else {
		Data.AddMember(Member)
		return &Data, nil
	}
}

func (Data *LiveStream) GetYtVideoDetail() {
	rows, err := DB.Query(`SELECT * FROM Vtuber.Youtube Where id=?`, Data.ID)
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers, &Data.Length, &Data.Member.ID)
		if err != nil {
			log.Error(err)
		}
	}
}

//Add new bilibili space video id into cache layer
func (Data LiveStream) AddYoutubeToCache(id int64) {
	key := fmt.Sprintf("cache-layer-%s", Data.VideoID)
	log.WithFields(log.Fields{
		"Key": key,
	}).Info("Add youtube video into cache")

	Data.ID = id
	bytelive, err := json.Marshal(Data)
	if err != nil {
		log.Fatal(err)
	}
	err = LiveCache.Set(context.Background(), key, bytelive, 1*time.Hour).Err()
	if err != nil {
		log.Fatal(err)
	}
}

//Check new video or not
func (Group GroupYtChannel) CheckYoutubeVideo(VideoID string) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT * FROM GroupVideos Where VideoID=? AND VtuberGroup_id=?`, VideoID, Group.GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers, &Data.Length, &Data.IsBiliLive, &Data.Group.ID)
		if err != nil {
			return nil, err
		}
	}

	if Data.VideoID == "" {
		return nil, errors.New("VideoID not found in database")
	} else {
		return &Data, nil
	}
}

//Update youtube data
func (Data *LiveStream) UpdateYt(Status string) error {
	_, err := DB.Exec(`Update Youtube set Type=?,Status=?,Title=?,Thumbnails=?,Description=?,PublishedAt=?,ScheduledStart=?,EndStream=?,Viewers=?,Length=? where id=?`, Data.Type, Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.ID)
	if err != nil {
		return err
	}
	return nil
}

//Update youtube data
func (Data *LiveStream) UpdateGroupYt(Status string) error {
	_, err := DB.Exec(`Update GroupVideos set Type=?,Status=?,Title=?,Thumbnails=?,Description=?,PublishedAt=?,ScheduledStart=?,EndStream=?,Viewers=?,Length=? where id=?`, Data.Type, Status, Data.Title, Data.Thumb, Data.Desc, Data.Published, Data.Schedul, Data.End, Data.Viewers, Data.Length, Data.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckReminder(Num int) bool {
	var count int
	err := DB.QueryRow(`SELECT id FROM User where Reminder=?`, Num).Scan(&count)
	if err != nil {
		log.Error(err)
	} else if err == sql.ErrNoRows {
		return false
	}
	return true
}
