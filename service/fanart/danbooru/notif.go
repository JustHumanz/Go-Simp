package danbooru

import (
	"github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

func (Data Danbooru) SendNotif(Group database.Group, Member database.Member) {
	ChannelData := database.ChannelTag(Member.ID, 0, "Lewd", Member.Region)
	for _, Channel := range ChannelData {
		Msg, err := Bot.ChannelMessageSend(Channel.ChannelID, Data.FileURL)
		if err != nil {
			log.Error(err, Msg)
		}
	}
}
