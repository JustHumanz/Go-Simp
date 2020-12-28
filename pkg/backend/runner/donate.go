package runner

import (
	config "github.com/JustHumanz/Go-simp/tools/config"
	"github.com/JustHumanz/Go-simp/tools/database"
	engine "github.com/JustHumanz/Go-simp/tools/engine"
)

func Donate() {
	Donation := config.BotConf.DonationLink
	Bot.ChannelMessageSendEmbed(database.GetRanChannel(), engine.NewEmbed().
		SetTitle("Donate").
		SetURL(Donation).
		SetThumbnail(BotInfo.AvatarURL("128")).
		SetImage(config.GoSimpIMG).
		SetColor(14807034).
		SetDescription("Enjoy the bot?,don't forget to support this bot and dev").
		AddField("Ko-Fi", "[Link]("+Donation+")").
		AddField("Or if you a broke gang,you can upvote "+BotInfo.Username, "[top.gg]("+config.BotConf.TopGG+")").MessageEmbed)
}
