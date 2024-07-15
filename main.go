package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"scholar-bot/apihandlers"
)

var botSession *discordgo.Session

func init() {
	var botToken string
	var err error
	botToken = os.Getenv("discord_bot_token")

	botSession, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot token: %v", err)
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "hello",
			Description: "A command to say back hello to the user",
		},
		{
			Name:        "gs",
			Description: "Get first study found on google scholar",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "search query for study",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate){
		"hello": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			botSession.InteractionRespond(
				botInteraction.Interaction,
				&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Hey there! Congratulations, your first command worked",
					},
				})
		},
		"gs": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			options := botInteraction.ApplicationCommandData().Options
			optionMap := make(
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
				len(options),
			)
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if query, ok := optionMap["query"]; ok {

				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: apihandlers.QueryFirstGs(query.StringValue()),
						},
					})
			} else {
				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error happened when retrieving the study",
						},
					})
			}

		},
	}
)

func init() {
	botSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	botSession.AddHandler(func(botSession *discordgo.Session, botReady *discordgo.Ready) {
		log.Printf(
			"Logged in as: %v#%v",
			botSession.State.User.Username,
			botSession.State.User.Discriminator,
		)
	})
	err := botSession.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Printf("Adding %d commands...\n", len(commands))
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := botSession.ApplicationCommandCreate(
			botSession.State.User.ID,
			"",
			v,
		)
		log.Printf("Created %v command and registered", v.Name)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer botSession.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Removing commands...")

	for _, v := range registeredCommands {
		err := botSession.ApplicationCommandDelete(
			botSession.State.User.ID,
			"",
			v.ID,
		)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	log.Println("Gracefully shutting down.")
}
