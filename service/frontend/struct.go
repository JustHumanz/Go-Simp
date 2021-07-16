package main

import (
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
)

//DynamicSvr for bilibili author
type DynamicSvr struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Card struct {
			Desc struct {
				UID         int   `json:"uid"`
				Type        int   `json:"type"`
				Rid         int   `json:"rid"`
				ACL         int   `json:"acl"`
				View        int   `json:"view"`
				Repost      int   `json:"repost"`
				Comment     int   `json:"comment"`
				Like        int   `json:"like"`
				IsLiked     int   `json:"is_liked"`
				DynamicID   int64 `json:"dynamic_id"`
				Timestamp   int   `json:"timestamp"`
				PreDyID     int   `json:"pre_dy_id"`
				OrigDyID    int   `json:"orig_dy_id"`
				OrigType    int   `json:"orig_type"`
				UserProfile struct {
					Info struct {
						UID   int    `json:"uid"`
						Uname string `json:"uname"`
						Face  string `json:"face"`
					} `json:"info"`
					Card struct {
						OfficialVerify struct {
							Type int    `json:"type"`
							Desc string `json:"desc"`
						} `json:"official_verify"`
					} `json:"card"`
					Vip struct {
						VipType       int    `json:"vipType"`
						VipDueDate    int    `json:"vipDueDate"`
						DueRemark     string `json:"dueRemark"`
						AccessStatus  int    `json:"accessStatus"`
						VipStatus     int    `json:"vipStatus"`
						VipStatusWarn string `json:"vipStatusWarn"`
						ThemeType     int    `json:"themeType"`
						Label         struct {
							Path string `json:"path"`
						} `json:"label"`
					} `json:"vip"`
					Pendant struct {
						Pid          int    `json:"pid"`
						Name         string `json:"name"`
						Image        string `json:"image"`
						Expire       int    `json:"expire"`
						ImageEnhance string `json:"image_enhance"`
					} `json:"pendant"`
					Rank      string `json:"rank"`
					Sign      string `json:"sign"`
					LevelInfo struct {
						CurrentLevel int    `json:"current_level"`
						CurrentMin   int    `json:"current_min"`
						CurrentExp   int    `json:"current_exp"`
						NextExp      string `json:"next_exp"`
					} `json:"level_info"`
				} `json:"user_profile"`
				UIDType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerID      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIDStr string `json:"dynamic_id_str"`
				PreDyIDStr   string `json:"pre_dy_id_str"`
				OrigDyIDStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
			} `json:"desc"`
			Card       string `json:"card"`
			ExtendJSON string `json:"extend_json"`
			Display    struct {
				TopicInfo struct {
					TopicDetails []struct {
						TopicID    int    `json:"topic_id"`
						TopicName  string `json:"topic_name"`
						IsActivity int    `json:"is_activity"`
						TopicLink  string `json:"topic_link"`
					} `json:"topic_details"`
				} `json:"topic_info"`
				Relation struct {
					Status     int `json:"status"`
					IsFollow   int `json:"is_follow"`
					IsFollowed int `json:"is_followed"`
				} `json:"relation"`
				ShowTip struct {
					DelTip string `json:"del_tip"`
				} `json:"show_tip"`
			} `json:"display"`
		} `json:"card"`
		Result int `json:"result"`
		Gt     int `json:"_gt_"`
	} `json:"data"`
}

type ChannelRegister struct {
	AdminID       string
	State         string
	MessageID     string
	RegionTMP     []string
	AddRegionVal  []string
	DelRegionVal  []string
	Gass          bool
	Emoji         bool
	DisableState  bool
	ChannelState  database.DiscordChannel
	ChannelStates []database.DiscordChannel
	Index         int
}

func (Data *ChannelRegister) SetLiveOnly(new bool) *ChannelRegister {
	Data.ChannelState.LiveOnly = new
	return Data
}

func (Data *ChannelRegister) EmojiTrue() *ChannelRegister {
	Data.Emoji = true
	return Data
}

