package pixiv

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/JustHumanz/Go-Simp/service/fanart/notif"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	Bot         *discordgo.Session
	VtubersData database.VtubersPayload
	lewd        bool
)

const (
	BaseURL = "https://www.pixiv.net/en/artworks/"
	Limit   = 10
)

//Start start twitter module
func Start(a *discordgo.Session, b *cron.Cron, c database.VtubersPayload, d config.ConfigFile, e bool) {
	Bot = a
	VtubersData = c
	lewd = e
	b.AddFunc(config.PixivFanart, CheckPixiv)
	CheckPixiv()
	CheckPixivLewd()
	if lewd {
		b.AddFunc(config.PixivFanartLewd, CheckPixivLewd)
		log.Info("Enable Pixiv lewd fanart module")
	} else {
		log.Info("Enable Pixiv fanart module")
	}
}

//CheckNew Check new fanart
func CheckPixiv() {
	for _, Group := range VtubersData.VtuberData {
		var wg sync.WaitGroup
		for i, Member := range Group.Members {
			wg.Add(1)
			go func(wg *sync.WaitGroup, Member database.Member) {
				defer wg.Done()
				FixFanArt := &database.DataFanart{
					Member: Member,
					Group:  Group,
					Lewd:   false,
				}

				if Member.JpName != "" {
					log.WithFields(log.Fields{
						"Member": Member.JpName,
						"Group":  Group.GroupName,
						"Lewd":   false,
					}).Info("Start curl pixiv")
					URLJP := GetPixivURL(url.QueryEscape(Member.JpName))
					err := Pixiv(URLJP, FixFanArt, false)
					if err != nil {
						log.Error(err)
					}
				}

				if Member.EnName == Member.Name {
					if Member.EnName != "" {
						log.WithFields(log.Fields{
							"Member": Member.EnName,
							"Group":  Group.GroupName,
							"Lewd":   false,
						}).Info("Start curl pixiv")
						URLEN := GetPixivURL(engine.UnderScoreName(Member.EnName))
						err := Pixiv(URLEN, FixFanArt, false)
						if err != nil {
							log.Error(err)
						}

					}
				} else {
					if Member.EnName != "" {
						log.WithFields(log.Fields{
							"Member": Member.EnName,
							"Group":  Group.GroupName,
							"Lewd":   false,
						}).Info("Start curl pixiv")
						URLEN := GetPixivURL(engine.UnderScoreName(Member.EnName))
						err := Pixiv(URLEN, FixFanArt, false)
						if err != nil {
							log.Error(err)
						}

					}
				}

				if Member.TwitterHashtags != "" {
					log.WithFields(log.Fields{
						"Member": Member.Name,
						"Group":  Group.GroupName,
						"Lewd":   false,
					}).Info("Start curl pixiv")
					URL := GetPixivURL(engine.UnderScoreName(Member.TwitterHashtags[1:]))
					err := Pixiv(URL, FixFanArt, false)
					if err != nil {
						log.Error(err)
					}
				} else if Member.BiliBiliHashtags != "" {
					log.WithFields(log.Fields{
						"Member": Member.Name,
						"Group":  Group.GroupName,
						"Lewd":   false,
					}).Info("Start curl pixiv")
					URL := GetPixivURL(engine.UnderScoreName(Member.BiliBiliHashtags[1 : len(Member.BiliBiliHashtags)-1]))
					err := Pixiv(URL, FixFanArt, false)
					if err != nil {
						log.Error(err)
					}
				}

			}(&wg, Member)
			if i%4 == 0 {
				wg.Wait()
			}
		}
		wg.Wait()
	}
}

