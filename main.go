package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"scholar-bot/apihandlers"

	"github.com/bwmarrin/discordgo"
)

var botSession *discordgo.Session

func init() {
	var botToken string
	var err error
	botToken = os.Getenv("scholar_bot")

	botSession, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot token: %v", err)
	}
}

func YearInputHelper(optionMap map[string]*discordgo.ApplicationCommandInteractionDataOption) string {
	if minYear, ok := optionMap["minyear"]; ok {
		return fmt.Sprint(minYear.IntValue())
	} else {
		return "2015"
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "gs",
			Description: "Get first study found on google scholar",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "google",
					Description: "What study should the bot google",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minyear",
					Description: "Minimum year for study (default 2015)",
					Required:    false,
				},
			},
		},
		{
			Name:        "gst10",
			Description: "Get first 10 studies found on google scholar",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "google",
					Description: "what should the bot google",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minyear",
					Description: "Minimum year for study (default 2015)",
					Required:    false,
				},
			},
		},
		{
			Name:        "pmc",
			Description: "Get first study found on PubMedCentral",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "google",
					Description: "What study should the bot google",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minyear",
					Description: "Minimum year for study (default 2015)",
					Required:    false,
				},
			},
		},
		{
			Name:        "pmct10",
			Description: "Get top 10 studies found on PubMedCentral",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "google",
					Description: "What study should the bot google",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "minyear",
					Description: "Minimum year for study (default 2015)",
					Required:    false,
				},
			},
		},
	}

	commandHandlers = map[string]func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate){
		"gs": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			options := botInteraction.ApplicationCommandData().Options
			optionMap := make(
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
				len(options),
			)
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if query, ok := optionMap["google"]; ok {
				var studyEmbed *apihandlers.StudyStruct
				studyEmbed, ok := apihandlers.QueryFirstGs(query.StringValue(),YearInputHelper(optionMap))
				if ok {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds: []*discordgo.MessageEmbed{
									{
										Title:       studyEmbed.Title,
										Description: studyEmbed.Abstract,
										URL:         studyEmbed.Url,
										Author: &discordgo.MessageEmbedAuthor{
											Name: studyEmbed.Authors,
										},
									},
								},
							},
						})
				} else {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error happened when retrieving the studies from google scholar",
								Flags:   1 << 6,
							},
						})
				}
			} else {
				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error happened when retrieving the query",
							Flags:   1 << 6,
						},
					})
			}

		},
		"gst10": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			options := botInteraction.ApplicationCommandData().Options
			optionMap := make(
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
				len(options),
			)
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if query, ok := optionMap["google"]; ok {
				var studySlice *[]apihandlers.StudyStruct
				studySlice, ok := apihandlers.QueryTopTenGs(query.StringValue(), YearInputHelper(optionMap))
				var studyTextList string
				for _, studyStruct := range *studySlice {
					studyTextList = studyTextList + fmt.Sprintf(
						"- [%s](<%s>)\n",
						studyStruct.Title,
						studyStruct.Url,
					)
				}
				if ok {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: studyTextList,
								Embeds:  nil,
							},
						})
				} else {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error happened when retrieving the studies from google scholar",
								Flags:   1 << 6,
							},
						})
				}
			} else {
				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error happened when retrieving the query",
							Flags:   1 << 6,
						},
					})
			}

		},
		"pmc": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			options := botInteraction.ApplicationCommandData().Options
			optionMap := make(
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
				len(options),
			)
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if query, ok := optionMap["google"]; ok {
				var studyEmbed *apihandlers.StudyStruct
				studyEmbed, ok := apihandlers.QueryFirstPMC(query.StringValue(), YearInputHelper(optionMap))
				if ok {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds: []*discordgo.MessageEmbed{
									{
										Title:       studyEmbed.Title,
										Description: studyEmbed.Abstract,
										URL:         studyEmbed.Url,
										Author: &discordgo.MessageEmbedAuthor{
											Name: studyEmbed.Authors,
										},
									},
								},
							},
						})
				} else {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error happened when retrieving the studies from PubMed",
								Flags:   1 << 6,
							},
						})
				}
			} else {
				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error happened when retrieving the query",
							Flags:   1 << 6,
						},
					})
			}

		},
		"pmct10": func(botSession *discordgo.Session, botInteraction *discordgo.InteractionCreate) {
			options := botInteraction.ApplicationCommandData().Options
			optionMap := make(
				map[string]*discordgo.ApplicationCommandInteractionDataOption,
				len(options),
			)
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			if query, ok := optionMap["google"]; ok {
				var studySlice *[]apihandlers.StudyStruct
				studySlice, ok := apihandlers.QueryTopTenPMC(query.StringValue(), YearInputHelper(optionMap))
				var studyTextList string
				for _, studyStruct := range *studySlice {
					studyTextList = studyTextList + fmt.Sprintf(
						"- [%s](<%s>)\n",
						studyStruct.Title,
						studyStruct.Url,
					)
				}
				if ok {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: studyTextList,
								Embeds:  nil,
							},
						})
				} else {
					botSession.InteractionRespond(
						botInteraction.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Content: "An error happened when retrieving the studies from google scholar",
								Flags:   1 << 6,
							},
						})
				}
			} else {
				botSession.InteractionRespond(
					botInteraction.Interaction,
					&discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: "An error happened when retrieving the query",
							Flags:   1 << 6,
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
