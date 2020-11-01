package engine

import (
	"context"
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	"github.com/bwmarrin/discordgo"
	"github.com/cenkalti/dominantcolor"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	//BotID      *discordgo.User
	//BotSession *discordgo.Session
	//db         *sql.DB
	GroupData  []database.GroupName
	GroupsName []string
	RegList    = make(map[string]string)
	//H3llcome   = []string{config.Bonjour, config.Howdy, config.Guten, config.Koni, config.Selamat, config.Assalamu, config.Approaching}
	//PathLiteDB = "./engine/guild.db"
	//GuildList  []string
)

//Start module
func Start() {
	GroupData = database.GetGroup()
	for _, Group := range GroupData {
		GroupsName = append(GroupsName, Group.NameGroup)
		list := []string{}
		keys := make(map[string]bool)
		for _, Member := range database.GetName(Group.ID) {
			if _, value := keys[Member.Region]; !value {
				keys[Member.Region] = true
				list = append(list, Member.Region)
			}
		}
		RegList[Group.NameGroup] = strings.Join(list, ",")
	}
	log.Info("Engine module ready")
}

//BruhMoment Error hanlder
func BruhMoment(err error, msg string, exit bool) {
	if err != nil {
		log.Info(msg)
		log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

//GetYtToken Get a valid token
func GetYtToken() string {
	FreshToken := config.YtToken[0]
	for _, Token := range config.YtToken {
		_, err := Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id=UCfuz6xYbYFGsWWBi3SpJI1w&key="+Token, nil)
		if err == nil {
			FreshToken = Token
		}
	}
	return FreshToken
}

//Curl make a http request
func Curl(url string, addheader []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body []byte
	spaceClient := http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("cache-control", "no-cache")
	request.Header.Set("User-Agent", RandomAgent())
	if addheader != nil {
		request.Header.Set(addheader[0], addheader[1])
	}

	response, err := spaceClient.Do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"Status": response.StatusCode,
			"Reason": response.Status,
			"URL":    url,
		}).Error("Status code not daijobu")
		return nil, errors.New(response.Status)
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return body, err

	}

	return body, nil
}

//CoolerCurl make a cooler http request *with multitor*
func CoolerCurl(urls string, addheader []string) ([]byte, error) {
	counter := 0
	for {
		counter++
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		proxyURL, err := url.Parse("http://multi_tor:16379")
		if err != nil || counter == 3 {
			return nil, err
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
				DialContext: (&net.Dialer{
					Timeout: 10 * time.Second,
				}).DialContext,
			},
		}

		request, err := http.NewRequest("GET", urls, nil)
		if err != nil || counter == 3 {
			return nil, err
		}
		request.Header.Set("cache-control", "no-cache")
		request.Header.Set("User-Agent", RandomAgent())
		if addheader != nil {
			request.Header.Set(addheader[0], addheader[1])
		}
		response, err := client.Do(request.WithContext(ctx))
		if err != nil || counter == 3 {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK && counter == 3 {
			return nil, errors.New("Multi Tor get Error")
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil || counter == 3 {
			return nil, err
		}

		if data != nil {
			return data, nil
		}
	}
}

//FixName change to Title format
func FixName(A string, B string) string {
	if A != "" && B != "" {
		return "【" + strings.Title(strings.Join([]string{A, B}, "/")) + "】"
	} else if B != "" {
		return "【" + strings.Title(B) + "】"
	} else {
		return "【" + strings.Title(A) + "】"
	}
}

//RanString Random string for tmp file
func RanString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < 3; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

//SaucenaoCheck Check image from bilibili to saucenao *cuz to many reupload on bilibili*
func SaucenaoCheck(url string) (bool, []string, error) {
	var (
		data    Sauce
		body    []byte
		curlerr error
		urls    = "https://saucenao.com/search.php?db=999&output_type=2&numres=1&url=" + url + "&api_key=" + config.SauceAPI
	)
	body, curlerr = Curl(urls, nil)
	if curlerr != nil {
		log.Error(curlerr, string(body))
		log.Info("Trying use tor")

		body, curlerr = CoolerCurl(urls, nil)
		if curlerr != nil {
			log.Error(curlerr)
			return true, nil, curlerr
		}
	}
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Error(err)
		return true, nil, err
	}

	for _, res := range data.Results {
		simi, err := strconv.Atoi(res.Header.Similarity[:1])
		if err != nil {
			log.Error(err)
			return true, nil, err
		}
		if simi > 8 {
			return true, res.Data.ExtUrls, nil
		}
	}
	return false, nil, nil
}

