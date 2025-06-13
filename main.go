package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	attendees      = make(map[string]*discordgo.User)
	RemoveCommands = false
)

func main() {
	fmt.Println("Starting BobBot...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	token := os.Getenv("TOKEN")
	guildId := os.Getenv("GUILD_ID")

	bot := NewBot(token, guildId, RemoveCommands)

	bot.AddIntents()
	bot.MapCommandHandlers()

	err = bot.Session.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer bot.Session.Close()

	// create commands after connection is opened
	bot.CreateCommands()

	fmt.Println("Bot is now running")
	// rest api stuff

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	fmt.Println("Shutting down...")
	if RemoveCommands {
		log.Println("Removing commands...")
		cmds := bot.GetCommands()
		for _, cmd := range cmds {
			err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, guildId, cmd.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
			}
		}
	}
}
