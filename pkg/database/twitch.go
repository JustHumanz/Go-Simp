package database

func GetTwitch(MemberID int64) (*LiveStream, error) {
	var Data LiveStream
	rows, err := DB.Query(`SELECT * FROM Vtuber.Twitch Where VtuberMember_id=?`, MemberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.Game, &Data.Status, &Data.Title, &Data.Thumb, &Data.Schedul, &Data.Viewers, &MemberID)
		if err != nil {
			return nil, err
		}
	}
	return &Data, nil
}

func (Data *LiveStream) UpdateTwitch() error {
	_, err := DB.Exec(`Update Twitch set Game=?,Status=?,Thumbnails=?,ScheduledStart=?,Viewers=? Where id=? AND VtuberMember_id=?`, Data.Game, Data.Status, Data.Thumb, Data.Schedul, Data.Viewers, Data.ID, Data.Member.ID)
	if err != nil {
		return err
	}
	return nil
}
