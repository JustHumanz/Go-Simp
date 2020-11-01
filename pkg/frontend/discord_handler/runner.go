package discordhandler

import (
	"fmt"
	"os"
	"strings"

	config "github.com/JustHumanz/Go-simp/tools/config"
	database "github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//Prefix command
const (
	Enable       = "enable"
	Disable      = "disable"
	Update       = "update"
	TagMe        = "tag me"
	DelTag       = "del tag"
	MyTags       = "my tags"
	TagRoles     = "tag roles"
	RolesTags    = "roles tags"
	DelRoles     = "del roles"
	ChannelState = "channel state"
	VtuberData   = "vtuber data"
	Subscriber   = "subscriber"
	Upcoming     = "upcoming"
	Past         = "past"
	Live         = "live"
)

func StartInit(path string) {
	conf, err := config.ReadConfig(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	db := conf.CheckSQL()

	Bot, _ := discordgo.New("Bot " + config.Token)
	err = Bot.Open()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	database.Start(db)
	engine.Start()

	//Bot.AddHandler(GuildJoin)
	Bot.AddHandler(Fanart)
	Bot.AddHandler(Tags)
	Bot.AddHandler(EnableState)
	Bot.AddHandler(Status)
	Bot.AddHandler(Help)
	Bot.AddHandler(BiliBiliMessage)
	Bot.AddHandler(BiliBiliSpace)
	Bot.AddHandler(YoutubeMessage)
	Bot.AddHandler(SubsMessage)
}

//ValidName Find a valid name from user input
func ValidName(Name string) Memberst {

	for _, Group := range engine.GroupData {
		for _, Member := range database.GetName(Group.ID) {
			if Name == strings.ToLower(Member.Name) || Name == strings.ToLower(Member.JpName) {
				return Memberst{
					VTName:     engine.FixName(Member.EnName, Member.JpName),
					ID:         Member.ID,
					YtChannel:  Member.YoutubeID,
					SpaceID:    Member.BiliBiliID,
					BiliAvatar: Member.BiliBiliAvatar,
					YtAvatar:   Member.YoutubeAvatar,
				}
			}
		}
	}
	return Memberst{}
}
