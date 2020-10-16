package bilibili

import (
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/JustHumanz/Go-simp/config"
	database "github.com/JustHumanz/Go-simp/database"
	engine "github.com/JustHumanz/Go-simp/engine"

	log "github.com/sirupsen/logrus"
)

type Notif struct {
	TBiliData   database.InputTBiliBili
	Group       database.GroupName
	PhotosImgur string
	PhotosCount int
	MemberID    int64
}

//Push Data to discord channel
func (NotifData Notif) PushNotif(Color int) {
	Data := NotifData.TBiliData
	Group := NotifData.Group
	wg := new(sync.WaitGroup)
	ID, DiscordChannelID := database.ChannelTag(NotifData.MemberID, 1)
	for i := 0; i < len(DiscordChannelID); i++ {
		UserTagsList := database.GetUserList(ID[i], NotifData.MemberID)
		wg.Add(1)
		go func(DiscordChannel string, wg *sync.WaitGroup) {
			defer wg.Done()

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
			if UserTagsList != nil {
				tags = strings.Join(UserTagsList, " ")
			} else {
				tags = "_"
			}

			tmp, err := BotSession.ChannelMessageSendEmbed(DiscordChannel, engine.NewEmbed().
				SetAuthor(strings.Title(Group.NameGroup), Group.IconURL).
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
					log.Info("Delete Discord Channel ", DiscordChannel)
					database.DelChannel(DiscordChannel, NotifData.MemberID)
				}
			}
			err = engine.Reacting(map[string]string{
				"ChannelID": DiscordChannel,
			})
			if err != nil {
				log.Error(err)
			}
		}(DiscordChannelID[i], wg)
	}
	wg.Wait()
}
