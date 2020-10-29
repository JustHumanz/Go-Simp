package bilibili

import (
	"time"

	database "github.com/JustHumanz/Go-simp/database"
	"github.com/bwmarrin/discordgo"
)

type Bilibili struct {
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
			Vlist []struct {
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
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Count int `json:" count"`
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
		} `json:"page"`
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

type getInfoByRoom struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		RoomInfo struct {
			UID            int    `json:"uid"`
			RoomID         int    `json:"room_id"`
			ShortID        int    `json:"short_id"`
			Title          string `json:"title"`
			Cover          string `json:"cover"`
			Tags           string `json:"tags"`
			Background     string `json:"background"`
			Description    string `json:"description"`
			LiveStatus     int    `json:"live_status"`
			LiveStartTime  int    `json:"live_start_time"`
			LiveScreenType int    `json:"live_screen_type"`
			LockStatus     int    `json:"lock_status"`
			LockTime       int    `json:"lock_time"`
			HiddenStatus   int    `json:"hidden_status"`
			HiddenTime     int    `json:"hidden_time"`
			AreaID         int    `json:"area_id"`
			AreaName       string `json:"area_name"`
			ParentAreaID   int    `json:"parent_area_id"`
			ParentAreaName string `json:"parent_area_name"`
			Keyframe       string `json:"keyframe"`
			SpecialType    int    `json:"special_type"`
			UpSession      string `json:"up_session"`
			PkStatus       int    `json:"pk_status"`
			IsStudio       bool   `json:"is_studio"`
			Pendants       struct {
				Frame struct {
					Name  string `json:"name"`
					Value string `json:"value"`
					Desc  string `json:"desc"`
				} `json:"frame"`
			} `json:"pendants"`
			OnVoiceJoin int `json:"on_voice_join"`
			Online      int `json:"online"`
			RoomType    struct {
				Four1 int `json:"4-1"`
			} `json:"room_type"`
		} `json:"room_info"`
		AnchorInfo struct {
			BaseInfo struct {
				Uname        string `json:"uname"`
				Face         string `json:"face"`
				Gender       string `json:"gender"`
				OfficialInfo struct {
					Role  int    `json:"role"`
					Title string `json:"title"`
					Desc  string `json:"desc"`
				} `json:"official_info"`
			} `json:"base_info"`
			LiveInfo struct {
				Level        int    `json:"level"`
				LevelColor   int    `json:"level_color"`
				Score        int    `json:"score"`
				UpgradeScore int    `json:"upgrade_score"`
				Current      []int  `json:"current"`
				Next         []int  `json:"next"`
				Rank         string `json:"rank"`
			} `json:"live_info"`
			RelationInfo struct {
				Attention int `json:"attention"`
			} `json:"relation_info"`
			MedalInfo struct {
				MedalName string `json:"medal_name"`
				MedalID   int    `json:"medal_id"`
				Fansclub  int    `json:"fansclub"`
			} `json:"medal_info"`
		} `json:"anchor_info"`
		NewsInfo struct {
			UID     int    `json:"uid"`
			Ctime   string `json:"ctime"`
			Content string `json:"content"`
		} `json:"news_info"`
		RankdbInfo struct {
			Roomid    int    `json:"roomid"`
			RankDesc  string `json:"rank_desc"`
			Color     string `json:"color"`
			H5URL     string `json:"h5_url"`
			WebURL    string `json:"web_url"`
			Timestamp int    `json:"timestamp"`
		} `json:"rankdb_info"`
		AreaRankInfo struct {
			AreaRank struct {
				Index int    `json:"index"`
				Rank  string `json:"rank"`
			} `json:"areaRank"`
			LiveRank struct {
				Rank string `json:"rank"`
			} `json:"liveRank"`
		} `json:"area_rank_info"`
		BattleRankEntryInfo struct {
			FirstRankImgURL string `json:"first_rank_img_url"`
			RankName        string `json:"rank_name"`
			ShowStatus      int    `json:"show_status"`
		} `json:"battle_rank_entry_info"`
		TabInfo struct {
			List []struct {
				Type      string `json:"type"`
				Desc      string `json:"desc"`
				IsFirst   int    `json:"isFirst"`
				IsEvent   int    `json:"isEvent"`
				EventType string `json:"eventType"`
				ListType  string `json:"listType"`
				APIPrefix string `json:"apiPrefix"`
				RankName  string `json:"rank_name"`
			} `json:"list"`
		} `json:"tab_info"`
		ActivityInitInfo struct {
			EventList []interface{} `json:"eventList"`
			WeekInfo  struct {
				BannerInfo interface{} `json:"bannerInfo"`
				GiftName   interface{} `json:"giftName"`
			} `json:"weekInfo"`
			GiftName interface{} `json:"giftName"`
			Lego     struct {
				Timestamp int    `json:"timestamp"`
				Config    string `json:"config"`
			} `json:"lego"`
		} `json:"activity_init_info"`
		VoiceJoinInfo struct {
			Status struct {
				Open        int    `json:"open"`
				AnchorOpen  int    `json:"anchor_open"`
				Status      int    `json:"status"`
				UID         int    `json:"uid"`
				UserName    string `json:"user_name"`
				HeadPic     string `json:"head_pic"`
				Guard       int    `json:"guard"`
				StartAt     int    `json:"start_at"`
				CurrentTime int    `json:"current_time"`
			} `json:"status"`
			Icons struct {
				IconClose    string `json:"icon_close"`
				IconOpen     string `json:"icon_open"`
				IconWait     string `json:"icon_wait"`
				IconStarting string `json:"icon_starting"`
			} `json:"icons"`
			WebShareLink string `json:"web_share_link"`
		} `json:"voice_join_info"`
		AdBannerInfo struct {
			Data []struct {
				ID       int    `json:"id"`
				Title    string `json:"title"`
				Location string `json:"location"`
				Position int    `json:"position"`
				Pic      string `json:"pic"`
				Link     string `json:"link"`
				Weight   int    `json:"weight"`
			} `json:"data"`
		} `json:"ad_banner_info"`
		SkinInfo struct {
			ID          int    `json:"id"`
			SkinName    string `json:"skin_name"`
			SkinConfig  string `json:"skin_config"`
			ShowText    string `json:"show_text"`
			SkinURL     string `json:"skin_url"`
			StartTime   int    `json:"start_time"`
			EndTime     int    `json:"end_time"`
			CurrentTime int    `json:"current_time"`
		} `json:"skin_info"`
		WebBannerInfo struct {
			ID               int    `json:"id"`
			Title            string `json:"title"`
			Left             string `json:"left"`
			Right            string `json:"right"`
			JumpURL          string `json:"jump_url"`
			BgColor          string `json:"bg_color"`
			HoverColor       string `json:"hover_color"`
			TextBgColor      string `json:"text_bg_color"`
			TextHoverColor   string `json:"text_hover_color"`
			LinkText         string `json:"link_text"`
			LinkColor        string `json:"link_color"`
			InputColor       string `json:"input_color"`
			InputTextColor   string `json:"input_text_color"`
			InputHoverColor  string `json:"input_hover_color"`
			InputBorderColor string `json:"input_border_color"`
			InputSearchColor string `json:"input_search_color"`
		} `json:"web_banner_info"`
		LolInfo struct {
			LolActivity struct {
				Status     int    `json:"status"`
				GuessCover string `json:"guess_cover"`
				VoteCover  string `json:"vote_cover"`
				VoteH5URL  string `json:"vote_h5_url"`
				VoteUseH5  bool   `json:"vote_use_h5"`
			} `json:"lol_activity"`
		} `json:"lol_info"`
		WishListInfo struct {
			List   interface{} `json:"list"`
			Status int         `json:"status"`
		} `json:"wish_list_info"`
		ScoreCardInfo  interface{} `json:"score_card_info"`
		PkInfo         interface{} `json:"pk_info"`
		BattleInfo     interface{} `json:"battle_info"`
		SilentRoomInfo struct {
			Type       string `json:"type"`
			Level      int    `json:"level"`
			Second     int    `json:"second"`
			ExpireTime int    `json:"expire_time"`
		} `json:"silent_room_info"`
		SwitchInfo struct {
			CloseGuard   bool `json:"close_guard"`
			CloseGift    bool `json:"close_gift"`
			CloseOnline  bool `json:"close_online"`
			CloseDanmaku bool `json:"close_danmaku"`
		} `json:"switch_info"`
		RecordSwitchInfo struct {
			RecordTab bool `json:"record_tab"`
		} `json:"record_switch_info"`
		RoomConfigInfo struct {
			DmText string `json:"dm_text"`
		} `json:"room_config_info"`
		GiftMemoryInfo struct {
			List interface{} `json:"list"`
		} `json:"gift_memory_info"`
		NewSwitchInfo struct {
			RoomSocket           int `json:"room-socket"`
			RoomPropSend         int `json:"room-prop-send"`
			RoomSailing          int `json:"room-sailing"`
			RoomInfoPopularity   int `json:"room-info-popularity"`
			RoomDanmakuEditor    int `json:"room-danmaku-editor"`
			RoomEffect           int `json:"room-effect"`
			RoomFansMedal        int `json:"room-fans_medal"`
			RoomReport           int `json:"room-report"`
			RoomFeedback         int `json:"room-feedback"`
			RoomPlayerWatermark  int `json:"room-player-watermark"`
			RoomRecommendLiveOff int `json:"room-recommend-live_off"`
			RoomActivity         int `json:"room-activity"`
			RoomWebBanner        int `json:"room-web_banner"`
			RoomSilverSeedsBox   int `json:"room-silver_seeds-box"`
			RoomWishingBottle    int `json:"room-wishing_bottle"`
			RoomBoard            int `json:"room-board"`
			RoomSupplication     int `json:"room-supplication"`
			RoomHourRank         int `json:"room-hour_rank"`
			RoomWeekRank         int `json:"room-week_rank"`
			RoomAnchorRank       int `json:"room-anchor_rank"`
			RoomInfoIntegral     int `json:"room-info-integral"`
			RoomSuperChat        int `json:"room-super-chat"`
			RoomTab              int `json:"room-tab"`
		} `json:"new_switch_info"`
		SuperChatInfo struct {
			Status      int           `json:"status"`
			JumpURL     string        `json:"jump_url"`
			Icon        string        `json:"icon"`
			RankedMark  int           `json:"ranked_mark"`
			MessageList []interface{} `json:"message_list"`
		} `json:"super_chat_info"`
		VideoConnectionInfo interface{} `json:"video_connection_info"`
		PlayerThrottleInfo  struct {
			Status              int `json:"status"`
			NormalSleepTime     int `json:"normal_sleep_time"`
			FullscreenSleepTime int `json:"fullscreen_sleep_time"`
			TabSleepTime        int `json:"tab_sleep_time"`
			PromptTime          int `json:"prompt_time"`
		} `json:"player_throttle_info"`
	} `json:"data"`
}

