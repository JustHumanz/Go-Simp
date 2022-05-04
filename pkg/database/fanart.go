package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"regexp"
	"strings"
	"time"

	config "github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/go-redis/redis/v8"
	twitterscraper "github.com/n0madic/twitter-scraper"
	log "github.com/sirupsen/logrus"
)

//GetRandomFanart Get Member fanart URL from TBiliBili and Twitter
func GetRandomFanart(GroupID, MemberID int64) (*DataFanart, error) {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
	)

	Twitter := func() error {
		rows, err := DB.Query(`SELECT Twitter.* FROM Vtuber.Twitter Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.Member.ID)
			if err != nil {
				return err
			}
		}

		if Data.ID == 0 {
			return errors.New("vtuber don't have any fanart in Twitter")
		}
		Data.State = config.TwitterArt
		return nil
	}
	Tbilibili := func() error {
		rows, err := DB.Query(`SELECT TBiliBili.* FROM Vtuber.TBiliBili Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.Dynamic_id, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		if Data.ID == 0 {
			return errors.New("vtuber don't have any fanart in BiliBili")
		}

		Data.State = config.BiliBiliArt
		return nil
	}

	Pixiv := func() error {
		rows, err := DB.Query(`SELECT Pixiv.* FROM Vtuber.Pixiv Inner Join Vtuber.VtuberMember on VtuberMember.id = Pixiv.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Text, &Data.PixivID, &Data.Member.ID)
			if err != nil {
				return err
			}
		}
		if Data.ID == 0 {
			return errors.New("vtuber don't have any fanart in Pixiv")
		}

		Data.State = config.PixivArt
		return nil
	}

	gachaint := rand.Intn(3-1) + 1
	if gachaint == 1 {
		err := Pixiv()
		if err != nil {
			err = Twitter()
			if err != nil {
				err = Tbilibili()
				if err != nil {
					return nil, errors.New("vtuber don't have any fanart in Pixiv/Twitter/BiliBili")
				}
			}
		}
	} else if gachaint == 2 {
		err := Tbilibili()
		if err != nil {
			err = Twitter()
			if err != nil {
				err = Pixiv()
				if err != nil {
					return nil, errors.New("vtuber don't have any fanart in Pixiv/Twitter/BiliBili")
				}
			}
		}
	} else {
		err := Twitter()
		if err != nil {
			err = Pixiv()
			if err != nil {
				err = Tbilibili()
				if err != nil {
					return nil, errors.New("vtuber don't have any fanart in Pixiv/Twitter/BiliBili")
				}
			}
		}
	}

	Data.Videos = Video.String
	Data.Photos = strings.Fields(PhotoTmp.String)
	return &Data, nil

}

//GetLewd Get Member lewd fanart URL from and Twitter and Pixiv
func GetLewd(GroupID, MemberID int64) (*DataFanart, error) {
	var (
		Data     DataFanart
		PhotoTmp sql.NullString
		Video    sql.NullString
	)

	rows, err := DB.Query(`SELECT Lewd.* FROM Vtuber.Lewd Inner Join Vtuber.VtuberMember on VtuberMember.id = Lewd.VtuberMember_id Inner Join Vtuber.VtuberGroup on VtuberGroup.id = VtuberMember.VtuberGroup_id where (VtuberGroup.id=? OR VtuberMember.id=?)  ORDER by RAND() LIMIT 1`, GroupID, MemberID)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, errors.New("vtuber don't have any lewd fanart in Twitter or Pixiv")
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.PixivID, &Data.Member.ID)
		if err != nil {
			return nil, err
		}
	}
	if Data.PixivID != "" {
		Data.State = config.PixivArt
	} else {
		Data.State = config.TwitterArt
	}
	Data.Photos = strings.Fields(PhotoTmp.String)
	Data.Videos = Video.String
	return &Data, nil

}

//DeleteFanart Delete fanart when get 404 error status
func (Data DataFanart) DeleteFanart(e string) error {
	if notfound, _ := regexp.MatchString("404", e); notfound {
		log.Info("Delete fanart metadata ", Data.PermanentURL)
		if Data.State == "Twitter" {
			stmt, err := DB.Prepare(`DELETE From Twitter WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		} else {
			stmt, err := DB.Prepare(`DELETE From TBiliBili WHERE id=?`)
			if err != nil {
				return err
			}
			defer stmt.Close()

			stmt.Exec(Data.ID)
			return nil
		}
	} else {
		return nil
	}
}

