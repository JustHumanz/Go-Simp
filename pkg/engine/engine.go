package engine

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	"github.com/bwmarrin/discordgo"
	"github.com/cenkalti/dominantcolor"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/nicklaw5/helix"
	"github.com/pbnjay/memory"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

//GetYtToken Get a valid token
func GetYtToken() *string {
	for _, Token := range config.GoSimpConf.YtToken {
		body, err := network.Curl("https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=GNkPJvVEm0s&key="+Token, nil)
		if err == nil && body != nil {
			return &Token
		}
	}
	log.Error("Youtube Token out of limit")
	return nil
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

//GetColor Get color from image
func GetColor(filepath, url string) (int, error) {
	def := 16770790

	if url == "" {
		return def, errors.New("urls img nill")
	}
	if url[len(url)-4:] == ".gif" {
		return def, nil
	}
	match, err := regexp.MatchString("http", url)
	if err != nil {
		return def, err
	}

	if match {
		filepath += RanString()
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

func Reacting(Data map[string]string, s *discordgo.Session) error {
	ChannelID := Data["ChannelID"]
	if Data["State"] == "Youtube" {
		EmojiList := config.GoSimpConf.Emoji.Livestream
		for _, Emoji := range EmojiList {
			err := s.MessageReactionAdd(ChannelID, Data["MessageID"], Emoji)
			if err != nil {
				return err
			}
		}
	} else if Data["State"] == "SelectType" {
		EmojiList := []string{config.Ok, config.No}
		for _, Emoji := range EmojiList {
			err := s.MessageReactionAdd(ChannelID, Data["MessageID"], Emoji)
			if err != nil {
				return err
			}
		}
	} else if Data["State"] == "Menu" {
		EmojiList := []string{config.One, config.Two, config.Three}
		for _, Emoji := range EmojiList {
			err := s.MessageReactionAdd(ChannelID, Data["MessageID"], Emoji)
			if err != nil {
				return err
			}
		}
	} else if Data["State"] == "Menu2" {
		EmojiList := []string{config.One, config.Two, config.Three, config.Four}
		for _, Emoji := range EmojiList {
			err := s.MessageReactionAdd(ChannelID, Data["MessageID"], Emoji)
			if err != nil {
				return err
			}
		}
	} else if Data["State"] == "TypeChannel" {
		EmojiList := []string{config.Art, config.Live, config.Lewd}
		for _, Emoji := range EmojiList {
			err := s.MessageReactionAdd(ChannelID, Data["MessageID"], Emoji)
			if err != nil {
				return err
			}
		}
	} else {
		MessID, err := s.Channel(ChannelID)
		if err != nil {
			return err
		}
		EmojiList := config.GoSimpConf.Emoji.Fanart
		for l := 0; l < len(EmojiList); l++ {
			if Data["Content"][len(Data["Prefix"]):] == "kanochi" {
				err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[0])
				if err != nil {
					return err
					//log.Error(err, ChannelID)
				}
				break
			} else if Data["Content"][len(Data["Prefix"]):] == "cleaire" {
				err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
				if err != nil {
					return err
					//log.Error(err, ChannelID)
				}
				if l == len(EmojiList)-1 {
					err = s.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":latom:767810745860751391")
					if err != nil {
						return err
						//log.Error(err, ChannelID)
					}
				}
			} else if Data["Content"][len(Data["Prefix"]):] == "senchou" {
				err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
				if err != nil {
					return err
					//log.Error(err, ChannelID)
				}

				if l == len(EmojiList)-1 {
					err = s.MessageReactionAdd(ChannelID, MessID.LastMessageID, ":hormny:768700671750176790")
					if err != nil {
						return err
						//log.Error(err, ChannelID)
					}
				}
			} else {
				err := s.MessageReactionAdd(ChannelID, MessID.LastMessageID, EmojiList[l])
				if err != nil {
					return err
					//log.Error(err, ChannelID)
					//break
				}
			}
		}
	}
	return nil
}

func Zawarudo(Region string) *time.Location {
	Default := func() *time.Location {
		loc, _ := time.LoadLocation("UTC")
		return loc
	}

	if Region == "ID" {
		loc, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "JP" {
		loc, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "CN" {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "KR" {
		loc, err := time.LoadLocation("Asia/Seoul")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "MY" {
		loc, err := time.LoadLocation("Asia/Kuala_Lumpur")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "IN" {
		loc, err := time.LoadLocation("Asia/Dhaka")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "PH" {
		loc, err := time.LoadLocation("Asia/Manila")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "AU" {
		loc, err := time.LoadLocation("Australia/Sydney")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else if Region == "FI" {
		loc, err := time.LoadLocation("Europe/Helsinki")
		if err != nil {
			log.Error(err)
			return Default()
		}
		return loc
	} else {
		return Default()
	}
}

func CountryCodetoUniCode(Region string) string {
	if Region == "ID" {
		return "🇮🇩"
	} else if Region == "JP" {
		return "🇯🇵"
	} else if Region == "CN" {
		return "🇨🇳"
	} else if Region == "KR" {
		return "🇰🇷"
	} else if Region == "MY" {
		return "🇲🇾"
	} else if Region == "IN" {
		return "🇮🇳"
	} else if Region == "PH" {
		return "🇵🇭"
	} else if Region == "AU" {
		return "🇦🇺"
	} else if Region == "US" {
		return "🇺🇸"
	} else if Region == "FL" {
		return "🇫🇮"
	} else if Region == "EN" {
		return "🇪🇺"
	} else if Region == "UK" {
		return "🇬🇧"
	}
	return ""
}

func UniCodetoCountryCode(Region string) string {
	if Region == "🇮🇩" {
		return "ID"
	} else if Region == "🇯🇵" {
		return "JP"
	} else if Region == "🇨🇳" {
		return "CN"
	} else if Region == "🇰🇷" {
		return "KR"
	} else if Region == "🇲🇾" {
		return "MY"
	} else if Region == "🇪🇺" {
		return "EN"
	} else if Region == "🇮🇳" {
		return "IN"
	} else if Region == "🇵🇭" {
		return "PH"
	} else if Region == "🇦🇺" {
		return "AU"
	} else if Region == "🇺🇸" {
		return "US"
	} else if Region == "🇫🇮" {
		return "FI"
	} else if Region == "🇬🇧" {
		return "UK"
	}
	return ""
}

func YtFindType(title string) string {
	yttype := ""
	title = strings.ToLower(title)
	if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|mv|covered)", title); Cover {
		yttype = "Covering"
	} else if Chat, _ := regexp.MatchString("(?m)(chat|room)", title); Chat {
		yttype = "ChatRoom"
	} else if Singing, _ := regexp.MatchString("(?m)(sing|歌枠)", title); Singing {
		yttype = "Singing"
	} else {
		yttype = "Streaming"
	}
	return yttype
}

//GetAuthorAvatar Get twitter avatar
func GetAuthorAvatar(username string) string {
	profile, err := InitTwitterScraper().GetProfile(username)
	if err != nil {
		log.Error(err)
	}
	return strings.Replace(profile.Avatar, "normal.jpg", "400x400.jpg", -1)
}

//GetTwitterFollow Scrapping twitter followers
func GetTwitterFollow(UserName string) (twitterscraper.Profile, error) {
	profile, err := InitTwitterScraper().GetProfile(UserName)
	if err != nil {
		return twitterscraper.Profile{}, err
	}
	return profile, nil
}

func RoundPrec(x float64, prec int) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow * sign
}

func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
	if math.IsNaN(number) || math.IsInf(number, 0) {
		number = 0
	}

	var ret string
	var negative bool

	if number < 0 {
		number *= -1
		negative = true
	}

	d, fract := math.Modf(number)

	if decimals <= 0 {
		fract = 0
	} else {
		pow := math.Pow(10, float64(decimals))
		fract = RoundPrec(fract*pow, 0)
	}

	if thousandsSep == "" {
		ret = strconv.FormatFloat(d, 'f', 0, 64)
	} else if d >= 1 {
		var x float64
		for d >= 1 {
			d, x = math.Modf(d / 1000)
			x = x * 1000
			ret = strconv.FormatFloat(x, 'f', 0, 64) + ret
			if d >= 1 {
				ret = thousandsSep + ret
			}
		}
	} else {
		ret = "0"
	}

	fracts := strconv.FormatFloat(fract, 'f', 0, 64)

	// "0" pad left
	for i := len(fracts); i < decimals; i++ {
		fracts = "0" + fracts
	}

	ret += decPoint + fracts

	if negative {
		ret = "-" + ret
	}
	return ret
}

func RoundInt(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

func FormatNumber(input float64) string {
	x := RoundInt(input)
	xFormatted := NumberFormat(float64(x), 2, ".", ",")
	return xFormatted
}

func NearestThousandFormat(num float64) string {

	if math.Abs(num) < 999.5 {
		xNum := FormatNumber(num)
		xNumStr := xNum[:len(xNum)-3]
		return string(xNumStr)
	}

	xNum := FormatNumber(num)
	// first, remove the .00 then convert to slice
	xNumStr := xNum[:len(xNum)-3]
	xNumCleaned := strings.Replace(xNumStr, ",", " ", -1)
	xNumSlice := strings.Fields(xNumCleaned)
	count := len(xNumSlice) - 2
	unit := [4]string{"K", "M", "B", "T"}
	xPart := unit[count]

	afterDecimal := ""
	if xNumSlice[1][0] != 0 {
		afterDecimal = "." + string(xNumSlice[1][0])
	}
	final := xNumSlice[0] + afterDecimal + xPart
	return final
}

func RemoveEmbed(VideoID string, Bot *discordgo.Session) {
	ChannelState, err := database.GetLiveNotifMsg(VideoID)
	if err != nil {
		log.Error(err)
	}
	for _, v := range ChannelState {
		log.WithFields(log.Fields{
			"VideoID": VideoID,
			"Status":  "Past",
		}).Info("Delete message from ", []string{v.TextMessageID, v.EmbedMessageID})
		err := Bot.ChannelMessagesBulkDelete(v.ChannelID, []string{v.TextMessageID, v.EmbedMessageID})
		if err != nil {
			log.Error(err)
		}
	}
}

func UnderScoreName(s string) string {
	return strings.ReplaceAll(s, " ", "_")
}

func MaintenanceIMG() string {
	return config.CdnDomain + "Maintenance/" + strconv.Itoa(RandomNum(1, 8)) + ".png"
}

func NotFoundIMG() string {
	return config.CdnDomain + "Command_Not_Found/" + strconv.Itoa(RandomNum(1, 7)) + ".png"
}

func LewdIMG() string {
	return config.CdnDomain + "Lewd/" + strconv.Itoa(RandomNum(1, 5)) + ".png"
}

func Gif() string {
	return config.CdnDomain + "Gif/" + strconv.Itoa(RandomNum(1, 4)) + ".gif"
}

func RandomNum(min, max int) int {
	return rand.Intn(max-min) + min
}

//RemoveTwitterShortLink remove twitter shotlink
func RemoveTwitterShortLink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}

//GetRSS GetRSS from Channel
func GetRSS(YtID string, proxy bool) ([]string, error) {
	var (
		DataXML YtXML
		VideoID []string
		URL     = "https://www.youtube.com/feeds/videos.xml?channel_id=" + YtID + "&q=searchterms"
		Data    []byte
		err     error
	)

	if proxy {
		Data, err = network.CoolerCurl(URL, nil)
	} else {
		Data, err = network.Curl(URL, nil)
	}

	if err != nil {
		log.Error(err, string(Data))
		return nil, err
	}

	xml.Unmarshal(Data, &DataXML)

	for i := 0; i < len(DataXML.Entry); i++ {
		VideoID = append(VideoID, DataXML.Entry[i].VideoId)
		if i == config.GoSimpConf.LimitConf.YoutubeLimit {
			break
		}
	}
	return VideoID, nil
}

var ExTknList []string

//YtAPI Get data from youtube api
func YtAPI(VideoID []string) (YtData, error) {
	var (
		Data YtData
	)
	log.WithFields(log.Fields{
		"VideoID": VideoID,
	}).Info("Checking from youtubeAPI")

	for i, Token := range config.GoSimpConf.YtToken {
		if ExTknList != nil {
			isExhaustion := false
			for _, v := range ExTknList {
				if v == Token {
					isExhaustion = true
					break
				}
			}

			if isExhaustion {
				continue
			}
		}
		url := "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(snippet(publishedAt,title,description,thumbnails(standard),channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=" + strings.Join(VideoID, ",") + "&key=" + Token

		bdy, curlerr := network.Curl(url, nil)
		if curlerr != nil {
			log.Error(curlerr)
			if curlerr.Error() == "403 Forbidden" {
				ExTknList = append(ExTknList, Token)
			} else {
				time.Sleep(10 * time.Second)
			}

			if i == len(config.GoSimpConf.YtToken)-1 {
				break
			}
			continue
		}

		err := json.Unmarshal(bdy, &Data)
		if err != nil {
			return Data, err
		}
		return Data, nil

	}
	return YtData{}, errors.New("exhaustion Token")
}

//GetWaiting get viwers by scraping yt video
func GetWaiting(VideoID string) (string, error) {
	var (
		bit     []byte
		curlerr error
		urls    = "https://www.youtube.com/watch?v=" + VideoID
	)
	bit, curlerr = network.Curl(urls, nil)
	if curlerr != nil || bit == nil {
		bit, curlerr = network.CoolerCurl(urls, nil)
		if curlerr != nil {
			return config.Ytwaiting, curlerr
		}
	}
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return config.Ytwaiting, err
	}
	for _, element := range regexp.MustCompile(`(?m)videoViewCountRenderer.*?text([0-9\s]+).+(isLive\strue)`).FindAllStringSubmatch(reg.ReplaceAllString(string(bit), " "), -1) {
		tmp := strings.Replace(element[1], " ", "", -1)
		if tmp != "" {
			config.Ytwaiting = tmp
		}
	}
	return config.Ytwaiting, nil
}

//ParseDuration Parse video duration
func ParseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`P(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?`)
	matches := durationRegex.FindStringSubmatch(str)

	years := ParseInt64(matches[1])
	months := ParseInt64(matches[2])
	days := ParseInt64(matches[3])
	hours := ParseInt64(matches[4])
	minutes := ParseInt64(matches[5])
	seconds := ParseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func ParseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

func GetMaxSqlConn() int {
	a := int((int64(memory.FreeMemory()/1024/1024) * 1024 * 1024) / 12582880)
	//Avoid discord Rate Limit
	if a > 50 {
		return 50
	}
	return a
}

//Init twitter scrapper with or without proxy
func InitTwitterScraper() *twitterscraper.Scraper {
	Scraper := twitterscraper.New()
	Scraper.SetSearchMode(twitterscraper.SearchLatest)
	if config.GoSimpConf.MultiTOR != "" {
		err := Scraper.SetProxy("socks5://" + config.GoSimpConf.MultiTOR)
		if err != nil {
			log.Error(err)
		}
	}
	return Scraper
}

//Start bot
func StartBot(torTransport bool) *discordgo.Session {
	i := config.GoSimpConf
	tmp, err := discordgo.New("Bot " + i.Discord)
	if err != nil {
		log.Panic(err)
	}
	if torTransport {
		dialSocksProxy, err := proxy.SOCKS5("tcp", config.GoSimpConf.MultiTOR, nil, proxy.Direct)
		if err != nil {
			log.Error(err)
		}
		tmp.Client.Transport = &http.Transport{Dial: dialSocksProxy.Dial}
	}
	return tmp
}

//Get twitch token
func GetTwitchTkn() *helix.Client {
	i := config.GoSimpConf
	TwitchClient, err := helix.NewClient(&helix.Options{
		ClientID:     i.Twitch.ClientID,
		ClientSecret: i.Twitch.ClientSecret,
	})
	if err != nil {
		log.Panic(err)
	}
	return TwitchClient
}