func (Data *ChannelRegister) SetChannels(a []database.DiscordChannel) *ChannelRegister {
	Data.ChannelStates = a
	return Data
}

func (Data *ChannelRegister) SetNewUpcoming(new bool) *ChannelRegister {
	Data.ChannelState.NewUpcoming = new
	return Data
}

func (Data *ChannelRegister) SetDynamic(new bool) *ChannelRegister {
	Data.ChannelState.Dynamic = new
	return Data
}

func (Data *ChannelRegister) SetLite(new bool) *ChannelRegister {
	Data.ChannelState.LiteMode = new
	return Data
}

func (Data *ChannelRegister) SetIndieNotif(new bool) *ChannelRegister {
	Data.ChannelState.IndieNotif = new
	return Data
}

func (Data *ChannelRegister) SetChannelID(new string) *ChannelRegister {
	Data.ChannelState.ChannelID = new
	return Data
}
func (Data *ChannelRegister) SetChannel(new database.DiscordChannel) *ChannelRegister {
	Data.ChannelState = new
	return Data
}

func (Data *ChannelRegister) UpdateState(new string) *ChannelRegister {
	Data.State = new
	return Data
}

func (Data *ChannelRegister) SetGroup(new database.Group) *ChannelRegister {
	Data.ChannelState.Group = new
	return Data
}

func (Data *ChannelRegister) DisableChannelState(new bool) *ChannelRegister {
	Data.DisableState = new
	return Data
}

func (Data *ChannelRegister) FixRegion(s string) *ChannelRegister {
	list := []string{}
	keys := make(map[string]bool)
	for _, Reg := range Data.RegionTMP {
		if _, value := keys[Reg]; !value {
			keys[Reg] = true
			list = append(list, Reg)
		}
	}
	if s == "add" {
		Data.ChannelState.Region = strings.Join(list, ",")
	} else {
		tmp := []string{}
		for _, v := range list {
			skip := false
			for _, v2 := range Data.DelRegionVal {
				if v2 == v {
					skip = true
					break
				}
			}
			if !skip {
				tmp = append(tmp, v)
			}
		}
		Data.ChannelState.Region = strings.Join(tmp, ",")
	}
	return Data
}

func (Data *ChannelRegister) AddNewRegion(new string) *ChannelRegister {
	Data.AddRegionVal = append(Data.AddRegionVal, new)
	Data.RegionTMP = append(Data.RegionTMP, new)
	return Data
}

func (Data *ChannelRegister) RemoveRegion(new string) *ChannelRegister {
	Data.DelRegionVal = append(Data.DelRegionVal, new)
	return Data
}

func (Data *ChannelRegister) UpdateType(new int) *ChannelRegister {
	Data.ChannelState.TypeTag = new
	return Data
}

func (Data *ChannelRegister) UpdateMessageID(new string) *ChannelRegister {
	Data.MessageID = new
	return Data
}

func (Data *ChannelRegister) Stop() *ChannelRegister {
	Data.Gass = false
	return Data
}

func (Data *ChannelRegister) Start() *ChannelRegister {
	Data.Gass = true
	return Data
}

func (Data *ChannelRegister) BreakPoint(num time.Duration) {
	for i := 0; i < 100; i++ {
		if Data.Gass {
			break
		}
		time.Sleep(num * time.Second)
	}
}

func (Data *ChannelRegister) ChangeLiveStream() *ChannelRegister {
	if Data.ChannelState.TypeTag == config.LiveType || Data.ChannelState.TypeTag == config.ArtNLiveType {
		Data.Stop()
		Data.LiveOnly()
		Data.BreakPoint(1)

		if !Data.ChannelState.LiveOnly {
			Data.Stop()
			Data.NewUpcoming()
			Data.BreakPoint(1)
		}

		Data.Stop()
		Data.Dynamic()
		Data.BreakPoint(1)

		Data.Stop()
		Data.Lite()
		Data.BreakPoint(1)
	}
	return Data
}
