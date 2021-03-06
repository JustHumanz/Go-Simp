package danbooru

import (
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	log "github.com/sirupsen/logrus"
)

func (Data Danbooru) SendNotif(Group database.Group, Member database.Member) {
	var Link = ""
	if Data.PixivID != 0 {
		Link = "https://www.pixiv.net/en/artworks/" + strconv.Itoa(Data.PixivID)
	} else if strings.HasPrefix(Data.Source, "https://twitter.com") {
		Link = Data.Source
	}

	ChannelData := database.ChannelTag(Member.ID, 0, config.LewdChannel, Member.Region)
	for _, Channel := range ChannelData {
		Msg, err := Bot.ChannelMessageSend(Channel.ChannelID, Link)
		if err != nil {
			log.Error(err, Msg)
		}
	}
}
