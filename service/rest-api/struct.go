package main

type MembersPayload struct {
	ID       int64
	NickName string
	EnName   string
	JpName   string
	Region   string
	Fanbase  string
	Status   string
	BiliBili interface{}
	Youtube  interface{}
	Twitter  interface{}
	Twitch   interface{}
	Group    interface{}
	IsLive   interface{}
}

type GroupPayload struct {
	ID        int64
	GroupName string
	GroupIcon string
	Youtube   interface{}
}

type Twitter struct {
	TwitterFanart   string `json:"Twitter_Fanart"`
	TwitterLewd     string `json:"Twitter_Lewd"`
	TwitterUsername string `json:"Twitter_Username"`
}
type Youtube struct {
	YtID string `json:"Youtube_ID"`
}
type BiliBili struct {
	BiliBiliFanart string `json:"BiliBili_Fanart"`
	BiliBiliID     int    `json:"BiliBili_ID"`
	BiliRoomID     int    `json:"BiliRoom_ID"`
}
type Twitch struct {
	TwitchUsername string `json:"Twitch_Username"`
}

type GroupChannel struct {
	Youtube  []interface{} `json:"Youtube"`
	BiliBili []interface{} `json:"BiliBili"`
}

type Agency struct {
	GroupName    string       `json:"GroupName"`
	GroupIcon    string       `json:"GroupIcon"`
	GroupChannel GroupChannel `json:"GroupChannel,omitempty"`
	Members      []struct {
		Name     string   `json:"Name"`
		EnName   string   `json:"EN_Name"`
		JpName   string   `json:"JP_Name"`
		GroupID  int      `json:"GroupID"`
		Twitter  Twitter  `json:"Twitter,omitempty"`
		Youtube  Youtube  `json:"Youtube,omitempty"`
		BiliBili BiliBili `json:"BiliBili,omitempty"`
		Status   string   `json:"Status"`
		Twitch   Twitch   `json:"Twitch"`
		Region   string   `json:"Region"`
		Fanbase  string   `json:"Fanbase"`
	}
}
