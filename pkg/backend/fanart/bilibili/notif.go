package bilibili

import (
	"context"
	"strconv"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	log "github.com/sirupsen/logrus"
)

//Notif struct
type Notif struct {
	TBiliData   database.InputTBiliBili
	Group       database.Group
	PhotosImgur string
	PhotosCount int
	Member      database.Member
}

//PushNotif Push Data to discord channel
func (NotifData Notif) PushNotif(Color int) {
	Data := NotifData.TBiliData
	Group := NotifData.Group
	ChannelData := database.ChannelTag(NotifData.Member.ID, 1, "")
	tags := ""
	for _, Channel := range ChannelData {
		UserTagsList := Channel.GetUserList(context.Background())
		if UserTagsList != nil {
			tags = strings.Join(UserTagsList, " ")
		} else {
			tags = "_"
		}

		if Group.GroupName == "Independen" {
			Group.IconURL = ""
		}
		if tags == "_" && Group.GroupName == "Independen" {
			//do nothing,like my life
		} else {
			tmp, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
				SetAuthor(strings.Title(Group.GroupName), Group.IconURL).
				SetTitle(Data.Author).
				SetURL(Data.URL).
				SetThumbnail(Data.Avatar).
				SetDescription(Data.Text).
				SetImage(NotifData.PhotosImgur).
				AddField("User Tags", tags).
				//AddField("Similar art", msg).
				SetFooter("1/"+strconv.Itoa(NotifData.PhotosCount)+" photos", config.BiliBiliIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.Error(tmp, err.Error())
				err = Channel.DelChannel(err.Error())
				if err != nil {
					log.Error(err)
				}
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": Channel.ChannelID,
			}, Bot)
			if err != nil {
				log.Error(err)
			}
		}
	}
	/*
		msg := ""
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
	*/
}