//Add new lewd fanart
func (Data DataFanart) AddLewd() (int64, error) {
	log.WithFields(log.Fields{
		"Vtuber": Data.Member.Name,
		"Img":    Data.Photos,
		"URL":    Data.PermanentURL,
	}).Info("New Lewd Fanart")
	stmt, err := DB.Prepare(`INSERT INTO Lewd (PermanentURL,Author,Photos,Videos,Text,TweetID,PixivID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(Data.PermanentURL, Data.Author, strings.Join(Data.Photos, "\n"), Data.Videos, Data.Text, Data.TweetID, Data.PixivID, Data.Member.ID)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

//CheckMemberFanart Check if `that` was a new fanart
func (FanArt DataFanart) CheckTweetFanArt(Update bool) (bool, error) {
	var (
		id int
	)

	FanartCache := CheckFanartFromCache(FanArt.PermanentURL)
	if FanartCache.ID == 0 {
		if FanArt.Lewd {
			err := DB.QueryRow(`SELECT id FROM Lewd WHERE TweetID=?`, FanArt.TweetID).Scan(&id)
			if err == sql.ErrNoRows {
				id, err := FanArt.AddLewd()
				if err != nil {
					return false, err
				}
				FanArt.AddFanartToCache(id)

				return true, nil
			}
		} else {
			err := DB.QueryRow(`SELECT id FROM Twitter WHERE TweetID=?`, FanArt.TweetID).Scan(&id)
			if err == sql.ErrNoRows {
				log.WithFields(log.Fields{
					"Name":    FanArt.Member.EnName,
					"Hashtag": FanArt.Member.TwitterHashtag,
					"Lewd":    FanArt.Lewd,
					"URL":     FanArt.PermanentURL,
				}).Info("New Fanart")

				stmt, err := DB.Prepare(`INSERT INTO Twitter (PermanentURL,Author,Likes,Photos,Videos,Text,TweetID,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
				if err != nil {
					return false, err
				}
				defer stmt.Close()

				res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.TweetID, FanArt.Member.ID)
				if err != nil {
					return false, err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return false, err
				}

				FanArt.AddFanartToCache(id)
				return true, nil
			} else if err != nil {
				return false, err
			}
		}
	} else if Update {
		//update like
		log.WithFields(log.Fields{
			"Name":    FanArt.Member.EnName,
			"Hashtag": FanArt.Member.TwitterHashtag,
			"Likes":   FanArt.Likes,
		}).Info("Update like")
		_, err := DB.Exec(`Update Twitter set Likes=? Where id=? `, FanArt.Likes, FanartCache.ID)
		if err != nil {
			return false, err
		}

	}
	return false, nil
}

