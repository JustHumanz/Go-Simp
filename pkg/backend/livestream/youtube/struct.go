package youtube

import (
	"encoding/xml"
	"time"

	"github.com/JustHumanz/Go-simp/tools/database"
)

type YtXML struct {
	XMLName xml.Name `xml:"feed"`
	Text    string   `xml:",chardata"`
	Link    []struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
		Type string `xml:"type,attr"`
		ID   string `xml:"id,attr"`
	} `xml:"link"`
	ID        string `xml:"id"`
	ChannelId string `xml:"channelId"`
	Title     string `xml:"title"`
	Author    struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
		URI  string `xml:"uri"`
	} `xml:"author"`
	Published string `xml:"published"`
	Entry     []struct {
		Text      string `xml:",chardata"`
		ID        string `xml:"id"`
		VideoId   string `xml:"videoId"`
		ChannelId string `xml:"channelId"`
		Title     string `xml:"title"`
		Link      struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Author struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
			URI  string `xml:"uri"`
		} `xml:"author"`
		Published string `xml:"published"`
		Updated   string `xml:"updated"`
		Group     struct {
			Text    string `xml:",chardata"`
			Title   string `xml:"title"`
			Content struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Type   string `xml:"type,attr"`
				Width  string `xml:"width,attr"`
				Height string `xml:"height,attr"`
			} `xml:"content"`
			Thumbnail struct {
				Text   string `xml:",chardata"`
				URL    string `xml:"url,attr"`
				Width  string `xml:"width,attr"`
				Height string `xml:"height,attr"`
			} `xml:"thumbnail"`
			Description string `xml:"description"`
			Community   struct {
				Text       string `xml:",chardata"`
				StarRating struct {
					Text    string `xml:",chardata"`
					Count   string `xml:"count,attr"`
					Average string `xml:"average,attr"`
					Min     string `xml:"min,attr"`
					Max     string `xml:"max,attr"`
				} `xml:"starRating"`
				Statistics struct {
					Text  string `xml:",chardata"`
					Views string `xml:"views,attr"`
				} `xml:"statistics"`
			} `xml:"community"`
		} `xml:"group"`
	} `xml:"entry"`
	Style struct {
		Text string `xml:",chardata"`
		Type string `xml:"type,attr"`
		ID   string `xml:"id,attr"`
	} `xml:"style"`
	Script string `xml:"script"`
}

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
		ContentDetails struct {
			Duration string `json:"duration"`
		} `json:"contentDetails"`
		Statistics struct {
			ViewCount string `json:"viewCount"`
		} `json:"statistics"`
		LiveDetails struct {
			StartTime       time.Time `json:"scheduledStartTime"`
			ActualStartTime time.Time `json:"actualStartTime"`
			EndTime         time.Time `json:"actualEndTime"`
			Viewers         string    `json:"concurrentViewers"`
		} `json:"liveStreamingDetails"`
	} `json:"items"`
}

type NotifStruct struct {
	YtData        *database.YtDbData
	ActuallyStart time.Time
	Member        database.Name
	Group         database.GroupName
}

func (Data *NotifStruct) AddData(new *database.YtDbData) *NotifStruct {
	Data.YtData = new
	return Data
}

func (Data *NotifStruct) ChangeYtStatus(new string) *NotifStruct {
	Data.YtData.Status = new
	return Data
}

func (Data *NotifStruct) SetActuallyStart(new time.Time) *NotifStruct {
	Data.YtData.Schedul = new
	return Data
}

func (Data *NotifStruct) UpdateYtDB() {
	Data.YtData.UpdateYt(Data.YtData.Status)
}

func (Data *NotifStruct) SendtoDB() error {
	err := Data.YtData.InputYt(Data.Member.ID)
	if err != nil {
		return err
	}
	return nil
}

func (Data *NotifStruct) UpYtView(new string) *NotifStruct {
	Data.YtData.Viewers = new
	return Data
}

func (Data *NotifStruct) UpYtEnd(new time.Time) *NotifStruct {
	Data.YtData.End = new
	return Data
}

func (Data *NotifStruct) UpYtLen(new string) *NotifStruct {
	Data.YtData.Length = new
	return Data
}

func (Data *NotifStruct) UpYtSchedul(new time.Time) *NotifStruct {
	Data.YtData.Schedul = new
	return Data
}
