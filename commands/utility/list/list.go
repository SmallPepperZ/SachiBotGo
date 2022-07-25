package list

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/smallpepperz/sachibotgo/api"
)

type Command struct {
	api.Command
}

func (*Command) Load(ds *discordgo.Session) {
	appCmd := &discordgo.ApplicationCommand{
		Name:        "list",
		Description: "Lists the bot's registered commands",
		Type:        discordgo.ChatApplicationCommand,
	}
	api.CreateGlobalCommand(ds, appCmd, runCommand)
}

func runCommand(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	commandstrings := []string{}
	for cmdName := range api.GetLoadedCommands() {
		commandstrings = append(commandstrings, "<:Success:865674863330328626> "+cmdName)
	}
	for cmdName := range api.GetAvailableCommands() {
		if api.GetLoadedCommands()[cmdName] == nil {
			commandstrings = append(commandstrings, "<:Failure:865674863031877663> "+cmdName)
		}
	}

	// add image to embed
	embed := &discordgo.MessageEmbed{
		Title:       "Commands",
		Description: strings.Join(commandstrings, "\n"),
	}

	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  uint64(discordgo.MessageFlagsEphemeral),
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		fmt.Printf("Failed to send message\n%s", err)
	}
}

func (Command) Name() string {
	return "utility/list"
}
