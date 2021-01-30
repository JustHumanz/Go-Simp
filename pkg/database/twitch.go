package database

func GetTwitch(MemberID int64) (*TwitchDB, error) {
	var Data TwitchDB
	rows, err := DB.Query(`SELECT * FROM Vtuber.Twitch Where VtuberMember_id=?`, MemberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.Game, &Data.Status, &Data.Title, &Data.Thumbnails, &Data.ScheduledStart, &Data.Viewers, &MemberID)
		if err != nil {
			return nil, err
		}
	}
	return &Data, nil
}

func (Data *TwitchDB) UpdateTwitch(MemberID int64) error {
	_, err := DB.Exec(`Update Twitch set Game=?,Status=?,Thumbnails=?,ScheduledStart=?,Viewers=? Where id=? AND VtuberMember_id=?`, Data.ID, MemberID)
	if err != nil {
		return err
	}
	return nil
}
