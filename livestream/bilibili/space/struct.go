package space

type CheckSctruct struct {
	SpaceID                                            int
	MemberID                                           int64
	VideoList                                          Vlist
	MemberName, MemberFace, MemberUrl, Name, GroupIcon string
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
