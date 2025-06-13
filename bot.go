package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

const (
	cID = "1352366645510930523"
)

type Bot struct {
	Session        *discordgo.Session
	GuildID        string
	RemoveCommands bool
	Votes          map[string][]int // how a user voted
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
		Votes:          make(map[string][]int),
	}
}

func (b *Bot) MapCommandHandlers() {
	b.Session.AddHandler(CommandLocation)
	b.Session.AddHandler(CommandDate)
	b.Session.AddHandler(CommandCheckPolls)
}

func (b *Bot) AddIntents() {
	b.Session.Identify.Intents |= discordgo.IntentsGuildMessages
	b.Session.Identify.Intents |= discordgo.IntentMessageContent
	b.Session.Identify.Intents |= discordgo.IntentGuildMessagePolls
}

func (b *Bot) CreateCommands() {
	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		if err != nil {
			log.Fatalf("Cannot create '%v' command %v", cmd.Name, err)
		}
	}
}

func (b *Bot) GetCommands() []*discordgo.ApplicationCommand {
	cmds, err := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	if err != nil {
		return nil
	}
	return cmds
}

func (b *Bot) PrintCommands() {
	cmds, _ := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	for _, c := range cmds {
		fmt.Printf("Command name: %s, id: %s\n", c.Name, c.ID)
	}
}
