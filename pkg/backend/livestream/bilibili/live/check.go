package live

import (
	"encoding/json"
	"strconv"

	network "github.com/JustHumanz/Go-simp/tools/network"
	log "github.com/sirupsen/logrus"
)

func GetRoomStatus(RoomID int) (getInfoByRoom, error) {
	var (
		body    []byte
		curlerr error
		tmp     getInfoByRoom
		url     = "https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id=" + strconv.Itoa(RoomID)
	)
	body, curlerr = network.Curl(url, nil)
	if curlerr != nil {
		body, curlerr = network.CoolerCurl(url, nil)
		if curlerr != nil {
			return getInfoByRoom{}, curlerr
		} else {
			log.Info("Successfully use tor")
		}
	}
	err := json.Unmarshal(body, &tmp)
	if err != nil {
		return getInfoByRoom{}, err
	}
	return tmp, nil
}

func (Data getInfoByRoom) CheckScheduleLive() bool {
	if Data.Data.RoomInfo.LiveStatus == 1 {
		return true
	} else {
		return false
	}
}
