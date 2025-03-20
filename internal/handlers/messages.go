package handlers

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		log.Printf("%s sent message in %s\n", m.Author.Username, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		log.Printf("%s sent message in %s\n", m.Author.Username, m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func MessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.BeforeUpdate == nil {
		return
	}

	fmt.Printf("%s changed their message from %s to %s\n", m.Author.Username, m.BeforeUpdate.Content, m.Content)
}

func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Printf("Emoji %s added\n", m.Emoji.Name)
	fmt.Printf("Member %s joined party\n", m.UserID)
}

func MessageReactionRemove(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	fmt.Printf("Member %s left party\n", m.UserID)
}