//Check if `this` was a new fanart
func (FanArt DataFanart) CheckTBiliBiliFanArt() (bool, error) {
	FanartCache := CheckFanartFromCache(FanArt.PermanentURL)
	if FanartCache.ID == 0 {
		var tmp int64
		row := DB.QueryRow("SELECT id FROM Vtuber.TBiliBili where Dynamic_id=?", FanArt.Dynamic_id)
		err := row.Scan(&tmp)
		if err == sql.ErrNoRows {
			log.WithFields(log.Fields{
				"Vtuber": FanArt.Member.EnName,
				"Img":    FanArt.Photos,
			}).Info("New Fanart")
			stmt, err := DB.Prepare(`INSERT INTO TBiliBili (PermanentURL,Author,Likes,Photos,Videos,Text,Dynamic_id,VtuberMember_id) values(?,?,?,?,?,?,?,?)`)
			if err != nil {
				return false, err
			}
			defer stmt.Close()

			res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, FanArt.Likes, strings.Join(FanArt.Photos, "\n"), FanArt.Videos, FanArt.Text, FanArt.Dynamic_id, FanArt.Member.ID)
			if err != nil {
				return false, err
			}

			id, err := res.LastInsertId()
			if err != nil {
				return false, err
			}

			FanArt.AddFanartToCache(id)
			return true, nil
		} else if err != nil {
			return false, err
		}
	}
	return false, nil
}

func CheckFanartFromCache(PermanentURL string) DataFanart {
	var Fanart DataFanart
	val2, err := FanartCache.Get(context.Background(), PermanentURL).Result()
	if err == redis.Nil {
		return Fanart
	} else if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal([]byte(val2), &Fanart)
	if err != nil {
		log.Fatal(err)
	}
	return Fanart

}

func (Data DataFanart) AddFanartToCache(id int64) {
	log.WithFields(log.Fields{
		"Key": Data.PermanentURL,
	}).Info("Add fanart into cache")

	Data.ID = id
	bytefanart, err := json.Marshal(Data)
	if err != nil {
		log.Fatal(err)
	}
	err = FanartCache.Set(context.Background(), Data.PermanentURL, bytefanart, 10*time.Minute).Err()
	if err != nil {
		log.Fatal(err)
	}
}

//Check if `this` was a new fanart
func (FanArt DataFanart) CheckPixivFanArt() (bool, error) {
	FanartCache := CheckFanartFromCache(FanArt.PermanentURL)
	if FanartCache.ID == 0 {
		var tmp int64
		if FanArt.Lewd {
			row := DB.QueryRow("SELECT id FROM Vtuber.Lewd where PixivID=?", FanArt.PixivID)
			err := row.Scan(&tmp)
			if err == sql.ErrNoRows {
				id, err := FanArt.AddLewd()
				if err != nil {
					return false, err
				}

				FanArt.AddFanartToCache(id)
				return true, nil
			} else if err != nil {
				return false, err
			}
		} else {
			row := DB.QueryRow("SELECT id FROM Vtuber.Pixiv where PixivID=?", FanArt.PixivID)
			err := row.Scan(&tmp)
			if err == sql.ErrNoRows {
				log.WithFields(log.Fields{
					"Vtuber": FanArt.Member.EnName,
					"Img":    FanArt.Photos,
					"URL":    FanArt.PermanentURL,
				}).Info("New Fanart")
				stmt, err := DB.Prepare(`INSERT INTO Pixiv (PermanentURL,Author,Photos,Text,PixivID,VtuberMember_id) values(?,?,?,?,?,?)`)
				if err != nil {
					return false, err
				}
				defer stmt.Close()

				res, err := stmt.Exec(FanArt.PermanentURL, FanArt.Author, strings.Join(FanArt.Photos, "\n"), FanArt.Text, FanArt.PixivID, FanArt.Member.ID)
				if err != nil {
					return false, err
				}

				id, err := res.LastInsertId()
				if err != nil {
					return false, err
				}

				FanArt.AddFanartToCache(id)
				return true, nil

			} else if err != nil {
				return false, err
			}
		}
	}
	return false, errors.New(FanArt.Member.Name + " Still same")
}

//DataFanart fanart struct
type DataFanart struct {
	ID           int64
	Member       Member
	Group        Group
	PermanentURL string
	Author       string
	AuthorAvatar string
	Photos       []string
	Videos       string
	Text         string
	Likes        int
	Dynamic_id   string
	TweetID      string
	PixivID      string
	Lewd         bool
	State        string
	FilePath     string
}

