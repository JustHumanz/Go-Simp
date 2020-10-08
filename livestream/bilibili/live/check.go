package bilibili

import (
	"encoding/json"
	"strconv"

	engine "github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func GetRoomStatus(RoomID int) getInfoByRoom {
	body, err := engine.Curl("https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id="+strconv.Itoa(RoomID), nil)
	if err != nil {
		log.Error(err, string(body))
	}
	var tmp getInfoByRoom
	jsonErr := json.Unmarshal(body, &tmp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return tmp
}

func (Data RoomID2) GetRoomStatus() getInfoByRoom {
	return GetRoomStatus(Data.Data.Roomid)
}

func GetRoom2(SpaceID int) RoomID2 {
	var (
		body    []byte
		curlerr error
		urls    = "https://api.live.bilibili.com/room/v1/Room/getRoomInfoOld?mid=" + strconv.Itoa(SpaceID)
	)
	body, curlerr = engine.Curl(urls, nil)
	if curlerr != nil {
		log.Error(curlerr, string(body))

		log.Info("Trying use tor")
		body, curlerr = engine.CoolerCurl(urls)
		if curlerr != nil {
			log.Error(curlerr)
		}
	}

	var tmp RoomID2
	jsonErr := json.Unmarshal(body, &tmp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return tmp
}

func (Data NewSchedule) CheckNewScheduleStatus() bool {
	if Data.DurationUp > Data.DurationPast {
		return true
	} else {
		return false
	}
}

func (Data getInfoByRoom) CheckScheduleLive() bool {
	if Data.Data.RoomInfo.LiveStatus == 1 {
		return true
	} else {
		return false
	}
}
