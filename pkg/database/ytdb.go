package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	log "github.com/sirupsen/logrus"
)

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

		if !rows.Next() {
			log.WithFields(log.Fields{
				"State": Payload["State"],
				"Key":   Key2,
			}).Warn("not found any schdule,set value to empty")
			err = LiveCache.LPush(ctx, Key2, LiveStream{}).Err()
			if err != nil {
				return nil, Key2, err
			}

		} else {
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
		}

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

		return id, nil
	}
}

//YtIsEmpty?
func (Data LiveStream) YtIsEmpty() bool {
	return Data.VideoID != ""
}

//Check new video or not
func (Member Member) CheckYoutubeVideo(VideoID string) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT * FROM Vtuber.Youtube Where VideoID=? AND VtuberMember_id=?`, VideoID, Member.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.VideoID, &Data.Type, &Data.Status, &Data.Title, &Data.Thumb, &Data.Desc, &Data.Published, &Data.Schedul, &Data.End, &Data.Viewers, &Data.Length, &Data.Member.ID)
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

//Check new video or not
func (Group GroupYtChannel) CheckYoutubeVideo(VideoID string) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT * FROM Vtuber.GroupVideos Where VideoID=? AND VtuberGroup_id=?`, VideoID, Group.ID)
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
	if Data.ID == 0 {
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
