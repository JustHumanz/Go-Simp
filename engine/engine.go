package engine

import (
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
	GCSDIR     string
	ImgDomain  string
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

	go BotSession.AddHandler(Fanart)
	go BotSession.AddHandler(Tags)
	go BotSession.AddHandler(Enable)
	go BotSession.AddHandler(Status)
	go BotSession.AddHandler(Help)
	go BotSession.AddHandler(Humanz)

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

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//make a http request
func Curl(url string, addheader []string) ([]byte, error) {
	var body []byte
	spaceClient := http.Client{
		Timeout: time.Second * 30, // Timeout after 30 seconds
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte(err.Error()), err
	}
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("User-Agent", RandomAgent())
	if addheader != nil {
		req.Header.Set(addheader[0], addheader[1])
	}

	res, err := spaceClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if res.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"Status": res.StatusCode,
			"Reason": res.Status,
			"URL":    url,
		}).Error("Status code not daijobu")
		return []byte{}, errors.New(res.Status)
	}

	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return body, err

	}

	return body, nil
}

//make a cooler http request *with multitor*
func CoolerCurl(urls string) ([]byte, error) {
	var counter int
	for {
		counter++
		proxyURL, err := url.Parse("http://multi_tor:16379")
		if err != nil {
			return []byte{}, err
		}

		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
			Timeout: 30 * time.Second,
		}

		request, err := http.NewRequest("GET", urls, nil)
		if err != nil {
			return []byte{}, err
		}

		response, err := client.Do(request)
		if err != nil && counter == 3 {
			return []byte{}, err
		}

		if response.StatusCode != http.StatusOK && counter == 3 {
			return []byte{}, errors.New("Tor get Status code " + strconv.Itoa(response.StatusCode))
		}

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return []byte{}, err
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
	DataGroup := database.GetGroup()
	for i := 0; i < len(DataGroup); i++ {
		DataMember := database.GetName(DataGroup[i].ID)
		for j := 0; j < len(DataMember); j++ {
			Name = strings.ToLower(Name)
			DataMember[j].Name = strings.ToLower(DataMember[j].Name)
			DataMember[j].JpName = strings.ToLower(DataMember[j].JpName)
			if Name == DataMember[j].Name || Name == DataMember[j].JpName {
				return Memberst{
					VTName:     FixName(DataMember[j].EnName, DataMember[j].JpName),
					ID:         DataMember[j].ID,
					YtChannel:  strings.Split(DataMember[j].YoutubeID, "\n"),
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

		body, curlerr = CoolerCurl(urls)
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
		"Mac OS X / Safari 12 [Desktop]: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4 Supplemental Update) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.1 Safari/605.1.15",
		"Windows / IE 11 [Desktop]: Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Windows / Edge 44 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763",
		"Windows / Chrome 83 [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36",
		"Windows / Firefox 60 ESR [Desktop]: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:60.8) Gecko/20100101 Firefox/60.8",
		"Android Phone / Firefox 77 [Mobile]: Mozilla/5.0 (Android 10; Mobile; rv:77.0) Gecko/77.0 Firefox/77.0",
		"Android Phone / Chrome 83 [Mobile]: Mozilla/5.0 (Linux; Android 10; Z832 Build/MMB29M) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Mobile Safari/537.36",
		"Android Tablet / Firefox 77 [Mobile]: Mozilla/5.0 (Android 10; Tablet; rv:77.0) Gecko/77.0 Firefox/77.0",
		"Android Tablet / Chrome 83 [Mobile]: Mozilla/5.0 (Linux; Android 10; SAMSUNG-SM-T377A Build/NMF26X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Mobile Safari/537.36",
		"iPhone / Safari 12.1.1 [Mobile]: Mozilla/5.0 (iPhone; CPU OS 10_15_4 Supplemental Update like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.1 Mobile/14E304 Safari/605.1.15",
		"iPad / Safari 12.1.1 [Mobile]: Mozilla/5.0 (iPad; CPU OS 10_15_4 Supplemental Update like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.1.1 Mobile/15E148 Safari/605.1.15",
		"Mozilla/5.0 (PlayStation 4 4.71) AppleWebKit/601.2 (KHTML, like Gecko)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.18",
		"Opera/9.80 (Linux armv7l) Presto/2.12.407 Version/12.51 , D50u-D1-UHD/V1.5.16-UHD (Vizio, D50u-D1, Wireless)",
		"Opera/9.80 (Windows NT 6.0) Presto/2.12.388 Version/12.14",
		"Opera/9.80 (Windows NT 5.1; WOW64) Presto/2.12.388 Version/12.17",
		"Opera/9.80 (Windows NT 5.1; U; ru) Presto/2.9.168 Version/11.50",
		"Opera/9.80 (Windows NT 6.1; WOW64) Presto/2.12.388 Version/12.16"}

	return Agent[rand.Intn(len(Agent))]
}