type RoomID struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		RoomStatus    int    `json:"roomStatus"`
		RoundStatus   int    `json:"roundStatus"`
		LiveStatus    int    `json:"liveStatus"`
		URL           string `json:"url"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		Online        int    `json:"online"`
		Roomid        int    `json:"roomid"`
		BroadcastType int    `json:"broadcast_type"`
		OnlineHidden  int    `json:"online_hidden"`
	} `json:"data"`
}

type NewSchedule struct {
	DurationUp   time.Duration
	DurationPast time.Duration
	RealTime     time.Time
	Desc         string
	Publish      time.Time
	ImgURL       string
}

type Card struct {
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

type TimeLine struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		HasMore int `json:"has_more"`
		Cards   []struct {
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
				UIDType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerID      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIDStr string `json:"dynamic_id_str"`
				PreDyIDStr   string `json:"pre_dy_id_str"`
				OrigDyIDStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
			} `json:"desc,omitempty"`
			Card       string `json:"card"`
			ExtendJSON string `json:"extend_json"`
			Extra      struct {
				IsSpaceTop int `json:"is_space_top"`
			} `json:"extra"`
			Display struct {
				EmojiInfo struct {
					EmojiDetails []struct {
						EmojiName string `json:"emoji_name"`
						ID        int    `json:"id"`
						PackageID int    `json:"package_id"`
						State     int    `json:"state"`
						Type      int    `json:"type"`
						Attr      int    `json:"attr"`
						Text      string `json:"text"`
						URL       string `json:"url"`
						Meta      struct {
							Size int `json:"size"`
						} `json:"meta"`
						Mtime int `json:"mtime"`
					} `json:"emoji_details"`
				} `json:"emoji_info"`
				Relation struct {
					Status     int `json:"status"`
					IsFollow   int `json:"is_follow"`
					IsFollowed int `json:"is_followed"`
				} `json:"relation"`
			} `json:"display,omitempty"`
			Desc1 struct {
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
				UIDType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerID      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIDStr string `json:"dynamic_id_str"`
				PreDyIDStr   string `json:"pre_dy_id_str"`
				OrigDyIDStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
				Bvid         string `json:"bvid"`
			} `json:"desc,omitempty"`
			Display1 struct {
				TopicInfo struct {
					TopicDetails []struct {
						TopicID    int    `json:"topic_id"`
						TopicName  string `json:"topic_name"`
						IsActivity int    `json:"is_activity"`
						TopicLink  string `json:"topic_link"`
					} `json:"topic_details"`
				} `json:"topic_info"`
				UsrActionTxt string `json:"usr_action_txt"`
				Relation     struct {
					Status     int `json:"status"`
					IsFollow   int `json:"is_follow"`
					IsFollowed int `json:"is_followed"`
				} `json:"relation"`
			} `json:"display,omitempty"`
			Desc2 struct {
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
				UIDType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerID      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIDStr string `json:"dynamic_id_str"`
				PreDyIDStr   string `json:"pre_dy_id_str"`
				OrigDyIDStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
				Bvid         string `json:"bvid"`
			} `json:"desc,omitempty"`
			Display2 struct {
				TopicInfo struct {
					TopicDetails []struct {
						TopicID    int    `json:"topic_id"`
						TopicName  string `json:"topic_name"`
						IsActivity int    `json:"is_activity"`
						TopicLink  string `json:"topic_link"`
					} `json:"topic_details"`
				} `json:"topic_info"`
				UsrActionTxt string `json:"usr_action_txt"`
				Relation     struct {
					Status     int `json:"status"`
					IsFollow   int `json:"is_follow"`
					IsFollowed int `json:"is_followed"`
				} `json:"relation"`
			} `json:"display,omitempty"`
			Desc3 struct {
				UID         int   `json:"uid"`
				Type        int   `json:"type"`
				Rid         int64 `json:"rid"`
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
				UIDType      int    `json:"uid_type"`
				Stype        int    `json:"stype"`
				RType        int    `json:"r_type"`
				InnerID      int    `json:"inner_id"`
				Status       int    `json:"status"`
				DynamicIDStr string `json:"dynamic_id_str"`
				PreDyIDStr   string `json:"pre_dy_id_str"`
				OrigDyIDStr  string `json:"orig_dy_id_str"`
				RidStr       string `json:"rid_str"`
			} `json:"desc,omitempty"`
		} `json:"cards"`
		NextOffset int64 `json:"next_offset"`
		Gt         int   `json:"_gt_"`
	} `json:"data"`
}

type BiliBiliSchedule struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		UserInfos            interface{} `json:"user_infos"`
		ProgramInfos         interface{} `json:"program_infos"`
		Timestamp            int         `json:"timestamp"`
		AnchorNotInWhiteList int         `json:"anchor_not_in_white_list"`
		GlobalSwitchClose    int         `json:"global_switch_close"`
	} `json:"data"`
}

type ProgramList struct {
	Ruid           int    `json:"ruid"`
	IsSubscription int    `json:"is_subscription"`
	StartTime      int    `json:"start_time"`
	IsRecommend    int    `json:"is_recommend"`
	SubscriptionID int    `json:"subscription_id"`
	Title          string `json:"title"`
	ProgramID      int    `json:"program_id"`
	Expired        int    `json:"expired"`
	RoomID         int    `json:"room_id"`
}

type UserInfo struct {
	Ruid        int    `json:"ruid"`
	Uname       string `json:"uname"`
	Face        string `json:"face"`
	Description string `json:"description"`
	RoomID      int    `json:"room_id"`
}

type ScheduleData struct {
	User []UserInfo
	List []ProgramList
}

type RoomID2 struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		RoomStatus    int    `json:"roomStatus"`
		RoundStatus   int    `json:"roundStatus"`
		LiveStatus    int    `json:"liveStatus"`
		URL           string `json:"url"`
		Title         string `json:"title"`
		Cover         string `json:"cover"`
		Online        int    `json:"online"`
		Roomid        int    `json:"roomid"`
		BroadcastType int    `json:"broadcast_type"`
		OnlineHidden  int    `json:"online_hidden"`
	} `json:"data"`
}

type LiveBili struct {
	RoomData *database.LiveBiliDB
	Member   database.Name
	Embed    *discordgo.MessageEmbed
}

func (Data *LiveBili) AddData(new *database.LiveBiliDB) *LiveBili {
	Data.RoomData = new
	return Data
}

func (Data *LiveBili) UpdateSchdule(new time.Time) *LiveBili {
	Data.RoomData.ScheduledStart = new
	return Data
}

func (Data *LiveBili) SetStatus(new string) *LiveBili {
	Data.RoomData.Status = new
	return Data
}

func (Data *LiveBili) UpdateOnline(new int) *LiveBili {
	Data.RoomData.Online = new
	return Data
}

func (Data *LiveBili) AddMember(new database.Name) *LiveBili {
	Data.Member = new
	return Data
}