//GetColor Get color from image
func GetColor(filepath, url string) (int, error) {
	def := 16770790

	if url == "" {
		return def, errors.New("urls nill ")
	}
	if url[len(url)-4:] == ".gif" {
		return def, nil
	}
	match, err := regexp.MatchString("http", url)
	if err != nil {
		return def, err
	}

	if match {
		filepath = filepath + RanString()
		out, err := os.Create(filepath)
		if err != nil {
			return def, err
		}

		defer out.Close()
		resp, err := http.Get(url)
		if err != nil {
			return def, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return def, errors.New("Server Error,status get " + resp.Status)
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return def, err
		}

	} else {
		filepath = url
	}
	f, err := os.Open(filepath)
	if err != nil {
		return def, err
	}

	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return def, err
	}

	Hex := dominantcolor.Hex(dominantcolor.Find(img))
	Hex = strings.Replace(Hex, "#", "0x", -1)
	Hex = strings.ToLower(Hex)
	Fix, err := strconv.ParseInt(Hex, 0, 64)
	if err != nil {
		return def, err
	}

	err = os.Remove(filepath)
	if err != nil {
		return def, err
	}

	return int(Fix), nil
}

//RandomAgent Create random useragent for bypass some IDS
func RandomAgent() string {
	Agent := []string{"Windows / Firefox 77 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Linux / Firefox 77 [Desktop]: Mozilla/5.0 (X11; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"}
	return Agent[rand.Intn(len(Agent))]
}

func Reacting(Data map[string]string, s *discordgo.Session) error {
	EmojiList := config.EmojiFanart
	ChannelID := Data["ChannelID"]
	MessID, err := s.Channel(ChannelID)
	if err != nil {
		return errors.New(err.Error() + " ChannelID: " + ChannelID)
	}
	for l := 0; l < len(EmojiList); l++ {
		if Data["Content"][len(Data["Prefix"]):] == "kanochi" {
			err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[0])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}
			break
		} else if Data["Content"][len(Data["Prefix"]):] == "cleaire" {
			err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}
			if l == len(EmojiList)-1 {
				err = s.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":latom:767810745860751391")
				if err != nil {
					return errors.New(err.Error() + " ChannelID: " + ChannelID)
					//log.Error(err, ChannelID)
				}
			}
		} else if Data["Content"][len(Data["Prefix"]):] == "senchou" {
			err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
			}

			if l == len(EmojiList)-1 {
				err = s.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":hormny:768700671750176790")
				if err != nil {
					return errors.New(err.Error() + " ChannelID: " + ChannelID)
					//log.Error(err, ChannelID)
				}
			}
		} else {
			err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
			if err != nil {
				return errors.New(err.Error() + " ChannelID: " + ChannelID)
				//log.Error(err, ChannelID)
				//break
			}
		}
	}
	return nil
}

func Zawarudo(Region string) *time.Location {
	if Region == "ID" {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		return loc
	} else if Region == "JP" {
		loc, _ := time.LoadLocation("Asia/Tokyo")
		return loc
	} else if Region == "CN" {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		return loc
	} else if Region == "KR" {
		loc, _ := time.LoadLocation("Asia/Seoul")
		return loc
	} else {
		loc, _ := time.LoadLocation("UTC")
		return loc
	}
}
