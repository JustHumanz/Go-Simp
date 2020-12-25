package engine

import (
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	network "github.com/JustHumanz/Go-simp/tools/network"
	"github.com/bwmarrin/discordgo"
	"github.com/cenkalti/dominantcolor"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	//BotID      *discordgo.User
	//BotSession *discordgo.Session
	//db         *sql.DB
	GroupData  []database.Group
	GroupsName []string
	RegList    = make(map[string]string)
	//PathLiteDB = "./engine/guild.db"
	//GuildList  []string
)

//Start module
func Start() {
	GroupData = database.GetGroups()
	for _, Group := range GroupData {
		GroupsName = append(GroupsName, Group.GroupName)
		list := []string{}
		keys := make(map[string]bool)
		for _, Member := range database.GetMembers(Group.ID) {
			if _, value := keys[Member.Region]; !value {
				keys[Member.Region] = true
				list = append(list, Member.Region)
			}
		}
		RegList[Group.GroupName] = strings.Join(list, ",")
	}
	log.Info("Engine module ready")
}

//GetYtToken Get a valid token
func GetYtToken() string {
	FreshToken := config.YtToken[0]
	for _, Token := range config.YtToken {
		_, err := network.Curl("https://www.googleapis.com/youtube/v3/channels?part=statistics&id=UCfuz6xYbYFGsWWBi3SpJI1w&key="+Token, nil)
		if err == nil {
			FreshToken = Token
		}
	}
	return FreshToken
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
	body, curlerr = network.Curl(urls, nil)
	if curlerr != nil {
		log.Error(curlerr, string(body))
		log.Info("Trying use tor")

		body, curlerr = network.CoolerCurl(urls, nil)
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
	scraper := twitterscraper.New()
	scraper.SetProxy(config.MultiTOR)
	profile, err := scraper.GetProfile(username)
	if err != nil {
		log.Error(err)
	}
	return strings.Replace(profile.Avatar, "normal.jpg", "400x400.jpg", -1)
}
