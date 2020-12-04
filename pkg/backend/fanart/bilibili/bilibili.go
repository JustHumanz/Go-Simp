package bilibili

import (
	"encoding/json"
	"net/url"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	network "github.com/JustHumanz/Go-simp/tools/network"
)

//CheckNew Start Check new fanart
func CheckNew() {
	for _, Group := range engine.GroupData {
		wg := new(sync.WaitGroup)
		for _, Member := range database.GetMembers(Group.ID) {
			wg.Add(1)
			go func(Group database.Group, Member database.Member, wg *sync.WaitGroup) {
				defer wg.Done()
				if Member.BiliBiliHashtags != "" {
					log.WithFields(log.Fields{
						"Group":  Group.NameGroup,
						"Vtuber": Member.EnName,
					}).Info("Start crawler bilibili")
					var (
						body    []byte
						errcurl error
						urls    = "https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name=" + url.QueryEscape(Member.BiliBiliHashtags)
					)
					body, err := network.Curl(urls, nil)
					if err != nil {
						log.Info("Trying use tor")

						body, errcurl = network.CoolerCurl(urls, nil)
						if errcurl != nil {
							log.Error(errcurl)
						} else {
							log.Info("Successfully use tor")
						}
					}
					var (
						TB              TBiliBili
						DynamicIDStrTmp string
					)
					_ = json.Unmarshal(body, &TB)
					if (len(TB.Data.Cards) > 0) && TB.Data.Cards[0].Desc.DynamicIDStr != DynamicIDStrTmp {
						DynamicIDStrTmp = TB.Data.Cards[0].Desc.DynamicIDStr
						for i := 0; i < len(TB.Data.Cards); i++ {
							var (
								STB  SubTbili
								img  []string
								nope bool
							)
							_ = json.Unmarshal([]byte(TB.Data.Cards[i].Card), &STB)
							if STB.Item.Pictures != nil && TB.Data.Cards[i].Desc.Type == 2 { //type 2 is picture post (prob,heheheh)
								niggerlist := []string{"解锁专属粉丝卡片", "twitter", "咖啡厅", "cd", "专辑", "pixiv", "遇", "marshmallow", "saucenao", "pid", "twi"}
								nope, _ = regexp.MatchString("(?m)("+strings.Join(niggerlist, "|")+")", strings.ToLower(STB.Item.Description))
								New := database.GetTBiliBili(TB.Data.Cards[i].Desc.DynamicIDStr)

								if New && !nope {
									link, color, err := STB.Mirroring()
									if err != nil {
										log.WithFields(log.Fields{
											"Group":  Group.NameGroup,
											"Vtuber": Member.EnName,
										}).Error(err)
										break
									}
									if link != "" {
										log.WithFields(log.Fields{
											"Vtuber": Member.EnName,
											"Img":    link,
										}).Info("New Fanart")
										for l := 0; l < len(STB.Item.Pictures); l++ {
											img = append(img, STB.Item.Pictures[l].ImgSrc)
										}

										log.Info("Send to database")
										log.WithFields(log.Fields{"Group": Group.NameGroup, "Vtuber": Member.EnName}).Info("Push to notif")
										Data := Notif{
											TBiliData: database.InputTBiliBili{
												URL:        "https://t.bilibili.com/" + TB.Data.Cards[i].Desc.DynamicIDStr + "?tab=2",
												Author:     TB.Data.Cards[i].Desc.UserProfile.Info.Uname,
												Avatar:     TB.Data.Cards[i].Desc.UserProfile.Info.Face,
												Like:       TB.Data.Cards[i].Desc.Like,
												Photos:     strings.Join(img, "\n"),
												Dynamic_id: TB.Data.Cards[i].Desc.DynamicIDStr,
												Text:       STB.Item.Description,
											},
											Group:       Group,
											PhotosCount: STB.Item.PicturesCount,
											PhotosImgur: link,
											MemberID:    Member.ID,
										}
										Data.TBiliData.InputTBiliBili(Member.ID)
										Data.PushNotif(color)
									}
								}
							}
						}
					} else {
						log.WithFields(log.Fields{
							"Group":  Group.NameGroup,
							"Vtuber": Member.EnName,
						}).Info("Still same")
					}
					//time.Sleep(time.Duration(int64(rand.Intn((7-1)+1))) * time.Second)
				}
			}(Group, Member, wg)
		}
		wg.Wait()
	}
}

//Mirroring fanart to imgur *sometime discord fail to load image because bilibili CDN(prob)*
func (Data SubTbili) Mirroring() (string, int, error) {
	link := Data.Item.Pictures[0].ImgSrc
	/*
		buf := &bytes.Buffer{}
		writer := multipart.NewWriter(buf)
		err := writer.WriteField("image", link)
		err = writer.WriteField("title", Data.Item.Title)
		err = writer.WriteField("name", Data.User.Name)
		err = writer.Close()
		if err != nil {
			log.Error(err)
		}
	*/
	color, err := engine.GetColor("/tmp/bilibili", link)
	if err != nil {
		log.Error(err)
	}
	/*
		req, err := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
		if err != nil {
			return "", 0, err
		}

		req.Header.Set("Authorization", "Client-ID "+config.ImgurClient)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Votre) AppleWebKit/601.2 (KHTML, like Gecko)")
		req.Header.Set("Content-Type", writer.FormDataContentType())

		htt := http.Client{Timeout: time.Second * 20}
		res, err := htt.Do(req)
		if err != nil {
			log.Error(err)
			log.Info("bypass it")
			return link, color, nil
		}

		defer res.Body.Close()

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
			log.Info("bypass it")
			return link, color, nil
		}

		if err != nil || res.StatusCode != 200 {
			log.Error(err, res.Status, string(b))
			log.Info("bypass it")
			return link, color, nil
		}

		var tmp Imgur
		err = json.Unmarshal(b, &tmp)
		if err != nil {
			log.Error(err)
			return "", 0, err
		}

		color, err = engine.GetColor("/tmp/bilibili", tmp.Data.Link)
		if err != nil {
			log.Error(err, " ", link)
			return "", 0, err
		}

	*/
	return link, color, nil
}

//Imgur struct
type Imgur struct {
	Data struct {
		ID          string        `json:"id"`
		Title       interface{}   `json:"title"`
		Description interface{}   `json:"description"`
		Datetime    int           `json:"datetime"`
		Type        string        `json:"type"`
		Animated    bool          `json:"animated"`
		Width       int           `json:"width"`
		Height      int           `json:"height"`
		Size        int           `json:"size"`
		Views       int           `json:"views"`
		Bandwidth   int           `json:"bandwidth"`
		Vote        interface{}   `json:"vote"`
		Favorite    bool          `json:"favorite"`
		Nsfw        interface{}   `json:"nsfw"`
		Section     interface{}   `json:"section"`
		AccountURL  interface{}   `json:"account_url"`
		AccountID   int           `json:"account_id"`
		IsAd        bool          `json:"is_ad"`
		InMostViral bool          `json:"in_most_viral"`
		Tags        []interface{} `json:"tags"`
		AdType      int           `json:"ad_type"`
		AdURL       string        `json:"ad_url"`
		InGallery   bool          `json:"in_gallery"`
		Deletehash  string        `json:"deletehash"`
		Name        string        `json:"name"`
		Link        string        `json:"link"`
	} `json:"data"`
	Success bool `json:"success"`
	Status  int  `json:"status"`
}
