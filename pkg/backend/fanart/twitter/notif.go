package twitter

import (
	"regexp"
	"strconv"
	"strings"

	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
)

//Public variable
var (
	URLTMP string
	Color  = 14807034
)

//SendNude Send to Discord channel
func (Data *TwitterFanart) SendNude() error {
	for _, Fanart := range Data.Fanart {
		url := Fanart.PermanentURL
		ID, DiscordChannelID := database.ChannelTag(Data.Member.ID, 1, "")
		Bot := runner.Bot
		for i := 0; i < len(DiscordChannelID); i++ {
			ChannelState := database.DiscordChannel{
				ChannelID:     DiscordChannelID[i],
				VtuberGroupID: Data.Group.ID,
			}
			UserTagsList := database.GetUserList(ID[i], Data.Member.ID)

			Color, _ = engine.GetColor("/tmp/tw", Fanart.Photos[0])
			var (
				tags      string
				GroupIcon string
				Media     string
				Msg       string
			)

			if len(Fanart.Videos) > 0 {
				Media = Fanart.Videos[0].Preview
				Msg = "1/1 Videos"
			} else {
				Media = Fanart.Photos[0]
				Msg = "1/" + strconv.Itoa(len(Fanart.Photos)) + " Photos"
			}

			if match, _ := regexp.MatchString("404.jpg", Data.Group.IconURL); match {
				GroupIcon = ""
			} else {
				GroupIcon = Data.Group.IconURL
			}
			if URLTMP != url {
				if UserTagsList != nil {
					tags = strings.Join(UserTagsList, " ")
				} else {
					tags = "_"
				}
				if tags == "_" && Data.Group.NameGroup == "Independen" {
					//do nothing,like my life
				} else {
					msg, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], engine.NewEmbed().
						SetAuthor(strings.Title(Data.Group.NameGroup), GroupIcon).
						SetTitle("(@"+Fanart.Username+")").
						SetURL(url).
						SetThumbnail(engine.GetAuthorAvatar(Fanart.Username)).
						SetDescription(RemoveTwitterShortLink(Fanart.Text)).
						SetImage(Media).
						AddField("User Tags", tags).
						SetColor(Color).
						SetFooter(Msg, config.TwitterIMG).MessageEmbed)
					if err != nil {
						log.Error(msg, err)
						err = ChannelState.DelChannel(err.Error())
						if err != nil {
							return err
						}
					}
					engine.Reacting(map[string]string{
						"ChannelID": DiscordChannelID[i],
					}, Bot)
				}
				URLTMP = url
			} else {
				log.WithFields(log.Fields{
					"Old URL": URLTMP,
					"New URL": url,
				}).Info("Same post,multiple hashtags")
			}
		}
	}
	return nil
}

//RemoveTwitterShortLink remove twitter shotlink
func RemoveTwitterShortLink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}
