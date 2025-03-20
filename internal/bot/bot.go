package bot

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/erobx/bobBot/internal/handlers"
)

var command = &discordgo.ApplicationCommand{
	Name:        "start",
	Description: "Start bot",
}

type Bot struct {
	Session        *discordgo.Session
	GuildID        string
	RemoveCommands bool
}

func NewBot(token, guildId string, rmvCmds bool) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	return &Bot{
		Session:        s,
		GuildID:        guildId,
		RemoveCommands: rmvCmds,
	}
}

func (b *Bot) MapMessageHandlers() {
	b.Session.AddHandler(handlers.MessageCreate)
	b.Session.AddHandler(handlers.MessageUpdate)
	b.Session.AddHandler(handlers.MessageReactionAdd)
	b.Session.AddHandler(handlers.MessageReactionRemove)
}

func (b *Bot) MapCommandHandlers() {
	b.Session.AddHandler(handlers.CommandStart)
}

func (b *Bot) AddIntents() {
	b.Session.Identify.Intents |= discordgo.IntentsGuildMessages
	b.Session.Identify.Intents |= discordgo.IntentMessageContent
}

func (b *Bot) CreateCommands() {
	_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, command)
	if err != nil {
		log.Fatalf("Cannot create '%v' command %v", command.Name, err)
	}
}

func (b *Bot) PrintCommands() {
	cmds, _ := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	for _, c := range cmds {
		fmt.Printf("Command name %s\n", c.Name)
	}
}
