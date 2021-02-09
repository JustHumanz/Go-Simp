package bilibili

//TBiliBili TopicBiliBili struct
type TBiliBili struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		OffsetTopicBoard string `json:"offset_topic_board"`
		Cards            []struct {
			Desc struct {
				UID         int   `json:"uid"`
				Type        int   `json:"type"`
				Rid         int   `json:"rid"`
				ACL         int   `json:"acl"`
				View        int   `json:"view"`
				Repost      int   `json:"repost"`
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
						VipDueDate    int64  `json:"vipDueDate"`
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
				UIDType       int `json:"uid_type"`
				RecommendInfo struct {
					IsAttention int `json:"is_attention"`
				} `json:"recommend_info"`
				Stype          int    `json:"stype"`
				RType          int    `json:"r_type"`
				InnerID        int    `json:"inner_id"`
				TopicBoard     string `json:"topic_board"`
				TopicBoardDesc string `json:"topic_board_desc"`
				Status         int    `json:"status"`
				DynamicIDStr   string `json:"dynamic_id_str"`
				PreDyIDStr     string `json:"pre_dy_id_str"`
				OrigDyIDStr    string `json:"orig_dy_id_str"`
				RidStr         string `json:"rid_str"`
			} `json:"desc"`
			Card       string `json:"card"`
			ExtendJSON string `json:"extend_json"`
			Display
		} `json:"cards"`
		IsDrawerTopic int    `json:"is_drawer_topic"`
		HasMore       int    `json:"has_more"`
		Offset        string `json:"offset"`
		Gt            int    `json:"_gt_"`
	} `json:"data"`
}

//Display Display struct
type Display []struct { //`json:"display,omitempty"`
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
	Tags []struct {
		TagType   int    `json:"tag_type"`
		Icon      string `json:"icon"`
		Text      string `json:"text"`
		Link      string `json:"link"`
		SubModule string `json:"sub_module"`
	} `json:"tags"`
	UpActButton struct {
		ReportTitle        string `json:"report_title"`
		FounderReportTitle string `json:"founder_report_title"`
		TopTitle           string `json:"top_title"`
		TopConfirmTitle    string `json:"top_confirm_title"`
		TopCancelTitle     string `json:"top_cancel_title"`
	} `json:"up_act_button"`
}

//SubTbili bilibili struct
type SubTbili struct {
	Item struct {
		ID          int           `json:"id"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		Category    string        `json:"category"`
		Role        []interface{} `json:"role"`
		Source      []interface{} `json:"source"`
		Pictures    []struct {
			ImgSrc    string      `json:"img_src"`
			ImgWidth  int         `json:"img_width"`
			ImgHeight int         `json:"img_height"`
			ImgSize   interface{} `json:"img_size"`
		} `json:"pictures"`
		PicturesCount int    `json:"pictures_count"`
		UploadTime    int    `json:"upload_time"`
		AtControl     string `json:"at_control"`
		Reply         int    `json:"reply"`
		Settings      struct {
			CopyForbidden int `json:"copy_forbidden"`
		} `json:"settings"`
		IsFav int `json:"is_fav"`
	} `json:"item"`
	User struct {
		UID     int    `json:"uid"`
		HeadURL string `json:"head_url"`
		Name    string `json:"name"`
		Vip     struct {
			VipType       int    `json:"vipType"`
			VipDueDate    int64  `json:"vipDueDate"`
			DueRemark     string `json:"dueRemark"`
			AccessStatus  int    `json:"accessStatus"`
			VipStatus     int    `json:"vipStatus"`
			VipStatusWarn string `json:"vipStatusWarn"`
			ThemeType     int    `json:"themeType"`
			Label         struct {
				Path string `json:"path"`
			} `json:"label"`
		} `json:"vip"`
	} `json:"user"`
}