func (Data *DataFanart) MarshallBin() []byte {
	bit, err := json.Marshal(Data)
	if err != nil {
		log.Error(err)
	}
	return bit
}

func (Data *DataFanart) AddMember(new Member) *DataFanart {
	Data.Member = new
	return Data
}

func (Data *DataFanart) AddGroup(new Group) *DataFanart {
	Data.Group = new
	return Data
}

func (Data *DataFanart) AddPermanentURL(new string) *DataFanart {
	Data.PermanentURL = new
	return Data
}

func (Data *DataFanart) AddAuthor(new string) *DataFanart {
	Data.Author = new
	return Data
}

func (Data *DataFanart) AddAuthorAvatar(new string) *DataFanart {
	Data.AuthorAvatar = new
	return Data
}

func (Data *DataFanart) AddPhotos(new []string) *DataFanart {
	Data.Photos = new
	return Data
}

func (Data *DataFanart) AddVideos(new string) *DataFanart {
	Data.Videos = new
	return Data
}

func (Data *DataFanart) AddText(new string) *DataFanart {
	Data.Text = new
	return Data
}

func (Data *DataFanart) AddDynamicID(new string) *DataFanart {
	Data.Dynamic_id = new
	return Data
}

func (Data *DataFanart) AddTweetID(new string) *DataFanart {
	Data.TweetID = new
	return Data
}

func (Data *DataFanart) AddPixivID(new string) *DataFanart {
	Data.PixivID = new
	return Data
}

func (Data *DataFanart) SetLewd(new bool) *DataFanart {
	Data.Lewd = new
	return Data
}

func (Data *DataFanart) SetState(new string) *DataFanart {
	Data.State = new
	return Data
}

