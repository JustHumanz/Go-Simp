package lewd

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/engine"
	"github.com/JustHumanz/Go-Simp/pkg/network"
	log "github.com/sirupsen/logrus"
)

func GetDan(Data database.Group) {
	wg := new(sync.WaitGroup)
	for i, Mem := range Data.Members {
		wg.Add(1)
		if i%10 == 0 {
			time.Sleep(5 * time.Second)
		}
		go func(Member database.Member, w *sync.WaitGroup) {
			defer w.Done()
			log.WithFields(log.Fields{
				"Group":   Data.GroupName,
				"Vtubers": Member.Name,
				"Site":    "Danbooru",
			}).Info("Check lewd pic")

			var (
				DanPayload []Danbooru
				SendFanart = func(FanArt database.DataFanart, AuthorImg string, Color int) {
					ChannelData, err := database.ChannelTag(Member.ID, 0, config.LewdChannel, Member.Region)
					if err != nil {
						log.Error(err)
					}

					for _, Channel := range ChannelData {
						Msg, err := Bot.ChannelMessageSendEmbed(Channel.ChannelID, engine.NewEmbed().
							SetAuthor(strings.Title(Data.GroupName), Data.IconURL).
							SetTitle(FanArt.Author).
							SetURL(FanArt.PermanentURL).
							SetThumbnail(AuthorImg).
							SetDescription(FanArt.Text).
							SetImage(FanArt.Photos[0]).
							SetColor(Color).MessageEmbed)
						if err != nil {
							log.Error(err, Msg)
						}
					}
				}
				TwitterHandler = func(ID []string) error {
					TweetRes, err := config.Scraper.GetTweet(ID[len(ID)-1])
					if err != nil {
						return err
					}
					var Video string
					if TweetRes.Videos != nil {
						Video = TweetRes.Videos[0].URL
					}

					FanArtData := database.DataFanart{
						PermanentURL: TweetRes.PermanentURL,
						Author:       TweetRes.Username,
						Photos:       TweetRes.Photos,
						Videos:       Video,
						Text:         TweetRes.Text,
						TweetID:      TweetRes.ID,
						Member:       Member,
					}
					err = database.AddLewd(FanArtData)
					if err != nil {
						return err
					}

					Color, err := engine.GetColor(config.TmpDir, TweetRes.Photos[0])
					if err != nil {
						return err
					}

					SendFanart(FanArtData, engine.GetAuthorAvatar(TweetRes.Username), Color)
					return nil
				}
			)

			if Member.TwitterLewd != "" {
				log.WithFields(log.Fields{
					"Group":   Data.GroupName,
					"Vtubers": Member.Name,
					"Site":    "Twitter",
				}).Info("Check lewd pic")

				for tweet := range config.Scraper.SearchTweets(context.Background(), Member.TwitterLewd+" AND -filter:replies -filter:retweets -filter:quote filter:media OR filter:link", 20) {
					if len(tweet.Photos) > 0 && !tweet.IsQuoted && !tweet.IsReply {
						if database.IsLewdNew("Twitter", tweet.PermanentURL) {
							log.WithFields(log.Fields{
								"Group":    Data.GroupName,
								"Vtubers":  Member.Name,
								"TweetURL": tweet.PermanentURL,
							}).Info("New Lewd pic from Twitter")
							TweetID := strings.Split(tweet.PermanentURL, "/")
							err := TwitterHandler(TweetID)
							if err != nil {
								log.Error(err)
							}
						}
					}
				}
			}

			databyte, err := network.Curl(config.DanbooruEndPoint+strings.Replace(Member.EnName, " ", "_", -1)+"&limit=10", nil)
			if err != nil {
				log.Error(err)
			}

			err = json.Unmarshal(databyte, &DanPayload)
			if err != nil {
				log.Error(err)
			}

			for _, Dan := range DanPayload {
				if Dan.CheckLewd() {
					if Dan.ParentID == nil {
						if Dan.IsPixiv() {
							if database.IsLewdNew("Pixiv", strconv.Itoa(Dan.PixivID)) {
								log.WithFields(log.Fields{
									"Group":      Data.GroupName,
									"Vtubers":    Member.Name,
									"DanbooruID": Dan.ID,
								}).Info("New Lewd pic from Pixiv")

								var (
									Illusts map[string]interface{}
									User    map[string]interface{}
								)
								illusbyte, err := network.Curl(config.PixivIllustsEnd+strconv.Itoa(Dan.PixivID), nil)
								if err != nil {
									log.Error(err)
									break
								}

								err = json.Unmarshal(illusbyte, &Illusts)
								if err != nil {
									log.Error(err)
								}

								Body := Illusts["body"].(map[string]interface{})
								Tags := Body["tags"].(map[string]interface{})
								Img := Body["urls"].(map[string]interface{})
								FixImg := config.PixivProxy + Img["original"].(string)
								MiniImg := config.PixivProxy + Img["mini"].(string)

								usrbyte, err := network.Curl(config.PixivUserEnd+Tags["authorId"].(string), nil)
								if err != nil {
									log.Error(err)
								}

								err = json.Unmarshal(usrbyte, &User)
								if err != nil {
									log.Error(err)
								}

								UserBody := User["body"].(map[string]interface{})
								FanArtData := database.DataFanart{
									PermanentURL: "https://www.pixiv.net/en/artworks/" + strconv.Itoa(Dan.PixivID),
									Author:       UserBody["name"].(string),
									Photos:       []string{FixImg},
									Text:         Body["title"].(string),
									PixivID:      strconv.Itoa(Dan.PixivID),
									Member:       Member,
								}
								err = database.AddLewd(FanArtData)
								if err != nil {
									log.Error(err)
								}

								Color, err := engine.GetColor(config.TmpDir, MiniImg)
								if err != nil {
									log.Error(err)
								}

								SendFanart(FanArtData, config.PixivProxy+UserBody["imageBig"].(string), Color)
							}
						} else if Dan.IsTwitter() {
							if database.IsLewdNew("Twitter", Dan.Source) {
								log.WithFields(log.Fields{
									"Group":      Data.GroupName,
									"Vtubers":    Member.Name,
									"DanbooruID": Dan.ID,
								}).Info("New Lewd pic from Twitter")
								TweetID := strings.Split(Dan.Source, "/")
								err := TwitterHandler(TweetID)
								if err != nil {
									log.Error(err)
								}
							}
						}
					}
				}
			}
		}(Mem, wg)
	}
	wg.Wait()
}
