package main

import (
	"time"
)

type YtData struct {
	Items []struct {
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"Standard"`
			} `json:"thumbnails"`
			ChannelTitle string `json:"channelTitle"`
			VideoStatus  string `json:"liveBroadcastContent"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount string `json:"viewCount"`
		} `json:"statistics"`
		LiveDetails struct {
			StartTime time.Time `json:"scheduledStartTime"`
			Viewers   string    `json:"concurrentViewers"`
			EndTime   time.Time `json:"actualEndTime"`
		} `json:"liveStreamingDetails"`
	} `json:"items"`
}

type Vtuber struct {
	VtuberData Data `json:"Data"`
}
type Twitter struct {
	TwitterFanart   string `json:"Twitter_Fanart"`
	TwitterLewd     string `json:"Twitter_Lewd"`
	TwitterUsername string `json:"Twitter_Username"`
}
type Youtube struct {
	YtID string `json:"Yt_ID"`
}
type BiliBili struct {
	BiliBiliFanart string `json:"BiliBili_Fanart"`
	BiliBiliID     int    `json:"BiliBili_ID"`
	BiliRoomID     int    `json:"BiliRoom_ID"`
}
type Twitch struct {
	TwitchUsername string `json:"Twitch_Username"`
}
type Members struct {
	Name     string   `json:"Name"`
	ENName   string   `json:"EN_Name"`
	JPName   string   `json:"JP_Name"`
	Twitter  Twitter  `json:"Twitter"`
	Youtube  Youtube  `json:"Youtube"`
	BiliBili BiliBili `json:"BiliBili"`
	Twitch   Twitch   `json:"Twitch"`
	Region   string   `json:"Region"`
	Fanbase  string   `json:"Fanbase"`
	Status   string   `json:"Status"`
}
type Independent struct {
	Members []Members `json:"Members"`
}
type Group struct {
	GroupName    string       `json:"GroupName"`
	GroupIcon    string       `json:"GroupIcon"`
	GroupChannel GroupChannel `json:"GroupChannel,omitempty"`
	Members      []Members    `json:"Members"`
}

type GroupChannel struct {
	Youtube  []interface{} `json:"Youtube"`
	BiliBili []interface{} `json:"BiliBili"`
}

type Data struct {
	Independent Independent `json:"independent"`
	Group       []Group     `json:"Group"`
}

type Subs struct {
	Kind     string   `json:"kind"`
	Etag     string   `json:"etag"`
	PageInfo PageInfo `json:"pageInfo"`
	Items    []Items  `json:"items"`
}
type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}
type Default struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type Medium struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type High struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
type Thumbnails struct {
	Default Default `json:"default"`
	Medium  Medium  `json:"medium"`
	High    High    `json:"high"`
}
type Localized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type Snippet struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	PublishedAt     time.Time  `json:"publishedAt"`
	Thumbnails      Thumbnails `json:"thumbnails"`
	DefaultLanguage string     `json:"defaultLanguage"`
	Localized       Localized  `json:"localized"`
	Country         string     `json:"country"`
}
type Statistics struct {
	ViewCount             string `json:"viewCount"`
	SubscriberCount       string `json:"subscriberCount"`
	HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
	VideoCount            string `json:"videoCount"`
}
type Items struct {
	Kind       string     `json:"kind"`
	Etag       string     `json:"etag"`
	ID         string     `json:"id"`
	Snippet    Snippet    `json:"snippet"`
	Statistics Statistics `json:"statistics"`
}

func (Data Subs) Default() Subs {
	Data.Items = append(Data.Items, Items{
		Statistics: Statistics{
			VideoCount:      "0",
			ViewCount:       "0",
			SubscriberCount: "0",
		},
	})

	return Data
}

