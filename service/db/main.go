package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"
	bilibili "github.com/JustHumanz/Go-simp/livestream/bilibili/live"
	youtube "github.com/JustHumanz/Go-simp/livestream/youtube"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	res          Vtuber
	Limit        int
	db           *sql.DB
	YtToken      string
	member       string
	Publish      time.Time
	Roomstatus   string
	BiliSession  string
	New          NewVtuber
	DiscordToken string
)

type NewVtuber struct {
	Member Member
	Group  database.GroupName
}

func init() {
	fmt.Println("Reading hashtag file...")
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	file, err := ioutil.ReadFile("./vtuber.json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(string(file))
	err = json.Unmarshal(file, &res)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	YtToken = os.Getenv("GTOKEN")
	BiliSession = os.Getenv("SBILI")
	DiscordToken = os.Getenv("DISCORD")
	Limit = 100

	err = CreateDB()
	if err != nil {
		log.Error(err)
	}
}

func main() {
	Service := flag.String("service", "bootstrapping", "select service mode[bootstrapping/twitter_scrap]")
	ScrapMember := flag.Bool("vtuber", false, "enable this if you want to scrap tweet(fanart) each member")
	flag.StringVar(&member, "member", "kano", "list of vtuber name (split by space)")
	flag.Parse()
	db = DBConn()

	if (*Service) == "bootstrapping" {
		AddData(res)
		go CheckYT()
		go CheckSchedule()
		go CheckVideoSpace()
		go CheckTBili()
		go youtube.CheckPrivate()
		go func() {
			go Tweet("Independen", 0, Limit)
			for i := 0; i < len(res.Vtuber.Group); i++ {
				Tweet(res.Vtuber.Group[i].GroupName, 0, Limit)
			}
		}()
		time.Sleep(5 * time.Minute)
		log.Info("Done")
		os.Exit(0)
	} else if (*Service) == "twitter_scrap" {
		Limit = 10000000
		if *ScrapMember {
			if len(flag.Args()) > 0 {
				for i := 0; i < len(flag.Args()); i++ {
					Data := engine.FindName(flag.Args()[i])
					Tweet(Data.GroupName, Data.MemberID, Limit)
				}
				log.Info("Done")
				os.Exit(0)
			} else {
				log.Error("No Vtuber Name found")
				os.Exit(1)
			}
		} else {
			for i := 0; i < len(res.Vtuber.Group); i++ {
				Tweet(res.Vtuber.Group[i].GroupName, 0, Limit)
			}
			log.Info("Done")
			os.Exit(0)
		}
	} else {
		AddData(res)
		//for _, NewData := range New {
		//	NewData.SendNotif(Bot)
		// }
		//Tweet("Independen", 0, Limit)
		log.Info("Done")
		os.Exit(0)
	}
}

func CheckYT() {
	Data := database.GetGroup()
	for i := 0; i < len(Data); i++ {
		Names := database.GetName(Data[i].ID)
		for _, Name := range Names {
			log.WithFields(log.Fields{
				"Vtube":        Name.EnName,
				"Youtube ID":   Name.YoutubeID,
				"Vtube Region": Name.Region,
			}).Info("Checking yt")
			FilterYt(Name)
		}
	}
}

func CheckTBili() {
	DataGroup := database.GetGroup()
	for k := 0; k < len(DataGroup); k++ {
		DataMember := database.GetName(DataGroup[k].ID)
		for z := 0; z < len(DataMember); z++ {
			if DataMember[z].BiliBiliHashtags != "" {
				log.WithFields(log.Fields{
					"Group":  DataGroup[k].NameGroup,
					"Vtuber": DataMember[z].EnName,
				}).Info("Start crawler T.bilibili")
				body, err := engine.Curl("https://api.vc.bilibili.com/topic_svr/v1/topic_svr/topic_new?topic_name="+url.QueryEscape(DataMember[z].BiliBiliHashtags), nil)
				if err != nil {
					log.Error(err)
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
							niggerlist := []string{"解锁专属粉丝卡片", "Official", "twitter.com", "咖啡厅", "CD", "专辑", "PIXIV", "遇", "marshmallow-qa.com"}
							for _, Nworld := range niggerlist {
								nope, _ = regexp.MatchString(Nworld, STB.Item.Description)
								if nope {
									break
								}
							}
							New := database.GetTBiliBili(TB.Data.Cards[i].Desc.DynamicIDStr)

							if New && !nope {
								log.WithFields(log.Fields{
									"Group":  DataGroup[k].NameGroup,
									"Vtuber": DataMember[z].EnName,
								}).Info("New Fanart")
								for l := 0; l < len(STB.Item.Pictures); l++ {
									img = append(img, STB.Item.Pictures[l].ImgSrc)
								}

								Data := database.InputTBiliBili{
									URL:        "https://t.bilibili.com/" + TB.Data.Cards[i].Desc.DynamicIDStr + "?tab=2",
									Author:     TB.Data.Cards[i].Desc.UserProfile.Info.Uname,
									Avatar:     TB.Data.Cards[i].Desc.UserProfile.Info.Face,
									Like:       TB.Data.Cards[i].Desc.Like,
									Photos:     strings.Join(img, "\n"),
									Dynamic_id: TB.Data.Cards[i].Desc.DynamicIDStr,
									Text:       STB.Item.Description,
								}
								log.Info("Send to database")
								Data.InputTBiliBili(DataMember[z].ID)
							} else {
								log.WithFields(log.Fields{
									"Group":  DataGroup[k].NameGroup,
									"Vtuber": DataMember[z].EnName,
								}).Info("Still same")
							}
						}
					}
				} else {
					log.WithFields(log.Fields{
						"Group":  DataGroup[k].NameGroup,
						"Vtuber": DataMember[z].EnName,
					}).Info("Still same")
				}
				time.Sleep(time.Duration(int64(rand.Intn((7-1)+1))) * time.Second)
			}
		}
	}
}

