package bilibili

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

//Notif struct
type Notif struct {
	TBiliData   database.InputTBiliBili
	Group       database.GroupName
	PhotosImgur string
	PhotosCount int
	MemberID    int64
}

//PushNotif Push Data to discord channel
func (NotifData Notif) PushNotif(Color int) {
	Data := NotifData.TBiliData
	Group := NotifData.Group
	Bot := runner.Bot
	ID, DiscordChannelID := database.ChannelTag(NotifData.MemberID, 1, "")
	GroupIcon := ""
	msg := ""
	tags := ""
	repost, url, err := engine.SaucenaoCheck(strings.Split(Data.Photos, "\n")[0])
	if err != nil {
		log.Error(err)
		msg = "??????"
	} else if repost && url != nil {
		log.WithFields(log.Fields{
			"Source Img": Data.URL,
			"Sauce Img":  url,
		}).Info("Repost")
		msg = url[0]
	} else {
		log.WithFields(log.Fields{
			"Source Img": Data.URL,
			"Sauce Img":  url,
		}).Info("Ntap,Anyar cok")
		msg = "_"
	}
	for i := 0; i < len(DiscordChannelID); i++ {
		ChannelState := database.DiscordChannel{
			ChannelID:     DiscordChannelID[i],
			VtuberGroupID: Group.ID,
		}
		UserTagsList := database.GetUserList(ID[i], NotifData.MemberID)
		if UserTagsList != nil {
			tags = strings.Join(UserTagsList, " ")
		} else {
			tags = "_"
		}

		if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
			GroupIcon = ""
		} else {
			GroupIcon = Group.IconURL
		}

		if tags == "_" && Group.NameGroup == "Independen" {

		} else {
			tmp, err := Bot.ChannelMessageSendEmbed(DiscordChannelID[i], engine.NewEmbed().
				SetAuthor(strings.Title(Group.NameGroup), GroupIcon).
				SetTitle(Data.Author).
				SetURL(Data.URL).
				SetThumbnail(Data.Avatar).
				SetDescription(Data.Text).
				SetImage(NotifData.PhotosImgur).
				AddField("User Tags", tags).
				AddField("Similar art", msg).
				SetFooter("1/"+strconv.Itoa(NotifData.PhotosCount)+" photos", config.BiliBiliIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.Error(tmp, err.Error())
				match, _ := regexp.MatchString("Unknown Channel", err.Error())
				if match {
					log.Info("Delete Discord Channel ", DiscordChannelID[i])
					err = ChannelState.DelChannel()
					if err != nil {
						log.Error(err)
					}
				}
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": DiscordChannelID[i],
			}, Bot)
			if err != nil {
				log.Error(err)
			}
		}
	}
}