func Pixiv(p string, FixFanArt *database.DataFanart, l bool) error {
	var Art PixivArtworks
	req, err := http.NewRequest("GET", p, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:87.0) Gecko/20100101 Firefox/87.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Connection", "keep-alive")
	if l {
		req.Header.Set("Cookie", "PHPSESSID="+config.GoSimpConf.PixivSession)
	}
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Te", "Trailers")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if l && response.StatusCode == http.StatusUnauthorized {
			pilot.ReportDeadService("Pixiv Session outdate")
		}
		log.WithFields(log.Fields{
			"Status":  response.StatusCode,
			"Reason":  response.Status,
			"Payload": p,
		}).Error("Status code not daijobu")
		return errors.New(response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &Art)
	if err != nil {
		return err
	}

	if Art.Body.Illustmanga.Data != nil {
		for i, v := range Art.Body.Illustmanga.Data {
			v2 := v.(map[string]interface{})
			IsVtuber := false
			IsNotLoli := true

			for _, tag := range v2["tags"].([]interface{}) {
				Tag := strings.ToLower(tag.(string))
				GrpName := strings.ToLower(FixFanArt.Group.GroupName)
				HashTag := "vtuber"
				if FixFanArt.Member.TwitterHashtags != "" {
					HashTag = FixFanArt.Member.TwitterHashtags[1:]
				} else if FixFanArt.Member.BiliBiliHashtags != "" {
					HashTag = FixFanArt.Member.BiliBiliHashtags[1 : len(FixFanArt.Member.BiliBiliHashtags)-1]
				}

				match, _ := regexp.MatchString("(youtube|vtuber|"+strings.ToLower(HashTag)+"|"+GrpName+")", Tag)
				if match {
					IsVtuber = true
				}

				if l {
					for _, black := range config.BlackList {
						if strings.ToLower(black) == Tag {
							log.WithFields(log.Fields{
								"URL": BaseURL + v2["id"].(string),
							}).Info("Lol,it's loli")
							IsNotLoli = false
						}
					}
				}
			}

			var (
				Illusts map[string]interface{}
				User    map[string]interface{}
				TextFix string
			)

			if IsVtuber {
				if v2["xRestrict"].(float64) == 0 && !l {
					illusbyte, err := network.Curl(config.PixivIllustsEnd+v2["id"].(string), nil)
					if err != nil {
						return err
					}

					err = json.Unmarshal(illusbyte, &Illusts)
					if err != nil {
						return err
					}

					Body := Illusts["body"].(map[string]interface{})
					Tags := Body["tags"].(map[string]interface{})
					Img := Body["urls"].(map[string]interface{})
					FixImg := Img["original"].(string)

					usrbyte, err := network.Curl(config.PixivUserEnd+Tags["authorId"].(string), nil)
					if err != nil {
						return err
					}

					err = json.Unmarshal(usrbyte, &User)
					if err != nil {
						return err
					}

					UserBody := User["body"].(map[string]interface{})

					Desc := RemoveHtmlTag(Body["description"].(string))
					if match, _ := regexp.MatchString("http://twitter.com", Desc); match {
						TextFix = ClearTwitterURL(Desc)
					} else {
						TextFix = "**" + Body["title"].(string) + "**\n" + Desc
					}

					FixFanArt.AddAuthor(v2["userName"].(string)).AddPermanentURL(BaseURL + v2["id"].(string)).
						AddAuthorAvatar(config.PixivProxy + UserBody["imageBig"].(string)).AddPhotos([]string{FixImg}).
						AddText(TextFix).AddPixivID(v2["id"].(string)).SetState(config.PixivArt)

					new, err := FixFanArt.CheckPixivFanArt()
					if err != nil {
						return err
					}

					if new {
						path, err := DownloadImg(Img["mini"].(string))
						if err != nil {
							log.Error(err)
						}

						Color, err := engine.GetColor("", path)
						if err != nil {
							return err
						}
						FixFanArt.Photos[0] = config.PixivProxy + FixImg
						notif.SendNude(*FixFanArt, Bot, Color)
					}
				} else if l && v2["xRestrict"].(float64) == 1 && IsNotLoli {
					illusbyte, err := network.Curl(config.PixivIllustsEnd+v2["id"].(string), nil)
					if err != nil {
						return err
					}

					err = json.Unmarshal(illusbyte, &Illusts)
					if err != nil {
						return err
					}

					Body := Illusts["body"].(map[string]interface{})
					Tags := Body["tags"].(map[string]interface{})
					Img := Body["urls"].(map[string]interface{})
					FixImg := Img["original"].(string)

					usrbyte, err := network.Curl(config.PixivUserEnd+Tags["authorId"].(string), nil)
					if err != nil {
						return err
					}

					err = json.Unmarshal(usrbyte, &User)
					if err != nil {
						return err
					}

					UserBody := User["body"].(map[string]interface{})

					Desc := RemoveHtmlTag(Body["description"].(string))
					if match, _ := regexp.MatchString("http://twitter.com", Desc); match {
						TextFix = ClearTwitterURL(Desc)
					} else {
						TextFix = "**" + Body["title"].(string) + "**\n" + Desc
					}

					FixFanArt.AddAuthor(v2["userName"].(string)).AddPermanentURL(BaseURL + v2["id"].(string)).
						AddAuthorAvatar(config.PixivProxy + UserBody["imageBig"].(string)).AddPhotos([]string{FixImg}).
						AddText(TextFix).AddPixivID(v2["id"].(string)).SetState(config.PixivArt)

					new, err := FixFanArt.CheckPixivFanArt()
					if err != nil {
						return err
					}
					if new {
						path, err := DownloadImg(Img["mini"].(string))
						if err != nil {
							log.Error(err)
						}

						Color, err := engine.GetColor("", path)
						if err != nil {
							return err
						}
						FixFanArt.Photos[0] = config.PixivProxy + FixImg
						notif.SendNude(*FixFanArt, Bot, Color)
					}
				}
				if i == Limit {
					break
				}
			}
		}
	}
	return nil
}

func GetPixivURL(str string) string {
	return "https://www.pixiv.net/ajax/search/artworks/" + str + "?word=" + str + "&order=date_d&mode=all&p=1&s_mode=s_tag&type=all&lang=en"
}

func GetPixivLewdURL(str string) string {
	return "https://www.pixiv.net/ajax/search/artworks/" + str + "?word=" + str + "&order=date_d&mode=r18&p=1&s_mode=s_tag&type=all&lang=en"
}

func ClearTwitterURL(str1 string) string {
	re := regexp.MustCompile(`(https\:\/\/twitter\.com\/.*)\<`)
	submatchall := re.FindStringSubmatch(str1)
	fix := ""
	for _, element := range submatchall {
		fix = element
	}
	return fix
}

func RemoveHtmlTag(in string) string {
	// regex to match html tag
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)
	// should replace long string first
	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})
	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}
	return in
}

func DownloadImg(u string) (string, error) {
	dir := config.TmpDir + engine.RanString()
	out, err := os.Create(dir)
	if err != nil {
		return "", err
	}

	defer out.Close()
	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", network.RandomAgent())
	request.Header.Set("Referer", "https://www.pixiv.net")

	spaceClient := http.Client{}
	resp, err := spaceClient.Do(request.WithContext(ctx))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return dir, nil
}
