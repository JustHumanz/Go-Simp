package main

import (
	"context"
	"encoding/json"
	"flag"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	network "github.com/JustHumanz/Go-Simp/pkg/network"
	pilot "github.com/JustHumanz/Go-Simp/service/pilot/grpc"
	"github.com/JustHumanz/Go-Simp/service/utility/runfunc"
	log "github.com/sirupsen/logrus"
)

var (
	Bot        *discordgo.Session
	configfile config.ConfigFile
	Payload    database.VtubersPayload
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: true})
}

func main() {
	Youtube := flag.Bool("Youtube", false, "Enable youtube module")
	BiliBili := flag.Bool("BiliBili", false, "Enable bilibili module")
	Twitter := flag.Bool("Twitter", false, "Enable twitter module")
	flag.Parse()

	gRCPconn := pilot.NewPilotServiceClient(network.InitgRPC(config.Pilot))

	RequestPay := func() {
		res, err := gRCPconn.ReqData(context.Background(), &pilot.ServiceMessage{
			Message: "Send me nude",
			Service: "Subscriber",
		})
		if err != nil {
			log.Fatalf("Error when request payload: %s", err)
		}
		err = json.Unmarshal(res.ConfigFile, &configfile)
		if err != nil {
			log.Panic(err)
		}

		err = json.Unmarshal(res.VtuberPayload, &Payload)
		if err != nil {
			log.Panic(err)
		}
	}
	RequestPay()

	var err error
	Bot, err = discordgo.New("Bot " + configfile.Discord)
	if err != nil {
		log.Error(err)
	}

	configfile.InitConf()
	database.Start(configfile)

	err = Bot.Open()
	if err != nil {
		log.Panic(err)
	}

	c := cron.New()
	c.Start()

	if *Youtube {
		c.AddFunc(config.YoutubeSubscriber, CheckYoutube)
		log.Info("Add youtube subscriber to cronjob")
	}

	if *BiliBili {
		c.AddFunc(config.BiliBiliFollowers, CheckBiliBili)
		log.Info("Add bilibili followers to cronjob")
	}

	if *Twitter {
		c.AddFunc(config.TwitterFollowers, CheckTwitter)
		log.Info("Add twitter followers to cronjob")
	}

	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "YoutubeSubscriber",
		Enabled: *Youtube,
	})
	if err != nil {
		log.Error(err)
	}
	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "BiliBiliFollowers",
		Enabled: *BiliBili,
	})
	if err != nil {
		log.Error(err)
	}
	_, err = gRCPconn.ModuleList(context.Background(), &pilot.ModuleData{
		Module:  "TwitterFollowers",
		Enabled: *Twitter,
	})
	if err != nil {
		log.Error(err)
	}

	go pilot.RunHeartBeat(gRCPconn, "Subscriber")
	runfunc.Run(Bot)
}

func SendNude(Embed *discordgo.MessageEmbed, Group database.Group, Member database.Member) {
	if match, _ := regexp.MatchString("404.jpg", Group.IconURL); match {
		Embed.Author.IconURL = ""
	}
	ChannelData := Group.GetChannelByGroup(Member.Region)
	for i, Channel := range ChannelData {
		Tmp := &Channel
		ctx := context.Background()
		UserTagsList, err := Tmp.SetMember(Member).SetGroup(Group).GetUserList(ctx)
		if err != nil {
			log.Error(err)
		}
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
		if i%config.Waiting == 0 && configfile.LowResources {
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
