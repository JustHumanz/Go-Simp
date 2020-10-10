package twitter

import (
	"regexp"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	"github.com/prometheus/common/log"
)

type PushData struct {
	Twitter    database.InputTW
	Image      string
	Msg        string
	ScreenName string
	UserName   string
	Text       string
	Avatar     string
	Group      database.MemberGroupID
}

var (
	UrlTMP string
	Color  = 14807034
)

//Send to Discord channel
func (Data PushData) SendNude() error {
	ID, DiscordChannelID := database.ChannelTag(Data.Group.MemberID, 1)
	wg := new(sync.WaitGroup)
	url := Data.Twitter.Url
	for i := 0; i < len(DiscordChannelID); i++ {
		wg.Add(1)
		UserTagsList := database.GetUserList(ID[i], Data.Group.MemberID)

		go func(DiscordChannel string, wg *sync.WaitGroup) {
			defer wg.Done()
			Color, _ = engine.GetColor("/tmp/tw", Data.Image)
			if UrlTMP != url {
				if UserTagsList != nil {
					Embed := engine.NewEmbed().
						SetAuthor(strings.Title(Data.Group.GroupName), Data.Group.GroupIcon).
						SetTitle(Data.UserName+"(@"+Data.ScreenName+")").
						SetURL(url).
						SetThumbnail(strings.Replace(Data.Avatar, "_normal.jpg", ".jpg", -1)).
						SetDescription(RemoveTwitterShotlink(Data.Text)).
						SetImage(Data.Image).
						AddField("User Tags", strings.Join(UserTagsList, " ")).
						SetColor(Color).
						SetFooter(Data.Msg, config.TwitterIMG).MessageEmbed
					msg, err := BotSession.ChannelMessageSendEmbed(DiscordChannel, Embed)
					if err != nil {
						log.Error(msg, err)
					}
					engine.Reacting(map[string]string{
						"ChannelID": DiscordChannel,
					})
				}
			} else {
				log.Info("Same post,multiple hashtags")
				log.Info("UrlTMP :", UrlTMP, " PermanentURL: ", url)
			}
		}(DiscordChannelID[i], wg)
	}
	wg.Wait()
	UrlTMP = url
	return nil
}

//remove twitter shotlink
func RemoveTwitterShotlink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}
