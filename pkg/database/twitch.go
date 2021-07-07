package database

import "errors"

func GetTwitch(MemberID int64) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT * FROM Vtuber.Twitch Where VtuberMember_id=?`, MemberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, errors.New("not found any schdule")
	}

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.Game, &Data.Status, &Data.Title, &Data.Thumb, &Data.Schedul, &Data.End, &Data.Viewers, &MemberID)
		if err != nil {
			return nil, err
		}
	}

	if Data.ID == 0 {
		return nil, errors.New("not found any schdule")
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

//TwitchGet Get LiveBiliBili by Status (live,past)
func TwitchGet(GroupID int64, MemberID int64, Status string) ([]LiveStream, error) {
	var (
		Limit int
	)

	if GroupID > 0 && Status != "Live" {
		Limit = 3
	} else {
		Limit = 2525
	}

	rows, err := DB.Query(`SELECT Twitch.* FROM Vtuber.Twitch Inner join Vtuber.VtuberMember on VtuberMember.id=VtuberMember_id Inner join Vtuber.VtuberGroup on VtuberGroup.id = VtuberGroup_id Where (VtuberGroup.id=? or VtuberMember.id=?) AND Twitch.Status=? Order by ScheduledStart ASC Limit ?`, GroupID, MemberID, Status, Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		Data []LiveStream
		list LiveStream
	)
	for rows.Next() {
		err = rows.Scan(&list.ID, &list.Game, &list.Status, &list.Title, &list.Thumb, &list.Schedul, &list.End, &list.Viewers, &list.Member.ID)
		if err != nil {
			return nil, err
		}
		Data = append(Data, list)
	}
	return Data, nil
}
