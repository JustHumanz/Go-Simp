package danbooru

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	log "github.com/sirupsen/logrus"
)

func (Data Danbooru) SendNotif(Group database.Group, Member database.Member) {
	ChannelData, err := database.ChannelTag(Member.ID, 0, config.LewdChannel, Member.Region)
	if err != nil {
		log.Error(err)
	}

	var Link = ""
	if Data.PixivID != 0 {
		Link = "https://www.pixiv.net/en/artworks/" + strconv.Itoa(Data.PixivID)
		var (
			Illusts map[string]interface{}
			User    map[string]interface{}
		)
		illusbyte, err := network.Curl(config.PixivIllustsEnd+strconv.Itoa(Data.PixivID), nil)
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(illusbyte, &Illusts)
		if err != nil {
			log.Error(err)
		}

		Body := Illusts["body"].(map[string]interface{})
		Tags := Body["tags"].(map[string]interface{})
		Img := Body["urls"].(map[string]interface{})
		FixImg := config.PixivProxy + Img["original"].(string)

		usrbyte, err := network.Curl(config.PixivUserEnd+Tags["authorId"].(string), nil)
		if err != nil {
			log.Error(err)
		}

		err = json.Unmarshal(usrbyte, &User)
		if err != nil {
			log.Error(err)
		}

		for _, Channel := range ChannelData {
			Msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
				SetAuthor(strings.Title(Group.GroupName), Group.IconURL).
				SetTitle(User["name"].(string)).
				SetURL(Link).
				SetThumbnail(config.PixivProxy+User["imageBig"].(string)).
				SetDescription(Body["title"].(string)).
				SetImage(FixImg).MessageEmbed)
			if err != nil {
				log.Error(err, Msg)
			}
		}

	} else if strings.HasPrefix(Data.Source, "https://twitter.com") {
		Link = Data.Source
	} else {
		Link = Data.FileURL
	}
	for _, Channel := range ChannelData {
		Msg, err := Bot.ChannelMessageSend(Channel.ChannelID, Link)
		if err != nil {
			log.Error(err, Msg)
		}
	}
}
