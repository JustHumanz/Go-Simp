package subscriber

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	runner "github.com/JustHumanz/Go-simp/pkg/backend/runner"
	"github.com/JustHumanz/Go-simp/tools/database"
	log "github.com/sirupsen/logrus"
)

func SendNude(Embed *discordgo.MessageEmbed, Group database.GroupName, MemberID int64) {
	Bot := runner.Bot
	ChannelID, DiscordChannelID := Group.GetChannelByGroup()
	for i, Channel := range DiscordChannelID {
		UserTagsList := database.GetUserList(ChannelID[i], MemberID)
		msg, err := Bot.ChannelMessageSendEmbed(Channel, Embed)
		if err != nil {
			log.Error(msg, err)
		}
		if UserTagsList != nil {
			msg, err = Bot.ChannelMessageSend(DiscordChannelID[i], "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
			}
		}
	}
}

type Subs struct {
	Kind     string `json:"kind"`
	Etag     string `json:"etag"`
	PageInfo struct {
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind       string `json:"kind"`
		Etag       string `json:"etag"`
		ID         string `json:"id"`
		Statistics struct {
			ViewCount             string `json:"viewCount"`
			CommentCount          string `json:"commentCount"`
			SubscriberCount       string `json:"subscriberCount"`
			HiddenSubscriberCount bool   `json:"hiddenSubscriberCount"`
			VideoCount            string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type BiliBiliStat struct {
	LikeView LikeView
	Follow   BiliFollow
	Videos   int
}

type LikeView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Archive struct {
			View int `json:"view"`
		} `json:"archive"`
		Article struct {
			View int `json:"view"`
		} `json:"article"`
		Likes int `json:"likes"`
	} `json:"data"`
}

type BiliFollow struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Mid       int `json:"mid"`
		Following int `json:"following"`
		Whisper   int `json:"whisper"`
		Black     int `json:"black"`
		Follower  int `json:"follower"`
	} `json:"data"`
}
