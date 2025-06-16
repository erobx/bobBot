package discord

import (
	"context"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/erobx/bobBot/internal/store"
)

type Bot struct {
	Session        	*discordgo.Session
	flightStore		*store.FlightStore
	pollStore		*store.PollStore
	guildID        	string
	channelID	   	string
	aeroKey			string
	removeCommands 	bool
}

func NewBot(token, guildID, channelID, aeroKey string, rmvCmds bool) *Bot {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	client := store.NewDynamoClient(context.Background())
	fs := store.NewFlightStore(client, "flights")
	ps := store.NewPollStore(client, "poll_results")

	return &Bot{
		Session:        s,
		guildID:        guildID,
		flightStore: 	fs,
		pollStore: 		ps,
		channelID: 		channelID,
		aeroKey: 		aeroKey,
		removeCommands: rmvCmds,
	}
}

func (b *Bot) MapCommandHandlers() {
	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		"flight": b.CommandAddFlight,
		"poll": b.CommandCheckPolls,
		"location": b.CommandLocation,
		"date": b.CommandDate,
	}

	b.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func (b *Bot) AddIntents() {
	b.Session.Identify.Intents |= discordgo.IntentsGuildMessages
	b.Session.Identify.Intents |= discordgo.IntentMessageContent
	b.Session.Identify.Intents |= discordgo.IntentGuildMessagePolls
}

func (b *Bot) CreateCommands() {
	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.guildID, cmd)
		if err != nil {
			log.Fatalf("Cannot create '%v' command %v", cmd.Name, err)
		}
	}
}

func (b *Bot) GetCommands() []*discordgo.ApplicationCommand {
	cmds, err := b.Session.ApplicationCommands(b.Session.State.User.ID, b.guildID)
	if err != nil {
		return nil
	}
	return cmds
}

func (b *Bot) PrintCommands() {
	cmds, _ := b.Session.ApplicationCommands(b.Session.State.User.ID, b.guildID)
	for _, c := range cmds {
		fmt.Printf("Command name: %s, id: %s\n", c.Name, c.ID)
	}
}
