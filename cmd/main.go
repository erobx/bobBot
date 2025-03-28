package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/erobx/bobBot/internal/bot"
	"github.com/joho/godotenv"
)

var (
	attendees      = make(map[string]*discordgo.User)
	RemoveCommands = true

	command = &discordgo.ApplicationCommand{
		Name:        "start-polls",
		Description: "Start polls",
	}
)

func main() {
	fmt.Println("Starting BobBot...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	token := os.Getenv("TOKEN")
	guildId := os.Getenv("GUILD_ID")

	bot := bot.NewBot(token, guildId, RemoveCommands)

	bot.MapMessageHandlers()
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

	go func() {
		votes := make(map[string][]int)
		for {
			answerIds := []int{1, 2}
			for _, id := range answerIds {
				voters := bot.GetVoters(id)
				for _, voter := range voters {
					if _, ok := votes[voter.GlobalName]; ok {
						votes[voter.GlobalName] = append(votes[voter.GlobalName], id)
					} else {
						votes[voter.GlobalName] = []int{id}
					}
				}
			}
			for username, ids := range votes {
				fmt.Printf("User %s voted for %v\n", username, ids)
			}
			time.Sleep(time.Second * 5)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit

	if RemoveCommands {
		log.Println("Removing commands...")
		err := bot.Session.ApplicationCommandDelete(bot.Session.State.User.ID, guildId, command.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", command.Name, err)
		}
	}
}
