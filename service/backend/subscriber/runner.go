package main

import (
	"context"
	"flag"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	engine "github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/service/backend/utility/runfunc"
	log "github.com/sirupsen/logrus"
)

var (
	Bot *discordgo.Session
)

func main() {
	Youtube := flag.Bool("Youtube", false, "Enable youtube module")
	BiliBili := flag.Bool("BiliBili", false, "Enable bilibili module")
	Twitter := flag.Bool("Twitter", false, "Enable twitter module")

	flag.Parse()

	conf, err := config.ReadConfig("../../../config.toml")
	if err != nil {
		log.Panic(err)
	}

	Bot, err = discordgo.New("Bot " + config.BotConf.Discord)
	if err != nil {
		log.Error(err)
	}

	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	database.Start(conf.CheckSQL())
	engine.Start()

	c := cron.New()
	c.Start()

	if *Youtube {
		c.AddFunc(config.YoutubeSubscriber, CheckYoutube)
		log.Info("Add youtube subscriber to cronjob")
		database.ModuleInfo("YoutubeSubscriber")
	}

	if *BiliBili {
		c.AddFunc(config.BiliBiliFollowers, CheckBiliBili)
		log.Info("Add bilibili followers to cronjob")
		database.ModuleInfo("BiliBiliFollowers")
	}

	if *Twitter {
		c.AddFunc(config.TwitterFollowers, CheckTwitter)
		log.Info("Add twitter followers to cronjob")
		database.ModuleInfo("TwitterFollowers")
	}
	runfunc.Run(Bot)
}

func SendNude(Embed *discordgo.MessageEmbed, Group database.Group, Member database.Member) {
	if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
		Embed.Author.IconURL = ""
	}
	ChannelData := Group.GetChannelByGroup()
	for i, Channel := range ChannelData {
		Tmp := &Channel
		UserTagsList := Tmp.SetMember(Member).SetGroup(Group).GetUserList(context.Background()) //database.GetUserList(Channel.ID, MemberID)
		msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, Embed)
		if err != nil {
			log.Error(msg, err)
		}
		if UserTagsList != nil {
			msg, err = Bot.ChannelMessageSend(Channel.ChannelID, "UserTags: "+strings.Join(UserTagsList, " "))
			if err != nil {
				log.Error(msg, err)
			}
		}
		if i%config.Waiting == 0 && config.BotConf.LowResources {
			log.WithFields(log.Fields{
				"Func":  "Subscriber",
				"Value": config.Waiting,
			}).Warn("Waiting send message")
			time.Sleep(100 * time.Millisecond)
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
