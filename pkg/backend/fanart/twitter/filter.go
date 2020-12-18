package twitter

import (
	"database/sql"
	"regexp"
	"strconv"
	"strings"
	"sync"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	log "github.com/sirupsen/logrus"
)

//CheckNew check if new fanart or not
func (Data TwitterStruct) CheckNew() []Statuses {
	var (
		tmp []Statuses
	)
	for _, TwData := range Data.Statuses {
		var (
			id int
		)
		err := database.DB.QueryRow(`SELECT id FROM Twitter WHERE PermanentURL=?`, "https://twitter.com/"+TwData.User.ScreenName+"/status/"+TwData.IDStr).Scan(&id)
		if err == sql.ErrNoRows {
			tmp = append(tmp, TwData)
		} else {
			//update
			_, err := database.DB.Exec(`Update Twitter set Likes=? Where id=? `, TwData.FavoriteCount, id)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return tmp
}

//CheckHashTag filter hashtag post
func (Data Statuses) CheckHashTag(Group []database.MemberGroup, wg *sync.WaitGroup) {
	defer wg.Done()
	rgx := "(?m)(.+free|leak|wrong|antihololive|asacoco|haachama|Lv[0-9]|Lv.+|taiwan|kson.+)"
	tiananmen, _ := regexp.MatchString(rgx, Data.Text)

	for _, hashtag := range Data.Entities.Hashtags {
		westtaiwan, _ := regexp.MatchString(rgx, hashtag.Text)
		if !westtaiwan && !tiananmen {
			for i := 0; i < len(Group); i++ {
				//just temporary rule
				if Group[i].JpName == "桐生ココ" {
					if Data.User.FollowersCount < 70 && Data.User.FriendsCount < 100 && Data.User.FavouritesCount < 100 && Data.User.StatusesCount < 100 && len(Data.Entities.Hashtags) > 4 {
						//fuck off dummy account
						log.WithFields(log.Fields{
							"Hashtags":   Group[i].TwitterHashtags,
							"MemberName": Group[i].EnName,
							"URL":        "https://twitter.com/" + Data.User.ScreenName + "/status/" + Data.IDStr,
						}).Info("dummy account(prob)")
					}
					return
				} else if Group[i].GroupName == "Hololive" {
					if Data.User.FollowersCount < 25 && Data.User.FriendsCount < 30 && Data.User.StatusesCount < 100 && len(Data.Entities.Hashtags) > 4 {
						//fuck off dummy account
						log.WithFields(log.Fields{
							"Hashtags":   Group[i].TwitterHashtags,
							"MemberName": Group[i].EnName,
							"URL":        "https://twitter.com/" + Data.User.ScreenName + "/status/" + Data.IDStr,
						}).Info("dummy account(prob)")
						return
					}
				}

				if "#"+hashtag.Text == Group[i].TwitterHashtags {
					//new
					log.WithFields(log.Fields{
						"Hashtags":   Group[i].TwitterHashtags,
						"MemberName": Group[i].EnName,
						"URL":        "https://twitter.com/" + Data.User.ScreenName + "/status/" + Data.IDStr,
					}).Info("Get new post")

					var (
						Photos    []string
						Video     string
						SendMedia string
						msg       string
					)
					for _, Media := range Data.Entities.Media {
						Photos = append(Photos, Media.MediaURLHTTPS)
					}
					for _, vid := range Data.ExtendedEntities.Media {
						if vid.VideoInfo.Variants != nil {
							Video = vid.VideoInfo.Variants[0].URL
						}
					}
					if Photos != nil && Video == "" {
						SendMedia = Photos[0]
						msg = "1/" + strconv.Itoa(len(Data.ExtendedEntities.Media)) + " photos"
					} else if Video != "" {
						SendMedia = Data.ExtendedEntities.Media[0].MediaURLHTTPS
						msg = "Video type,check original post"
					} else {
						SendMedia = config.NotFound
						msg = "Image or Video oversize,check original post"
					}
					TwitterData := PushData{
						Twitter: database.InputTW{
							Url:      "https://twitter.com/" + Data.User.ScreenName + "/status/" + Data.IDStr,
							Author:   Data.User.Name,
							Like:     Data.FavoriteCount,
							Photos:   strings.Join(Photos, "\n"),
							Video:    Video,
							Text:     Data.Text,
							TweetID:  Data.IDStr,
							MemberID: Group[i].MemberID,
						},
						Image:      SendMedia,
						Msg:        msg,
						ScreenName: Data.User.ScreenName,
						UserName:   Data.User.Name,
						Text:       RemoveTwitterShortLink(Data.Text),
						Avatar:     (strings.Replace(Data.User.ProfileImageURLHTTPS, "_normal.jpg", ".jpg", -1)),
						Group:      Group[i],
					}
					TwitterData.Twitter.InputTwitter()
					TwitterData.SendNude()
				}
			}
		}
	}
}
