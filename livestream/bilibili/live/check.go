package bilibili

import (
	"encoding/json"
	"strconv"

	engine "github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

func GetRoomStatus(RoomID int) getInfoByRoom {
	var (
		body    []byte
		curlerr error
		tmp     getInfoByRoom
		url     = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id=" + strconv.Itoa(RoomID)
	)
	body, curlerr = engine.Curl(url, nil)
	if curlerr != nil {
		body, curlerr = engine.CoolerCurl(url)
		if curlerr != nil {
			log.Error(curlerr)
		} else {
			log.Info("Successfully use tor")
		}
	}
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
		log.Info("Trying use tor")
		body, curlerr = engine.CoolerCurl(urls)
		if curlerr != nil {
			log.Error(curlerr)
		} else {
			log.Info("Successfully use tor")
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
