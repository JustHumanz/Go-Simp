package live

import (
	"encoding/json"
	"strconv"

	network "github.com/JustHumanz/Go-simp/tools/network"
)

func GetRoomStatus(RoomID int) (getInfoByRoom, error) {
	var (
		tmp getInfoByRoom
	)
	body, curlerr := network.CoolerCurl("https://api.live.bilibili.com/xlive/web-room/v1/index/getInfoByRoom?room_id="+strconv.Itoa(RoomID), nil)
	if curlerr != nil {
		return getInfoByRoom{}, curlerr
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
