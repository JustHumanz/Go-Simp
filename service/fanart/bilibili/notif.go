package bilibili

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	database "github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	log "github.com/sirupsen/logrus"
)

//PushNotif Push Data to discord channel
func PushNotif(Data database.TBiliBili) error {
	Group := Data.Group
	ChannelData := database.ChannelTag(Data.Member.ID, 1, "", Data.Member.Region)
	Color, err := engine.GetColor(config.TmpDir, Data.Photos[0])
	if err != nil {
		return err
	}
	tags := ""
	for i, Channel := range ChannelData {
		Channel.SetMember(Data.Member)
		ctx := context.Background()
		UserTagsList, err := Channel.GetUserList(ctx)
		if err != nil {
			log.Error(err)
			break
		}
		if UserTagsList != nil {
			tags = strings.Join(UserTagsList, " ")
		} else {
			tags = "_"
		}

		if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
			Group.IconURL = ""
		}

		if tags == "_" && Group.GroupName == "Independen" && !Channel.IndieNotif {
			//do nothing,like my life
		} else {
			tmp, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
				SetAuthor(strings.Title(Group.GroupName), Group.IconURL).
				SetTitle(Data.Author).
				SetURL(Data.URL).
				SetThumbnail(Data.Avatar).
				SetDescription(Data.Text).
				SetImage(Data.Photos[0]).
				AddField("User Tags", tags).
				//AddField("Similar art", msg).
				SetFooter("1/"+strconv.Itoa(len(Data.Photos))+" photos", config.BiliBiliIMG).
				InlineAllFields().
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.Error(tmp, err.Error())
				err = Channel.DelChannel(err.Error())
				if err != nil {
					return err
				}
				return err
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": Channel.ChannelID,
			}, Bot)
			if err != nil {
				return err
			}
		}
		if i%config.Waiting == 0 && config.GoSimpConf.LowResources {
			log.WithFields(log.Fields{
				"Func": "BiliBili Fanart",
			}).Warn(config.FanartSleep)
			time.Sleep(config.FanartSleep)
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
	return nil
}
