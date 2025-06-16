package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/erobx/bobBot/internal/discord"
	"github.com/joho/godotenv"
)

var (
	RemoveCommands = true
)

func main() {
	fmt.Println("Starting BobBot...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	token := os.Getenv("TOKEN")
	guildID := os.Getenv("GUILD_ID")
	channelID := os.Getenv("BOT_CHANNEL_ID")
	aeroKey := os.Getenv("AERO_API_KEY")

	bot := discord.NewBot(token, guildID, channelID, aeroKey, RemoveCommands)
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	fmt.Println("Shutting down...")
	if RemoveCommands {
		log.Println("Removing commands...")
		cmds := bot.GetCommands()
		for _, cmd := range cmds {
			err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, guildID, cmd.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
			}
		}
	}
}
