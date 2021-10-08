package engine

import (
	"sync"

	database "github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

//Start slash command
func InitSlash(Bot *discordgo.Session, GroupsPayload []database.Group, NewGuild *discordgo.Guild) {
	var (
		VtuberGroupChoices []*discordgo.ApplicationCommandOptionChoice
	)
	for _, v := range GroupsPayload {
		VtuberGroupChoices = append(VtuberGroupChoices, &discordgo.ApplicationCommandOptionChoice{
			Name:  v.GroupName,
			Value: v.ID,
		})
	}

	var (
		commands = []*discordgo.ApplicationCommand{
			{
				Name:        "setup",
				Description: "Setup bot",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
						Name:        "channel-type",
						Description: "select channel type",
						Options: []*discordgo.ApplicationCommandOption{
							{
								Name:        "livestream",
								Description: "Enable livestream notif on this channel",
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Options: []*discordgo.ApplicationCommandOption{
									{
										Name:        "channel-name",
										Description: "Setup channel",
										Required:    true,
										Type:        discordgo.ApplicationCommandOptionChannel,
									},
									{
										Type:        discordgo.ApplicationCommandOptionInteger,
										Name:        "vtuber-group",
										Description: "select vtuber-group",
										Choices:     VtuberGroupChoices,
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "liveonly",
										Description: "Set livestreams in strict mode(ignoring covering or regular video) notification",
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "newupcoming",
										Description: "Bot will send new upcoming livestream",
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "dynamic",
										Description: "Livestream message will disappear after livestream ended",
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "lite-mode",
										Description: "Disabling ping user/role function",
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "indie-notif",
										Description: "Send all independent vtubers notification **Ignore this if you not enable indie group**",
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "fanart",
										Description: "Enable vtuber fanart",
										Required:    false,
									},
								},
							},
							{
								Name:        "fanart",
								Description: "Enable fanart notif on this channel",
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Options: []*discordgo.ApplicationCommandOption{
									{
										Name:        "channel-name",
										Description: "Setup channel",
										Required:    true,
										Type:        discordgo.ApplicationCommandOptionChannel,
									},
									{
										Type:        discordgo.ApplicationCommandOptionInteger,
										Name:        "vtuber-group",
										Description: "select vtuber-group",
										Choices:     VtuberGroupChoices,
										Required:    true,
									},
									{
										Type:        discordgo.ApplicationCommandOptionBoolean,
										Name:        "lewd",
										Description: "Enable lewd vtuber fanart",
										Required:    false,
									},
								},
							},
							{
								Name:        "lewd",
								Description: "Enable lewd fanart notif on this channel",
								Type:        discordgo.ApplicationCommandOptionSubCommand,
								Options: []*discordgo.ApplicationCommandOption{
									{
										Name:        "channel-name",
										Description: "Setup channel",
										Required:    true,
										Type:        discordgo.ApplicationCommandOptionChannel,
									},
									{
										Type:        discordgo.ApplicationCommandOptionInteger,
										Name:        "vtuber-group",
										Description: "select vtuber-group",
										Choices:     VtuberGroupChoices,
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			{
				Name:        "channel-update",
				Description: "update channel-state information",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel-name",
						Description: "Select channel",
						Required:    true,
					},
				},
			},
			{
				Name:        "channel-state",
				Description: "get channel-state information",
				Options: []*discordgo.ApplicationCommandOption{

					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel-name",
						Description: "Select channel",
						Required:    true,
					},
				},
			},
			{
				Name:        "channel-delete",
				Description: "delete channel",
				Options: []*discordgo.ApplicationCommandOption{

					{
						Type:        discordgo.ApplicationCommandOptionChannel,
						Name:        "channel-delete",
						Description: "Select channel",
						Required:    true,
					},
				},
			},
			{
				Name:        "art",
				Description: "get random vtuber fanart",
				Options: []*discordgo.ApplicationCommandOption{

					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "group-name",
						Description: "Select vtuber GroupName",
						Choices:     VtuberGroupChoices,
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "Select vtuber",
						Required:    false,
					},
				},
			},
			{
				Name:        "lewd",
				Description: "get random vtuber lewd",
				Options: []*discordgo.ApplicationCommandOption{

					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "group-name",
						Description: "Select vtuber GroupName",
						Choices:     VtuberGroupChoices,
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "Select vtuber",
						Required:    false,
					},
				},
			},
			{
				Name:        "livestream",
				Description: "get random vtuber fanart",
				Options: []*discordgo.ApplicationCommandOption{

					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "state",
						Description: "Select livestream platform",
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Youtube",
								Value: 1,
							},
							{
								Name:  "BiliBili",
								Value: 2,
							},
							{
								Name:  "Twitch",
								Value: 3,
							},
						},
						Required: true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "status",
						Description: "Select livestream status",
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Live",
								Value: 1,
							},
							{
								Name:  "Upcoming",
								Value: 2,
							},
							{
								Name:  "Past",
								Value: 3,
							},
						},
						Required: true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "group-name",
						Description: "Select vtuber GroupName",
						Choices:     VtuberGroupChoices,
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "Select vtuber",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "region",
						Description: "Select region of vtuber group",
						Required:    false,
					},
				},
			},
			{
				Name:        "info",
				Description: "get subscriber/followers/viwers information",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "select vtuber",
						Required:    true,
					},
				},
			},
			{
				Name:        "tag-me",
				Description: "Add you to the tag list if any new fan art or live stream is uploaded",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "vtuber-group",
						Description: "select vtuber",
						Choices:     VtuberGroupChoices,
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "select vtuber",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "reminder",
						Description: "remind you before livestream started",
						Required:    false,
					},
				},
			},
			{
				Name:        "del-tag",
				Description: "delete a vtuber from your tag list",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "vtuber-group",
						Description: "select vtuber",
						Choices:     VtuberGroupChoices,
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "select vtuber",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "reminder",
						Description: "remind you before livestream started",
						Required:    false,
					},
				},
			},
			{
				Name:        "mytags",
				Description: "Shows all your info on this bot",
			},
			{
				Name:        "tag-role",
				Description: "Same like tag-me command,but this will tag roles",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role-name",
						Description: "Role",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "vtuber-group",
						Description: "select vtuber",
						Choices:     VtuberGroupChoices,
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "select vtuber",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "reminder",
						Description: "remind you before livestream started",
						Required:    false,
					},
				},
			},
			{
				Name:        "del-role",
				Description: "Same like del-tag command,but this will tag roles",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role-name",
						Description: "Role",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "vtuber-group",
						Description: "select vtuber",
						Choices:     VtuberGroupChoices,
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "vtuber-name",
						Description: "select vtuber",
						Required:    false,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "reminder",
						Description: "remind you before livestream started",
						Required:    false,
					},
				},
			},
			{
				Name:        "role-info",
				Description: "Get role info",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role-name",
						Description: "Role",
						Required:    true,
					},
				},
			},
			{
				Name:        "prediction",
				Description: "prediction vtuber subs/followers in next week",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "platform",
						Description: "Select platform",
						Required:    true,
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{
								Name:  "Twitter",
								Value: "tw",
							},
							{
								Name:  "BiliBili",
								Value: "bl",
							},
							{
								Name:  "Youtube",
								Value: "yt",
							},
						},
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Required:    true,
						Name:        "vtuber-name",
						Description: "input vtuber name",
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Required:    false,
						Name:        "prediction-count",
						Description: "imput prediction-count(days)\nexample,1 is tomorrow or 7 is next week",
					},
				},
			},
		}
	)

	if NewGuild == nil {
		log.Info("Start init slash command")
		var wg sync.WaitGroup
		for i, G := range Bot.State.Guilds {
			wg.Add(1)
			go func(wg *sync.WaitGroup, Guild *discordgo.Guild) {
				defer wg.Done()
				log.WithFields(log.Fields{
					"GuildName":    Guild.Name,
					"GuildID":      Guild.ID,
					"GuildOwnerID": Guild.OwnerID,
				}).Info("Create bot command")
				for _, v := range commands {
					_, err := Bot.ApplicationCommandCreate(Bot.State.User.ID, Guild.ID, v)
					if err != nil {
						log.Errorf("Cannot create '%v' command: %v guild: %v", v.Name, err, Guild.ID)
						if err.Error() == "HTTP 403 Forbidden" {
							return
						} else {
							continue
						}
					}
				}

			}(&wg, G)
			if (i % 20) == 0 {
				wg.Wait()
			}
		}
		wg.Wait()
		log.Info("Done init slash command")
	} else {
		log.Info("Register slash command to new guild")
		log.WithFields(log.Fields{
			"GuildName":    NewGuild.Name,
			"GuildID":      NewGuild.ID,
			"GuildOwnerID": NewGuild.OwnerID,
		}).Info("Create bot command")
		for _, v := range commands {
			_, err := Bot.ApplicationCommandCreate(Bot.State.User.ID, NewGuild.ID, v)
			if err != nil {
				log.Errorf("Cannot create '%v' command: %v guild: %v", v.Name, err, NewGuild.ID)
				continue
			}
		}
	}
}
