package discord

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/erobx/bobBot/internal/types"
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "flight",
			Description: "Add flight info",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "number",
					Description: "Flight number",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "airport-code",
					Description: "Airport Code",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "arrival-date",
					Description: "Date flight arrives (YYYY-MM-DD)",
					Required: true,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "arrival-time",
					Description: "Flight arrival time in 24hr (HH:MM)",
					Required: true,
				},
			},
		},
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

func (b *Bot) CommandAddFlight(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	flightNumber := ""
	airportCode := ""
	arrivalDate := ""
	arrivalTime := ""

	if opt, ok := optionMap["number"]; ok {
		flightNumber = opt.StringValue()
	}

	if opt, ok := optionMap["airport-code"]; ok {
		airportCode = opt.StringValue()
	}
	
	if opt, ok := optionMap["arrival-date"]; ok {
		arrivalDate = opt.StringValue()
	}

	if opt, ok := optionMap["arrival-time"]; ok {
		arrivalTime = opt.StringValue()
	}

	invokee := i.Member
	userID := ""
	if invokee != nil {
		userID = i.Member.User.ID
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Adding flight...",
		},
	})

	//aeroUrl := "https://aeroapi.flightaware.com/aeroapi"

	//req, err := http.NewRequest(
	//	"GET",
	//	fmt.Sprintf(aeroUrl + "/flights/%s", flightNumber),
	//	bytes.NewReader([]byte{}),
	//)

	//if err != nil {
	//	return
	//}
	//req.Header.Set("x-apikey", b.aeroKey)

	//res, err := http.DefaultClient.Do(req)
	//if err != nil {
	//	return
	//}

	if _, err := time.Parse("2006-01-02", arrivalDate); err != nil {
		return
	}

	if _, err := time.Parse("15:04", arrivalTime); err != nil {
		return
	}

	flight := types.Flight{
		UserID: userID,
		Number: flightNumber,
		AirportCode: airportCode,
		ArrivalDate: arrivalDate,
		ArrivalTime: arrivalTime,
	}

	log.Printf("Adding flight number %s for user %s", flightNumber, userID)

	err := b.flightStore.Put(context.TODO(), flight)
	if err != nil {
		log.Printf("error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Sorry something went wrong, please try again.",
		})
		return
	}

	content := "Flight added."
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func (b *Bot) CommandCheckPolls(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	msgs, err := s.ChannelMessages(channelID, 10, "", "", "")
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

		item := types.Poll{
			Question: poll.Question.Text,
			Answer: result,
		}

		err := b.pollStore.Put(context.Background(), item)
		if err != nil {
			log.Printf("error: %v", err)
			ErrorRespond(s, i)
			return
		}
	}

	time.Sleep(time.Second * 2)
	content := fmt.Sprintf("Finished checking for <#%s>", channelID)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

func (b *Bot) CommandLocation(s *discordgo.Session, i *discordgo.InteractionCreate) {
	poll, err := b.pollStore.Get(context.Background(), "Location")
	if err != nil {
		log.Printf("error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Location hasn't been decided.",
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Location: %s", poll.Answer),
		},
	})
}

func (b *Bot) CommandDate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	poll, err := b.pollStore.Get(context.TODO(), "Date")
	if err != nil {
		log.Printf("error: %v", err)
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Date hasn't been decided.",
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Date: %s", poll.Answer),
		},
	})
}

func ErrorRespond(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Sorry something went wrong, please try again.",
	})
}
