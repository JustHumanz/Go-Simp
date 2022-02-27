package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/go-redis/redis/v8"
)

//Get data of member twitch
func GetTwitch(MemberID int64) (*LiveStream, error) {
	var Data LiveStream
	err := DB.QueryRow("SELECT * FROM Vtuber.Twitch Where VtuberMember_id=?", MemberID).Scan(&Data.ID, &Data.Game, &Data.Status, &Data.Title, &Data.Thumb, &Data.Schedul, &Data.End, &Data.Viewers, &MemberID)
	if err != nil {
		return nil, err
	}

	return &Data, nil
}

func (Data *LiveStream) UpdateTwitch() error {
	_, err := DB.Exec(`Update Twitch set Game=?,Status=?,Thumbnails=?,ScheduledStart=?,EndStream=?,Viewers=? Where id=? AND VtuberMember_id=?`, Data.Game, Data.Status, Data.Thumb, Data.Schedul, Data.End, Data.Viewers, Data.ID, Data.Member.ID)
	if err != nil {
		return err
	}
	return nil
}

func (Data *Member) GetTwitchLiveStream(status string) (LiveStream, error) {
	var (
		Query string
		ctx   = context.Background()
		Live  LiveStream
		Key   = fmt.Sprintf("%d-%s-%s-%s", Data.ID, Data.Name, status, Data.TwitchName)
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
			Query = "SELECT Twitch.* FROM Vtuber.Twitch Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Where  VtuberMember.id=? AND Twitch.Status=? Order by ScheduledStart ASC Limit ?"
		} else {
			Query = "SELECT Twitch.* FROM Vtuber.Twitch Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Where  VtuberMember.id=? AND Twitch.Status=? Order by ScheduledStart DESC Limit ?"
		}

		rows, err := DB.Query(Query, Data.ID, status, Limit)
		if err != nil {
			return LiveStream{}, err
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&Live.ID, &Live.Game, &Live.Status, &Live.Title, &Live.Thumb, &Live.Schedul, &Live.End, &Live.Viewers, &Live.Member.ID)
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

	return Live, nil
}

func (Group Group) GetTwitchLiveStream(status string) ([]LiveStream, error) {
	var LiveData []LiveStream
	for _, Member := range Group.Members {
		Live, err := Member.GetBlLiveStream(status)
		if err != nil {
			return nil, err
		}

		if Live.ID == 0 {
			continue
		}

		LiveData = append(LiveData, Live)
	}
	return LiveData, nil
}