type SpaceVideo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		List struct {
			Tlist struct {
				Num1 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"1"`
				Num3 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"3"`
				Num4 struct {
					Tid   int    `json:"tid"`
					Count int    `json:"count"`
					Name  string `json:"name"`
				} `json:"4"`
			} `json:"tlist"`
			Vlist `json:"vlist"`
		} `json:"list"`
		Page struct {
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
			Count int `json:"count"`
		} `json:"page"`
		EpisodicButton struct {
			Text string `json:"text"`
			URI  string `json:"uri"`
		} `json:"episodic_button"`
	} `json:"data"`
}

type Avatar struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int    `json:"mid"`
		Name      string `json:"name"`
		Sex       string `json:"sex"`
		Face      string `json:"face"`
		Sign      string `json:"sign"`
		Rank      int    `json:"rank"`
		Level     int    `json:"level"`
		Jointime  int    `json:"jointime"`
		Moral     int    `json:"moral"`
		Silence   int    `json:"silence"`
		Birthday  string `json:"birthday"`
		Coins     int    `json:"coins"`
		FansBadge bool   `json:"fans_badge"`
		Official  struct {
			Role  int    `json:"role"`
			Title string `json:"title"`
			Desc  string `json:"desc"`
			Type  int    `json:"type"`
		} `json:"official"`
		Vip struct {
			Type      int `json:"type"`
			Status    int `json:"status"`
			ThemeType int `json:"theme_type"`
			Label     struct {
				Path       string `json:"path"`
				Text       string `json:"text"`
				LabelTheme string `json:"label_theme"`
			} `json:"label"`
			AvatarSubscript int    `json:"avatar_subscript"`
			NicknameColor   string `json:"nickname_color"`
		} `json:"vip"`
		Pendant struct {
			Pid          int    `json:"pid"`
			Name         string `json:"name"`
			Image        string `json:"image"`
			Expire       int    `json:"expire"`
			ImageEnhance string `json:"image_enhance"`
		} `json:"pendant"`
		Nameplate struct {
			Nid        int    `json:"nid"`
			Name       string `json:"name"`
			Image      string `json:"image"`
			ImageSmall string `json:"image_small"`
			Level      string `json:"level"`
			Condition  string `json:"condition"`
		} `json:"nameplate"`
		IsFollowed bool   `json:"is_followed"`
		TopPhoto   string `json:"top_photo"`
		Theme      struct {
		} `json:"theme"`
		SysNotice struct {
		} `json:"sys_notice"`
	} `json:"data"`
}

type Vlist []struct {
	Comment      int    `json:"comment"`
	Typeid       int    `json:"typeid"`
	Play         int    `json:"play"`
	Pic          string `json:"pic"`
	Subtitle     string `json:"subtitle"`
	Description  string `json:"description"`
	Copyright    string `json:"copyright"`
	Title        string `json:"title"`
	Review       int    `json:"review"`
	Author       string `json:"author"`
	Mid          int    `json:"mid"`
	Created      int    `json:"created"`
	Length       string `json:"length"`
	VideoReview  int    `json:"video_review"`
	Aid          int    `json:"aid"`
	Bvid         string `json:"bvid"`
	HideClick    bool   `json:"hide_click"`
	IsPay        int    `json:"is_pay"`
	IsUnionVideo int    `json:"is_union_video"`
	VideoType    string
}

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

type SubTbili struct {
	Item struct {
		ID          int           `json:"id"`
		Title       string        `json:"title"`
		Description string        `json:"description"`
		Category    string        `json:"category"`
		Role        []interface{} `json:"role"`
		Source      []interface{} `json:"source"`
		Pictures    []struct {
			ImgSrc    string `json:"img_src"`
			ImgWidth  int    `json:"img_width"`
			ImgHeight int    `json:"img_height"`
			ImgSize   int    `json:"img_size"`
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

type DataTBiliBili struct {
	URL             string
	Author          string
	Avatar          string
	Like            int
	Photos          string
	Videos          string
	Text            string
	Dynamic_id      string
	VtuberMember_id int64
}