func CheckSchedule() {
	log.Info("Start check Schedule")
	Group := database.GetGroup()
	for z := 0; z < len(Group); z++ {
		Name := database.GetName(Group[z].ID)
		for k := 0; k < len(Name); k++ {
			if Name[k].BiliBiliID != 0 {
				log.WithFields(log.Fields{
					"Group":   Group[z].NameGroup,
					"SpaceID": Name[k].EnName,
				}).Info("Check Room")
				var (
					ScheduledStart time.Time
				)
				DataDB := database.GetRoomData(Name[k].ID, Name[k].BiliRoomID)
				Status := bilibili.GetRoomStatus(Name[k].BiliRoomID)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				if Status.Data.RoomInfo.LiveStartTime != 0 {
					ScheduledStart = time.Unix(int64(Status.Data.RoomInfo.LiveStartTime), 0).In(loc)
				} else {
					ScheduledStart = time.Time{}
				}
				Data := map[string]interface{}{
					"LiveRoomID":     Name[k].BiliRoomID,
					"Status":         "",
					"Title":          Status.Data.RoomInfo.Title,
					"Thumbnail":      Status.Data.RoomInfo.Cover,
					"Description":    Status.Data.NewsInfo.Content,
					"PublishedAt":    time.Time{},
					"ScheduledStart": ScheduledStart,
					"Face":           Status.Data.AnchorInfo.BaseInfo.Face,
					"Online":         Status.Data.RoomInfo.Online,
					"BiliBiliID":     Name[k].BiliBiliID,
					"MemberID":       Name[k].ID,
				}
				if Status.CheckScheduleLive() {
					//Live
					Data["Status"] = "Live"
					LiveBiliBili(Data)
				} else if !Status.CheckScheduleLive() && DataDB.Status == "Live" {
					//prob past
					Data["Status"] = "Past"
					LiveBiliBili(Data)
				} else if DataDB == (database.LiveBiliDB{}) {
					Data["Status"] = "Unknown"
					LiveBiliBili(Data)
				}
			}
		}
	}
}

func CheckVideoSpace() {
	Group := database.GetGroup()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for z := 0; z < len(Group); z++ {
		Name := database.GetName(Group[z].ID)
		for k := 0; k < len(Name); k++ {
			if Name[k].BiliBiliID != 0 {
				log.WithFields(log.Fields{
					"Group":   Group[z].NameGroup,
					"SpaceID": Name[k].EnName,
				}).Info("Check Space")
				var (
					PushVideo SpaceVideo
					videotype string
					url       []string
				)
				baseurl := "https://api.bilibili.com/x/space/arc/search?mid=" + strconv.Itoa(Name[k].BiliBiliID) + "&ps=100"
				url = []string{baseurl + "&tid=1", baseurl + "&tid=3", baseurl + "&tid=4"}
				for f := 0; f < len(url); f++ {
					body, err := engine.Curl(url[f], nil)
					if err != nil {
						log.Error(err, string(body))
					}
					var tmp SpaceVideo
					err = json.Unmarshal(body, &tmp)
					if err != nil {
						log.Error(err)
					}
					for _, Vlist := range tmp.Data.List.Vlist {
						PushVideo.Data.List.Vlist = append(PushVideo.Data.List.Vlist, Vlist)
					}
				}

				for _, video := range PushVideo.Data.List.Vlist {
					if Cover, _ := regexp.MatchString("(?m)(cover|song|feat|music|翻唱|mv)", strings.ToLower(video.Title)); Cover {
						videotype = "Covering"
					} else {
						videotype = "Streaming"
					}
					database.InputSpaceVideo(database.InputBiliBili{
						VideoID:  video.Bvid,
						Type:     videotype,
						Title:    video.Title,
						Thum:     "https:" + video.Pic,
						Desc:     video.Description,
						Update:   time.Unix(int64(video.Created), 0).In(loc),
						Viewers:  video.Play,
						MemberID: Name[k].ID,
					})
				}
			}
		}
	}
}

func (Data NewVtuber) SendNotif(Bot *discordgo.Session) {
	Avatar := Data.Member.YtAvatar()
	Color, err := engine.GetColor("/tmp/notf.gg", Avatar)
	if err != nil {
		log.Error(err)
	}
	for _, Channel := range GetChannelByGroup(Data.Group.ID) {
		msg, err := Bot.ChannelMessageSendEmbed(Channel, engine.NewEmbed().
			SetAuthor(Data.Group.NameGroup, Data.Group.IconURL).
			SetTitle(engine.FixName(Data.Member.ENName, Data.Member.JPName)).
			SetImage(Avatar).
			SetThumbnail("https://justhumanz.me/update.png").
			SetDescription("New Vtuber has been added to list").
			AddField("Nickname", Data.Member.Name).
			AddField("Twitter", "@"+Data.Member.TwitterName).
			InlineAllFields().
			AddField("Region", Data.Member.Region).
			SetURL("https://www.youtube.com/channel/"+Data.Member.YtID[0]+"?sub_confirmation=1").
			SetColor(Color).MessageEmbed)
		if err != nil {
			log.Error(msg, err)
		}
		msg, err = Bot.ChannelMessageSend(Channel, "New Update!!!! @here")
		if err != nil {
			log.Error(msg, err)
		}
	}
}
