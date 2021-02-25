package twitch

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

func (Data TwitchNotif) SendNotif() error {
	Color, err := engine.GetColor(config.TmpDir, Data.TwitchData.Thumbnails)
	if err != nil {
		return err
	}
	//id, DiscordChannelID
	var (
		wg          sync.WaitGroup
		ChannelData = database.ChannelTag(Data.Member.ID, 2, "", Data.Member.Region)
		VtuberName  = engine.FixName(Data.Member.EnName, Data.Member.JpName)
		ImgURL      = "https://www.twitch.tv/" + Data.Member.TwitchName
		loc         = engine.Zawarudo(Data.Member.Region)
		expiresAt   = time.Now().In(loc)
		User        = &database.UserStruct{
			Human:    true,
			Reminder: 0,
		}
	)
	for i, v := range ChannelData {
		v.SetMember(Data.Member)

		wg.Add(1)
		go func(Channel database.DiscordChannel, wg *sync.WaitGroup) error {
			defer wg.Done()
			ctx := context.Background()
			UserTagsList, err := Channel.GetUserList(ctx)
			if err != nil {
				log.Error(err)
			}
			if UserTagsList == nil && Data.Group.GroupName != "Independen" {
				UserTagsList = []string{"_"}
			} else if UserTagsList == nil && Data.Group.GroupName == "Independen" && !Channel.IndieNotif {
				return nil
			}

			MsgEmbed, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
				SetAuthor(VtuberName, Data.Member.YoutubeAvatar, ImgURL).
				SetTitle("Live right now").
				SetDescription(Data.TwitchData.Title).
				SetImage(Data.TwitchData.Thumbnails).
				SetThumbnail(Data.Group.IconURL).
				SetURL(ImgURL).
				AddField("Start live", durafmt.Parse(expiresAt.Sub(Data.TwitchData.ScheduledStart.In(loc))).LimitFirstN(1).String()+" Ago").
				AddField("Viewers", strconv.Itoa(Data.TwitchData.Viewers)+" simps").
				SetFooter(Data.TwitchData.ScheduledStart.In(loc).Format(time.RFC822), config.TwitchIMG).
				SetColor(Color).MessageEmbed)
			if err != nil {
				log.WithFields(log.Fields{
					"Message":          MsgEmbed,
					"ChannelID":        Channel.ID,
					"DiscordChannelID": Channel.ChannelID,
				}).Error(err)
				err = Channel.DelChannel(err.Error())
				if err != nil {
					return err
				}
				return err
			}

			if Channel.Dynamic {
				log.WithFields(log.Fields{
					"DiscordChannel": Channel.ChannelID,
					"VtuberGroupID":  Data.Group.ID,
					"TwitchID":       "Twitch" + Data.Member.TwitchName,
				}).Info("Set dynamic mode")
				Channel.SetVideoID("Twitch" + Data.Member.TwitchName).
					SetMsgEmbedID(MsgEmbed.ID)
			}

			if !Channel.LiteMode {
				Msg := "Push " + configfile.Emoji.Livestream[0] + " to add you in `" + Data.Member.Name + "` ping list\nPush " + configfile.Emoji.Livestream[1] + " to remove you from ping list"
				msg, err := Bot.ChannelMessageSend(Channel.ChannelID, "`"+Data.Member.Name+"` Live right now\nUserTags: "+strings.Join(UserTagsList, " ")+"\n"+Msg)
				if err != nil {
					return err
				}
				User.SetDiscordChannelID(Channel.ChannelID).
					SetGroup(Data.Group).
					SetMember(Data.Member)

				err = User.SendToCache(msg.ID)
				if err != nil {
					return err
				}

				Channel.SetMsgTextID(msg.ID)

				err = engine.Reacting(map[string]string{
					"ChannelID": Channel.ChannelID,
					"State":     "Twitch",
					"MessageID": msg.ID,
				}, Bot)
				if err != nil {
					return err
				}
			}

			Channel.PushReddis()

			return nil
		}(v, &wg)
		//Wait every ge 5 discord channel
		if i%config.Waiting == 0 && configfile.LowResources {
			log.WithFields(log.Fields{
				"Func":  "Twitch",
				"Value": config.Waiting,
			}).Warn("Waiting send message")
			wg.Wait()
			expiresAt = time.Now().In(loc)
		}
	}
	wg.Wait()
	return nil
}
