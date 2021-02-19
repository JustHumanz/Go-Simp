package twitter

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	log "github.com/sirupsen/logrus"
)

func SendFanart(Data []Fanart, Group database.Group) {
	for _, MemberFanart := range Data {
		url := MemberFanart.Tweet.PermanentURL
		ChannelData := database.ChannelTag(MemberFanart.Member.ID, 1, "", MemberFanart.Member.Region)
		var (
			tags  string
			Media string
			Msg   string
		)

		if len(MemberFanart.Tweet.Videos) > 0 {
			Media = MemberFanart.Tweet.Videos[0].Preview
			Msg = "1/1 Videos"
		} else if len(MemberFanart.Tweet.Photos) > 0 {
			Media = MemberFanart.Tweet.Photos[0]
			Msg = "1/" + strconv.Itoa(len(MemberFanart.Tweet.Photos)) + " Photos"
		} else {
			Media = config.NotFound
			Msg = "Photos/Video oversize,check original post"
		}

		Color, err := engine.GetColor(config.TmpDir, Media)
		if err != nil {
			log.Error(err)
		}

		if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
			Group.IconURL = ""
		}
		for i, Channel := range *ChannelData {
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
			if tags == "_" && Group.GroupName == "Independen" {
				//do nothing,like my life
			} else {
				msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
					SetAuthor(strings.Title(Group.GroupName), Group.IconURL).
					SetTitle("@"+MemberFanart.Tweet.Username).
					SetURL(url).
					SetThumbnail(engine.GetAuthorAvatar(MemberFanart.Tweet.Username)).
					SetDescription(RemoveTwitterShortLink(MemberFanart.Tweet.Text)).
					SetImage(Media).
					AddField("User Tags", tags).
					SetColor(Color).
					SetFooter(Msg, config.TwitterIMG).MessageEmbed)
				if err != nil {
					log.Error(msg, err)
					err = Channel.DelChannel(err.Error())
					if err != nil {
						log.Error(err)
					}
				}
				engine.Reacting(map[string]string{
					"ChannelID": Channel.ChannelID,
				}, Bot)
			}
			if i%config.Waiting == 0 && configfile.LowResources {
				log.WithFields(log.Fields{
					"Func": "Twitter Fanart",
				}).Warn(config.FanartSleep)
				time.Sleep(config.FanartSleep)
			}
		}
	}
}

/*
//SendNude Send to Discord channel
func (Data *TwitterFanart) SendNude() {
	for _, Fanart := range Data.Fanart {
		url := Fanart.PermanentURL
		ID, DiscordChannelID := database.ChannelTag(Data.Member.ID, 1, "")
		Bot := runner.Bot
		wg := new(sync.WaitGroup)
		for i := 0; i < len(DiscordChannelID); i++ {
			wg.Add(1)
			go func(DiscordChannel string, ID int, Data *TwitterFanart, wg *sync.WaitGroup) {
				defer wg.Done()
				ChannelState := database.DiscordChannel{
					ChannelID:     DiscordChannel,
					VtuberGroupID: Data.Group.ID,
				}
				UserTagsList := database.GetUserList(ID, Data.Member.ID)

				var (
					tags      string
					GroupIcon string
					Media     string
					Msg       string
				)

				if len(Fanart.Videos) > 0 {
					Media = Fanart.Videos[0].Preview
					Msg = "1/1 Videos"
				} else if len(Fanart.Photos) > 0 {
					Media = Fanart.Photos[0]
					Msg = "1/" + strconv.Itoa(len(Fanart.Photos)) + " Photos"
				} else {
					Media = config.NotFound
					Msg = "Photos/Video oversize,check original post"
				}

				Color, err := engine.GetColor("/tmp/tw", Media)
				if err != nil {
					log.Error(err)
				}

				if match, _ := regexp.MatchString("404.jpg", Data.Group.IconURL); match {
					GroupIcon = ""
				} else {
					GroupIcon = Data.Group.IconURL
				}
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
						SetTitle("@"+Fanart.Username).
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
							log.Error(err)
						}
					}
					engine.Reacting(map[string]string{
						"ChannelID": DiscordChannel,
					}, Bot)
				}
			}(DiscordChannelID[i], ID[i], Data, wg)
		}
		wg.Wait()
	}
}

*/
//RemoveTwitterShortLink remove twitter shotlink
func RemoveTwitterShortLink(text string) string {
	return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(text, "${1}$2")
}