//Get random fanart from group struct
func (p *Group) GetRandomFanart() (*DataFanart, error) {
	b, err := GetRandomFanart(p.ID, 0)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Get random fanart from member struct
func (p *Member) GetRandomFanart() (*DataFanart, error) {
	b, err := GetRandomFanart(0, p.ID)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Get random lewd fanart from group struct
func (p *Group) GetRandomLewd() (*DataFanart, error) {
	b, err := GetLewd(p.ID, 0)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//Get random lewd fanart from member struct
func (p *Member) GetRandomLewd() (*DataFanart, error) {
	b, err := GetLewd(0, p.ID)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (Member Member) ScrapTwitterFanart(Scraper *twitterscraper.Scraper, Lewd bool, update bool) ([]DataFanart, error) {

	var (
		FanartList []DataFanart
	)

	Query := func() string {
		if Lewd {
			return Member.TwitterLewd + " AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)"
		} else {
			return Member.TwitterHashtag + " AND (-filter:replies -filter:retweets -filter:quote) AND (filter:media OR filter:link)"
		}
	}()

	for tweet := range Scraper.SearchTweets(context.Background(), Query, config.GoSimpConf.LimitConf.TwitterFanart) {
		if tweet.Error != nil {
			log.Error(tweet.Error)
			continue
		}

		for _, TweetHashtag := range tweet.Hashtags {
			if (strings.EqualFold("#"+TweetHashtag, Member.TwitterHashtag) || strings.EqualFold("#"+TweetHashtag, Member.TwitterLewd)) && !tweet.IsQuoted && !tweet.IsReply && len(tweet.Photos) > 0 {
				TweetArt := DataFanart{
					PermanentURL: tweet.PermanentURL,
					Author:       tweet.Username,
					AuthorAvatar: func() string {
						profile, err := Scraper.GetProfile(tweet.Username)
						if err != nil {
							log.Error(err)
						}
						return strings.Replace(profile.Avatar, "normal.jpg", "400x400.jpg", -1)
					}(),
					TweetID: tweet.ID,
					Text: func() string {
						return regexp.MustCompile(`(?m)^(.*?)https:\/\/t.co\/.+`).ReplaceAllString(tweet.Text, "${1}$2")
					}(),
					Photos: tweet.Photos,
					Likes:  tweet.Likes,
					Member: Member,
					State:  config.TwitterArt,
					Lewd:   Lewd,
					Group:  Member.Group,
				}
				if tweet.Videos != nil {
					TweetArt.Videos = tweet.Videos[0].Preview
				}

				New, err := TweetArt.CheckTweetFanArt(update)
				if err != nil {
					return nil, err
				}

				if New {
					FanartList = append(FanartList, TweetArt)
				}
			}
		}
	}

	if len(FanartList) > 0 {
		return FanartList, nil
	} else {
		return nil, errors.New("still same")
	}
}

func (Member Member) GetFanartData(State string, Limit int) ([]DataFanart, error) {
	var (
		Datafanart []DataFanart
		PhotoTmp   sql.NullString
		Video      sql.NullString
	)

	Twitter := func() error {
		rows, err := DB.Query(`SELECT Twitter.* FROM Vtuber.Twitter Inner Join Vtuber.VtuberMember on VtuberMember.id = Twitter.VtuberMember_id where  VtuberMember.id=?  ORDER by id desc LIMIT ?`, Member.ID, Limit)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.TweetID, &Data.Member.ID)
			if err != nil {
				return err
			}

			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}

			Data.State = config.TwitterArt
			Data.Videos = Video.String
			Data.Photos = strings.Fields(PhotoTmp.String)
			Data.AddMember(Member)
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}
	Tbilibili := func() error {
		rows, err := DB.Query(`SELECT TBiliBili.* FROM Vtuber.TBiliBili Inner Join Vtuber.VtuberMember on VtuberMember.id = TBiliBili.VtuberMember_id where VtuberMember.id=?  ORDER by id desc LIMIT ?`, Member.ID, Limit)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &Data.Likes, &PhotoTmp, &Video, &Data.Text, &Data.Dynamic_id, &Data.Member.ID)
			if err != nil {
				return err
			}
			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}

			Data.State = config.BiliBiliArt
			Data.Photos = strings.Fields(PhotoTmp.String)
			Data.AddMember(Member)
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}

	Pixiv := func() error {
		rows, err := DB.Query(`SELECT Pixiv.* FROM Vtuber.Pixiv Inner Join Vtuber.VtuberMember on VtuberMember.id = Pixiv.VtuberMember_id where VtuberMember.id=? ORDER by id desc LIMIT ?`, Member.ID, Limit)
		if err != nil {
			return err
		}

		defer rows.Close()
		for rows.Next() {
			var Data DataFanart
			err = rows.Scan(&Data.ID, &Data.PermanentURL, &Data.Author, &PhotoTmp, &Data.Text, &Data.PixivID, &Data.Member.ID)
			if err != nil {
				return err
			}
			if Data.ID == 0 {
				return errors.New("vtuber don't have any fanart in Twitter")
			}

			Data.State = config.PixivArt
			Data.Photos = strings.Fields(PhotoTmp.String)
			Data.AddMember(Member)
			Datafanart = append(Datafanart, Data)
		}
		return nil
	}

	if State == config.PixivArt {
		err := Pixiv()
		if err != nil {
			return nil, err
		}
	} else if State == config.BiliBiliArt {
		err := Tbilibili()
		if err != nil {
			return nil, err
		}
	} else {
		err := Twitter()
		if err != nil {
			return nil, err
		}
	}
	return Datafanart, nil
}

func (Group Group) GetFanartData(State string, Limit int) ([]DataFanart, error) {
	var Data []DataFanart
	for i := 1; i < Limit; i++ {
		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(Group.Members)
		DataFanart, err := Group.Members[n].GetFanartData(State, 1)
		if err != nil {
			return nil, err
		}
		Data = append(Data, DataFanart...)
	}
	return Data, nil
}
