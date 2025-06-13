package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "poll",
			Description: "Manually check for poll results",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionChannel,
					Name: "check",
					Description: "Channel option",
					ChannelTypes: []discordgo.ChannelType{
						discordgo.ChannelTypeGuildText,
					},
					Required: true,
				},
			},
		},
		{
			Name: "location",
			Description: "Location of the Bob Party",
		},
		{
			Name: "date",
			Description: "Date of the Bob Party",
		},
	}
)

func CommandCheckPolls(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	channelID := ""
	if opt, ok := optionMap["check"]; ok {
		channelID = opt.ChannelValue(nil).ID
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Checking polls in <#%s>...", channelID),
		},
	})

	msgs, err := s.ChannelMessages(channelID, 100, "", "", "")
	if err != nil {
		return
	}

	for _, msg := range msgs {
		poll := msg.Poll
		result := ""
		if poll != nil {
			if poll.Results.Finalized {
				maxCount := -1
				winningAnswerID := -1
				for _, ans := range poll.Results.AnswerCounts {
					if ans.Count > maxCount {
						winningAnswerID = ans.ID
					}
				}

				for _, ans := range poll.Answers {
					if ans.AnswerID == winningAnswerID {
						result = ans.Media.Text
					}
				}
			}
		}

		if result == "" {
			continue
		}

		// Store in KV
		// { "location": "Atlanta" }
		// { "date": "2026-07-12" }
		fmt.Printf("For msg %s, result was %s.\n", msg.ID, result)
	}

	time.Sleep(time.Second * 2)
	content := fmt.Sprintf("Finished checking for <#%s>", channelID)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func CommandLocation(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "TBD",
		},
	})
}

func CommandDate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "TBD",
		},
	})
}
