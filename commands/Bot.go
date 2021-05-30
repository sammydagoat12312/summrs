package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/summrs-dev-team/summrs/events"
	"github.com/summrs-dev-team/summrs/utils"

	"github.com/bwmarrin/discordgo"
)

func (cmd *Commands) BotInfo(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Bot Info",

		Fields: []*discordgo.MessageEmbedField{
			{Name: "Name:", Value: s.State.User.Username, Inline: true},
			{Name: "Server Count:", Value: fmt.Sprint(events.GuildCount), Inline: true},
			{Name: "User Count:", Value: fmt.Sprint(events.MemberCount), Inline: true},
			{Name: "Ping:", Value: fmt.Sprintf("%s", s.HeartbeatLatency().Round(1*time.Millisecond)), Inline: true},
			{Name: "discordgo Version", Value: "v0.22.0", Inline: true},
			{Name: "Shard", Value: fmt.Sprintf("%d/%d", s.ShardID, s.ShardCount), Inline: true},
		},

		Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", m.Author.Username)},
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: s.State.User.AvatarURL("500")},
		Color:     0x36393F,
	})
}

func (cmd *Commands) Credits(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title: "Credits",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "__Main Developer__", Value: "[!fishgang Cy](https://github.com/Not-Cyrus)"},
			{Name: "__Developer/Hoster__", Value: "[lxi](https://github.com/lxi1400)\n[HK](https://giphy.com/gifs/brother-want-sister-lRQTVaje6kCzK)"},
			{Name: "__Networkers__,", Value: "[four](https://tenor.com/view/bearded-bear-guy-slay-gay-pride-super-gay-lgbt-gif-16465293)\n[Jinx](https://media1.giphy.com/media/F1wf27zzepXoZCfRkY/giphy-downsized-large.gif)\n[Sirus](https://giphy.com/gifs/kiss-spongebob-squarepants-lTQF0ODLLjhza)"},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", m.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmds *Commands) Fox(s *discordgo.Session, m *discordgo.Message, ctx *Context) {
	rand.Seed(time.Now().Unix())

	resBody, err := utils.MakeRequest("https://raw.githubusercontent.com/Not-Cyrus/fox-pic-repo/main/count.txt")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error: could not fetch the amount of fox pics, try re-running the command.")
		return
	}

	maxcount, err := strconv.Atoi(strings.TrimSuffix(string(resBody), "\n"))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong when attempting to get the amount")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("https://raw.githubusercontent.com/Not-Cyrus/fox-pic-repo/main/%d.jpg", rand.Intn(maxcount-0)+0))
}

func (cmd *Commands) Invite(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Bot Invite", Value: fmt.Sprintf("[Click Here](https://discord.com/api/oauth2/authorize?client_id=%s&permissions=8&scope=bot)", s.State.User.ID), Inline: true},
		},
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", m.Author.Username)},
		Color:  0x36393F,
	})
}

func (cmd *Commands) Ping(s *discordgo.Session, m *discordgo.Message, ctx *Context) {

	s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  fmt.Sprintf("Bot Ping: `%s`\nCurrent shard `%d/%d`", s.HeartbeatLatency().Round(1*time.Millisecond), s.ShardID, s.ShardCount),
		Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Requested by: %s", m.Author.Username)},
		Color:  0x36393F,
	})
}
