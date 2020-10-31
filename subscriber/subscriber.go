package subscriber

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/robfig/cron.v2"

	"github.com/JustHumanz/Go-simp/config"
	"github.com/JustHumanz/Go-simp/database"
	"github.com/JustHumanz/Go-simp/engine"
	log "github.com/sirupsen/logrus"
)

var (
	BiliSession string
	yttoken     string
)

func Start(c *cron.Cron) {
	BiliSession = config.BiliBiliSes
	yttoken = config.YtToken[1]
	if BiliSession == "" {
		log.Error("BiliBili Session not found")
		os.Exit(1)
	}
	c.AddFunc("@every 2h0m0s", CheckYtSubsCount)
	c.AddFunc("@every 0h15m0s", CheckTwFollowCount)
	c.AddFunc("@every 0h30m0s", CheckBiliFollowCount)
	log.Info("Subs&Follow Checker module ready")
}

func SendNude(Embed *discordgo.MessageEmbed, Group database.GroupName) {
	for _, Channel := range Group.GetChannelByGroup() {
		msg, err := engine.BotSession.ChannelMessageSendEmbed(Channel, Embed)
		if err != nil {
			log.Error(msg, err)
		}
		/*
			msg, err = engine.BotSession.ChannelMessageSend(Channel, "@here")
			if err != nil {
				log.Error(msg, err)
			}
		*/
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
