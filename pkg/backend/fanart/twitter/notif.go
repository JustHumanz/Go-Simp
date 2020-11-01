package twitter

import (
	"regexp"
	"strings"
	"sync"

	config "github.com/JustHumanz/Go-simp/config"
	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
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
	Bot := runner.Bot
	for i := 0; i < len(DiscordChannelID); i++ {
		wg.Add(1)
		UserTagsList := database.GetUserList(ID[i], Data.Group.MemberID)

		go func(DiscordChannel string, wg *sync.WaitGroup) {
			defer wg.Done()
			Color, _ = engine.GetColor("/tmp/tw", Data.Image)
			var tags string
			if UrlTMP != url {
				if UserTagsList != nil {
					tags = strings.Join(UserTagsList, " ")
				} else {
					tags = "_"
				}
				msg, err := Bot.ChannelMessageSendEmbed(DiscordChannel, engine.NewEmbed().
					SetAuthor(strings.Title(Data.Group.GroupName), Data.Group.GroupIcon).
					SetTitle(Data.UserName+"(@"+Data.ScreenName+")").
					SetURL(url).
					SetThumbnail(strings.Replace(Data.Avatar, "_normal.jpg", ".jpg", -1)).
					SetDescription(RemoveTwitterShotlink(Data.Text)).
					SetImage(Data.Image).
					AddField("User Tags", tags).
					SetColor(Color).
					SetFooter(Data.Msg, config.TwitterIMG).MessageEmbed)
				if err != nil {
					log.Error(msg, err)
					match, _ := regexp.MatchString("Unknown Channel", err.Error())
					if match {
						log.Info("Delete Discord Channel ", DiscordChannel)
						database.DelChannel(DiscordChannel, Data.Group.GroupID)
					}
				}
				engine.Reacting(map[string]string{
					"ChannelID": DiscordChannel,
				}, Bot)
			} else {
				log.WithFields(log.Fields{
					"Old URL": UrlTMP,
					"New URL": url,
				}).Info("Same post,multiple hashtags")
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
