package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/summrs-dev-team/summrs/database"
	"github.com/summrs-dev-team/summrs/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (cmds *Commands) AddOwner(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if m.Mentions[0].ID == utils.GetGuildOwner(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, "You can not change the status of this user")
		return
	}

	if whitelisted, err := database.Database.SetOwner(m.GuildID, m.Mentions[0], true); !whitelisted {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Made that user an owner (They can now run owner-only commands).")
}

func (cmd *Commands) AntiInvite(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if !(ctx.Fields[0] == "on" || ctx.Fields[0] == "off") {
		return
	}

	if _, err := database.Database.SetData(m.GuildID, "anti-invite", ctx.Fields[0]); err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Set Anti-Invite to %s", ctx.Fields[0]))
}

func (cmds *Commands) DelOwner(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	if m.Mentions[0].ID == utils.GetGuildOwner(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, "You can not change the status of this user")
		return
	}

	if whitelisted, err := database.Database.SetOwner(m.GuildID, m.Mentions[0], false); !whitelisted {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Revoked that users owner status (They can no longer run owner-only commands).")
}

func (cmd *Commands) LoggingChannel(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if set, err := database.Database.SetData(message.GuildID, "log-channel", message.ChannelID); !set {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(message.ChannelID, "Set the logging channel to the current channel")
}

func (cmd *Commands) Prefix(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if set, err := database.Database.SetData(message.GuildID, "prefix", ctx.Fields[0]); !set {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("Prefix has been set to `%s`", ctx.Fields[0]),
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", message.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmd *Commands) Settings(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	data, err := database.Database.FindData(message.GuildID)
	guild, _ := s.State.Guild(message.GuildID)

	if err != nil {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	var (
		embed = &discordgo.MessageEmbed{
			Title:  fmt.Sprintf("%s current settings", guild.Name),
			Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", message.Author.Username)},
			Color:  0x36393F,
		}
		tempValue string
	)

	for index, value := range data {
		if index == "users" || index == "owners" || index == "_id" || index == "guild_id" {
			continue
		}

		switch value.(string) {
		case "on":
			tempValue = "<:enabled:799507631274197022>"

		case "off":
			tempValue = "<:disabled:799507673648594954>"

		case "nil":
			tempValue = "<:disabled:799507673648594954>"

		default:
			tempValue = value.(string)
			if index == "log-channel" {
				tempValue = fmt.Sprintf("<#%s>", value.(string))
			}
		}

		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   index,
			Value:  tempValue,
			Inline: false,
		})

	}
	s.ChannelMessageSendEmbed(message.ChannelID, embed)
}

func (cmd *Commands) Whitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if whitelisted, err := database.Database.SetWhitelist(message.GuildID, message.Mentions[0], true); !whitelisted {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(message.ChannelID, "Whitelisted that user.")
}

func (cmd *Commands) Unwhitelist(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	if whitelisted, err := database.Database.SetWhitelist(message.GuildID, message.Mentions[0], false); !whitelisted {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(message.ChannelID, "Unwhitelisted that user.")
}

func (cmd *Commands) ViewWhitelisted(s *discordgo.Session, message *discordgo.Message, ctx *Context) {
	data, err := database.Database.FindData(message.GuildID)

	if err != nil {
		s.ChannelMessageSend(message.ChannelID, err.Error())
		return
	}

	var whitelistedUsers []string

	for _, userID := range data["owners"].(bson.A) {
		user, err := s.User(userID.(string))
		if err != nil {
			continue
		}

		whitelistedUsers = append(whitelistedUsers, fmt.Sprintf("👑 | %s#%s", user.Username, user.Discriminator))
	}

	for _, userID := range data["users"].(bson.A) {
		user, err := s.User(userID.(string))
		if err != nil {
			continue
		}

		if !database.Database.IsOwner(message.GuildID, user.ID) {
			whitelistedUsers = append(whitelistedUsers, fmt.Sprintf("📋 | %s#%s", user.Username, user.Discriminator))
		}
	}

	s.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
		Title:       "Whitelisted Members",
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", message.Author.Username)},
		Description: strings.Join(whitelistedUsers, "\n"),
		Color:       0x36393F,
	})
}
