package twitter

import (
	"regexp"
	"strings"
	"sync"

	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
)

//PushData Push data to discord channel struct
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

//Public variable
var (
	URLTMP string
	Color  = 14807034
)

//SendNude Send to Discord channel
func (Data PushData) SendNude() error {
	ID, DiscordChannelID := database.ChannelTag(Data.Group.MemberID, 1, "")
	wg := new(sync.WaitGroup)
	url := Data.Twitter.Url
	Bot := runner.Bot
	for i := 0; i < len(DiscordChannelID); i++ {
		wg.Add(1)

		ChannelState := database.DiscordChannel{
			ChannelID:     DiscordChannelID[i],
			VtuberGroupID: Data.Group.GroupID,
		}
		UserTagsList := database.GetUserList(ID[i], Data.Group.MemberID)

		go func(DiscordChannel string, wg *sync.WaitGroup) {
			defer wg.Done()
			Color, _ = engine.GetColor("/tmp/tw", Data.Image)
			var (
				tags      string
				GroupIcon string
			)

			if match, _ := regexp.MatchString("404.jpg", Data.Group.GroupIcon); match {
				GroupIcon = ""
			} else {
				GroupIcon = Data.Group.GroupIcon
			}
			if URLTMP != url {
				if UserTagsList != nil {
					tags = strings.Join(UserTagsList, " ")
				} else {
					tags = "_"
				}
				if tags == "_" && Data.Group.GroupName == "Independen" {
					//do nothing,like my life
				} else {
					msg, err := Bot.ChannelMessageSendEmbed(DiscordChannel, engine.NewEmbed().
						SetAuthor(strings.Title(Data.Group.GroupName), GroupIcon).
						SetTitle(Data.UserName+"(@"+Data.ScreenName+")").
						SetURL(url).
						SetThumbnail(strings.Replace(Data.Avatar, "_normal.jpg", ".jpg", -1)).
						SetDescription(RemoveTwitterShortLink(Data.Text)).
						SetImage(Data.Image).
						AddField("User Tags", tags).
						SetColor(Color).
						SetFooter(Data.Msg, config.TwitterIMG).MessageEmbed)
					if err != nil {
						log.Error(msg, err)
						err = ChannelState.DelChannel(err.Error())
						if err != nil {
							log.Error(err)
						}
					}
					engine.Reacting(map[string]string{
						"ChannelID": DiscordChannel,
					}, Bot)
				}
			} else {
				log.WithFields(log.Fields{
					"Old URL": URLTMP,
					"New URL": url,
				}).Info("Same post,multiple hashtags")
			}
		}(DiscordChannelID[i], wg)
	}
	wg.Wait()
	URLTMP = url
	return nil
}

//RemoveTwitterShortLink remove twitter shotlink
func RemoveTwitterShortLink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}
