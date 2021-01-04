package discordhandler

import database "github.com/JustHumanz/Go-simp/tools/database"

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

//Memberst Vtuber member struct
type Memberst struct {
	ID         int64
	VTName     string
	YtChannel  string
	SpaceID    int
	BiliAvatar string
	YtData     database.YtDbData
	SpaceData  database.SpaceBiliDB
	LiveData   database.LiveBiliDB
	Msg        string
	Msg1       string
	Msg2       string
	Msg3       string
	View       string
	Length     string
}
