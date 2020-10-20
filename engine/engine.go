package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
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
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/bwmarrin/discordgo"
	"github.com/cenkalti/dominantcolor"
	log "github.com/sirupsen/logrus"
)

var (
	BotID      string
	BotSession *discordgo.Session
	db         *sql.DB
	debug      bool
	GroupData  []database.GroupName
	GroupsName []string
	GCSDIR     string
	ImgDomain  string
	RegList    = make(map[string]string)
)

//Start module
func Start(b *discordgo.Session, m string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	BotSession = b
	db = database.DB
	BotID = m
	if strings.ToLower(config.Logging) == "debug" {
		debug = true
	} else {
		debug = false
	}
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

	go BotSession.AddHandler(Fanart)
	go BotSession.AddHandler(Tags)
	go BotSession.AddHandler(Enable)
	go BotSession.AddHandler(Status)
	go BotSession.AddHandler(Help)
	go BotSession.AddHandler(BiliBiliMessage)
	go BotSession.AddHandler(BiliBiliSpace)
	go BotSession.AddHandler(YoutubeMessage)
	go BotSession.AddHandler(SubsMessage)
	//go BotSession.AddHandler(Humanz)

	log.Info("Engine module ready")
}
func Debugging(a, b, c interface{}) {
	if debug {
		f, err := os.OpenFile("debug.log", os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println(err)
		}
		log.SetOutput(f)
		log.WithFields(log.Fields{
			"Func":   a,
			"Status": b,
			"Value":  c,
		}).Debug(c)
	}
}

//Bruh moment
func BruhMoment(err error, msg string, exit bool) {
	if err != nil {
		log.Info(msg)
		log.Error(err)
		if exit {
			os.Exit(1)
		}
	}
}

func GetYtToken() string {
	FreshToken := config.YtToken[0]
	for _, Token := range config.YtToken {
		body, err := Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id=UCfuz6xYbYFGsWWBi3SpJI1w&key="+Token, nil)
		if err == nil || body != nil {
			FreshToken = Token
		}
	}
	return FreshToken
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//make a http request
func Curl(url string, addheader []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body []byte
	spaceClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("User-Agent", RandomAgent())
	if addheader != nil {
		req.Header.Set(addheader[0], addheader[1])
	}

	res, err := spaceClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"Status": res.StatusCode,
			"Reason": res.Status,
			"URL":    url,
		}).Error("Status code not daijobu")
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return body, err

	}

	return body, nil
}

//make a cooler http request *with multitor*
func CoolerCurl(urls string, addheader []string) ([]byte, error) {
	var counter int
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for {
		counter++
		proxyURL, err := url.Parse("http://multi_tor:16379")
		if err != nil && counter == 2 {
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
		if err != nil && counter == 2 {
			return nil, err
		}
		request.Header.Set("cache-control", "no-cache")
		request.Header.Set("User-Agent", RandomAgent())
		if addheader != nil {
			request.Header.Set(addheader[0], addheader[1])
		}
		response, err := client.Do(request.WithContext(ctx))
		if err != nil && counter == 2 {
			return nil, err
		}

		if response.StatusCode != http.StatusOK && counter == 2 {
			return nil, errors.New("Tor get Status code " + response.Status)
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil && counter == 2 {
			return nil, err
		}
		return data, nil
	}

}

//change to Title format
func FixName(A string, B string) string {
	funcvar := GetFunctionName(FixName)
	Debugging(funcvar, "In", fmt.Sprint(A, B))
	if A != "" && B != "" {
		P := []string{A, B}
		VName := strings.Join(P, "/")
		return strings.Title(VName)

	} else if B != "" {
		return strings.Title(B)
	} else {
		return strings.Title(A)
	}
}

//Find a valid name
func ValidName(Name string) Memberst {
	for i := 0; i < len(GroupData); i++ {
		DataMember := database.GetName(GroupData[i].ID)
		for j := 0; j < len(DataMember); j++ {
			Name = strings.ToLower(Name)
			DataMember[j].Name = strings.ToLower(DataMember[j].Name)
			DataMember[j].JpName = strings.ToLower(DataMember[j].JpName)
			if Name == DataMember[j].Name || Name == DataMember[j].JpName {
				return Memberst{
					VTName:     FixName(DataMember[j].EnName, DataMember[j].JpName),
					ID:         DataMember[j].ID,
					YtChannel:  DataMember[j].YoutubeID,
					SpaceID:    DataMember[j].BiliBiliID,
					BiliAvatar: DataMember[j].BiliBiliAvatar,
					YtAvatar:   DataMember[j].YoutubeAvatar,
				}
			}
		}
	}
	return Memberst{}
}

//Random string for tmp file
func RanString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 3
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}

//Check image from bilibili to saucenao *cuz to many reupload on bilibili*
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
		} else {
			return false, res.Data.ExtUrls, nil
		}
	}
	return false, nil, nil
}

//Get color from image
func GetColor(filepath, url string) (int, error) {
	if url == "" {
		return 0, errors.New("Url nill")
	}
	ex := url[len(url)-4:]

	if ex == ".gif" {
		return 16770790, nil
	}
	match, err := regexp.MatchString("http", url)
	if err != nil {
		log.Error(err)
		return 0, nil
	}

	if match {
		filepath = filepath + RanString()
		out, err := os.Create(filepath)
		if err != nil {
			log.Error(err)
			return 0, err
		}

		defer out.Close()
		resp, err := http.Get(url)
		if err != nil {
			log.Error(err)
			return 0, err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error("Server Error")
			return 0, nil
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Error(err)
			return 0, err
		}

	} else {
		filepath = url
	}
	f, err := os.Open(filepath)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	Hex := dominantcolor.Hex(dominantcolor.Find(img))
	Hex = strings.Replace(Hex, "#", "0x", -1)
	Hex = strings.ToLower(Hex)
	Fix, err := strconv.ParseInt(Hex, 0, 64)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	err = os.Remove(filepath)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return int(Fix), nil
}

//Create random useragent for bypass some IDS
func RandomAgent() string {
	Agent := []string{"Windows / Firefox 77 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Linux / Firefox 77 [Desktop]: Mozilla/5.0 (X11; Linux x86_64; rv:77.0) Gecko/20100101 Firefox/77.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
		"Windows / IE 11 [Desktop]: Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Windows / Edge 44 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
		"Windows / Chrome 83 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36",
		"Windows / Firefox 60 ESR [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.8) Gecko/20100101 Firefox/60.8",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.157 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:75.0) Gecko/20100101 Firefox/75.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:80.0) Gecko/20100101 Firefox/80.0",
		"Mozilla/5.0 (PlayStation 4 4.71) AppleWebKit/601.2 (KHTML, like Gecko)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8"}
	return Agent[rand.Intn(len(Agent))]
}
